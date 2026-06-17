# Architecture (devconv)

Org template layout for services on [ecfg](https://github.com/omcrgnt/ecfg), [builder](https://github.com/omcrgnt/builder), [res](https://github.com/omcrgnt/res), [sdi](https://github.com/omcrgnt/sdi), [srv-http](https://github.com/omcrgnt/srv-http), and [runner](https://github.com/omcrgnt/runner).

## Layers

| Layer | Path | Role |
|-------|------|------|
| **app** | `cmd/<app>/`, `internal/config/` | Composition root: config, wiring, lifecycle |
| **domain** | `internal/domain/` | Models, ports, business logic |
| **api** | `internal/api/<transport>/` | Inbound adapters: invoke domain (HTTP, gRPC, CLI, cron, …) |
| **data** | `internal/data/<kind>/<aggregate>/<backend>/` | Data movement: repos, queues, remote storage |

**Dependency rule:** `api → domain ← data`. Domain imports neither `api` nor `data`.

## `api` (broad sense)

Everything that **initiates a use case** lives under `internal/api/`:

```
internal/api/
  http/       # REST (this demo)
  grpc/       # future
  graphql/    # future
  cli/        # future: in-process commands (same binary)
  cron/       # future: scheduled jobs
```

Handlers: parse input → call `domain/service/<name>` → map response. No business rules in `api`.

## Domain

### Shared models

Models used by more than one service or layer:

```
internal/domain/model/
  item.go
  order.go
```

No `json` / `db` / transport tags — pure domain types.

### Per-subdomain service

Even in a narrow microservice, split by subdomain when needed:

```
internal/domain/service/<name>/
  service.go      # Config, Build, Deps, Inject at top; then Service + methods
  interface.go    # all ports/deps (driven + outbound); mockgen-friendly
  model.go        # models used only by this subdomain (optional; often empty)
```

- **`interface.go`** — ports the service needs (repos, event publishers, …). Interface is declared by the consumer; implementation lives in `data/`.
- **`service.go`** — wiring (`Config`, `Build`, `Deps`, `Inject`) stays at the **top** of the file, not in a separate `config.go`.
- Service-local models stay in `model.go`; shared ones go to `domain/model/`.

### Errors

Domain errors live locally (e.g. `domain/errors.go` or per-service) for now.

**Backlog:** org-wide shared errors package — not used yet.

## Cross-cutting: logging

**Backlog:** [github.com/omcrgnt/logger](https://github.com/omcrgnt/logger) + meta [`res/core/use`](https://github.com/omcrgnt/res/tree/main/core/use) — wired in demo via single blank import.

### System package (zero-config)

- Blank import in the binary: `import _ "github.com/omcrgnt/res/core/use"` pulls logger and telemetry system defaults.
- In `init()`, logger registers in `res` with a **production-ready default** (stderr, structured output, sensible level). No `AppConfig` field required for most services.
- Optional override: `logger.Config` in `AppConfig` → env prefix → `Build()` → `res`. Zero value = leave the default from `init` unchanged.

### Public API

- Only `*Ctx` functions (`InfoCtx`, `ErrorCtx`, …). No legacy calls without `context.Context`.
- Attrs: `...any` (slog-style). Prefer benchmarks over custom attr types when choosing the surface.

### Context and engine

- **Singleton** log engine for the process. Do not fork or allocate a logger per request.
- `context.Context` carries **observation attrs** (correlation_id, trace_id via `obs`), not a `*Logger`. Merge ctx attrs at write time on each call.

### Where to log

| Layer | Rule |
|-------|------|
| **domain** | Return errors; no logger in struct fields or `interface.go`. `ctx` for cancel/timeout/deadlines only. |
| **api** | Log at boundaries (handlers, middleware): `logger.*Ctx(ctx, …)` after obs is attached. |
| **data** | Log I/O failures and adapter diagnostics at the repo boundary. |

### Bootstrap vs runtime

- Before `ecfg.Parse` or on fatal startup paths: stdlib `log` / stderr is fine (no structured pipeline yet).
- After successful startup: structured logging for the **full process lifetime** via `logger.*Ctx`.

## `data`

Everything needed to **store, fetch, or move data**. Naming: **`aggregate → backend`** (universal for template):

```
internal/data/
  sync/
    item/
      memory/       # in-memory ItemRepository (this demo)
      postgres/     # future
      remote/       # future: REST/gRPC client to another service as repo
    order/
      postgres/     # future
  async/
    item/
      kafka/        # future: consumers, producers
    order/
      nats/         # future
```

Standard files per repo module:

```
internal/data/sync/<aggregate>/<backend>/
  repo.go         # Config, Build at top; Repo + methods
  model.go        # persistence/adapter models (db tags, etc.); optional, often empty
```

- **In/out:** domain models from `domain/model/` (map inside repo if persistence model differs).
- **Async** stays under `data/`, not `api` — data stream, not a public call contract.

## App pipeline

[`cmd/demo/app.go`](cmd/demo/app.go):

1. `ecfg.Parse[AppConfig]()` — load env
2. `builder.Build(cfg, res.Default)` — `Config.Build()` per module
3. `sdi.Resolve(res.Default)` — dedupe (`DefaultDedupPolicy` + `Remove`), then inject deps
4. `runner.New(runnerPool{reg: res.Default}).Run / Stop` — start stoppable resources (HTTP server)

`runnerPool` adapts `res.Registry` (`WalkEntries`) to `runner.Pool` (`Walk`). System defaults (e.g. logger) register via `AddWithTags(..., TagReplaceable)` in library `init` — not in app code.

See [docs/res-sdi-coupling.md](docs/res-sdi-coupling.md) for res↔sdi design variants (ADR).

Each module exposes `Config` with `Build() (any, error)`; wired resources implement `Deps()` / `Inject()` where needed.

## Where to put new code

| Adding | Location |
|--------|----------|
| Shared domain model | `internal/domain/model/<name>.go` |
| Subdomain-only model | `internal/domain/service/<name>/model.go` |
| Port (interface) | `internal/domain/service/<name>/interface.go` |
| Business logic | `internal/domain/service/<name>/service.go` |
| Domain errors | `internal/domain/errors.go` (local; org package — backlog) |
| REST endpoint | `internal/api/http/` |
| gRPC service | `internal/api/grpc/` |
| DTO / request-response types | same package as transport under `internal/api/` |
| Sync repo | `internal/data/sync/<aggregate>/<backend>/` |
| Async producer/consumer | `internal/data/async/<aggregate>/<broker>/` |
| AppConfig field | `internal/config/config.go` |
| Binary entrypoint | `cmd/<app>/` |

## AppConfig naming

Field names match modules (`HTTP`, `ItemRepo`, `ItemService`, …), not generic names like `Controller`.

```go
type AppConfig struct {
    ItemRepo    memory.Config   // data/sync/item/memory
    ItemService item.Config     // domain/service/item
    HTTP        http.Config
    Metrics     http.MetricsConfig
    HTTPServer  *srvhttp.Config[*http.API]
}
```

## Tests

- Domain/service: unit tests with `builder` + `sdi` on a minimal config struct
- Data repos: test repo behaviour against the port contract
- `go test ./...` from repo root
