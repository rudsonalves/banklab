---
description: "Use when creating or editing API service files, DTOs, RestClient requests, response parsing, ApiEnvelope, auth API, account APIs. Covers endpoint-level behavior, DTO rules, error mapping, and registration in services.dart."
applyTo: "mobile/lib/data/services/apis/**"
---
# API Services Agent Guide

This guide applies to API service code under `mobile/lib/data/services/apis`.

API services are the lowest app-owned layer above the shared HTTP client. They
build requests, call `RestClient`, parse backend payloads, and return typed
`Result` values to repositories.

## Role In The Architecture

Current architectural flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

API services belong between repositories and `RestClient`.

They should:

- know endpoint paths, methods, headers, bodies, and query parameters
- parse backend envelopes and DTOs
- map transport or parsing failures into `AppError`
- return `AsyncResult<T>`

They should not:

- know about widgets or navigation
- coordinate app sessions or selected account state
- store tokens directly
- expose Dio details to repositories

## Folder Structure

Current layout:

- `auth/`: authentication and profile-related endpoints
- `auth/dtos/`: auth request and response DTOs
- `account/`: account, balance, and statement endpoints
- `account/dtos/`: account request and response DTOs
- `core/`: API parsing helpers and shared backend envelope types

Prefer this layout for new API areas:

```text
data/services/apis/<feature>/<operation>_api.dart
data/services/apis/<feature>/dtos/<operation>_request_dto.dart
data/services/apis/<feature>/dtos/<operation>_response_dto.dart
```

Keep DTOs close to the API that owns the backend contract.

## Dependency Rules

API services may depend on:

- `RestClient`, `RestClientRequest`, and `RestClientResponse`
- `Result`, `AsyncResult`, `AppError`, and `Unit`
- API helpers from `data/services/apis/core`
- DTOs from the same API feature folder
- domain models only when the project already uses them as app-facing API output

API services must not depend on:

- repositories
- view models
- pages or widgets
- `BuildContext`
- GoRouter
- secure storage directly, except for a deliberate infrastructure refactor
- Dio directly

## Request Construction

Use `RestClientRequest` for all HTTP calls.

Guidelines:

- Put endpoint paths in the API service method.
- Use `body` for request payloads.
- Use `queryParameters` for filters and pagination.
- Use `headers` only for endpoint-specific headers.
- Let core interceptors handle shared auth headers when possible.
- Keep request DTO conversion in `toMap()` methods.

Example shape:

```dart
final response = await _client.post(
  RestClientRequest(
    path: '/auth/login',
    headers: {'X-App-Token': AppEnv.appToken},
    body: dto.toMap(),
  ),
);
```

## Response Handling

API services should follow this sequence:

1. Call `RestClient`.
2. If the client result is failure, return `Result.failure(response.error!)`.
3. Read the `RestClientResponse`.
4. Validate successful HTTP status when needed.
5. Parse the backend payload into an `ApiEnvelope<T>` or expected list shape.
6. If the envelope contains an error, return `Failure(AppError(...))`.
7. If required data is missing, return `Failure(AppError(...))`.
8. Return `Success(parsedValue)`.

Do not return raw backend maps to repositories unless the receiving contract is
explicitly designed for raw data.

## ApiEnvelope And API Core Helpers

Shared API helpers live in `apis/core`.

Use:

- `ApiEnvelope<T>` for standard responses shaped like `{ data, error }`
- `ApiError` for backend error payloads
- `ApiParse` for parsing backend scalar values into app types such as `Money`
  or `BigInt`

Rules:

- Extend `apis/core` only for reusable API parsing concerns.
- Do not put feature-specific DTOs in `apis/core`.
- Do not put repository state or UI formatting in `apis/core`.
- For money transport scalars, always use `ApiParse` conversions:
  `ApiParse.toInt` for `Money -> int64` and `ApiParse.toMoney` for
  backend numeric scalar -> `Money`.

## DTO Rules

DTOs represent transport contracts.

Guidelines:

- Keep DTOs immutable in practice with final fields where possible.
- Use `fromMap` for response DTO parsing.
- Use `toMap` for request DTO serialization.
- Keep DTO field names aligned with backend payloads when practical.
- Do not hand-roll money scalar conversions in DTOs; use
  `ApiParse.toInt` and `ApiParse.toMoney`.
- Validate required fields during parsing by failing loudly inside the API
  method's `try/catch`, so parsing failures become `AppErrorCode.parsingError`.
- Avoid leaking DTOs into UI when a stable domain model exists.

## Error Handling

Expected error codes:

- `AppErrorCode.httpError` for unsuccessful HTTP responses or backend envelope
  errors
- `AppErrorCode.parsingError` for invalid or unexpected response shapes
- use other `AppErrorCode` values only when they clearly match the API failure

Rules:

- Wrap parsing in `try/catch`.
- Log parsing failures with `ConsoleLog` when the API already follows that
  pattern.
- Return user-facing messages through `AppError.message`.
- Do not throw raw exceptions from API service methods.
- Do not hide backend error messages when the current API convention preserves
  them.

## Auth API Conventions

`AuthApi` currently owns:

- `/auth/register`
- `/auth/login`
- `profile/me` until profile is split into a dedicated service
- auth request/response parsing
- app token header for auth entrypoints

Keep token persistence out of `AuthApi`; `AuthRepositoryImpl` owns session state
and secure storage coordination.

## Account API Conventions

Account APIs currently split operations by endpoint:

- `BalanceApi`
- `ListAccountsApi`
- `StatementApi`

Keep this style when endpoint behavior differs enough to deserve a focused
class. For small related endpoint groups, follow the nearest existing feature
pattern.

## Registration

When adding an API service:

1. Create the service under `data/services/apis/<feature>/`.
2. Add request/response DTOs under the feature's `dtos/` folder.
3. Register the service in `data/services/services.dart`.
4. Inject it into the relevant repository in `data/data.dart`.

Do not instantiate API services directly inside view models or pages.

## Testing

Add or update tests when changing:

- request path, body, query, or headers
- response parsing
- envelope error handling
- list parsing behavior
- money or numeric parsing
- AppError mapping

Use fake or mock `RestClient` implementations. Do not make real network calls in
unit tests.

## Do Not

- Do not import Dio in API services.
- Do not navigate or show UI feedback here.
- Do not read `BuildContext`.
- Do not store tokens or selected account state here.
- Do not return backend envelope objects to UI.
- Do not duplicate parsing helpers when `apis/core` already has one.
- Do not introduce a second networking abstraction beside `RestClient`.
