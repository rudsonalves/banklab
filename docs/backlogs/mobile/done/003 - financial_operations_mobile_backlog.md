# Internal Transfer Mobile Backlog

## 1. Objective

Extend the Flutter mobile application so authenticated users can execute
internal transfers from the app and view the resulting transaction details.

This backlog defines the application changes required to turn internal transfer
into a usable mobile flow. It intentionally comes before transactional password
implementation so the team can validate API contracts, mobile state management,
operation results, operation errors, and navigation behavior first.

This backlog is strictly about the mobile application implementation. Backend
changes remain outside its scope. The mobile work should consume the REST
contract documented in `api/docs/07-api-rest.md`.

## 2. Product Scope

This backlog covers the first mobile version of internal transfer.

Included:

- home page entry point for internal transfer;
- dedicated internal transfer flow;
- request DTO and API service class for the updated transfer endpoint;
- repository method that exposes transfer to view models;
- view model with command-based execution state;
- local input validation for required fields and positive amount;
- mobile-generated idempotency key for every transfer request;
- backend error handling through the existing `AppError` flow;
- loading, disabled, success, and failure states;
- transaction detail/receipt screen backed by `transaction_reference`;
- receipt view and share actions after successful transfer;
- dependency injection and routing updates;
- unit and focused widget tests for the new behavior.

Product decisions:

- Deposit should not be exposed as a mobile user operation in this backlog. It
  may make sense later for a cashier or back-office flow, but not for the mobile
  app.
- Withdrawal should not be exposed as a mobile user operation in this backlog.
- Internal transfers should be entered with banking identifiers, not internal
  account IDs. The bank is currently implicit in the backend contract. The
  selected source account provides `from_branch` and `from_account_number`.
  The transfer request must also send `to_branch` and `to_account_number`.
- For the first mobile version, the destination branch may default to the
  selected source account branch if the UI does not expose a destination branch
  field. Even when hidden, the mobile DTO must send `to_branch` explicitly.
- Financial operation APIs consumed by mobile should use public banking
  identifiers instead of internal account IDs. The backend should resolve bank,
  branch, and account number to internal IDs inside the application layer.
- Internal transfer should remain separate from future external transfer rails.
  TED, DOC, Pix, or other external operations should be planned later as their
  own transfer types or endpoints, not as flags ignored by the internal transfer
  flow.
- The mobile app should generate an idempotency key for every internal transfer
  request.
- After a successful transfer, the app should present receipt actions instead of
  immediately returning to the home page. Users should be able to view the
  receipt, share it, and then return to the home page.
- Successful transfer responses return `transaction_reference` plus public
  account identifiers and balances. The mobile app should use
  `transaction_reference` as the public lookup key for receipt/detail data.
- The first receipt version should include a basic trusted set of fields:
  operation type, amount, status, transaction reference, operation date/time,
  source account, transfer destination account number, recipient name when
  applicable, and optional description when provided by the user.

Not included in this backlog:

- deposit mobile flow;
- withdrawal mobile flow;
- transactional password;
- biometric confirmation;
- push notifications;
- external transfers;
- Pix, TED, DOC, boleto, card operations, or other payment rails;
- balance visibility behavior when the home page cannot load balance data;
- statement refresh or statement validation after internal transfer;
- full statement screen implementation;
- full statement-to-detail navigation, except keeping the transaction detail
  route reusable for that future flow;
- backend implementation changes.

## 3. Technical Direction

The implementation should preserve the current mobile architecture:

```text
UI -> ViewModel -> Repository -> API/Service -> RestClient
```

The first version should extend the existing account feature instead of
introducing a new application architecture. Operation screens should depend on
view models. View models should call `AccountRepository`. Repository
implementations should call typed account API services. UI code must not call
API services, `RestClient`, or Dio directly.

Use cases under `mobile/lib/domain/usecases` should not be introduced for this
backlog unless a later task explicitly asks for that refactor.

Current REST contract consumed by mobile:

- Transfer endpoint: `POST /accounts/transfer`
- Transfer request:
  - `from_branch`
  - `from_account_number`
  - `to_branch`
  - `to_account_number`
  - `amount`
  - `idempotency_key`
- Transfer response includes:
  - `transaction_reference`
  - `from_branch`
  - `from_account_number`
  - `to_branch`
  - `to_account_number`
  - `amount`
  - `from_balance`
  - `to_balance`
- Receipt/detail endpoint:
  `GET /accounts/transfer/{transaction_reference}/receipt`
- Receipt/detail response includes:
  - `operation_type`
  - `amount`
  - `status`
  - `transaction_reference`
  - `operation_date`
  - `source_branch`
  - `source_account_number`
  - `destination_branch`
  - `destination_account_number`
  - `recipient_name`
  - optional `description`
- Transfer errors that should be mapped to user-facing states include
  `INVALID_REQUEST`, `INVALID_DATA`, `INVALID_AMOUNT`,
  `SAME_ACCOUNT_TRANSFER`, `FORBIDDEN`, `ACCOUNT_NOT_FOUND`,
  `INSUFFICIENT_FUNDS`, and `ACCOUNT_INACTIVE`.
- Receipt/detail errors that should be mapped to user-facing states include
  `INVALID_DATA`, `FORBIDDEN`, and `TRANSACTION_NOT_FOUND`.

## 4. Epic 1: Transfer API Layer

### Goal

Add typed mobile API service and DTOs for the updated internal transfer
contract.

### Backlog Items

- Add an internal transfer request DTO for the updated transfer endpoint using
  public banking identifiers.
- Transfer request DTO must serialize `from_branch`, `from_account_number`,
  `to_branch`, `to_account_number`, `amount`, and `idempotency_key`.
- Transfer requests should identify the source and destination accounts with
  public branch/account-number pairs. The mobile app must not send
  `from_account_id` or `to_account_id`.
- Add a typed transfer operation response DTO containing at least
  `transaction_reference`.
- Add a typed receipt/detail response DTO that fetches persisted receipt data by
  `transaction_reference`.
- Add a typed `TransferApi` under the existing account API service structure.
- Add a typed method for `GET /accounts/transfer/{transaction_reference}/receipt`.
- Ensure request bodies are serialized through DTOs, not built inline by view
  models or widgets.
- Parse success and failure responses consistently with the existing account
  APIs.
- Map backend failures into the existing `AppError` model.

### Acceptance Criteria

- Internal transfer has a typed API service class.
- Internal transfer request has a typed DTO.
- Internal transfer response has a typed DTO containing `transaction_reference`.
- Receipt/detail response has a typed DTO containing persisted operation data.
- Endpoint path, HTTP method, and request body are covered by tests.
- Receipt/detail endpoint path, HTTP method, path parameter, and response parsing
  are covered by tests.
- API service tests cover success, validation failure, and generic backend
  failure.
- No UI or view model imports API service classes directly.

## 5. Epic 2: Account Repository Extension

### Goal

Expose internal transfer through the account repository boundary used by the UI
layer.

### Backlog Items

- Add `transfer` to `AccountRepository`.
- Keep repository transfer method expressed in public banking identifiers, not
  internal account IDs.
- Repository transfer method should receive or derive source branch/account
  number from the selected account and receive destination branch/account number
  from the transfer flow.
- Implement transfer in `AccountRepositoryImpl`.
- Inject the new transfer API service into `AccountRepositoryImpl`.
- Return a typed transfer operation result containing at least `transaction_reference`.
- Add repository support for fetching transfer receipt/detail by
  `transaction_reference`.
- Require a selected account for transfer source.
- Return a typed failure when no account is selected.
- Keep all UI messaging out of the repository.

### Acceptance Criteria

- View models can execute transfer through `AccountRepository`.
- Transfer view models pass branch/account-number identifiers to the repository
  instead of passing internal account IDs.
- Transaction detail view models can load receipt/detail data through
  `AccountRepository` using only `transaction_reference`.
- Missing selected account is handled as a failure, not an uncaught exception.
- Backend transfer failures are propagated as `AppError`.
- Repository tests cover transfer success, receipt/detail success, backend
  failure, receipt/detail not found, and missing selected account.

## 6. Epic 3: Dependency Injection And Routing

### Goal

Wire the transfer API, repository dependency, view model, and page into the
existing application setup.

### Backlog Items

- Register the transfer API service in the existing service registration module.
- Update repository construction to receive the transfer API.
- Register the transfer view model in the UI dependency module.
- Register the transaction detail/receipt view model in the UI dependency
  module.
- Add transfer route definition using the existing routing conventions.
- Add transaction detail/receipt route definition that receives
  `transaction_reference`.
- Add route builder for the internal transfer page.
- Add route builder for the transaction detail/receipt page.
- Ensure the transfer page receives dependencies from the injector.
- Ensure the transaction detail page receives dependencies from the injector.
- Keep route paths centralized in the routing layer.

### Acceptance Criteria

- The app starts with all new transfer dependencies registered.
- Transfer page can be opened without manual object construction in widgets.
- Transaction detail page can be opened with a `transaction_reference` without
  manual object construction in widgets.
- Route follows the existing route enum and route module style.
- Navigation remains consistent with the current GoRouter setup.

## 7. Epic 4: Home Page Entry Point

### Goal

Expose internal transfer from the authenticated home page.

### Backlog Items

- Replace the current pending transfer placeholder with real transfer
  navigation.
- Keep the action area readable on small screens.
- Preserve the existing balance tile and refresh action behavior.
- Decide whether the home action layout should remain a simple grid/row or move
  to a more scalable quick-actions section.

### Acceptance Criteria

- Authenticated users can open the internal transfer flow from the home page.
- Transfer entry point is visually consistent with the existing app style.
- Tapping transfer does not execute financial work directly from the home page.
- Returning to the home page after a successful transfer follows the existing
  home balance loading behavior.

## 8. Epic 5: Shared Operation Form Foundation

### Goal

Create a small, reusable foundation for amount-based operation forms without
overbuilding a generic transaction framework.

### Backlog Items

- Standardize currency amount input behavior for internal transfer.
- Define positive amount validation.
- Reuse existing text field, scaffold, theme, and snackbar patterns.
- Create shared helpers or widgets only where duplication becomes meaningful.
- Standardize submit button loading and disabled states.
- Ensure the form remains usable on small mobile screens.
- Ensure command execution cannot be triggered multiple times while already
  running.

### Acceptance Criteria

- Amount input behavior is consistent with existing app conventions.
- Validation messages are clear and transfer-specific.
- Loading and disabled states are consistent.
- Shared code reduces real duplication without hiding transfer-specific rules.

## 9. Epic 6: Internal Transfer Flow

### Goal

Implement internal transfer between BankLab accounts from the mobile app using
the destination account number entered by the user.

### Backlog Items

- Build an internal transfer page with selected source account context.
- Add destination account number input.
- Use the selected source account context to populate `from_branch` and
  `from_account_number`.
- Populate `to_branch` explicitly. For the first version, either default it to
  the selected source branch or expose a destination branch input if that fits
  the current UI.
- Keep this flow dedicated to internal transfers only. Do not add TED, DOC, Pix,
  or external-transfer flags to this mobile flow.
- Add amount input.
- Add optional description input if supported by the backend contract.
- Validate that amount is present and greater than zero.
- Validate that destination account number is present.
- Prevent same-account transfer locally when the app has enough account data to
  detect it, comparing the selected source account number with the destination
  account number within the same bank and branch context.
- Submit transfer data using `from_branch`, `from_account_number`, `to_branch`,
  and `to_account_number`. Do not resolve source or destination accounts to
  internal IDs in the mobile app.
- Generate and send an idempotency key for every internal transfer request.
- Execute transfer through a transfer view model command.
- Show loading state while the command is running.
- Show success feedback after the operation completes.
- Surface insufficient-funds, invalid-destination, same-account, and
  account-status errors from the backend.
- Present receipt actions after success, including view receipt and share
  receipt.

### Acceptance Criteria

- User can transfer from the selected account to another internal account.
- Invalid local input does not call the repository.
- Same-account transfer is blocked locally when detectable.
- The UI does not ask the user to type an internal account UUID.
- The user-facing destination input is the account number.
- The request payload contains branch/account-number identifiers and never
  contains internal account IDs.
- The mobile API payload does not expose internal destination account IDs.
- The flow does not expose external transfer types or flags.
- Backend business errors are specific enough for validation.
- Transfer requests always include a mobile-generated idempotency key.
- Transfer idempotency behavior matches the confirmed backend contract.
- Transfer tests cover invalid amount, missing destination, same-account
  validation, success, insufficient funds, invalid destination, backend failure,
  and idempotency key generation.

## 10. Epic 7: Transfer Result And Receipt Behavior

### Goal

Define a consistent post-transfer experience for the first mobile version.

### Backlog Items

- Show a transfer success state after internal transfer using the returned
  `transaction_reference`.
- Fetch receipt details by `transaction_reference` after a successful transfer.
- Add a dedicated transaction detail/receipt screen.
- Provide a view transaction/receipt action after successful transfer.
- Provide a share receipt action from the transaction detail screen.
- Render the basic receipt fields for internal transfer from persisted backend
  data.
- Use a dedicated receipt/detail endpoint as the source of receipt data.
- Standardize transfer success message.
- Standardize error message presentation using existing snackbar patterns.
- Keep a clear path back to the home page after the user reviews or shares the
  receipt.

### Acceptance Criteria

- Internal transfer has a clear success state.
- Successful transfer offers a view transaction/receipt action.
- Transaction detail screen loads by `transaction_reference`.
- Transaction detail screen supports sharing the receipt once data is loaded.
- Receipt content is generated from trustworthy operation data, not from
  unconfirmed form input alone.
- Receipt content is loaded by `transaction_reference` from the receipt/detail
  endpoint.
- Receipt shows operation type, amount, status, transaction reference,
  operation date/time, source account, recipient name when applicable,
  destination account number, and optional description when available.
- Internal transfer has a clear failure state.
- Users are not left unsure whether a transfer was submitted.
- Users can return to the home page after reviewing or sharing the receipt.

## 11. Epic 8: Transaction Detail Screen

### Goal

Provide a reusable mobile screen for viewing a persisted financial transaction
by public transaction reference.

### Backlog Items

- Build a transaction detail page that receives `transaction_reference` from
  navigation.
- Load receipt/detail data from
  `GET /accounts/transfer/{transaction_reference}/receipt`.
- Show loading, success, not-found, forbidden, and generic failure states.
- Render operation type, amount, status, transaction reference, operation
  date/time, source account, destination account, recipient name, and optional
  description.
- Keep the screen internal-transfer aware for now, but structure route and
  view-model naming so statement navigation can reuse it later.
- Allow navigation from the transfer success state to this screen.
- Keep share behavior on this screen instead of relying only on the transfer
  success state.
- Do not require an internal transaction UUID anywhere in navigation, state, or
  DTOs.

### Acceptance Criteria

- Transaction detail can be opened with only `transaction_reference`.
- Transaction detail fetches persisted data from the backend before rendering
  receipt fields.
- Invalid or unknown references show a clear failure state.
- Forbidden references show an access-denied state.
- The screen does not render receipt content from form input alone.
- The route can later be reused from statement items without changing the public
  receipt lookup contract.

## 12. Epic 9: Test Coverage

### Goal

Add enough automated coverage to safely evolve internal transfer before adding
transactional password.

### Backlog Items

- Add DTO serialization tests for transfer requests.
- Add DTO parsing tests for transfer responses and receipt/detail responses.
- Add API service tests for endpoint, method, request body, success, and
  failure.
- Add receipt/detail API service tests for endpoint, path parameter, success,
  not found, forbidden, and generic backend failure.
- Add repository tests for selected account behavior, transfer failure, and
  receipt/detail lookup.
- Add view model tests for validation, command state, success, and failure.
- Add transaction detail view model tests for loading, success, not found,
  forbidden, and generic failure.
- Add focused widget tests for transfer form behavior where practical.
- Add focused widget tests for transaction detail rendering where practical.
- Confirm the transfer tests run with the current mobile test command.

### Acceptance Criteria

- Core transfer behavior is covered before transactional password work starts.
- Success and failure paths are both tested.
- Receipt/detail lookup and rendering are covered.
- Tests follow the existing mobile test style.
- The transfer test suite can run as part of the regular Flutter tests.

## 13. Suggested Delivery Order

1. Read `api/docs/07-api-rest.md` and align the mobile DTOs with the documented
   transfer and receipt/detail contracts.
2. Add transfer DTOs and API service.
3. Add receipt/detail DTOs and API service method.
4. Extend `AccountRepository` and `AccountRepositoryImpl`.
5. Register new services and view models in dependency injection.
6. Add transfer and transaction detail routes.
7. Add home page entry point.
8. Implement the transfer form foundation.
9. Implement internal transfer flow.
10. Implement transaction detail/receipt screen.
11. Standardize transfer result behavior and share actions.
12. Strengthen tests around the full transfer and receipt path.
13. Create a separate transactional password backlog after this flow is stable.

## 14. Open Decisions

- Decide whether the first transfer UI exposes destination branch input or
  defaults `to_branch` from the selected source account branch.
- Decide whether the transaction detail screen is named "Receipt" in the UI,
  "Transaction detail", or uses both labels depending on navigation context.

## 15. Definition Of Ready For Task Breakdown

This backlog is ready to be split into implementation tasks when:

- the mobile team has reviewed the transfer and receipt/detail contracts in
  `api/docs/07-api-rest.md`;
- destination branch UI behavior is decided;
- the minimum test depth for the first implementation pass is agreed.
