# Login Requires Approved Account Backlog

## 1. Objective

Prevent customer users from logging in before their onboarding approval has
created at least one bank account.

In BankLab, a customer without an account is not a valid user state for the
mobile app. Account creation is part of the admin approval flow:

```http
POST /admin/users/{id}/approve
```

Therefore, login must not succeed for a customer that has valid credentials but
has not yet been approved and provisioned with an account.

## 2. Problem

The current login flow validates email and password, generates tokens, creates a
refresh session, and returns authenticated user data.

This allows a user whose onboarding is still pending, or whose account was not
provisioned correctly, to receive valid tokens. The mobile app can then enter an
authenticated state even though the user cannot actually use the banking
features that require an account.

This creates an ambiguous product state:

- the credentials are correct, but the user is not ready to use the app;
- the mobile app cannot reliably distinguish invalid credentials from pending
  approval;
- the user receives no clear guidance that admin approval is required;
- downstream account, balance, statement, and transfer screens may fail later
  with less helpful errors.

The API must return a specific, stable error that the mobile app can identify and
turn into an appropriate message.

## 3. Fixed Premises

This backlog should not reopen decisions that already belong to the current
architecture.

Fixed premises:

- the first account is created during admin approval;
- a customer user without an account must not be allowed to complete login;
- account approval remains an admin action through
  `POST /admin/users/{id}/approve`;
- the mobile app needs a specific API error code for this condition;
- invalid credentials must remain indistinguishable from wrong email/password
  cases where appropriate;
- admin users are not governed by the customer account requirement;
- this change should not move account creation back into a customer-facing login
  or account endpoint.

## 4. Scope

Included:

- add a login-time validation that customer users must be approved and have at
  least one account;
- define a specific auth/domain error for users who cannot login because account
  approval/provisioning is incomplete;
- map the new error to a stable API error code and HTTP status;
- ensure the mobile app can identify the condition from the error envelope;
- update login tests for pending users, approved users with accounts, and
  inconsistent users without accounts;
- update API documentation for the login error contract.

Not included:

- mobile implementation;
- admin UI implementation;
- changing the onboarding approval endpoint;
- creating accounts during login;
- allowing customers to create their own first account;
- redesigning account provisioning;
- changing refresh-token behavior for already authenticated sessions;
- changing account, balance, statement, or transfer endpoints.

## 5. Current Behavior

Current login behavior:

- receives email and password;
- normalizes email;
- finds the user by email;
- compares password hash;
- generates access and refresh tokens;
- creates a refresh session;
- returns user identity fields, including role and optional `customer_id`.

Current violation:

- login does not verify whether the customer user has an account;
- login does not return a specific error for pending approval or missing account;
- the mobile app cannot distinguish this business condition from other failures.

## 6. Target Behavior

For customer users, login should succeed only when the user is in a business
state that can use the mobile app:

- credentials are valid;
- user status allows login;
- user has a `customer_id`;
- at least one account exists for that `customer_id`.

If the user still needs admin approval, or if provisioning did not create an
account, login should fail before token generation and before refresh-session
creation.

The API should return a specific error code, for example:

```json
{
  "data": null,
  "error": {
    "code": "ACCOUNT_APPROVAL_REQUIRED",
    "message": "Account approval required"
  }
}
```

The HTTP status for this condition is:

- `403 Forbidden`: credentials may be valid, but access to the app is not
  allowed until account approval/provisioning is complete.

## 7. Architectural Positioning

The rule belongs to the login application flow because login is the boundary
that decides whether an authenticated session may be created.

The rule should not be enforced only in mobile, because mobile clients cannot be
trusted to protect server-side session creation. It should also not be delayed
until account listing or balance retrieval, because by then the API has already
created valid tokens and a refresh session.

The implementation will likely require the auth login use case to check account
existence for customer users. This introduces a dependency from auth application
logic to an account-facing contract, so the dependency should be explicit and
minimal.

Preferred direction:

- define a small account lookup/existence interface needed by login;
- inject it into `LoginUserUseCase`;
- keep the dependency at the application boundary, not in HTTP delivery;
- return an auth-domain error that is mapped by the shared error registry.

## 8. Epic 1: Login Eligibility Rule

### Goal

Define and enforce the business rule that a customer user can only login after
account approval/provisioning is complete.

### Backlog Items

- Define what user statuses may complete login.
- Ensure pending or blocked customer users do not receive tokens.
- Require customer users to have a `customer_id`.
- Require customer users to have at least one account for their `customer_id`.
- Ensure admin users are not blocked by missing customer/account data.
- Ensure token generation and refresh-session creation happen only after the
  eligibility rule passes.

### Acceptance Criteria

- A pending customer with valid credentials cannot login.
- An approved customer with at least one account can login.
- A customer with valid credentials but no account cannot login.
- Admin login behavior remains valid.
- Failed eligibility does not create refresh sessions.

## 9. Epic 2: Specific API Error Contract

### Goal

Expose a stable error that mobile can identify and translate into a user-facing
message.

### Backlog Items

- Add a domain/application error for account approval required.
- Add a shared error code such as `ACCOUNT_APPROVAL_REQUIRED`.
- Register the error in the auth error registry.
- Define the HTTP status for this condition.
- Document the login error response.
- Keep invalid credentials behavior unchanged.

### Acceptance Criteria

- The login endpoint returns a stable error code for approval/account-required
  failures.
- Mobile can distinguish this condition from invalid credentials.
- The response follows the existing `data`/`error` envelope.
- The error message does not expose sensitive credential-validation details.

## 10. Epic 3: Account Existence Check

### Goal

Allow the auth login use case to verify account provisioning without taking a
large dependency on account internals.

### Backlog Items

- Define the minimal account dependency required by login.
- Reuse the existing account repository method if it is appropriate.
- Ensure the check is based on `customer_id`, not on client-provided data.
- Handle repository errors without turning them into misleading approval errors.
- Keep the query efficient and deterministic.

### Acceptance Criteria

- Login checks account existence using server-side user/customer data.
- Missing account maps to the approval-required error.
- Account repository errors remain internal/server errors.
- No client input can bypass the check.

## 11. Epic 4: Tests And Documentation

### Goal

Cover the new lifecycle boundary and document the mobile-facing contract.

### Backlog Items

- Add login use-case tests for pending customer users.
- Add login use-case tests for active customer users without accounts.
- Add login use-case tests for active customer users with accounts.
- Add tests proving tokens and sessions are not created on approval-required
  failures.
- Add delivery tests for the new error mapping.
- Update REST API documentation for login errors.

### Acceptance Criteria

- Tests fail if login creates tokens for an unapproved/no-account customer.
- Tests fail if the approval-required error is mapped incorrectly.
- Documentation tells mobile which error code to handle.
- Existing successful login tests continue to pass.

## 12. Open Implementation Questions

The following implementation decisions are resolved:

- `ACCOUNT_APPROVAL_REQUIRED` returns `403 Forbidden`.
- Pending customer users and active customer users without an account return the
  same public error code: `ACCOUNT_APPROVAL_REQUIRED`.
- Blocked users are not part of this backlog. User blocking is a different
  lifecycle state and should be handled separately later.
- Refresh token revalidation is not part of this backlog. This backlog only
  changes initial login/session creation.
- The auth login application should depend on a small auth-local account
  provisioning interface, for example:

```go
type AccountProvisioningChecker interface {
    ExistsByCustomerID(ctx context.Context, customerID uuid.UUID) (bool, error)
}
```

The existing account repository can satisfy this contract because it already
exposes `ExistsByCustomerID`.

There are no remaining open implementation questions for this backlog.
