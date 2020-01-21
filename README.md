# Nylas API v2 Go client

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/teamwork/nylas-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamwork/nylas-go)](https://goreportcard.com/report/github.com/teamwork/nylas-go)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)

Provides access to the [Nylas Platform](https://docs.nylas.com/reference) v2 REST API.

## Installation

```
go get github.com/teamwork/nylas-go
```

## Usage

```go
package main

import (
    "context"
    "log"
    "net/http"
    "time"

    nylas "github.com/teamwork/nylas-go"
)

const (
    clientID     = "..."
    clientSecret = "..."
    accessToken  = "..."
)

func main() {
    client := nylas.NewClient(clientID, clientSecret,
        nylas.WithHTTPClient(&http.Client{
            Timeout: 3 * time.Second,
        }),
        nylas.WithAccessToken(accessToken),
    )

    ctx := context.Background()
    deltaResp, err := client.Deltas(ctx, cursor, &nylas.DeltasOptions{
        View:         nylas.ViewExpanded,
        IncludeTypes: []string{"thread"},
    })
    if err != nil {
        log.Fatalf("get thread deltas: %v", err)
    }

    for _, d := range deltaResp.Deltas {
        log.Printf("%+v", d)
    }
}
```

## Status

The following features are implemented in the client:

### Authentication

#### Hosted Authentication

- [ ] GET	/oauth/authorize
- [ ] POST	/oauth/token
- [ ] POST	/oauth/revoke

#### Native Authentication

- [ ] POST	/connect/authorize
- [ ] POST	/connect/token

### Apps & Accounts

#### Accounts

- [ ] GET	/account

#### Account Management

- [ ] GET	/a/{client_id}/accounts
- [ ] GET	/a/{client_id}/accounts/{id}
- [ ] POST	/a/{client_id}/accounts/{id}/downgrade
- [ ] POST	/a/{client_id}/accounts/{id}/upgrade
- [ ] POST	/a/{client_id}/accounts/{id}/revoke-all
- [ ] GET	/a/{client_id}/ip_addresses
- [ ] POST	/a/{client_id}/accounts/{id}/token-info

#### Application Management

- [ ] GET	/a/{client_id}
- [ ] POST	/a/{client_id}

### Threads

- [x] GET	/threads
- [ ] GET	/threads/{id}
- [ ] PUT	/threads/{id}

### Messages

- [x] GET	/messages
- [ ] GET	/messages/{id}
- [ ] PUT	/messages/{id}
- [ ] GET	/messages/{id} (raw message content)

### Folders

- [ ] GET	/folders
- [ ] GET	/folders/{id}
- [ ] POST	/folders
- [ ] PUT	/folders/{id}
- [ ] DEL	/folders/{id}

### Labels

- [ ] GET	/labels
- [ ] GET	/labels/{id}
- [ ] POST	/labels
- [ ] PUT	/labels/{id}
- [ ] DEL	/labels/{id}

### Drafts

- [ ] GET	/drafts
- [ ] GET	/drafts/{id}
- [ ] POST	/drafts
- [ ] PUT	/drafts/{id}
- [ ] DEL	/drafts/{id}

### Sending

- [ ] POST	/send#drafts
- [x] POST	/send#directly
- [ ] POST	/send#raw

### Files

- [ ] GET	/files
- [ ] GET	/files/{id}
- [ ] GET	/files/{id}/download
- [ ] POST	/files
- [ ] DEL	/files/{id}

### Calendars

- [ ] GET	/calendars
- [ ] GET	/calendars/{id}
- [ ] POST	/calendars/free-busy


### Events

- [ ] GET	/events
- [ ] GET	/events/{id}
- [ ] POST	/events
- [ ] PUT	/events/{id}
- [ ] DEL	/events/{id}
- [ ] POST	/send-rsvp

### Room Resources

- [ ] GET	/resources

### Contacts

- [ ] GET	/contacts
- [ ] GET	/contacts/{id}
- [ ] POST	/contacts
- [ ] PUT	/contacts/{id}
- [ ] DEL	/contacts/{id}
- [ ] GET	/contacts/{id}/picture
- [ ] GET	/contacts/groups

### Search

- [ ] GET	/threads/search
- [ ] GET	/messages/search

### Webhooks

- [x] Listener client
- [ ] GET	/webhooks
- [ ] POST	/webhooks
- [ ] GET	/webhooks/{id}
- [ ] PUT	/webhooks/{id}
- [ ] DEL	/webhooks/{id}


### Deltas

- [x] POST	/delta/latest_cursor
- [x] GET	/delta
- [ ] GET	/delta/longpoll
- [x] GET	/delta/streaming
