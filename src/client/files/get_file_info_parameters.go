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
)

// NewGetFileInfoParams creates a new GetFileInfoParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetFileInfoParams() *GetFileInfoParams {
	return &GetFileInfoParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetFileInfoParamsWithTimeout creates a new GetFileInfoParams object
// with the ability to set a timeout on a request.
func NewGetFileInfoParamsWithTimeout(timeout time.Duration) *GetFileInfoParams {
	return &GetFileInfoParams{
		timeout: timeout,
	}
}

// NewGetFileInfoParamsWithContext creates a new GetFileInfoParams object
// with the ability to set a context for a request.
func NewGetFileInfoParamsWithContext(ctx context.Context) *GetFileInfoParams {
	return &GetFileInfoParams{
		Context: ctx,
	}
}

// NewGetFileInfoParamsWithHTTPClient creates a new GetFileInfoParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetFileInfoParamsWithHTTPClient(client *http.Client) *GetFileInfoParams {
	return &GetFileInfoParams{
		HTTPClient: client,
	}
}

/* GetFileInfoParams contains all the parameters to send to the API endpoint
   for the get file info operation.

   Typically these are written to a http.Request.
*/
type GetFileInfoParams struct {

	// FileID.
	//
	// Format: uint
	FileID uint64

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get file info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetFileInfoParams) WithDefaults() *GetFileInfoParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get file info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetFileInfoParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get file info params
func (o *GetFileInfoParams) WithTimeout(timeout time.Duration) *GetFileInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get file info params
func (o *GetFileInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get file info params
func (o *GetFileInfoParams) WithContext(ctx context.Context) *GetFileInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get file info params
func (o *GetFileInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get file info params
func (o *GetFileInfoParams) WithHTTPClient(client *http.Client) *GetFileInfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get file info params
func (o *GetFileInfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFileID adds the fileID to the get file info params
func (o *GetFileInfoParams) WithFileID(fileID uint64) *GetFileInfoParams {
	o.SetFileID(fileID)
	return o
}

// SetFileID adds the fileId to the get file info params
func (o *GetFileInfoParams) SetFileID(fileID uint64) {
	o.FileID = fileID
}

// WriteToRequest writes these params to a swagger request
func (o *GetFileInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param fileID
	if err := r.SetPathParam("fileID", swag.FormatUint64(o.FileID)); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
