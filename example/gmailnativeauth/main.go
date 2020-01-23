package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/teamwork/nylas-go"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	nylasClientID      = "nylas client id"
	nylasClientSecret  = "nylas client secret"
	googleClientID     = "google client id"
	googleClientSecret = "google client secret"
	name               = "your name"
	emailAddress       = "youremail@gmail.com"
)

func main() {
	conf := &oauth2.Config{
		ClientID:     googleClientID,
		ClientSecret: googleClientSecret,
		RedirectURL:  "http://localhost:8080",
		Scopes: []string{
			"openid",
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://mail.google.com/",
			"https://www.googleapis.com/auth/calendar",
			"https://www.googleapis.com/auth/contacts",
		},
		Endpoint: google.Endpoint,
	}

	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.ApprovalForce)
	log.Printf("Visit the URL for the auth dialog:\n%v", url)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			defer cancel()

			code := r.URL.Query().Get("code")
			tok, err := conf.Exchange(ctx, code)
			if err != nil {
				log.Printf("Exchange: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Print("Connecting account to Nylas, this might take a while")
			client := nylas.NewClient(nylasClientID, nylasClientSecret)
			acc, err := client.ConnectAccount(ctx, nylas.AuthorizeRequest{
				Name:         name,
				EmailAddress: emailAddress,
				Settings: nylas.GmailAuthorizeSettings{
					GoogleClientID:     googleClientID,
					GoogleClientSecret: googleClientSecret,
					GoogleRefreshToken: tok.RefreshToken,
				},
				Scopes: []string{"email", "calendar"},
			})
			if err != nil {
				log.Printf("ConnectAccount: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Printf("Success, account: %+v", acc)
			_ = json.NewEncoder(w).Encode(acc)
		})
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	<-ctx.Done()
}
