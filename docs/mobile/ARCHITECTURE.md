# Mobile Architecture

## Overview

The BankFlow mobile app follows a layered architecture organized by responsibility. The goal is to keep user interface code simple, isolate networking and persistence concerns, and make business flows predictable and testable.

Current source root:
- [mobile/lib](../../mobile/lib)

Main layers:
- UI layer: screens, widgets, and view models
- Data layer: repositories and API services
- Core layer: dependency injection, routing, HTTP client, secure storage, result model
- Domain layer: core models and enums

## Architectural principles

- Clear separation of concerns
- Dependency flow from outer layers to inner abstractions
- Explicit error handling via Result types instead of throwing through the app
- Infrastructure details (Dio, secure storage) hidden behind interfaces
- Constructor-based dependency injection

## Project structure

- [mobile/lib/main.dart](../../mobile/lib/main.dart): app bootstrap and dependency setup
- [mobile/lib/core](../../mobile/lib/core): config, routing, result model, platform services
- [mobile/lib/data](../../mobile/lib/data): API services and repository implementations
- [mobile/lib/domain](../../mobile/lib/domain): domain models and enums
- [mobile/lib/uis](../../mobile/lib/uis): app widget, pages, themes, and view models

## Dependency graph

At startup:
1. [mobile/lib/main.dart](../../mobile/lib/main.dart) calls setupDependencies
2. [mobile/lib/core/config/dependencies.dart](../../mobile/lib/core/config/dependencies.dart) registers all modules
3. [mobile/lib/uis/app_widget.dart](../../mobile/lib/uis/app_widget.dart) builds MaterialApp.router

Registration order in the injector:
1. CoreServices
2. Services
3. Data
4. Uis

This ensures UI view models can resolve repositories, and repositories can resolve APIs and platform services.

## Layer responsibilities

### Core layer

Contains cross-cutting infrastructure:
- Environment and app mode configuration: [mobile/lib/core/config/app_env.dart](../../mobile/lib/core/config/app_env.dart)
- Routing setup and route declarations: [mobile/lib/core/routing/router.dart](../../mobile/lib/core/routing/router.dart), [mobile/lib/core/routing/routes.dart](../../mobile/lib/core/routing/routes.dart)
- Typed result and command execution helpers: [mobile/lib/core/result/result.dart](../../mobile/lib/core/result/result.dart), [mobile/lib/core/result/command.dart](../../mobile/lib/core/result/command.dart)
- HTTP abstraction and Dio implementation
- Secure storage abstraction and implementation

### Data layer

Contains integration and persistence orchestration:
- API clients map transport payloads into app DTOs/models
- Repositories implement use-oriented operations and local token persistence

Auth example:
- API: [mobile/lib/data/services/apis/auth/auth_api.dart](../../mobile/lib/data/services/apis/auth/auth_api.dart)
- Repository contract: [mobile/lib/data/repositories/auth/auth_repository.dart](../../mobile/lib/data/repositories/auth/auth_repository.dart)
- Repository implementation: [mobile/lib/data/repositories/auth/auth_repository_impl.dart](../../mobile/lib/data/repositories/auth/auth_repository_impl.dart)

### Domain layer

Contains domain-centric models used across layers:
- [mobile/lib/domain/auth/models/auth_user.dart](../../mobile/lib/domain/auth/models/auth_user.dart)
- [mobile/lib/domain/auth/models/user_profile.dart](../../mobile/lib/domain/auth/models/user_profile.dart)
- [mobile/lib/domain/enums/user_role.dart](../../mobile/lib/domain/enums/user_role.dart)

### UI layer

Contains presentation and interaction state:
- App shell: [mobile/lib/uis/app_widget.dart](../../mobile/lib/uis/app_widget.dart)
- Page routes: auth and home pages
- View models expose Commands for async actions

View model examples:
- [mobile/lib/uis/pages/auth/login/viewmodel/login_viewmodel.dart](../../mobile/lib/uis/pages/auth/login/viewmodel/login_viewmodel.dart)
- [mobile/lib/uis/pages/auth/register/viewmodel/register_viewmodel.dart](../../mobile/lib/uis/pages/auth/register/viewmodel/register_viewmodel.dart)
- [mobile/lib/uis/pages/home/viewmodel/home_viewmodel.dart](../../mobile/lib/uis/pages/home/viewmodel/home_viewmodel.dart)

## Request and authentication flow

1. ViewModel executes a Command
2. Command invokes repository method
3. Repository calls API service
4. API service performs request via RestClient
5. DioRestClient returns Success or Failure mapped to AppError
6. Command updates state to success/failure and notifies listeners

Token behavior:
- Access and refresh tokens are stored through LocalSecureStorage
- AuthInterceptor appends Authorization header when available
- On HTTP 401 (non-refresh endpoint), interceptor attempts token refresh
- If refresh succeeds, original request is retried
- If refresh fails, session tokens are cleared

Relevant files:
- [mobile/lib/core/services/client_http/interceptors/auth/auth_interceptor.dart](../../mobile/lib/core/services/client_http/interceptors/auth/auth_interceptor.dart)
- [mobile/lib/core/services/client_http/dio/dio_rest_client.dart](../../mobile/lib/core/services/client_http/dio/dio_rest_client.dart)

## Routing model

Routing is handled by GoRouter:
- Router entry: [mobile/lib/core/routing/router.dart](../../mobile/lib/core/routing/router.dart)
- Route enums: [mobile/lib/core/routing/routes.dart](../../mobile/lib/core/routing/routes.dart)
- Route groups: [mobile/lib/core/routing/routes/auth_routes.dart](../../mobile/lib/core/routing/routes/auth_routes.dart), [mobile/lib/core/routing/routes/home_routes.dart](../../mobile/lib/core/routing/routes/home_routes.dart)

Current initial location:
- Login route

## State and error model

Asynchronous operations are represented by:
- Result<T>: success/failure wrapper
- AppError: typed application error
- Command: stateful execution wrapper for UI actions

Command states:
- idle
- running
- success
- failure

This pattern keeps side effects explicit and allows pages to react to command state transitions consistently.

## Configuration model

Environment values are compile-time defines consumed by AppEnv:
- BASE_URL
- CONNECT_TIMEOUT
- RECEIVE_TIMEOUT
- APP_MODE

If BASE_URL is missing or invalid, app startup fails fast with a StateError.

## Known constraints and future improvements

- Auth refresh currently has a known concurrency risk when many requests fail with 401 at the same time (multiple refresh attempts may happen)
- A refresh lock strategy should be introduced to serialize token refresh
- Profile concerns are currently mixed into AuthApi and can be split into a dedicated profile API service

## Suggested evolution path

- Introduce dedicated use case layer between UI and repositories
- Add feature module boundaries for accounts and transactions
- Add navigation guards for authenticated routes
- Expand automated tests around interceptor refresh behavior and repository caching semantics
