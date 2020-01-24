package nylas

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewError(t *testing.T) {
	body := func(s string) io.ReadCloser {
		return ioutil.NopCloser(bytes.NewBufferString(s))
	}
	tests := map[string]struct {
		in  *http.Response
		out error
	}{
		"nylas error": {
			in: &http.Response{
				StatusCode: 400,
				Body: body(`{
					"message": "Invalid datetime value z for start_time",
					"type": "invalid_request_error"
				}`),
			},
			out: &Error{
				StatusCode: 400,
				Body: []byte(`{
					"message": "Invalid datetime value z for start_time",
					"type": "invalid_request_error"
				}`),
				Message: "Invalid datetime value z for start_time",
				Type:    "invalid_request_error",
			},
		},
		"invalid json": {
			in: &http.Response{
				StatusCode: 403,
				Body:       body(`something went wrong`),
			},
			out: &Error{
				StatusCode: 403,
				Body:       []byte(`something went wrong`),
				Message:    "something went wrong",
				Type:       "unknown_error_format",
			},
		},
	}

	for desc, tt := range tests {
		t.Run(desc, func(t *testing.T) {
			err := NewError(tt.in)
			if diff := cmp.Diff(err, tt.out); diff != "" {
				t.Errorf("req body: (-got +want):\n%s", diff)
			}
		})
	}
}
