# Internal Transfer Payment Flow Tasks

These tasks split the internal transfer payment flow into implementation-sized
steps.

The flow is page-based. Temporary transfer state must be passed through
serializable route extras instead of being preserved in a singleton view model.
Internal transfer execution must use `from_account_id` and `to_account_id`.

Tests should be added after the flow behavior and route contracts are stable.

## Task 1/11: Add serializable transfer flow drafts

### Goal

Define the data objects passed between transfer flow pages.

### Scope

- Add `TransferRecipientDraft`.
- Add `TransferConfirmationDraft`.
- Add `toMap` and `fromMap` helpers compatible with `ExtraCodec`.
- Serialize `Money` as integer minor units.
- Rebuild `Money` from minor units when decoding confirmation draft data.
- Keep drafts free of API-only envelope objects and repository state.

### Acceptance Criteria

- Recipient draft carries `toAccountId`, masked `document`, `holderName`,
  `branch`, and `accountNumber`.
- Confirmation draft carries `fromAccountId`, `toAccountId`,
  `recipientDocument`, `recipientName`, `recipientBranch`,
  `recipientAccountNumber`, `amount`, and optional `description`.
- Draft maps contain only primitives, lists, maps, and null values accepted by
  `mobile/lib/core/routing/extra_codec.dart`.
- Drafts do not include customer ID, balance, full document, phone, e-mail, or
  unrelated internal fields.

### Depends On

- None.

## Task 2/11: Add transfer flow routes

### Goal

Centralize navigation entries for the page-based internal transfer flow.

### Scope

- Add route entries for recipient lookup, payment data, confirmation, status,
  and receipt.
- Keep paths centralized in the existing route enum/module style.
- Decode route extras into transfer draft objects in route builders.
- Fail gracefully when required route extras are missing or invalid.
- Do not pass internal transaction UUIDs or account IDs except the account IDs
  required by transfer execution drafts.

### Acceptance Criteria

- Recipient lookup page can be opened from app navigation.
- Payment data page receives only a serializable recipient draft.
- Confirmation page receives only a serializable confirmation draft.
- Status page receives a serializable transfer result/status payload.
- Receipt page can be opened with only `transaction_reference`.
- Invalid route extras follow the current app error/fallback conventions.

### Depends On

- Task 1.

## Task 3/11: Build recipient lookup page state and view model

### Goal

Support recipient search and selection before payment data is entered.

### Scope

- Add recipient lookup page view model.
- Add CPF search state.
- Add branch + account number search state.
- Use the existing recipient lookup use case/repository API.
- Store lookup loading, success, empty, and failure states.
- Store selected recipient account.
- Auto-select the recipient when lookup returns exactly one account.
- Require explicit selection when lookup returns multiple accounts.
- Do not construct API services or repositories in the view model.

### Acceptance Criteria

- View model can search by CPF.
- View model can search by branch + account number.
- CPF lookup is CPF-only for now.
- Empty lookup results are represented without failure.
- Multiple lookup results require user selection.
- A selected recipient can be converted to `TransferRecipientDraft`.

### Depends On

- Task 1.

## Task 4/11: Build recipient lookup page UI

### Goal

Let the user search and select the destination account.

### Scope

- Add page with CPF input.
- Add page with branch and account number inputs.
- Show lookup loading state.
- Show no-result state.
- Show one-result selected state.
- Show dropdown/list selection when multiple accounts are returned.
- Enable continuing only when a recipient is selected.
- Route to the payment data page with `TransferRecipientDraft`.

### Acceptance Criteria

- User can search recipient by CPF.
- User can search recipient by branch + account number.
- User can select one account from multiple results.
- User cannot continue without a selected recipient.
- UI shows masked document only.
- UI does not show customer ID, full document, balance, phone, e-mail, or
  unrelated internal fields.

### Depends On

- Task 2.
- Task 3.

## Task 5/11: Build payment data page state and view model

### Goal

Collect source account, amount, and optional description after recipient
selection.

### Scope

- Add payment data page view model.
- Receive `TransferRecipientDraft`.
- Load or receive available source accounts using existing account state.
- Store selected source account.
- Store amount and optional description.
- Expose selected source account balance for display.
- Calculate whether the local balance appears sufficient.
- Build `TransferConfirmationDraft`.
- Do not execute the transfer from this page.

### Acceptance Criteria

- View model requires a recipient draft.
- View model can select a source account.
- View model validates positive amount.
- View model reports insufficient local balance as a UI guard.
- View model can build confirmation draft with `fromAccountId` and
  `toAccountId`.
- View model does not call transfer execution.

### Depends On

- Task 1.

## Task 6/11: Build payment data page UI

### Goal

Let the user review the recipient and enter payment details.

### Scope

- Display recipient summary card.
- Display source account selector when needed.
- Display selected source account balance.
- Add amount input.
- Add optional description input.
- Disable continue when required data is missing.
- Disable continue when local balance appears insufficient.
- Route to confirmation with `TransferConfirmationDraft`.

### Acceptance Criteria

- Page displays recipient name, masked document, branch, and account number.
- User can choose the paying account.
- User can enter amount.
- User can enter optional description.
- Continue button is enabled only with source account, valid amount, and
  sufficient local balance.
- The page routes forward with serializable confirmation draft data.

### Depends On

- Task 2.
- Task 5.

## Task 7/11: Build confirmation page state and view model

### Goal

Execute the internal transfer only after explicit user approval.

### Scope

- Add confirmation page view model.
- Receive `TransferConfirmationDraft`.
- Generate a non-empty idempotency key for the current transfer attempt.
- Keep the same idempotency key while retrying the same attempt.
- Build the final transfer execution input/request.
- Call the transfer use case.
- Store loading, success, and failure states.
- Do not construct API services or repositories in the view model.

### Acceptance Criteria

- View model does not execute transfer before confirmation.
- Transfer request uses `from_account_id` and `to_account_id`.
- Transfer request includes `amount`, optional `description`, and
  `idempotency_key`.
- Retry from the same confirmation page instance reuses the idempotency key.
- Success exposes `transaction_reference` for the status/receipt flow.
- Failure exposes the propagated `AppError`.

### Depends On

- Task 1.

## Task 8/11: Build confirmation page UI

### Goal

Show the operation summary and let the user approve the transfer.

### Scope

- Display source account summary.
- Display recipient summary.
- Display amount.
- Display optional description when present.
- Add explicit confirm action.
- Prevent duplicate submission while command is running.
- Navigate to status page with success or failure result.

### Acceptance Criteria

- Page clearly shows the operation before execution.
- User can go back before executing.
- Confirm action executes the transfer once per tap.
- Loading state prevents duplicate submissions.
- Success and failure both navigate to the status page with serializable data.

### Depends On

- Task 2.
- Task 7.

## Task 9/11: Build operation status page

### Goal

Present the result of the transfer execution.

### Scope

- Add status page route payload.
- Show success state with transfer reference.
- Show failure state with backend/application error information.
- Allow restarting the transfer process after failure.
- Allow opening the receipt after success.
- Ensure restart does not reuse stale transfer drafts.

### Acceptance Criteria

- Status page handles success.
- Status page handles failure.
- Failure can restart the flow at recipient lookup.
- Success can navigate to receipt with only `transaction_reference`.
- Status route data is serializable by `ExtraCodec`.

### Depends On

- Task 2.
- Task 8.

## Task 10/11: Build transfer receipt page

### Goal

Show transfer receipt details and allow sharing a receipt image.

### Scope

- Add receipt page or adapt the existing transaction detail page for this flow.
- Load receipt data by `transaction_reference`.
- Display public confirmation data only.
- Add receipt image capture/share behavior.
- Return to home/start page when the receipt flow is closed.

### Acceptance Criteria

- Receipt lookup requires only `transaction_reference`.
- Receipt page does not require internal transaction UUIDs.
- Receipt page does not require account IDs.
- Receipt can be rendered for sharing as an image.
- Closing the receipt clears the transfer flow.

### Depends On

- Task 2.
- Task 9.

## Task 11/11: Add transfer flow tests

### Goal

Cover the stabilized page-based internal transfer flow.

### Scope

- Add draft serialization tests.
- Add route extra decoding tests where current routing test patterns allow it.
- Add recipient lookup view model tests.
- Add payment data view model tests.
- Add confirmation view model tests.
- Add status/receipt view model tests where behavior is non-trivial.
- Add widget tests for the critical page states if the existing project style
  supports them.

### Acceptance Criteria

- Draft tests prove route maps are `ExtraCodec` compatible.
- Recipient lookup tests cover CPF search, branch/account search, empty
  results, one result, multiple results, selection, and backend failure.
- Payment data tests cover source selection, amount validation, insufficient
  local balance, and confirmation draft creation.
- Confirmation tests cover idempotency creation, idempotency retry stability,
  success, duplicate-submit prevention, and failure propagation.
- Status tests cover success, failure, restart, and receipt navigation data.
- Receipt tests cover lookup by `transaction_reference` and backend failure.

### Depends On

- Task 1.
- Task 2.
- Task 3.
- Task 4.
- Task 5.
- Task 6.
- Task 7.
- Task 8.
- Task 9.
- Task 10.

