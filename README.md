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
