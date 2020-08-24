package nylas

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Draft is a special kind of message which has not been sent, and therefore
// it's body contents and recipients are still mutable.
// See: https://docs.nylas.com/reference#drafts
type Draft struct {
	Message
	ReplyToMessageID string `json:"reply_to_message_id"`
	Version          int    `json:"version"`
}

// DraftRequest contains the request parameters required to create a draft or
// send it directly.
// See: https://docs.nylas.com/reference#drafts
type DraftRequest struct {
	Subject          string        `json:"subject"`
	From             []Participant `json:"from"`
	To               []Participant `json:"to"`
	CC               []Participant `json:"cc"`
	BCC              []Participant `json:"bcc"`
	ReplyTo          []Participant `json:"reply_to"`
	ReplyToMessageID string        `json:"reply_to_message_id,omitempty"`
	Body             string        `json:"body"`
	FileIDs          []string      `json:"file_ids"`
}

// UpdateDraftRequest contains the request parameters requiredto update a draft.
//
// Version is required to specify the version of the draft you wish to update,
// other fields are optional and will overwrite previous values if given.
type UpdateDraftRequest struct {
	Subject          *string        `json:"subject,omitempty"`
	From             *[]Participant `json:"from,omitempty"`
	To               *[]Participant `json:"to,omitempty"`
	CC               *[]Participant `json:"cc,omitempty"`
	BCC              *[]Participant `json:"bcc,omitempty"`
	ReplyTo          *[]Participant `json:"reply_to,omitempty"`
	ReplyToMessageID *string        `json:"reply_to_message_id,omitempty"`
	Body             *string        `json:"body,omitempty"`
	FileIDs          *[]string      `json:"file_ids,omitempty"`

	Version int `json:"version"`
}

// DraftsOptions provides optional parameters to the Drafts method.
type DraftsOptions struct {
	View   string `url:"view,omitempty"`
	Limit  int    `url:"limit,omitempty"`
	Offset int    `url:"offset,omitempty"`

	// Return messages belonging to a specific thread
	ThreadID string `url:"thread_id,omitempty"`

	// Return drafts that have been sent or received from the list of
	// email addresses. A maximum of 25 emails may be specified
	AnyEmail []string `url:"any_email,comma,omitempty"`
}

// Drafts returns drafts which match the filter specified by parameters.
// See: https://docs.nylas.com/reference#get-drafts
func (c *Client) Drafts(ctx context.Context, opts *DraftsOptions) ([]Draft, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/drafts", nil)
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

	var resp []Draft
	return resp, c.do(req, &resp)
}

// DraftsCount returns the count of drafts which match the filter specified by
// parameters.
// See: https://docs.nylas.com/reference#get-drafts
func (c *Client) DraftsCount(ctx context.Context, opts *DraftsOptions) (int, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/drafts", nil)
	if err != nil {
		return 0, err
	}

	if opts == nil {
		opts = &DraftsOptions{}
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

// Draft returns a draft by id.
// See: https://docs.nylas.com/reference#get-draft
func (c *Client) Draft(ctx context.Context, id string) (Draft, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/drafts/"+id, nil)
	if err != nil {
		return Draft{}, err
	}

	var resp Draft
	return resp, c.do(req, &resp)
}

// CreateDraft creates a new draft.
// See: https://docs.nylas.com/reference#post-draft
func (c *Client) CreateDraft(ctx context.Context, draftReq DraftRequest) (Draft, error) {
	req, err := c.newUserRequest(ctx, http.MethodPost, "/drafts", &draftReq)
	if err != nil {
		return Draft{}, err
	}

	var resp Draft
	return resp, c.do(req, &resp)
}

// UpdateDraft updates a draft with the id.
//
// Updating a draft returns a draft with the same ID but different Version.
// When submitting subsequent send or save actions, you must use this new version.
// See: https://docs.nylas.com/reference#put-draft
func (c *Client) UpdateDraft(
	ctx context.Context, id string, updateReq UpdateDraftRequest,
) (Draft, error) {
	req, err := c.newUserRequest(ctx, http.MethodPut, "/drafts/"+id, &updateReq)
	if err != nil {
		return Draft{}, err
	}

	var resp Draft
	return resp, c.do(req, &resp)
}

// DeleteDraft deletes draft matching the id, version must be the latest version
// of the draft.
// See: https://docs.nylas.com/reference#draftsid
func (c *Client) DeleteDraft(ctx context.Context, id string, version int) error {
	endpoint := fmt.Sprintf("/drafts/%s", id)
	req, err := c.newUserRequest(ctx, http.MethodDelete, endpoint, &map[string]interface{}{
		"version": version,
	})
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// SendDraft sends an existing drafted with the given id and version.
// Version must be the most recent version of the draft or the request will fail.
// See: https://docs.nylas.com/reference#sending-drafts
func (c *Client) SendDraft(ctx context.Context, id string, version int) (Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodPost, "/send", &map[string]interface{}{
		"draft_id": id,
		"version":  version,
	})
	if err != nil {
		return Message{}, err
	}

	var resp Message
	return resp, c.do(req, &resp)
}

// SendDirectly a message without creating a draft first.
// See: https://docs.nylas.com/reference#sending-directly
func (c *Client) SendDirectly(ctx context.Context, draftRequest DraftRequest) (Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodPost, "/send", &draftRequest)
	if err != nil {
		return Message{}, err
	}

	var resp Message
	return resp, c.do(req, &resp)
}
