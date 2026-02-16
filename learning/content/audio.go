package content

import (
	"context"
	"errors"
	"fmt"
	"time"

	"encore.dev/beta/errs"
	"encore.dev/storage/objects"
	"github.com/google/uuid"
)

// AudioBucket stores pre-recorded and TTS-generated audio files.
var AudioBucket = objects.NewBucket("lesson-audio", objects.BucketConfig{
	Versioned: false,
})

// AudioUploadRequest contains the audio file data.
type AudioUploadRequest struct {
	Data        []byte `json:"data"`         // Base64-decoded audio bytes
	ContentType string `json:"content_type"` // e.g., "audio/mpeg", "audio/wav"
}

// AudioUploadResponse contains the result of an audio upload.
type AudioUploadResponse struct {
	AudioKey string `json:"audio_key"`
}

// UploadQuestionAudio uploads audio for a question (admin only).
//
//encore:api auth method=POST path=/admin/questions/:questionID/audio tag:admin
func UploadQuestionAudio(ctx context.Context, questionID string, req *AudioUploadRequest) (*AudioUploadResponse, error) {
	uid, err := uuid.Parse(questionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	// Verify question exists
	_, err = FindQuestionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch question"}
	}

	if len(req.Data) == 0 {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "audio data is required"}
	}

	// Determine file extension from content type
	ext := ".mp3"
	switch req.ContentType {
	case "audio/wav":
		ext = ".wav"
	case "audio/ogg":
		ext = ".ogg"
	case "audio/mpeg", "audio/mp3":
		ext = ".mp3"
	}

	// Generate storage key
	key := fmt.Sprintf("questions/%s/prompt%s", questionID, ext)

	// Upload to bucket
	writer := AudioBucket.Upload(ctx, key)

	_, err = writer.Write(req.Data)
	if err != nil {
		writer.Abort(err)
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to write audio data"}
	}

	err = writer.Close()
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to complete upload"}
	}

	// Update question with audio key
	err = UpdateQuestionAudioKey(ctx, uid, key)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update question"}
	}

	return &AudioUploadResponse{AudioKey: key}, nil
}

// AudioURLResponse contains a URL for audio playback.
type AudioURLResponse struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// GetQuestionAudioURL returns a signed URL for audio playback.
//
//encore:api auth method=GET path=/questions/:questionID/audio
func GetQuestionAudioURL(ctx context.Context, questionID string) (*AudioURLResponse, error) {
	uid, err := uuid.Parse(questionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	question, err := FindQuestionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch question"}
	}

	if question.PromptAudioKey == "" {
		return nil, &errs.Error{Code: errs.NotFound, Message: "no audio available for this question"}
	}

	// Generate signed URL (valid for 1 hour)
	ttl := 1 * time.Hour
	expiry := time.Now().Add(ttl)

	signedURL, err := AudioBucket.SignedDownloadURL(ctx, question.PromptAudioKey, objects.WithTTL(ttl))
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to generate audio URL"}
	}

	return &AudioURLResponse{
		URL:       signedURL.URL,
		ExpiresAt: expiry,
	}, nil
}

// DeleteQuestionAudio deletes audio for a question (admin only).
//
//encore:api auth method=DELETE path=/admin/questions/:questionID/audio tag:admin
func DeleteQuestionAudio(ctx context.Context, questionID string) (*DeleteQuestionResponse, error) {
	uid, err := uuid.Parse(questionID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid question ID"}
	}

	question, err := FindQuestionByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrQuestionNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "question not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch question"}
	}

	if question.PromptAudioKey == "" {
		return nil, &errs.Error{Code: errs.NotFound, Message: "no audio to delete"}
	}

	// Delete from bucket
	err = AudioBucket.Remove(ctx, question.PromptAudioKey)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete audio"}
	}

	// Clear audio key from question
	err = UpdateQuestionAudioKey(ctx, uid, "")
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update question"}
	}

	return &DeleteQuestionResponse{Success: true}, nil
}
