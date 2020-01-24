package nylas

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestConnectAuthorize(t *testing.T) {
	clientSecret := "clientSecret"
	wantBody := []byte(`{"client_id":"clientid","email_address":"email@example.org","name":"Name","provider":"imap","scopes":"email,calendar","settings":{"imap_host":"imap.host","imap_port":993,"imap_username":"imap.user","imap_password":"imap.pass","smtp_host":"smtp.host","smtp_port":465,"smtp_username":"smtp.user","smtp_password":"smtp.pass","ssl_required":true}}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/connect/authorize")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("req body: (-got +want):\n%s", diff)
		}
		_, _ = w.Write([]byte(`{"code":"code"}`))
	}))
	defer ts.Close()

	client := NewClient("clientid", clientSecret, withTestServer(ts))
	code, err := client.connectAuthorize(context.Background(), AuthorizeRequest{
		Name:         "Name",
		EmailAddress: "email@example.org",
		Settings: IMAPAuthorizeSettings{
			IMAPHost:     "imap.host",
			IMAPPort:     993,
			IMAPUsername: "imap.user",
			IMAPPassword: "imap.pass",
			SMTPHost:     "smtp.host",
			SMTPPort:     465,
			SMTPUsername: "smtp.user",
			SMTPPassword: "smtp.pass",
			SSLRequired:  true,
		},
		Scopes: []string{"email", "calendar"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if code != "code" {
		t.Errorf("response code: got %v; want code", code)
	}
}

func TestConnectToken(t *testing.T) {
	clientSecret := "clientSecret"
	wantBody := []byte(`{"client_id":"clientid","client_secret":"clientSecret","code":"code"}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertMethodPath(t, r, http.MethodPost, "/connect/token")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			fmt.Println(string(body))
			t.Errorf("req body: (-got +want):\n%s", diff)
		}
		_, _ = w.Write([]byte(`{"access_token":"accessToken"}`))
	}))
	defer ts.Close()

	client := NewClient("clientid", clientSecret, withTestServer(ts))
	acc, err := client.connectExchangeCode(context.Background(), "code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if acc.AccessToken != "accessToken" {
		t.Errorf("response accessToken: got %v; want accessToken", acc.AccessToken)
	}
}
