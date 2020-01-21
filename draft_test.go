package nylas

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSendDirectly(t *testing.T) {
	accessToken := "accessToken"
	wantBody := []byte(`{"subject":"Subject","from":[{"email":"from@example.org","name":"From Name"}],"to":[{"email":"to@example.org","name":"To Name"}],"cc":[{"email":"to@example.org","name":"To Name"}],"bcc":[{"email":"to@example.org","name":"To Name"}],"reply_to":[{"email":"replyto@example.org","name":"ReplyTo Name"}],"body":"body","file_ids":["fileid1","fileid2"]}`)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertBasicAuth(t, r, accessToken, "")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if diff := cmp.Diff(body, wantBody); diff != "" {
			t.Errorf("Message: (-got +want):\n%s", diff)
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
