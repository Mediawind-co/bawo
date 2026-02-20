package lesson

import (
	"time"

	"encore.dev/pubsub"
)

// LessonCompletedEvent is published when a user completes a lesson.
type LessonCompletedEvent struct {
	UserID    string    `json:"user_id"`
	LessonID  string    `json:"lesson_id"`
	SessionID string    `json:"session_id"`
	XPEarned  int       `json:"xp_earned"`
	Correct   int       `json:"correct"`
	Incorrect int       `json:"incorrect"`
	Timestamp time.Time `json:"timestamp"`
}

// LessonCompleted is the topic for lesson completion events.
var LessonCompleted = pubsub.NewTopic[*LessonCompletedEvent]("lesson-completed", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// PublishLessonCompleted publishes a lesson completed event.
func PublishLessonCompleted(session *LessonSession) error {
	_, err := LessonCompleted.Publish(nil, &LessonCompletedEvent{
		UserID:    session.UserID.String(),
		LessonID:  session.LessonID.String(),
		SessionID: session.ID.String(),
		XPEarned:  session.XPEarned,
		Correct:   session.CorrectCount,
		Incorrect: session.IncorrectCount,
		Timestamp: time.Now(),
	})
	return err
}

// MistakeRecordedEvent is published when a user answers a question incorrectly.
type MistakeRecordedEvent struct {
	UserID        string    `json:"user_id"`
	QuestionID    string    `json:"question_id"`
	LanguageID    string    `json:"language_id"`
	UserAnswer    string    `json:"user_answer"`
	CorrectAnswer string    `json:"correct_answer"`
	Timestamp     time.Time `json:"timestamp"`
}

// MistakeRecorded is the topic for mistake events.
var MistakeRecorded = pubsub.NewTopic[*MistakeRecordedEvent]("mistake-recorded", pubsub.TopicConfig{
	DeliveryGuarantee: pubsub.AtLeastOnce,
})

// PublishMistakeRecorded publishes a mistake recorded event.
func PublishMistakeRecorded(userID, questionID, languageID, userAnswer, correctAnswer string) error {
	_, err := MistakeRecorded.Publish(nil, &MistakeRecordedEvent{
		UserID:        userID,
		QuestionID:    questionID,
		LanguageID:    languageID,
		UserAnswer:    userAnswer,
		CorrectAnswer: correctAnswer,
		Timestamp:     time.Now(),
	})
	return err
}
