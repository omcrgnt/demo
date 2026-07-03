# demo — AppResources reference app

Reference application for the target architecture: single `go.mod`, org libs from `github.com/omcrgnt/*`, explicit pipeline (no legacy `Resourcer`).

See [ARCHITECTURE.md](ARCHITECTURE.md) for contracts, roles, and backlog.

## Layout

```text
cmd/app/           main, AppResources catalog
internal/          domain, api, data
.env.template      ecfg-gen output (repo root)
env.md             ecfg usage docs
```

## Pipeline

Handled by [`github.com/omcrgnt/app`](https://github.com/omcrgnt/app) v0.21:

```text
app.Run(&appResources, pipeline)
  → fill → LoadEnv → materialize → merge → Transform → Resolve → Serve
```

`AppResources` holds **resources** only: each field is [ResourceFactory] or [Configurable]. Configurable resources use two types — resource + Spec/Config (`BuildConfig()` → spec, `Build()` → resource). ecfg walks the spec, not wire fields.

Blank-import [`github.com/omcrgnt/app/use`](https://github.com/omcrgnt/app) for default `*app.App` and `*runner.Runner`. Add `App *app.App \`ecfg:"APP"\`` only to override app config from env.

## AppResources

Fields follow `{type}{subject}` (e.g. `RepoOrder`, `ServiceItem`).

| Field | Mechanism |
|-------|-----------|
| `RepoItem` | [ResourceFactory] → `*memory.Repo` |
| `ServiceItem` | `item.ServiceRoot` [Configurable] → `item.Spec` |
| `RepoOrder` | [ResourceFactory] → `*ordermemory.Repo` |
| `ServiceOrder` | [ResourceFactory] → `*order.Service` |
| `ServerHTTPOps` | `ophttp.Config` [Configurable] → same (Build → systemServer) |
| `ServerHTTPItem` | `srvhttp.Server[*http.API]` [Configurable] → `srvhttp.Config[*http.API]` |
| `APIItem` | [ResourceFactory] |
| `ServerHTTPOrder` | `srvhttp.Server[*orderhttp.API]` [Configurable] → `srvhttp.Config[*orderhttp.API]` |
| `APIOrder` | [ResourceFactory] |

HTTP metrics: [`srv-http/use`](https://github.com/omcrgnt/srv-http) registers shared `HTTPMetrics` (slok recorder); [`ops/metrics/use`](https://github.com/omcrgnt/ops) provides registry + actuator. Scrape via ops `:9090/metrics` (default).

### app.Pipeline (set in main)

- `Registry` — `unique.Global()` or `unique.New()` in tests
- `EnvPrefix` — ecfg prefix, e.g. `"DEMO"`
- `Transforms` — e.g. `[]res.TransformFunc{obs.ApplyTransform}`; empty skips Transform

### Blank imports (main)

```go
_ "github.com/omcrgnt/app/use"
_ "github.com/omcrgnt/logger/use"
_ "github.com/omcrgnt/telemetry/use"
_ "github.com/omcrgnt/srv-http/use"           // HTTPMetrics singleton
_ "github.com/omcrgnt/ops/metrics/use"       // registry + metrics actuator
_ "github.com/omcrgnt/ops/transport/http/use" // probe + ops HTTP (:9090)
```

## Commands

From repo root:

```bash
cp .env.example .env
task r-app      # build + run with .env
task test       # go test ./...
task gen        # go generate (obsgen, ecfg-gen → .env.template + env.md)
```

- `.env.template` — generated keys only (`KEY=`)
- `env.md` — generated usage docs (tables by ecfg block)
- `.env.example` — sample values for local run (`cp .env.example .env`)

Ops probes (with app running):

```bash
curl -s :9090/livez
curl -s :9090/readyz
curl -s :9090/metrics | head
```

## Notes

- Stack: `app`, `ecfg`, `res`, `sdi`, `runner`, `obs`, `srv-http`, `ops`, `logger`, `telemetry` at org v0.21 / v0.22.
- Catalog slots with `BuildConfig` live in domain/ops/srv-http libs (`item.ServiceRoot`, `ophttp.Config`, `srvhttp.Server[T]`). No `catalog_wire.go` in demo.
- Temporary `replace` for local dev: `builder`, `ecfg`, `ops` — see org [backlog](https://github.com/omcrgnt/backlog).
- External require: `github.com/omcrgnt/proto/gen/go` (srv-http Label/Host/Port).
