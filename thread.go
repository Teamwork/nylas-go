package nylas

import (
	"context"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// Thread combines multiple messages from the same conversation into a single
// first-class object that is similar to what users expect from email clients.
// See: https://docs.nylas.com/reference#threads
type Thread struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	Folders        []Folder `json:"folders"`
	HasAttachments bool     `json:"has_attachments"`

	FirstMessageTimestamp        int64 `json:"first_message_timestamp"`
	LastMessageReceivedTimestamp int64 `json:"last_message_received_timestamp"`
	LastMessageSentTimestamp     int64 `json:"last_message_sent_timestamp"`
	LastMessageTimestamp         int64 `json:"last_message_timestamp"`

	MessageIDs []string `json:"message_ids"`
	DraftIDs   []string `json:"draft_ids"`

	// Only available in expanded view and the body will be missing, see:
	// https://docs.nylas.com/reference#views
	Messages []Message `json:"messages"`
	Drafts   []Message `json:"drafts"`

	Participants []Participant `json:"participants"`
	Labels       []Label       `json:"labels"`
	Snippet      string        `json:"snippet"`
	Starred      bool          `json:"starred"`
	Subject      string        `json:"subject"`
	Unread       bool          `json:"unread"`
	Version      int           `json:"version"`
}

// ThreadsOptions provides optional parameters to the Threads method.
type ThreadsOptions struct {
	View   string `url:"view,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
	// Return threads with a matching literal subject
	Subject string `url:"subject,omitempty"`
	// Return threads that have been sent or received from the list of email
	// addresses. A maximum of 25 emails may be specified
	AnyEmail []string `url:"any_email,comma,omitempty"`
	// Return threads containing messages sent to this email address
	To string `url:"to,omitempty"`
	// Return threads containing messages sent from this email address
	From string `url:"from,omitempty"`
	// Return threads containing messages that were CC'd to this email address
	CC string `url:"cc,omitempty"`
	// Return threads containing messages that were BCC'd to this email
	// address, likely sent from the parent account. (Most SMTP gateways
	// remove BCC information.)
	BCC string `url:"bcc,omitempty"`
	// Return threads in a given folder, or with a given label.
	// This parameter supports the name, display_name, or id of a folder or
	// label.
	In string `url:"in,omitempty"`
	// Return threads with one or more unread messages
	Unread   *bool  `url:"unread,omitempty"`
	Filename string `url:"filename,omitempty"`
	// Return threads whose most recent message was received before this
	// Unix-based timestamp.
	LastMessageBefore int64 `url:"last_message_before,omitempty"`
	// Return threads whose most recent message was received after this
	// Unix-based timestamp.
	LastMessageAfter int64 `url:"last_message_after,omitempty"`
	// Return threads whose first message was received before this
	// Unix-based timestamp.
	StartedBefore int64 `url:"started_before,omitempty"`
	// Return threads whose first message was received after this
	// Unix-based timestamp.
	StartedAfter int64 `url:"started_after,omitempty"`
}

// Threads returns threads which match the filter specified by parameters.
// See: https://docs.nylas.com/reference#get-threads
func (c *Client) Threads(ctx context.Context, opts *ThreadsOptions) ([]Thread, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/threads", nil)
	if err != nil {
		return nil, err
	}

	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return nil, err
		}
		appendQueryValues(req, vs)
	}

	var resp []Thread
	return resp, c.do(req, &resp)
}

// ThreadsCount returns the count of threads which match the filter specified by
// parameters.
// See: https://docs.nylas.com/reference#get-threads
func (c *Client) ThreadsCount(ctx context.Context, opts *ThreadsOptions) (int, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/threads", nil)
	if err != nil {
		return 0, err
	}

	if opts == nil {
		opts = &ThreadsOptions{}
	}
	vs, err := query.Values(opts)
	if err != nil {
		return 0, err
	}
	vs.Set("view", ViewCount)
	appendQueryValues(req, vs)

	var resp countResponse
	return resp.Count, c.do(req, &resp)
}

// Thread returns a thread by id.
// See: https://docs.nylas.com/reference#threadsid
func (c *Client) Thread(ctx context.Context, id string, expanded bool) (Thread, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/threads/"+id, nil)
	if err != nil {
		return Thread{}, err
	}

	if expanded {
		appendQueryValues(req, url.Values{"view": {ViewExpanded}})
	}

	var resp Thread
	return resp, c.do(req, &resp)
}

// UpdateThreadRequest contains the request parameters required to update a
// thread.
type UpdateThreadRequest struct {
	Unread  *bool `json:"unread,omitempty"`
	Starred *bool `json:"starred,omitempty"`
	// FolderID to move this thread to.
	FolderID *string `json:"folder_id,omitempty"`
	// LabelIDs to overwrite any previous labels with, you must provide
	// existing labels such as sent/drafts.
	LabelIDs *[]string `json:"label_ids,omitempty"`
}

// UpdateThread updates a thread with the id.
// See: https://docs.nylas.com/reference#threadsid-1
func (c *Client) UpdateThread(
	ctx context.Context, id string, updateReq UpdateThreadRequest,
) (Thread, error) {
	req, err := c.newUserRequest(ctx, http.MethodPut, "/threads/"+id, &updateReq)
	if err != nil {
		return Thread{}, err
	}

	var resp Thread
	return resp, c.do(req, &resp)
}
