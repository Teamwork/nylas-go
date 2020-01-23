package main

import (
	"log"

	"github.com/teamwork/nylas-go"
)

const (
	clientSecret = "nylas client secret"
)

func main() {
	listener := nylas.NewWebhookListener(clientSecret)
	log.Fatal(listener.Listen(":8080", func(delta nylas.WebhookDelta) error {
		// Handle the webhook, if doing long running work do it in
		// another routine as the request will timeout
		log.Printf("%s\nDate: %v\nData: %+v", delta.Type, delta.Date, delta.ObjectData)

		return nil // non-nil error will be returned to Nylas as 500
	}))
}
