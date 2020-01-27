package nylas

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const apiURL = "https://api.nylas.com"

// View constants for more info, see:
// https://docs.nylas.com/reference#views
const (
	ViewCount    = "count"
	ViewExpanded = "expanded"
	ViewIDs      = "ids"
)

// Label/Folder mailbox name cosntants, for more info see:
// https://docs.nylas.com/reference#get-labels
// https://docs.nylas.com/reference#get-folders
// https://tools.ietf.org/html/rfc6154
const (
	MailboxInbox     = "inbox"
	MailboxAll       = "all"
	MailboxTrash     = "trash"
	MailboxArchive   = "archive"
	MailboxDrafts    = "drafts"
	MailboxSent      = "sent"
	MailboxSpam      = "spam"
	MailboxImportant = "important"
)

// ErrAccessTokenNotSet is returned when Client methods are called that require
// an access token to be set.
var ErrAccessTokenNotSet = errors.New("access token not set on client")

// Client for working with the Nylas API.
type Client struct {
	clientID, clientSecret string
	accessToken            string

	baseURL string
	c       *http.Client
}

// Option sets an optional setting on the Client.
type Option func(*Client)

// NewClient returns a new client for working with the Nylas API.
func NewClient(clientID, clientSecret string, opts ...Option) *Client {
	client := &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		baseURL:      apiURL,
		c:            http.DefaultClient,
	}

	for _, opt := range opts {
		opt(client)
	}
	return client
}

// WithHTTPClient returns an Option to set the http.Client to be used.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.c = httpClient
	}
}

// WithBaseURL returns an Option to set the base URL to be used.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithAccessToken returns an option to set the access token to be used.
// This token is used for user mailbox specific methods.
func WithAccessToken(token string) Option {
	return func(c *Client) {
		c.accessToken = token
	}
}

// As returns a copy of the Client with the given access token set.
func (c *Client) As(accessToken string) *Client {
	as := *c
	WithAccessToken(accessToken)(c)
	return &as
}

func (c *Client) newUserRequest(
	ctx context.Context, method, endpoint string, body interface{},
) (*http.Request, error) {
	if c.accessToken == "" {
		return nil, ErrAccessTokenNotSet
	}

	req, err := c.newRequest(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.accessToken, "")
	return req, nil
}

func (c *Client) newAccountRequest(
	ctx context.Context, method, endpoint string, body interface{},
) (*http.Request, error) {
	req, err := c.newRequest(ctx, method, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.clientSecret, "")
	return req, nil
}

func (c *Client) newRequest(
	ctx context.Context, method, endpoint string, body interface{},
) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body: %v", err)
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
		req.Header.Add("Content-Type", "application/json; charset=utf")
	}
	return req, nil
}

func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode >= 299 {
		return NewError(resp)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}
	return nil
}

func appendQueryValues(req *http.Request, values url.Values) {
	q := req.URL.Query()
	for k, vs := range values {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()
}

type countResponse struct {
	Count int `json:"count"`
}

// Bool returns a pointer to the given bool value.
func Bool(v bool) *bool { return &v }

// String returns a pointer to the given string value.
func String(v string) *string { return &v }
