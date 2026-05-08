# Use Cases Agent Guide

Use cases live in `mobile/lib/domain/usecases`.

This folder owns application workflows that are too complex for a view model to
coordinate directly. The current reference implementation is
`transfer/transfer_usecase.dart`, registered in `usecases.dart`.

## Role In The Architecture

Current complex flow:

`UI -> ViewModel -> UseCase -> Repository -> API/Service -> RestClient -> Dio`

Use a use case to keep reusable application orchestration out of widgets and
out of repositories that should remain focused on their own data boundary.

## Current Structure

- `usecases.dart`: dependency registration entrypoint for use cases
- `<feature>/<feature>_usecase.dart`: workflow orchestration
- `<feature>/inputs/`: app-facing input objects owned by the use case when a
  repository DTO is not the right API for the UI

Current example:

- `transfer/transfer_usecase.dart`
- `transfer/inputs/transfer_draft.dart`

## When To Create A Use Case

Create a use case when:

- a view model would call more than one repository for one user action
- a workflow has multiple steps or branching rules
- orchestration is not a UI concern but also does not belong inside one
  repository

Keep direct repository injection in the view model when the command simply
delegates to one repository method with minimal mapping.

## Dependency Rules

Use cases may depend on:

- repository contracts
- app-facing DTOs or domain types already used by the surrounding feature
- small framework-independent helpers from `core`, such as `AsyncResult`

Use cases must not depend on:

- Flutter widgets or `BuildContext`
- navigation, snackbars, dialogs, or page lifecycle concerns
- API services or `RestClient`
- secure storage directly
- Dio or transport-layer response objects

## Implementation Rules

- Inject repositories through constructors.
- Register each use case in `usecases.dart`.
- Return `Result` or `AsyncResult`.
- Keep API calls inside repositories and API services, not use cases.
- Keep UI feedback and navigation in pages.
- Keep request mapping inside the use case when it is part of application
  orchestration, such as converting a `TransferDraft` into a repository/API
  request DTO.
- Expose read-only workflow state only when it is needed by the consuming view
  model, such as selected account or available accounts during a transfer flow.

## Do Not

- Do not put `TextEditingController`, forms, or validators here.
- Do not instantiate repositories manually.
- Do not move repository-owned persistence or cache rules into a use case.
- Do not create a use case for every repository method by default.

