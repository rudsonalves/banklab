# Installation Identity MVP Tasks

These tasks split the mobile implementation of installation identity into
implementation-sized steps aligned with the API contract.

The first cut covers local installation identity, HTTP propagation, login,
restricted installation registration, and approved error states. Installation
management with listing and revocation is intentionally deferred.

## Task 1/12: Add installation identity storage keys and local marker contract

### Goal

Define the local keys and local marker contract used to distinguish app updates
from new installations.

### Scope

- Add `StorageKeys.installationId` with value `banklab.installation.id`.
- Define the local marker mechanism used outside durable secure storage.
- Ensure the marker is not treated as a secret or identity.
- Ensure restore without a valid marker results in a new installation identity.
- Document that logout and credential cleanup must not delete
  `StorageKeys.installationId`.

### Acceptance Criteria

- The installation id key is available through `StorageKeys`.
- The local marker behavior is documented in code or tests.
- The implementation can distinguish an app update from a reinstall/clear data
  scenario.
- Session cleanup does not target the installation id key.

### Depends On

- None.

## Task 2/12: Implement installation identity service

### Goal

Create a service that resolves a stable UUID v4 for the current app
installation.

### Scope

- Create `InstallationIdentityService`.
- Read the existing installation id from `LocalSecureStorage`.
- Validate that the stored value is a canonical UUID v4.
- Generate a new UUID v4 when the local marker is missing, the value is
  missing, or the value is invalid.
- Persist the generated UUID v4 in `flutter_secure_storage`.
- Create or refresh the local installation marker after a successful resolve.
- Return an `AsyncResult<String>`.
- Do not log the full installation id.

### Acceptance Criteria

- First run creates and persists a UUID v4.
- App update with local marker present reuses the same UUID.
- Reinstall/clear data/restore without local marker generates a new UUID even if
  secure storage still returns an older value.
- Invalid stored values are replaced.
- Read/write failures return a failure result and do not silently continue.
- Logs do not include the complete UUID.

### Depends On

- Task 1.

## Task 3/12: Block bootstrap when installation identity fails

### Goal

Prevent login and API calls when the app cannot resolve a stable
`installation_id`.

### Scope

- Resolve installation identity during app bootstrap or before the first API
  workflow.
- Show a recoverable error when identity resolution fails.
- Provide a retry path.
- Do not send API requests without `X-Installation-Id`.

### Acceptance Criteria

- Login is unavailable while installation identity is unresolved.
- A storage/read/write/generation failure blocks the flow.
- Retrying can resolve the identity and unblock the app.
- No request is intentionally sent without `X-Installation-Id`.

### Depends On

- Task 2.

## Task 4/12: Add installation interceptor

### Goal

Attach `X-Installation-Id` to all requests handled by the main HTTP client.

### Scope

- Create `InstallationInterceptor`.
- Resolve the current installation id through `InstallationIdentityService`.
- Add `X-Installation-Id` to request headers.
- Register the interceptor in the main Dio/RestClient setup.
- Keep the existing auth interceptor behavior.
- Remove or leave unused the old commented device interceptor reference without
  reintroducing device terminology.

### Acceptance Criteria

- Login requests include `X-Installation-Id`.
- Authenticated requests include `X-Installation-Id`.
- The header value is canonical lowercase UUID v4.
- Existing `Authorization` behavior remains unchanged.
- Interceptor failures block the request through normal app error handling
  instead of sending a headerless request.

### Depends On

- Task 2.
- Task 3.

## Task 5/12: Add installation header to token refresh

### Goal

Ensure `/auth/refresh` sends the same installation identity as normal requests.

### Scope

- Update `AuthInterceptor` refresh flow.
- Inject or otherwise access `InstallationIdentityService` for refresh.
- Add `X-Installation-Id` to the dedicated refresh Dio request.
- Preserve refresh de-duplication and retry behavior.
- Preserve token cleanup behavior on refresh failure.

### Acceptance Criteria

- Refresh requests include `X-Installation-Id`.
- Refresh still writes rotated access and refresh tokens.
- Refresh failure still clears session tokens.
- No refresh request is sent without an installation id when identity resolution
  fails.

### Depends On

- Task 2.
- Task 4.

## Task 6/12: Model operational and restricted login responses

### Goal

Represent login responses that can either create an operational session or
require installation registration.

### Scope

- Replace direct login parsing as only `LoggedUser` with a typed login result.
- Support operational response with `access_token` and `refresh_token`.
- Support restricted response with `restricted_access_token`,
  `restricted_token_type`, `restricted_scope`, and `restricted_expires_at`.
- Support `INSTALLATION_LIMIT_REACHED` as a typed app error.
- Preserve current handling of account approval and contact verification.
- Do not persist tokens for restricted or limit-reached outcomes.

### Acceptance Criteria

- Operational login continues to persist tokens and load profile.
- Restricted login does not persist operational tokens.
- Limit reached does not persist tokens.
- Existing login error tests still pass after updates.
- Parsing failures return `AppErrorCode.parsingError`.

### Depends On

- Task 4.

## Task 7/12: Add installation registration API

### Goal

Call `POST /security/installations` to exchange a restricted authorization for
an operational session.

### Scope

- Add API DTOs for installation registration request/response as needed.
- Send `Authorization: Bearer <restricted_access_token>`.
- Send `X-Step-Up-Token`.
- Send the same `X-Installation-Id` used by the app.
- Parse operational `access_token` and `refresh_token`.
- Parse installation metadata returned by the API.

### Acceptance Criteria

- Successful registration returns operational tokens.
- Missing or invalid step-up token maps to an app error.
- Installation mismatch and invalid installation id map to app errors.
- Registration does not send the transaction password.

### Depends On

- Task 5.
- Task 6.

## Task 8/12: Extend step-up operation for installation registration

### Goal

Allow the existing transaction password step-up flow to authorize installation
registration.

### Scope

- Add `StepUpOperation.installationRegistration`.
- Use method `POST`.
- Use path `/security/installations`.
- Reuse existing `/security/step-up/authorize` API.
- Ensure the transaction password is sent only to the step-up endpoint.

### Acceptance Criteria

- The app can request a step-up token for `POST /security/installations`.
- Step-up request body matches the API public operation contract.
- Transaction password is not sent to login or installation registration.
- Existing internal transfer step-up behavior remains unchanged.

### Depends On

- None.

## Task 9/12: Implement restricted installation certification flow

### Goal

Complete login for a known user on a new installation by asking for the
transaction password and registering the installation.

### Scope

- When login returns restricted installation registration, navigate to a
  certification state/screen.
- Ask for the transaction password.
- Authorize step-up for `POST /security/installations`.
- Call installation registration with the restricted token and step-up token.
- Persist operational tokens returned by registration.
- Load auth session and continue the normal post-login flow.
- Clear restricted token, step-up token, and password on success, failure, or
  cancellation.

### Acceptance Criteria

- Restricted login leads to the certification flow.
- Successful certification ends in an authenticated operational session.
- No operational tokens are persisted before registration succeeds.
- Cancel returns to login and requires a new login.
- Expired restricted token returns to login and requires a new login.

### Depends On

- Task 6.
- Task 7.
- Task 8.

## Task 10/12: Handle installation limit and transaction password blockers

### Goal

Show approved blocking states for limit reached and transaction password
preconditions.

### Scope

- Map `INSTALLATION_LIMIT_REACHED` to a typed app state or error.
- Show the approved limit reached message:
  `Esta conta já possui 3 instalações cadastradas. A instalação atual ainda não está autorizada.`
- Continue the message with:
  `Acesse sua conta por uma instalação já autorizada e remova uma instalação antiga para liberar espaço. Depois, tente entrar novamente neste app.`
- Use button text `Entendi`.
- Return to login after acknowledgement.
- Do not start step-up on limit reached.
- On `TRANSACTION_PASSWORD_NOT_SET` or `TRANSACTION_PASSWORD_LOCKED`, do not
  register the installation.
- Discard temporary restricted flow state for these blockers.

### Acceptance Criteria

- Limit reached shows the approved title, body, and button.
- Limit reached does not ask for transaction password.
- Transaction password not set does not call `POST /security/installations`.
- Transaction password locked does not call `POST /security/installations`.
- All blocker outcomes return the user to login or an informational state that
  requires login again.

### Depends On

- Task 6.
- Task 9.

## Task 11/12: Add safe telemetry and tests

### Goal

Make failures observable without leaking identifiers or credentials, and cover
the critical identity lifecycle.

### Scope

- Log storage/read/write/validation failures without full installation id.
- Do not log tokens.
- Do not log transaction password.
- Add tests for first run identity creation.
- Add tests for local-marker-present reuse.
- Add tests for local-marker-missing replacement of an old secure-storage value.
- Add tests for failure blocking behavior.
- Add tests for interceptor header injection.
- Add tests for refresh header injection.
- Add tests for restricted login and registration success.
- Add tests for limit reached and transaction password blockers.

### Acceptance Criteria

- Logs contain event context but no complete UUID, token, or password.
- Identity lifecycle tests pass.
- HTTP header tests pass.
- Auth repository/API tests cover operational, restricted, and blocked outcomes.

### Depends On

- Task 2.
- Task 4.
- Task 5.
- Task 9.
- Task 10.

## Task 12/12: Run mobile verification

### Goal

Confirm the mobile app remains healthy after the installation identity change.

### Scope

- Run `dart format` on changed Dart files.
- Run focused tests for identity, interceptors, auth parsing, and registration.
- Run `flutter analyze`.
- Run `flutter test` if feasible.
- Fix regressions introduced by the implementation.

### Acceptance Criteria

- Changed Dart files are formatted.
- Focused tests pass.
- `flutter analyze` passes.
- `flutter test` passes or any skipped/unavailable test run is documented.
- Existing login, refresh, logout, and step-up transfer behavior remain
  unchanged.

### Depends On

- Task 11.

## Deferred

- `GET /security/installations`.
- `DELETE /security/installations/{installation_resource_id}`.
- Installation management screen.
- Revoking an existing installation from mobile.
- Identifying and disabling removal of the current installation in management.
