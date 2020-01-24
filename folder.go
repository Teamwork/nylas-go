package nylas

import (
	"context"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Folder represents a folder in the Nylas system.
type Folder struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	// Localized name of the folder
	DisplayName string `json:"display_name"`

	// Standard categories type, based on RFC-6154, can be one of the
	// Mailbox* constants, e.g MailboxInbox or empty if user created.
	// See: https://tools.ietf.org/html/rfc6154
	Name string `json:"name"`
}

// FoldersOptions provides optional parameters to the Folders method.
type FoldersOptions struct {
	Limit  int `url:"limit,omitempty"`
	Offset int `url:"offset,omitempty"`
}

// Folders returns folders which match the filter specified by parameters.
// See: https://docs.nylas.com/reference#get-folders
func (c *Client) Folders(ctx context.Context, opts *FoldersOptions) ([]Folder, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/folders", nil)
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

	var resp []Folder
	return resp, c.do(req, &resp)
}

// FoldersCount returns the count of folders.
// See: https://docs.nylas.com/reference#get-folders
func (c *Client) FoldersCount(ctx context.Context) (int, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/folders?view=count", nil)
	if err != nil {
		return 0, err
	}

	var resp countResponse
	return resp.Count, c.do(req, &resp)
}
