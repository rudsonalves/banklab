# Data Layer Agent Guide

This directory owns external integration and app data orchestration.

Use this file as the top-level guide before editing code under
`mobile/lib/data`.

More specific guides:

- `repository/AGENT.md`: repository guidance. Note that production code
  currently lives under `repositories/` plural.
- `services/apis/AGENT.md`: API service and DTO guidance.

## Role In The Architecture

Current architectural flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

The data layer sits between UI-facing view models and core infrastructure. It
turns app operations into API calls, local persistence actions, typed DTOs, and
`Result` values.

## Responsibilities

- Expose repositories consumed by view models
- Implement repository contracts
- Coordinate API services, secure storage, local cache, and lightweight data
  state
- Build and parse API requests/responses through API service classes
- Keep DTOs close to the backend contract that owns them
- Convert transport, backend, and parsing failures into `AppError`
- Return `AsyncResult<T>` instead of throwing raw exceptions through the app

## Folder Map

- `repositories.dart`: registers repository implementations in the dependency injector
- `repositories/`: app-facing repository contracts and implementations
- `repository/`: documentation-only guide requested for repository behavior
- `services/services.dart`: registers API service classes in the injector
- `services/apis/`: endpoint-oriented API services and DTOs
- `services/apis/core/`: shared API parsing helpers and envelope types
- `services/storage/`: reserved for data-specific storage services if needed

## Dependency Rules

Data may depend on:

- `core` contracts and services, such as `Result`, `RestClient`, and
  `LocalSecureStorage`
- `domain` models and enums
- data-layer DTOs and API services

Data must not depend on:

- `ui`
- pages, widgets, themes, or view models
- `BuildContext`
- GoRouter navigation behavior
- Dio directly outside the core HTTP implementation

Repositories may depend on APIs and core services. API services may depend on
`RestClient` and API DTOs. Keep those directions one-way.

## Separation Inside Data

### `services/apis/...`

API services own endpoint-level behavior:

- request paths and HTTP methods
- headers, bodies, and query parameters
- request DTO serialization
- response DTO parsing
- backend envelope handling
- low-level API error mapping

Rules:

- API services should not know about widgets, navigation, or view models.
- API services should not coordinate session state or selected account state.
- API services should not read or write secure storage directly in normal
  feature work.
- For money transport scalars, always convert with `ApiParse`: use
  `ApiParse.toInt` for `Money -> int64` and `ApiParse.toMoney` for
  backend numeric scalar -> `Money`.
- API services should return `AsyncResult<T>`.

### `repositories/...`

Repositories own app-facing data operations:

- combining one or more API calls
- coordinating local persistence
- exposing app-oriented methods like `login`, `logout`, `profile`,
  `listAccounts`, `loadBalance`, and `getStatement`
- maintaining small data state, such as current user, selected account, cached
  profile, or last balance

Rules:

- Repositories should depend on APIs and core services, not widgets.
- Repositories should expose interfaces where the codebase already expects
  injection through contracts.
- Repositories should hide API and storage details from view models.
- Repositories should not parse backend envelopes when an API service should own
  that work.

## Dependency Injection

Data dependencies are registered in two module entrypoints:

- `services/services.dart`: API service registrations
- `repositories.dart`: repository registrations

Current bootstrap order from `core/config/dependencies.dart`:

1. `CoreServices`
2. `Services`
3. `Repositories`
4. `Usecases`
5. `Viewmodels`

When adding a new API or repository:

1. Register the API in `services/services.dart`.
2. Inject the API into the repository implementation.
3. Register the repository in `repositories.dart`.
4. Inject the repository into a view model through the existing UI module.

Do not instantiate APIs or repositories directly in pages.

## Error Handling

Use the shared `Result` model consistently.

Rules:

- Return `Success(value)` or `Result.success(value)` for success.
- Return `Failure(AppError(...))` or `Result.failure(error)` for failure.
- Use `Unit` for successful operations with no meaningful payload.
- Convert parsing failures to `AppErrorCode.parsingError`.
- Convert backend envelope errors or non-success HTTP status responses to
  `AppErrorCode.httpError` unless a more specific code already exists.
- Preserve explicit backend/user-facing messages when the current code does so.
- Do not add another async result pattern beside `Result`.

## DTO And Model Placement

- Request and response DTOs belong under the relevant
  `services/apis/<feature>/dtos/` folder.
- Shared API envelope helpers belong under `services/apis/core/`.
- Stable app concepts belong under `lib/domain`.
- A repository may return DTOs when that is the current convention for the
  feature.
- Prefer domain models when a type represents app meaning beyond one backend
  payload.

DTOs may be used by repositories and view models when they are intentionally
app-facing contracts for this mobile API and already expose idiomatic Dart
fields and app types such as `Money`, `DateTime`, or enums. Do not create a
domain model that merely duplicates a DTO with the same fields and meaning.

Do not leak raw backend envelopes, JSON maps, HTTP status handling, Dio types,
or snake_case transport payloads into UI or view models.

## Current Auth Conventions

Current split:

- `AuthApi` performs auth request/response parsing.
- `AuthRepositoryImpl` stores access and refresh tokens.
- `AuthRepositoryImpl` owns session state such as `currentUser`, `userProfile`,
  and `isLoggedIn`.
- `AuthInterceptor` in `core` appends tokens and handles refresh at transport
  level.

Keep new auth work aligned with this split unless the task explicitly asks for
an architecture refactor.

## Current Account Conventions

Current split:

- account APIs are endpoint-focused: `BalanceApi`, `ListAccountsApi`, and
  `StatementApi`
- `AccountRepositoryImpl` owns selected account state
- `AccountRepositoryImpl` caches last balance and broadcasts balance updates
- statement requests use the currently selected account

Keep API DTO parsing in APIs and selected-account orchestration in the
repository.

## State And Caching

Small repository-owned data state is acceptable when it belongs to app data
behavior.

Examples already present:

- current authenticated user
- cached user profile
- selected account
- last balance cache
- balance stream

Guidelines:

- Keep state private where possible.
- Expose read-only getters or streams when view models need to observe data.
- Avoid storing UI lifecycle, controllers, form values, or navigation state.
- Provide cleanup when a repository owns resources that need disposal.

## Testing

Add or update tests when changing:

- API request construction
- response parsing
- DTO conversion
- repository orchestration
- token persistence or cleanup
- cache semantics
- stream behavior
- error mapping

Use fakes or mocks for `RestClient`, APIs, and storage. Do not make real network
calls in unit tests.

## Do Not

- Do not call Dio directly from repositories, APIs, pages, or view models.
- Do not place widgets, themes, or UI formatting in `data`.
- Do not make repositories depend on route classes or `BuildContext`.
- Do not store tokens in pages or view models.
- Do not bypass repositories by calling API services directly from UI.
- Do not introduce a use case layer unless the task explicitly asks for that
  architecture change.
- Do not perform broad folder migrations while implementing a narrow feature.
