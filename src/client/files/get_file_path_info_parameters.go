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
)

// NewGetFilePathInfoParams creates a new GetFilePathInfoParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetFilePathInfoParams() *GetFilePathInfoParams {
	return &GetFilePathInfoParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetFilePathInfoParamsWithTimeout creates a new GetFilePathInfoParams object
// with the ability to set a timeout on a request.
func NewGetFilePathInfoParamsWithTimeout(timeout time.Duration) *GetFilePathInfoParams {
	return &GetFilePathInfoParams{
		timeout: timeout,
	}
}

// NewGetFilePathInfoParamsWithContext creates a new GetFilePathInfoParams object
// with the ability to set a context for a request.
func NewGetFilePathInfoParamsWithContext(ctx context.Context) *GetFilePathInfoParams {
	return &GetFilePathInfoParams{
		Context: ctx,
	}
}

// NewGetFilePathInfoParamsWithHTTPClient creates a new GetFilePathInfoParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetFilePathInfoParamsWithHTTPClient(client *http.Client) *GetFilePathInfoParams {
	return &GetFilePathInfoParams{
		HTTPClient: client,
	}
}

/* GetFilePathInfoParams contains all the parameters to send to the API endpoint
   for the get file path info operation.

   Typically these are written to a http.Request.
*/
type GetFilePathInfoParams struct {

	// FilePath.
	FilePath string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get file path info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetFilePathInfoParams) WithDefaults() *GetFilePathInfoParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get file path info params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetFilePathInfoParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get file path info params
func (o *GetFilePathInfoParams) WithTimeout(timeout time.Duration) *GetFilePathInfoParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get file path info params
func (o *GetFilePathInfoParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get file path info params
func (o *GetFilePathInfoParams) WithContext(ctx context.Context) *GetFilePathInfoParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get file path info params
func (o *GetFilePathInfoParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get file path info params
func (o *GetFilePathInfoParams) WithHTTPClient(client *http.Client) *GetFilePathInfoParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get file path info params
func (o *GetFilePathInfoParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithFilePath adds the filePath to the get file path info params
func (o *GetFilePathInfoParams) WithFilePath(filePath string) *GetFilePathInfoParams {
	o.SetFilePath(filePath)
	return o
}

// SetFilePath adds the filePath to the get file path info params
func (o *GetFilePathInfoParams) SetFilePath(filePath string) {
	o.FilePath = filePath
}

// WriteToRequest writes these params to a swagger request
func (o *GetFilePathInfoParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param filePath
	if err := r.SetPathParam("filePath", o.FilePath); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}