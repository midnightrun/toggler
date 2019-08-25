// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/toggler-io/toggler/lib/go/models"
)

// NewWebsocketParams creates a new WebsocketParams object
// with the default values initialized.
func NewWebsocketParams() *WebsocketParams {
	var ()
	return &WebsocketParams{

		timeout: cr.DefaultTimeout,
	}
}

// NewWebsocketParamsWithTimeout creates a new WebsocketParams object
// with the default values initialized, and the ability to set a timeout on a request
func NewWebsocketParamsWithTimeout(timeout time.Duration) *WebsocketParams {
	var ()
	return &WebsocketParams{

		timeout: timeout,
	}
}

// NewWebsocketParamsWithContext creates a new WebsocketParams object
// with the default values initialized, and the ability to set a context for a request
func NewWebsocketParamsWithContext(ctx context.Context) *WebsocketParams {
	var ()
	return &WebsocketParams{

		Context: ctx,
	}
}

// NewWebsocketParamsWithHTTPClient creates a new WebsocketParams object
// with the default values initialized, and the ability to set a custom HTTPClient for a request
func NewWebsocketParamsWithHTTPClient(client *http.Client) *WebsocketParams {
	var ()
	return &WebsocketParams{
		HTTPClient: client,
	}
}

/*WebsocketParams contains all the parameters to send to the API endpoint
for the websocket operation typically these are written to a http.Request
*/
type WebsocketParams struct {

	/*Body*/
	Body *models.WebsocketRequestPayload

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithTimeout adds the timeout to the websocket params
func (o *WebsocketParams) WithTimeout(timeout time.Duration) *WebsocketParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the websocket params
func (o *WebsocketParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the websocket params
func (o *WebsocketParams) WithContext(ctx context.Context) *WebsocketParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the websocket params
func (o *WebsocketParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the websocket params
func (o *WebsocketParams) WithHTTPClient(client *http.Client) *WebsocketParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the websocket params
func (o *WebsocketParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the websocket params
func (o *WebsocketParams) WithBody(body *models.WebsocketRequestPayload) *WebsocketParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the websocket params
func (o *WebsocketParams) SetBody(body *models.WebsocketRequestPayload) {
	o.Body = body
}

// WriteToRequest writes these params to a swagger request
func (o *WebsocketParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
