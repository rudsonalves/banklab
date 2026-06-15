# Pages Agent Guide

This guide applies to screens, page-level widgets, feature widgets, and page
view models under `mobile/lib/ui/pages`.

Pages are the presentation entrypoints for user workflows. They should compose
widgets, own widget lifecycle state, trigger view model commands, react to
command results, and perform navigation or user feedback.

## Role In The Architecture

Current architectural flow:

Simple flow:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

Complex flow:

`UI -> ViewModel -> UseCase -> Repository -> API/Service -> RestClient -> Dio`

Pages belong at the UI edge of the app. They should never skip directly to API
services or HTTP clients.

## Folder Structure

Current page layout:

- `auth/login/login_page.dart`
- `auth/login/viewmodel/login_viewmodel.dart`
- `register/register_cpf_page.dart`
- `register/register_name_page.dart`
- `register/register_birthdate_page.dart`
- `register/register_email_page.dart`
- `register/register_token_page.dart`
- `register/register_phone_page.dart`
- `register/register_password_page.dart`
- `register/register_status_page.dart`
- `register/viewmodel/register_viewmodel.dart`
- `home/home_page.dart`
- `home/viewmodel/home_viewmodel.dart`
- `home/widgets/...`
- `transfer/models/...`

Prefer this layout for new pages:

```text
pages/<feature>/<page_name>/<page_name>_page.dart
pages/<feature>/<page_name>/viewmodel/<page_name>_viewmodel.dart
pages/<feature>/<page_name>/models/<presentation_model>.dart
pages/<feature>/<page_name>/widgets/<local_widget>.dart
```

For simple feature areas that already use a flatter structure, follow the
nearest existing pattern instead of reorganizing files.

## Page And ViewModel Split

Pages own Flutter concerns:

- widget tree
- `BuildContext`
- form keys
- text controllers
- focus handling
- animation/listener subscriptions
- snackbars, dialogs, and lightweight user feedback
- navigation
- route lifecycle handling

View models own presentation-facing orchestration:

- injected repositories for simple workflows
- injected use cases for complex workflows
- `Command0` and `Command1` instances
- read-only getters that expose repository state
- simple lifecycle methods when needed, such as starting or stopping refresh
  behavior

Do not put `BuildContext`, `TextEditingController`, `GlobalKey<FormState>`, or
widget references inside view models.

## Dependencies

Pages may depend on:

- Flutter widgets and material components
- GoRouter navigation helpers
- route enums from `core/routing/routes.dart`
- page view models
- feature-local presentation models from `pages/<feature>/.../models`
- shared UI components from `ui/components`
- DTOs needed to submit user input to view model commands

View models may depend on:

- repositories
- use cases from `domain/usecases`
- `Command`, `Result`, and `Unit`
- feature-local presentation models when they represent UI workflow state rather
  than domain or backend contracts
- DTOs needed by repository methods
- domain models or app-facing DTOs exposed by repositories

Pages and view models must not depend on:

- API services
- `RestClient`
- Dio
- secure storage directly
- backend envelope parsing

## Presentation Models

Use a feature-local `models/` folder under `ui/pages` for small types that
belong to the presentation flow rather than to a use case or backend contract.

Good examples:

- route `extra` payloads that need typed serialization
- confirmation or review screen data assembled from form input and lookup DTOs
- immutable UI snapshots used to keep later screens independent from mutable
  selection state
- page-flow state that combines display fields with values later submitted to a
  command

Current example:

- `transfer/models/transfer_confirmation_data.dart`

Keep these models out of `domain/usecases/inputs` unless they are the actual
input contract consumed by a use case. For example, `TransferDraft` belongs to
the transfer use case because it represents the data needed to execute the
transfer. `TransferConfirmationData` belongs to the transfer UI flow because it
contains display-oriented recipient fields and route/confirmation data.

Keep these models out of `data/services/apis/.../dtos` unless they represent a
backend request or response shape.

## Use Cases In View Models

Use cases live under `mobile/lib/domain/usecases`.

Use a use case when a view model would otherwise coordinate several repositories
or contain a multi-step workflow that is not purely presentation logic.

Good use case triggers:

- one screen action needs data from multiple repositories
- one command performs several repository calls in sequence
- multiple screens should reuse the same app workflow
- the view model starts making non-UI branching decisions
- tests for the view model become mostly orchestration setup

Keep direct repository injection for simple commands that call one repository
method and expose the result to the page.

Use case rules:

- Inject repositories into the use case, not into the page.
- Inject the use case into the view model.
- Wrap use case execution in `Command0` or `Command1`.
- Keep use cases free of `BuildContext`, controllers, widgets, snackbars, and
  navigation.
- Return `Result`/`AsyncResult` so pages can keep reacting through command state.

## Commands And Async UI

Use `Command0` and `Command1` for async actions triggered by UI.

Current patterns:

- `LoginViewModel.login` wraps `AuthRepository.login`
- `RegisterViewmodel` coordinates the multi-step register use case and its
  final `register` command wraps `RegisterUsecase.register`
- `HomeViewmodel.initialize` wraps account loading behavior
- More complex view models should wrap use case methods instead of coordinating
  multiple repositories directly.

Page rules:

- Disable submit buttons while the relevant command is running.
- Drive loading indicators from `command.isRunning`.
- Use `AnimatedBuilder`, `ListenableBuilder`, or listeners to react to command
  state.
- Show errors from `command.error?.message`.
- Avoid launching the same command twice while it is already running.
- Keep async request orchestration out of button callbacks when a view model
  command already owns it.

## Forms And Validation

For page-specific forms:

- keep `GlobalKey<FormState>` in the page state
- keep `TextEditingController`s in the page state
- dispose controllers in `dispose`
- keep simple validators close to the form
- create request DTOs in the submit method after validation passes
- unfocus before submitting when it improves keyboard behavior

Move validation out of the page only when it becomes shared behavior or a domain
rule.

## Navigation

Routing uses GoRouter.

Guidelines:

- Prefer named routes such as `context.goNamed(...)`.
- Use route enums from `core/routing/routes.dart`.
- Keep route decisions in pages or presentation flow, not repositories.
- Add new routes in `core/routing/routes/...`, not inside page files.
- Do not scatter raw path strings through widgets.

## Feedback

Use the existing feedback helpers when they fit the behavior.

Current convention:

- `AppSnackbar.show(...)` for success, error, and informational messages
- command listeners decide when to show feedback after async operations

Guidelines:

- Show user-facing messages from `AppError.message` when available.
- Keep backend parsing details out of UI messages.
- Avoid showing snackbars from view models or repositories.

## Route Lifecycle

Use `RouteAware` and `routeObserver` when a page needs route lifecycle events.

The current home page uses route lifecycle callbacks to start and stop periodic
balance refresh behavior through its view model.

Guidelines:

- Subscribe in `didChangeDependencies` when a route is available.
- Unsubscribe in `dispose`.
- Keep timer ownership in the view model or repository only when that ownership
  is already clear.
- Stop refresh or polling work when the page is no longer visible.

## Local Widgets

Put page-specific widgets under the page or feature `widgets/` folder.

Guidelines:

- Use local widgets for readability when the page becomes too dense.
- Promote widgets to `ui/components` only when reuse across features is clear.
- Keep local widgets presentation-only.
- Do not make local widgets call repositories or APIs.

## Styling

Follow the app's existing Material theme and shared UI components.

Use:

- `SafeScaffold` for scaffold consistency when it fits
- `BasicTextFormField` for text inputs when it fits
- theme colors and text styles from `Theme.of(context)`
- existing icon and spacing conventions from nearby pages

Avoid one-off styling systems inside page files.

## Registration And Construction

View models are registered in `ui/viewmodels.dart` and resolved in route builders.

Rules:

- Inject repositories through view model constructors for simple workflows.
- Inject use cases through view model constructors for workflows that coordinate
  multiple repositories or contain reusable application orchestration.
- Resolve view models in route builders with the central injector.
- Do not instantiate repositories inside pages.
- Do not create global singleton view models manually.

## Testing

Add or update tests when changing:

- view model command behavior
- form validation logic
- navigation decisions
- command listener side effects
- lifecycle-driven refresh or polling behavior
- widget rendering for important states

Prefer focused widget tests for pages and unit tests for view models.

## Do Not

- Do not call API services directly from pages or view models.
- Do not import Dio or `RestClient` in UI code.
- Do not store tokens or parse backend envelopes in UI.
- Do not put route declarations inside page files.
- Do not put app-wide theme changes in a single page.
- Do not put widget controllers or `BuildContext` in view models.
- Do not promote every small widget to shared UI before it is reused.
