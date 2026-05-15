# Internal Transfer Mobile Tasks

These tasks split the mobile backlog into implementation-sized steps.

Backend implementation is outside this task list. The mobile app must consume
the REST contract documented in `api/docs/07-api-rest.md`.

## Task 1/18: Align mobile transfer contract DTOs

### Goal

Add typed mobile DTOs for the current internal transfer REST contract.

### Scope

- Add a transfer request DTO for `POST /accounts/transfer`.
- Serialize `from_branch`, `from_account_number`, `to_branch`,
  `to_account_number`, `amount`, and `idempotency_key`.
- Add a transfer response DTO containing `transaction_reference`,
  public source/destination identifiers, amount, and balances.
- Ensure DTOs do not include `from_account_id` or `to_account_id`.
- Follow the existing mobile DTO naming, serialization, and folder conventions.

### Acceptance Criteria

- Transfer request JSON matches `api/docs/07-api-rest.md`.
- Transfer response parsing includes `transaction_reference`.
- DTO tests prove internal account IDs are not serialized.
- DTO tests cover required fields and representative success payloads.

### Depends On

- None.

## Task 2/18: Add receipt/detail contract DTOs

### Goal

Add typed mobile DTOs for transfer receipt/detail lookup by
`transaction_reference`.

### Scope

- Add a receipt/detail response DTO for
  `GET /accounts/transfer/{transaction_reference}/receipt`.
- Parse operation type, amount, status, transaction reference, operation date,
  source account, destination account, recipient name, and optional description.
- Keep the DTO public-reference based. Do not introduce internal transaction or
  account IDs.
- Follow the existing mobile date/time parsing conventions.

### Acceptance Criteria

- Receipt/detail response parsing matches `api/docs/07-api-rest.md`.
- Optional description is handled safely when absent.
- DTO tests cover a complete receipt payload and a payload without description.
- DTOs require only `transaction_reference` as the lookup identifier.

### Depends On

- None.

## Task 3/18: Add transfer API service methods

### Goal

Expose transfer and receipt/detail endpoints through a typed mobile API service.

### Scope

- Add or extend a typed account/transfer API service.
- Implement `POST /accounts/transfer`.
- Implement `GET /accounts/transfer/{transaction_reference}/receipt`.
- Use DTOs for request and response bodies.
- Parse existing API envelopes consistently with current account APIs.
- Map backend error envelopes into the existing `AppError` flow.

### Acceptance Criteria

- API service tests cover transfer endpoint path, method, body, success, and
  backend failure.
- API service tests cover receipt endpoint path, method, path parameter,
  success, not found, forbidden, and backend failure.
- API service code is not imported by UI or view models.
- Unknown backend failures still map to the existing generic error behavior.

### Depends On

- Task 1.
- Task 2.

## Task 4/18: Extend account repository for transfer

### Goal

Expose internal transfer through `AccountRepository`.

### Scope

- Add a repository method for internal transfer.
- Keep repository inputs expressed as branch/account-number identifiers.
- Derive or receive source branch/account number from the selected source
  account.
- Receive destination branch/account number from the transfer flow.
- Generate no UI messages in the repository.
- Propagate backend transfer failures as `AppError`.

### Acceptance Criteria

- View models can execute transfer through `AccountRepository`.
- Repository code does not accept or require internal account IDs for transfer.
- Missing selected source account returns a typed failure.
- Repository tests cover success, backend failure, and missing selected account.

### Depends On

- Task 3.

## Task 5/18: Extend account repository for receipt/detail lookup

### Goal

Expose transfer receipt/detail lookup through `AccountRepository`.

### Scope

- Add a repository method that receives `transaction_reference`.
- Call the typed receipt/detail API service method.
- Return a typed receipt/detail model to view models.
- Propagate `TRANSACTION_NOT_FOUND`, `FORBIDDEN`, and generic backend failures
  through the existing `AppError` flow.

### Acceptance Criteria

- Transaction detail view models can load receipt data through
  `AccountRepository`.
- Repository receipt lookup requires only `transaction_reference`.
- Repository tests cover success, not found, forbidden, and generic backend
  failure.

### Depends On

- Task 3.

## Task 6/18: Register transfer dependencies

### Goal

Wire the transfer API, repository dependencies, and transfer view model into the
existing dependency injection setup.

### Scope

- Register the transfer API service.
- Update `AccountRepositoryImpl` construction to receive the transfer API
  dependency.
- Register the transfer view model.
- Follow the existing service registration and UI dependency module style.

### Acceptance Criteria

- The app starts with transfer dependencies registered.
- Transfer page dependencies are resolved by the injector.
- No widget manually constructs API services or repositories.

### Depends On

- Task 4.

## Task 7/18: Register transaction detail dependencies

### Goal

Wire the transaction detail/receipt view model into the existing dependency
injection setup.

### Scope

- Register the transaction detail view model.
- Ensure it receives `AccountRepository`.
- Keep constructor and registration style consistent with existing view models.

### Acceptance Criteria

- Transaction detail page dependencies are resolved by the injector.
- No widget manually constructs repository or API service dependencies.

### Depends On

- Task 5.

## Task 8/18: Add transfer and transaction detail routes

### Goal

Add centralized navigation entries for the transfer flow and transaction detail
screen.

### Scope

- Add an internal transfer route using existing routing conventions.
- Add a transaction detail/receipt route that receives `transaction_reference`.
- Add route builders for both pages.
- Keep route paths centralized.
- Do not pass internal transaction UUIDs or account IDs through navigation.

### Acceptance Criteria

- Transfer page can be opened from app navigation.
- Transaction detail page can be opened with only `transaction_reference`.
- Routes follow the existing GoRouter and route enum/module style.
- Invalid or missing route parameters fail gracefully according to current app
  conventions.

### Depends On

- Task 6.
- Task 7.

## Task 9/18: Add home page transfer entry point

### Goal

Expose internal transfer from the authenticated home page.

### Scope

- Replace the pending transfer placeholder with navigation to the transfer page.
- Preserve existing balance tile and refresh behavior.
- Keep the action area readable on small screens.
- Do not execute financial work directly from the home page.

### Acceptance Criteria

- Authenticated users can open the internal transfer flow from the home page.
- Transfer entry point is visually consistent with the current app style.
- Returning from transfer does not break current home balance loading behavior.

### Depends On

- Task 8.

## Task 10/18: Add transfer form foundation

### Goal

Create the small shared form behavior needed by internal transfer.

### Scope

- Standardize currency amount input for this flow.
- Add positive amount validation.
- Reuse existing text field, scaffold, theme, snackbar, and button patterns.
- Add reusable helpers/widgets only where they remove real duplication.
- Ensure submit cannot run multiple times while the command is already running.

### Acceptance Criteria

- Amount input behavior matches existing app conventions.
- Validation messages are clear and transfer-specific.
- Loading and disabled states are consistent.
- The form remains usable on small mobile screens.

### Depends On

- None.

## Task 11/18: Implement transfer view model

### Goal

Add command-based transfer execution state for the internal transfer page.

### Scope

- Add transfer view model state for source account, destination branch/account
  number, amount, idempotency key, loading, success, and failure.
- Use the selected source account for `from_branch` and `from_account_number`.
- Decide locally whether `to_branch` is entered by the user or defaults from the
  selected source branch.
- Generate a mobile idempotency key for each transfer submission.
- Validate positive amount and required destination account data before calling
  the repository.
- Block same-account transfer locally when source and destination can be
  compared.
- Execute transfer through `AccountRepository`.

### Acceptance Criteria

- Invalid local input does not call the repository.
- Successful transfer exposes `transaction_reference` to the UI.
- Transfer requests always include a generated idempotency key.
- Same-account validation uses public branch/account-number identifiers.
- View model tests cover validation, success, backend failure, command state,
  and idempotency key generation.

### Depends On

- Task 4.
- Task 10.

## Task 12/18: Implement internal transfer page

### Goal

Build the user-facing internal transfer flow.

### Scope

- Show selected source account context.
- Add destination account number input.
- Add destination branch input or use the agreed default branch behavior.
- Add amount input.
- Add optional description input only if supported by the current mobile model
  and backend contract.
- Bind UI state to the transfer view model.
- Show loading, disabled, validation, success, and failure states.
- Surface backend business errors for invalid destination, same account,
  insufficient funds, inactive account, forbidden, and not found.

### Acceptance Criteria

- User can submit an internal transfer using public banking identifiers.
- UI never asks for internal account UUIDs.
- UI does not expose external transfer rails or flags.
- Success state offers a path to transaction detail/receipt.
- Widget tests cover invalid amount, missing destination, same-account
  validation, loading, success, and representative failure states.

### Depends On

- Task 8.
- Task 9.
- Task 11.

## Task 13/18: Implement transfer success behavior

### Goal

Provide a clear post-transfer state that uses the returned
`transaction_reference`.

### Scope

- Standardize the successful transfer message.
- Store or expose the returned `transaction_reference` in the success state.
- Provide a view transaction/receipt action.
- Provide a clear path back to the home page.
- Avoid building receipt content from unconfirmed form input.

### Acceptance Criteria

- Users are not left unsure whether the transfer was submitted.
- Success state can navigate to transaction detail using `transaction_reference`.
- Returning to home follows existing navigation and balance loading behavior.

### Depends On

- Task 12.

## Task 14/18: Implement transaction detail view model

### Goal

Load persisted transaction detail/receipt data by public transaction reference.

### Scope

- Add a transaction detail view model.
- Receive `transaction_reference` from navigation.
- Fetch receipt/detail data through `AccountRepository`.
- Model loading, success, invalid reference, not found, forbidden, and generic
  failure states.
- Keep the view model reusable for future statement-to-detail navigation.

### Acceptance Criteria

- View model loads receipt data using only `transaction_reference`.
- Invalid or unknown references produce clear failure states.
- Forbidden references produce an access-denied state.
- View model tests cover loading, success, invalid reference, not found,
  forbidden, and generic failure.

### Depends On

- Task 5.
- Task 8.

## Task 15/18: Implement transaction detail screen

### Goal

Build the reusable transaction detail/receipt screen.

### Scope

- Render loading, success, not-found, forbidden, and generic failure states.
- Render operation type, amount, status, transaction reference, operation
  date/time, source account, destination account, recipient name, and optional
  description.
- Load data from the transaction detail view model.
- Do not render final receipt fields from transfer form input alone.
- Keep the screen internal-transfer aware for now, while preserving route and
  view-model naming that can support statement navigation later.

### Acceptance Criteria

- Transaction detail can be opened with only `transaction_reference`.
- The screen displays persisted receipt/detail data from the backend.
- The screen never requires or displays internal account IDs as lookup values.
- Widget tests cover success, not found, forbidden, and generic failure states.

### Depends On

- Task 14.

## Task 16/18: Add receipt sharing behavior

### Goal

Allow users to share receipt/detail content after the persisted data is loaded.

### Scope

- Add share action to the transaction detail screen.
- Build share content from persisted receipt/detail data.
- Include transaction reference, amount, operation date/time, source account,
  destination account, recipient name, and status.
- Disable or hide share action until receipt/detail data is loaded.
- Follow the existing app sharing/plugin conventions if already present.

### Acceptance Criteria

- Share action is available only after successful receipt/detail load.
- Shared content includes `transaction_reference`.
- Shared content does not include internal account IDs.
- Tests cover share action availability and generated share content where
  practical.

### Depends On

- Task 15.

## Task 17/18: Strengthen end-to-end mobile regression coverage

### Goal

Ensure the full mobile transfer and receipt path is covered before
transactional password work starts.

### Scope

- Review DTO, API service, repository, view model, and widget coverage.
- Add missing tests for transfer success and failure paths.
- Add missing tests for receipt/detail lookup and rendering.
- Confirm tests run with the current mobile test command.
- Keep tests aligned with existing mobile test style.

### Acceptance Criteria

- Tests cover transfer success, validation failure, authorization/forbidden
  failure, not found, inactive account, insufficient funds, same-account
  transfer, idempotency key generation, and receipt/detail lookup.
- Tests confirm clients do not need internal account IDs.
- Receipt tests confirm lookup by `transaction_reference`.
- The mobile test suite passes with the regular Flutter test command.

### Depends On

- Task 1.
- Task 2.
- Task 3.
- Task 4.
- Task 5.
- Task 11.
- Task 12.
- Task 15.
- Task 16.

## Task 18/18: Prepare transactional password follow-up backlog

### Goal

Create a separate planning artifact for transactional password after the first
internal transfer flow is stable.

### Scope

- Summarize the implemented mobile transfer flow.
- Identify where transactional password should enter the flow.
- Keep biometric confirmation and other rails out of the first follow-up unless
  explicitly decided.
- Do not modify the current transfer implementation as part of this task.

### Acceptance Criteria

- A separate transactional password backlog can be created from the findings.
- The current internal transfer backlog remains focused on the first transfer
  and receipt flow.

### Depends On

- Task 17.
