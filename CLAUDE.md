<!-- Canonical 6022-labs Go engineering rules. Kept consistent across all agentic Go repos; minor per-repo drift is OK. -->
# CLAUDE.md — agentic-ens-gateway

A CCIP-Read gateway (ERC-3668) that reads ENS data cross-chain (ENS lives on Ethereum/Sepolia only). A 6022-labs Go service — follow the canonical rules below, which are identical across all 6022-labs Go repos. This file is **self-contained**: it depends on no external plugin, skill, or shared file.

## 1. Comments — write fewer, make them sharp
- **One line, max.** Exception: an interface or interface method may use two. Never a multi-line/multi-paragraph doc block.
- **Only the non-obvious** — a *why*, an invariant, a cross-package contract. If the code is clear, write **no** comment (the common case).
- Go doc form: start with the identifier name (`AgentMinter mints ...`), then stop.
- **Document the contract once, on the port** (the `*Interface`). The implementation's methods carry **no** doc comment — it only drifts. The impl keeps a type-level comment and comments on private helpers.
- Assume your first draft over-comments. Cut before presenting.

## 2. Clear boundaries — isolate tier-services, not everything
A **tier-service** is any dependency that crosses a boundary: it reaches an external system, does I/O, or is a side-effecting / cross-cutting infrastructure concern. Each tier-service has **exactly one owner package + an interface**; nothing else imports the underlying library — callers depend on the port and mock it.

**Pure utility libraries are not boundaries** — `samber/lo`, stdlib helpers, math, slices, etc. may be imported anywhere. Don't wrap them.

Litmus: *external system / I/O / side effect / cross-cutting tier* → isolate behind a service + port. *Pure in-memory computation* → use freely.

Hexagonal layering is the structure that enforces this:
- Domain logic lives in the domain module only (`src/<domain>` with `configurations/ services/ services/interfaces/ repositories/ models/ settings/ use_cases/`). Adapters (`*_blockchain/_http_*/_mvc/_db/_vault`, plus each module's `metrics/` + `traces/`) only translate tier-service ↔ domain types.
- **Inbound adapters are transport only — controllers AND listeners.** An HTTP `*_mvc` controller, a blockchain/event listener, or a queue consumer are all just *entry points into the business logic*: decode/deserialize the request/message → call exactly **one** use case → map the result back (status, ack/nack). **No business logic and no validation rules in them.** Message/body/field validation carries meaning → it's a business concern and lives in the use case, never in the transport adapter. Never return framework types from use cases.
- **Use cases are the domain entry points**: return `(*response, error)` (or `error`), depend only on ports (`services/interfaces/`, `repositories/`).
- Domain packages never import `*_host/`; domain-consumed settings live in `<domain>/settings/`.
- Adding behavior = define/extend a port → domain logic in `services/`/`use_cases/` → implement the port in the owner adapter → wire in DI.

Tier-services and their sole owners — depend on the port, not the package:

| Tier-service | Sole owner package | Everyone else uses |
|---|---|---|
| Metrics / traces (OpenTelemetry) | each module's `metrics/` + `traces/` sub-package (+ telemetry bootstrap) | the local `metrics`/`traces` `interfaces/` port |
| Persistence (gorm) | the `*_db` adapter | repository ports; models carry tags only |
| HTTP transport (fiber) | `*_mvc`, `common/mvc`, `*_host` | plain `(*response, error)` from use cases — never framework types |
| Blockchain RPC (go-ethereum) | `*_blockchain` | service/repository ports (address/hash value-types are the only exception) |
| Secrets (Vault) | the `*_vault` adapter | credential-store ports |
| IPFS | `*_http_ipfs` | port |
| LLM / MCP SDKs | the `*_http_inference` / `*_http_mcp` adapters | reasoning ports |

A tier-service import found outside its owner package is a boundary leak — fix it, don't spread it.

## 3. Tests — mandatory, not an afterthought
- **Every new behavior ships with its test in the same change.** Don't hand work back without tests.
- Naming/structure: `when_<behavior>_test.go`; nested `Given.../Should...` `t.Run` blocks; suite `BeforeEach` setup; `go.uber.org/mock` (gomock) + `testify`; reuse `model_builders`.
- Unit tests in `tests/<package>_unit_tests/`. Integration tests in `tests/<package>_integration_tests/` spin real infra (Postgres/Vault/Redis) via the shared `container_builders` helper, need Docker, and stay **out** of the unit-test sweep.

## 4. Mocks — generated, never hand-written
- Regenerate after **any** interface change: `deepsearch-mockgen -S ./src -O ./tests -A -P` (https://github.com/FournyP/deepsearch-mockgen-cli).
- Mock destination mirrors source: `tests/<domain>_mocks/...` (e.g. repository mocks under `tests/<domain>_mocks/repositories_mocks/`, metrics under `tests/<domain>_mocks/metrics_mocks/`).
- Every tier-service port is mocked in unit tests — including metrics/traces. No real telemetry, DB, RPC, or Vault in a unit test.

## 5. Metrics & traces — the tier-service to copy
otel is touched in exactly one place per module: a dedicated `metrics/` (and, where needed, `traces/`) sub-package that wraps the otel meter/tracer and exposes a `…/interfaces` port (e.g. `MintMetrics` → `MintMetricsInterface`). Services, use cases and controllers depend on that **port**, get it injected via dig, and mock it in `tests/<module>_mocks/metrics_mocks/`. **No service, use case, or controller imports `go.opentelemetry.io`.** New metric/trace = add a method to the module's `metrics/`/`traces/` package + its interface, then regenerate mocks. (Constructors fall back to a noop meter so tests and collector-less runs never panic.)

## 6. Cross-cutting plumbing
- **DI**: `go.uber.org/dig`. Each module exposes `Add<Module>Configuration(container)`; wire it in the host's `configurations/di_configuration.go`.
- **Config**: `koanf`. Every setting configurable from **both** `appsettings*.json` **and** env/`.env` (double-underscore lowercased → dotted, e.g. `JWT__SIGNING_KEY` → `jwt.signing_key`). Settings structs in each module's `settings/` with `koanf:"..."` tags + defaults in their `New...` constructor. New setting = tag + default + verified env override.
- **Logging**: `go.uber.org/zap`.
- **Errors → HTTP** (services with controllers): use cases fail with concrete types from `common/errors`, mapped centrally (e.g. `common/mvc.WriteError`): Validation→400, Unauthorized→401, NotFound→404, Conflict→409, Unavailable→502, Internal→500. Controllers pick the success status.
- (Secrets, persistence, and transport are tier-services — see §2.)

## 7. Naming
Name dependency fields/params after their **full type** (drop the package qualifier and `Interface` suffix, camelCase the rest): `EnsPublisherInterface` → `ensPublisher`; `AgentRegistrationStateRepositoryInterface` → `agentRegistrationStateRepository`; `*settings.AgentIdentitySettings` → `agentIdentitySettings`. Idiomatic infra names (`logger`, `client`, `ctx`) and local domain-value variables keep their short names.

---

## Before you hand work back — run this checklist
- [ ] **Tests** added/updated for every new behavior (unit; integration if it touches real infra).
- [ ] **Mocks** regenerated after any interface change.
- [ ] **Boundaries**: no tier-service package imported outside its owner — otel only in `metrics/`/`traces/`, fiber only in `*_mvc`/`*_host`, gorm only in `*_db`, go-ethereum only in `*_blockchain`, Vault only in `*_vault`. Pure utilities (lo, stdlib) are exempt.
- [ ] **Metrics/traces** behind their port and **mocked** in unit tests — no real telemetry in tests.
- [ ] **Layers clean**: no business logic OR validation rules in any inbound adapter (controllers, listeners, consumers) — body/message/field validation lives in the use case; use cases depend only on ports.
- [ ] **Comments cut** to one-line, non-obvious-only; contract docs only on the port.
- [ ] **Settings** reachable from both JSON and env; new settings have tag + default.
- [ ] **Secrets** in Vault, not DB/config rows.
- [ ] `go build ./...` and the unit tests pass.

If any box is unchecked, the work isn't done — finish it before presenting.

---

## Module specifics
_Add anything specific to this repo here (entrypoints, ports, external services, commands). The rules above stay canonical; this section may drift per repo._
