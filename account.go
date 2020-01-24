package nylas

import (
	"context"
	"net/http"
)

// Account contains the details of an account which corresponds to an email
// address, mailbox, and optionally a calendar.
type Account struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	AccountID string `json:"account_id"`

	Name             string `json:"name"`
	EmailAddress     string `json:"email_address"`
	Provider         string `json:"provider"`
	OrganizationUnit string `json:"organization_unit"`
	SyncState        string `json:"sync_state"`
	LinkedAt         int    `json:"linked_at"`

	// Only populated after a call to ConnectAccount
	AccessToken  string `json:"access_token"`
	BillingState string `json:"billing_state"`
}

// Account returns the account information for the user the client is
// authenticated as.
// See: https://docs.nylas.com/reference#account
func (c *Client) Account(ctx context.Context) (Account, error) {
	req, err := c.newUserRequest(ctx, http.MethodGet, "/account", nil)
	if err != nil {
		return Account{}, err
	}

	var resp Account
	return resp, c.do(req, &resp)
}
