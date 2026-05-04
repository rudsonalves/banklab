# Use Cases Agent Guide

Use cases live in `mobile/lib/domain/usecases`.

Use this folder for application workflows that are too complex for a view model
to coordinate directly, especially when a screen action needs multiple
repositories.

## When To Create A Use Case

Create a use case when:

- a view model would call more than one repository for one user action
- a workflow has multiple steps or branching rules
- the same workflow should be reused by more than one screen
- orchestration is not a UI concern but also does not belong inside one
  repository

## Rules

- Keep use cases free of Flutter widgets and `BuildContext`.
- Inject repositories through constructors.
- Return `Result` or `AsyncResult`.
- Keep UI feedback and navigation in pages.
- Keep API calls inside repositories and API services, not use cases.

