// Package httpapi provides API on HTTP layer to the toggler service.
//
// The purpose of this application is to provide API over HTTP to toggler service,
// in which you can interact with the service in a programmatic way.
//
//     BasePath: /api/v1
//     Version: v1
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package httpapi

// EnrollmentResponse returns information about the requester's rollout feature enrollment status.
// swagger:response enrollmentResponse
type EnrollmentResponse struct {
	// in: body
	Body EnrollmentResponseBody
}

// EnrollmentResponse will be returned when release flag status is requested.
// The content will be always given, regardless if the flag exists or not.
// This helps the developers to use it as a null object, regardless the toggler service state.
type EnrollmentResponseBody struct {
	// release flag enrollment status.
	Enrollment bool `json:"enrollment"`
}
