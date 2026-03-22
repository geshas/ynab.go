package oauth

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type blockingTokenTransport struct {
	started chan struct{}
	release chan struct{}
	body    string
}

func (t *blockingTokenTransport) RoundTrip(*http.Request) (*http.Response, error) {
	close(t.started)
	<-t.release

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(t.body)),
	}, nil
}

func TestTokenManagerGetTokenDoesNotBlockReadersDuringRefresh(t *testing.T) {
	config := NewOAuthConfig(Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURI:  "https://example.com/callback",
	})

	storage := NewMemoryStorage()
	require.NoError(t, storage.SaveToken(&Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(-time.Minute),
	}))

	transport := &blockingTokenTransport{
		started: make(chan struct{}),
		release: make(chan struct{}),
		body:    `{"access_token":"fresh-access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":7200}`,
	}

	tm := NewTokenManager(config, storage).WithHTTPClient(&http.Client{Transport: transport})

	getDone := make(chan error, 1)
	go func() {
		_, err := tm.GetToken(context.Background())
		getDone <- err
	}()

	select {
	case <-transport.started:
	case <-time.After(time.Second):
		t.Fatal("expected refresh request to start")
	}

	authDone := make(chan bool, 1)
	go func() {
		authDone <- tm.IsAuthenticated()
	}()

	select {
	case authenticated := <-authDone:
		assert.False(t, authenticated)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("IsAuthenticated blocked while refresh was in progress")
	}

	close(transport.release)

	select {
	case err := <-getDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("GetToken did not finish after refresh response was released")
	}
}

func TestTokenManagerGetTokenRefreshCallbackCanReadUpdatedToken(t *testing.T) {
	config := NewOAuthConfig(Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURI:  "https://example.com/callback",
	})

	storage := NewMemoryStorage()
	require.NoError(t, storage.SaveToken(&Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(-time.Minute),
	}))

	tm := NewTokenManager(config, storage).WithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(
					`{"access_token":"fresh-access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":7200}`,
				)),
			}, nil
		}),
	})

	callbackDone := make(chan string, 1)
	tm.WithTokenRefreshCallback(func(*Token) {
		accessToken, err := tm.GetAccessToken(context.Background())
		if err != nil {
			callbackDone <- "error: " + err.Error()
			return
		}

		callbackDone <- accessToken
	})

	_, err := tm.GetToken(context.Background())
	require.NoError(t, err)

	select {
	case accessToken := <-callbackDone:
		assert.Equal(t, "fresh-access-token", accessToken)
	case <-time.After(time.Second):
		t.Fatal("refresh callback did not complete")
	}
}

func TestTokenManagerRefreshTokenSerializesConcurrentRefreshes(t *testing.T) {
	config := NewOAuthConfig(Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURI:  "https://example.com/callback",
	})

	storage := NewMemoryStorage()
	require.NoError(t, storage.SaveToken(&Token{
		AccessToken:  "expired-access-token",
		RefreshToken: "refresh-token",
		ExpiresAt:    time.Now().Add(-time.Minute),
	}))

	transport := &countingBlockingTransport{
		started:       make(chan struct{}),
		release:       make(chan struct{}),
		secondAttempt: make(chan struct{}),
		body:          `{"access_token":"fresh-access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":7200}`,
	}

	tm := NewTokenManager(config, storage).WithHTTPClient(&http.Client{Transport: transport})

	getDone := make(chan error, 1)
	go func() {
		_, err := tm.GetToken(context.Background())
		getDone <- err
	}()

	select {
	case <-transport.started:
	case <-time.After(time.Second):
		t.Fatal("expected refresh request to start")
	}

	refreshDone := make(chan error, 1)
	refreshWaiting := make(chan struct{})
	tm.setBeforeRefreshLockHook(func() {
		close(refreshWaiting)
	})
	go func() {
		_, err := tm.RefreshToken(context.Background())
		refreshDone <- err
	}()
	<-refreshWaiting

	select {
	case <-transport.secondAttempt:
		t.Fatal("observed a second refresh request while first refresh was in progress")
	case <-time.After(200 * time.Millisecond):
	}
	assert.Equal(t, int32(1), transport.calls.Load())

	close(transport.release)

	select {
	case err := <-getDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("GetToken did not finish")
	}

	select {
	case err := <-refreshDone:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("RefreshToken did not finish")
	}

	assert.Equal(t, int32(1), transport.calls.Load())
}

func TestTokenManagerGetTokenDoesNotMutateStateOnStorageFailure(t *testing.T) {
	config := NewOAuthConfig(Config{
		ClientID:     "test-client",
		ClientSecret: "test-secret",
		RedirectURI:  "https://example.com/callback",
	})

	storage := &failingStorage{
		token: &Token{
			AccessToken:  "expired-access-token",
			RefreshToken: "refresh-token",
			ExpiresAt:    time.Now().Add(-time.Minute),
		},
		saveErr: fmt.Errorf("disk full"),
	}

	tm := NewTokenManager(config, storage).WithHTTPClient(&http.Client{
		Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Header:     make(http.Header),
				Body: io.NopCloser(strings.NewReader(
					`{"access_token":"fresh-access-token","refresh_token":"refresh-token","token_type":"Bearer","expires_in":7200}`,
				)),
			}, nil
		}),
	})

	_, err := tm.GetToken(context.Background())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to save refreshed token")

	require.NotNil(t, tm.token)
	assert.Equal(t, "expired-access-token", tm.token.AccessToken)
	require.NotNil(t, storage.lastSavedToken)
	assert.Equal(t, "fresh-access-token", storage.lastSavedToken.AccessToken)
	assert.Equal(t, "expired-access-token", storage.token.AccessToken)
	assert.Equal(t, 1, storage.saveCalls)
}

type countingBlockingTransport struct {
	started       chan struct{}
	release       chan struct{}
	secondAttempt chan struct{}
	body          string
	calls         atomic.Int32
}

func (t *countingBlockingTransport) RoundTrip(*http.Request) (*http.Response, error) {
	callNum := t.calls.Add(1)
	if callNum == 1 {
		close(t.started)
		<-t.release
	} else if callNum == 2 {
		close(t.secondAttempt)
	}

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(t.body)),
	}, nil
}

type failingStorage struct {
	token          *Token
	lastSavedToken *Token
	saveCalls      int
	saveErr        error
}

func (s *failingStorage) SaveToken(token *Token) error {
	s.saveCalls++
	s.lastSavedToken = token
	if s.saveErr != nil {
		return s.saveErr
	}
	s.token = token
	return nil
}

func (s *failingStorage) LoadToken() (*Token, error) {
	if s.token == nil {
		return nil, fmt.Errorf("no token stored")
	}
	return s.token, nil
}

func (s *failingStorage) ClearToken() error {
	s.token = nil
	return nil
}

func (s *failingStorage) HasToken() bool {
	return s.token != nil
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
