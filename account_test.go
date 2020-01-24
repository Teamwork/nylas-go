package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestAccount(t *testing.T) {
	accessToken := "accessToken"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/account")

		_, _ = w.Write(accountJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Account(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Account{
		ID:               "awa6ltos76vz5hvphkp8k17nt",
		AccountID:        "awa6ltos76vz5hvphkp8k17nt",
		Object:           "account",
		Name:             "Ben Bitdiddle",
		EmailAddress:     "benbitdiddle@gmail.com",
		Provider:         "gmail",
		OrganizationUnit: "label",
		SyncState:        "running",
		LinkedAt:         1470231381,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Message: (-got +want):\n%s", diff)
	}
}

var accountJSON = []byte(`{
    "id": "awa6ltos76vz5hvphkp8k17nt",
    "account_id": "awa6ltos76vz5hvphkp8k17nt",
    "object": "account",
    "name": "Ben Bitdiddle",
    "email_address": "benbitdiddle@gmail.com",
    "provider": "gmail",
    "organization_unit": "label",
    "sync_state": "running",
    "linked_at": 1470231381
}`)
