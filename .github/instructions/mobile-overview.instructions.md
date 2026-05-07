---
description: "Use when creating or editing any file under the mobile Flutter app. Covers top-level architecture, layered structure, dependency injection, routing, error handling, and global coding rules for the banklab mobile project."
applyTo: "mobile/**"
---
# Mobile Agent Guide

This directory contains the Flutter app for `banklab`.

Use this file as the top-level instruction set before editing code anywhere under `mobile/`.

## Architecture

The app currently follows a layered structure:

- `lib/core`: cross-cutting infrastructure and app plumbing
- `lib/data`: APIs and repository implementations
- `lib/domain`: app-facing models and enums
- `lib/uis`: pages, view models, themes, and UI building blocks
- `test`: unit tests for core services and adapters

The practical flow in the current codebase is:

`UI -> ViewModel -> Repository -> API/Service -> RestClient`

## Global Rules

- Preserve the current layered architecture. Do not make UI call API classes directly.
- Prefer constructor injection and align new dependencies with `lib/core/config/dependencies.dart`.
- Keep error handling explicit with `Result`, `AppError`, and `Command`. Do not spread raw exceptions through the app.
- Standardize money scalar conversions across the app: use `ApiParse.toInt` for `Money -> int64` and `ApiParse.toMoney` for backend numeric scalar -> `Money`.
- Reuse existing abstractions before creating new ones.
- Match the naming and file placement already used by neighboring code.
- Avoid introducing a new "use case" layer unless the task explicitly includes that refactor. The architecture doc lists it as future work, not current structure.
- Keep imports layer-appropriate:
  - `uis` may depend on `data`, `domain`, and `core`
  - `data` may depend on `domain` and `core`
  - `domain` should stay lightweight and framework-agnostic
- Keep changes focused. Do not mix architectural rewrites into a feature task unless clearly requested.

## Dependency Injection

- Dependency registration is centralized in `lib/core/config/dependencies.dart`.
- Current registration order matters:
  1. `CoreServices`
  2. `Services`
  3. `Data`
  4. `Uis`
- When adding a new repository, API, or view model, register it in the corresponding module entrypoint instead of instantiating it ad hoc in pages.

## Routing

- Routing is handled by GoRouter under `lib/core/routing`.
- Prefer named navigation with the existing route types instead of raw path strings when possible.
- Keep route declarations centralized in the routing layer.

## State and Async Work

- UI-triggered async operations should usually go through `Command0` or `Command1`.
- View models expose commands; pages observe them with `AnimatedBuilder`, listeners, or notifiers.
- Avoid putting async request orchestration directly inside widgets when a view model already owns that flow.

## Error Handling

- Network and parsing failures should become `AppError`.
- Repositories and APIs should return `AsyncResult<T>`.
- UI should present user-facing messages from `AppError.message` and avoid parsing transport payloads itself.

## Style Expectations

- Stay consistent with the current code style instead of introducing a new architecture pattern.
- Keep widgets readable and extracted only when it improves clarity or reuse.
- Prefer explicit, boring code over clever abstractions.

## Before Creating New Files

Check whether the project already has an established place for the change:

- New HTTP contract: `lib/data/services/apis/...`
- New repository behavior: `lib/data/repositories/...`
- New app model/enum: `lib/domain/...`
- New screen/view model: `lib/uis/pages/...`
- New shared UI primitive: `lib/uis/core/...`
- New cross-cutting infra utility: `lib/core/...`

## Testing

- Add or update tests when changing behavior in core services, HTTP adapters, parsers, repositories, or command/state infrastructure.
- Follow the existing direct style in `test/`:
  - small arrange/act/assert tests
  - descriptive test names
  - no unnecessary test harness complexity
