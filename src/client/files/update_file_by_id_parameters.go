// Code generated by go-swagger; DO NOT EDIT.

package files

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/xy3/synche/src/models"
)

// NewUpdateFileByIDParams creates a new UpdateFileByIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewUpdateFileByIDParams() *UpdateFileByIDParams {
	return &UpdateFileByIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewUpdateFileByIDParamsWithTimeout creates a new UpdateFileByIDParams object
// with the ability to set a timeout on a request.
func NewUpdateFileByIDParamsWithTimeout(timeout time.Duration) *UpdateFileByIDParams {
	return &UpdateFileByIDParams{
		timeout: timeout,
	}
}

// NewUpdateFileByIDParamsWithContext creates a new UpdateFileByIDParams object
// with the ability to set a context for a request.
func NewUpdateFileByIDParamsWithContext(ctx context.Context) *UpdateFileByIDParams {
	return &UpdateFileByIDParams{
		Context: ctx,
	}
}

// NewUpdateFileByIDParamsWithHTTPClient creates a new UpdateFileByIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewUpdateFileByIDParamsWithHTTPClient(client *http.Client) *UpdateFileByIDParams {
	return &UpdateFileByIDParams{
		HTTPClient: client,
	}
}

/* UpdateFileByIDParams contains all the parameters to send to the API endpoint
   for the update file by ID operation.

   Typically these are written to a http.Request.
*/
type UpdateFileByIDParams struct {

	// FileID.
	//
	// Format: uint
	FileID uint64

	// FileUpdate.
	FileUpdate *models.FileUpdate

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the update file by ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateFileByIDParams) WithDefaults() *UpdateFileByIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the update file by ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *UpdateFileByIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the update file by ID params
func (o *UpdateFileByIDParams) WithTimeout(timeout time.Duration) *UpdateFileByIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the update file by ID params
func (o *UpdateFileByIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the update file by ID params
func (o *UpdateFileByIDParams) WithContext(ctx context.Context) *UpdateFileByIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the update file by ID params
func (o *UpdateFileByIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the update file by ID params
func (o *UpdateFileByIDParams) WithHTTPClient(client *http.Client) *UpdateFileByIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the update file by ID params
func (o *UpdateFileByIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFileID adds the fileID to the update file by ID params
func (o *UpdateFileByIDParams) WithFileID(fileID uint64) *UpdateFileByIDParams {
	o.SetFileID(fileID)
	return o
}

// SetFileID adds the fileId to the update file by ID params
func (o *UpdateFileByIDParams) SetFileID(fileID uint64) {
	o.FileID = fileID
}

// WithFileUpdate adds the fileUpdate to the update file by ID params
func (o *UpdateFileByIDParams) WithFileUpdate(fileUpdate *models.FileUpdate) *UpdateFileByIDParams {
	o.SetFileUpdate(fileUpdate)
	return o
}

// SetFileUpdate adds the fileUpdate to the update file by ID params
func (o *UpdateFileByIDParams) SetFileUpdate(fileUpdate *models.FileUpdate) {
	o.FileUpdate = fileUpdate
}

// WriteToRequest writes these params to a swagger request
func (o *UpdateFileByIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param fileID
	if err := r.SetPathParam("fileID", swag.FormatUint64(o.FileID)); err != nil {
		return err
	}
	if o.FileUpdate != nil {
		if err := r.SetBodyParam(o.FileUpdate); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}