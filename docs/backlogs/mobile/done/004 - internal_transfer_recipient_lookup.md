# Task: Add internal transfer recipient lookup API

## Goal

Allow the mobile app to search and confirm eligible recipient accounts before
executing an internal transfer between accounts inside the same bank.

This lookup is part of the internal transfer flow only. It is not a Pix, TED,
DOC, or interbank account discovery endpoint.

## Flow

1. The user selects the paying account in the app.
2. The user searches the recipient by branch + account number, or by CPF.
3. The API returns one or more eligible internal recipient accounts with the
   minimum data needed for user confirmation.
4. If more than one account is returned, the user must select the recipient
   account.
5. The app executes the transfer using `from_account_id` from the selected
   paying account and `to_account_id` from the selected recipient account.

Branch, account number, and document are lookup/confirmation inputs. Transfer
execution must use account IDs.

## Scope

- Add an authenticated endpoint to search accounts inside the same bank.
- Support lookup by branch + account number.
- Support lookup by CPF document.
- Keep CNPJ and legal person account lookup out of scope for this backlog.
- Return only transfer-safe public confirmation data.
- Return one or more eligible recipient accounts.
- Do not return balances, customer IDs, internal transaction IDs, e-mail,
  phone, full document, or unrelated account/customer data.
- Keep this endpoint scoped to internal transfers.
- Keep transfer execution based on `from_account_id` and `to_account_id`.

## Implementation Workflow

Use implementation-first sequencing for this backlog.

1. Implement the endpoint and supporting application/repository code until the
   main behavior compiles and follows the intended contract.
2. Review the resulting design before writing tests: endpoint path, query
   modes, response shape, error behavior, normalization, and privacy limits.
3. Adjust the implementation after the design review.
4. Write or update tests only after the behavior and contract are stable.
5. Run the focused package tests and then the broader backend suite.

During implementation, ignore test creation and test maintenance unless a
compile error blocks progress. Tests are acceptance proof, not the mechanism
used to shape the first implementation pass.

## Endpoint

`GET /accounts/internal-transfers/recipients`

## Query Parameters

- `branch`
- `account_number`
- `document`

Rules:

- `branch` and `account_number` must be provided together.
- `document` can be used alone and currently accepts CPF only.
- Reject requests with neither search mode.
- Reject ambiguous mixed modes unless explicitly supported.
- Normalize branch, account number, and CPF before querying.

## Response

```json
{
  "accounts": [
    {
      "account_id": "acc_123",
      "holder_name": "Maria Silva",
      "document": "***.456.789-**",
      "branch": "0001",
      "account_number": "12345-6",
      "account_type": "checking"
    }
  ]
}
```

Response rules:

- `account_id` is the account identifier used later as `to_account_id`.
- `document` must be masked.
- Lookup by branch + account number should return zero or one account.
- Lookup by CPF may return zero, one, or many accounts.
- If multiple accounts are returned, the mobile UI must require user selection.

## Transfer Execution

The execution endpoint for internal transfers should use account IDs, not
branch + account number:

`POST /accounts/internal-transfers`

```json
{
  "from_account_id": "acc_001",
  "to_account_id": "acc_123",
  "amount": 2500,
  "description": "Aluguel",
  "idempotency_key": "01HY..."
}
```

Backend validation:

- Validate that the authenticated user can operate `from_account_id`.
- Validate that `to_account_id` exists and is eligible to receive an internal
  transfer.
- Validate that `from_account_id` and `to_account_id` are not the same account.
- Validate account status, balance, limits, and idempotency.
- Do not use branch + account number as execution identifiers.

## Backend Requirements

- Require authenticated user.
- Search only internal bank accounts.
- Return only accounts eligible to receive transfers.
- Exclude closed, blocked, inactive, or non-transferable accounts.
- Decide whether the sender’s own accounts should appear in lookup results.
- Mask CPF in the response.
- Add rate limiting or abuse protection to reduce account/document enumeration.
- Add audit logging for lookup attempts if the project already has audit conventions.
- Keep transfer execution endpoint using `from_account_id` and `to_account_id`.

## Acceptance Criteria

- Endpoint is `GET /accounts/internal-transfers/recipients`.
- Lookup by branch + account number returns zero or one eligible account.
- Lookup by CPF returns zero, one, or many eligible accounts.
- Response includes `account_id`, `holder_name`, masked `document`, `branch`,
  `account_number`, and optional account type.
- Response does not expose customer ID, balance, full document, phone, e-mail,
  or internal-only fields.
- Multiple accounts for the same CPF can be returned for user selection.
- Invalid query combinations return typed API errors.
- Unauthorized requests are rejected.
- Internal transfer execution remains ID-based with `from_account_id` and
  `to_account_id`.
- Tests cover successful lookup by account, successful lookup by document,
  multiple accounts for one CPF, no results, invalid parameters,
  unauthorized request, and ineligible accounts.

## Test Timing

Tests should be implemented at the end of the task, after the endpoint behavior
has been reviewed and stabilized.
