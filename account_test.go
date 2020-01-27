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

func TestAccounts(t *testing.T) {
	clientID := "clientID"
	clientSecret := "clientSecret"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, clientSecret, "")
		wantPath := fmt.Sprintf("/a/%s/accounts", clientID)
		assertMethodPath(t, r, http.MethodGet, wantPath)

		_, _ = w.Write(managementAccountsJSON)
	}))
	defer ts.Close()

	client := NewClient(clientID, clientSecret, withTestServer(ts))
	got, err := client.Accounts(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []BillingAccount{
		{
			ID:           "622x1k5v1ujh55t6ucel7av4",
			AccountID:    "622x1k5v1ujh55t6ucel7av4",
			BillingState: "free",
			Email:        "example@example.com",
			Provider:     "yahoo",
			SyncState:    "running",
		},
		{
			ID:           "123rvgm1iccsgnjj7nn6jwu1",
			AccountID:    "123rvgm1iccsgnjj7nn6jwu1",
			BillingState: "paid",
			Email:        "example@example.com",
			Provider:     "gmail",
			SyncState:    "running",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Message: (-got +want):\n%s", diff)
	}
}

func TestCancelAccount(t *testing.T) {
	clientID := "clientID"
	clientSecret := "clientSecret"
	accountID := "accountID"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, clientSecret, "")
		wantPath := fmt.Sprintf("/a/%s/accounts/%s/downgrade", clientID, accountID)
		assertMethodPath(t, r, http.MethodPost, wantPath)

		_, _ = w.Write([]byte(`{"success":"true"}`))
	}))
	defer ts.Close()

	client := NewClient(clientID, clientSecret, withTestServer(ts))
	err := client.CancelAccount(context.Background(), accountID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReactivateAccount(t *testing.T) {
	clientID := "clientID"
	clientSecret := "clientSecret"
	accountID := "accountID"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, clientSecret, "")
		wantPath := fmt.Sprintf("/a/%s/accounts/%s/upgrade", clientID, accountID)
		assertMethodPath(t, r, http.MethodPost, wantPath)

		_, _ = w.Write([]byte(`{"success":"true"}`))
	}))
	defer ts.Close()

	client := NewClient(clientID, clientSecret, withTestServer(ts))
	err := client.ReactivateAccount(context.Background(), accountID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRevokeAccountTokens(t *testing.T) {
	clientID := "clientID"
	clientSecret := "clientSecret"
	accountID := "accountID"
	wantBody := []byte(`{"keep_access_token":"keep_me"}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, clientSecret, "")
		wantPath := fmt.Sprintf("/a/%s/accounts/%s/revoke-all", clientID, accountID)
		assertMethodPath(t, r, http.MethodPost, wantPath)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("req body: (-got +want):\n%s", diff)
		}

		_, _ = w.Write([]byte(`{"success":"true"}`))
	}))
	defer ts.Close()

	client := NewClient(clientID, clientSecret, withTestServer(ts))
	err := client.RevokeAccountTokens(context.Background(), accountID, String("keep_me"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

var managementAccountsJSON = []byte(`[
  {
    "account_id": "622x1k5v1ujh55t6ucel7av4",
    "billing_state": "free",
    "email": "example@example.com",
    "id": "622x1k5v1ujh55t6ucel7av4",
    "provider": "yahoo",
    "sync_state": "running",
    "trial": false
  },
  {
    "account_id": "123rvgm1iccsgnjj7nn6jwu1",
    "billing_state": "paid",
    "email": "example@example.com",
    "id": "123rvgm1iccsgnjj7nn6jwu1",
    "provider": "gmail",
    "sync_state": "running",
    "trial": false
  }
]`)
