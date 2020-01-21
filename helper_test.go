package nylas

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func withTestServer(ts *httptest.Server) Option {
	return WithBaseURL(ts.URL)
}

func assertBasicAuth(t *testing.T, r *http.Request, user, pass string) {
	t.Helper()
	gotUser, gotPass, ok := r.BasicAuth()
	if !ok {
		t.Errorf("basic auth not provided")
	}
	if user != gotUser || pass != gotPass {
		t.Errorf("basic auth: got %q:%q; want %q;%q", gotUser, gotPass, user, pass)
	}
}
