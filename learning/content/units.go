package content

import (
	"context"
	"errors"

	"encore.dev/beta/errs"
	"github.com/google/uuid"
)

// ========== Public Endpoints ==========

// ListUnits returns all units for a language.
//
//encore:api auth method=GET path=/languages/:languageID/units
func ListUnits(ctx context.Context, languageID string) (*UnitsResponse, error) {
	langID, err := uuid.Parse(languageID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	units, err := ListUnitsByLanguage(ctx, langID)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to list units"}
	}

	return &UnitsResponse{Units: units}, nil
}

// GetUnit returns a specific unit by ID.
//
//encore:api auth method=GET path=/units/:id
func GetUnit(ctx context.Context, id string) (*UnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid unit ID"}
	}

	unit, err := FindUnitByID(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrUnitNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "unit not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to fetch unit"}
	}

	return &UnitResponse{Unit: unit}, nil
}

// ========== Admin Endpoints ==========

// CreateUnitAdmin creates a new unit (admin only).
//
//encore:api auth method=POST path=/admin/units tag:admin
func CreateUnitAdmin(ctx context.Context, params *CreateUnitParams) (*UnitResponse, error) {
	if params.LanguageID == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "language_id is required"}
	}
	if params.Title == "" {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "title is required"}
	}

	langID, err := uuid.Parse(params.LanguageID)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid language ID"}
	}

	unit, err := CreateUnit(ctx, langID, *params)
	if err != nil {
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to create unit"}
	}

	return &UnitResponse{Unit: unit}, nil
}

// UpdateUnitAdmin updates an existing unit (admin only).
//
//encore:api auth method=PUT path=/admin/units/:id tag:admin
func UpdateUnitAdmin(ctx context.Context, id string, params *UpdateUnitParams) (*UnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid unit ID"}
	}

	unit, err := UpdateUnit(ctx, uid, *params)
	if err != nil {
		if errors.Is(err, ErrUnitNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "unit not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to update unit"}
	}

	return &UnitResponse{Unit: unit}, nil
}

// DeleteUnitResponse is returned after deleting a unit.
type DeleteUnitResponse struct {
	Success bool `json:"success"`
}

// DeleteUnitAdmin deletes a unit (admin only).
//
//encore:api auth method=DELETE path=/admin/units/:id tag:admin
func DeleteUnitAdmin(ctx context.Context, id string) (*DeleteUnitResponse, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, &errs.Error{Code: errs.InvalidArgument, Message: "invalid unit ID"}
	}

	err = DeleteUnit(ctx, uid)
	if err != nil {
		if errors.Is(err, ErrUnitNotFound) {
			return nil, &errs.Error{Code: errs.NotFound, Message: "unit not found"}
		}
		return nil, &errs.Error{Code: errs.Internal, Message: "failed to delete unit"}
	}

	return &DeleteUnitResponse{Success: true}, nil
}
