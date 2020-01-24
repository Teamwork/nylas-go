package nylas

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
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

func assertQueryParams(t *testing.T, r *http.Request, want url.Values) {
	t.Helper()
	got := r.URL.Query()
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("query params: (-got +want):\n%s", diff)
	}
}

func assertMethodPath(t *testing.T, r *http.Request, method, path string) {
	t.Helper()
	if r.Method != method || r.URL.Path != path {
		t.Errorf("req: got %v %v; want %v %v", r.Method, r.URL.Path, method, path)
	}
}
