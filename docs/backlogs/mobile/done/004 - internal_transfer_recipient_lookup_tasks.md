# Internal Transfer Recipient Lookup Tasks

These tasks split the internal transfer recipient lookup backlog into
implementation-sized steps.

The goal is to keep recipient discovery separate from transfer execution:
branch, account number, and document are lookup/confirmation inputs; internal
transfer execution uses `from_account_id` and `to_account_id`.

## Task 1/10: Add backend recipient lookup endpoint

### Goal

Allow authenticated users to search eligible internal recipient accounts before
executing an internal transfer.

### Scope

- Add `GET /accounts/internal-transfers/recipients`.
- Support lookup by `branch` + `account_number`.
- Support lookup by CPF through `document`.
- Keep CNPJ and legal person account lookup out of scope for this backlog.
- Normalize branch, account number, and CPF before querying.
- Return zero, one, or many eligible recipient accounts.
- Mask CPF in the response.
- Exclude closed, blocked, inactive, or non-transferable accounts.
- Keep the response limited to transfer confirmation data.

### Acceptance Criteria

- Lookup by branch + account number returns zero or one eligible account.
- Lookup by CPF returns zero, one, or many eligible accounts.
- Response includes `account_id`, `holder_name`, masked `document`, `branch`,
  `account_number`, and optional `account_type`.
- Response does not include customer ID, balance, full document, phone, e-mail,
  or unrelated internal fields.
- Invalid query combinations return typed API errors.
- Unauthorized requests are rejected.
- Backend tests cover successful lookup by account, successful lookup by
  CPF, multiple accounts for one CPF, no results, invalid parameters,
  unauthorized request, and ineligible accounts.

### Depends On

- None.

## Task 2/10: Restore ID-based internal transfer execution

### Goal

Make internal transfer execution operate on account IDs resolved before
confirmation.

### Scope

- Update `POST /accounts/internal-transfers` to receive `from_account_id` and
  `to_account_id`.
- Keep `amount`, optional `description`, and required `idempotency_key`.
- Validate that the authenticated user can operate `from_account_id`.
- Validate that `to_account_id` exists and can receive internal transfers.
- Reject transfers where source and destination are the same account.
- Preserve balance, status, limit, transaction, and idempotency validations.
- Keep receipt response public/confirmation-oriented.

### Acceptance Criteria

- Transfer execution no longer uses branch/account-number identifiers.
- Backend authorization checks ownership/operation rights for
  `from_account_id`.
- Backend rejects invalid or ineligible `to_account_id`.
- Backend rejects same-account transfer attempts.
- Existing idempotency behavior is preserved.
- Backend tests cover success, unauthorized source account, forbidden source
  account, missing destination, ineligible destination, same account,
  insufficient funds, and idempotent retry.

### Depends On

- None.

## Task 3/10: Update internal transfer REST contract

### Goal

Document the implemented internal transfer flow as recipient lookup followed by
ID-based transfer execution.

### Scope

- Document `GET /accounts/internal-transfers/recipients` after the endpoint is
  implemented.
- Document `POST /accounts/internal-transfers` after the execution contract is
  restored to account IDs.
- Document receipt/detail endpoint naming for internal transfers if it changes.
- Define lookup query modes: branch + account number, or CPF document.
- Define ID-based execution with `from_account_id` and `to_account_id`.
- Define response fields, masking rules, errors, and authorization behavior.
- Remove branch/account-number execution from the current transfer contract
  documentation.

### Acceptance Criteria

- API docs describe recipient lookup before transfer execution.
- Transfer execution docs use `from_account_id` and `to_account_id`.
- Lookup docs clearly state that the endpoint is internal-transfer-only.
- Error cases include invalid query, no results, forbidden, unauthorized, and
  ineligible recipient account.
- Docs state that full document, customer ID, balances, phone, and e-mail are
  not returned by lookup.

### Depends On

- Task 1.
- Task 2.

## Task 4/10: Add mobile recipient lookup DTOs

### Goal

Add typed mobile DTOs for internal transfer recipient lookup.

### Scope

- Add query/input DTOs for lookup by branch + account number or CPF document.
- Treat `document` as CPF-only for now; CNPJ/legal person lookup is future
  scope.
- Add response DTO for the lookup envelope.
- Add recipient account DTO with `account_id`, `holder_name`, masked
  `document`, `branch`, `account_number`, and optional `account_type`.
- Keep DTOs aligned with `api/docs/07-api-rest.md`.
- Do not expose customer ID, balances, phone, e-mail, or full document.

### Acceptance Criteria

- Query serialization supports branch + account number.
- Query serialization supports CPF through `document`.
- DTO parsing handles zero, one, and multiple accounts.
- DTO tests prove internal/prohibited fields are not exposed.
- DTO tests cover representative account and CPF lookup payloads.

### Depends On

- Task 3.

## Task 5/10: Update mobile internal transfer DTOs

### Goal

Align mobile transfer execution DTOs with ID-based internal transfer execution.

### Scope

- Update transfer request DTO to serialize `from_account_id` and
  `to_account_id`.
- Keep `amount`, optional `description`, and required `idempotency_key`.
- Remove branch/account-number execution fields from transfer request DTO.
- Keep response and receipt DTOs public/confirmation-oriented.

### Acceptance Criteria

- Transfer request JSON matches the ID-based internal transfer contract.
- Request DTO tests prove branch/account-number execution fields are not
  serialized.
- Required fields include `from_account_id`, `to_account_id`, `amount`, and
  `idempotency_key`.
- Optional `description` is omitted when absent.

### Depends On

- Task 3.

## Task 6/10: Add mobile API service methods

### Goal

Expose recipient lookup and ID-based internal transfer execution through typed
mobile API services.

### Scope

- Add or update typed API service for internal transfers.
- Implement `GET /accounts/internal-transfers/recipients`.
- Implement `POST /accounts/internal-transfers`.
- Use DTOs for lookup query, lookup response, transfer request, and transfer
  response.
- Parse API envelopes consistently with existing API services.
- Map backend error envelopes into the existing `AppError` flow.

### Acceptance Criteria

- API service tests cover recipient lookup path, query parameters, success,
  no results, invalid query, forbidden/unauthorized, and generic backend
  failure.
- API service tests cover transfer path, method, body, success, backend
  failure, and generic backend failure.
- API service code is not imported by UI widgets.

### Depends On

- Task 4.
- Task 5.

## Task 7/10: Update transaction repository for recipient lookup

### Goal

Expose internal transfer recipient lookup and ID-based transfer execution
through the transaction repository.

### Scope

- Add repository method for recipient lookup.
- Keep lookup inputs expressed as branch/account-number or CPF document.
- Keep transfer execution input ID-based.
- Propagate lookup and transfer backend failures as `AppError`.
- Keep no UI messages or navigation in the repository.
- Cache only data that is intentionally useful to the transfer flow.

### Acceptance Criteria

- View models/usecases can search internal recipients through
  `TransactionRepository`.
- Repository transfer execution uses account IDs, not branch/account-number
  identifiers.
- Repository tests cover lookup success, multiple results, no results, invalid
  query, forbidden/unauthorized, and backend failure.
- Repository tests cover ID-based transfer success and backend failure.

### Depends On

- Task 6.

## Task 8/10: Update transfer usecase workflow

### Goal

Coordinate selected source account, recipient lookup, recipient selection, and
ID-based internal transfer execution.

### Scope

- Add usecase input objects for recipient lookup.
- Add usecase input object for confirmed internal transfer execution.
- Use selected source account ID as `from_account_id`.
- Use selected recipient account ID as `to_account_id`.
- Preserve idempotency key behavior supplied by the view model.
- Return typed failures for missing selected recipient or missing idempotency
  key.
- Keep API calls inside repositories, not the usecase.

### Acceptance Criteria

- Usecase can search recipients by branch + account number.
- Usecase can search recipients by CPF document.
- Usecase can execute transfer only after a recipient account is selected.
- Usecase does not use branch/account-number as execution identifiers.
- Usecase tests cover lookup, multiple recipient selection, missing recipient,
  missing idempotency key, success, and backend failure propagation.

### Depends On

- Task 7.

## Task 9/10: Update transfer view model

### Goal

Support recipient lookup, recipient selection, and ID-based transfer submission
from the transfer screen state.

### Scope

- Add command/state for recipient lookup.
- Store lookup results for user confirmation.
- Store the selected recipient account.
- Keep idempotency key stable for a transfer attempt.
- Generate a new idempotency key after successful transfer or explicit reset.
- Prevent transfer submission while lookup/transfer commands are running.
- Do not let the view model construct API services or repositories.

### Acceptance Criteria

- View model can search by branch + account number.
- View model can search by CPF document.
- View model requires recipient selection when lookup returns multiple
  accounts.
- View model submits transfer with `from_account_id`, `to_account_id`, amount,
  optional description, and idempotency key.
- View model tests cover successful lookup, no results, multiple results,
  recipient selection, transfer success, transfer failure, idempotency retry,
  and idempotency reset after success.

### Depends On

- Task 8.

## Task 10/10: Update transfer UI flow

### Goal

Let users find, confirm, and transfer to an internal recipient account.

### Scope

- Add search mode for branch + account number.
- Add search mode for CPF document.
- Show lookup loading, empty, error, one-result, and multiple-result states.
- Require user selection when multiple recipient accounts are returned.
- Show confirmation data: holder name, masked document, branch, account
  number, and account type when available.
- Submit transfer using the confirmed recipient.
- Keep route/navigation behavior centralized and free of internal transaction
  IDs.

### Acceptance Criteria

- User can search recipient by branch + account number.
- User can search recipient by CPF.
- User can select among multiple returned accounts.
- User cannot submit transfer before confirming/selecting a recipient.
- UI does not show full document, customer ID, balance, phone, e-mail, or
  unrelated internal fields.
- Successful transfer navigates or presents the next state using
  `transaction_reference`, not internal transaction UUIDs.

### Depends On

- Task 9.
