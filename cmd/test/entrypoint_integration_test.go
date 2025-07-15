package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"context"

	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/app"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/getenv"
	"github.com/felipeazsantos/ratelimiter-fullcycle-challenge/config/redisdb"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// Constantes de teste
const (
	TestIP1 = "1.2.3.4:12345"
	TestIP2 = "7.7.7.7:3333"
	TestIP1Isolation = "1.2.3.4:2222"
	TestTokenIntegration = "integration-token-123"
	TestTokenIsolation1 = "isolation-token-1"
	TestTokenIsolation2 = "isolation-token-2"
)

// Limpa as chaves de teste do Redis
func cleanupTestKeys(rdb *redis.Client) {
	ctx := context.Background()
	ipKeys := []string{
		"1.2.3.4:3333",
		"7.7.7.7:4444",
		"1.2.3.4:5555",
	}
	tokenKeys := []string{
		"integration-token-123",
		"isolation-token-1",
		"isolation-token-2",
	}
	for _, key := range ipKeys {
		rdb.Del(ctx, key)
	}
	for _, key := range tokenKeys {
		rdb.Del(ctx, key)
	}
}

func setupTestServer(t *testing.T) (*httptest.Server, *redis.Client) {
	err := getenv.InitConfig("../../.env")
	if err != nil {
		assert.Fail(t, fmt.Sprintf("Failed to load .env: %v", err.Error()))
	}
	if getenv.AppConfig == nil {
		assert.Fail(t, "AppConfig is nil after InitConfig")
	}

	router := chi.NewRouter()
	rdb := redisdb.NewRedisClient(getenv.AppConfig.RedisAddr)
	err = app.StartDependencies(router, rdb)
	if err != nil {
		t.Fatalf("Failed to start dependencies: %v", err)
	}
	return httptest.NewServer(router), rdb
}

func TestRateLimiter_IP(t *testing.T) {
	ts, rdb := setupTestServer(t)
	defer ts.Close()
	defer rdb.Close()
	defer cleanupTestKeys(rdb)

	client := &http.Client{}
	maxRequests := getenv.AppConfig.IPRateLimiterMaxRequest
	blockTime, _ := time.ParseDuration(getenv.AppConfig.IPRateLimiterBlockTime)

	for i := 0; i < maxRequests; i++ {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.RemoteAddr = TestIP1
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Exceeding the limit
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.RemoteAddr = TestIP1
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "maximum number of requests")
	resp.Body.Close()

	// After block time
	t.Logf("Waiting %v to test unblock...", blockTime)
	time.Sleep(blockTime + time.Second)

	// should allow again
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.RemoteAddr = TestIP1
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestRateLimiter_Token(t *testing.T) {
	ts, rdb := setupTestServer(t)
	defer ts.Close()
	defer rdb.Close()
	defer cleanupTestKeys(rdb)

	client := &http.Client{}
	maxRequests := getenv.AppConfig.TokenRateLimiterMaxRequest
	blockTime, _ := time.ParseDuration(getenv.AppConfig.TokenRateLimiterBlockTime)
	token := TestTokenIntegration

	for range maxRequests {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.Header.Set("API_KEY", token)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Exceeding the limit
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "maximum number of requests")
	resp.Body.Close()

	// After block time
	t.Logf("Waiting %v to test unblock...", blockTime)
	time.Sleep(blockTime + time.Second)

	// should allow again
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}

func TestRateLimiter_Isolation(t *testing.T) {
	ts, rdb := setupTestServer(t)
	defer ts.Close()
	defer rdb.Close()
	defer cleanupTestKeys(rdb)

	client := &http.Client{}

	// IP1 hits the limit
	maxRequests := getenv.AppConfig.IPRateLimiterMaxRequest
	for range maxRequests {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.RemoteAddr = TestIP1Isolation
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// IP1 blocked
	req, _ := http.NewRequest("GET", ts.URL+"/", nil)
	req.RemoteAddr = TestIP1Isolation
	resp, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	resp.Body.Close()

	// IP2 should pass normally
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.RemoteAddr = TestIP2
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()

	// Token1 hits the limit
	maxToken := getenv.AppConfig.TokenRateLimiterMaxRequest
	token1 := TestTokenIsolation1
	token2 := TestTokenIsolation2

	for range maxToken {
		req, _ := http.NewRequest("GET", ts.URL+"/", nil)
		req.Header.Set("API_KEY", token1)
		resp, err := client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Token1 blocked
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token1)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	resp.Body.Close()

	// Token2 should pass normalmente
	req, _ = http.NewRequest("GET", ts.URL+"/", nil)
	req.Header.Set("API_KEY", token2)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	resp.Body.Close()
}
