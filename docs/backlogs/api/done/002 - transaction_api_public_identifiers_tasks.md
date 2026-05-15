# Transfer API Public Identifiers Tasks

Parent backlog:

- `002 - transaction_api_public_identifiers_backlog.md`

Suggested project fields for all tasks:

- Status: Backlog
- Priority: High
- Area: Account
- Type: Improvement

## Task 1: Define the public internal transfer API contract

### Goal

Define the new public request and response contract for
`POST /accounts/transfer` before implementation starts.

### Scope

- Define the request fields using public banking identifiers.
- Remove `from_account_id` and `to_account_id` from the public request.
- Keep `amount` and `idempotency_key`.
- Define how the selected/source account is identified from authenticated
  context.
- Define the successful response shape with at least `transaction_reference`.
- Define validation and error behavior for missing or invalid destination account
  number.

### Acceptance Criteria

- The contract no longer requires internal account IDs in the public request.
- The request clearly identifies the destination by account number.
- The response includes `transaction_reference`.
- Error cases are listed for invalid amount, invalid destination, same-account
  transfer, insufficient funds, inactive account, unauthorized, forbidden, and
  not found.
- The contract remains explicitly internal-transfer only.

### Depends On

- None.

## Task 2: Add backend account lookup by public banking identifiers

### Goal

Add backend support to resolve accounts from public banking identifiers before
executing the transfer.

### Scope

- Add or expose lookup behavior for account number in the known bank and branch
  context.
- Ensure destination account lookup returns clear not-found or invalid-data
  errors.
- Keep internal UUIDs inside repository/application boundaries.
- Avoid exposing internal account IDs in public request parsing.

### Acceptance Criteria

- Destination account can be resolved by account number within the selected
  source account bank/branch context.
- Missing, malformed, or unknown destination account numbers return clear API
  errors.
- Lookup behavior is covered by unit or integration tests.
- No mobile/client-facing contract requires `account_id`.

### Depends On

- Task 1.

## Task 3: Resolve the transfer source account from authenticated context

### Goal

Resolve the transfer source account server-side instead of accepting
`from_account_id` from the public request.

### Scope

- Determine how the selected/source account is obtained from authenticated user
  context and request data.
- Enforce ownership and authorization after source account resolution.
- Return clear errors when no usable source account can be resolved.
- Preserve customer/admin access rules already used by account operations.

### Acceptance Criteria

- Public transfer requests do not include `from_account_id`.
- Source account resolution is performed server-side.
- Unauthorized source account access returns `FORBIDDEN`.
- Missing or invalid source context returns a typed API error.
- Existing authorization rules remain enforced.

### Depends On

- Task 1.

## Task 4: Refactor the transfer handler request parsing

### Goal

Update the `POST /accounts/transfer` handler to parse and validate the new
public transfer request.

### Scope

- Replace UUID-based request parsing with the new public request body.
- Validate amount.
- Validate destination account number.
- Preserve strict JSON parsing behavior where applicable.
- Map validation failures to existing API envelope and error codes.
- Keep the route as `POST /accounts/transfer`.

### Acceptance Criteria

- Handler accepts the new request shape.
- Handler rejects `from_account_id` and `to_account_id` if unknown fields are
  rejected by the API style.
- Invalid JSON returns `INVALID_REQUEST`.
- Invalid amount returns `INVALID_AMOUNT`.
- Invalid destination data returns `INVALID_DATA` or the agreed error code.
- Handler tests cover successful parsing and validation failures.

### Depends On

- Task 1.
- Task 2.
- Task 3.

## Task 5: Execute transfer using resolved internal account IDs

### Goal

Connect the new public request flow to the existing transfer domain behavior.

### Scope

- Resolve source and destination accounts before calling the transfer use case.
- Pass internal account IDs only inside backend application/domain layers.
- Preserve same-account validation.
- Preserve inactive-account validation.
- Preserve insufficient-funds validation.
- Preserve atomic debit/credit behavior.

### Acceptance Criteria

- Successful transfers still debit source and credit destination atomically.
- Same-account transfer returns `SAME_ACCOUNT_TRANSFER`.
- Inactive account returns `ACCOUNT_INACTIVE`.
- Insufficient funds returns `INSUFFICIENT_FUNDS`.
- Source or destination not found returns `ACCOUNT_NOT_FOUND` or the agreed
  error code.
- Existing domain invariants remain covered by tests.

### Depends On

- Task 4.

## Task 6: Return public transaction reference on transfer success

### Goal

Return a public `transaction_reference` from successful internal transfer
responses.

### Scope

- Define or reuse the persisted transaction reference field.
- Ensure successful transfer response includes `transaction_reference`.
- Avoid exposing internal transaction UUID as the public receipt lookup
  identifier.
- Keep response envelope shape unchanged.

### Acceptance Criteria

- Successful transfer response includes `transaction_reference`.
- Response does not require the mobile app to use internal transaction UUIDs.
- Idempotent replay returns the same historical transaction reference.
- Response tests cover normal success and replay success.

### Depends On

- Task 5.

## Task 7: Add receipt/detail lookup by transaction reference

### Goal

Expose receipt details through a dedicated endpoint that accepts
`transaction_reference`.

### Scope

- Add or expose receipt/detail route.
- Lookup persisted transaction data by `transaction_reference`.
- Return basic receipt fields:
  operation type, amount, status, transaction reference, operation date/time,
  source account, destination account number, recipient name when applicable,
  and optional description when available.
- Enforce authorization for receipt access.

### Acceptance Criteria

- Receipt details can be fetched by `transaction_reference`.
- Receipt response is built from persisted transaction data.
- Unauthorized receipt access is rejected.
- Unknown reference returns a typed not-found error.
- Receipt endpoint tests cover success, forbidden access, and not found.

### Depends On

- Task 6.

## Task 8: Align transfer idempotency with public request fields

### Goal

Keep idempotent transfer behavior correct after removing internal account IDs
from the public request.

### Scope

- Preserve `idempotency_key` support.
- Reconfirm scope as resolved source account plus idempotency key.
- Ensure duplicate public requests do not duplicate financial effects.
- Ensure replay returns historical transfer result and transaction reference.
- Update docs to explain idempotency without requiring `from_account_id`.

### Acceptance Criteria

- Same source account and same `idempotency_key` replay the historical result.
- Replays return the same `transaction_reference`.
- Different source account with same `idempotency_key` remains independently
  scoped.
- Idempotency tests use the new public request fields.

### Depends On

- Task 5.
- Task 6.

## Task 9: Update REST API documentation

### Goal

Update `api/docs/07-api-rest.md` to reflect the new transfer and receipt
contracts.

### Scope

- Update `POST /accounts/transfer` request example.
- Remove `from_account_id` and `to_account_id` from public docs.
- Add `transaction_reference` to success response examples.
- Update transfer notes and idempotency notes.
- Add receipt/detail endpoint documentation.
- Update error scenarios by endpoint.

### Acceptance Criteria

- Documentation matches the implemented transfer request.
- Documentation includes receipt/detail lookup by `transaction_reference`.
- Documentation no longer instructs clients to submit internal account IDs for
  transfer.
- Error examples match implemented error codes and envelope shape.

### Depends On

- Task 1.
- Task 6.
- Task 7.
- Task 8.

## Task 10: Update transfer tests and regression coverage

### Goal

Ensure the transfer contract refactor is covered end to end.

### Scope

- Update handler tests for the new public request body.
- Update integration tests for successful internal transfer.
- Add validation tests for missing/invalid destination account number.
- Add authorization tests.
- Add same-account transfer tests using public request fields.
- Add insufficient-funds and inactive-account tests using public request fields.
- Add idempotent replay tests.
- Add receipt/detail endpoint tests.

### Acceptance Criteria

- Tests cover success, validation failure, authorization failure, not found,
  inactive account, insufficient funds, same-account transfer, and idempotent
  replay.
- Tests confirm clients do not need internal account IDs.
- Receipt tests confirm lookup by `transaction_reference`.
- Backend transfer and authorization tests pass.

### Depends On

- Task 4.
- Task 5.
- Task 7.
- Task 8.

## Suggested GitHub Project Order

1. Define the public internal transfer API contract.
2. Add backend account lookup by public banking identifiers.
3. Resolve the transfer source account from authenticated context.
4. Refactor the transfer handler request parsing.
5. Execute transfer using resolved internal account IDs.
6. Return public transaction reference on transfer success.
7. Add receipt/detail lookup by transaction reference.
8. Align transfer idempotency with public request fields.
9. Update transfer tests and regression coverage.
10. Update REST API documentation.
