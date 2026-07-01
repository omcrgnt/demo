# demo architecture

Reference app on **app v0.21**: explicit `AppResources` catalog, `unique.Global()` registry, blank-import defaults.

## Pipeline

```text
app.Run(&appResources, Pipeline{Registry: unique.Global(), ...})
  fill → ecfg.LoadEnv → materialize → unique.Merge → Transform → sdi.Resolve → Serve
```

## Metrics (v0.21)

```text
srv-http/use     → HTTPMetrics (MetricsContributor + slok Recorder), TagFixed
ops/metrics/use  → *prometheus.Registry + metrics.Actuator
ops/transport/http/use → probe + ops Handler + DefaultServer

Resolve:
  metrics.Actuator.Inject → RegisterMetrics(reg) on HTTPMetrics + domain contributors
  srv-http servers        → shared metrics.Recorder (HTTPMetrics)
  ops Handler             → /livez /readyz /healthz /metrics on OPS_HTTP port
```

One registry, one slok recorder, N srv-http servers (distinct `Service` label).

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
