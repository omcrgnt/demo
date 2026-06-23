# demo — AppResources reference app

Reference application for the target architecture: single `go.mod`, org libs from `github.com/omcrgnt/*`, pipeline without legacy `Resourcer`.

See [ARCHITECTURE.md](ARCHITECTURE.md) for contracts, roles, and backlog.

## Pipeline

Handled by [`github.com/omcrgnt/app`](https://github.com/omcrgnt/app):

```text
app.Run(&appResources, pipeline)  // pipeline configured explicitly in main
  → Seed → Apply → Build → Transform → Resolve → runner
```

`AppResources` holds **resources** only: each field is [NewResourceer] or [BuildConfiger]. Configurable resources use two types — resource + Spec/Config (`BuildConfig()` → spec, `Build()` → resource). ecfg walks the spec, not wire fields.

## AppResources

Fields follow `{type}{subject}` (e.g. `RepoOrder`, `ServiceItem`). Each field is [NewResourceer] or [BuildConfiger].

| Field | Mechanism |
|-------|-----------|
| `App` | `*app.App` [BuildConfiger] → `app.Spec` |
| `Runner` | [NewResourceer] → `*runner.Runner` |
| `RepoItem` | [NewResourceer] → `*memory.Repo` |
| `ServiceItem` | `*item.Service` [BuildConfiger] → `item.Spec` |
| `RepoOrder` | [NewResourceer] → `*ordermemory.Repo` |
| `ServiceOrder` | [NewResourceer] → `*order.Service` |
| `Metrics` | [NewResourceer] |
| `ServerHTTPItem` | `*http.Server` [BuildConfiger] → `srvhttp.Config[*http.API]` |
| `APIItem` | [NewResourceer] |
| `ServerHTTPOrder` | `*srvhttp.Config[...]` [BuildConfiger] → same config type |
| `APIOrder` | [NewResourceer] |

### app.Pipeline (set in main)

- `Registry` — e.g. `res.Global()` or `res.New()` in tests
- `EnvPrefix` — ecfg prefix, e.g. `"DEMO"`
- `Transforms` — e.g. `[]res.TransformFunc{obs.ApplyTransform}`; empty skips Transform

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

## Notes

- Stack deps are published modules (`app`, `builder`, `ecfg`, `res`, `sdi`, …) at `v0.20.x`.
- `logger/use` and `telemetry/use` are blank-imported in `main` so `res.Global()` registers defaults.
- External require: `github.com/omcrgnt/proto/gen/go` (srv-http Label/Host/Port).
