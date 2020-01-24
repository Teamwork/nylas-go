package nylas

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
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
	Labels []Label       `json:"labels"`

	Starred bool `json:"starred"`
	Unread  bool `json:"unread"`
}

// MessagesOptions provides optional parameters to the Messages method.
type MessagesOptions struct {
	View   string `url:"view,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`
	// Return messages with a matching literal subject
	Subject string `url:"subject,omitempty"`
	// Return messages that have been sent or received from the list of
	// email addresses. A maximum of 25 emails may be specified
	AnyEmail []string `url:"any_email,comma,omitempty"`
	// Return  messages sent to this email address
	To string `url:"to,omitempty"`
	// Return  messages sent from this email address
	From string `url:"from,omitempty"`
	// Return  messages that were CC'd to this email address
	CC string `url:"cc,omitempty"`
	// Return messages that were BCC'd to this email address, likely sent
	// from the parent account.
	// (Most SMTP gateways remove BCC information.)
	BCC string `url:"bcc,omitempty"`
	// Return messages in a given folder, or with a given label.
	// This parameter supports the name, display_name, or id of a folder or
	// label.
	In      string `url:"in,omitempty"`
	Unread  *bool  `url:"unread,omitempty"`
	Starred *bool  `url:"starred,omitempty"`
	// Return messages belonging to a specific thread
	ThreadID string `url:"thread_id,omitempty"`
	Filename string `url:"filename,omitempty"`
	// Return messages received before this Unix-based timestamp.
	ReceivedBefore int64 `url:"received_before,omitempty"`
	// Return messages received after this Unix-based timestamp.
	ReceivedAfter int64 `url:"received_after,omitempty"`
	HasAttachment *bool `url:"has_attachment,omitempty"`
}

// Messages returns messages which match the filter specified by parameters.
// See: https://docs.nylas.com/reference#messages-1
func (c *Client) Messages(ctx context.Context, opts *MessagesOptions) ([]Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/messages", nil)
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

	var resp []Message
	return resp, c.do(req, &resp)
}

// MessagesCount returns the count of messages which match the filter specified by
// parameters.
// See: https://docs.nylas.com/reference#messages-1
func (c *Client) MessagesCount(ctx context.Context, opts *MessagesOptions) (int, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/messages", nil)
	if err != nil {
		return 0, err
	}

	if opts == nil {
		opts = &MessagesOptions{}
	}
	opts.View = ViewCount
	vs, err := query.Values(opts)
	if err != nil {
		return 0, err
	}
	appendQueryValues(req, vs)

	var resp countResponse
	return resp.Count, c.do(req, &resp)
}

// Message returns a message by id.
// See: https://docs.nylas.com/reference#messagesid
func (c *Client) Message(ctx context.Context, id string, expanded bool) (Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/messages/"+id, nil)
	if err != nil {
		return Message{}, err
	}

	if expanded {
		appendQueryValues(req, url.Values{"view": {ViewExpanded}})
	}

	var resp Message
	return resp, c.do(req, &resp)
}

// RawMessage returns the raw message in RFC-2822 format.
// See: https://docs.nylas.com/reference#raw-message-contents
func (c *Client) RawMessage(ctx context.Context, id string) ([]byte, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/messages/"+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "message/rfc822")

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode >= 299 {
		return nil, NewError(resp)
	}

	return ioutil.ReadAll(resp.Body)
}

// UpdateMessageRequest contains the request parameters required to update a
// message.
type UpdateMessageRequest struct {
	Unread  *bool `json:"unread,omitempty"`
	Starred *bool `json:"starred,omitempty"`
	// FolderID to move this message to.
	FolderID *string `json:"folder_id,omitempty"`
	// LabelIDs to overwrite any previous labels with, you must provide
	// existing labels such as sent/drafts.
	LabelIDs *[]string `json:"label_ids,omitempty"`
}

// UpdateMessage updates a message with the id.
// See: https://docs.nylas.com/reference#messagesid-1
func (c *Client) UpdateMessage(
	ctx context.Context, id string, updateReq UpdateMessageRequest,
) (Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodPut, "/messages/"+id, &updateReq)
	if err != nil {
		return Message{}, err
	}

	var resp Message
	return resp, c.do(req, &resp)
}
