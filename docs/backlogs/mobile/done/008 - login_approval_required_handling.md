# Login Approval Required Handling

## Problem Statement

The API now blocks customer login when the user has valid credentials but has not
yet completed the admin approval/account provisioning flow. In this state,
`POST /auth/login` returns:

```json
{
  "data": null,
  "error": {
    "code": "ACCOUNT_APPROVAL_REQUIRED",
    "message": "Account approval required"
  }
}
```

The HTTP status is `403 Forbidden`.

The mobile app must not treat this condition as an invalid password or generic
login failure. It should identify the API error code and show guidance that the
user still needs account approval before entering the app.

## Goal

Handle `ACCOUNT_APPROVAL_REQUIRED` in the mobile authentication flow and present
a clear, specific message when login is blocked because the account has not been
approved/provisioned yet.

## Existing Context

- Full login and short login both call `AuthRepository.login()`.
- `AuthRepositoryImpl.login()` now loads the user profile after successful
  authentication and updates the remembered login cache.
- Short login uses the remembered identifier and asks only for the password.
- API errors are mapped through the mobile HTTP/error layer before reaching UI
  state.
- The backend contract for this condition is stable:
  - HTTP status: `403`;
  - error code: `ACCOUNT_APPROVAL_REQUIRED`;
  - message: `Account approval required`.

## Non-Goals

- Do not implement the admin approval flow in mobile.
- Do not allow the user to create or request an account from mobile.
- Do not change the backend login endpoint.
- Do not change refresh-token behavior.
- Do not implement blocked-user messaging in this backlog.
- Do not remove or clear remembered login cache automatically unless a concrete
  UX reason is defined.

## Product Behavior

When login fails with `ACCOUNT_APPROVAL_REQUIRED`, the app should communicate:

- the credentials were accepted enough for the API to identify the account state;
- access is not available yet;
- the user's account still needs approval/provisioning;
- the user should wait for approval or contact support/admin, depending on the
  final product text.

The message should be different from the invalid credentials message.

Approved user-facing copy:

```text
Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você
poderá acessar o app.
```

The UI must not show a "wrong password" style message for this condition.

## Epic 1: Mobile Error Code Mapping

### Goal

Represent `ACCOUNT_APPROVAL_REQUIRED` as a first-class mobile error condition.

### Scope

- Add `ACCOUNT_APPROVAL_REQUIRED` to the mobile API/app error mapping as
  `AppErrorCode.accountApprovalRequired`.
- Preserve the raw backend error code where the current architecture supports
  it.
- Map it to a semantic app error that authentication UI can inspect through
  `AppErrorCode.accountApprovalRequired`.
- Keep `INVALID_CREDENTIALS` behavior unchanged.
- Avoid collapsing this error into generic forbidden or generic network failure.

### Acceptance Criteria

- The mobile error layer can distinguish `ACCOUNT_APPROVAL_REQUIRED` through
  `AppErrorCode.accountApprovalRequired`.
- Invalid credentials still map to the current invalid credentials behavior.
- Generic `403` errors do not automatically become approval-required errors
  unless the backend code is `ACCOUNT_APPROVAL_REQUIRED`.

## Epic 2: Full Login UI Handling

### Goal

Show a specific approval-pending message on the full login screen.

### Scope

- Update the full login view model/state handling for the new error.
- Display approval-required copy when `AuthRepository.login()` fails with the
  mapped error.
- Use `AppSnackbar` as the UI feedback mechanism.
- Do not navigate to home.
- Do not save remembered login cache on this failure.
- Preserve existing behavior for invalid credentials, invalid input, and generic
  failures.

### Acceptance Criteria

- Full login displays a specific approval-pending message.
- Full login uses `AppSnackbar` for this feedback.
- Full login does not show a password/credential error for this condition.
- Full login does not navigate to authenticated areas.
- Full login keeps existing behavior for other errors.

## Epic 3: Short Login UI Handling

### Goal

Handle approval-required failures correctly when the user logs in from the
remembered login screen.

### Scope

- Update short login view model/state handling for the new error.
- Display approval-required copy using the same semantic condition as full
  login.
- Use `AppSnackbar` as the UI feedback mechanism.
- Keep the remembered identity visible unless a later UX decision says
  otherwise.
- Keep the "use another account" path available.
- Do not navigate to home.

### Acceptance Criteria

- Short login displays a specific approval-pending message.
- Short login uses `AppSnackbar` for this feedback.
- Short login does not treat this as a wrong password.
- Short login keeps the user able to switch to full login.
- Remembered login cache is not cleared by this condition.

## Epic 4: Repository Flow Safety

### Goal

Ensure failed login due to approval-required state does not run post-login
success steps.

### Scope

- Confirm `AuthRepositoryImpl.login()` only loads profile after login success.
- Confirm profile loading is not attempted after `ACCOUNT_APPROVAL_REQUIRED`.
- Confirm remembered login cache is updated only after profile loading succeeds.
- Confirm auth tokens are not persisted when login fails.

### Acceptance Criteria

- Approval-required failures do not save tokens.
- Approval-required failures do not load profile.
- Approval-required failures do not update remembered login cache.
- Existing successful login still loads profile and updates remembered identity.

## Epic 5: Validation

### Scenarios

Full login with user still pending approval:

```text
Login -> API returns ACCOUNT_APPROVAL_REQUIRED -> show approval-pending message
```

Short login with remembered identity still pending approval:

```text
ShortLogin -> API returns ACCOUNT_APPROVAL_REQUIRED -> show approval-pending message
```

Invalid password:

```text
Login/ShortLogin -> API returns INVALID_CREDENTIALS -> show invalid credentials message
```

Successful login after approval:

```text
Login/ShortLogin -> login succeeds -> profile loads -> cache updates -> Home
```

## Acceptance Criteria

- The mobile app recognizes `ACCOUNT_APPROVAL_REQUIRED`.
- The mobile app maps it to `AppErrorCode.accountApprovalRequired`.
- Full login and short login show approval-specific feedback.
- Full login and short login use `AppSnackbar` for this feedback.
- Invalid credentials still show credential-specific feedback.
- Generic failures still show the existing generic error behavior.
- Approval-required failures do not persist tokens, load profile, update cache,
  or navigate to home.
- Successful login behavior remains unchanged after approval.

## Resolved Decisions

- Use the approved copy: `Sua conta ainda está aguardando aprovação. Assim que
  ela for liberada, você poderá acessar o app.`
- Do not clear remembered login cache on `ACCOUNT_APPROVAL_REQUIRED`.
- Add `AppErrorCode.accountApprovalRequired`.
- Use `AppSnackbar` from
  `mobile/lib/uis/core/messages/app_snackbar.dart` as the standard feedback
  mechanism for this message.
- Blocked-user handling remains outside this backlog.
