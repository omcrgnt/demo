# demo architecture

Reference app on **app v0.21**: `_appResources` catalog in `cmd/app`, `unique.Global()` registry, blank-import system defaults.

## Pipeline

```text
app.Run(&appResources, Pipeline{Registry: unique.Global(), ...})
  fill → ecfg.LoadEnv → materialize → unique.Merge → Transform → sdi.Resolve → Serve
```

## Catalog vs system modules

| Layer | Where | Example |
|-------|--------|---------|
| **User catalog** | `cmd/app` `_appResources` | domain servers, handlers, services, repos |
| **System** | `*/use` init on `unique.Global()` | ops probe/metrics HTTP, HTTPMetrics, default App |

Ops HTTP is **not** in the user catalog: `ops/transport/http/use` registers `DefaultServer()` with org defaults (`0.0.0.0:8080`). Domain HTTP ports come from ecfg (`DEMO_SERVER_HTTP_*`).

## Metrics (v0.21)

```text
srv-http/use     → HTTPMetrics (MetricsContributor + slok Recorder), TagFixed
ops/metrics/use  → *prometheus.Registry + metrics.Actuator
ops/transport/http/use → probe + ops Handler + DefaultServer (:8080)

Resolve:
  metrics.Actuator.Inject → RegisterMetrics(reg) on HTTPMetrics + domain contributors
  srv-http servers        → shared metrics.Recorder (HTTPMetrics)
  ops Handler             → /livez /readyz /healthz /metrics on :8080
```

One registry, one slok recorder, N srv-http domain servers (distinct `Service` label).

## HTTP layout

```text
internal/api/http/item/   package item — item REST API (chi)
internal/api/http/order/  package order — order REST API
```

Catalog pairs: `srvhttp.Server[*handler.API]` + `*handler.API` per domain slice.

## Known gaps (local)

Demo-only gaps: [backlog/items/demo-reference-gaps.md](https://github.com/omcrgnt/backlog/blob/main/items/demo-reference-gaps.md).

## Org backlog

Cross-repo themes live in [github.com/omcrgnt/backlog](https://github.com/omcrgnt/backlog):

| Theme | Item |
|-------|------|
| Drop `builder`, srv-http v0.21 publish | [drop-builder-app-v21](https://github.com/omcrgnt/backlog/blob/main/items/drop-builder-app-v21.md) |
| Configurable / Blueprint naming | [app-catalog-naming](https://github.com/omcrgnt/backlog/blob/main/items/app-catalog-naming.md) |
| ecfg-gen typing, CustomTag | [ecfg-res-custom-tags](https://github.com/omcrgnt/backlog/blob/main/items/ecfg-res-custom-tags.md) |
| sdi DependencyOrder, Many warn | [sdi-v21-followups](https://github.com/omcrgnt/backlog/blob/main/items/sdi-v21-followups.md) |
| ops probe + metrics | [ops-probe-v1-followups](https://github.com/omcrgnt/backlog/blob/main/items/ops-probe-v1-followups.md) |
| srv-http defer listen | [srv-http-defer-listen](https://github.com/omcrgnt/backlog/blob/main/items/srv-http-defer-listen.md) |
| shared Taskfiles | [org-devtools-taskfiles](https://github.com/omcrgnt/backlog/blob/main/items/org-devtools-taskfiles.md) |
| local dev composer | [decompose-local-dev](https://github.com/omcrgnt/backlog/blob/main/items/decompose-local-dev.md) |
