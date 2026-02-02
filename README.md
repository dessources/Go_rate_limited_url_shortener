# URL shortener behind Rate Limiter

## Overview

This is a url shortening service protected by two layers of rate limiting, written in Go. The purpose of this project was to improve my systems engineering skills and familiarize myself with the Go programming language.

## Live Website

Check out the live production version here: [https://pety.to](https://www.pety.to)

"pety" is derived from the french word "petit" which means small.

## Features

- Token Bucket algorithm for Global rate limiting
- Sliding Window log algorithm for per-client rate limiting
- URL shortener
- SSE live metrics
- Isolated stress testing
- Dockerized deployment

## Architecture

This project uses two different rate limiting algorithms:

- **Token Bucket** for global rate limiting, protecting the API against a surge in requests
- **Sliding Window Log** for per-client rate limiting

Both rate limit algorithms are applied through middlewares and for certain routes are combined through middleware composition.

Live metrics are streamed to the frontend with data about the global rate limit, active users, and active URLs. This is achieved through Server Side Events (SSE). In SSE, a client opens a connection to the server and maintains it open. The server then pushes updates to the client until either one closes the connection. This is better than polling and WebSockets for these reasons:

1. **No wasted resources** — Unlike polling, the server pushes updates only when there are new updates, rather than the client constantly querying
2. **Unidirectional communication** — Once connected, the client lets the server do all the talking, which is more efficient than the constant back-and-forth of WebSockets

## Tech Stack

Go, React/Next.js, Shadcn UI, Docker, Bash

## Getting Started

Prerequisites:

- **pnpm** is preferred but you can use any package manager
- **Go** must be installed on your system

With Docker:

```bash
# Using docker-compose (recommended)
docker-compose up --build

# Or manually with docker
docker build -t rate-limiter .
docker run -p 8090:8090 -p 8091:8091 --env-file .env rate-limiter
```

Without Docker:

In the project's root directory, install dependencies with

```bash
pnpm install
```

Then, run:

```bash
go run .
```

The app should start running on http://localhost:8090.

## Configuration

All settings are configurable via environment variables. Create a `.env` file in the project root.

### Server

| Variable               | Description                         | Default                                       |
| ---------------------- | ----------------------------------- | --------------------------------------------- |
| `BASE_URL`             | Base URL for generated short links  | `https://pety.to`                             |
| `SERVER_ADDR`          | Server listen address               | `:8090`                                       |
| `TEST_SERVER_ADDR`     | Isolated stress test server address | `:8091`                                       |
| `CORS_ALLOWED_ORIGINS` | Comma-separated allowed origins     | `http://localhost:3000,http://localhost:8090` |

### Global Rate Limiter (Token Bucket)

| Variable              | Description             | Default  |
| --------------------- | ----------------------- | -------- |
| `GLOBAL_LIMITER_CAP`  | Maximum token capacity  | `50000`  |
| `GLOBAL_LIMITER_RATE` | Tokens added per minute | `600000` |

### Per-Client Rate Limiter (Sliding Window)

| Variable                        | Description                            | Default |
| ------------------------------- | -------------------------------------- | ------- |
| `PER_CLIENT_LIMITER_CAP`        | Max number of tracked clients          | `50000` |
| `PER_CLIENT_LIMITER_LIMIT`      | Requests allowed per window            | `10`    |
| `PER_CLIENT_WINDOW_SECONDS`     | Window duration in seconds             | `60`    |
| `PER_CLIENT_LIMITER_CLIENT_TTL` | Inactive client cleanup time (seconds) | `1800`  |

### URL Shortener

| Variable              | Description                     | Default  |
| --------------------- | ------------------------------- | -------- |
| `SHORTENER_CAP`       | Max stored URLs                 | `100000` |
| `SHORTENER_TTL_HOURS` | URL expiration time (seconds)   | `3600`   |
| `SHORT_CODE_LENGTH`   | Length of generated short codes | `4`      |
| `MAX_URL_LENGTH`      | Maximum allowed URL length      | `4096`   |

## Live Demo

[Go to pety.to for live demo](https://pety.to)

## What I learned

- **Concurrency patterns in Go:** Using RWMutex for read-heavy operations (per-client limiter) vs regular Mutex for write-heavy operations (token bucket refills). Understanding when to use channels vs shared memory with locks.

- **Rate limiting algorithm trade-offs:** Token bucket is memory-efficient but allows bursts; sliding window log provides precise rate limiting but requires more memory per client. Production systems need both layers to handle different threat models.

- **Server-Sent Events (SSE) for real-time metrics:** SSE is ideal for one-way server→client streaming. No polling overhead (unlike REST), no bidirectional complexity (unlike WebSockets). Perfect for metrics dashboards where the server dictates update frequency.

- **Graceful shutdown patterns:** Using context propagation and signal handling to ensure all resources (rate limiter goroutines, HTTP connections, test servers) clean up properly on SIGINT/SIGTERM.

## License

MIT
