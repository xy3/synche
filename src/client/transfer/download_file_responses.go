// Code generated by go-swagger; DO NOT EDIT.

package transfer

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/xy3/synche/src/models"
)

// DownloadFileReader is a Reader for the DownloadFile structure.
type DownloadFileReader struct {
	formats strfmt.Registry
	writer  io.Writer
}

// ReadResponse reads a server response into the received o.
func (o *DownloadFileReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewDownloadFileOK(o.writer)
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewDownloadFileUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewDownloadFileForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewDownloadFileNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		result := NewDownloadFileDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewDownloadFileOK creates a DownloadFileOK with default headers values
func NewDownloadFileOK(writer io.Writer) *DownloadFileOK {
	return &DownloadFileOK{

		Payload: writer,
	}
}

/* DownloadFileOK describes a response with status code 200, with default header values.

OK
*/
type DownloadFileOK struct {
	ContentDisposition string
	ContentLength      uint64

	Payload io.Writer
}

func (o *DownloadFileOK) Error() string {
	return fmt.Sprintf("[GET /download/{fileID}][%d] downloadFileOK  %+v", 200, o.Payload)
}
func (o *DownloadFileOK) GetPayload() io.Writer {
	return o.Payload
}

func (o *DownloadFileOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// hydrates response header Content-Disposition
	hdrContentDisposition := response.GetHeader("Content-Disposition")

	if hdrContentDisposition != "" {
		o.ContentDisposition = hdrContentDisposition
	}

	// hydrates response header Content-Length
	hdrContentLength := response.GetHeader("Content-Length")

	if hdrContentLength != "" {
		valcontentLength, err := swag.ConvertUint64(hdrContentLength)
		if err != nil {
			return errors.InvalidType("Content-Length", "header", "uint64", hdrContentLength)
		}
		o.ContentLength = valcontentLength
	}

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewDownloadFileUnauthorized creates a DownloadFileUnauthorized with default headers values
func NewDownloadFileUnauthorized() *DownloadFileUnauthorized {
	return &DownloadFileUnauthorized{}
}

/* DownloadFileUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type DownloadFileUnauthorized struct {
}

func (o *DownloadFileUnauthorized) Error() string {
	return fmt.Sprintf("[GET /download/{fileID}][%d] downloadFileUnauthorized ", 401)
}

func (o *DownloadFileUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDownloadFileForbidden creates a DownloadFileForbidden with default headers values
func NewDownloadFileForbidden() *DownloadFileForbidden {
	return &DownloadFileForbidden{}
}

/* DownloadFileForbidden describes a response with status code 403, with default header values.

Forbidden
*/
type DownloadFileForbidden struct {
}

func (o *DownloadFileForbidden) Error() string {
	return fmt.Sprintf("[GET /download/{fileID}][%d] downloadFileForbidden ", 403)
}

func (o *DownloadFileForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDownloadFileNotFound creates a DownloadFileNotFound with default headers values
func NewDownloadFileNotFound() *DownloadFileNotFound {
	return &DownloadFileNotFound{}
}

/* DownloadFileNotFound describes a response with status code 404, with default header values.

File not found
*/
type DownloadFileNotFound struct {
}

func (o *DownloadFileNotFound) Error() string {
	return fmt.Sprintf("[GET /download/{fileID}][%d] downloadFileNotFound ", 404)
}

func (o *DownloadFileNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewDownloadFileDefault creates a DownloadFileDefault with default headers values
func NewDownloadFileDefault(code int) *DownloadFileDefault {
	return &DownloadFileDefault{
		_statusCode: code,
	}
}

/* DownloadFileDefault describes a response with status code -1, with default header values.

error
*/
type DownloadFileDefault struct {
	_statusCode int

	Payload models.Error
}

// Code gets the status code for the download file default response
func (o *DownloadFileDefault) Code() int {
	return o._statusCode
}

func (o *DownloadFileDefault) Error() string {
	return fmt.Sprintf("[GET /download/{fileID}][%d] downloadFile default  %+v", o._statusCode, o.Payload)
}
func (o *DownloadFileDefault) GetPayload() models.Error {
	return o.Payload
}

func (o *DownloadFileDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
