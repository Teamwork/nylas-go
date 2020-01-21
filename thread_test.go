package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestThreads(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"any_email":           {"a@example.com,b@example.com"},
		"bcc":                 {"f@example.com"},
		"cc":                  {"e@example.com"},
		"filename":            {"filename"},
		"from":                {"d@example.com"},
		"in":                  {"in"},
		"last_message_after":  {"4"},
		"last_message_before": {"3"},
		"limit":               {"1"},
		"offset":              {"2"},
		"started_after":       {"6"},
		"started_before":      {"5"},
		"subject":             {"subject"},
		"to":                  {"c@example.com"},
		"unread":              {"true"},
		"view":                {"ids"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		if gotQuery := r.URL.Query(); !reflect.DeepEqual(gotQuery, wantQuery) {
			t.Errorf("query params:\ngot:\n%#v\nwant:\n%+v", gotQuery, wantQuery)
		}
		_, _ = w.Write(threadsJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Threads(context.Background(), &ThreadsOptions{
		AnyEmail:          []string{"a@example.com", "b@example.com"},
		BCC:               "f@example.com",
		CC:                "e@example.com",
		Filename:          "filename",
		From:              "d@example.com",
		In:                "in",
		LastMessageAfter:  4,
		LastMessageBefore: 3,
		Limit:             1,
		Offset:            2,
		StartedAfter:      6,
		StartedBefore:     5,
		Subject:           "subject",
		To:                "c@example.com",
		Unread:            true,
		View:              "ids",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Thread{
		{
			AccountID:             "1234***",
			DraftIDs:              []string{},
			FirstMessageTimestamp: 1557950729,
			Folders: []Folder{
				{
					DisplayName: "Inbox",
					ID:          "4567****",
					Name:        "inbox",
				},
			},
			HasAttachments:               false,
			ID:                           "4312****",
			LastMessageReceivedTimestamp: 1557950729,
			LastMessageSentTimestamp:     0,
			LastMessageTimestamp:         1557950729,
			MessageIDs: []string{
				"5634***",
			},
			Object: "thread",
			Participants: []Participant{
				{
					Email: "no-reply@cc.yahoo-inc.com",
					Name:  "Yahoo",
				},
				{
					Email: "james****@yahoo.com",
					Name:  "",
				},
			},
			Snippet: "Hi James, welcome.",
			Starred: false,
			Subject: "Security settings changed on your Yahoo account",
			Unread:  false,
			Version: 1,
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Threads: (-got +want):\n%s", diff)
	}
}

var threadsJSON = []byte(`[
    {
	"account_id": "1234***",
	"draft_ids": [],
	"first_message_timestamp": 1557950729,
	"folders": [
	    {
		"display_name": "Inbox",
		"id": "4567****",
		"name": "inbox"
	    }
	],
	"has_attachments": false,
	"id": "4312****",
	"last_message_received_timestamp": 1557950729,
	"last_message_sent_timestamp": null,
	"last_message_timestamp": 1557950729,
	"message_ids": [
	    "5634***"
	],
	"object": "thread",
	"participants": [
	    {
		"email": "no-reply@cc.yahoo-inc.com",
		"name": "Yahoo"
	    },
	    {
		"email": "james****@yahoo.com",
		"name": ""
	    }
	],
	"snippet": "Hi James, welcome.",
	"starred": false,
	"subject": "Security settings changed on your Yahoo account",
	"unread": false,
	"version": 1
    }
]`)
