package main

import (
	"fmt"
	"strings"

	"github.com/MikeRoetgers/deploi/protobuf"
)

// responseError can be used to package ResponseHeader errors into a Go error
type responseError struct {
	Errors []*protobuf.Error
}

// newResponseError takes a ResponseHeader and extracts the errors
func newResponseError(header *protobuf.ResponseHeader) *responseError {
	return &responseError{
		Errors: header.Errors,
	}
}

func (r *responseError) Error() string {
	var output []string
	for _, e := range r.Errors {
		output = append(output, fmt.Sprintf("Code: %s | Message: %s", e.Code, e.Message))
	}
	return strings.Join(output, " || ")
}
