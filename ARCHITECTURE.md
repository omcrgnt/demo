# Architecture

Reference app assembly: one `go.mod`, org libs from `github.com/omcrgnt/*`, explicit pipeline (no legacy `Resourcer`).

## Pipeline

```text
app.Run(&appResources, pipeline)
  → Seed → Apply → Build → Transform → Resolve → App.Serve
```

| Step | Package | Role |
|------|---------|------|
| Seed | `builder` | Walk `AppResources`; register specs (`BuildConfig`) or deferred `NewResource` |
| Apply | `ecfg` | Load env into specs from seed map |
| Build | `builder` | `Spec.Build()` / `NewResource()` → resources in registry |
| Transform | `res` + `obs` | Optional wrappers (metrics, tracing) |
| Resolve | `sdi` | Inject deps from `Deps()` / `Inject()` |
| Serve | `app` + `runner` | Signal → run pool → graceful stop |

## AppResources

Single struct in `main` — catalog of everything in the process.

### Field naming

`{type}{subject}` — e.g. `RepoOrder`, `ServiceItem`, `ServerHTTPItem`.

### Field kinds

Each field is **one resource slot** — either:

- **NewResourceer** — `NewResource()` at Build (no env block)
- **BuildConfiger** — `BuildConfig()` → Spec/Config; env applies to spec; `Build()` → resource

No third type. No duplicate config fields on the resource for ecfg.

```text
AppResources field = resource (*item.Service, *app.App, …)
  BuildConfiger:  resource.BuildConfig() → Spec/Config  →  Spec.Build() → resource
  NewResourceer:  NewResource() → resource
```

### ecfg

- Tagged fields: `ecfg:"BLOCK_NAME"` on the **resource** slot (for block prefix only).
- Env shape comes from **Spec/Config** returned by `BuildConfig()`, not from resource fields.
- Codegen (`ecfg-gen`) resolves spec type via `BuildConfig()` AST; runtime Apply writes into spec in seed map.

Domain suffix in env blocks where relevant: `SERVICE_ITEM`, `SERVER_HTTP_ITEM`, `SERVER_HTTP_ORDER`.

## Domain layout

```text
internal/
  domain/          models, ports, services
  data/sync/       repos (NewResourceer)
  api/http/        HTTP adapters (NewResourceer), thin over domain ports
cmd/demo/          AppResources + main
```

Org stack: `app`, `builder`, `res`, `sdi`, `ecfg`, `runner`, `obs`, `srv-http`, `logger`, `telemetry` — versioned `github.com/omcrgnt/*` modules.

Codegen on services: `sdigen` (Deps/Inject), `obsgen` (Observe).

## Roles (examples)

| Concern | Owner |
|---------|--------|
| Shutdown grace period | `app.Spec` → built `*app.App` |
| Run/stop pool | stateless `runner.Runner` |
| List max length | `item.Spec` → `*item.Service` |
| HTTP listen addr | `srvhttp.Config` via `BuildConfig` on `*http.Server` or `*Config` |

## Known gaps (backlog)

- [ ] **sdi dedup** — pool-wide Replaceable dedup, not only types from `Deps()` stubs
- [ ] **BuildConfig typing** — spec type inferred by AST; fragile for non-literal returns
- [ ] **Symmetry** — order service has no Spec yet; HTTP item uses `http.Server` alias, order uses `srvhttp.Config` directly
- [ ] **Full stack in demo** — logger/telemetry registered via `use` imports; not yet configurable through AppResources

## References

- [README.md](README.md) — commands, AppResources table
- [env.md](env.md) — generated env docs
