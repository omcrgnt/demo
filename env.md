# Environment configuration

Copy `.env.template` to `.env` and set values. See `.env.example` for local defaults.

**Prefix:** `DEMO`

## OPS_HTTP

| Variable | Description |
|----------|-------------|
| `DEMO_OPS_HTTP_HOST` | URL или IP |
| `DEMO_OPS_HTTP_LABEL` | Метка сущности: по шаблону ^[a-z][a-z0-9_]{1,31}$ |
| `DEMO_OPS_HTTP_PORT` | Порт сервера: 1-65535 |

## SERVER_HTTP_ITEM

| Variable | Description |
|----------|-------------|
| `DEMO_SERVER_HTTP_ITEM_HOST` | URL или IP |
| `DEMO_SERVER_HTTP_ITEM_LABEL` | Метка сущности: по шаблону ^[a-z][a-z0-9_]{1,31}$ |
| `DEMO_SERVER_HTTP_ITEM_PORT` | Порт сервера: 1-65535 |

## SERVER_HTTP_ORDER

| Variable | Description |
|----------|-------------|
| `DEMO_SERVER_HTTP_ORDER_HOST` | URL или IP |
| `DEMO_SERVER_HTTP_ORDER_LABEL` | Метка сущности: по шаблону ^[a-z][a-z0-9_]{1,31}$ |
| `DEMO_SERVER_HTTP_ORDER_PORT` | Порт сервера: 1-65535 |

## SERVICE_ITEM

| Variable | Description |
|----------|-------------|
| `DEMO_SERVICE_ITEM_MAX_LIST_LEN` | Maximum items returned by List (0 = default 100) |
