package errors

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrorResponse Contains information about http error status code as json
type errorResponse struct {
	Err            error  `json:"Runtime error,omitempty"` // low-level runtime error
	HTTPStatusCode int    `json:"Status code,omitempty"`   // http response status code
	StatusText     string `json:"Message"`                 // user-level status message
	AppCode        int64  `json:"Code,omitempty"`          // application-specific error code
	ErrorText      string `json:"Error,omitempty"`         // application-level error message, for debugging
}

var (
	//Err400 http status: Bad request
	Err400 = &errorResponse{
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest)}

	//Err401 http status: Unauthorized
	Err401 = &errorResponse{
		HTTPStatusCode: http.StatusUnauthorized,
		StatusText:     http.StatusText(http.StatusUnauthorized)}

	//Err403 http status: Forbidden
	Err403 = &errorResponse{
		HTTPStatusCode: http.StatusForbidden,
		StatusText:     http.StatusText(http.StatusForbidden)}

	//Err404 http status: Not found
	Err404 = &errorResponse{
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound)}

	//Err422 http status: Unprocessable entity
	Err422 = &errorResponse{
		HTTPStatusCode: http.StatusUnprocessableEntity,
		StatusText:     http.StatusText(http.StatusUnprocessableEntity)}

	//Err500 http status: internal server error
	Err500 = &errorResponse{
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError)}

	//Err503 http status: Service unavailable
	Err503 = &errorResponse{
		HTTPStatusCode: http.StatusServiceUnavailable,
		StatusText:     http.StatusText(http.StatusServiceUnavailable)}
)

// Render Shows the errorResponse pointer in request
func (e *errorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
