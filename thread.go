package nylas

import (
	"context"
	"net/http"
)

// Thread combines multiple messages from the same conversation into a single
// first-class object that is similar to what users expect from email clients.
// See: https://docs.nylas.com/reference#threads
type Thread struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	DraftIDs       []string `json:"draft_ids"`
	Folders        []Folder `json:"folders"`
	HasAttachments bool     `json:"has_attachments"`

	FirstMessageTimestamp        int64 `json:"first_message_timestamp"`
	LastMessageReceivedTimestamp int64 `json:"last_message_received_timestamp"`
	LastMessageSentTimestamp     int64 `json:"last_message_sent_timestamp"`
	LastMessageTimestamp         int64 `json:"last_message_timestamp"`

	MessageIDs []string `json:"message_ids"`
	// Only available in expanded view and the body will be missing, see:
	// https://docs.nylas.com/reference#views
	Messages []Message `json:"messages"`

	Participants []Participant `json:"participants"`
	Snippet      string        `json:"snippet"`
	Starred      bool          `json:"starred"`
	Subject      string        `json:"subject"`
	Unread       bool          `json:"unread"`
	Version      int           `json:"version"`
}

// Threads returns threads which match the filter specified by parameters.
// TODO: params
// See: https://docs.nylas.com/reference#get-threads
func (c *Client) Threads(ctx context.Context) ([]Thread, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/threads", nil)
	if err != nil {
		return nil, err
	}

	var resp []Thread
	return resp, c.do(req, &resp)
}
