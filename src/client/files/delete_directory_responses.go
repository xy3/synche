// Code generated by go-swagger; DO NOT EDIT.

package files

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/xy3/synche/src/models"
)

// DeleteDirectoryReader is a Reader for the DeleteDirectory structure.
type DeleteDirectoryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *DeleteDirectoryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDeleteDirectoryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewDeleteDirectoryUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 501:
		result := NewDeleteDirectoryNotImplemented()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDeleteDirectoryDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDeleteDirectoryOK creates a DeleteDirectoryOK with default headers values
func NewDeleteDirectoryOK() *DeleteDirectoryOK {
	return &DeleteDirectoryOK{}
}

/* DeleteDirectoryOK describes a response with status code 200, with default header values.

OK
*/
type DeleteDirectoryOK struct {
}

func (o *DeleteDirectoryOK) Error() string {
	return fmt.Sprintf("[DELETE /directory/{id}][%d] deleteDirectoryOK ", 200)
}

func (o *DeleteDirectoryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDirectoryUnauthorized creates a DeleteDirectoryUnauthorized with default headers values
func NewDeleteDirectoryUnauthorized() *DeleteDirectoryUnauthorized {
	return &DeleteDirectoryUnauthorized{}
}

/* DeleteDirectoryUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type DeleteDirectoryUnauthorized struct {
}

func (o *DeleteDirectoryUnauthorized) Error() string {
	return fmt.Sprintf("[DELETE /directory/{id}][%d] deleteDirectoryUnauthorized ", 401)
}

func (o *DeleteDirectoryUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDirectoryNotImplemented creates a DeleteDirectoryNotImplemented with default headers values
func NewDeleteDirectoryNotImplemented() *DeleteDirectoryNotImplemented {
	return &DeleteDirectoryNotImplemented{}
}

/* DeleteDirectoryNotImplemented describes a response with status code 501, with default header values.

Not implemented
*/
type DeleteDirectoryNotImplemented struct {
}

func (o *DeleteDirectoryNotImplemented) Error() string {
	return fmt.Sprintf("[DELETE /directory/{id}][%d] deleteDirectoryNotImplemented ", 501)
}

func (o *DeleteDirectoryNotImplemented) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDeleteDirectoryDefault creates a DeleteDirectoryDefault with default headers values
func NewDeleteDirectoryDefault(code int) *DeleteDirectoryDefault {
	return &DeleteDirectoryDefault{
		_statusCode: code,
	}
}

/* DeleteDirectoryDefault describes a response with status code -1, with default header values.

Error
*/
type DeleteDirectoryDefault struct {
	_statusCode int

	Payload models.Error
}

// Code gets the status code for the delete directory default response
func (o *DeleteDirectoryDefault) Code() int {
	return o._statusCode
}

func (o *DeleteDirectoryDefault) Error() string {
	return fmt.Sprintf("[DELETE /directory/{id}][%d] deleteDirectory default  %+v", o._statusCode, o.Payload)
}
func (o *DeleteDirectoryDefault) GetPayload() models.Error {
	return o.Payload
}

func (o *DeleteDirectoryDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
