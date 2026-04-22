# Core Layer Agent Guide

This directory contains shared infrastructure and app plumbing.

Examples already present:

- dependency injection setup
- routing and route declarations
- result and command abstractions
- HTTP client interfaces and Dio implementation
- secure storage adapters
- app environment/config values
- logging helpers

## What Belongs Here

- Cross-cutting infrastructure used by multiple features
- App-wide abstractions that should not live in a single repository or page
- Generic helpers that are still architectural, not business-specific

## What Must Not Go Here

- Screen-specific UI logic
- Feature-specific repository orchestration
- Auth/business rules that belong in repositories or models
- Page form validation that is only used in one screen

## Existing Patterns To Preserve

- `Result<T>` and `AppError` are the standard async/error model.
- `Command0` and `Command1` wrap async actions for UI consumption.
- HTTP code flows through `RestClient`, request/response types, and Dio adapters.
- `AppEnv` reads compile-time config and should fail fast when mandatory config is missing.

## Commands

- Use `Command0` and `Command1` for reusable UI async orchestration.
- Keep command classes generic and infra-agnostic.
- Do not embed feature-specific behavior into `Command`.

## HTTP

- New transport concerns should extend the current `client_http` structure rather than bypass it.
- Keep Dio-specific code under the Dio implementation folders.
- Keep request/response objects simple and serializable.
- Interceptors should handle transport concerns, not UI concerns.

## Routing

- Keep route definitions centralized.
- Prefer extending the current route groups instead of scattering route strings in pages.
- Shared navigation helpers should live in routing extensions/utilities, not in random widgets.

## Storage and Config

- Secure storage access should go through the existing abstraction.
- Shared storage keys belong in `resources/storage_keys.dart`.
- Environment and app mode behavior should stay centralized in `resources/app_env.dart`.

## When Editing Core

- Be conservative: this layer has broad impact.
- Prefer additive changes over rewrites.
- If a change affects many features, keep backward compatibility where practical.
- Add tests for behavior changes in this layer.

