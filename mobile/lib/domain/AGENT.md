# Domain Layer Agent Guide

This directory contains stable app-facing models, enums, and lightweight domain
types shared across layers.

Use these instructions before creating or editing code under
`mobile/lib/domain`.

## Role In The Architecture

Current architectural flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

The domain layer sits outside transport and presentation details. It defines
concepts the app can reason about regardless of where the data came from.

Current examples:

- `common/auth/models/auth_state.dart`
- `common/auth/models/user_profile.dart`
- `common/user/enums/user_role.dart`
- `common/receipt/enums/transfer_receipt_status.dart`

Current domain root structure:

- `common/`: stable app-facing models and enums grouped by area/context
- `usecases/`: domain workflow orchestration use cases

## Responsibilities

- Represent stable app concepts
- Provide enums and small type hierarchies used across layers
- Keep app-facing models easy to construct, compare, and reason about
- Provide lightweight parsing or conversion only when it clearly belongs to the
  model and matches existing style

## Dependency Rules

Domain may depend on:

- Dart standard library
- other domain models or enums
- small, framework-independent helpers from `core` when already established

Domain models under `common/` must not depend on:

- Flutter widgets or UI packages
- GoRouter or navigation
- Dio or HTTP request/response types
- secure storage
- backend envelope types

Use cases under `usecases/` are the exception:

- they may depend on repository contracts
- they may be registered through `usecases/usecases.dart`
- they must still stay free of Flutter widgets, navigation, and transport
  clients

Keep domain as framework-agnostic as practical, especially under `common/`.

## Current Model Style

The current domain style uses plain Dart classes and enums:

- `AuthState` is a sealed class hierarchy.
- `OperationalAuthState`, `RestrictedInstallationAuthState`, and
  `AnonymousAuthState` model authentication state.
- `UserProfile` is a plain model with required fields.
- `UserRole` is an enum with a `byName` factory fallback.

Guidelines:

- Prefer final fields.
- Prefer explicit constructors with required fields.
- Keep models small and focused.
- Avoid code generation unless the project explicitly adopts it.
- Keep naming aligned with app meaning first, backend names second.

## Parsing In Domain Models

Simple `fromMap` factories are acceptable when:

- the project already uses that style for the model
- the payload maps directly to a stable app concept
- parsing does not require backend envelope handling
- the conversion is small and obvious

Do not put these concerns in domain models:

- HTTP status handling
- request construction
- API envelope parsing
- endpoint-specific list parsing
- transport-only DTO conversion
- storage reads or writes

If parsing logic is only meaningful for one endpoint payload, put the type under
`data/services/apis/<feature>/dtos/` instead.

## Domain Models Vs DTOs

Use a domain model when the type represents a stable app concept, for example:

- authenticated user state
- user profile
- role or permission concepts
- app-level account or transaction concepts, if added later

Use a DTO when the type represents a backend contract, for example:

- login request payload
- registration request payload
- one endpoint response shape
- query parameters
- transport-specific field grouping

Repositories and APIs may map DTOs into domain models when that improves the
boundary for the UI.

Do not create domain models only for architectural purity. If a DTO is already
an app-facing contract from this mobile API, uses idiomatic Dart names and app
types, and matches what view models need, it can remain the type crossing the
repository/view-model boundary. Add a domain model when it adds behavior,
combines multiple sources, hides an unstable/non-app-specific contract, or
represents meaning beyond one endpoint payload.

## Auth Domain Conventions

Current auth domain conventions:

- `AuthState` represents the current authentication subject, including
  anonymous, operational, and restricted installation states.
- `OperationalAuthState` contains token-bearing login response data currently
  needed by the repository.
- `RestrictedInstallationAuthState` contains restricted authorization data for
  the installation registration flow and must not be treated as an operational
  session.
- `AuthRepository.login` may return either an operational state or a restricted
  installation state. Only the operational state may persist session tokens.
- Installation certification promotes a restricted state to
  `OperationalAuthState` only after step-up authorization and successful
  registration.
- `AnonymousAuthState` is the anonymous/default state.
- `UserProfile` represents profile details fetched after login.
- `UserRole.byName` maps unknown role strings to `UserRole.none`.

Keep session orchestration out of domain. The repository owns current user
state, profile caching, and token persistence.

## Folder Placement

Keep stable app-facing concepts grouped under `domain/common` to avoid growing
the `domain` root with many small top-level feature folders.

```text
domain/common/<area>/models/<model>.dart
domain/common/<area>/enums/<enum>.dart
```

Rules:

- Put broad or shared app enums under the closest `domain/common/<area>/enums`.
- Put app-facing models under the closest `domain/common/<area>/models`.
- Keep workflow orchestration under `domain/usecases`.
- Keep `domain` top-level limited to `common` and `usecases` unless a migration
  is explicitly requested.
- Do not create catch-all files for unrelated models.
- Avoid moving existing models unless the task explicitly asks for a migration.

## Error Handling

Domain models should not own application error handling.

Guidelines:

- Use simple Dart exceptions only inside parsing when invalid input should fail
  fast and be caught by the API layer.
- Convert parsing failures into `AppError` in API services or repositories, not
  in domain models.
- Do not import `AppError` into domain models unless there is a deliberate
  architecture decision to make domain produce app errors.

## Testing

Add or update tests when changing:

- enum parsing/fallback behavior
- `fromMap` parsing
- sealed model variants
- equality or value semantics, if introduced later
- required field behavior that other layers rely on

Keep tests small and independent from HTTP, storage, and widgets.

## Do Not

- Do not put API request DTOs in domain.
- Do not parse backend envelopes here.
- Do not make domain models under `common/` call repositories or services.
- Do not add widget, routing, storage, or Dio dependencies.
- Do not store mutable UI state in domain models.
- Do not turn `domain/usecases` into a second repository or transport layer.
