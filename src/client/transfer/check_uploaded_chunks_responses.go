// Code generated by go-swagger; DO NOT EDIT.

package transfer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/xy3/synche/src/models"
)

// CheckUploadedChunksReader is a Reader for the CheckUploadedChunks structure.
type CheckUploadedChunksReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CheckUploadedChunksReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCheckUploadedChunksOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewCheckUploadedChunksUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewCheckUploadedChunksDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewCheckUploadedChunksOK creates a CheckUploadedChunksOK with default headers values
func NewCheckUploadedChunksOK() *CheckUploadedChunksOK {
	return &CheckUploadedChunksOK{}
}

/* CheckUploadedChunksOK describes a response with status code 200, with default header values.

OK
*/
type CheckUploadedChunksOK struct {
	Payload *models.ExistingChunks
}

func (o *CheckUploadedChunksOK) Error() string {
	return fmt.Sprintf("[POST /upload/check][%d] checkUploadedChunksOK  %+v", 200, o.Payload)
}
func (o *CheckUploadedChunksOK) GetPayload() *models.ExistingChunks {
	return o.Payload
}

func (o *CheckUploadedChunksOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.ExistingChunks)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewCheckUploadedChunksUnauthorized creates a CheckUploadedChunksUnauthorized with default headers values
func NewCheckUploadedChunksUnauthorized() *CheckUploadedChunksUnauthorized {
	return &CheckUploadedChunksUnauthorized{}
}

/* CheckUploadedChunksUnauthorized describes a response with status code 401, with default header values.

Unauthorised
*/
type CheckUploadedChunksUnauthorized struct {
}

func (o *CheckUploadedChunksUnauthorized) Error() string {
	return fmt.Sprintf("[POST /upload/check][%d] checkUploadedChunksUnauthorized ", 401)
}

func (o *CheckUploadedChunksUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewCheckUploadedChunksDefault creates a CheckUploadedChunksDefault with default headers values
func NewCheckUploadedChunksDefault(code int) *CheckUploadedChunksDefault {
	return &CheckUploadedChunksDefault{
		_statusCode: code,
	}
}

/* CheckUploadedChunksDefault describes a response with status code -1, with default header values.

Error
*/
type CheckUploadedChunksDefault struct {
	_statusCode int

	Payload models.Error
}

// Code gets the status code for the check uploaded chunks default response
func (o *CheckUploadedChunksDefault) Code() int {
	return o._statusCode
}

func (o *CheckUploadedChunksDefault) Error() string {
	return fmt.Sprintf("[POST /upload/check][%d] checkUploadedChunks default  %+v", o._statusCode, o.Payload)
}
func (o *CheckUploadedChunksDefault) GetPayload() models.Error {
	return o.Payload
}

func (o *CheckUploadedChunksDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}