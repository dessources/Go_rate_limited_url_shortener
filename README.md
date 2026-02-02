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

## Live Demo

[Go to pety.to for live demo](https://pety.to)

## What I learned

- Making APIs thread-safe with mutexes. When different threads might need to access the same data entity as often is the case in APIs, it is necessary to use mutexes to ensure only on routine accesses the data at a time preventing race conditions.

- SSE is a pretty cool communication method for metrics streaming. It is pretty useful when we want the server to send data to a client on its own time.

## License

MIT
