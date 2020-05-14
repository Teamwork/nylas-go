# Nylas API v2 Go client

[![Build Status](https://travis-ci.com/Teamwork/nylas-go.svg?token=o9pKscgKamFB17WDSzzf&branch=master)](https://travis-ci.com/Teamwork/nylas-go)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://pkg.go.dev/github.com/teamwork/nylas-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/teamwork/nylas-go)](https://goreportcard.com/report/github.com/teamwork/nylas-go)
[![License MIT](https://img.shields.io/badge/license-MIT-lightgrey.svg?style=flat)](LICENSE)

Provides access to the [Nylas Platform](https://docs.nylas.com/reference) v2 REST API.

This library does not currently cover all endpoints and will be implemented as and when they are needed by our internal projects. Check the [GoDoc](https://godoc.org/github.com/teamwork/nylas-go) or [Status](#status) for current implementation status.

This is not an official SDK for the Nylas Platform, for the official SDKs visit [Nylas SDKs](https://www.nylas.com/sdks/).

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
        nylas.WithAccessToken(accessToken), // not required for account management endpoints
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

## Examples

See the [example directory](example).

## Contributing

We would like to make this library feature complete with the offical SDK projects and contributions are welcome.

Following the existing code style and conventions and submit a PR.

## Status

The following features are implemented in the client:

### Authentication

#### Hosted Authentication

- [ ] GET	/oauth/authorize
- [ ] POST	/oauth/token
- [ ] POST	/oauth/revoke

#### Native Authentication

- [x] POST	/connect/authorize
- [x] POST	/connect/token

### Apps & Accounts

#### Accounts

- [x] GET	/account

#### Account Management

- [x] GET	/a/{client_id}/accounts
- [ ] GET	/a/{client_id}/accounts/{id}
- [x] POST	/a/{client_id}/accounts/{id}/downgrade
- [x] POST	/a/{client_id}/accounts/{id}/upgrade
- [x] POST	/a/{client_id}/accounts/{id}/revoke-all
- [ ] GET	/a/{client_id}/ip_addresses
- [ ] POST	/a/{client_id}/accounts/{id}/token-info

#### Application Management

- [ ] GET	/a/{client_id}
- [ ] POST	/a/{client_id}

### Threads

- [x] GET	/threads
- [x] GET	/threads/{id}
- [x] PUT	/threads/{id}

### Messages

- [x] GET	/messages
- [x] GET	/messages/{id}
- [x] PUT	/messages/{id}
- [x] GET	/messages/{id} (raw message content)

### Folders

- [x] GET	/folders
- [ ] GET	/folders/{id}
- [ ] POST	/folders
- [ ] PUT	/folders/{id}
- [ ] DEL	/folders/{id}

### Labels

- [x] GET	/labels
- [ ] GET	/labels/{id}
- [ ] POST	/labels
- [ ] PUT	/labels/{id}
- [ ] DEL	/labels/{id}

### Drafts

- [x] GET	/drafts
- [x] GET	/drafts/{id}
- [x] POST	/drafts
- [x] PUT	/drafts/{id}
- [x] DEL	/drafts/{id}

### Sending

- [x] POST	/send#drafts
- [x] POST	/send#directly
- [ ] POST	/send#raw

### Files

- [ ] GET	/files
- [x] GET	/files/{id}
- [x] GET	/files/{id}/download
- [x] POST	/files
- [x] DEL	/files/{id}

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
