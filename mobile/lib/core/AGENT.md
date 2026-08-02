# Core Layer Agent Guide

This directory contains the shared infrastructure for the Flutter application.

Use these instructions before creating or editing code under `mobile/lib/core`.

## Role In The Architecture

`core` is the application's cross-cutting layer. It provides contracts, adapters,
and plumbing used by `data`, `domain`, and `ui`, but it must not contain
screen-specific rules or feature-specific business flows.

Current architectural flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

The `core` layer supports this flow with:

- dependency injection
- routing
- result and error modeling
- commands for UI-triggered async work
- abstract HTTP client and Dio implementation
- transport interceptors
- secure storage
- environment configuration
- logging and shared extensions

## Folder Map

- `config/`: dependency bootstrap with `AutoInjector`
- `resources/`: global constants and configuration, such as `AppEnv`,
  `StorageKeys`, and currencies
- `result/`: `Result`, `AsyncResult`, `AppError`, `Unit`, and `Command`
- `routing/`: GoRouter setup, typed routes, animations, observer, and codecs
- `services/client_http/`: `RestClient` contract, request/response DTOs, Dio
  implementation, factory, and error mapping
- `services/client_http/interceptors/`: HTTP interceptors, including auth and
  device concerns
- `services/secure_storage/`: secure storage abstraction and adapter
- `services/logging/`: simple logging used by infrastructure
- `extensions/`: shared, feature-independent extensions

## Dependency Rules

- `core` may depend on infrastructure libraries such as Flutter, Dio, GoRouter,
  secure storage, and AutoInjector.
- `core` must not depend on `data`, `domain`, or `ui`.
- Shared contracts belong in `core`; concrete implementations belong in the
  appropriate subdirectory.
- Avoid importing Dio details outside `services/client_http/dio` and
  `services/client_http/interceptors`.

## Dependency Injection

The bootstrap entrypoint is `config/dependencies.dart`.

Current registration order:

1. `CoreServices`
2. `Services`
3. `Repositories`
4. `Usecases`
5. `Viewmodels`

`CoreServices.add` registers:

- `FlutterSecureStorage`
- `LocalSecureStorage`
- main `Dio` instance
- `AuthInterceptor`
- `RestClient` implemented by `DioRestClient`

When adding a new infrastructure dependency:

- register it in `services/core_services.dart` when it is a cross-cutting
  service
- preserve registration order when services depend on each other
- prefer constructors and `injector.get<T>()`; do not instantiate global
  services ad hoc in pages, repositories, or APIs
- keep `injector.commit()` only in the central bootstrap

## Result, AppError, And Command

`Result<T>` is the standard contract for operations that can fail.

Use:

- `AsyncResult<T>` for async methods that return `Future<Result<T>>`
- `Success(value)` or `Result.success(value)` for success
- `Failure(AppError(...))` or `Result.failure(error)` for failure
- `Unit` for operations with no meaningful payload

Rules:

- Do not let raw exceptions cross layers.
- Expected failures must become `AppError` with an appropriate `AppErrorCode`.
- `Command0` and `Command1` are the bridge between UI and async work.
- `Command` must remain generic and infrastructure-agnostic; do not put auth,
  account, screen, route, or widget behavior inside it.
- When changing `Result`, `AppError`, or `Command`, add or update tests.

## HTTP

The public networking contract is `RestClient`.

Expected flow:

1. An API under `data/services/apis/...` builds a `RestClientRequest`
2. `RestClient` executes the call
3. `DioRestClient` converts the response into `RestClientResponse`
4. Dio errors are converted by `mapHttpError`
5. The caller receives `Success` or `Failure(AppError)`

Rules:

- Do not use `Dio` directly outside `core` or infrastructure-specific tests.
- New HTTP methods must be added to the `RestClient` contract first.
- Request and response objects must stay simple, serializable, and independent
  from UI.
- Transport logging belongs in the implementation or interceptors.
- HTTP error mapping must remain centralized in `dio/dio_error_mapper.dart`.

## Interceptors

Interceptors handle transport concerns, never presentation concerns.

`AuthInterceptor` currently:

- adds `AppHttpHeaders.authorization` with bearer value when an access token is
  stored
- skips requests that already include `AppHttpHeaders.authorization`
- attempts refresh on `401` responses, except on the refresh endpoint
- persists new tokens when refresh succeeds
- retries the original request with the new access token
- clears tokens when refresh fails

Guidelines:

- Do not use `BuildContext`, snackbars, or navigation inside an interceptor.
- Do not make interceptors depend on repositories, to avoid dependency cycles.
- Refresh uses a single in-flight operation when multiple requests receive
  `401` at the same time; keep this behavior covered by dedicated interceptor
  tests.
- The Dio instance used for refresh must remain separate enough to avoid
  interceptor recursion.

## Routing

Routing uses GoRouter in `routing/router.dart`.

Current patterns:

- typed/enumerated routes in `routing/routes.dart`
- route groups in `routing/routes/auth_routes.dart` and
  `routing/routes/home_routes.dart`
- `routeObserver` for pages that need route lifecycle awareness
- `AppCustomTransactionPage` for custom transitions where already used
- `ExtraCodec` for `extra` serialization

Rules:

- Centralize new routes under `routing/routes/...`.
- Prefer `goNamed` and route names instead of scattering path strings through
  the UI.
- Pages may receive view models resolved by the injector in route builders.
- Auth guards do not exist yet; do not simulate guards by spreading ad hoc
  checks across multiple pages.

## Storage And Configuration

Secure storage access must go through `LocalSecureStorage`.

Rules:

- Shared keys belong in `resources/storage_keys.dart`.
- Environment configuration belongs in `resources/app_env.dart`.
- `AppEnv` uses `String.fromEnvironment` and should fail fast for invalid
  required configuration, such as `BASE_URL`.
- Do not read or write tokens directly with `FlutterSecureStorage` outside the
  adapter or infrastructure-specific tests.

## Extensions And Utilities

Only put an extension in `core/extensions` when it is genuinely shared and
feature-independent.

Avoid:

- formatting helpers used by a single screen
- business rules disguised as extensions
- utilities that depend on repositories, APIs, or specific widgets

## What Belongs Here

- infrastructure used by multiple features
- cross-cutting contracts
- adapters for external libraries
- architectural helpers, such as result, commands, routing, and config
- small, stable global resources

## What Does Not Belong Here

- screen-specific widgets
- form validation for a single page
- auth or account rules that belong in repositories
- API DTOs
- backend envelope parsing
- feature-specific visual formatting
- direct calls to product endpoints

## When Editing Core

- Be conservative: changes here have broad impact.
- Preserve compatibility when practical.
- Prefer small, additive changes over broad rewrites.
- Check imports to keep `core` independent from `data`, `domain`, and `ui`.
- Update tests when changing Result, Command, HTTP, storage, config, or
  interceptor behavior.
- Run at least the related tests under `mobile/test/core/...` when the change
  touches infrastructure.
