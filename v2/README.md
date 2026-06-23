# demo/v2 — AppResources spike

Isolated sandbox for the target architecture: single `go.mod`, org libs under `pkg/`, pipeline without `Resourcer`.

See [ARCHITECTURE.md](ARCHITECTURE.md) for contracts, roles, and v2 → prod checklist.

## Pipeline

Handled by [`pkg/app`](/opt/github/demo/v2/pkg/app):

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

- `Registry` — e.g. `res.Default` or `res.New()` in tests
- `EnvPrefix` — ecfg prefix, e.g. `"DEMO"`
- `Transforms` — e.g. `[]res.TransformFunc{obs.ApplyTransform}`; empty skips Transform

## Commands

From `demo/v2/`:

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

- v1 root (`/opt/github/demo`) is unchanged.
- Origin repos (`github.com/omcrgnt/builder`, etc.) are not modified — copies live in `pkg/`.
- External require: `github.com/omcrgnt/proto/gen/go` (srv-http Label/Host/Port).
