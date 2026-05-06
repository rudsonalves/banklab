---
description: "Use when creating or editing shared UI core files: SafeScaffold, AppSnackbar, BasicTextFormField, CpfInputFormatter, MaterialTheme, input formatters, layout shells, or theme helpers. Covers promotion rules, widget design, feedback helpers, and dependency constraints for uis/core."
applyTo: "mobile/lib/uis/core/**"
---
# Shared UI Core Agent Guide

This guide applies to reusable widgets, input elements, feedback helpers,
formatters, and theme primitives under `mobile/lib/uis/core`.

Use this area for UI elements with real potential to be reused across multiple
screens. Keep page-specific widgets under `mobile/lib/uis/pages/.../widgets`
until reuse is clear.

## Role In The Architecture

Current architectural flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

`uis/core` belongs to the presentation layer. It provides shared visual and
interaction building blocks for pages, but it must not know about repositories,
API services, storage, or networking.

## Current Areas

- `base/`: layout shells such as `SafeScaffold`
- `feedback/`: shared user feedback helpers such as `AppSnackbar`
- `input_formatters/`: reusable text input formatters such as
  `CpfInputFormatter`
- `text_form_field/`: reusable input widgets such as `BasicTextFormField`
- `themes/`: Material theme and text theme helpers

## What Belongs Here

Put code in `uis/core` when it is:

- reusable across more than one page or likely to become app-wide
- presentation-only
- independent of a specific feature workflow
- configurable through constructor parameters
- aligned with the app's Material theme
- useful for consistency, accessibility, or repeated interaction patterns

Examples:

- shared scaffold wrappers
- shared text fields
- shared snackbars or feedback helpers
- input formatters
- theme helpers
- generic buttons, tiles, loading states, or empty states when reused

## What Does Not Belong Here

Do not put these in `uis/core`:

- widgets used by only one page with no clear reuse
- feature-specific cards, rows, or sections
- widgets that call repositories, APIs, or view models directly
- route-specific navigation behavior
- form validation that belongs to one page
- backend DTO parsing or formatting
- token, session, or storage logic

Keep local feature widgets under `pages/<feature>/.../widgets`.

## Dependency Rules

Shared UI core may depend on:

- Flutter and Material widgets
- shared UI theme primitives
- small formatting utilities when they are presentation-focused
- core extensions only when they are framework-appropriate and already shared

Shared UI core must not depend on:

- repositories
- API services
- `RestClient`
- Dio
- secure storage
- `BuildContext` outside widget/build/helper contexts where Flutter requires it
- page view models
- route-specific enums unless the component is explicitly navigation-focused

Prefer data-in/data-out widgets. Pages should provide callbacks and values.

## Promotion Rules

Before moving a widget from a page into `uis/core`, check:

- Is it needed by at least two screens, or is reuse very likely?
- Can it be named without referencing one feature?
- Can feature-specific text, icons, callbacks, and values be parameters?
- Does it avoid importing page or repository code?
- Does it improve consistency without hiding important page behavior?

If the answer is no, keep it local to the page.

## Widget Design Guidelines

Shared widgets should:

- have clear, small constructor APIs
- expose required values explicitly
- accept callbacks instead of owning workflow decisions
- use `Theme.of(context)` and existing `ColorScheme`/`TextTheme`
- support enabled, disabled, loading, empty, and error states when relevant
- avoid hard-coded copy unless it is genuinely app-wide
- avoid hidden side effects
- remain easy to test in isolation

Do not make shared widgets too generic too early. Prefer a small component that
matches current app needs over a large configuration surface.

## Form Inputs And Formatters

Use `text_form_field/` for shared input widgets and `input_formatters/` for
formatting rules.

Guidelines:

- Keep controllers owned by pages, not shared widgets.
- Shared fields should receive controllers and callbacks as parameters.
- Page-specific validators should stay on the page.
- Reusable formatters should not know about API payloads or repositories.
- Preserve cursor behavior carefully when editing input formatters.
- Add tests for non-trivial formatter behavior.

Current examples:

- `BasicTextFormField`
- `CpfInputFormatter`

## Feedback Helpers

Use `feedback/` for app-wide feedback patterns.

Current example:

- `AppSnackbar.show(...)`

Guidelines:

- Feedback helpers may use `BuildContext` because they operate in Flutter UI.
- Pages should decide when feedback is shown.
- View models and repositories should not show feedback directly.
- Keep feedback messages passed in by callers unless the message is truly
  reusable across the whole app.

## Layout And Base Widgets

Use `base/` for layout wrappers that define repeated screen structure.

Current example:

- `SafeScaffold`

Guidelines:

- Keep base widgets flexible enough for different pages.
- Avoid embedding feature content or route decisions.
- Respect safe areas and common spacing conventions.
- Do not create nested scaffold patterns unless a page explicitly needs them.

## Themes

Theme code lives under `themes/`.

Guidelines:

- Prefer app-wide theme changes in theme files, not individual pages.
- Use Material `ColorScheme` and `TextTheme` consistently.
- Avoid one-off colors in shared widgets unless there is no theme token.
- Keep generated or large theme structures stable unless the task is theme work.
- When adding shared widgets, style from the theme instead of introducing local
  design systems.

## Naming And Files

Use descriptive names based on generic UI purpose:

- `SafeScaffold`
- `BasicTextFormField`
- `AppSnackbar`
- `<Thing>InputFormatter`

Avoid feature names in `uis/core` unless the component is truly part of the
whole app brand or design language.

Keep one main widget/helper per file unless tightly coupled private helpers make
the file clearer.

## Testing

Add or update tests when changing:

- input formatter behavior
- shared widget enabled/disabled/loading states
- feedback helper behavior that can be verified safely
- layout constraints that many pages rely on
- theme helpers that affect app-wide rendering

Prefer widget tests for shared components and small unit tests for pure
formatters.

## Do Not

- Do not call repositories or APIs from shared UI components.
- Do not import Dio, `RestClient`, or secure storage here.
- Do not put page-specific copy or workflow logic in shared widgets.
- Do not move a widget here only to reduce one page file's length.
- Do not hide navigation decisions inside generic UI elements.
- Do not bypass the app theme with scattered hard-coded styles.
