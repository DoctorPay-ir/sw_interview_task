# Software Engineering Interview Task

## Overview

This repository contains a minimal HTTP service written in Go. It exposes two endpoints that return greeting messages. Your task is to extend this service with observability, rate limiting, and a structure suitable for a larger production codebase.

You may use standard library packages, third-party libraries, or both. Document any assumptions you make.

## Starting Point

The service currently:

- Listens on a configurable port (default `8080`)
- Exposes `GET /greet` and `GET /goodbye`
- Accepts optional query parameters: `first_name`, `last_name`
- Uses a `GreetingService` interface with a default English implementation
- Starts a background goroutine that prints `ok` every 5 seconds

Run locally:

```bash
go run main.go
go run main.go -port 9090
```

Example requests:

```bash
curl "http://localhost:8080/greet?first_name=Jane&last_name=Doe"
curl "http://localhost:8080/goodbye?first_name=Jane"
```

## Requirements

Implement the following. Point values indicate relative weight for evaluation.

### 1. Basic Metrics — Success / Failure Rate (5 pts)

Track and expose metrics for HTTP request outcomes:

- Count successful responses (2xx)
- Count failed responses (4xx, 5xx)
- Make metrics accessible (e.g. `/metrics` endpoint, structured logs, or another reasonable approach)

At minimum, a reviewer should be able to determine success vs. failure rates per endpoint.

### 2. Config Management (5 pts)

Centralize application configuration with sensible defaults. At minimum, make these configurable without code changes (flags, environment variables, config file, or a combination):

- **Port** — HTTP server listen port
- **OK print interval** — how often the background goroutine prints `ok`
- **Serve metrics** — whether the metrics endpoint (or equivalent metrics surface) is enabled

Validate configuration at startup where reasonable, and document defaults plus how to override each setting.

### 3. Basic Traceability (5 pts)

Add request-level traceability so that:

- Each incoming request can be correlated across logs (and optionally metrics)
- A unique identifier is propagated through the request lifecycle (e.g. request ID / trace ID in context)
- Logs (or equivalent output) include enough context to follow a single request

### 4. Basic Rate Limiting (5 pts)

Protect the service with a simple rate limiter:

- Reject or throttle excessive requests from a client
- Return an appropriate HTTP status when the limit is exceeded (e.g. `429 Too Many Requests`)
- Apply rate limiting in a way that does not require changes to core handler business logic

### 5. Configurable Rate Limiting (15 pts)

Extend rate limiting so limits can be configured without code changes:

- Support configuration via flags, environment variables, config file, or a combination
- At minimum, make these configurable:
  - Whether rate limiting is enabled
  - The rate limit threshold (e.g. requests per time window)
  - The time window (if applicable)
  - Support distinct rate-limit keys per route or handler (e.g. `/greet` vs `/goodbye`), derived at bootstrap via explicit registration or code generation from handler signatures

### 6. Module Design & Dependency Injection (15 pts)

Restructure the project for maintainability as it grows:

- Split code into logical packages/modules (e.g. config, handlers, middleware, services)
- Use dependency injection (constructor injection preferred) instead of wiring everything in `main`
- Keep interfaces at boundaries where it improves testability and swap-ability
- `main` should primarily compose and start the application

This item is evaluated on clarity, separation of concerns, and whether the design would scale to a larger team and codebase — not on matching a specific folder layout.

### 7. Graceful Shutdown (5 pts)

Shut down the service cleanly when it receives a termination signal (e.g. `SIGINT`, `SIGTERM`):

- Stop accepting new HTTP requests
- Allow in-flight requests to complete (within a reasonable timeout)
- Stop background work started in `main` (including the goroutine that prints `ok`) without leaking goroutines or leaving the process in a hung state
- Exit with a clear, predictable lifecycle: startup → running → draining → stopped

Document how shutdown is triggered and any timeout or behavior choices in `NOTES.md`.

## Deliverables

1. Your implementation (code changes in this repo or a fork)
2. A brief `NOTES.md` (or section in your PR description) covering:
   - How to build and run the service
   - How to configure the service (port, OK print interval, metrics, rate limits)
   - Where to find metrics and how to interpret them
   - How graceful shutdown works (signals, timeouts, background goroutines)
   - Any trade-offs or shortcuts you took due to time constraints

## Evaluation Criteria

| Area | Weight | What we look for |
|------|--------|------------------|
| Metrics | 5 | Correct success/failure tracking, usable exposure |
| Config management | 5 | Port, OK interval, metrics toggle; documented defaults |
| Traceability | 5 | Request correlation, context propagation |
| Basic rate limit | 5 | Works end-to-end, clean integration |
| Configurable rate limit | 15 | Flexible config, sensible defaults, documented |
| Module design & DI | 15 | Clear boundaries, testable, idiomatic Go |
| Graceful shutdown | 5 | Signal handling, drain in-flight work, clean exit |
| Code quality | — | Readability, error handling, idiomatic patterns |


## Out of Scope (Optional)

You do **not** need to implement:

- Authentication / authorization
- Persistence or databases
- Deployment manifests or CI pipelines
- Full OpenTelemetry / Prometheus integration (unless you choose to)

Focus on the requirements above within the time you have available.

