# ynab.go Project Overview

## Purpose
Unofficial Go client library for the YNAB (You Need A Budget) API.
Covers 100% of the YNAB API v1 endpoints.

## Tech Stack
- Go 1.24
- github.com/stretchr/testify for testing
- gopkg.in/jarcoal/httpmock.v1 for HTTP mocking
- golangci-lint for linting
- gofumpt for formatting

## Key Packages
- Root package `ynab`: Main client (client.go)
- `api/`: Core HTTP client, error types, interfaces, rate limiting, token provider
- `api/account/`, `api/category/`, `api/transaction/`, etc.: Service packages per resource
- `oauth/`: OAuth 2.0 support (token.go, flow.go, config.go, storage.go, client.go)

## Commands
- `go test -race -short ./...` - run tests
- `make lint` - run linting
- `make coverage` - generate coverage report
- `make go-check` - full validation

## Architecture
- `ClientServicer` interface: main entry point
- Services instantiated from `client` struct
- `TokenProvider` interface: abstracts static vs OAuth tokens
- `RateLimitTracker`: rolling window rate limiting
- `HTTPClient`: handles request preparation and response parsing
