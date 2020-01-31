package nylas

import (
	"context"
	"fmt"
	"net/http"
)

// BillingState constants, for more info see:
// https://docs.nylas.com/reference#aclient_idaccounts
const (
	BillingStateCancelled = "cancelled"
	BillingStatePaid      = "paid"
	BillingStateDeleted   = "deleted"
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

// ManagementAccount contains the details of an account and is used when working
// with the Account Management endpoints.
// See: https://docs.nylas.com/reference#account-management
type ManagementAccount struct {
	ID        string `json:"id"`
	AccountID string `json:"account_id"`

	BillingState string `json:"billing_state"`
	Email        string `json:"email"`
	Provider     string `json:"provider"`
	SyncState    string `json:"sync_state"`
	Trial        bool   `json:"trial"`
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

// Accounts returns the account information for all accounts.
// See: https://docs.nylas.com/reference#aclient_idaccounts
func (c *Client) Accounts(ctx context.Context) ([]ManagementAccount, error) {
	endpoint := fmt.Sprintf("/a/%s/accounts", c.clientID)
	req, err := c.newAccountRequest(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	var resp []ManagementAccount
	return resp, c.do(req, &resp)
}

// CancelAccount cancels a paid account.
// See: https://docs.nylas.com/reference#cancel-an-account
func (c *Client) CancelAccount(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/a/%s/accounts/%s/downgrade", c.clientID, id)
	req, err := c.newAccountRequest(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// ReactivateAccount re-enables a cancelled account to make it activate again.
// See: https://docs.nylas.com/reference#re-activate-an-account.
func (c *Client) ReactivateAccount(ctx context.Context, id string) error {
	endpoint := fmt.Sprintf("/a/%s/accounts/%s/upgrade", c.clientID, id)
	req, err := c.newAccountRequest(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// RevokeAccountTokens revokes all account tokens, optionally excluding one.
// See: https://docs.nylas.com/reference#revoke-all
func (c *Client) RevokeAccountTokens(ctx context.Context, id string, keepToken *string) error {
	endpoint := fmt.Sprintf("/a/%s/accounts/%s/revoke-all", c.clientID, id)
	var body map[string]interface{}
	if keepToken != nil {
		body = map[string]interface{}{
			"keep_access_token": *keepToken,
		}
	}
	req, err := c.newAccountRequest(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
