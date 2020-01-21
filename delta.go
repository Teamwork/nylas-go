package nylas

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
)

// DeltaResponse contains the response of a Delta API request.
type DeltaResponse struct {
	CursorStart string  `json:"cursor_start"`
	CursorEnd   string  `json:"cursor_end"`
	Deltas      []Delta `json:"deltas"`
}

// Delta represents a change in the Nylas system.
// See: https://docs.nylas.com/reference#deltas
type Delta struct {
	ID     string `json:"id"`
	Object string `json:"object"`

	Event      string          `json:"event"`
	Cursor     string          `json:"cursor"`
	Attributes json.RawMessage `json:"attributes"`
}

// Message unmarshals the receivers Attributes field into a Message.
func (d Delta) Message() (Message, error) {
	var message Message
	return message, json.Unmarshal(d.Attributes, &message)
}

// Thread unmarshals the receivers Attributes field into a Thread.
func (d Delta) Thread() (Thread, error) {
	var thread Thread
	return thread, json.Unmarshal(d.Attributes, &thread)
}

// LatestDeltaCursor returns latest delta cursor for a users mailbox.
// See: https://docs.nylas.com/reference#obtaining-a-delta-cursor
func (c *Client) LatestDeltaCursor(ctx context.Context) (string, error) {
	req, err := c.newUserRequest(ctx, http.MethodPost, "/delta/latest_cursor", nil)
	if err != nil {
		return "", err
	}

	latestCursor := struct {
		Cursor string `json:"cursor"`
	}{}
	return latestCursor.Cursor, c.do(req, &latestCursor)
}

// DeltasOptions provides optional parameters to the Deltas method.
type DeltasOptions struct {
	IncludeTypes []string `url:"include_types,comma,omitempty"`
	ExcludeTypes []string `url:"exclude_types,comma,omitempty"`
	View         string   `url:"view,omitempty"`
}

// Deltas requests a set of changes starting at cursor for a users mailbox.
//
// Note: this may not return all the changes that have happened since the start
// of the cursor and so you should keep requesting using DeltaResponse.CursorEnd
// until a response is given with CursorStart equal to CursorEnd.
//
// See: https://docs.nylas.com/reference#requesting-a-set-of-deltas
func (c *Client) Deltas(
	ctx context.Context, cursor string, opts *DeltasOptions,
) (DeltaResponse, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/delta", nil)
	if err != nil {
		return DeltaResponse{}, err
	}

	appendQueryValues(req, url.Values{"cursor": {cursor}})
	if opts != nil {
		vs, err := query.Values(opts)
		if err != nil {
			return DeltaResponse{}, err
		}
		appendQueryValues(req, vs)
	}

	var resp DeltaResponse
	return resp, c.do(req, &resp)
}

// StreamDeltas streams deltas for a users mailbox with a long lived connection
// calling the provided function with each delta.
//
// This method will block until the context is cancelled or an error occurs.
// Ensure you set a http.Client with appropriate timeout settings, e.g:
//   &http.Client{
//	Transport: &http.Transport{
//		Dial: (&net.Dialer{
//			Timeout: 5 * time.Second,
//		}).Dial,
//		ResponseHeaderTimeout: 10 * time.Second,
//		TLSHandshakeTimeout:   5 * time.Second,
//	},
//   }
//
// See: https://docs.nylas.com/reference#streaming-delta-updates
func (c *Client) StreamDeltas(ctx context.Context, cursor string, fn func(Delta)) error {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/delta/streaming", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("cursor", cursor)
	req.URL.RawQuery = q.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close() // nolint: errcheck

	if resp.StatusCode > 299 {
		return NewError(resp)
	}

	reader := bufio.NewReader(resp.Body)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			line, err := reader.ReadBytes('\n')
			if err != nil {
				return err
			}
			if len(line) == 1 { // keep alive
				continue
			}

			var delta Delta
			if err := json.Unmarshal(line, &delta); err != nil {
				return fmt.Errorf("unmarshal delta: %q: %v", line, err)
			}
			fn(delta)
		}
	}
}
