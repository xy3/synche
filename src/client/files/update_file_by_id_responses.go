// Code generated by go-swagger; DO NOT EDIT.

package files

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/xy3/synche/src/client/models"
)

// UpdateFileByIDReader is a Reader for the UpdateFileByID structure.
type UpdateFileByIDReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UpdateFileByIDReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUpdateFileByIDOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUpdateFileByIDUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUpdateFileByIDNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewUpdateFileByIDDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUpdateFileByIDOK creates a UpdateFileByIDOK with default headers values
func NewUpdateFileByIDOK() *UpdateFileByIDOK {
	return &UpdateFileByIDOK{}
}

/* UpdateFileByIDOK describes a response with status code 200, with default header values.

OK
*/
type UpdateFileByIDOK struct {
	Payload *models.File
}

func (o *UpdateFileByIDOK) Error() string {
	return fmt.Sprintf("[PATCH /files/{fileID}][%d] updateFileByIdOK  %+v", 200, o.Payload)
}
func (o *UpdateFileByIDOK) GetPayload() *models.File {
	return o.Payload
}

func (o *UpdateFileByIDOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.File)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUpdateFileByIDUnauthorized creates a UpdateFileByIDUnauthorized with default headers values
func NewUpdateFileByIDUnauthorized() *UpdateFileByIDUnauthorized {
	return &UpdateFileByIDUnauthorized{}
}

/* UpdateFileByIDUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type UpdateFileByIDUnauthorized struct {
}

func (o *UpdateFileByIDUnauthorized) Error() string {
	return fmt.Sprintf("[PATCH /files/{fileID}][%d] updateFileByIdUnauthorized ", 401)
}

func (o *UpdateFileByIDUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateFileByIDNotFound creates a UpdateFileByIDNotFound with default headers values
func NewUpdateFileByIDNotFound() *UpdateFileByIDNotFound {
	return &UpdateFileByIDNotFound{}
}

/* UpdateFileByIDNotFound describes a response with status code 404, with default header values.

File not found
*/
type UpdateFileByIDNotFound struct {
}

func (o *UpdateFileByIDNotFound) Error() string {
	return fmt.Sprintf("[PATCH /files/{fileID}][%d] updateFileByIdNotFound ", 404)
}

func (o *UpdateFileByIDNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUpdateFileByIDDefault creates a UpdateFileByIDDefault with default headers values
func NewUpdateFileByIDDefault(code int) *UpdateFileByIDDefault {
	return &UpdateFileByIDDefault{
		_statusCode: code,
	}
}

/* UpdateFileByIDDefault describes a response with status code -1, with default header values.

Error
*/
type UpdateFileByIDDefault struct {
	_statusCode int

	Payload models.Error
}

// Code gets the status code for the update file by ID default response
func (o *UpdateFileByIDDefault) Code() int {
	return o._statusCode
}

func (o *UpdateFileByIDDefault) Error() string {
	return fmt.Sprintf("[PATCH /files/{fileID}][%d] updateFileByID default  %+v", o._statusCode, o.Payload)
}
func (o *UpdateFileByIDDefault) GetPayload() models.Error {
	return o.Payload
}

func (o *UpdateFileByIDDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
