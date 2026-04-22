# Data Layer Agent Guide

This directory owns API integration and repository implementation.

Current examples:

- `services/apis/...` for transport contracts and remote calls
- `repositories/...` for app-facing orchestration and local persistence

## Responsibilities

- Talk to external APIs through the `RestClient`
- Parse API envelopes and DTOs
- Convert transport responses into domain/app models
- Coordinate local token persistence and caching when needed
- Return `AsyncResult<T>` instead of throwing through the app

## Separation Inside Data

### `services/apis/...`

Use this area for:

- HTTP requests
- headers, path, body, query parameters
- API envelope parsing
- low-level response mapping

Rules:

- Services should not know about widgets or navigation.
- Services should not depend on view models.
- Keep DTOs close to the API that uses them.
- DTOs are transport structures; do not leak them into the UI if a domain/app model already exists.

### `repositories/...`

Use this area for:

- combining API calls and local storage
- exposing app-oriented methods like `login`, `logout`, `profile`
- caching or memoization that belongs to app behavior

Rules:

- Repositories are the place to coordinate remote calls plus local persistence.
- Repositories should depend on APIs and core services, not on widgets.
- Repositories should expose interfaces where the codebase already expects them.

## Error Handling

- Convert failures to `Failure(AppError(...))`.
- Preserve explicit error messages for user feedback when the current code does so.
- Parsing failures should become `AppErrorCode.parsingError`.
- Transport failures should stay in the `Result` model, not become uncaught exceptions.

## Current Auth Conventions

Based on the current auth implementation:

- `AuthApi` performs request/response parsing.
- `AuthRepositoryImpl` stores access and refresh tokens.
- `AuthRepositoryImpl` owns simple session state like `currentUser`, `userProfile`, and `isLoggedIn`.

Keep new code aligned with this split unless the task explicitly asks for refactoring.

## DTO and Model Placement

- Request/response DTOs go under the relevant API folder.
- Domain/app models stay under `lib/domain`.
- If the backend response already matches an app model closely, mapping directly into a domain model is acceptable, as done with `LoggedUser.fromMap`.

## Do Not

- Do not call `dio` directly from pages or view models.
- Do not place feature widgets or UI formatting here.
- Do not make repositories depend on route classes or `BuildContext`.
- Do not add a second async/result pattern beside `Result`.

