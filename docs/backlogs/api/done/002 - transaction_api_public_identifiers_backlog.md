# Transfer API Public Identifiers Backlog

## 1. Objective

Refactor the BankLab internal transfer API so public requests use banking
identifiers instead of internal account IDs.

This backlog should be delivered before the mobile internal transfer backlog.
The mobile app should submit transfer instructions using the user-facing banking
domain: source account context and destination account number. Internal UUIDs
should remain an implementation detail resolved inside the backend application
layer.

## 2. Scope

Included:

- internal transfer request contract refactor;
- backend account resolution for transfer source and destination;
- idempotency alignment after source account resolution;
- successful operation response containing at least `transaction_reference`;
- receipt/detail lookup by `transaction_reference`;
- REST API documentation update;
- backend tests for the updated transfer contract.

Not included:

- deposit endpoint refactor;
- withdrawal endpoint refactor;
- balance or statement endpoint refactor;
- mobile screens;
- transactional password;
- external transfer implementation;
- TED, DOC, Pix, boleto, or card operations;
- replacing internal database IDs;
- removing internal IDs from backend repositories, domain services, ledger, or
  persistence models.

## 3. Current API Surface

Current internal transfer endpoint:

- `POST /accounts/transfer`

Current request fields:

- `from_account_id`
- `to_account_id`
- `amount`
- `idempotency_key`

Current success response focuses on account IDs and balances. It does not
provide a public transaction reference that the mobile app can use to fetch
receipt details.

## 4. Target Direction

The internal transfer API should accept public banking identifiers. For the
first mobile version:

- the source account is the selected account in the authenticated user's
  context;
- bank and branch are known from the selected source account context;
- the user enters only the destination account number;
- the backend resolves the destination account number within the known bank and
  branch context;
- the backend resolves internal IDs before executing the existing domain
  transfer behavior.

The backend should:

- validate the authenticated user and selected/source account;
- resolve the source and destination accounts to internal account IDs;
- enforce authorization after resolution;
- execute the existing transfer domain operation using internal IDs;
- preserve same-account, inactive-account, insufficient-funds, and forbidden
  validations;
- return at least a public `transaction_reference` on success;
- expose receipt/detail data through a lookup endpoint by
  `transaction_reference`.

Product decisions:

- Keep the existing `POST /accounts/transfer` route for the internal transfer
  contract refactor. The project is still in development, so a versioned or
  replacement route is not required for this change.
- Use `transaction_reference` as the public identifier returned by transfer and
  accepted by the receipt/detail endpoint. Do not expose the internal
  transaction UUID as the public receipt lookup identifier.

## 5. Epic 1: Public Transfer Contract

### Goal

Define the updated internal transfer request and response contract.

### Backlog Items

- Define the updated transfer request contract using destination account number
  instead of `to_account_id`.
- Define how the source account is identified from authenticated context and the
  selected account.
- Remove `from_account_id` and `to_account_id` from the public transfer request.
- Keep `amount` and `idempotency_key` in the transfer request.
- Define a successful transfer response containing at least
  `transaction_reference`.
- Preserve the existing error envelope format.
- Update REST API documentation for transfer request, response, notes, and error
  examples.

### Acceptance Criteria

- Transfer requests no longer require internal account IDs.
- Successful transfer responses include `transaction_reference`.
- Documentation includes updated request examples, success examples, and error
  examples.
- The contract remains explicitly internal-transfer only.

## 6. Epic 2: Backend Account Resolution

### Goal

Resolve transfer accounts inside the backend before executing the existing
domain operation.

### Backlog Items

- Resolve the selected/source account from authenticated context.
- Resolve the destination account by account number in the known bank and branch
  context.
- Keep authorization checks server-side after account resolution.
- Ensure same-account transfer validation still works after resolution.
- Ensure inactive account and insufficient-funds rules remain unchanged.
- Return clear errors for invalid, missing, or not-found destination account
  numbers.

### Acceptance Criteria

- The mobile app does not need internal account IDs to execute an internal
  transfer.
- Backend application code resolves account IDs before calling existing transfer
  domain behavior.
- Existing domain invariants remain protected.
- Invalid banking identifiers return clear business/API errors.

## 7. Epic 3: Operation Response And Receipt Lookup

### Goal

Return a public transaction reference from transfer calls and expose receipt
details by transaction reference.

### Backlog Items

- Ensure transfer success returns `transaction_reference`.
- Add or expose a receipt/detail endpoint that accepts `transaction_reference`.
- Include basic receipt fields in the receipt/detail response:
  operation type, amount, status, transaction reference, operation date/time,
  source account, destination account number, recipient name when applicable,
  and optional description when available.
- Ensure receipt lookup is authorized for the authenticated user.

### Acceptance Criteria

- Transfer success responses provide enough data for the mobile app to request a
  receipt.
- Receipt details can be loaded by `transaction_reference`.
- Receipt data is generated from persisted transaction data, not from client form
  input.
- Unauthorized receipt access is rejected.

## 8. Epic 4: Idempotency Alignment

### Goal

Keep internal transfer idempotency correct after replacing account IDs in the
public request.

### Backlog Items

- Keep `idempotency_key` support for internal transfer.
- Reconfirm idempotency scope after source account resolution.
- Ensure replay behavior still returns the historical transfer result.
- Update docs to explain idempotency using public request fields while preserving
  backend scoping by resolved source account.

### Acceptance Criteria

- Duplicate internal transfer requests with the same source account and
  idempotency key do not duplicate financial effects.
- Replay behavior remains consistent with the existing ledger design.
- Documentation explains the idempotency behavior clearly.

## 9. Epic 5: Tests And Compatibility

### Goal

Cover the transfer contract refactor with backend tests before the mobile
implementation depends on it.

### Backlog Items

- Update handler tests for transfer request parsing.
- Update integration tests for successful transfer.
- Add tests for invalid destination account number.
- Add authorization tests for attempts to transfer from an unavailable source
  account.
- Add same-account transfer tests using public request fields.
- Add idempotent replay tests using public request fields.
- Add receipt/detail endpoint tests.
- Update API documentation examples used by tests or docs.

### Acceptance Criteria

- Existing transfer behavior remains valid through public banking identifiers.
- Tests cover success, validation failure, authorization failure, not-found
  account, inactive account, insufficient funds, same-account transfer, and
  idempotent replay.
- Mobile can rely on the updated transfer contract without sending internal
  account IDs.

## 10. Suggested Delivery Order

1. Define updated transfer API contract and examples.
2. Implement backend source account resolution from authenticated context.
3. Implement backend destination account lookup by account number.
4. Refactor internal transfer endpoint request handling.
5. Add `transaction_reference` to successful transfer responses.
6. Add receipt/detail lookup by `transaction_reference`.
7. Update idempotency documentation and tests.
8. Update REST API documentation.
9. Run backend transfer and authorization tests.

## 11. Open Decisions

No open product decisions remain for this backlog.

## 12. Definition Of Done

- Internal transfer can be executed without internal account IDs in the public
  API request.
- Successful transfer responses include `transaction_reference`.
- Receipt details can be fetched by transaction reference.
- API documentation reflects the updated transfer contract.
- Backend tests cover the updated transfer contract and receipt lookup.
