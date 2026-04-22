# UI Layer Agent Guide

This directory contains presentation code: app shell, pages, view models, themes, and shared UI primitives.

Current examples:

- `app_widget.dart`
- `pages/auth/...`
- `pages/home/...`
- `core/themes/...`
- `core/base/safe_scaffold.dart`
- `core/text_form_field/basic_text_form_field.dart`

## Responsibilities

- Build screens and reusable widgets
- Hold presentation-oriented state in view models
- React to `Command` state for loading, success, and failure
- Handle navigation and user feedback

## Page and ViewModel Split

Follow the current pattern:

- `page.dart` owns widget tree, controllers, form state, snackbars, and navigation
- `viewmodel.dart` exposes dependencies and `Command`s

The current `LoginPage` and `LoginViewModel` are the reference pattern.

## UI Rules

- Pages should not call API classes directly.
- Prefer consuming repositories through injected view models.
- Keep validation close to the form when it is page-specific.
- Reuse shared widgets like `BasicTextFormField` and `SafeScaffold` when they fit the task.
- Use the existing theme and color/text systems instead of inventing one-off styling patterns.

## Async UX

- Drive loading and button disabling from `Command.isRunning`.
- React to results using command listeners and/or `AnimatedBuilder`, consistent with the existing pages.
- Use `ScaffoldMessenger` for lightweight feedback when that matches current behavior.

## Navigation

- Use GoRouter helpers and named routes where the project already does.
- Keep route decisions in UI/presentation code, not repositories.

## ViewModel Rules

- View models should be thin.
- They should mainly wrap repository methods into `Command`s.
- Avoid putting widget-specific concepts like `BuildContext`, `TextEditingController`, or `GlobalKey<FormState>` into view models.

## Shared UI Components

- Put reusable visual primitives under `uis/core/...`.
- Do not promote a widget to shared too early; only extract when reuse or consistency clearly benefits.

## Do Not

- Do not bypass the view model/repository flow for convenience.
- Do not store tokens, make raw HTTP calls, or parse envelopes in UI.
- Do not put app-wide theme or routing logic inside a single page file.
