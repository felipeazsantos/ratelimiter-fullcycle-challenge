# Rate Limiter FullCycle Challenge

## What is a Rate Limiter?
A **rate limiter** is a mechanism used to control the number of requests that a user, service, or IP can make to an API or system within a certain time interval. This prevents abuse, denial-of-service (DoS) attacks, and ensures fairer use of resources.

## How it works in this project
In this project, the rate limiter was implemented to limit the number of requests each client can make within a configurable time window. If the limit is exceeded, new requests are blocked until the window resets.

The basic operation is:
- Each client (identified by IP or token) has a request counter.
- This counter is reset every time window (for example, every 1 minute).
- If the request limit is reached before the end of the window, subsequent requests receive an error response (e.g., HTTP 429 Too Many Requests).

## Rate Limiter Configuration
Configuration is done via environment variables, which can be set in the `.env` file or directly in the container/host environment.

### Available variables:

- `RATE_LIMITER_IP_MAX_REQUESTS`: Maximum number of requests per IP before blocking (default: 3)
- `RATE_LIMITER_IP_BLOCK_TIME`: Time period (in seconds or duration string) to block IP after exceeding the limit (default: 10s)
- `RATE_LIMITER_TOKEN_MAX_REQUESTS`: Maximum number of requests per token before blocking (default: 5)
- `RATE_LIMITER_TOKEN_BLOCK_TIME`: Time period (in seconds or duration string) to block token after exceeding the limit (default: 10s)

## Usage Example
Assuming the server is running on `localhost:8080`, simply make requests to the protected endpoint:

```bash
curl http://localhost:8080/
```

If the number of requests exceeds the configured limit within the time window, the response will be:

```json
{
  "error": "you have reached the maximum number of requests or actions allowed within a certain time frame"
}
```

And the HTTP status will be 429.

## How to run the project
1. Clone the repository:
   ```bash
   git clone <repo-url>
   cd ratelimiter-fullcycle-challenge
   ```
2. Configure the `.env` file as desired.
3. Run the project:
   ```bash
   go run ./cmd/server/main.go
   ```
   Or use Docker:
   ```bash
   docker-compose up --build
   ```

---


