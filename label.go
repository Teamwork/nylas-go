package nylas

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Label represents a label in the Nylas system.
type Label struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	// Localized name of the label
	DisplayName string `json:"display_name"`

	// Standard categories type, based on RFC-6154, can be one of the
	// Mailbox* constants, e.g MailboxInbox or empty if user created.
	// See: https://tools.ietf.org/html/rfc6154
	Name string `json:"name"`
}

// LabelsOptions provides optional parameters to the Labels method.
type LabelsOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

// Labels returns labels which match the filter specified by parameters.
// See: https://docs.nylas.com/reference#get-labels
func (c *Client) Labels(ctx context.Context, opts *LabelsOptions) ([]Label, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/labels", nil)
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

	var resp []Label
	return resp, c.do(req, &resp)
}

// LabelsCount returns the count of labels.
// See: https://docs.nylas.com/reference#get-labels
func (c *Client) LabelsCount(ctx context.Context) (int, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/labels?view=count", nil)
	if err != nil {
		return 0, err
	}

	var resp countResponse
	return resp.Count, c.do(req, &resp)
}
