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

	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
)

// NewListFilesParams creates a new ListFilesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewListFilesParams() *ListFilesParams {
	return &ListFilesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewListFilesParamsWithTimeout creates a new ListFilesParams object
// with the ability to set a timeout on a request.
func NewListFilesParamsWithTimeout(timeout time.Duration) *ListFilesParams {
	return &ListFilesParams{
		timeout: timeout,
	}
}

// NewListFilesParamsWithContext creates a new ListFilesParams object
// with the ability to set a context for a request.
func NewListFilesParamsWithContext(ctx context.Context) *ListFilesParams {
	return &ListFilesParams{
		Context: ctx,
	}
}

// NewListFilesParamsWithHTTPClient creates a new ListFilesParams object
// with the ability to set a custom HTTPClient for a request.
func NewListFilesParamsWithHTTPClient(client *http.Client) *ListFilesParams {
	return &ListFilesParams{
		HTTPClient: client,
	}
}

/* ListFilesParams contains all the parameters to send to the API endpoint
   for the list files operation.

   Typically these are written to a http.Request.
*/
type ListFilesParams struct {

	// Directory.
	Directory *models.DirectoryListRequest

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the list files params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListFilesParams) WithDefaults() *ListFilesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the list files params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ListFilesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the list files params
func (o *ListFilesParams) WithTimeout(timeout time.Duration) *ListFilesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the list files params
func (o *ListFilesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the list files params
func (o *ListFilesParams) WithContext(ctx context.Context) *ListFilesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the list files params
func (o *ListFilesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the list files params
func (o *ListFilesParams) WithHTTPClient(client *http.Client) *ListFilesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the list files params
func (o *ListFilesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithDirectory adds the directory to the list files params
func (o *ListFilesParams) WithDirectory(directory *models.DirectoryListRequest) *ListFilesParams {
	o.SetDirectory(directory)
	return o
}

// SetDirectory adds the directory to the list files params
func (o *ListFilesParams) SetDirectory(directory *models.DirectoryListRequest) {
	o.Directory = directory
}

// WriteToRequest writes these params to a swagger request
func (o *ListFilesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Directory != nil {
		if err := r.SetBodyParam(o.Directory); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
