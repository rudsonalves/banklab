---
description: "Use when creating or editing files under mobile/lib/ui: app_widget.dart, viewmodels.dart, pages, view models, themes, feedback helpers, shared UI components. Covers the UI/ViewModel split, Command usage, navigation, routing, shared widgets, and dependency rules for the Flutter presentation layer."
applyTo: "mobile/lib/ui/**"
---
# UI Layer Agent Guide

This directory contains the presentation layer for the Flutter application:
app shell, pages, view models, themes, feedback helpers, shared UI components,
and page-level widgets.

More specific guides:

- `mobile/lib/ui/pages/**` — `mobile-pages.instructions.md`: page, route lifecycle, and page view model guidance
- `mobile/lib/ui/components/**` — `mobile-ui-components.instructions.md`: shared/transversal UI widget and visual primitive guidance

## Role In The Architecture

Current architectural flow:

Simple flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

Complex flow:

`UI -> ViewModel -> UseCase -> Repository -> API/Service -> RestClient -> Dio`

The UI layer is the only layer that should know about Flutter widgets,
`BuildContext`, form controllers, visual feedback, and navigation decisions.
It should consume repositories through view models and should not bypass the
data layer.

## Current Structure

- `app_widget.dart`: root `MaterialApp.router`, theme selection, and app shell
- `viewmodels.dart`: registers UI view models in the dependency injector
- `pages/`: screens, page-local widgets, and page view models
- `core/`: shared UI components, input formatters, feedback helpers, and themes

Current page areas:

- `pages/auth/...`
- `pages/home/...`

Current shared UI areas:

- `core/base/...`
- `core/feedback/...`
- `core/input_formatters/...`
- `core/text_form_field/...`
- `core/themes/...`

## Responsibilities

- Build screens and reusable presentation widgets
- Own widget lifecycle state in pages
- Hold presentation-oriented orchestration in view models
- Trigger async work through `Command0` and `Command1`
- React to command state for loading, success, and failure
- Handle navigation and user feedback
- Use the app theme and shared UI components consistently

## Dependency Rules

UI may depend on:

- Flutter and Material widgets
- GoRouter for navigation
- `core/routing` route enums and helpers
- `core/result/command.dart`
- use cases from `domain/usecases` when a workflow is too complex for a view
  model to coordinate directly
- repositories through injected view models
- DTOs needed to submit user input to repository-backed commands
- domain models or app-facing DTOs exposed by repositories

UI must not depend on:

- API services
- `RestClient`
- Dio
- secure storage directly
- backend envelope parsing
- low-level HTTP response objects

Do not call APIs directly from pages or view models. Keep the flow through
view models and repositories or through use cases when the workflow requires
coordination across multiple repositories.

## Page And ViewModel Split

Pages own Flutter-specific concerns:

- widget tree
- `BuildContext`
- form keys
- text controllers
- focus handling
- snackbars, dialogs, and user feedback
- navigation
- listener setup and disposal
- route lifecycle subscription

View models own presentation-facing coordination:

- injected repositories for simple workflows
- injected use cases for complex workflows
- `Command0` and `Command1` instances
- read-only getters derived from repositories
- simple lifecycle methods when needed by pages

Reference pattern:

- `LoginPage` owns form state, validation, snackbars, and navigation.
- `LoginViewModel` exposes `login` as a `Command1`.

Do not put `BuildContext`, `TextEditingController`, `GlobalKey<FormState>`, or
widget-specific state into view models.

## Use Cases

Use cases live under `mobile/lib/domain/usecases`.

Use a use case when a view model would otherwise need to:

- coordinate multiple repositories
- perform a multi-step app workflow
- combine repository results into one app action
- make branching decisions that are not UI concerns
- keep orchestration that should be reusable across screens

Keep simple view models simple. If a command only delegates to one repository
method with minimal presentation mapping, injecting the repository directly is
still acceptable.

Use case rules:

- Keep use cases independent from Flutter widgets and `BuildContext`.
- Return `Result`/`AsyncResult` just like repositories.
- Inject repositories into use cases through constructors.
- Inject use cases into view models through constructors.
- Register use cases in the dependency setup when they are introduced.

## Commands And Async UX

UI-triggered async work should usually go through `Command0` or `Command1`.

Rules:

- Drive loading indicators and disabled buttons from `command.isRunning`.
- Use `AnimatedBuilder`, `ListenableBuilder`, or explicit command listeners to
  react to state changes.
- Show user-facing errors from `command.error?.message`.
- Keep success/failure side effects in pages.
- Keep repository calls wrapped by view model commands when the action belongs
  to a simple user workflow.
- For complex workflows, wrap use case calls in view model commands instead of
  making the view model coordinate multiple repositories directly.

## Navigation

Routing is handled by GoRouter in the core routing layer.

Guidelines:

- Prefer named navigation, such as `context.goNamed(...)`.
- Use route enums from `core/routing/routes.dart`.
- Keep route decisions in pages or presentation flow.
- Add route declarations in `core/routing/routes/...`, not in page files.
- Do not make repositories or API services navigate.

## Shared UI Primitives

Use `ui/components` for UI elements with real cross-screen potential.

`ui/components` is also the staging area for a possible future internal mobile
widget/feature package. Treat additions here as part of the app's shared UI
surface: they should be reusable, presentation-only, themed, and independent of
feature workflow decisions.

Current examples:

- `SafeScaffold`
- `AppSnackbar`
- `BasicTextFormField`
- `CpfInputFormatter`
- `MaterialTheme`
- `createTextTheme`

Guidelines:

- Keep page-specific widgets under `pages/<feature>/.../widgets`.
- Promote widgets to `ui/components` only when reuse or consistency is clear.
- Shared widgets should receive values and callbacks; they should not own app
  workflow decisions.
- Use `Theme.of(context)`, `ColorScheme`, and `TextTheme` instead of one-off
  styling.

## Forms And Feedback

For forms:

- keep controllers and `GlobalKey<FormState>` in pages
- keep page-specific validation close to the form
- create request DTOs after validation passes
- dispose controllers in `dispose`

For feedback:

- use `AppSnackbar.show(...)` for the current convention
- show errors from `AppError.message`
- keep feedback decisions in pages, not view models

## Registration

View models are registered in `ui/viewmodels.dart`.

Rules:

- Register each view model with its injected repositories or use cases.
- Resolve view models in route builders using the central injector.
- Do not instantiate repositories inside pages or view models directly.

## Testing

Add or update tests when changing:

- view model command behavior
- shared input formatters or UI primitives
- route lifecycle behavior such as refresh/polling start and stop

Prefer unit tests for view models and widget tests for screens or shared
components.

## Do Not

- Do not bypass the view model/repository flow for convenience.
- Do not store tokens in UI.
- Do not make raw HTTP calls from UI.
- Do not parse backend envelopes in UI.
- Do not put API services behind button callbacks.
- Do not put app-wide route declarations inside page files.
- Do not put page-specific widgets in `ui/components` before reuse is clear.
- Do not put widget lifecycle objects inside view models.
