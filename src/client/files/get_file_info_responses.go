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

// GetFileInfoReader is a Reader for the GetFileInfo structure.
type GetFileInfoReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetFileInfoReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetFileInfoOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewGetFileInfoUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewGetFileInfoNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewGetFileInfoOK creates a GetFileInfoOK with default headers values
func NewGetFileInfoOK() *GetFileInfoOK {
	return &GetFileInfoOK{}
}

/* GetFileInfoOK describes a response with status code 200, with default header values.

OK
*/
type GetFileInfoOK struct {
	Payload *models.File
}

func (o *GetFileInfoOK) Error() string {
	return fmt.Sprintf("[GET /files/{fileID}][%d] getFileInfoOK  %+v", 200, o.Payload)
}
func (o *GetFileInfoOK) GetPayload() *models.File {
	return o.Payload
}

func (o *GetFileInfoOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.File)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetFileInfoUnauthorized creates a GetFileInfoUnauthorized with default headers values
func NewGetFileInfoUnauthorized() *GetFileInfoUnauthorized {
	return &GetFileInfoUnauthorized{}
}

/* GetFileInfoUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type GetFileInfoUnauthorized struct {
}

func (o *GetFileInfoUnauthorized) Error() string {
	return fmt.Sprintf("[GET /files/{fileID}][%d] getFileInfoUnauthorized ", 401)
}

func (o *GetFileInfoUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewGetFileInfoNotFound creates a GetFileInfoNotFound with default headers values
func NewGetFileInfoNotFound() *GetFileInfoNotFound {
	return &GetFileInfoNotFound{}
}

/* GetFileInfoNotFound describes a response with status code 404, with default header values.

File not found
*/
type GetFileInfoNotFound struct {
}

func (o *GetFileInfoNotFound) Error() string {
	return fmt.Sprintf("[GET /files/{fileID}][%d] getFileInfoNotFound ", 404)
}

func (o *GetFileInfoNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}