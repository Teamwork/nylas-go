package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLabels(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"limit":  {"3"},
		"offset": {"1"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/labels")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(labelsJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Labels(context.Background(), &LabelsOptions{
		Offset: 1,
		Limit:  3,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Label{
		{
			AccountID:   "crkr5ct7aa3edvipotb****",
			DisplayName: "Inbox",
			ID:          "atamsqdb355jqyj0zhhatu3ao",
			Name:        "inbox",
			Object:      "label",
		},
		{
			AccountID:   "crkr5ct7aa3edvipotb****",
			DisplayName: "Trash",
			ID:          "558oz6v9e558so4najbpaonn8",
			Name:        "trash",
			Object:      "label",
		},
		{
			AccountID:   "crkr5ct7aa3edvipotb****",
			DisplayName: "Sent Mail",
			ID:          "3myg0x7i45lkn5xmel3m9d3v1",
			Name:        "sent",
			Object:      "label",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Labels: (-got +want):\n%s", diff)
	}
}

func TestLabelsCount(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"view": {ViewCount},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/labels")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write([]byte(`{"count":5}`))
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.LabelsCount(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 5

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("count: (-got +want):\n%s", diff)
	}
}

var labelsJSON = []byte(`[
  {
    "account_id": "crkr5ct7aa3edvipotb****",
    "display_name": "Inbox",
    "id": "atamsqdb355jqyj0zhhatu3ao",
    "name": "inbox",
    "object": "label"
  },
  {
    "account_id": "crkr5ct7aa3edvipotb****",
    "display_name": "Trash",
    "id": "558oz6v9e558so4najbpaonn8",
    "name": "trash",
    "object": "label"
  },
  {
    "account_id": "crkr5ct7aa3edvipotb****",
    "display_name": "Sent Mail",
    "id": "3myg0x7i45lkn5xmel3m9d3v1",
    "name": "sent",
    "object": "label"
  }
]`)
