# Task: Build internal transfer payment flow

## Goal

Build the mobile internal transfer flow as a sequence of focused pages:
recipient lookup, payment data, confirmation, operation status, and receipt.

The flow must keep navigation state explicit and serializable, avoid singleton
view models for temporary form state, and execute transfers using account IDs.

## Product Flow

1. The user searches and selects the destination account.
2. The user enters payment data: source account, amount, and optional
   description.
3. The user reviews the operation on a confirmation page.
4. The app executes the transfer and shows the operation status.
5. On success, the user can open the receipt and share it as an image.

## Architectural Decisions

- Use one view model per page.
- Do not rely on singleton view models to persist temporary transfer state.
- Pass flow data between pages through route extras.
- Keep route extras serializable through `mobile/lib/core/routing/extra_codec.dart`.
- Use draft objects with `toMap`/`fromMap` helpers for page-to-page data.
- Do not pass custom Dart objects directly through route extras.
- Serialize `Money` as integer minor units when routing between pages.
- Execute the transfer with `from_account_id` and `to_account_id`.
- Do not execute transfers with branch/account-number identifiers.
- Generate and preserve the idempotency key in the confirmation step for the
  current transfer attempt.

## Draft Objects

### TransferRecipientDraft

Carries the selected destination account from the recipient lookup page to the
payment data page.

```dart
class TransferRecipientDraft {
  final String toAccountId;
  final String document;
  final String holderName;
  final String branch;
  final String accountNumber;
}
```

Rules:

- `toAccountId` is the value later used as `to_account_id`.
- `document` must be the masked document returned by the lookup endpoint.
- The draft must not include customer ID, balance, full document, phone,
  e-mail, or unrelated internal fields.

### TransferConfirmationDraft

Carries the complete user-approved operation from the payment data page to the
confirmation page.

```dart
class TransferConfirmationDraft {
  final String fromAccountId;
  final String toAccountId;
  final String recipientDocument;
  final String recipientName;
  final String recipientBranch;
  final String recipientAccountNumber;
  final Money amount;
  final String? description;
}
```

Routing serialization rules:

- Serialize `amount` as minor units.
- Rebuild `Money` from minor units when decoding the route extra.
- Keep `description` absent/null when the user does not provide it.
- Do not include `idempotencyKey` in this draft unless the confirmation page
  needs to preserve it across rebuilds within the same transfer attempt.

## Page 1: Recipient Lookup

### Purpose

Search for eligible internal recipient accounts and let the user select the
destination account.

### Inputs

- CPF
- Branch
- Account number

### Behavior

- The user may search by CPF or by branch + account number.
- CPF lookup has priority when CPF is provided.
- CPF currently means individual document only; CNPJ and legal person account
  lookup are out of scope.
- Lookup uses `GET /accounts/internal-transfers/recipients`.
- If no accounts are returned, show an empty result state.
- If one account is returned, select it automatically.
- If more than one account is returned, show a dropdown/list and require user
  selection.
- The next page receives a `TransferRecipientDraft`.

### Acceptance Criteria

- The page can search by CPF.
- The page can search by branch + account number.
- The page handles no results.
- The page auto-selects a single result.
- The page requires selection when multiple results are returned.
- The page routes forward with only serializable recipient draft data.

## Page 2: Payment Data

### Purpose

Collect payment details after the destination account has been selected.

### Inputs

- Selected source account.
- Amount.
- Optional description.

### Behavior

- Show a summary card with the selected recipient.
- Let the user select the source account when more than one account is
  available.
- Show the selected source account balance.
- Enable confirmation only when the local balance appears sufficient.
- Treat local balance validation as a UX guard only; the backend remains the
  source of truth.
- Route to the confirmation page with a `TransferConfirmationDraft`.

### Acceptance Criteria

- The page receives and displays `TransferRecipientDraft`.
- The page lets the user select a source account.
- The page requires a positive transfer amount.
- The page prevents local confirmation when the selected account balance is
  insufficient.
- The page sends `fromAccountId`, `toAccountId`, recipient summary, amount, and
  optional description to the confirmation page.

## Page 3: Confirmation

### Purpose

Show the final operation details and execute the transfer only after explicit
user confirmation.

### Behavior

- Receive `TransferConfirmationDraft`.
- Display source account, recipient, amount, and optional description.
- Generate an idempotency key for the transfer attempt when the confirmation
  view model starts.
- Keep the same idempotency key while retrying the same confirmation attempt.
- Call the transfer use case with `from_account_id`, `to_account_id`, `amount`,
  optional `description`, and `idempotency_key`.
- Navigate to the status page with the transfer result.

### Acceptance Criteria

- The page does not execute transfer before user confirmation.
- The submitted request uses account IDs, not branch/account-number execution
  fields.
- The idempotency key is non-empty.
- Retry from the same confirmation attempt reuses the same idempotency key.
- Successful execution routes to the status page with success data.
- Failed execution routes to the status page with failure data.

## Page 4: Operation Status

### Purpose

Present the result of the transfer execution.

### Behavior

- Show success or failure state.
- On failure, allow the user to restart the transfer process.
- On success, allow the user to open the receipt.
- Use `transaction_reference` as the receipt lookup identifier.

### Acceptance Criteria

- The page clearly handles success.
- The page clearly handles failure.
- Failure restart does not reuse stale draft data.
- Success can navigate to receipt using only `transaction_reference`.

## Page 5: Receipt

### Purpose

Show transfer receipt details and allow sharing a receipt image.

### Behavior

- Load receipt details by `transaction_reference`.
- Display public confirmation data only.
- Allow sharing the rendered receipt as an image.
- After closing the receipt flow, return to the home/start page.

### Acceptance Criteria

- Receipt lookup uses only `transaction_reference`.
- The page does not require internal transaction UUIDs or account IDs.
- The page can render a receipt suitable for sharing as an image.
- Returning from the receipt clears the payment flow.

## Routing Requirements

- Keep paths centralized in the existing route enum/module style.
- Add routes for recipient lookup, payment data, confirmation, status, and
  receipt when they do not already exist.
- Route extras must use maps/lists/primitives/null accepted by
  `ExtraCodec`.
- Each draft object must provide deterministic `toMap` and `fromMap` helpers.
- Invalid or missing route extras must fail gracefully according to the current
  app routing conventions.

Suggested route shape:

- `/transfer/recipient`
- `/transfer/payment`
- `/transfer/confirmation`
- `/transfer/status`
- `/transfer/receipt`

Exact names may follow the existing `HomeRoutes` conventions.

## View Model Requirements

- Keep one view model per page.
- Do not make the transfer flow view model a singleton.
- Do not let view models construct API services or repositories directly.
- Recipient lookup page view model calls the transfer use case or repository
  method for recipient lookup.
- Payment data page view model manages local form state only.
- Confirmation page view model owns execution state and idempotency for the
  current attempt.
- Status and receipt page view models manage their own loading/result state.

## Validation Rules

- Recipient lookup accepts CPF or branch + account number.
- CPF lookup is CPF-only for now.
- The user must select a recipient account before payment data can be entered.
- The user must select a source account before confirmation.
- The user must enter a positive amount.
- Local insufficient balance prevents continuing, but backend validation is
  still authoritative.
- Same-account transfer attempts must not be submitted when the UI can detect
  them.
- Missing or empty idempotency key must prevent transfer execution.

## Out Of Scope

- Pix, TED, DOC, or interbank transfer behavior.
- CNPJ and legal person account lookup.
- Transactional password page.
- Persisting unfinished transfer drafts across app restarts.
- Backend changes beyond the already documented internal transfer contract.

## Test Timing

Use implementation-first sequencing for this backlog.

1. Build the pages, routing, drafts, and view models until the flow compiles and
   can be exercised manually.
2. Review the resulting user flow and route data contracts.
3. Adjust page boundaries, route extras, and idempotency behavior after review.
4. Add tests after the behavior is stable.

During the first implementation pass, avoid writing tests in the middle of the
feature unless a compile error or contract uncertainty blocks progress.

## Acceptance Criteria

- The flow starts at recipient lookup.
- The user can search recipients by CPF.
- The user can search recipients by branch + account number.
- Multiple recipient results require selection.
- A single recipient result can be auto-selected.
- Payment data receives a serializable recipient draft.
- Confirmation receives a serializable confirmation draft.
- Transfer execution uses `from_account_id` and `to_account_id`.
- Idempotency is stable for the same confirmation attempt.
- Status page handles success and failure.
- Receipt page loads by `transaction_reference`.
- Route extras remain compatible with `ExtraCodec`.
- No flow page depends on a singleton view model to preserve temporary state.

