# Kodama Engineering Guide

Kodama is a personal AI assistant written in Go. Optimize for code that is easy to read, debug, and evolve by one experienced backend engineer.

## Development Style

- Build outside-in: start from the behavior at the boundary, then work inward.
- Prefer small, reviewable implementation chunks.
- Use TDD where it clarifies behavior, especially at boundaries and tricky logic.
- Keep tests lean. This is a single-contributor project, so avoid broad coverage for obvious wiring or framework-free boilerplate.
- Add tests when behavior could regress silently, when parsing/auth decisions matter, or when domain/use-case logic becomes non-trivial.

## Architecture

- Follow DDD pragmatically. Use the language and boundaries when they improve readability; avoid ceremony for its own sake.
- Keep use cases explicit and easy to trace.
- A use case may live in its own folder when it has meaningful behavior, request/response objects, or ports.
- Keep transport and infrastructure details out of use-case request/response types.
- Infrastructure packages adapt external systems such as HTTP, Telegram, databases, queues, and LLM providers.
- Domain/use-case packages define the interfaces they need. Infrastructure implements them.

Expected package direction:

```text
cmd/api                  process lifecycle and top-level wiring
internal/config          runtime configuration
internal/infra/api       HTTP server, routing, auth, request parsing, serialization
internal/infra/logging   slog setup
internal/infra/*         external adapters
internal/usecase/*       application use cases and their ports
```

## Go Conventions

- Always pass `context.Context` as the first parameter for operations that may block, call I/O, depend on request scope, or may need cancellation.
- Use `log/slog` for logging.
- Set the global logger during startup through `internal/infra/logging`.
- Prefer explicit constructors and simple structs over framework magic.
- Use interfaces only when they provide clear value, usually at package boundaries where a use case depends on an external capability.
- Keep configuration environment-based until there is a concrete need for files or a richer config system.

## Operational Direction

- Start with one executable.
- It is fine to add more executables later, such as a worker for background jobs.
- Do not mix deployment/admin actions into normal runtime startup unless there is a clear reason.
- For Telegram webhooks, serving the webhook is runtime behavior; registering the webhook URL with Telegram is deployment setup.
