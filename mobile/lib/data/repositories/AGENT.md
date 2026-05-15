# Repository Layer Agent Guide

This guide applies to repository code in the mobile data layer.

## Role In The Architecture

Repositories are the app-facing boundary of the data layer. They coordinate API
services, local persistence, lightweight cache/state, and conversion between
transport details and the rest of the app.

Current architectural flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

Repositories are consumed by view models and should hide API and storage details
from the UI.

## Responsibilities

- Expose app-oriented operations, such as `login`, `logout`, `profile`,
  `listAccounts`, `loadBalance`, and `getStatement`
- Coordinate one or more API services
- Coordinate local persistence when needed, such as secure token storage
- Keep small in-memory state or cache when it belongs to data behavior
- Return `AsyncResult<T>` instead of throwing raw exceptions through the app
- Convert lower-level failures into the shared `Result` model

## Current Patterns

Repository contracts and implementations live together by feature:

- auth: `data/repositories/auth/auth_repository.dart` and
  `auth_repository_impl.dart`
- account: `data/repositories/account/account_repository.dart` and
  `account_repository_impl.dart`

Prefer this layout for new repositories:

```text
data/repositories/<feature>/<feature>_repository.dart
data/repositories/<feature>/<feature>_repository_impl.dart
```

## Contracts And Implementations

Use a repository interface when the surrounding code already follows that
pattern or when the repository is injected into view models.

Guidelines:

- Keep the interface focused on app use cases, not HTTP endpoints.
- Name methods after user/app operations, not transport details.
- Keep implementation details in `*RepositoryImpl`.
- Register implementations in `data/repositories.dart`.
- Inject repositories into view models through `AutoInjector`.

## Dependency Rules

Repositories may depend on:

- API services from `data/services/apis/...`
- core services such as `LocalSecureStorage`
- `Result`, `AppError`, and `Unit` from `core/result`
- domain models when a stable app model exists
- DTOs when the current app-facing contract still uses transport-shaped data

Repositories must not depend on:

- Flutter widgets
- `BuildContext`
- GoRouter route classes
- view models
- pages or shared UI components
- Dio directly

## Error Handling

Repositories should keep failures explicit.

Rules:

- Return `Success(...)` for completed operations.
- Return `Failure(AppError(...))` for repository-owned validation failures.
- Pass through API failures with `Result.failure(result.error!)` when no extra
  mapping is needed.
- Do not catch and swallow unexpected errors without converting them to
  `AppError`.
- Do not introduce a second async/error abstraction beside `Result`.

Use `Unit` for operations that complete without a meaningful payload, such as
`logout`, `register`, or `loadBalance`.

## Auth Repository Conventions

`AuthRepositoryImpl` currently owns:

- current user session state
- cached user profile
- login state via `isLoggedIn`
- access and refresh token persistence
- logout token cleanup

Keep this split:

- `AuthApi` performs auth HTTP calls and response parsing.
- `AuthRepositoryImpl` coordinates session state and secure storage.
- UI should only interact with auth through `AuthRepository` and view models.

Do not move token storage into pages, view models, or API classes unless the
architecture is intentionally refactored.

## Account Repository Conventions

`AccountRepositoryImpl` currently owns:

- selected account state
- last balance cache
- balance stream updates
- account listing orchestration
- statement requests for the selected account

Keep this split:

- account APIs perform endpoint calls and DTO parsing
- the repository chooses or remembers the selected account
- the repository coordinates balance refresh behavior that belongs to data state
- UI reads state through the view model instead of calling APIs directly

## State And Caching

Small repository-owned state is acceptable when it represents app data behavior,
for example:

- current authenticated user
- cached profile
- selected account
- last known balance
- broadcast streams for data changes

Guidelines:

- Keep state private and expose only the minimum needed surface.
- Avoid storing widget state, form state, or navigation state here.
- If a repository creates streams or timers, provide a clear lifecycle strategy.
- Prefer simple cache invalidation over hidden background behavior.

## DTOs And Domain Models

Use the existing boundary intentionally:

- Request and response DTOs belong under `data/services/apis/.../dtos`.
- Stable app models belong under `domain`.
- Repositories may return DTOs when that is the current local convention for a
  feature.
- Prefer domain models when the type represents app meaning beyond a single
  endpoint payload.

DTOs are acceptable repository outputs when the backend contract is designed for
the app and the DTO is already a curated app-facing type. Avoid adding a domain
model that simply mirrors the DTO without changing meaning, behavior, or
stability.

Do not leak low-level transport envelopes, raw maps, HTTP status handling, Dio
types, or snake_case backend payloads into UI or view models.

## Registration

When adding a repository:

1. Create the contract and implementation under `data/repositories/<feature>/`.
2. Add constructor dependencies for required APIs or core services.
3. Register the repository in `data/repositories.dart`.
4. Inject the repository into the relevant view model in `ui/viewmodels.dart` or via
   the existing constructor injection pattern.

Never instantiate repository implementations directly inside pages.

## Testing

Add or update tests when repository behavior changes:

- token persistence or cleanup
- cache semantics
- selected account behavior
- error mapping
- API orchestration
- stream updates

Keep tests focused on repository behavior. Mock or fake APIs and storage instead
of making real HTTP calls.

## Do Not

- Do not call Dio directly from repositories.
- Do not parse backend envelopes in repositories when an API service should own
  that work.
- Do not show snackbars, dialogs, or navigate from repositories.
- Do not pass `BuildContext` into repositories.
- Do not store form controllers, focus nodes, or widget lifecycle state here.
- Do not create feature-wide rewrites while adding a narrow repository method.
