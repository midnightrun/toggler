// Code generated by go-swagger; DO NOT EDIT.

package pilot

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

// NewGetPilotConfigParams creates a new GetPilotConfigParams object
// with the default values initialized.
func NewGetPilotConfigParams() *GetPilotConfigParams {
	var ()
	return &GetPilotConfigParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewGetPilotConfigParamsWithTimeout creates a new GetPilotConfigParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewGetPilotConfigParamsWithTimeout(timeout time.Duration) *GetPilotConfigParams {
	var ()
	return &GetPilotConfigParams{

		timeout: timeout,
	}
}

// NewGetPilotConfigParamsWithContext creates a new GetPilotConfigParams object
// with the default values initialized, and the ability to set a context for a request
func NewGetPilotConfigParamsWithContext(ctx context.Context) *GetPilotConfigParams {
	var ()
	return &GetPilotConfigParams{

		Context: ctx,
	}
}

// NewGetPilotConfigParamsWithHTTPClient creates a new GetPilotConfigParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewGetPilotConfigParamsWithHTTPClient(client *http.Client) *GetPilotConfigParams {
	var ()
	return &GetPilotConfigParams{
		HTTPClient: client,
	}
}

/*GetPilotConfigParams contains all the parameters to send to the API endpoint
for the get pilot config operation typically these are written to a http.Request
*/
type GetPilotConfigParams struct {

	/*Body*/
	Body GetPilotConfigBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the get pilot config params
func (o *GetPilotConfigParams) WithTimeout(timeout time.Duration) *GetPilotConfigParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get pilot config params
func (o *GetPilotConfigParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get pilot config params
func (o *GetPilotConfigParams) WithContext(ctx context.Context) *GetPilotConfigParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get pilot config params
func (o *GetPilotConfigParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get pilot config params
func (o *GetPilotConfigParams) WithHTTPClient(client *http.Client) *GetPilotConfigParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get pilot config params
func (o *GetPilotConfigParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the get pilot config params
func (o *GetPilotConfigParams) WithBody(body GetPilotConfigBody) *GetPilotConfigParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the get pilot config params
func (o *GetPilotConfigParams) SetBody(body GetPilotConfigBody) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *GetPilotConfigParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if err := r.SetBodyParam(o.Body); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}