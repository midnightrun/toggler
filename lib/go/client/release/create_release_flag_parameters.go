// Code generated by go-swagger; DO NOT EDIT.

package release

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

// NewCreateReleaseFlagParams creates a new CreateReleaseFlagParams object
// with the default values initialized.
func NewCreateReleaseFlagParams() *CreateReleaseFlagParams {
	var ()
	return &CreateReleaseFlagParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewCreateReleaseFlagParamsWithTimeout creates a new CreateReleaseFlagParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewCreateReleaseFlagParamsWithTimeout(timeout time.Duration) *CreateReleaseFlagParams {
	var ()
	return &CreateReleaseFlagParams{

		timeout: timeout,
	}
}

// NewCreateReleaseFlagParamsWithContext creates a new CreateReleaseFlagParams object
// with the default values initialized, and the ability to set a context for a request
func NewCreateReleaseFlagParamsWithContext(ctx context.Context) *CreateReleaseFlagParams {
	var ()
	return &CreateReleaseFlagParams{

		Context: ctx,
	}
}

// NewCreateReleaseFlagParamsWithHTTPClient creates a new CreateReleaseFlagParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewCreateReleaseFlagParamsWithHTTPClient(client *http.Client) *CreateReleaseFlagParams {
	var ()
	return &CreateReleaseFlagParams{
		HTTPClient: client,
	}
}

/*CreateReleaseFlagParams contains all the parameters to send to the API endpoint
for the create release flag operation typically these are written to a http.Request
*/
type CreateReleaseFlagParams struct {

	/*Body*/
	Body CreateReleaseFlagBody

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the create release flag params
func (o *CreateReleaseFlagParams) WithTimeout(timeout time.Duration) *CreateReleaseFlagParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the create release flag params
func (o *CreateReleaseFlagParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the create release flag params
func (o *CreateReleaseFlagParams) WithContext(ctx context.Context) *CreateReleaseFlagParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the create release flag params
func (o *CreateReleaseFlagParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the create release flag params
func (o *CreateReleaseFlagParams) WithHTTPClient(client *http.Client) *CreateReleaseFlagParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the create release flag params
func (o *CreateReleaseFlagParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the create release flag params
func (o *CreateReleaseFlagParams) WithBody(body CreateReleaseFlagBody) *CreateReleaseFlagParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the create release flag params
func (o *CreateReleaseFlagParams) SetBody(body CreateReleaseFlagBody) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *CreateReleaseFlagParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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