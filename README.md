# demo

In-memory REST CRUD demo built with [ecfg](https://github.com/omcrgnt/ecfg), [builder](https://github.com/omcrgnt/builder), [res](https://github.com/omcrgnt/res), [sdi](https://github.com/omcrgnt/sdi), [srv-http](https://github.com/omcrgnt/srv-http), and [runner](https://github.com/omcrgnt/runner).

## Domain

`Item` with fields `id` and `title`, stored in memory.

| Method | Path | Description |
|--------|------|-------------|
| GET | `/items` | List all items |
| GET | `/items/{id}` | Get item by id |
| POST | `/items` | Create item (`{"title":"..."}`) |
| PUT | `/items/{id}` | Update item (`{"title":"..."}`) |
| DELETE | `/items/{id}` | Delete item |

## Configuration

Environment variables (prefix `DEMO`):

```bash
export DEMO_HTTP_SERVER_LABEL=demo
export DEMO_HTTP_SERVER_HOST=127.0.0.1
export DEMO_HTTP_SERVER_PORT=8080
```

See [env.template](env.template) for generated documentation.

Regenerate template:

```bash
go generate ./internal/config/...
```

## Run

```bash
cp .env.example .env
task r-app
```

Build only:

```bash
task b-app
# -> bin/demo
```

## Examples

```bash
# list (empty)
curl -s localhost:8080/items

# create
curl -s -X POST localhost:8080/items \
  -H 'Content-Type: application/json' \
  -d '{"title":"first item"}'

# get (replace ID)
curl -s localhost:8080/items/<id>

# update
curl -s -X PUT localhost:8080/items/<id> \
  -H 'Content-Type: application/json' \
  -d '{"title":"updated"}'

# delete
curl -s -X DELETE localhost:8080/items/<id> -w '\n'
```

## Architecture

`AppConfig` lives in [`internal/config`](internal/config/config.go). Application entrypoint is [`cmd/demo/app.go`](cmd/demo/app.go).

```go
type AppConfig struct {
    Store      store.Config
    Service    service.Config
    Controller httpapi.Config
    Metrics    httpapi.MetricsConfig
    HTTPServer *srvhttp.Config[*httpapi.API]
}
```

1. **ecfg** — `Parse[AppConfig]`: load config from environment
2. **builder** — `Build(cfg, res.Default)`: run each component `Config.Build()`, register resources
3. **sdi** — `Resolve(res.Default)`: wire deps into built resources
4. **runner** — `New(res.Default)` then `Run` / `Stop`

Each component `Config` implements `Build()` (the shared builder contract used by **builder**).

## Test

```bash
go test ./...
```
