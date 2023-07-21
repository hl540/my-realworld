package errors

import (
	"fmt"
	"errors"
	stdhttp "net/http"

	"github.com/go-kratos/kratos/v2/transport/http"
)

// HTTPError is an HTTP error.
type HTTPError struct {
	code int
	Errors map[string][]string `json:"errors"`
}

func NewHTTPError(code int, field string, msg string) *HTTPError {
	return &HTTPError{
		code: code,
		Errors: map[string][]string{
			field: {msg},
		},
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("HTTPError code: %d message: %v", e.code, e.Errors)
}

// FromError try to convert an error to *HTTPError.
func FromError(err error) *HTTPError {
	if err == nil {
		return nil
	}
	if se := new(HTTPError); errors.As(err, &se) {
		return se
	}
	return &HTTPError{code: 500}
}

func ErrorEncoder(w stdhttp.ResponseWriter, r *stdhttp.Request, err error) {
	se := FromError(err)
	codec, _ := http.CodecForRequest(r, "Accept")
	body, err := codec.Marshal(se)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/"+codec.Name())
	w.WriteHeader(se.code)
	_, _ = w.Write(body)
}