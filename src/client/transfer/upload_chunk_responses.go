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

// UploadChunkReader is a Reader for the UploadChunk structure.
type UploadChunkReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UploadChunkReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 201:
		result := NewUploadChunkCreated()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewUploadChunkDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewUploadChunkCreated creates a UploadChunkCreated with default headers values
func NewUploadChunkCreated() *UploadChunkCreated {
	return &UploadChunkCreated{}
}

/* UploadChunkCreated describes a response with status code 201, with default header values.

OK
*/
type UploadChunkCreated struct {
	Payload *models.FileChunk
}

func (o *UploadChunkCreated) Error() string {
	return fmt.Sprintf("[POST /upload/chunk][%d] uploadChunkCreated  %+v", 201, o.Payload)
}
func (o *UploadChunkCreated) GetPayload() *models.FileChunk {
	return o.Payload
}

func (o *UploadChunkCreated) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.FileChunk)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewUploadChunkDefault creates a UploadChunkDefault with default headers values
func NewUploadChunkDefault(code int) *UploadChunkDefault {
	return &UploadChunkDefault{
		_statusCode: code,
	}
}

/* UploadChunkDefault describes a response with status code -1, with default header values.

Error
*/
type UploadChunkDefault struct {
	_statusCode int

	Payload models.Error
}

// Code gets the status code for the upload chunk default response
func (o *UploadChunkDefault) Code() int {
	return o._statusCode
}

func (o *UploadChunkDefault) Error() string {
	return fmt.Sprintf("[POST /upload/chunk][%d] uploadChunk default  %+v", o._statusCode, o.Payload)
}
func (o *UploadChunkDefault) GetPayload() models.Error {
	return o.Payload
}

func (o *UploadChunkDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
