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
3. `sdi.Resolve(res.Default)` — inject deps
4. `runner.New(res.Default).Run / Stop` — start stoppable resources (HTTP server)

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
