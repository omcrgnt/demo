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
| **System** | `meta/core/use` or `*/use` init; library `init` on `unique.Global()` | ops probe/metrics HTTP, HTTPMetrics, default App |

Ops HTTP is **not** in the user catalog: `meta/core/use` pulls in `ops/transport/http/use` → `DefaultServer()` with org defaults (`0.0.0.0:8080`). Domain HTTP ports come from ecfg (`DEMO_SERVER_HTTP_*`).

## Metrics (v0.21)

```text
srv-http (init)  → HTTPMetrics (MetricsContributor + slok Recorder), TagFixed
meta/core/use      → app, logger, telemetry, ops/transport/http/use (probe + metrics + ops HTTP)

Resolve:
  metrics.Actuator.Inject → RegisterMetrics(reg) on HTTPMetrics + domain contributors
  srv-http servers        → shared metrics.Recorder (HTTPMetrics)
  ops Handler             → /livez /readyz /metrics on :8080
```

One registry, one slok recorder, N srv-http domain servers (distinct `Service` label). `srv-grpc` mirrors this with `GRPCMetrics` and grpc-prometheus interceptors.

## Probes (ops v0.25)

`/readyz` on `:8080` — `probe.Actuator` aggregates `ProbeReadiness` from domain `srv-http.Server[T]`, `srv-grpc.Server[T]` (catalog), and ops `transport/http.Server` (system). `/livez` — liveness only.

## HTTP layout

```text
internal/api/http/item/   package item — item REST API (chi)
internal/api/http/order/  package order — order REST API
```

Catalog pairs: `srvhttp.Server[*handler.API]` + `*handler.API` per domain slice.

## gRPC layout

```text
proto/demo/v1/              local OrderService + ProductService protos
internal/api/grpc/gen/demo/v1/  generated stubs (buf)
internal/api/grpc/order/        order gRPC handler
internal/api/grpc/product/      product gRPC handler
internal/api/grpc/bundle/       composite RegisterGRPC (one port, two services)
```

Catalog: `srvgrpc.Server[*bundle.Bundle]` + `Bundle` + `GRPCAPIOrder` + `GRPCAPIProduct` + `ServiceProduct` + `RepoProduct`. Order reuses `ServiceOrder` / `RepoOrder` (HTTP + gRPC share domain layer). Item remains HTTP-only.

## Known gaps (local)

Demo-only gaps: [backlog/items/demo-reference-gaps.md](https://github.com/omcrgnt/backlog/blob/main/items/demo-reference-gaps.md).

## Org backlog

Cross-repo themes live in [github.com/omcrgnt/backlog](https://github.com/omcrgnt/backlog):

| Theme | Item |
|-------|------|
| Drop `builder`, srv-http v0.21 publish | [drop-builder-app-v21](https://github.com/omcrgnt/backlog/blob/main/items/drop-builder-app-v21.md) |
| Configurable / Blueprint naming | [app-catalog-naming](https://github.com/omcrgnt/backlog/blob/main/items/app-catalog-naming.md) |
| ecfg-gen typing, CustomTag | [ecfg-res-custom-tags](https://github.com/omcrgnt/backlog/blob/main/items/ecfg-res-custom-tags.md) |
| sdi CheckCycles opt-in, Many warn | [sdi-v21-followups](https://github.com/omcrgnt/backlog/blob/main/items/sdi-v21-followups.md) |
| ops probe + metrics | [ops-probe-v1-followups](https://github.com/omcrgnt/backlog/blob/main/items/ops-probe-v1-followups.md) |
| srv-http defer listen | [srv-http-defer-listen](https://github.com/omcrgnt/backlog/blob/main/items/srv-http-defer-listen.md) |
| shared Taskfiles | [org-devtools-taskfiles](https://github.com/omcrgnt/backlog/blob/main/items/org-devtools-taskfiles.md) |
| local dev composer | [decompose-local-dev](https://github.com/omcrgnt/backlog/blob/main/items/decompose-local-dev.md) |
