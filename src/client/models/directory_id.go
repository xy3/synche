// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// DirectoryID the id of the directory to list
// Example: d290f1ee-6c54-4b01-90e6-d701748f0851
//
// swagger:model DirectoryId
type DirectoryID strfmt.UUID

// Validate validates this directory Id
func (m DirectoryID) Validate(formats strfmt.Registry) error {
	var res []error

	if err := validate.FormatOf("", "body", "uuid", strfmt.UUID(m).String(), formats); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this directory Id based on context it is used
func (m DirectoryID) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
