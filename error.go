package nylas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Error returned from the Nylas API.
// See: https://docs.nylas.com/reference#errors
// See: https://docs.nylas.com/reference#section-sending-errors
type Error struct {
	StatusCode int    `json:"-"`
	Body       []byte `json:"-"`

	Message     string `json:"message"`
	Type        string `json:"type"`
	ServerError string `json:"server_error"`
}

// Error implements the error interface.
func (e Error) Error() string {
	s := fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
	if e.ServerError != "" {
		s += ": " + e.ServerError
	}
	return s
}

// NewError creates a new Error from an API response.
func NewError(resp *http.Response) error {
	apiErr := Error{StatusCode: resp.StatusCode}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		apiErr.Body = data
		if err := json.Unmarshal(data, &apiErr); err != nil {
			apiErr.Type = "unknown_error_format"
			apiErr.Message = string(data)
		}
	}
	return &apiErr
}
