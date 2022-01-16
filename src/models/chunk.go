// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Chunk A file chunk
//
// swagger:model Chunk
type Chunk struct {

	// hash
	Hash string `json:"Hash,omitempty"`

	// ID
	ID uint64 `json:"ID,omitempty"`

	// size
	Size int64 `json:"Size,omitempty"`
}

// Validate validates this chunk
func (m *Chunk) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this chunk based on context it is used
func (m *Chunk) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *Chunk) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Chunk) UnmarshalBinary(b []byte) error {
	var res Chunk
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
