// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// User User profile information
//
// swagger:model User
type User struct {

	// email
	// Required: true
	Email *string `json:"email"`

	// email verified
	EmailVerified bool `json:"emailVerified,omitempty"`

	// id
	// Required: true
	ID *uint64 `json:"id"`

	// name
	// Required: true
	Name *string `json:"name"`

	// picture
	// Required: true
	Picture *string `json:"picture"`

	// role
	// Required: true
	// Enum: [member admin]
	Role *string `json:"role"`
}

// Validate validates this user
func (m *User) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateEmail(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateID(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateName(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validatePicture(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateRole(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *User) validateEmail(formats strfmt.Registry) error {

	if err := validate.Required("email", "body", m.Email); err != nil {
		return err
	}

	return nil
}

func (m *User) validateID(formats strfmt.Registry) error {

	if err := validate.Required("id", "body", m.ID); err != nil {
		return err
	}

	return nil
}

func (m *User) validateName(formats strfmt.Registry) error {

	if err := validate.Required("name", "body", m.Name); err != nil {
		return err
	}

	return nil
}

func (m *User) validatePicture(formats strfmt.Registry) error {

	if err := validate.Required("picture", "body", m.Picture); err != nil {
		return err
	}

	return nil
}

var userTypeRolePropEnum []interface{}

func init() {
	var res []string
	if err := json.Unmarshal([]byte(`["member","admin"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		userTypeRolePropEnum = append(userTypeRolePropEnum, v)
	}
}

const (

	// UserRoleMember captures enum value "member"
	UserRoleMember string = "member"

	// UserRoleAdmin captures enum value "admin"
	UserRoleAdmin string = "admin"
)

// prop value enum
func (m *User) validateRoleEnum(path, location string, value string) error {
	if err := validate.EnumCase(path, location, value, userTypeRolePropEnum, true); err != nil {
		return err
	}
	return nil
}

func (m *User) validateRole(formats strfmt.Registry) error {

	if err := validate.Required("role", "body", m.Role); err != nil {
		return err
	}

	// value enum
	if err := m.validateRoleEnum("role", "body", *m.Role); err != nil {
		return err
	}

	return nil
}

// ContextValidate validates this user based on context it is used
func (m *User) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *User) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *User) UnmarshalBinary(b []byte) error {
	var res User
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}