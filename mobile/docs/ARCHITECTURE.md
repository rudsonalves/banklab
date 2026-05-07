# Mobile Architecture

## Table of Contents

- [Mobile Architecture](#mobile-architecture)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Architectural Principles](#architectural-principles)
  - [Project Structure](#project-structure)
  - [Agent Documentation](#agent-documentation)
  - [Copilot Instructions Mirror (.github/instructions)](#copilot-instructions-mirror-githubinstructions)
  - [Dependency Graph](#dependency-graph)
  - [Flow Models](#flow-models)
  - [Layer Responsibilities](#layer-responsibilities)
    - [Core Layer](#core-layer)
    - [Data Layer](#data-layer)
    - [Domain Layer](#domain-layer)
    - [UI Layer](#ui-layer)
  - [Request And Authentication Flow](#request-and-authentication-flow)
  - [Routing Model](#routing-model)
  - [State And Error Model](#state-and-error-model)
  - [Configuration Model](#configuration-model)
  - [Known Constraints And Future Improvements](#known-constraints-and-future-improvements)
  - [Suggested Evolution Path](#suggested-evolution-path)

## Overview

The BankFlow mobile app follows a layered architecture organized by
responsibility. The goal is to keep user interface code simple, isolate
networking and persistence concerns, and make business flows predictable and
testable.

Current source root:

- [mobile/lib](../../mobile/lib)

Main layers:

- UI layer: app shell, screens, widgets, view models, and shared UI primitives
- Domain layer: app-facing models, enums, and use cases for complex workflows
- Data layer: repositories, API services, DTOs, and data orchestration
- Core layer: dependency injection, routing, HTTP client, secure storage,
  result model, commands, config, and shared infrastructure

## Architectural Principles

- Clear separation of concerns
- Constructor-based dependency injection
- Explicit error handling via `Result` types instead of throwing through the app
- Infrastructure details such as Dio and secure storage hidden behind interfaces
- UI code talks to repositories through view models for simple workflows
- UI code talks to use cases through view models for complex workflows that
  coordinate multiple repositories
- API services own endpoint calls and DTO/backend envelope parsing

## Project Structure

- [mobile/lib/main.dart](../../mobile/lib/main.dart): app bootstrap and
  dependency setup
- [mobile/lib/core](../../mobile/lib/core): config, routing, result model,
  commands, HTTP client, secure storage, logging, and platform services
- [mobile/lib/data](../../mobile/lib/data): API services, DTOs, and repository
  implementations
- [mobile/lib/domain](../../mobile/lib/domain): domain root for app-facing
  types and workflow orchestration
- [mobile/lib/domain/common](../../mobile/lib/domain/common): stable
  app-facing models and enums grouped by context
- [mobile/lib/domain/usecases](../../mobile/lib/domain/usecases): use cases for
  complex application workflows
- [mobile/lib/uis](../../mobile/lib/uis): app widget, pages, shared UI
  primitives, themes, and view models
- [mobile/test](../../mobile/test): unit and widget tests

## Agent Documentation

Folder-specific `AGENT.md` files document implementation rules for AI-assisted
code changes:

- [mobile/AGENT.md](../../mobile/AGENT.md): top-level mobile guidance
- [mobile/lib/core/AGENT.md](../../mobile/lib/core/AGENT.md): core layer
- [mobile/lib/data/AGENT.md](../../mobile/lib/data/AGENT.md): data layer
- [mobile/lib/data/repositories/AGENT.md](../../mobile/lib/data/repositories/AGENT.md): repositories
- [mobile/lib/data/services/apis/AGENT.md](../../mobile/lib/data/services/apis/AGENT.md): API services and DTOs
- [mobile/lib/domain/AGENT.md](../../mobile/lib/domain/AGENT.md): domain models and enums
- [mobile/lib/domain/usecases/AGENT.md](../../mobile/lib/domain/usecases/AGENT.md): use cases
- [mobile/lib/uis/AGENT.md](../../mobile/lib/uis/AGENT.md): UI layer
- [mobile/lib/uis/pages/AGENT.md](../../mobile/lib/uis/pages/AGENT.md): pages and view models
- [mobile/lib/uis/core/AGENT.md](../../mobile/lib/uis/core/AGENT.md): shared UI primitives

## Copilot Instructions Mirror (.github/instructions)

The project also includes instruction files under
[.github/instructions](../../.github/instructions) that mirror the guidance
from the mobile `AGENT.md` files. These files are used by coding agents to
apply folder-specific rules automatically via `applyTo` patterns.

Current mobile instruction files:

- [mobile-overview.instructions.md](../../.github/instructions/mobile-overview.instructions.md): top-level rules for all files under `mobile/**`
- [mobile-core.instructions.md](../../.github/instructions/mobile-core.instructions.md): core infrastructure under `mobile/lib/core/**`
- [mobile-data.instructions.md](../../.github/instructions/mobile-data.instructions.md): data layer under `mobile/lib/data/**`
- [mobile-api-services.instructions.md](../../.github/instructions/mobile-api-services.instructions.md): API services and DTO rules under `mobile/lib/data/services/apis/**`
- [mobile-repositories.instructions.md](../../.github/instructions/mobile-repositories.instructions.md): repository rules under `mobile/lib/data/repositories/**`
- [mobile-domain.instructions.md](../../.github/instructions/mobile-domain.instructions.md): domain model rules under `mobile/lib/domain/**`
- [mobile-usecases.instructions.md](../../.github/instructions/mobile-usecases.instructions.md): use case rules under `mobile/lib/domain/usecases/**`
- [mobile-uis.instructions.md](../../.github/instructions/mobile-uis.instructions.md): UI layer rules under `mobile/lib/uis/**`
- [mobile-pages.instructions.md](../../.github/instructions/mobile-pages.instructions.md): page and view model rules under `mobile/lib/uis/pages/**`
- [mobile-uis-core.instructions.md](../../.github/instructions/mobile-uis-core.instructions.md): shared UI primitives under `mobile/lib/uis/core/**`

Maintenance note:

- When updating any mobile `AGENT.md`, update the corresponding instruction
  file in `.github/instructions` to keep both sources aligned.
- Keep `applyTo` patterns specific to avoid leaking rules between layers.

## Dependency Graph

At startup:

1. [mobile/lib/main.dart](../../mobile/lib/main.dart) calls `setupDependencies`
2. [mobile/lib/core/config/dependencies.dart](../../mobile/lib/core/config/dependencies.dart)
   registers all modules
3. [mobile/lib/uis/app_widget.dart](../../mobile/lib/uis/app_widget.dart)
   builds `MaterialApp.router`

Current registration order in the injector:

1. `CoreServices`
2. `Services`
3. `Repositories`
4. `Usecases`
5. `Viewmodels`

This ensures UI view models can resolve their dependencies, repositories can
resolve APIs and platform services, and API services can resolve the shared
`RestClient`.

This also ensures use cases resolve repository dependencies before UI view
models that consume them are constructed.

## Flow Models

Simple user flows:

```text
UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio
```

Complex user flows:

```text
UI -> ViewModel -> UseCase -> Repository -> API/Service -> RestClient -> Dio
```

Use a use case when a view model would otherwise coordinate multiple
repositories, perform a multi-step application workflow, or contain reusable
non-UI orchestration.

Keep direct repository injection for simple commands that call one repository
method and expose the result to the page.

## Layer Responsibilities

### Core Layer

Contains cross-cutting infrastructure:

- Environment and app mode configuration:
  [mobile/lib/core/resources/app_env.dart](../../mobile/lib/core/resources/app_env.dart)
- Dependency bootstrap:
  [mobile/lib/core/config/dependencies.dart](../../mobile/lib/core/config/dependencies.dart)
- Routing setup and route declarations:
  [mobile/lib/core/routing/router.dart](../../mobile/lib/core/routing/router.dart),
  [mobile/lib/core/routing/routes.dart](../../mobile/lib/core/routing/routes.dart)
- Typed result and command execution helpers:
  [mobile/lib/core/result/result.dart](../../mobile/lib/core/result/result.dart),
  [mobile/lib/core/result/command.dart](../../mobile/lib/core/result/command.dart)
- HTTP abstraction and Dio implementation:
  [mobile/lib/core/services/client_http](../../mobile/lib/core/services/client_http)
- Secure storage abstraction and implementation:
  [mobile/lib/core/services/secure_storage](../../mobile/lib/core/services/secure_storage)

Core must not depend on `data`, `domain`, or `uis`.

### Data Layer

Contains integration and persistence orchestration:

- Dependency entrypoint:
  [mobile/lib/data/repositories.dart](../../mobile/lib/data/repositories.dart)
- API services map transport payloads into DTOs or app models
- Repositories implement app-oriented operations and local data state
- DTOs stay close to their owning API service
- API and repository methods return `AsyncResult<T>`

DTOs may cross from repositories into view models when they are intentionally
app-facing contracts for this mobile API and already expose idiomatic Dart
fields and app types. Domain models are introduced when they add meaning,
combine sources, or shield the app from an unstable/non-app-specific contract,
not merely to duplicate DTO fields. Raw JSON maps, backend envelopes, HTTP
status handling, Dio types, and snake_case transport payloads stay inside the
data/API layers.

Auth example:

- API:
  [mobile/lib/data/services/apis/auth/auth_api.dart](../../mobile/lib/data/services/apis/auth/auth_api.dart)
- Repository contract:
  [mobile/lib/data/repositories/auth/auth_repository.dart](../../mobile/lib/data/repositories/auth/auth_repository.dart)
- Repository implementation:
  [mobile/lib/data/repositories/auth/auth_repository_impl.dart](../../mobile/lib/data/repositories/auth/auth_repository_impl.dart)

Account example:

- APIs:
  [mobile/lib/data/services/apis/account/balance_api.dart](../../mobile/lib/data/services/apis/account/balance_api.dart),
  [mobile/lib/data/services/apis/account/list_accounts_api.dart](../../mobile/lib/data/services/apis/account/list_accounts_api.dart),
  [mobile/lib/data/services/apis/account/statement_api.dart](../../mobile/lib/data/services/apis/account/statement_api.dart)
- Repository:
  [mobile/lib/data/repositories/account/account_repository_impl.dart](../../mobile/lib/data/repositories/account/account_repository_impl.dart)

Data may depend on `core` and `domain`, but must not depend on `uis`.

### Domain Layer

Contains domain-centric models, enums, and use cases used across layers.
For growth, the domain root is organized by context under `common/` plus
workflow orchestration under `usecases/`:

- Dependency entrypoint:
  [mobile/lib/domain/usecases/usecases.dart](../../mobile/lib/domain/usecases/usecases.dart)
- [mobile/lib/domain/common/auth/models/auth_user.dart](../../mobile/lib/domain/common/auth/models/auth_user.dart)
- [mobile/lib/domain/common/auth/models/user_profile.dart](../../mobile/lib/domain/common/auth/models/user_profile.dart)
- [mobile/lib/domain/common/user/enums/user_role.dart](../../mobile/lib/domain/common/user/enums/user_role.dart)
- [mobile/lib/domain/common/receipt/enums/transfer_receipt_status.dart](../../mobile/lib/domain/common/receipt/enums/transfer_receipt_status.dart)
- [mobile/lib/domain/usecases](../../mobile/lib/domain/usecases)

Domain models should remain lightweight and framework-agnostic. Use cases should
coordinate application workflows, not UI concerns.

Use case rules:

- Use cases receive repositories through constructors
- Use cases return `Result` or `AsyncResult`
- Use cases do not import Flutter widgets, `BuildContext`, GoRouter, Dio, or
  secure storage directly
- Use cases do not make direct API calls; repositories and API services own that
  work

### UI Layer

Contains presentation and interaction state:

- App shell:
  [mobile/lib/uis/app_widget.dart](../../mobile/lib/uis/app_widget.dart)
- UI dependency registration:
  [mobile/lib/uis/viewmodels.dart](../../mobile/lib/uis/viewmodels.dart)
- Pages and page view models:
  [mobile/lib/uis/pages](../../mobile/lib/uis/pages)
- Shared UI primitives:
  [mobile/lib/uis/core](../../mobile/lib/uis/core)

View model examples:

- [mobile/lib/uis/pages/auth/login/viewmodel/login_viewmodel.dart](../../mobile/lib/uis/pages/auth/login/viewmodel/login_viewmodel.dart)
- [mobile/lib/uis/pages/auth/register/viewmodel/register_viewmodel.dart](../../mobile/lib/uis/pages/auth/register/viewmodel/register_viewmodel.dart)
- [mobile/lib/uis/pages/home/viewmodel/home_viewmodel.dart](../../mobile/lib/uis/pages/home/viewmodel/home_viewmodel.dart)

UI rules:

- Pages own widget tree, controllers, forms, feedback, and navigation
- View models expose `Command0` or `Command1`
- Simple view models may inject repositories directly
- Complex view models should inject use cases instead of coordinating multiple
  repositories directly
- UI must not call API services, `RestClient`, or Dio directly

## Request And Authentication Flow

Simple request flow:

1. Page triggers a view model `Command`
2. Command invokes a repository method
3. Repository calls an API service
4. API service performs a request through `RestClient`
5. `DioRestClient` returns `Success` or `Failure` mapped to `AppError`
6. Command updates state to success/failure and notifies listeners
7. Page reacts with navigation, loading state, or feedback

Complex request flow:

1. Page triggers a view model `Command`
2. Command invokes a use case method
3. Use case coordinates one or more repositories
4. Repositories call API services or local persistence as needed
5. Result flows back through the command
6. Page reacts to command state

Token behavior:

- Access and refresh tokens are stored through `LocalSecureStorage`
- `AuthInterceptor` appends `Authorization` header when available
- On HTTP `401` from a non-refresh endpoint, the interceptor attempts token
  refresh
- If refresh succeeds, the original request is retried
- If refresh fails, session tokens are cleared

Relevant files:

- [mobile/lib/core/services/client_http/interceptors/auth/auth_interceptor.dart](../../mobile/lib/core/services/client_http/interceptors/auth/auth_interceptor.dart)
- [mobile/lib/core/services/client_http/dio/dio_rest_client.dart](../../mobile/lib/core/services/client_http/dio/dio_rest_client.dart)

## Routing Model

Routing is handled by GoRouter:

- Router entry:
  [mobile/lib/core/routing/router.dart](../../mobile/lib/core/routing/router.dart)
- Route enums:
  [mobile/lib/core/routing/routes.dart](../../mobile/lib/core/routing/routes.dart)
- Route groups:
  [mobile/lib/core/routing/routes/auth_routes.dart](../../mobile/lib/core/routing/routes/auth_routes.dart),
  [mobile/lib/core/routing/routes/home_routes.dart](../../mobile/lib/core/routing/routes/home_routes.dart)

Current initial location:

- Login route

## State And Error Model

Asynchronous operations are represented by:

- `Result<T>`: success/failure wrapper
- `AppError`: typed application error
- `Command`: stateful execution wrapper for UI actions

Command states:

- idle
- running
- success
- failure

This pattern keeps side effects explicit and allows pages to react to command
state transitions consistently.

## Configuration Model

Environment values are compile-time defines consumed by `AppEnv`:

- `BASE_URL`
- `CONNECT_TIMEOUT`
- `RECEIVE_TIMEOUT`
- `APP_MODE`
- `APP_ACCESS_TOKEN`

If `BASE_URL` is missing or invalid, app startup fails fast with a `StateError`.

## Known Constraints And Future Improvements

- Auth refresh currently has a known concurrency risk when many requests fail
  with `401` at the same time; multiple refresh attempts may happen
- A refresh lock strategy should be introduced to serialize token refresh
- Profile concerns are currently mixed into `AuthApi` and can be split into a
  dedicated profile API service
- Use cases exist as the preferred place for complex view model orchestration,
  but no production use case implementation has been added yet

## Suggested Evolution Path

- Add concrete use cases when view models start coordinating multiple
  repositories
- Add feature module boundaries for accounts and transactions
- Add navigation guards for authenticated routes
- Expand automated tests around interceptor refresh behavior, repository cache
  semantics, and use case orchestration
