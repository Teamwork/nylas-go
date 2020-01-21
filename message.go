package nylas

import (
	"context"
	"net/http"
)

// Message contains all the details of a single email message.
// See: https://docs.nylas.com/reference#messages
type Message struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`
	ThreadID  string `json:"thread_id"`

	From    []Participant `json:"from"`
	To      []Participant `json:"to"`
	CC      []Participant `json:"cc"`
	BCC     []Participant `json:"bcc"`
	ReplyTo []Participant `json:"reply_to"`

	// Only available in expanded view, see:
	// https://docs.nylas.com/reference#views
	Headers struct {
		InReplyTo  string   `json:"In-Reply-To"`
		MessageID  string   `json:"Message-Id"`
		References []string `json:"References"`
	} `json:"headers"`

	Subject string `json:"subject"`
	Date    int64  `json:"date"`
	Body    string `json:"body"`
	Snippet string `json:"snippet"`

	Events []interface{} `json:"events"`
	Files  []File        `json:"files"`
	Folder Folder        `json:"folder"`
	Labels []Label       `json:"label"`

	Starred bool `json:"starred"`
	Unread  bool `json:"unread"`
}

// Messages returns messages which match the filter specified by parameters.
// TODO: params
// See: https://docs.nylas.com/reference#messages-1
func (c *Client) Messages(ctx context.Context) ([]Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/messages", nil)
	if err != nil {
		return nil, err
	}

	var resp []Message
	return resp, c.do(req, &resp)
}
