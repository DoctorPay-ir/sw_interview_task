# Software Engineering Interview Task

## Overview

This repository contains a minimal HTTP service written in Go. It exposes two endpoints that return greeting messages. Your task is to extend this service with observability, rate limiting, and a structure suitable for a larger production codebase.

You may use standard library packages, third-party libraries, or both. Document any assumptions you make.

### Bonus signals (not required)

These are not part of the scored requirements, but strong candidates often demonstrate:

| Signal | Weight |
|--------|--------|
| **Using code generation tools** (e.g. `go generate`, `stringer`, OpenAPI/protobuf codegen) where they reduce boilerplate | **+** |
| **Writing your own code generator** for repetitive wiring (e.g. rate-limit keys, handler registration, config structs) | **++** |
| **Reusable abstractions** — packages/modules generic enough to drop into other services without greeting-specific coupling | **++** |

If you use or build generators, briefly explain what they produce and how to run them in `NOTES.md`.

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

Add a simple rate limit to protect the HTTP server from abuse (e.g. request floods, basic DoS-style traffic):

- Limit incoming HTTP requests — no abstraction required; middleware or handler-level limiting is fine
- Reject excessive traffic with an appropriate status (e.g. `429 Too Many Requests`)
- Keep it straightforward: a fixed or hard-coded limit is acceptable for this item

This is intentionally separate from the configurable interface rate limiter in the next section.

### 5. Configurable Rate Limiting (15 pts)

Build a **configurable, reusable rate limiter** that can wrap arbitrary Go interfaces with minimal boilerplate. Use `GreetingService` as the concrete integration in this repo.

**Design goals**

- **Interface-first** — the limiter package must not depend on `net/http`, routes, or handlers; it operates on interface methods only
- **Minimal wrapping cost** — adding rate limits to a new interface should require little or no hand-written glue per method (e.g. generic wrapper, decorator registry, or code generation)
- **Per-method limits** — distinct limits/metadata for each wrapped method (e.g. `Greet` vs `Goodbye`), not a single global bucket for the whole service

**Rate limit metadata**

Each limited method should be configurable with metadata such as:

| Field | Description |
|-------|-------------|
| **Key** | Stable identifier for the limit bucket (e.g. `GreetingService.Greet`) |
| **Enabled** | Whether limiting applies to this method |
| **Threshold** | Max allowed calls in the window |
| **Window** | Time window for the threshold (if applicable) |

Metadata may be loaded from flags, environment variables, config file, or a combination — without code changes.

**Bootstrap / registration**

- Register limits at startup by associating metadata with interface methods
- Keys may be derived via explicit bootstrap registration or code generation from interface/handler signatures
- Wiring in `main` (or a composition root) should compose: `RateLimited(GreetingService)` → inject into handlers

Document the metadata schema, defaults, and an example config in `NOTES.md`.

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

### 8. Persian & Bilingual Greeters (5 pts)

Extend `GreetingService` with additional implementations:

- **Persian greeter** — returns greet/goodbye messages in Persian (Farsi), following the same name rules as the English implementation (first name only vs first + last name)
- **Bilingual greeter** — returns a response that includes both English and Persian (composition/wrapping is up to you; keep it readable)

Wire greeters through the existing interface and DI boundaries (do not hard-code language logic inside HTTP handlers). Language selection may be driven by query parameter, header, config, or another reasonable mechanism — document your choice and example requests in `NOTES.md`.

## Deliverables

1. Your implementation (code changes in this repo or a fork)
2. A brief `NOTES.md` (or section in your PR description) covering:
   - How to build and run the service
   - How to configure the service (port, OK print interval, metrics)
   - How to configure rate limits (metadata schema, per-method keys, example config)
   - Where to find metrics and how to interpret them
   - How graceful shutdown works (signals, timeouts, background goroutines)
   - How to use Persian and bilingual greeters (language selection, example requests)
   - Any trade-offs or shortcuts you took due to time constraints

## Evaluation Criteria

| Area | Weight | What we look for |
|------|--------|------------------|
| Metrics | 5 | Correct success/failure tracking, usable exposure |
| Config management | 5 | Port, OK interval, metrics toggle; documented defaults |
| Traceability | 5 | Request correlation, context propagation |
| Basic rate limit | 5 | Simple HTTP-level protection against abusive traffic |
| Configurable rate limit | 15 | Per-method metadata, minimal-wrap design, `GreetingService` integration |
| Module design & DI | 15 | Clear boundaries, testable, idiomatic Go |
| Graceful shutdown | 5 | Signal handling, drain in-flight work, clean exit |
| Persian & bilingual greeters | 5 | Correct Persian text, bilingual composition, clean DI |
| Code quality | — | Readability, error handling, idiomatic patterns |
| Using code generators | bonus (+) | Appropriate use of existing codegen tools |
| Writing code generators | bonus (++) | Custom generator that reduces repetitive wiring |
| Reusable abstractions | bonus (++) | Modules usable across services, not tied to this app |

## Out of Scope (Optional)

You do **not** need to implement:

- Authentication / authorization
- Persistence or databases
- Deployment manifests or CI pipelines
- Full OpenTelemetry / Prometheus integration (unless you choose to)

Focus on the requirements above within the time you have available.

## Scoring & How to Approach the Task

**Minimum score to pass: 20 points** (out of 60 available across the scored items above).

All requirements are **optional individually**, but each one matters. You do **not** need to implement everything — choose the items that best show your strengths and aim for at least 20 points total.

**Quality over breadth:** two areas implemented **very well** (e.g. configurable rate limiting + module design) is better than touching every item at an average level.

Pick what fits your time and expertise, document your choices in `NOTES.md`, and make the parts you submit polished and easy to review.
