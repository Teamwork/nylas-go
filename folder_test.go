package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFolders(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"limit":  {"3"},
		"offset": {"1"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/folders")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(foldersJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Folders(context.Background(), &FoldersOptions{
		Offset: 1,
		Limit:  3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Folder{
		{
			ID:          "4zv7p****",
			Object:      "folder",
			Name:        "inbox",
			DisplayName: "INBOX",
			AccountID:   "awa6lt****",
		},
		{
			ID:          "76zrf****",
			Name:        "archive",
			DisplayName: "2015 Archive",
			AccountID:   "awa6lto****",
			Object:      "folder",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Folders: (-got +want):\n%s", diff)
	}
}

func TestFoldersCount(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"view": {ViewCount},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/folders")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write([]byte(`{"count":5}`))
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.FoldersCount(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 5

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("count: (-got +want):\n%s", diff)
	}
}

var foldersJSON = []byte(`[
    {
	"id": "4zv7p****",
	"object": "folder",
	"name": "inbox",
	"display_name": "INBOX",
	"account_id": "awa6lt****"
    },
    {
	"id": "76zrf****",
	"name": "archive",
	"display_name": "2015 Archive",
	"account_id": "awa6lto****",
	"object": "folder"
    }
]`)
