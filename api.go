package deploi

import (
	"fmt"
	"strings"

	"github.com/MikeRoetgers/deploi/protobuf"
)

// ResponseError can be used to package ResponseHeader errors into a Go error
type ResponseError struct {
	Errors []*protobuf.Error
}

// NewResponseError takes a ResponseHeader and extracts the errors
func NewResponseError(header *protobuf.ResponseHeader) *ResponseError {
	return &ResponseError{
		Errors: header.Errors,
	}
}

func (r *ResponseError) Error() string {
	var output []string
	for _, e := range r.Errors {
		output = append(output, fmt.Sprintf("Code: %s | Message: %s", e.Code, e.Message))
	}
	return strings.Join(output, " || ")
}
