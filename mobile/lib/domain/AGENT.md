# Domain Layer Agent Guide

This directory contains app-facing models and enums shared across layers.

Current examples:

- `auth/models/auth_user.dart`
- `auth/models/user_profile.dart`
- `enums/user_role.dart`

## Responsibilities

- Represent stable app concepts
- Provide lightweight parsing or conversion logic when it clearly belongs to the model
- Hold enums and sealed model variants used across repositories and UI

## Keep Domain Lightweight

- Avoid dependencies on Flutter widgets, navigation, Dio, or storage.
- Keep models easy to construct and reason about.
- Prefer immutable-like usage patterns, even if the current code does not use formal immutable tooling.

## What Is Acceptable Here

- Enums
- Sealed class hierarchies like `AuthUser`
- `fromMap` factories when they are simple and already part of the current model style

## What Should Not Go Here

- HTTP request construction
- API envelope parsing concerns
- dependency injection
- widget concerns
- repository orchestration

## Current Model Style

- Models are plain Dart classes.
- Factories such as `fromMap` are acceptable for direct API payload conversion when the project already uses that style.
- Keep naming explicit and aligned with backend contracts when appropriate.

## When Adding New Models

- Put them under a feature-oriented folder if one already exists.
- Prefer small, focused models over oversized transport-shaped blobs.
- If a type is only meaningful for one API request/response, it likely belongs in `data/services/apis/.../dtos`, not here.

