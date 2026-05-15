# Task: Open transaction details from statement

## Problem Statement

The mobile app has the infrastructure to retrieve account statements and already
has a shared `DetailsPage` capable of presenting a transfer receipt by
transaction reference. However, the statement entry point is not available from
the home screen yet, and statement transactions are not connected to the details
route.

In a banking app, users expect to inspect a transaction from the statement to
confirm amount, date, status, counterparty, description, and transaction
reference. The statement list should therefore become an actionable flow where
tapping an eligible transaction opens the existing details page.

## Goal

Add a mobile statement flow that lets the user open account transactions and
tap a transaction to navigate to the existing `DetailsPage` through the shared
details route.

## Existing Context

- `StatementApi` already retrieves statements for an account.
- `StatementItemDto` includes:

```dart
transactionId
type
amount
balanceAfter
referenceId
createdAt
```

- `DetailsPage` currently receives a transaction reference and loads the
  transfer receipt through `DetailsViewmodel`.
- `SharedRoutes.details` already routes to `DetailsPage` using `state.extra as
  String`.
- The transfer success flow already opens `DetailsPage` with a transaction
  reference.

## Non-Goals

- Do not create a new transaction details page.
- Do not replace `DetailsPage`.
- Do not add a new details endpoint unless the current receipt endpoint cannot
  support the selected transaction.
- Do not implement advanced statement filtering in this backlog.
- Do not implement pagination unless it is required to render the first useful
  version.
- Do not support details for transaction types that do not have a usable
  reference yet.

## Product Flow

1. The user opens the statement entry point from the home screen.
2. The app loads the statement for the selected account.
3. The user sees a list of transactions.
4. The user taps an eligible transaction.
5. The app navigates to `SharedRoutes.details`.
6. `DetailsPage` loads and presents the transaction details using the reference
   passed by the statement item.

## Routing Decision

Use the existing shared details route:

```dart
context.pushNamed(
  SharedRoutes.details.name,
  extra: transaction.referenceId,
);
```

Only transactions with a non-empty `referenceId` should navigate to details.
Transactions without a reference should either be disabled or show a clear
"details unavailable" message.

## Epic 1: Statement Entry Point

Enable the user to open the statement flow from the home screen.

### Scope

- Replace the current pending behavior for the "Extrato" action.
- Add or register a statement route.
- Navigate from the home action tile to the statement route.
- Keep the selected account as the statement source.

### Acceptance Criteria

- The "Extrato" action opens the statement flow.
- The statement flow uses the currently selected account.
- The home screen no longer shows "coming soon" for statement access.

## Epic 2: Statement State and Loading

Create the state needed to load and present statement items.

### Scope

- Add a statement page view model or extend the existing home/account state in a
  focused way.
- Use `AccountRepository.getStatement()`.
- Load statement items for the selected account.
- Represent loading, success, empty, and failure states.
- Preserve the current account/balance behavior.

### Acceptance Criteria

- Statement loading calls the existing repository method.
- Loading state is visible.
- Empty state is visible when there are no transactions.
- Failure state is visible and retryable.
- Statement data is not fetched directly from the page widget.

## Epic 3: Statement Page UI

Build the statement list UI.

### Scope

- Add a statement page if one does not exist yet.
- Display transaction type.
- Display amount.
- Display creation date.
- Display resulting balance when useful.
- Make eligible transactions tappable.
- Indicate unavailable details when a transaction has no `referenceId`.

### Acceptance Criteria

- The page renders a readable list of statement transactions.
- Credit and debit amounts are visually distinguishable.
- Dates are formatted consistently with the app.
- Items with a valid reference are tappable.
- Items without a valid reference do not navigate silently.

## Epic 4: Navigate to Existing Details Page

Connect statement items to the existing shared details route.

### Scope

- On transaction tap, read `StatementItemDto.referenceId`.
- Navigate to `SharedRoutes.details` with the reference as route extra.
- Keep `DetailsPage` as the destination.
- Reuse existing `DetailsViewmodel` and receipt loading behavior.
- Handle missing or blank references gracefully.

### Acceptance Criteria

- Tapping a transaction with `referenceId` opens `DetailsPage`.
- `DetailsPage` receives the selected transaction reference.
- The existing receipt loading flow runs unchanged.
- Missing reference does not crash the app.

## Epic 5: Details Route Hardening

Make the shared details route safer for statement-driven navigation.

### Scope

- Review `SharedRoutes.details`.
- Avoid unsafe casts when route extra is missing or not a `String`.
- Route to a safe fallback or show an error state when the reference is invalid.
- Keep compatibility with the existing transfer success flow.

### Acceptance Criteria

- Opening details without a valid reference does not crash.
- Existing transfer success navigation still opens details correctly.
- Statement navigation opens details correctly.

## Epic 6: Validation

Validate the statement-to-details journey.

### Scenarios

Open statement from home:

```text
Home -> Statement
```

Load statement:

```text
Statement -> AccountRepository.getStatement() -> transaction list
```

Open transaction details:

```text
Statement item tap -> SharedRoutes.details -> DetailsPage
```

Transaction without reference:

```text
Statement item without reference -> details unavailable feedback
```

Receipt load failure:

```text
DetailsPage -> receipt load fails -> retry/error state
```

## Acceptance Criteria

- The statement entry point is available from home.
- Statement transactions are loaded from the existing account statement flow.
- Statement items render enough information for the user to identify the
  transaction.
- Tapping an eligible transaction opens the existing `DetailsPage`.
- The details route receives the transaction reference from the statement item.
- Transactions without a usable reference are handled gracefully.
- The app does not introduce a duplicate details screen.
