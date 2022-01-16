// Code generated by go-swagger; DO NOT EDIT.

package tokens

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/xy3/synche/src/models"
)

// RefreshTokenReader is a Reader for the RefreshToken structure.
type RefreshTokenReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *RefreshTokenReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewRefreshTokenOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewRefreshTokenDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewRefreshTokenOK creates a RefreshTokenOK with default headers values
func NewRefreshTokenOK() *RefreshTokenOK {
	return &RefreshTokenOK{}
}

/* RefreshTokenOK describes a response with status code 200, with default header values.

Token refreshed successfully
*/
type RefreshTokenOK struct {
	Payload *models.AccessToken
}

func (o *RefreshTokenOK) Error() string {
	return fmt.Sprintf("[POST /token/refresh][%d] refreshTokenOK  %+v", 200, o.Payload)
}
func (o *RefreshTokenOK) GetPayload() *models.AccessToken {
	return o.Payload
}

func (o *RefreshTokenOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.AccessToken)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewRefreshTokenDefault creates a RefreshTokenDefault with default headers values
func NewRefreshTokenDefault(code int) *RefreshTokenDefault {
	return &RefreshTokenDefault{
		_statusCode: code,
	}
}

/* RefreshTokenDefault describes a response with status code -1, with default header values.

Error
*/
type RefreshTokenDefault struct {
	_statusCode int

	Payload models.Error
}

// Code gets the status code for the refresh token default response
func (o *RefreshTokenDefault) Code() int {
	return o._statusCode
}

func (o *RefreshTokenDefault) Error() string {
	return fmt.Sprintf("[POST /token/refresh][%d] refreshToken default  %+v", o._statusCode, o.Payload)
}
func (o *RefreshTokenDefault) GetPayload() models.Error {
	return o.Payload
}

func (o *RefreshTokenDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}