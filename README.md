# demo

In-memory REST CRUD **service template** built with [ecfg](https://github.com/omcrgnt/ecfg), [builder](https://github.com/omcrgnt/builder), [res](https://github.com/omcrgnt/res), [sdi](https://github.com/omcrgnt/sdi), [srv-http](https://github.com/omcrgnt/srv-http), and [runner](https://github.com/omcrgnt/runner).

Layer layout and conventions: **[ARCHITECTURE.md](ARCHITECTURE.md)**.

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

## Test

```bash
go test ./...
```
