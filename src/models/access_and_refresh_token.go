// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// AccessAndRefreshToken Contains both the AccessToken and RefreshToken for a user.
//
// swagger:model AccessAndRefreshToken
type AccessAndRefreshToken struct {

	// access token
	AccessToken string `json:"accessToken,omitempty"`

	// the unix timestamp expiry of the access token
	AccessTokenExpiry int64 `json:"accessTokenExpiry,omitempty"`

	// refresh token
	RefreshToken string `json:"refreshToken,omitempty"`
}

// Validate validates this access and refresh token
func (m *AccessAndRefreshToken) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this access and refresh token based on context it is used
func (m *AccessAndRefreshToken) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *AccessAndRefreshToken) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *AccessAndRefreshToken) UnmarshalBinary(b []byte) error {
	var res AccessAndRefreshToken
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}