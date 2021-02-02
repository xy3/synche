// Code generated by go-swagger; DO NOT EDIT.

package testing

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

// CheckGetReader is a Reader for the CheckGet structure.
type CheckGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *CheckGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewCheckGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewCheckGetOK creates a CheckGetOK with default headers values
func NewCheckGetOK() *CheckGetOK {
	return &CheckGetOK{}
}

/* CheckGetOK describes a response with status code 200, with default header values.

OK
*/
type CheckGetOK struct {
	Payload *models.Message
}

func (o *CheckGetOK) Error() string {
	return fmt.Sprintf("[GET /check][%d] checkGetOK  %+v", 200, o.Payload)
}
func (o *CheckGetOK) GetPayload() *models.Message {
	return o.Payload
}

func (o *CheckGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Message)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
