# Statement Transaction Details Tasks

These tasks split the statement-to-details flow into implementation-sized
steps.

The goal is to let users open the account statement from the home screen and
tap an eligible transaction to open the existing `DetailsPage` through
`SharedRoutes.details`.

## Task 1/8: Add statement route

### Goal

Create a dedicated route for the statement page.

### Scope

- Add a statement route enum/path using the existing route organization style.
- Register the route in the app router.
- Keep the route separate from `SharedRoutes.details`.
- Use the current selected account from repository/application state rather than
  passing account data through route extras.

### Acceptance Criteria

- The app has a named statement route.
- The statement route can be opened from navigation.
- The route is registered in the main router.
- The route does not require custom route extras for the first version.

### Depends On

- None.

## Task 2/8: Add statement view model

### Goal

Provide state and actions for loading account statement items.

### Scope

- Add a statement page view model.
- Inject `AccountRepository`.
- Use `AccountRepository.getStatement()`.
- Load statement for the current selected account.
- Expose loading, success, empty, and failure states through the existing
  command/result pattern where possible.
- Expose the loaded `StatementItemDto` list to the page.
- Do not call API services directly from the view model.

### Acceptance Criteria

- Statement loading goes through `AccountRepository`.
- The view model does not instantiate repositories or API services.
- The page can read loaded statement items.
- Empty statement results are distinguishable from failures.
- Loading can be retried.

### Depends On

- Task 1.

## Task 3/8: Register statement view model in DI

### Goal

Make the statement view model available to the statement route.

### Scope

- Register the statement view model in the existing view model dependency setup.
- Resolve the view model from the route builder.
- Keep dependency construction out of widgets.

### Acceptance Criteria

- The statement route obtains its view model through DI.
- The view model receives the existing `AccountRepository`.
- No page manually constructs the statement view model.

### Depends On

- Task 2.

## Task 4/8: Build statement page shell and states

### Goal

Create the statement page UI structure and state rendering.

### Scope

- Add a statement page.
- Trigger statement loading when the page starts.
- Show loading state.
- Show empty state when there are no statement items.
- Show error state with retry.
- Show the selected account context when useful.
- Keep the UI consistent with the existing mobile app style.

### Acceptance Criteria

- Statement page renders without data while loading.
- Empty state is clear and non-error.
- Failure state shows a retry action.
- The page does not fetch data directly from API services.

### Depends On

- Task 3.

## Task 5/8: Render statement transaction items

### Goal

Display statement transactions in a readable banking-style list.

### Scope

- Render each `StatementItemDto`.
- Display transaction type.
- Display amount.
- Display creation date.
- Display resulting balance when useful.
- Visually distinguish credit and debit movements.
- Keep item layout stable on mobile widths.
- Avoid exposing internal-only identifiers as primary user-facing text.

### Acceptance Criteria

- Each transaction item shows enough information to identify the movement.
- Amounts are formatted consistently with the app.
- Dates are formatted consistently with the app.
- Credit and debit movements are visually distinguishable.
- Long text does not overflow or overlap.

### Depends On

- Task 4.

## Task 6/8: Navigate from statement item to details

### Goal

Connect eligible statement transactions to the existing details route.

### Scope

- Make statement items with a non-empty `referenceId` tappable.
- Navigate to `SharedRoutes.details`.
- Pass `referenceId` as route extra.
- Do not create a new details page.
- Do not bypass `DetailsViewmodel`.
- Do not navigate when `referenceId` is missing or blank.
- Show clear feedback or disabled state for transactions without details.

### Acceptance Criteria

- Tapping an item with `referenceId` opens `DetailsPage`.
- `DetailsPage` receives the selected reference.
- Transactions without `referenceId` do not crash or navigate silently.
- The existing transfer receipt details flow remains unchanged.

### Depends On

- Task 5.

## Task 7/8: Harden shared details route input

### Goal

Prevent crashes when the shared details route is opened without a valid
transaction reference.

### Scope

- Review `SharedRoutes.details`.
- Replace unsafe `state.extra as String` behavior with safe validation.
- Handle missing, blank, or non-string extras.
- Route to a safe fallback or render a friendly invalid-reference state.
- Preserve existing transfer success navigation behavior.

### Acceptance Criteria

- Opening `SharedRoutes.details` without a valid string does not crash.
- Opening details with a blank string does not attempt receipt loading.
- Existing transfer success navigation still works.
- Statement navigation still works.

### Depends On

- Task 6.

## Task 8/8: Validate statement-to-details flow

### Goal

Verify the end-to-end user journeys.

### Scope

- Run focused mobile static checks.
- Validate home-to-statement navigation.
- Validate statement loading success.
- Validate statement empty state.
- Validate statement failure and retry.
- Validate details navigation for item with reference.
- Validate unavailable-details behavior for item without reference.
- Validate receipt load failure behavior in `DetailsPage`.

### Acceptance Criteria

- Home opens statement:

```text
Home -> Statement
```

- Statement loads transactions:

```text
Statement -> AccountRepository.getStatement() -> transaction list
```

- Eligible transaction opens details:

```text
Statement item tap -> SharedRoutes.details -> DetailsPage
```

- Transaction without reference is handled gracefully:

```text
Statement item without reference -> details unavailable feedback
```

- Mobile analysis/checks pass.

### Depends On

- Task 7.
