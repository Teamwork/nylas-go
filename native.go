package nylas

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
)

// AuthorizeSettings provides settings for a native authentication connect
// request and should JSON marshal into the desired object.
// See: https://docs.nylas.com/reference#section-provider-specific-settings
type AuthorizeSettings interface {
	// Provider returns the provider value to be used in a connect request.
	Provider() string
}

// AuthorizeRequest used to start the process of connecting an account to Nylas.
// See: https://docs.nylas.com/reference#connectauthorize
type AuthorizeRequest struct {
	clientID     string
	Name         string
	EmailAddress string
	Settings     AuthorizeSettings
	Scopes       []string
}

// MarshalJSON implements the json.Marshaler interface.
func (r AuthorizeRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"client_id":     r.clientID,
		"name":          r.Name,
		"email_address": r.EmailAddress,
		"provider":      r.Settings.Provider(),
		"settings":      r.Settings,
		"scopes":        strings.Join(r.Scopes, ","),
	})
}

// ConnectAccount to Nylas with Native Authentication.
// See: https://docs.nylas.com/docs/native-authentication
func (c *Client) ConnectAccount(ctx context.Context, authReq AuthorizeRequest) (Account, error) {
	code, err := c.connectAuthorize(ctx, authReq)
	if err != nil {
		return Account{}, err
	}
	return c.connectExchangeCode(ctx, code)
}

func (c *Client) connectAuthorize(ctx context.Context, authReq AuthorizeRequest) (string, error) {
	authReq.clientID = c.clientID
	req, err := c.newRequest(ctx, http.MethodPost, "/connect/authorize", &authReq)
	if err != nil {
		return "", err
	}

	connectResp := struct {
		Code string `json:"code"`
	}{}
	return connectResp.Code, c.do(req, &connectResp)
}

func (c *Client) connectExchangeCode(ctx context.Context, authorizeCode string) (Account, error) {
	req, err := c.newRequest(ctx, http.MethodPost, "/connect/token", &map[string]interface{}{
		"client_id":     c.clientID,
		"client_secret": c.clientSecret,
		"code":          authorizeCode,
	})
	if err != nil {
		return Account{}, err
	}

	var resp Account
	return resp, c.do(req, &resp)
}

// GmailAuthorizeSettings implements AuthorizeSettings.
type GmailAuthorizeSettings struct {
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
	GoogleRefreshToken string `json:"google_refresh_token"`
}

// Provider returns the provider value to be used in a connect request.
func (GmailAuthorizeSettings) Provider() string { return "gmail" }

// IMAPAuthorizeSettings implements AuthorizeSettings.
type IMAPAuthorizeSettings struct {
	IMAPHost     string `json:"imap_host"`
	IMAPPort     int    `json:"imap_port"`
	IMAPUsername string `json:"imap_username"`
	IMAPPassword string `json:"imap_password"`
	SMTPHost     string `json:"smtp_host"`
	SMTPPort     int    `json:"smtp_port"`
	SMTPUsername string `json:"smtp_username"`
	SMTPPassword string `json:"smtp_password"`
	SSLRequired  bool   `json:"ssl_required"`
}

// Provider returns the provider value to be used in a connect request.
func (IMAPAuthorizeSettings) Provider() string { return "imap" }

// ExchangeAuthorizeSettings implements AuthorizeSettings.
type ExchangeAuthorizeSettings struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Provider returns the provider value to be used in a connect request.
func (ExchangeAuthorizeSettings) Provider() string { return "exchange" }

// Office365AuthorizeSettings implements AuthorizeSettings.
type Office365AuthorizeSettings struct {
	MicrosoftClientSecret string `json:"microsoft_client_secret"`
	MicrosoftRefreshToken string `json:"microsoft_refresh_token"`
	RedirectURI           string `json:"redirect_uri"`
}

// Provider returns the provider value to be used in a connect request.
func (Office365AuthorizeSettings) Provider() string { return "office365" }

// OutlookAuthorizeSettings implements AuthorizeSettings.
type OutlookAuthorizeSettings struct {
	Password string `json:"password"`
}

// Provider returns the provider value to be used in a connect request.
func (OutlookAuthorizeSettings) Provider() string { return "outlook" }
