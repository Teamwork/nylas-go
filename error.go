package nylas

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Error returned from the Nylas API.
// See: https://docs.nylas.com/reference#errors
type Error struct {
	StatusCode int    `json:"-"`
	Body       []byte `json:"-"`

	Message string `json:"message"`
	Type    string `json:"type"`
}

// Error implements the error interface.
func (e Error) Error() string {
	return fmt.Sprintf("%d: %s", e.StatusCode, e.Message)
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
