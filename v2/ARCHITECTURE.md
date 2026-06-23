# v2 architecture

Spike for the target app assembly: one `go.mod`, org libs under `pkg/`, explicit pipeline (no legacy `Resourcer`).

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
pkg/               reusable org libs (builder, res, sdi, ecfg, runner, obs, …)
```

Codegen on services: `sdigen` (Deps/Inject), `obsgen` (Observe).

## Roles (examples)

| Concern | Owner |
|---------|--------|
| Shutdown grace period | `app.Spec` → built `*app.App` |
| Run/stop pool | stateless `runner.Runner` |
| List max length | `item.Spec` → `*item.Service` |
| HTTP listen addr | `srvhttp.Config` via `BuildConfig` on `*http.Server` or `*Config` |

## Known gaps (spike backlog)

See also `pkg/sdi/KNOWN_ISSUES.md`.

- [ ] **sdi dedup** — pool-wide Replaceable dedup, not only types from `Deps()` stubs
- [ ] **Lib boundaries** — copies in `pkg/` vs published modules (`github.com/omcrgnt/...`)
- [ ] **BuildConfig typing** — spec type inferred by AST; fragile for non-literal returns
- [ ] **Symmetry** — order service has no Spec yet; HTTP item uses `http.Server` alias, order uses `srvhttp.Config` directly
- [ ] **Full stack in demo** — logger/telemetry packages exist but not wired in `AppResources`

## v2 → prod checklist

1. Freeze contracts above (AppResources, resource+Spec, ecfg-on-spec, naming).
2. Close sdi dedup backlog; add regression test for duplicate concretes.
3. Decide module strategy: extract `pkg/*` to versioned repos or monorepo root.
4. Harden ecfg-gen (typed `BuildConfig` return or explicit spec link) if AST becomes a pain.
5. Add one more configurable domain (e.g. `order.Spec`) to prove the pattern copies cleanly.
6. Wire logger/telemetry through AppResources when validating full org stack.
7. Document env blocks per service in runbooks; keep `go generate` in CI for `.env.template` / `env.md`.

## References

- [README.md](README.md) — commands, AppResources table
- [env.md](env.md) — generated env docs
- v1 root (`/opt/github/demo`) — unchanged legacy demo
