package nylas

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// WebhookListener receives requests from a Nylas webhook.
// See: https://docs.nylas.com/reference#webhooks
type WebhookListener struct {
	clientSecret string
}

// NewWebhookListener returns a new WebhookListener..
func NewWebhookListener(clientSecret string) *WebhookListener {
	return &WebhookListener{clientSecret}
}

// Listen for webhooks on the given address.
// The callback will be called for each webhook and a non-nil error will respond
// with a 500 including the error message.
//
// Note: the callback is handled synchronously, so if you need to do slow work
// in response to a webhook return nil and handle it in another routine.
//
// See: https://docs.nylas.com/reference#receiving-notifications
func (l *WebhookListener) Listen(addr string, fn func(WebhookDelta) error) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			challenge := r.URL.Query().Get("challenge")
			if r.Method == http.MethodGet && challenge != "" {
				_, _ = io.WriteString(w, challenge)
				return
			}

			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		signature := r.Header.Get("X-Nylas-Signature")
		if err := checkSignature(l.clientSecret, signature, data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		resp := struct {
			Deltas []WebhookDelta `json:"deltas"`
		}{}
		if err := json.Unmarshal(data, &resp); err != nil {
			msg := fmt.Sprintf("unmarshal delta: %v", err)
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		for _, delta := range resp.Deltas {
			if err := fn(delta); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})
	return http.ListenAndServe(addr, mux)
}

// WebhookDelta represents a change in a users mailbox from a webhook request.
type WebhookDelta struct {
	Date       int    `json:"date"`
	Object     string `json:"object"`
	Type       string `json:"type"`
	ObjectData struct {
		ID          string `json:"id"`
		Object      string `json:"object"`
		AccountID   string `json:"account_id"`
		NamespaceID string `json:"namespace_id"`

		Attributes struct { // only included for message.created event
			ThreadID     string `json:"thread_id"`
			ReceivedDate int    `json:"received_date"`
		} `json:"attributes"`

		// used for tracking, see:
		// https://docs.nylas.com/reference#understanding-tracking-notifications
		Metadata map[string]interface{} `json:"metadata"`
	} `json:"object_data"`
}

func checkSignature(secret, signature string, body []byte) error {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write(body)
	if err != nil {
		return err
	}
	if signature != hex.EncodeToString(mac.Sum(nil)) {
		return errors.New("signature mismatch")
	}
	return nil
}
