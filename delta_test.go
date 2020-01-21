package nylas

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLastestDeltaCursor(t *testing.T) {
	accessToken := "accessToken"
	wantCursor := "aqb0llc2ioo0***"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		fmt.Fprintf(w, `{"cursor": "%s"}`, wantCursor)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	cursor, err := client.LatestDeltaCursor(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cursor != wantCursor {
		t.Errorf("want cursor: %q; got cursor %q", wantCursor, cursor)
	}
}

func TestDeltas(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"cursor":        {"cursor"},
		"include_types": {"a,b"},
		"exclude_types": {"c,d"},
		"view":          {"ids"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(deltaJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Deltas(context.Background(), "cursor", &DeltasOptions{
		// Not valid to combine these two, just doing it for test
		IncludeTypes: []string{"a", "b"},
		ExcludeTypes: []string{"c", "d"},
		View:         "ids",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	threadAttributes := []byte(`{
		"draft_ids": [
		    "diu1tytx7p9***"
		],
		"first_message_timestamp": 1414778436,
		"id": "71ormxuivtg52p141tpgjk3vi",
		"last_message_timestamp": 1414778436,
		"message_ids": [],
		"account_id": "f3b0j663wmm***",
		"object": "thread",
		"participants": [],
		"snippet": "",
		"subject": "Hello World!",
		"folders": [
		    {
			"name": "drafts",
			"id": "e3b0j663wmm2***"
		    }
		],
		"unread": false,
		"starred": false
	    }`)
	want := DeltaResponse{
		CursorStart: "aqb0llc2i***",
		CursorEnd:   "5u9kwbgyq8***",
		Deltas: []Delta{
			{
				ID:         "aqb0llc2io***",
				Event:      "modify",
				Object:     "thread",
				Cursor:     "7ciyf89wp***",
				Attributes: threadAttributes,
			},
			{
				Cursor: "9vsuralamr***",
				Event:  "delete",
				ID:     "5oly0nmkf***",
				Object: "folder",
			},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("DeltaResponse: (-got +want):\n%s", diff)
	}
}

var deltaJSON = []byte(`{
    "cursor_start": "aqb0llc2i***",
    "cursor_end": "5u9kwbgyq8***",
    "deltas": [
	{
	    "id": "aqb0llc2io***",
	    "event": "modify",
	    "object": "thread",
	    "cursor": "7ciyf89wp***",
	    "attributes": {
		"draft_ids": [
		    "diu1tytx7p9***"
		],
		"first_message_timestamp": 1414778436,
		"id": "71ormxuivtg52p141tpgjk3vi",
		"last_message_timestamp": 1414778436,
		"message_ids": [],
		"account_id": "f3b0j663wmm***",
		"object": "thread",
		"participants": [],
		"snippet": "",
		"subject": "Hello World!",
		"folders": [
		    {
			"name": "drafts",
			"id": "e3b0j663wmm2***"
		    }
		],
		"unread": false,
		"starred": false
	    }
	},
	{
	    "cursor": "9vsuralamr***",
	    "event": "delete",
	    "id": "5oly0nmkf***",
	    "object": "folder"
	}
    ]
}`)
