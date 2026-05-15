# Login Approval Required Handling Tasks

These tasks split the mobile handling of `ACCOUNT_APPROVAL_REQUIRED` into
implementation-sized steps.

The goal is to let full login and short login show a specific approval-pending
message when the API blocks login because the customer account has not been
approved/provisioned yet.

## Task 1/7: Add account approval required app error code

### Goal

Represent the backend `ACCOUNT_APPROVAL_REQUIRED` code as a typed mobile error.

### Scope

- Add `AppErrorCode.accountApprovalRequired`.
- Map backend error code `ACCOUNT_APPROVAL_REQUIRED` to this app error code.
- Preserve the approved user-facing message:
  `Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você poderá acessar o app.`
- Keep `INVALID_CREDENTIALS` mapping unchanged.
- Do not map generic `403` responses to approval-required unless the backend
  error code is `ACCOUNT_APPROVAL_REQUIRED`.

### Acceptance Criteria

- The app can distinguish `ACCOUNT_APPROVAL_REQUIRED` through
  `AppErrorCode.accountApprovalRequired`.
- Invalid credentials still map to the existing invalid credentials error.
- Generic forbidden errors still use the existing forbidden/generic behavior.
- The mapped message uses the approved copy.

### Depends On

- None.

## Task 2/7: Add focused error mapping tests

### Goal

Prove that API error mapping distinguishes approval-required from invalid
credentials and generic forbidden errors.

### Scope

- Add or update tests around the HTTP/API error mapper.
- Cover backend code `ACCOUNT_APPROVAL_REQUIRED`.
- Cover backend code `INVALID_CREDENTIALS`.
- Cover a generic `403` without `ACCOUNT_APPROVAL_REQUIRED`.
- Verify the resulting `AppErrorCode` and message.

### Acceptance Criteria

- `ACCOUNT_APPROVAL_REQUIRED` maps to
  `AppErrorCode.accountApprovalRequired`.
- `INVALID_CREDENTIALS` behavior remains unchanged.
- Generic `403` does not map to account approval required.
- Tests cover the approved user-facing message.

### Depends On

- Task 1.

## Task 3/7: Handle approval-required in full login

### Goal

Show the approval-pending message on the full login screen when login fails with
`AppErrorCode.accountApprovalRequired`.

### Scope

- Update the full login page error handling.
- Use `AppSnackbar.show(...)` from
  `mobile/lib/uis/core/messages/app_snackbar.dart`.
- Show the approved copy for `AppErrorCode.accountApprovalRequired`.
- Do not navigate to home on this error.
- Do not show the invalid credentials message for this condition.
- Preserve existing full login behavior for invalid credentials, validation
  errors, and generic failures.

### Acceptance Criteria

- Full login shows the approved approval-pending message through `AppSnackbar`.
- Full login stays on the login screen.
- Full login does not treat this error as a wrong password.
- Other login errors continue to behave as before.

### Depends On

- Task 1.

## Task 4/7: Handle approval-required in short login

### Goal

Show the approval-pending message on the remembered short login screen when login
fails with `AppErrorCode.accountApprovalRequired`.

### Scope

- Update the short login page error handling.
- Use `AppSnackbar.show(...)`.
- Show the approved copy for `AppErrorCode.accountApprovalRequired`.
- Keep the remembered identity visible.
- Keep the "use another account" path available.
- Do not clear remembered login cache.
- Do not navigate to home on this error.
- Preserve existing short login behavior for invalid credentials and generic
  failures.

### Acceptance Criteria

- Short login shows the approved approval-pending message through `AppSnackbar`.
- Short login keeps the remembered user identity visible.
- Short login does not clear the cache.
- Short login does not treat this error as a wrong password.
- The user can still navigate to full login to use another account.

### Depends On

- Task 1.

## Task 5/7: Verify repository success-only side effects

### Goal

Ensure approval-required failures do not execute post-login success behavior.

### Scope

- Review `AuthRepositoryImpl.login()`.
- Confirm profile loading happens only after login success.
- Confirm remembered login cache is updated only after profile loading succeeds.
- Confirm auth tokens are not persisted when login fails.
- Add or update repository tests if the current test structure supports it.
- Do not change backend contracts.

### Acceptance Criteria

- Approval-required login failure does not save access or refresh tokens.
- Approval-required login failure does not call `getProfile`.
- Approval-required login failure does not update remembered login cache.
- Successful login still loads profile and updates remembered login identity.

### Depends On

- Task 1.

## Task 6/7: Validate login UI feedback behavior

### Goal

Verify the visible behavior of full login and short login for the new error.

### Scope

- Validate full login with `AppErrorCode.accountApprovalRequired`.
- Validate short login with `AppErrorCode.accountApprovalRequired`.
- Validate invalid credentials still show credential-specific feedback.
- Validate generic failures still show the existing generic feedback.
- Use the current app testing style; add widget tests if practical and useful.

### Acceptance Criteria

- Full login shows approval-specific feedback.
- Short login shows approval-specific feedback.
- Both screens use `AppSnackbar` for this feedback.
- Invalid credentials still look and behave differently.
- Generic failures still follow existing behavior.

### Depends On

- Task 3.
- Task 4.

## Task 7/7: Run mobile checks

### Goal

Confirm the mobile app remains healthy after the error handling change.

### Scope

- Run `dart format` on changed Dart files.
- Run focused Flutter tests for error mapping/auth UI if available.
- Run `flutter analyze`.
- Run `flutter test` if feasible.
- Fix regressions introduced by the implementation.

### Acceptance Criteria

- Changed Dart files are formatted.
- Focused tests pass.
- `flutter analyze` passes.
- `flutter test` passes or any skipped/unavailable test run is documented.
- Successful login behavior remains unchanged.

### Depends On

- Task 2.
- Task 5.
- Task 6.

