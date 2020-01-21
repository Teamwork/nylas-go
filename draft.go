package nylas

import (
	"context"
	"net/http"
)

// Draft is a special kind of message which has not been sent, and therefore
// it's body contents and recipients are still mutable.
// See: https://docs.nylas.com/reference#drafts
type Draft Message

// DraftRequest contains the request parameters required to create a draft or
// send it directly.
// See: https://docs.nylas.com/reference#drafts
type DraftRequest struct {
	Subject string        `json:"subject"`
	From    []Participant `json:"from"`
	To      []Participant `json:"to"`
	CC      []Participant `json:"cc"`
	BCC     []Participant `json:"bcc"`
	ReplyTo []Participant `json:"reply_to"`
	Body    string        `json:"body"`
	FileIDs []string      `json:"file_ids"`
}

// SendDirectly a message without creating a draft first.
// See: https://docs.nylas.com/reference#sending-directly
func (c *Client) SendDirectly(ctx context.Context, draftRequest DraftRequest) (Message, error) {
	req, err := c.newUserRequest(ctx, http.MethodPost, "/send#directly", &draftRequest)
	if err != nil {
		return Message{}, err
	}

	var resp Message
	return resp, c.do(req, &resp)
}
