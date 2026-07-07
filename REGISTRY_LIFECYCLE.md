# Registry lifecycle (NewResourceer / BuildConfiger)

> **НИ В КОЕМ СЛУЧАЕ НЕ ПРАВИТЬ**
>
> Зафиксированная модель pipeline: `builder` с симметрией через `newResourceSpec` → **Builder** в pool; `res` — strict `Add` (`MustAddToGlobalWithTags`, только `NewResourceer` / `BuildConfiger` на `use.init`).
>
> Колонки **NewResourceer** / **BuildConfiger** — два параллельных пути; значения — что лежит в **pool** на этом этапе.

## Таблица

| Фаза | Действие | NewResourceer | BuildConfiger |
|------|----------|---------------|---------------|
| **`res`** | создать пустой `Global` | пусто | пусто |
| **`use.init`** | `MustAddToGlobalWithTags` | **NewResourceer** | **BuildConfiger** |
| **`builder.Seed`** (1) | обернуть в `newResourceSpec` / `BuildConfig()` | **Builder** | **Builder** |
| **`builder.Seed`** (2) | `Add` | **Builder** | **Builder** |
| **`ecfg.Apply`** | env → **Builder** в SeedMap | **Builder** | **Builder** |
| **`builder.Build`** | `Builder.Build()`; `Add` **resource**; `Remove` **Builder** | **resource** | **resource** |
| **`res.Transform`** | `Transform` | **resource** | **resource** |
| **`sdi.Resolve`** | `Inject` | **resource** | **resource** |

На **`builder.Seed`** (1) в колонке «Действие» — *как* получили **Builder**; в обеих колонках типов уже **Builder** (`newResourceSpec` — concrete-тип, реализующий `builder.Builder`).
