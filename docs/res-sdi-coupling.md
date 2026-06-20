# ADR: res + sdi coupling

Design record for registry ↔ wiring integration. Chosen approach: **variant 5 + import axis A**.

## Context

- **res** — resource registry (`Registry`), tags on `Entry` (`TagReplaceable` for library defaults via `AddWithTags`)
- **sdi** — dependency graph, topo-sort, `Inject`
- **demo** — composition root: `ecfg.Register` → `builder.Build(res)` → `sdi.Resolve(res.Default)`

Interface ports (e.g. logger `Output`) may have a **Replaceable** default + **explicit** override as **different concretes** until Resolve.

## Import axis (compile-time glue)

| ID | Link | Glue |
|----|------|------|
| A | duck, demo only | `runnerPool` adapter over `res.Registry` for `runner.Pool` |
| B | sdi → res | `Resolve(reg res.Registry)` |
| C | res → sdi | none (res does not import sdi) |
| D | both → `pool` | shared tiny module |
| E | one package / `core` module | merge |

**Chosen: A** for app↔runner; **sdi imports res** for `Registry` and `Entry.Replaceable`.

## Dedupe axis

| # | API | Semantics | Retry |
|---|-----|-----------|-------|
| 1 | `WalkEntries` only | none | — |
| 2 | `Remove(pair)` | res | outer `for` |
| 3 | `Dedup(candidates)` | res | outer `for` |
| 4 | `Dedup(ifaces)` | res | upfront |
| **5** | **`applyPolicy` + `DefaultDedupPolicy`** | **sdi** | **upfront** |
| 6 | import res | varies | varies |
| 7 | merge packages | internal | varies |

**Chosen: 5** — policy and behavior tests live in **sdi**; **res** lists implementors and executes `Remove` via `Registry.Remove`.

## Resolve pipeline (variant 5)

1. `collectDeps(reg)` — concrete and interface stubs from all `Deps()`
2. `cleanupConcretes(reg, concreteTypes, DefaultDedupPolicy)` — dedupe by exact type via `GetByType`
3. `validateInterfaces(reg, ifaces, DefaultDedupPolicy)` — dedupe by interface via `GetByInterface`
4. `wire(reg)` — topo-sort → `Inject`

### `DefaultDedupPolicy` (sdi)

| implementors for port I | action |
|-------------------------|--------|
| 0, 1 | ok |
| 2 | exactly one `Replaceable` → `Remove` it; 0 → ambiguous; 2 → multiple replaceable |
| ≥3 | too many implementations |

### Concrete stubs

Exact type match; dedupe via `cleanupConcretes` with the same policy. Dedupe rules apply to **both** concrete and interface stubs.

## Semantic drift

Behavior changes version with **sdi** (policy). **res** executor (`Remove`) stays stable. Guarded by:

- table tests on `DefaultDedupPolicy` (sdi)
- `applyPolicy` executor tests (sdi)
- demo `config_test` integration (`TestAppConfig_Resolve`)

## Lint (app)

Forbid `AddWithTags(..., TagReplaceable)` in app code — system modules only. Only `sdi.Resolve` triggers dedupe.

## References

- [ARCHITECTURE.md](../ARCHITECTURE.md) — app pipeline
- Plan: `sdi_res_variant_5` (implementation)
