package nylas

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMessages(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"any_email":       {"a@example.com,b@example.com"},
		"bcc":             {"f@example.com"},
		"cc":              {"e@example.com"},
		"filename":        {"filename"},
		"from":            {"d@example.com"},
		"has_attachment":  {"true"},
		"in":              {"in"},
		"limit":           {"1"},
		"offset":          {"2"},
		"received_after":  {"6"},
		"received_before": {"5"},
		"starred":         {"true"},
		"subject":         {"subject"},
		"thread_id":       {"threadid"},
		"to":              {"c@example.com"},
		"unread":          {"true"},
		"view":            {"ids"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(messagesJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Messages(context.Background(), &MessagesOptions{
		AnyEmail:       []string{"a@example.com", "b@example.com"},
		BCC:            "f@example.com",
		CC:             "e@example.com",
		Filename:       "filename",
		From:           "d@example.com",
		HasAttachment:  true,
		In:             "in",
		Limit:          1,
		Offset:         2,
		ReceivedAfter:  6,
		ReceivedBefore: 5,
		Starred:        true,
		Subject:        "subject",
		ThreadID:       "threadid",
		To:             "c@example.com",
		Unread:         true,
		View:           "ids",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Message{
		{

			AccountID: "43jf****",
			BCC:       []Participant{},
			Body:      "<html>\n<head>\n <meta charset=\"UTF-8\">\n <style type=\"text/css\">\n html {\n -webkit-text-size-adjust:none;\n }\n body {\n width:100%;\n margin:0 auto;\n padding:0;\n}\n  p {\n width:280px;\n line-height: 16px;\n letter-spacing: 0.5px;\n }\n </style>\n <title>Welcome  ...  </html>",
			CC:        []Participant{},
			Date:      1557950729,
			Events:    []interface{}{},
			Files:     []File{},
			Folder: Folder{
				DisplayName: "Inbox",
				ID:          "7hcg****",
				Name:        "inbox",
			},
			From: []Participant{
				{
					Email: "no-reply@cc.yahoo-inc.com",
					Name:  "Yahoo",
				},
			},
			ID:     "7a8939****",
			Object: "message",
			ReplyTo: []Participant{
				{
					Email: "no-reply@cc.yahoo-inc.com",
					Name:  "Yahoo",
				},
			},
			Snippet:  "Hi James, james****@yahoo.com. Welcome.",
			Starred:  false,
			Subject:  "Welcome",
			ThreadID: "cvsp****",
			To: []Participant{
				{
					Email: "james****@yahoo.com",
					Name:  "",
				},
			},
			Unread: true,
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Messages: (-got +want):\n%s", diff)
	}
}

var messagesJSON = []byte(`[
    {
	"account_id": "43jf****",
	"bcc": [],
	"body": "<html>\n<head>\n <meta charset=\"UTF-8\">\n <style type=\"text/css\">\n html {\n -webkit-text-size-adjust:none;\n }\n body {\n width:100%;\n margin:0 auto;\n padding:0;\n}\n  p {\n width:280px;\n line-height: 16px;\n letter-spacing: 0.5px;\n }\n </style>\n <title>Welcome  ...  </html>",
	"cc": [],
	"date": 1557950729,
	"events": [],
	"files": [],
	"folder": {
	    "display_name": "Inbox",
	    "id": "7hcg****",
	    "name": "inbox"
	},
	"from": [
	    {
		"email": "no-reply@cc.yahoo-inc.com",
		"name": "Yahoo"
	    }
	],
	"id": "7a8939****",
	"object": "message",
	"reply_to": [
	    {
		"email": "no-reply@cc.yahoo-inc.com",
		"name": "Yahoo"
	    }
	],
	"snippet": "Hi James, james****@yahoo.com. Welcome.",
	"starred": false,
	"subject": "Welcome",
	"thread_id": "cvsp****",
	"to": [
	    {
		"email": "james****@yahoo.com",
		"name": ""
	    }
	],
	"unread": true
    }
]`)
