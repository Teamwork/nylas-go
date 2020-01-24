package nylas

import (
	"context"
	"io/ioutil"
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
		assertMethodPath(t, r, http.MethodGet, "/messages")
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
		HasAttachment:  Bool(true),
		In:             "in",
		Limit:          1,
		Offset:         2,
		ReceivedAfter:  6,
		ReceivedBefore: 5,
		Starred:        Bool(true),
		Subject:        "subject",
		ThreadID:       "threadid",
		To:             "c@example.com",
		Unread:         Bool(true),
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

func TestMessagesCount(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"unread": {"true"},
		"view":   {ViewCount},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/messages")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write([]byte(`{"count":1}`))
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.MessagesCount(context.Background(), &MessagesOptions{
		Unread: Bool(true),
		View:   "ids",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 1

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("count: (-got +want):\n%s", diff)
	}
}

func TestMessage(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"view": {"expanded"},
	}
	id := "br57kcekhf1hsjq04y8aonkit"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/messages/"+id)
		assertQueryParams(t, r, wantQuery)

		_, _ = w.Write(getMessageJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Message(context.Background(), id, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Message{
		AccountID: "crkr5ct7aa3edvipotbj****",
		BCC:       []Participant{},
		Body:      "<div dir=\"ltr\">Body</div>",
		CC:        []Participant{},
		Date:      1579611155,
		Events:    []interface{}{},
		Files:     []File{},
		From: []Participant{
			{
				Email: "from@example.org",
				Name:  "From Name",
			},
		},
		Headers: struct {
			InReplyTo  string   `json:"In-Reply-To"`
			MessageID  string   `json:"Message-Id"`
			References []string `json:"References"`
		}{
			InReplyTo:  "",
			MessageID:  "<CAGkcA6KLq8q4bETj8+BhMLms1JrvaJ+5SvJVVz+u_Ok0y=iEoA@mail.gmail.com>",
			References: []string{},
		},
		ID: "br57kcekhf1hsjq04y8aonkit",
		Labels: []Label{
			{
				DisplayName: "Important",
				ID:          "a1ytpbvawxfaqua671478g1q0",
				Name:        "important",
			},
			{
				DisplayName: "Inbox",
				ID:          "atamsqdb355jqyj0zhhatu3ao",
				Name:        "inbox",
			},
		},
		Object:   "message",
		ReplyTo:  []Participant{},
		Snippet:  "Body",
		Starred:  true,
		Subject:  "Subject",
		ThreadID: "8r5awu0esbg8ct3wg5rj5sifp",
		To: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		Unread: true,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Message: (-got +want):\n%s", diff)
	}
}

func TestRawMessage(t *testing.T) {
	accessToken := "accessToken"
	id := "br57kcekhf1hsjq04y8aonkit"
	want := []byte(`Delivered-To: to@example.org
Received: by 2002:ab3:5e90:0:0:0:0:0 with SMTP id k16csp2294558TLC;
	Tue, 21 Jan 2020 04:52:47 -0800 (PST)
.......`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/messages/"+id)

		if v := r.Header.Get("Accept"); v != "message/rfc822" {
			t.Errorf("missing/incorrect accept header: %v", v)
		}
		_, _ = w.Write(want)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.RawMessage(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Message: (-got +want):\n%s", diff)
	}
}

func TestUpdateMessage(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{}
	id := "8r5awu0esbg8ct3wg5rj5sifp"
	wantBody := []byte(`{"unread":true,"starred":false,"folder_id":"folderid","label_ids":["label1","label2"]}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodPut, "/messages/"+id)
		assertQueryParams(t, r, wantQuery)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("req body: (-got +want):\n%s", diff)
		}
		_, _ = w.Write(getMessageJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	_, err := client.UpdateMessage(context.Background(), id, UpdateMessageRequest{
		Unread:   Bool(true),
		Starred:  Bool(false),
		FolderID: String("folderid"),
		LabelIDs: &[]string{"label1", "label2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

var getMessageJSON = []byte(`{
  "account_id": "crkr5ct7aa3edvipotbj****",
  "bcc": [],
  "body": "<div dir=\"ltr\">Body</div>",
  "cc": [],
  "date": 1579611155,
  "events": [],
  "files": [],
  "from": [
    {
      "email": "from@example.org",
      "name": "From Name"
    }
  ],
  "headers": {
    "In-Reply-To": null,
    "Message-Id": "<CAGkcA6KLq8q4bETj8+BhMLms1JrvaJ+5SvJVVz+u_Ok0y=iEoA@mail.gmail.com>",
    "References": []
  },
  "id": "br57kcekhf1hsjq04y8aonkit",
  "labels": [
    {
      "display_name": "Important",
      "id": "a1ytpbvawxfaqua671478g1q0",
      "name": "important"
    },
    {
      "display_name": "Inbox",
      "id": "atamsqdb355jqyj0zhhatu3ao",
      "name": "inbox"
    }
  ],
  "object": "message",
  "reply_to": [],
  "snippet": "Body",
  "starred": true,
  "subject": "Subject",
  "thread_id": "8r5awu0esbg8ct3wg5rj5sifp",
  "to": [
    {
      "email": "to@example.org",
      "name": "To Name"
    }
  ],
  "unread": true
}`)
