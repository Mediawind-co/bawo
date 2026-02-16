package language

import (
	"context"
	"errors"

	"encore.dev/beta/errs"
	"github.com/google/uuid"
)

// ========== Public Endpoints ==========

// ListLanguages returns all active languages available for learning.
//
//encore:api auth method=GET path=/languages
func ListLanguages(ctx context.Context) (*LanguagesResponse, error) {
	languages, err := List(ctx, false) // only active languages
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list languages"}
	}
	return &LanguagesResponse{Languages: languages}, nil
}

// GetLanguage returns a specific language by ID.
//
//encore:api auth method=GET path=/languages/:id
func GetLanguage(ctx context.Context, id string) (*LanguageResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	lang, err := FindByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrLanguageNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "language not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch language"}
	}

	return &LanguageResponse{Language: lang}, nil
}

// ========== Admin Endpoints ==========

// ListAllLanguages returns all languages including inactive ones (admin only).
//
//encore:api auth method=GET path=/admin/languages tag:admin
func ListAllLanguages(ctx context.Context) (*LanguagesResponse, error) {
	languages, err := List(ctx, true) // include inactive
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list languages"}
	}
	return &LanguagesResponse{Languages: languages}, nil
}

// CreateLanguage creates a new language (admin only).
//
//encore:api auth method=POST path=/admin/languages tag:admin
func CreateLanguage(ctx context.Context, params *CreateLanguageParams) (*LanguageResponse, error) {
	if params.Name == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "name is required"}
	}
	if params.Code == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "code is required"}
	}

	lang, err := Create(ctx, *params)
	if err != nil {
		if errors.Is(err, ErrLanguageExists) {
			return nil, &errs.Error{Code: errs.AlreadyExists, Message: "language already exists"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to create language"}
	}

	return &LanguageResponse{Language: lang}, nil
}

// UpdateLanguage updates an existing language (admin only).
//
//encore:api auth method=PUT path=/admin/languages/:id tag:admin
func UpdateLanguage(ctx context.Context, id string, params *UpdateLanguageParams) (*LanguageResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	lang, err := Update(ctx, uid, *params)
	if err != nil {
		if errors.Is(err, ErrLanguageNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "language not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update language"}
	}

	return &LanguageResponse{Language: lang}, nil
}

// DeleteLanguageResponse is returned after deleting a language.
type DeleteLanguageResponse struct {
	Success bool `json:"success"`
}

// DeleteLanguage deletes a language (admin only).
//
//encore:api auth method=DELETE path=/admin/languages/:id tag:admin
func DeleteLanguage(ctx context.Context, id string) (*DeleteLanguageResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	err = Delete(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrLanguageNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "language not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete language"}
	}

	return &DeleteLanguageResponse{Success: true}, nil
}
