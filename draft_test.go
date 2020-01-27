package nylas

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDrafts(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"any_email": {"a@example.com,b@example.com"},
		"limit":     {"1"},
		"offset":    {"2"},
		"view":      {"count"},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/drafts")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write(draftsJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Drafts(context.Background(), &DraftsOptions{
		AnyEmail: []string{"a@example.com", "b@example.com"},
		Limit:    1,
		Offset:   2,
		View:     "count",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []Draft{
		{
			Message: Message{
				ID:        "43vfrmdu1***",
				Object:    "draft",
				AccountID: "43jf3n4***",
				ThreadID:  "46wnzkxa***",
				From: []Participant{
					{Email: "nylastest***@yahoo.com", Name: "John Doe"},
				},
				To: []Participant{
					{Email: "{{email}}", Name: "{{name}}"},
				},
				CC:      []Participant{},
				BCC:     []Participant{},
				ReplyTo: []Participant{},
				Subject: "ugh?",
				Date:    1559763005,
				Body:    "Hello, how are you?",
				Snippet: "Hello, how are you?",
				Events:  []interface{}{},
				Files:   []File{},
			},
			Version:          0,
			ReplyToMessageID: "",
		},
		{
			Message: Message{
				ID:        "92c7gucghzh***",
				Object:    "draft",
				AccountID: "43jf3n4***",
				ThreadID:  "e48pmw615r2****",
				From: []Participant{
					{Email: "nylastest***@yahoo.com", Name: "John Doe"},
				},
				To: []Participant{
					{Email: "{{email}}", Name: "{{name}}"},
				},
				CC:      []Participant{},
				BCC:     []Participant{},
				ReplyTo: []Participant{},
				Subject: "Hello",
				Date:    1559762902,
				Body:    "Hello, how are you?",
				Snippet: "Hello, how are you?",
				Events:  []interface{}{},
				Files:   []File{},
				Folder: Folder{
					ID:          "eeangfw9vm***",
					DisplayName: "Draft",
					Name:        "drafts",
				},
			},
			Version:          1,
			ReplyToMessageID: "132i3h12u3hadw",
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Drafts: (-got +want):\n%s", diff)
	}
}

func TestDraftsCount(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{
		"any_email": {"a@example.com,b@example.com"},
		"view":      {ViewCount},
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/drafts")
		assertQueryParams(t, r, wantQuery)
		_, _ = w.Write([]byte(`{"count":1}`))
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.DraftsCount(context.Background(), &DraftsOptions{
		AnyEmail: []string{"a@example.com", "b@example.com"},
		View:     "dont use this value",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := 1

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("count: (-got +want):\n%s", diff)
	}
}

func TestDraft(t *testing.T) {
	accessToken := "accessToken"
	id := "br57kcekhf1hsjq04y8aonkit"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodGet, "/drafts/"+id)

		_, _ = w.Write(getDraftJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	got, err := client.Draft(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := Draft{
		Message: Message{
			ID:        "43vfrmdu1***",
			Object:    "draft",
			AccountID: "43jf3n4e***",
			ThreadID:  "46wnzkxaz***",
			From:      []Participant{{Email: "nylastest***@yahoo.com", Name: "John Doe"}},
			To:        []Participant{{Email: "{{email}}", Name: "{{name}}"}},
			CC:        []Participant{},
			BCC:       []Participant{},
			ReplyTo:   []Participant{},
			Subject:   "ugh?",
			Date:      1559763005,
			Body:      "Hello, how are you?",
			Snippet:   "Hello, how are you?",
			Events:    []interface{}{},
			Files:     []File{},
			Folder:    Folder{ID: "eeangfw9***", DisplayName: "Draft", Name: "drafts"},
		},
		Version: 2,
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Draft: (-got +want):\n%s", diff)
	}
}

func TestCreateDraft(t *testing.T) {
	accessToken := "accessToken"
	wantBody := []byte(`{"subject":"Subject","from":[{"email":"from@example.org","name":"From Name"}],"to":[{"email":"to@example.org","name":"To Name"}],"cc":[{"email":"to@example.org","name":"To Name"}],"bcc":[{"email":"to@example.org","name":"To Name"}],"reply_to":[{"email":"replyto@example.org","name":"ReplyTo Name"}],"body":"body","file_ids":["fileid1","fileid2"]}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodPost, "/drafts")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("req body: (-got +want):\n%s", diff)
		}
		_, _ = w.Write(getDraftJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	_, err := client.CreateDraft(context.Background(), DraftRequest{
		Subject: "Subject",
		From: []Participant{
			{
				Email: "from@example.org",
				Name:  "From Name",
			},
		},
		To: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		CC: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		BCC: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		ReplyTo: []Participant{
			{
				Email: "replyto@example.org",
				Name:  "ReplyTo Name",
			},
		},
		Body:    "body",
		FileIDs: []string{"fileid1", "fileid2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateDraft(t *testing.T) {
	accessToken := "accessToken"
	wantQuery := url.Values{}
	id := "8r5awu0esbg8ct3wg5rj5sifp"
	wantBody := []byte(strings.Join([]string{
		`{"subject":"Subject","from":[{"email":"from@example.org","name":`,
		`"from"}],"to":[{"email":"to@example.org","name":"to"}],"cc":[{"e`,
		`mail":"cc@example.org","name":"cc"}],"bcc":[{"email":"bcc@exampl`,
		`e.org","name":"bcc"}],"reply_to":[{"email":"replyto@example.org"`,
		`,"name":"replyto"}],"reply_to_message_id":"replytomessageid","bo`,
		`dy":"body","file_ids":["fileid1","fileid2"],"version":0}`,
	}, ""))
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodPut, "/drafts/"+id)
		assertQueryParams(t, r, wantQuery)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("req body: (-got +want):\n%s", diff)
		}
		_, _ = w.Write(getDraftJSON)
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	_, err := client.UpdateDraft(context.Background(), id, UpdateDraftRequest{
		Subject: String("Subject"),
		From: &[]Participant{
			{Name: "from", Email: "from@example.org"},
		},
		To: &[]Participant{
			{Name: "to", Email: "to@example.org"},
		},
		CC: &[]Participant{
			{Name: "cc", Email: "cc@example.org"},
		},
		BCC: &[]Participant{
			{Name: "bcc", Email: "bcc@example.org"},
		},
		ReplyTo: &[]Participant{
			{Name: "replyto", Email: "replyto@example.org"},
		},
		ReplyToMessageID: String("replytomessageid"),
		Body:             String("body"),
		FileIDs:          &[]string{"fileid1", "fileid2"},
		Version:          0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteDraft(t *testing.T) {
	accessToken := "accessToken"
	draftID := "draftID"
	wantBody := []byte(`{"version":1}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodDelete, "/drafts/"+draftID)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("req body: (-got +want):\n%s", diff)
		}
	}))
	defer ts.Close()

	client := NewClient("", "", withTestServer(ts), WithAccessToken(accessToken))
	err := client.DeleteDraft(context.Background(), draftID, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendDraft(t *testing.T) {
	accessToken := "accessToken"
	wantBody := []byte(`{"draft_id":"draftID","version":5}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodPost, "/send")

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
	_, err := client.SendDraft(context.Background(), "draftID", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSendDirectly(t *testing.T) {
	accessToken := "accessToken"
	wantBody := []byte(`{"subject":"Subject","from":[{"email":"from@example.org","name":"From Name"}],"to":[{"email":"to@example.org","name":"To Name"}],"cc":[{"email":"to@example.org","name":"To Name"}],"bcc":[{"email":"to@example.org","name":"To Name"}],"reply_to":[{"email":"replyto@example.org","name":"ReplyTo Name"}],"body":"body","file_ids":["fileid1","fileid2"]}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")
		assertMethodPath(t, r, http.MethodPost, "/send")

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
	_, err := client.SendDirectly(context.Background(), DraftRequest{
		Subject: "Subject",
		From: []Participant{
			{
				Email: "from@example.org",
				Name:  "From Name",
			},
		},
		To: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		CC: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		BCC: []Participant{
			{
				Email: "to@example.org",
				Name:  "To Name",
			},
		},
		ReplyTo: []Participant{
			{
				Email: "replyto@example.org",
				Name:  "ReplyTo Name",
			},
		},
		Body:    "body",
		FileIDs: []string{"fileid1", "fileid2"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

var getDraftJSON = []byte(`{
    "account_id": "43jf3n4e***",
    "bcc": [],
    "body": "Hello, how are you?",
    "cc": [],
    "date": 1559763005,
    "events": [],
    "files": [],
    "folder": {
	"display_name": "Draft",
	"id": "eeangfw9***",
	"name": "drafts"
    },
    "from": [
	{
	    "email": "nylastest***@yahoo.com",
	    "name": "John Doe"
	}
    ],
    "id": "43vfrmdu1***",
    "object": "draft",
    "reply_to": [],
    "reply_to_message_id": null,
    "snippet": "Hello, how are you?",
    "starred": false,
    "subject": "ugh?",
    "thread_id": "46wnzkxaz***",
    "to": [
	{
	    "email": "{{email}}",
	    "name": "{{name}}"
	}
    ],
    "unread": false,
    "version": 2
}`)

var draftsJSON = []byte(`[
    {
	"account_id": "43jf3n4***",
	"bcc": [],
	"body": "Hello, how are you?",
	"cc": [],
	"date": 1559763005,
	"events": [],
	"files": [],
	"folder": null,
	"from": [
	    {
		"email": "nylastest***@yahoo.com",
		"name": "John Doe"
	    }
	],
	"id": "43vfrmdu1***",
	"object": "draft",
	"reply_to": [],
	"reply_to_message_id": null,
	"snippet": "Hello, how are you?",
	"starred": false,
	"subject": "ugh?",
	"thread_id": "46wnzkxa***",
	"to": [
	    {
		"email": "{{email}}",
		"name": "{{name}}"
	    }
	],
	"unread": false,
	"version": 0
    },
    {
	"account_id": "43jf3n4***",
	"bcc": [],
	"body": "Hello, how are you?",
	"cc": [],
	"date": 1559762902,
	"events": [],
	"files": [],
	"folder": {
	    "display_name": "Draft",
	    "id": "eeangfw9vm***",
	    "name": "drafts"
	},
	"from": [
	    {
		"email": "nylastest***@yahoo.com",
		"name": "John Doe"
	    }
	],
	"id": "92c7gucghzh***",
	"object": "draft",
	"reply_to": [],
	"reply_to_message_id": "132i3h12u3hadw",
	"snippet": "Hello, how are you?",
	"starred": false,
	"subject": "Hello",
	"thread_id": "e48pmw615r2****",
	"to": [
	    {
		"email": "{{email}}",
		"name": "{{name}}"
	    }
	],
	"unread": false,
	"version": 1
    }
]`)
