# Login Requires Approved Account Tasks

Parent backlog:

- `004 - login_requires_approved_account_backlog.md`

Suggested project fields for all tasks:

- Status: Backlog
- Priority: High
- Area: API
- Type: Architecture/Security

## Task 1/6: Add the account approval required error contract

Status: Done

### Goal

Create the API error contract that mobile will use to identify users who cannot
login because approval/account provisioning is incomplete.

### Scope

- Add an auth domain/application error for account approval required.
- Add the shared error code `ACCOUNT_APPROVAL_REQUIRED`.
- Register the error in the auth error registry.
- Map the error to `403 Forbidden`.
- Use the response message `Account approval required` unless a better product
  text is defined before implementation.
- Keep invalid credentials behavior unchanged.

### Acceptance Criteria

- The shared error registry maps the new error to:
  - code: `ACCOUNT_APPROVAL_REQUIRED`;
  - status: `403 Forbidden`;
  - message: `Account approval required`.
- Invalid email/password still returns the existing invalid credentials error.
- The response follows the existing `data`/`error` envelope.

### Depends On

- None.

## Task 2/6: Add a login account provisioning dependency

Status: Done

### Goal

Allow the login use case to check whether a customer has at least one account
without depending on account delivery or account infrastructure details.

### Scope

- Define a small auth-local interface for account provisioning checks.
- The interface should be compatible with the existing account repository method:
  `ExistsByCustomerID(ctx, customerID)`.
- Inject this dependency into `LoginUserUseCase`.
- Wire the existing account repository into the login use case from API startup.
- Preserve admin login behavior when no customer/account is expected.

### Acceptance Criteria

- `LoginUserUseCase` can verify account existence by `customer_id`.
- The dependency is expressed as a minimal interface in the auth application
  boundary.
- API startup passes the existing account repository to login.
- No HTTP handler performs this business rule directly.

### Depends On

- None.

## Task 3/6: Enforce customer login eligibility

Status: Done

### Goal

Prevent customer users from receiving tokens or refresh sessions until approval
and account provisioning are complete.

### Scope

- After password validation and before token generation, validate customer login
  eligibility.
- Customer users must have an allowed status for login.
- Customer users must have a non-nil `customer_id`.
- Customer users must have at least one account for that `customer_id`.
- Pending customer users return `ACCOUNT_APPROVAL_REQUIRED`.
- Active customer users without an account return `ACCOUNT_APPROVAL_REQUIRED`.
- Blocked user handling is outside this backlog and should not be redesigned
  here.
- Admin users are not blocked by missing `customer_id` or missing account.
- Do not create access tokens, refresh tokens, or refresh sessions when the
  approval-required rule fails.

### Acceptance Criteria

- Pending customer users with valid credentials cannot login.
- Active customer users without accounts cannot login.
- Active customer users with at least one account can login.
- Admin users can login without customer/account data.
- Approval-required failures happen before token/session creation.
- Repository errors from the account check do not get converted into
  `ACCOUNT_APPROVAL_REQUIRED`.

### Depends On

- Task 1.
- Task 2.

## Task 4/6: Add login use-case tests for approval/account requirements

Status: Done

### Goal

Cover the new login lifecycle rule at the application layer.

### Scope

- Add tests for pending customer users with valid credentials.
- Add tests for active customer users with no account.
- Add tests for active customer users with an account.
- Add tests for admin users without customer/account data.
- Add tests proving no token generation occurs when approval is required.
- Add tests proving no refresh session is created when approval is required.
- Add tests proving account repository errors remain internal/wrapped errors.

### Acceptance Criteria

- The test suite fails if login issues tokens for a pending customer.
- The test suite fails if login issues tokens for a customer without an account.
- The test suite fails if a refresh session is created on approval-required
  failure.
- Existing successful login scenarios continue to pass.

### Depends On

- Task 3.

## Task 5/6: Add delivery/error mapping coverage

Status: Backlog

### Goal

Prove that the login endpoint exposes the new mobile-facing error contract.

### Scope

- Add or update auth delivery tests for the new error mapping.
- Validate the HTTP status is `403 Forbidden`.
- Validate the response error code is `ACCOUNT_APPROVAL_REQUIRED`.
- Validate invalid credentials still map to the existing invalid credentials
  response.
- Avoid exposing whether email/password were valid in invalid credentials cases.

### Acceptance Criteria

- Login delivery tests cover the approval-required response envelope.
- Mobile can rely on `error.code == "ACCOUNT_APPROVAL_REQUIRED"`.
- Invalid credentials behavior remains unchanged.

### Depends On

- Task 1.
- Task 3.

## Task 6/6: Update API documentation and run affected tests

Status: Backlog

### Goal

Document the login approval requirement and verify the API still behaves
correctly.

### Scope

- Update REST API documentation for `POST /auth/login`.
- Document the `ACCOUNT_APPROVAL_REQUIRED` error response.
- Mention that approval happens through `POST /admin/users/{id}/approve`.
- Do not edit implemented API/mobile reports unless explicitly requested.
- Run auth application tests.
- Run auth delivery tests.
- Run affected account/admin tests if wiring changes require it.
- Run the full API test suite if feasible.

### Acceptance Criteria

- API documentation tells mobile which error code to handle.
- Documentation states that users without approved/provisioned accounts cannot
  login.
- Auth tests pass.
- Affected API tests pass.
- Existing successful login, admin approval, and account provisioning behavior
  are not regressed.

### Depends On

- Task 4.
- Task 5.

## Suggested GitHub Project Order

1. Add the account approval required error contract.
2. Add a login account provisioning dependency.
3. Enforce customer login eligibility.
4. Add login use-case tests for approval/account requirements.
5. Add delivery/error mapping coverage.
6. Update API documentation and run affected tests.
