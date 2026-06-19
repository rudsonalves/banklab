# Changelog

## 2026/06/19 - mobile/installation-identity-04

This change advances the mobile installation identity flow by introducing typed authentication states, distinguishing operational sessions from restricted installation authorization responses.

The update also adds the installation registration API client, propagates installation identity during token refresh, maps installation-related backend errors, and extends tests around login parsing, restricted authentication, registration headers, and refresh behavior.

1. **docs/backlogs/mobile/013 - installation-identity-mvp.md**

   * Updated the restricted login flow to explicitly represent restricted responses as `RestrictedInstallationAuthState`.
   * Clarified that restricted authentication must not create an operational session.
   * Reordered the installation registration steps to keep restricted tokens, step-up authorization, operational token persistence, and token cleanup explicit.

2. **docs/backlogs/mobile/013 - installation-identity-mvp_tasks.md**

   * Updated the task scope from parsing login only as an operational user to using an `AuthState` hierarchy.
   * Documented support for operational, anonymous, and restricted installation authentication states.
   * Clarified that restricted login data is represented by `RestrictedInstallationAuthState`.

3. **mobile/lib/domain/common/auth/models/auth_state.dart**

   * Added the `AuthState` sealed hierarchy.
   * Added `OperationalAuthState` for token-bearing authenticated sessions.
   * Added `RestrictedInstallationAuthState` for restricted installation registration authorization.
   * Added `AnonymousAuthState` as the default unauthenticated state.
   * Added login response shape detection through `AuthState.fromLoginMap`.

4. **mobile/lib/domain/common/auth/models/auth_user.dart**

   * Removed the previous `AuthUser`, `LoggedUser`, and `NotLoggedUser` model hierarchy.
   * Replaced it with the new `AuthState` model structure.

5. **mobile/lib/domain/AGENT.md**

   * Updated domain documentation to reference `auth_state.dart`.
   * Documented the new authentication state model with operational, restricted installation, and anonymous states.
   * Clarified that restricted installation authentication must not be treated as an operational session.

6. **mobile/lib/core/result/errors/app_error_code.dart**

   * Added `installationRegistrationRequired`.
   * Added `installationLimitReached`.
   * Introduced explicit app-level error codes for installation identity login outcomes.

7. **mobile/lib/core/services/client_http/dio/dio_error_mapper.dart**

   * Added mapping for backend `INSTALLATION_LIMIT_REACHED`.
   * Converted this backend error into `AppErrorCode.installationLimitReached`.
   * Preserved backend error details in the mapped `AppError`.

8. **mobile/lib/core/services/client_http/interceptors/auth/auth_interceptor.dart**

   * Injected `InstallationIdentityService` into `AuthInterceptor`.
   * Resolved the installation identity before refreshing tokens.
   * Added `X-Installation-Id` to refresh token requests.
   * Prevented refresh requests when installation identity resolution fails.

9. **mobile/lib/core/services/core_services.dart**

   * Updated dependency registration for `AuthInterceptor`.
   * Injected `InstallationIdentityService` into the authentication interceptor factory.

10. **mobile/lib/data/repositories/auth/auth_repository.dart**

* Updated `currentUser` from `AuthUser` to `AuthState`.
* Updated `login` to return `OperationalAuthState`.
* Kept repository login semantics restricted to operational sessions.

11. **mobile/lib/data/repositories/auth/auth_repository_impl.dart**

* Replaced `NotLoggedUser` with `AnonymousAuthState`.
* Replaced `LoggedUser` with `OperationalAuthState`.
* Added handling for `RestrictedInstallationAuthState` by returning `installationRegistrationRequired` without persisting tokens.
* Added fallback handling for unknown login result types.
* Preserved token persistence and profile loading only for operational authentication.

12. **mobile/lib/data/services/apis/auth/auth_api.dart**

* Updated login parsing to return `AuthState`.
* Replaced direct `LoggedUser` parsing with `AuthState.fromLoginMap`.
* Added login-specific backend error mapping for `INSTALLATION_LIMIT_REACHED`.
* Preserved API envelope error details in login failures.

13. **mobile/lib/data/services/apis/installation/dtos/installation_registration_response_dto.dart**

* Added DTO for installation registration responses.
* Mapped operational tokens, installation resource id, and installation status from backend response data.

14. **mobile/lib/data/services/apis/installation/installation_api.dart**

* Added `InstallationApi`.
* Implemented `POST /security/installations`.
* Sent restricted access token through `Authorization`.
* Sent `X-Installation-Id` and `X-Step-Up-Token` headers.
* Added response envelope parsing, malformed response handling, and backend error propagation.

15. **mobile/lib/data/services/services.dart**

* Registered `InstallationApi` in the data service composition.
* Wired the API with `RestClient` and `InstallationIdentityService`.

16. **mobile/lib/ui/pages/auth/viewmodel/login_viewmodel.dart**

* Updated login command typing from `LoggedUser` to `OperationalAuthState`.
* Kept the login flow dependent on installation identity resolution before repository login.

17. **mobile/test/core/services/client_http/dio/dio_error_mapper_test.dart**

* Added coverage for mapping `INSTALLATION_LIMIT_REACHED`.
* Verified app error code, message, HTTP status, and backend error code preservation.

18. **mobile/test/core/services/client_http/interceptors/auth/auth_interceptor_test.dart**

* Updated interceptor tests to provide a fake installation identity service.
* Verified that refresh requests include `X-Installation-Id`.
* Added coverage to ensure refresh is not called when installation identity resolution fails.
* Verified token cleanup behavior after failed refresh preparation.

19. **mobile/test/core/services/client_http/interceptors/installation/installation_interceptor_test.dart**

* Updated test setup to satisfy the new `AuthInterceptor` dependency on `InstallationIdentityService`.

20. **mobile/test/data/repositories/auth/auth_repository_impl_test.dart**

* Updated repository tests to use `AuthState`, `OperationalAuthState`, and `AnonymousAuthState`.
* Added coverage for restricted login responses.
* Verified that restricted login does not persist tokens, load profile data, update session state, or mark the repository as logged in.
* Added coverage for installation limit failures without token persistence.

21. **mobile/test/data/services/apis/auth/auth_api_login_test.dart**

* Added login API tests for operational authentication responses.
* Added login API tests for restricted installation authentication responses.
* Added coverage for `INSTALLATION_LIMIT_REACHED`.
* Added coverage for unknown successful response shapes returning parsing errors.

22. **mobile/test/data/services/apis/installation/installation_api_test.dart**

* Added installation registration API tests.
* Verified required headers for restricted token, installation identity, and step-up token.
* Verified that no transaction password is sent during installation registration.
* Added coverage for installation identity resolution failures.
* Added coverage for backend registration errors and malformed response data.

23. **mobile/test/ui/pages/auth/post_login_destination_test.dart**

* Updated fake auth repository types to use `AuthState`, `AnonymousAuthState`, and `OperationalAuthState`.

24. **mobile/test/ui/pages/auth/viewmodel/login_viewmodel_test.dart**

* Updated login view model tests to use `OperationalAuthState`.
* Replaced old logged user model references with the new auth state model.

25. **mobile/test/ui/pages/splash/viewmodel/splash_viewmodel_test.dart**

* Updated fake auth repository types to use `AuthState` and `AnonymousAuthState`.
* Updated login contract typing to return `OperationalAuthState`.

### Conclusion

This change establishes a typed authentication model capable of representing operational sessions, anonymous state, and restricted installation authorization without mixing them in the same runtime behavior.

The mobile app now recognizes installation registration requirements, avoids persisting restricted login tokens, includes installation identity during refresh, and provides a dedicated API client for registering installations.

The accompanying tests validate the new login parsing, error mapping, refresh header behavior, restricted login handling, and installation registration contract.


## 2026/06/18 - mobile/installation-identity-03

This change introduces centralized HTTP header definitions and integrates installation identity propagation into the mobile HTTP client pipeline.

The update replaces scattered literal header names with `AppHttpHeaders`, adds automatic `X-Installation-Id` injection through a dedicated interceptor, improves interceptor error mapping, and updates API services, documentation, dependency injection, and tests to use the new contract.

1. **mobile/lib/core/resources/app_http_headers.dart**

   * Added `AppHttpHeaders` as the centralized source for HTTP header names used by the mobile app.
   * Added constants for common, authentication, installation, step-up, trace, and app-token headers.
   * Added `sensitiveLowercase` to define headers that must be redacted from logs.
   * Added `bearer()` helper to format bearer authorization values consistently.

2. **mobile/lib/core/services/client_http/dio**

   * Updated `DioFactory` to use `AppHttpHeaders` for default `Accept` and `Content-Type` headers.
   * Updated `DioRestClient` to use the centralized sensitive header list when sanitizing request logs.
   * Updated `DioErrorMapper` to preserve `AppError` instances propagated by interceptors through `DioException.error`.

3. **mobile/lib/core/services/client_http/interceptors/auth**

   * Refactored `AuthInterceptor` to use `AppHttpHeaders.authorization` and `AppHttpHeaders.bearer()`.
   * Preserved the existing behavior of skipping requests that already define an authorization header.
   * Updated retry logic after token refresh to write the refreshed authorization header through the centralized header helper.

4. **mobile/lib/core/services/client_http/interceptors/installation**

   * Added `InstallationInterceptor`.
   * Resolved the installation identity before outgoing requests.
   * Added `X-Installation-Id` to request headers when identity resolution succeeds.
   * Blocked the request with a `DioException` carrying the original `AppError` when installation identity resolution fails.

5. **mobile/lib/core/services/client_http/interceptors/device**

   * Removed the old commented `DeviceInterceptor` implementation.
   * Updated interceptor exports to expose the new installation interceptor.

6. **mobile/lib/core/services/core_services.dart**

   * Registered `InstallationInterceptor` in dependency injection.
   * Updated the main `RestClient` setup to attach the installation interceptor before the auth interceptor.
   * Clarified comments around the main Dio instance and interceptor registration order.

7. **mobile/lib/data/services/apis**

   * Updated auth, contact verification, registration, and transfer APIs to use `AppHttpHeaders`.
   * Replaced literal `X-App-Token` and `X-Step-Up-Token` strings with centralized constants.
   * Updated API documentation examples to reference `AppHttpHeaders.appToken`.

8. **mobile/lib/core/AGENT.md and mobile/lib/data/services/apis/AGENT.md**

   * Updated documentation to describe authorization and app-token headers through `AppHttpHeaders`.
   * Aligned interceptor and API examples with the centralized header contract.

9. **mobile/test/core/resources/app_http_headers_test.dart**

   * Added tests covering exposed header names.
   * Added coverage for bearer token formatting.
   * Added coverage for sensitive header classification used by log redaction.

10. **mobile/test/core/services/client_http**

* Updated client request, response, Dio client, and auth interceptor tests to use `AppHttpHeaders`.
* Adjusted log redaction assertions to validate sensitive headers through centralized constants.
* Updated one response copy test expectation from `Accepted` to `Created`.

11. **mobile/test/core/services/client_http/interceptors/installation**

* Added tests for installation header injection.
* Added coverage for interaction between installation and auth interceptors.
* Added failure-path coverage to ensure unresolved installation identity blocks outgoing requests.
* Added RestClient-level coverage to ensure interceptor failures are returned as `Failure` without reaching the HTTP adapter.

12. **mobile/test/data/services/apis**

* Updated contact verification and transfer API tests to assert headers through `AppHttpHeaders`.
* Adjusted test descriptions to reference app-token behavior instead of literal header strings.

13. **mobile/test/data/repositories/transaction_password**

* Updated step-up authorization failure message expectation to `Step-up auth failed`.

### Conclusion

The mobile HTTP layer now has a single source of truth for header names, sensitive header redaction, and bearer token formatting.

Installation identity is now part of the standard request pipeline, and requests are blocked when the identity cannot be resolved, preserving the underlying application error for consistent handling.

The API layer, dependency injection setup, documentation, and tests were aligned with this behavior, reducing duplicated header literals and strengthening request consistency.


## 2026/06/18 - mobile/installation-identity-02

This change introduces the first mobile implementation layer for installation identity in BankLab. It adds durable installation ID storage, a non-secret local marker, identity resolution during app bootstrap, and login blocking when the installation identity cannot be resolved.

The update affects core services, dependency injection, splash initialization, login flow, storage contracts, documentation, and automated tests for the installation identity lifecycle.

1. **docs/backlogs/mobile/013 - installation-identity-mvp_tasks.md**

   * Translated the installation identity MVP task document to Portuguese.
   * Preserved the original task structure, scope, dependencies, acceptance criteria, and deferred items.
   * Clarified the mobile implementation plan for local identity, HTTP propagation, restricted login, installation registration, and future management flows.

2. **mobile/lib/core/resources/storage_keys.dart**

   * Added `StorageKeys.installationId` for the durable installation identity stored in secure storage.
   * Added `StorageKeys.installationLocalMarker` as a non-secret marker stored outside secure storage.
   * Documented that logout, credential cleanup, and user switching must not delete the installation identity.

3. **mobile/lib/core/services/installation_identity/**

   * Added the installation identity module exports.
   * Created `InstallationMarkerStore` as the abstraction for checking and writing the local installation marker.
   * Implemented `FileInstallationMarkerStore` using the application support directory.
   * Implemented `InstallationIdentityService` to resolve, validate, generate, persist, and refresh the local installation identity.
   * Added canonical UUID v4 validation and replacement behavior for missing markers, missing values, or invalid stored values.
   * Added failure handling for marker and secure storage errors without silently continuing.

4. **mobile/lib/core/services/core_services.dart**

   * Registered `InstallationMarkerStore` and `InstallationIdentityService` in core dependency injection.
   * Inserted installation identity setup before the main HTTP client and authentication interceptor registration.
   * Preserved the existing Dio, AuthInterceptor, and RestClient setup sequence.

5. **mobile/lib/ui/pages/splash/viewmodel/splash_viewmodel.dart**

   * Updated splash initialization to resolve installation identity before loading remembered login data.
   * Added `SplashBootstrapState` to carry the optional remembered login identity.
   * Changed missing last-login data into a successful bootstrap state instead of a blocking failure.
   * Preserved blocking behavior when installation identity resolution fails.

6. **mobile/lib/ui/pages/splash/models/splash_bootstrap_state.dart**

   * Added a bootstrap state model to represent splash initialization output.
   * Encapsulated the optional `LastLoginIdentity` used to decide between short login and full login navigation.

7. **mobile/lib/ui/pages/splash/splash_page.dart**

   * Updated the splash animation listener to react to initialization command state.
   * Added a recoverable error UI when installation identity bootstrap fails.
   * Added a retry action that re-executes initialization.
   * Adjusted navigation to route to short login only when a remembered identity is available.

8. **mobile/lib/ui/pages/auth/viewmodel/login_viewmodel.dart**

   * Injected `InstallationIdentityService` into `LoginViewModel`.
   * Wrapped login execution with installation identity resolution.
   * Prevented login API calls when installation identity resolution fails.
   * Preserved post-login destination resolution based on the current app session.

9. **mobile/test/core/resources/storage_keys_test.dart**

   * Added tests for the durable installation ID storage key.
   * Added tests ensuring the local marker key is separate from the installation ID key.

10. **mobile/test/core/services/installation_identity/installation_identity_service_test.dart**

* Added lifecycle tests for first-run identity creation.
* Added tests for UUID reuse when the local marker exists.
* Added tests for replacing secure-storage values when the marker is missing.
* Added tests for invalid stored identity replacement.
* Added failure-path tests for marker lookup, secure storage read, secure storage write, and marker refresh errors.

11. **mobile/test/data/repositories/auth/auth_repository_impl_test.dart**

* Updated logout test coverage to confirm installation identity is not deleted during session cleanup.
* Added tracking of deleted storage keys in the fake secure storage implementation.
* Verified logout still clears access and refresh tokens while preserving the installation ID.

12. **mobile/test/ui/pages/auth/post_login_destination_test.dart**

* Updated `LoginViewModel` test setup to provide a fake installation identity service.
* Added local fake marker and secure storage implementations required by the new constructor dependency.
* Preserved existing post-login destination behavior tests.

13. **mobile/test/ui/pages/auth/viewmodel/login_viewmodel_test.dart**

* Added tests confirming login is blocked when installation identity resolution fails.
* Added tests confirming the login repository is called only after installation identity resolves successfully.

14. **mobile/test/ui/pages/splash/viewmodel/splash_viewmodel_test.dart**

* Added tests confirming splash bootstrap blocks when installation identity resolution fails.
* Added tests confirming remembered login identity is returned after successful identity resolution.
* Added tests confirming bootstrap continues successfully when no remembered identity exists.

### Conclusion

This change establishes the mobile foundation for installation identity by introducing durable identity storage, a local marker contract, identity resolution, bootstrap blocking, and login protection.

The mobile app now avoids entering authentication flows when it cannot safely resolve the current installation identity, while preserving existing logout behavior and short-login navigation.

Test coverage was expanded across storage keys, identity lifecycle, splash initialization, login blocking, and logout cleanup behavior.


## 2026/06/17 - mobile/installation-identity-01

This change finalizes the mobile planning backlog for the Installation Identity MVP and converts the previous research-oriented document into an implementation-ready plan.

The update closes the pending mobile decisions around local storage, reinstall detection, bootstrap failure behavior, restricted login, installation registration, installation limits, and deferred management capabilities. It also adds a dedicated task breakdown to guide the first mobile implementation cut.

1. **docs/backlogs/mobile/013 - installation-identity-mvp.md**

   * Changed the backlog type from research to planning and marked the state as ready for tasks.
   * Updated the mobile objective to block login and API calls when a stable installation identity cannot be resolved.
   * Expanded the installation lifecycle rules with a local marker contract to distinguish app updates from reinstall, clear-data, or incomplete restore scenarios.
   * Defined `banklab.installation.id` as the installation identity storage key.
   * Clarified that logout, credential cleanup, and user switching must not delete the installation identity.
   * Documented that secure storage alone must not be treated as proof that the current installation is still the same.
   * Reworked failure handling so bootstrap, login, and API calls are blocked when identity resolution fails.
   * Defined the expected retry behavior and prohibited sending requests without `X-Installation-Id`.
   * Refined the restricted installation registration flow, including token cleanup, expiration handling, cancellation handling, and transaction password preconditions.
   * Added the approved UX message for `installation_limit_reached`.
   * Deferred installation listing and revocation management to a future mobile implementation.
   * Updated the decision checklist to mark storage, reinstall detection, backup/restore behavior, restricted token handling, step-up integration, limit UX, telemetry, and bootstrap tests as resolved.
   * Adjusted the out-of-scope section to exclude mobile listing and revocation instead of API association concerns.
   * Updated references to point to the completed API backlog and the new task document.

2. **docs/backlogs/mobile/013 - installation-identity-mvp_tasks.md**

   * Added a new implementation task breakdown for the Installation Identity MVP.
   * Defined storage key and local marker tasks for distinguishing app updates from reinstall or restore cases.
   * Added the `InstallationIdentityService` task for read-or-create UUID v4 resolution, validation, persistence, marker refresh, and safe logging.
   * Added bootstrap blocking requirements for unresolved or failed installation identity resolution.
   * Added the installation interceptor task to propagate `X-Installation-Id` through the main HTTP client.
   * Added refresh-flow requirements to ensure `/auth/refresh` also sends the installation identity.
   * Added login response modeling tasks for operational login, restricted installation registration, and installation limit outcomes.
   * Added the installation registration API task for `POST /security/installations`.
   * Added the step-up operation task for authorizing installation registration with the transaction password flow.
   * Added the restricted installation certification flow task, covering password entry, step-up authorization, registration, token persistence, and cleanup.
   * Added blocker handling tasks for installation limits and transaction password preconditions.
   * Added safe telemetry and test coverage tasks for lifecycle, headers, restricted login, registration, and blocking outcomes.
   * Added the final mobile verification task covering formatting, focused tests, analysis, and regression checks.
   * Documented deferred work for installation listing, revocation, management UI, and current-installation removal protection.

### Conclusion

The mobile installation identity backlog is now implementation-ready and aligned with the API contract.

The change establishes the expected local identity lifecycle, failure behavior, restricted registration flow, and user-facing blocker states, while explicitly deferring installation management to a later cut.

The new task file provides a structured execution path for implementing the MVP without mixing it with future list and revocation features.


## 2026/06/17 - api/installation-identity-09

This change completes the audit, retention, and operational cleanup layer for the Installation Identity MVP. It documents the security boundaries of installation identity, clarifies revoked-installation behavior, and adds automated cleanup for restricted installation registration authorizations.

The update affects REST documentation, authentication implementation notes, database documentation, PostgreSQL migrations, the concrete installation repository, integration tests, and backlog organization.

1. **api/docs/07-api-rest.md**

   * Clarified that `X-Installation-Id` is only a weak installation signal.
   * Documented that installation identity does not replace JWT validation or step-up authorization.
   * Expanded revocation behavior for installations.
   * Documented that revoked installations invalidate bound refresh sessions, while already issued access tokens remain valid until their short expiration.
   * Clarified that revoked installations remain in historical storage and do not consume known-installation slots.

2. **api/docs/08-auth_implementation.md**

   * Updated the token model to include restricted access tokens.
   * Documented the restricted access token used during new installation registration.
   * Described its scope, token type, expiration, persisted `jti`, and operational limitations.
   * Added refresh-session behavior when sessions are bound to installations.
   * Added protected-route requirements for matching `X-Installation-Id`.
   * Added the installation registration flow using restricted access tokens and step-up authorization.
   * Added audit, retention, and safe logging rules for installation identity events.
   * Documented the cleanup policy for restricted installation registration authorizations.
   * Renumbered authorization and design rationale sections after the new operational documentation.

3. **api/docs/09-database.md**

   * Added `app_installations` and `installation_registration_authorizations` to the support table list.
   * Documented `installation_id` on refresh sessions.
   * Added notes describing refresh-session invalidation when an installation is revoked.
   * Added the `app_installations` table documentation, including status, slot behavior, metadata, retention, and privacy limits.
   * Added the `installation_registration_authorizations` table documentation, including states, scope, expiration, retention, and cleanup behavior.

4. **api/internal/installation/infrastructure/postgres_restricted_authorization_repository.go**

   * Added `CleanupExpired` to remove restricted installation registration authorizations outside the retention window.
   * Implemented cleanup rules for expired active authorizations, old consumed authorizations, and old revoked authorizations.
   * Added validation for invalid cleanup inputs.
   * Returned the number of deleted rows for operational visibility and testability.

5. **api/internal/installation/infrastructure/postgres_repository_test.go**

   * Added an integration test for restricted authorization cleanup.
   * Covered removal of old active, consumed, and revoked authorizations.
   * Verified that recent active authorizations are preserved.
   * Added a helper to create restricted authorizations with controlled timestamps.

6. **api/migrations/000015_installation_authorizations_cleanup.up.sql**

   * Added the `pg_cron` extension requirement.
   * Added partial indexes for active, consumed, and revoked restricted authorization cleanup.
   * Added `cleanup_installation_registration_authorizations()`.
   * Scheduled the cleanup job to run daily at 03:30.
   * Executed the cleanup function once during migration application.

7. **api/migrations/000015_installation_authorizations_cleanup.down.sql**

   * Added rollback logic for the scheduled cleanup job.
   * Removed the cleanup function.
   * Dropped the partial indexes created for cleanup support.

8. **docs/backlogs/api/done/**

   * Moved installation identity backlogs 010 through 019 to the `done` folder.
   * Marked the audit and retention backlog as concluded.
   * Added implemented decisions covering revoked installation retention, restricted authorization cleanup, refresh-session invalidation, safe logging, and weak installation identity semantics.
   * Added completed task documentation for the audit, retention, cleanup, and operational documentation work.

9. **docs/backlogs/mobile/pendencias.md**

   * Added mobile pending decisions for installation identity integration.
   * Documented open UX and implementation questions for reinstall detection, backup and restore, storage failure, concurrency, restricted login, installation limit handling, transaction password edge cases, token expiration, installation management, telemetry, and first delivery scope.

### Conclusion

This change finalizes the operational retention policy for installation identity and restricted installation registration authorizations.

The API documentation now presents installation identity as a weak contextual signal, with explicit limits around authentication, authorization, revocation, logging, and audit behavior.

The database and repository layers now include a versioned cleanup mechanism, backed by migration, scheduled execution, and integration test coverage.


## 2026/06/17 - api/installation-identity-08

This change completes the API delivery and enforcement layer for installation identity management. It connects restricted installation registration, operational installation enforcement, installation listing, and revocation into the real API wiring.

The implementation introduces dedicated installation use cases, HTTP handlers, middleware paths for operational and restricted tokens, refresh-token installation validation, and REST contract updates. It also extends step-up authorization to support installation registration under restricted authentication.

1. **api/cmd/api/main.go**

   * Wired the installation application and delivery modules into the API bootstrap.
   * Added use cases for registering, listing, and revoking installations.
   * Registered the installation handler with the main API router.
   * Replaced generic authenticated routing with operational, restricted, and operational-or-restricted middleware flows.
   * Added routes for:

     * `POST /security/installations`
     * `GET /security/installations`
     * `DELETE /security/installations/{installation_resource_id}`
   * Configured the JWT middleware with restricted access token verification.

2. **api/cmd/api/routes_test.go**

   * Updated router tests to account for the new middleware signatures.
   * Verified that step-up authorization uses the operational-or-restricted authentication path.
   * Added route-level tests confirming that installation registration uses restricted authentication, while listing and revocation use operational authentication.

3. **api/internal/auth/delivery/jwt_middleware.go**

   * Added operational authentication enforcement requiring `X-Installation-Id`.
   * Added restricted authentication support using restricted access token verification.
   * Added combined operational-or-restricted authentication for flows that may accept either context.
   * Validated consistency between the installation header and token claims.
   * Populated `OperationalSession` and `RestrictedSession` contexts for downstream application use.
   * Extracted principal creation into a shared helper.

4. **api/internal/auth/delivery/handler.go**

   * Updated refresh-token handling to require and parse `X-Installation-Id`.
   * Propagated the validated installation identifier to the refresh access token use case.

5. **api/internal/auth/application/refresh_access_token.go**

   * Extended refresh input with `InstallationID`.
   * Added validation to ensure the refresh session belongs to the installation presented in the request.
   * Returned installation mismatch when the header and persisted session installation differ.

6. **api/internal/auth/delivery/handler_test.go**

   * Updated refresh tests to include `X-Installation-Id`.
   * Verified that the installation identifier is propagated to the refresh use case.
   * Preserved invalid-token behavior under the new refresh contract.

7. **api/internal/auth/application/errors_registry.go**

   * Registered installation mismatch as a public domain error.
   * Registered additional restricted authorization states as invalid-token responses.

8. **api/internal/bootstrap/errors.go**

   * Added installation application error registration to the API bootstrap.
   * Preserved the existing registration order for admin, customer, auth, and security errors.

9. **api/internal/installation/application/errors_registry.go**

   * Added installation-specific error registration.
   * Mapped installation mismatch, invalid resource identifiers, missing installations, and restricted authorization failures to stable HTTP error responses.

10. **api/internal/installation/application/register_installation.go**

* Added the installation registration use case.
* Validated restricted session scope, user, installation identifier, and authorization `jti`.
* Enforced step-up for `installation.register`.
* Reserved a known installation slot using the installation repository.
* Consumed restricted authorization inside the registration transaction.
* Issued operational access and refresh tokens after successful registration.
* Persisted the refresh session bound to the installation.

11. **api/internal/installation/application/register_installation_test.go**

* Added tests for successful installation registration.
* Verified token generation, refresh-token hashing, session persistence, step-up enforcement, authorization consumption, and installation reservation.
* Covered mismatch scenarios where registration must not consume restricted authorization or step-up.
* Covered validation of consumed restricted authorization against the current restricted context.

12. **api/internal/installation/application/list_installations.go**

* Added the use case for listing installations owned by the authenticated user.
* Returned safe installation summaries using public `resource_id`, status, and timestamps.
* Avoided exposing the raw client-generated `installation_id`.

13. **api/internal/installation/application/list_installations_test.go**

* Added tests confirming that known and revoked installations are returned as safe summaries.
* Verified that listing is scoped to the operational session user.

14. **api/internal/installation/application/revoke_installation.go**

* Added the installation revocation use case.
* Validated operational session context.
* Revoked installations by public resource identifier.
* Prevented revocation of the current installation.
* Invalidated refresh sessions associated with the revoked installation inside the transaction.

15. **api/internal/installation/application/revoke_installation_test.go**

* Added tests for successful revocation and session invalidation.
* Verified that the current installation cannot revoke itself.
* Ensured invalid revocation attempts do not revoke installations or invalidate sessions.

16. **api/internal/installation/delivery/handler.go**

* Added HTTP handlers for installation registration, listing, and revocation.
* Parsed and validated canonical UUID values from headers and path parameters.
* Exposed safe response DTOs for installation management.
* Ensured list and revoke responses do not expose raw installation identifiers.

17. **api/internal/installation/delivery/handler_test.go**

* Added handler tests for successful registration.
* Verified step-up token and installation header propagation.
* Confirmed that listing does not expose `installation_id`.
* Covered invalid installation resource identifiers during revocation.

18. **api/internal/security/delivery/handler.go**

* Updated step-up authorization to support restricted sessions for installation registration.
* Preserved operational authentication requirements for other sensitive operations.
* Rejected restricted-context step-up requests targeting operations other than `POST /security/installations`.

19. **api/internal/security/domain/step_up_endpoint_policy.go**

* Added `installation.register` as an allowed step-up endpoint key.
* Included the installation registration endpoint in the default step-up policy.

20. **api/internal/security/domain/step_up_endpoint_policy_test.go**

* Added coverage confirming that the default step-up policy allows installation registration.

21. **api/internal/security/domain/step_up_public_operation_resolver.go**

* Added the public operation mapping for `POST /security/installations`.
* Resolved that route to the `installation.register` endpoint key.

22. **api/internal/security/domain/step_up_public_operation_resolver_test.go**

* Added coverage confirming that `POST /security/installations` resolves to the installation registration step-up endpoint.

23. **api/internal/shared/errors/codes.go**

* Added the public error code `INSTALLATION_MISMATCH`.

24. **api/docs/07-api-rest.md**

* Updated authentication rules to require `X-Installation-Id` for refresh and operational authenticated routes.
* Documented restricted login responses for new installations.
* Documented installation registration, listing, and revocation endpoints.
* Added installation-related error codes and behavior.
* Updated step-up documentation to include installation registration as a supported public operation.

25. **docs/backlogs/api/017 - installation-identity-management-usecases_tasks.md**

* Added completed task documentation for installation identity management use cases.
* Documented registration, restricted context validation, installation matching, slot reservation, authorization consumption, session creation, listing, revocation, and session invalidation.

26. **docs/backlogs/api/018 - installation-identity-delivery-enforcement_tasks.md**

* Added completed task documentation for delivery and enforcement of installation identity.
* Documented refresh enforcement, operational middleware, restricted middleware, step-up support for restricted registration, installation handlers, router wiring, and REST contract validation.

### Conclusion

This change set completes the installation identity management flow across application, delivery, middleware, routing, error handling, tests, and documentation.

The API now distinguishes operational access from restricted installation-registration access, binds authenticated flows to `X-Installation-Id`, supports secure registration of new installations through step-up, and allows users to list and revoke known installations without exposing raw installation identifiers.


## 2026/06/17 - api/installation-identity-07

This change prepares the authentication flow to carry installation identity across operational sessions, access tokens, refresh rotation, and shared authentication context.

It adds persistence support for `installation_id` in user sessions, extends JWT claims, introduces restricted access token infrastructure for installation registration, and documents the related backlog tasks.

1. **api/internal/auth/domain/interfaces.go**

   * Extended `TokenClaims` with optional `InstallationID`.
   * Added session creation and lookup models with installation metadata.
   * Expanded `SessionRepository` with installation-aware create, lookup, and revocation methods.

2. **api/internal/auth/application/login_user.go**

   * Included `installation_id` in generated access tokens when available.
   * Persisted refresh sessions using the installation-aware session creation flow.
   * Added helper logic to keep legacy logins compatible when no installation ID is provided.

3. **api/internal/auth/application/login_user_test.go**

   * Updated session repository mocks for installation-aware session creation.
   * Added assertions ensuring legacy login keeps `installation_id` absent.
   * Added assertions ensuring classified installation login propagates the installation ID into both token claims and session persistence.

4. **api/internal/auth/application/refresh_access_token.go**

   * Replaced legacy session lookup with installation-aware session records.
   * Preserved `installation_id` during access token generation and refresh token rotation.
   * Created rotated refresh sessions with the same installation binding.

5. **api/internal/auth/application/refresh_access_token_test.go**

   * Updated refresh session mocks and stateful test repositories to support installation-aware records.
   * Added coverage for preserving `installation_id` across refresh rotation.
   * Added revocation support by user and installation in test doubles.

6. **api/internal/auth/infrastructure/jwt_token_service.go**

   * Added optional `installation_id` claim emission in access tokens.
   * Added parsing and UUID validation for `installation_id`.
   * Preserved compatibility with access tokens that do not contain installation metadata.

7. **api/internal/auth/infrastructure/jwt_token_service_test.go**

   * Added coverage for generating and parsing access tokens with `installation_id`.
   * Added coverage for access tokens without installation identity.
   * Added validation coverage for malformed `installation_id` claims.

8. **api/internal/auth/infrastructure/postgres_session_repository.go**

   * Implemented installation-aware session creation and lookup.
   * Added session revocation by `user_id` and `installation_id`.
   * Implemented `InstallationSessionInvalidator` through the session repository.

9. **api/internal/auth/infrastructure/postgres_session_repository_test.go**

   * Added integration coverage for creating, reading, and revoking sessions by installation.
   * Verified that revocation affects only the target installation session.
   * Verified that sessions from other installations remain active.

10. **api/internal/auth/infrastructure/postgres_user_repository_test.go**

* Extended the shared auth repository test schema with `user_sessions`.
* Added `installation_id` schema support and index setup for repository tests.
* Adjusted test users to use an admin role where customer dependencies are not required.

11. **api/internal/auth/delivery/auth_authorization_integration_test.go**

* Updated integration schema setup to include the optional `installation_id` column in `user_sessions`.

12. **api/internal/installation/domain/restricted_token.go**

* Added restricted access token claims for installation registration.
* Defined validation rules for token type, scope, user, installation, JTI, and expiration.
* Added signer and verifier interfaces for restricted access tokens.

13. **api/internal/installation/infrastructure/jwt_restricted_access_token_service.go**

* Added JWT-based restricted access token signer and verifier.
* Validated restricted token claims against persisted restricted authorization records.
* Added rejection paths for missing, consumed, revoked, expired, mismatched, or invalid authorizations.

14. **api/internal/installation/infrastructure/jwt_restricted_access_token_service_test.go**

* Added coverage for signing and verifying restricted access tokens.
* Added validation for authorization mismatch.
* Added validation for consumed authorization rejection.

15. **api/internal/shared/authctx/context.go**

* Added operational session context carrying user, role, customer, and installation identity.
* Added restricted session context carrying user, installation, JTI, and scope.
* Added required and optional accessors with explicit missing-context errors.

16. **api/internal/shared/authctx/context_test.go**

* Added tests for operational session context storage and retrieval.
* Added tests for restricted session context storage and retrieval.
* Added missing-context error coverage for both session types.

17. **api/migrations/000014_user_sessions_installation_id.up.sql**

* Added `installation_id` to `user_sessions`.
* Added partial index on `(user_id, installation_id)` for installation-bound sessions.

18. **api/migrations/000014_user_sessions_installation_id.down.sql**

* Added rollback for the installation session index.
* Added rollback for the `installation_id` column.

19. **docs/backlogs/api/015 - installation-identity-session-tokens-context.md**

* Linked the backlog to its dedicated task breakdown document.

20. **docs/backlogs/api/015 - installation-identity-session-tokens-context_tasks.md**

* Added task breakdown for session schema, repository contracts, access token claims, restricted tokens, authentication contexts, installation-based invalidation, and validation without enforcement.

### Conclusion

This change establishes the infrastructure required to bind authentication sessions and tokens to installation identity while preserving compatibility with legacy flows.

The API can now persist installation-aware refresh sessions, propagate installation identity through access token and refresh rotation flows, validate restricted installation registration tokens, and expose operational and restricted authentication contexts for future enforcement work.


## 2026/06/17 - api/installation-identity-05

This change implements the PostgreSQL persistence layer for the installation identity module, adding schema support, repository implementations, and integration coverage for installation records and restricted registration authorizations.

The update introduces the database structures required to persist known and revoked app installations, enforce the three-known-installation limit per user, and support future restricted authorization flows for installation registration.

1. **api/migrations/000013_installation_identity_schema.up.sql**

   * Added the `app_installations` table to persist installation identity records.
   * Added integrity constraints for installation status, timestamps, revocation consistency, and known slot consistency.
   * Added unique indexes for `resource_id`, `(user_id, installation_id)`, and active known slots.
   * Added query indexes for user-based lookup, resource lookup, status filtering, and known installation counting.
   * Added the `installation_registration_authorizations` table for restricted installation registration authorizations.
   * Added constraints for `jti`, scope, status, expiration, and consumed timestamp consistency.
   * Added unique and query indexes for authorization lookup, active authorization enforcement, and expiration filtering.

2. **api/migrations/000013_installation_identity_schema.down.sql**

   * Added rollback support for the installation identity schema.
   * Drops authorization indexes and table before removing installation indexes and table.
   * Keeps rollback order compatible with table dependencies.

3. **api/internal/installation/infrastructure/postgres_repository.go**

   * Added `PostgresInstallationRepository` implementing `domain.InstallationRepository`.
   * Added transaction-aware executor support using `database.TxFromContext`.
   * Implemented lookup by `(user_id, installation_id)` and `(user_id, resource_id)`.
   * Implemented known installation counting, historical installation detection, and user installation listing.
   * Implemented atomic first-installation bootstrap with user-level locking.
   * Implemented atomic known-installation reservation with slot allocation and limit enforcement.
   * Implemented logical revocation by public `resource_id`, preserving history and releasing the known slot.
   * Added scan helpers to restore domain installation entities from PostgreSQL rows.
   * Added PostgreSQL constraint and uniqueness error mapping to domain errors.

4. **api/internal/installation/infrastructure/postgres_restricted_authorization_repository.go**

   * Added `PostgresRestrictedAuthorizationRepository` implementing `domain.RestrictedAuthorizationRepository`.
   * Implemented creation, lookup by `jti`, atomic consumption, and revocation of restricted authorizations.
   * Added active authorization revocation by `(user_id, installation_id, scope)`.
   * Added single-use consumption behavior based on `status = active` and `expires_at`.
   * Added failure mapping for consumed, revoked, expired, invalid, and missing authorizations.
   * Added scan helpers to restore restricted authorization domain entities from database rows.
   * Added PostgreSQL error mapping for duplicate active authorizations and invalid persisted data.

5. **api/internal/installation/infrastructure/postgres_repository_test.go**

   * Added PostgreSQL integration tests for installation repository behavior.
   * Covered first-installation bootstrap, historical bootstrap rejection, known installation limit enforcement, revoked slot reuse, and historical listing.
   * Added concurrency tests ensuring only one bootstrap succeeds for a new user.
   * Added concurrency tests ensuring known installation reservations do not exceed the configured limit.
   * Added PostgreSQL integration tests for restricted authorization repository behavior.
   * Covered authorization creation, lookup, active uniqueness, single-use consumption, expiration handling, and concurrent consumption.
   * Added local test schema helpers, test user creation, cleanup helpers, and domain value constructors.

6. **docs/backlogs/api/013 - installation-identity-database-schema_tasks.md**

   * Marked all database schema tasks as completed.
   * Updated task status for installation migrations, constraints, indexes, known slot support, restricted authorization schema, and rollback validation.

7. **docs/backlogs/api/014 - installation-identity-repositories.md**

   * Added a reference from the repository backlog to its dedicated task document.
   * Linked the repository planning document to the new task breakdown.

8. **docs/backlogs/api/014 - installation-identity-repositories_tasks.md**

   * Added the task breakdown for the installation identity repository backlog.
   * Documented eight completed tasks covering infrastructure structure, installation reads, atomic bootstrap, atomic reservation, logical revocation, restricted authorization persistence, concurrency tests, and build validation.
   * Captured scope, acceptance criteria, and dependencies for each repository implementation step.

### Conclusion

This change delivers the persistence foundation for installation identity in the API. The database schema, repository implementations, and integration tests now support installation history, known installation limits, logical revocation, and restricted registration authorization storage.

The module remains isolated from handlers and login wiring, keeping the implementation ready for later application and delivery integration without changing current HTTP behavior.


## 2026/06/17 - api/installation-identity-04

This change introduces the domain contracts for Installation Identity in the API. It defines the core installation aggregate, installation identifiers, login classifications, restricted registration authorizations, domain errors, repository ports, and unit tests.

The update also closes the domain-contract backlog tasks and prepares the next database-schema backlog with detailed migration tasks.

1. **api/internal/installation/domain/errors.go**

   * Added domain errors for installation identity validation, lookup, mismatch, revocation, limit enforcement, and first-installation bootstrap conflicts.
   * Added restricted authorization errors for invalid, missing, expired, consumed, revoked, and already-active authorization states.

2. **api/internal/installation/domain/installation_id.go**

   * Added `InstallationID` as a strict value object for client-provided installation identifiers.
   * Enforced canonical UUID v4 format for installation IDs.
   * Added `InstallationResourceID` as a separate value object for public installation management resources.
   * Added parsing, string conversion, UUID access, and zero-value checks for both identifiers.

3. **api/internal/installation/domain/installation.go**

   * Added the `Installation` domain entity with persisted fields for user association, resource identity, installation identity, status, platform metadata, lifecycle timestamps, and revocation state.
   * Added `known` and `revoked` installation statuses.
   * Added `NewKnownInstallation` and `RestoreInstallation` constructors with invariant validation.
   * Added validation rules for required identifiers, timestamps, status consistency, and revocation state.
   * Added `Revoke` behavior to transition a known installation to revoked.
   * Added login classifications for known, first, new, revoked, and limit-reached installation scenarios.
   * Added `LoginDecision` validation around known installation count, associated installation history, and the maximum known installation limit.

4. **api/internal/installation/domain/restricted_authorization.go**

   * Added `RestrictedAuthorization` for installation registration authorization.
   * Defined the initial restricted authorization scope as `installation.register`.
   * Added the default authorization duration of five minutes.
   * Added active, consumed, and revoked authorization statuses.
   * Added constructors and restoration logic with invariant validation.
   * Added expiration detection, consume behavior, and revoke behavior.
   * Enforced single-use behavior through consumed-state validation.

5. **api/internal/installation/domain/repository.go**

   * Added read ports for finding installations by user and installation ID, finding by resource ID, counting known installations, checking user installation history, and listing installations.
   * Added write ports for bootstrapping the first installation, reserving a known installation, and revoking an installation by resource ID.
   * Added a session invalidation port for invalidating sessions by installation ID.
   * Added restricted authorization repository ports for create, lookup, consume, revoke, and revoke-active operations.

6. **api/internal/installation/domain/installation_id_test.go**

   * Added tests for canonical UUID v4 installation ID parsing.
   * Covered invalid cases such as blank values, non-UUID strings, UUID v1, non-canonical UUIDs without hyphens, and uppercase UUIDs.
   * Verified that `InstallationResourceID` remains a separate value object from `InstallationID`.
   * Added validation for rejecting nil installation resource IDs.

7. **api/internal/installation/domain/installation_test.go**

   * Added tests for creating known installations.
   * Added revocation behavior tests, including repeated revocation rejection.
   * Added restoration validation tests for known and revoked installation state consistency.
   * Added login decision tests covering known, revoked, first, new, limit-reached, and invalid classification scenarios.

8. **api/internal/installation/domain/restricted_authorization_test.go**

   * Added tests for creating restricted authorizations with trimmed JTI values, default scope, active status, and default expiration.
   * Added consume behavior tests, including repeated consumption rejection.
   * Added expiration validation for restricted authorization consumption.
   * Added revoke behavior tests, including repeated revocation rejection.
   * Added restoration validation tests for active, consumed, revoked, invalid scope, blank JTI, and inconsistent consumed timestamp states.

9. **docs/backlogs/api/012 - installation-identity-domain-contracts_tasks.md**

   * Marked all eight domain-contract tasks as completed.
   * Updated the backlog status for installation identifiers, installation modeling, login classifications, domain errors, read ports, write ports, restricted authorization modeling, and restricted authorization ports.

10. **docs/backlogs/api/013 - installation-identity-database-schema.md**

* Added a direct reference to the new database schema task file.
* Linked the parent backlog to the detailed migration task breakdown.

11. **docs/backlogs/api/013 - installation-identity-database-schema_tasks.md**

* Added the database schema task backlog for Installation Identity.
* Defined tasks for creating the `app_installations` table, integrity constraints, query indexes, support for the three-known-installation limit, restricted authorization storage, authorization constraints/indexes, and migration validation.
* Documented acceptance criteria and dependencies for each database task.

### Conclusion

This change establishes the Installation Identity domain layer with explicit entities, value objects, lifecycle rules, authorization rules, errors, ports, and unit coverage.

It completes the domain-contract backlog and creates the next implementation path for database schema work, keeping persistence concerns separated from the domain model.


## 2026/06/17 - api/installation-identity-03

This change consolidates the installation identity entry contract and reorganizes the API backlog structure around technical dependency order. The login REST contract now explicitly documents the required `X-Installation-Id` header, including canonical UUID v4 validation and the `INVALID_INSTALLATION_ID` error behavior.

The backlog documentation was refactored to separate planning by architectural dependency instead of by isolated flow. This creates a clearer sequence from entry contract, domain contracts, database schema, repositories, session/token context, use cases, delivery enforcement, and operational documentation.

1. **api/docs/07-api-rest.md**

   * Documented that `POST /auth/login` requires both `X-App-Token` and `X-Installation-Id`.
   * Added request header documentation for `X-Installation-Id`.
   * Defined the expected canonical lowercase UUID v4 format with hyphens.
   * Clarified that the current contract only validates and propagates the installation identifier to the login application layer.
   * Added `INVALID_INSTALLATION_ID` as a possible login error.
   * Added the common error code description for missing or malformed installation identifiers.

2. **api/internal/auth/delivery/handler_test.go**

   * Expanded invalid installation identifier test coverage using table-driven cases.
   * Added validation scenarios for non-UUID values, UUID v1, UUID without hyphens, uppercase UUID, and blank values.
   * Preserved the assertion that invalid headers return `400 INVALID_INSTALLATION_ID`.
   * Preserved the assertion that the login use case is not called when the installation identifier is invalid.

3. **docs/backlogs/README.md**

   * Updated the installation identity backlog index to reflect the new dependency-based split.
   * Added the new split document as the operational reference for backlog sequencing.
   * Replaced the previous flow-oriented backlog names with the new dependency-oriented structure.
   * Updated backlog descriptions for domain contracts, database schema, repositories, session/token context, login use cases, management use cases, delivery enforcement, and audit/retention documentation.

4. **docs/backlogs/api/010 - installation-identity-mvp.md**

   * Removed the detailed implementation order from the main MVP backlog.
   * Replaced the embedded implementation sequence with a reference to the dedicated split document.
   * Simplified the derived backlog section to point to the dependency-based split.
   * Kept the main MVP backlog focused on product decision, threat model, and scope instead of operational task sequencing.

5. **docs/backlogs/api/010 - split-installation-identity-by-dependency.md**

   * Added a new planning document defining the dependency-based implementation order for installation identity.
   * Established the technical sequence from entry contract through audit, retention, and operational documentation.
   * Defined separation rules for domain, interfaces, migrations, repositories, use cases, delivery, and enforcement.
   * Documented the allowed exception for the minimal `X-Installation-Id` login entry contract before persistence exists.
   * Listed the resulting backlog files from 011 through 019.

6. **docs/backlogs/api/011 - installation-identity-entry-contract.md**

   * Refined the backlog objective to emphasize that the entry contract does not depend on persistent domain, installation tables, or linked sessions.
   * Clarified that this backlog is the initial layered-order exception because the HTTP header can be validated without consulting state.
   * Updated the task preparation guidance to avoid fake repositories or premature operational classification.
   * Added a link to the dedicated task file.
   * Added the split document as a reference.

7. **docs/backlogs/api/011 - installation-identity-entry-contract_tasks.md**

   * Added the task breakdown for the completed entry contract backlog.
   * Documented tasks for centralizing the installation header constant, defining `INVALID_INSTALLATION_ID`, validating canonical UUID v4, propagating the validated value to the application layer, covering tests, and updating REST documentation.
   * Marked all six entry contract tasks as completed.
   * Captured acceptance criteria and dependencies for each task.

8. **docs/backlogs/api/012 - installation-identity-domain-contracts.md**

   * Added the new domain contracts backlog.
   * Defined the objective of creating stable installation identity language before database, use case, or delivery implementation.
   * Scoped domain entities, persisted states, derived login classifications, value objects, internal errors, and application/domain ports.
   * Explicitly excluded migrations, Postgres repositories, login changes, refresh changes, middleware, and HTTP handlers.

9. **docs/backlogs/api/012 - installation-identity-domain-contracts_tasks.md**

   * Added the task breakdown for domain contracts.
   * Defined tasks for the installation identifier value object, app installation entity, login classifications, domain errors, read ports, write ports, restricted authorization model, and restricted authorization ports.
   * Captured acceptance criteria focused on pure domain behavior and stable interfaces without database or HTTP dependencies.

10. **docs/backlogs/api/012 - installation-identity-persistence-foundation.md**

* Removed the previous persistence foundation backlog.
* Replaced the combined persistence-oriented planning with the new split between domain contracts and database schema.

11. **docs/backlogs/api/013 - installation-identity-database-schema.md**

* Added the new database schema backlog.
* Scoped the relational foundation for `app_installations` and `installation_registration_authorizations`.
* Documented required states, constraints, indexes, public management identifiers, authorization expiration semantics, and migration expectations.
* Positioned schema creation after domain contracts and before repository implementation.

12. **docs/backlogs/api/013 - installation-identity-domain-repositories.md**

* Removed the previous combined domain and repositories backlog.
* Replaced it with separate domain contracts and repository implementation backlogs.

13. **docs/backlogs/api/014 - installation-identity-repositories.md**

* Added the new repositories backlog.
* Scoped Postgres implementations for installation and restricted authorization ports.
* Documented atomic operations for first-installation bootstrap, slot reservation, authorization consumption, authorization revocation, and logical installation revocation.
* Added concurrency requirements for preventing duplicate first installations and exceeding the known-installation limit.

14. **docs/backlogs/api/014 - installation-identity-session-tokens.md**

* Removed the previous session and tokens backlog.
* Replaced it with a more explicit session, tokens, and context backlog aligned with the new dependency sequence.

15. **docs/backlogs/api/015 - installation-identity-session-tokens-context.md**

* Added the new session, tokens, and context backlog.
* Scoped operational session binding to user and installation.
* Documented operational access token claims, restricted access token claims, restricted authorization validation, authenticated context, restricted context, and session invalidation interfaces.
* Kept login classification, installation registration, revocation, and HTTP handlers out of scope.

16. **docs/backlogs/api/015 - installation-identity-login-flow.md**

* Removed the previous login flow backlog.
* Replaced it with the new login use cases backlog aligned after repositories and session/token context.

17. **docs/backlogs/api/016 - installation-identity-login-usecases.md**

* Added the login use cases backlog.
* Scoped installation classification during login, known-installation sessions, first-installation bootstrap, revoked-installation blocking, limit handling, restricted authorization emission, and atomicity requirements.
* Explicitly kept schema, repositories, new handlers, refresh enforcement, authenticated route enforcement, and transactional password out of scope.

18. **docs/backlogs/api/016 - installation-identity-registration-flow.md**

* Removed the previous registration flow backlog.
* Replaced it with the new registration and management use cases backlog.

19. **docs/backlogs/api/017 - installation-identity-management-usecases.md**

* Added the registration and management use cases backlog.
* Scoped restricted-context registration, step-up validation, installation identifier matching, atomic slot confirmation, known-installation creation, restricted authorization consumption, operational session creation, listing, revocation, and session/token invalidation effects.
* Kept schema, repositories, final HTTP DTOs, administrative panels, and step-up for revocation out of scope.

20. **docs/backlogs/api/017 - installation-identity-management.md**

* Removed the previous management backlog.
* Replaced it with the broader registration and management use cases backlog.

21. **docs/backlogs/api/018 - installation-identity-delivery-enforcement.md**

* Added the delivery and enforcement backlog.
* Scoped HTTP connection for login results, refresh header enforcement, authenticated route header enforcement, operational middleware, restricted middleware, step-up support for installation registration, installation management handlers, HTTP error standardization, REST documentation, and test collection updates.
* Positioned delivery after the domain, persistence, repository, token, and application layers are ready.

22. **docs/backlogs/api/018 - installation-identity-refresh-enforcement.md**

* Removed the previous refresh and enforcement backlog.
* Replaced it with a broader delivery and enforcement backlog that covers handlers, middleware, REST contract, and route integration.

23. **docs/backlogs/api/019 - installation-identity-audit-retention.md**

* Renamed the backlog focus to audit, retention, and operational documentation.
* Raised the priority to High.
* Expanded scope to include retention for restricted authorizations, audited events, safe logging, revocation effects over sessions and tokens, final technical documentation, and operational validation.
* Clarified that the final documentation must not imply strong trust in `installation_id`.
* Updated dependencies to require the preceding installation identity backlogs or their final contracts.

### Conclusion

This change completes the current entry-contract documentation for installation identity by making `X-Installation-Id` an explicit login requirement and tightening validation test coverage for malformed identifiers.

It also restructures the installation identity planning model into a cleaner dependency-based backlog sequence. The result is a more coherent implementation path that separates domain contracts, persistence, repositories, session/token context, application use cases, delivery enforcement, and operational documentation.


## 2026/06/16 - api/installation-identity-02

This change introduces the first implementation layer for installation identity handling during authentication. The login flow now requires a canonical installation identifier, propagates it through the application layer, and introduces the initial classification model for installation-aware authentication decisions.

The update also establishes shared infrastructure for HTTP headers, adds installation-specific error handling, and prepares the authentication use case for future installation registration, restricted authorization, and session binding flows. In parallel, the installation identity backlog was reorganized to follow dependency-driven implementation phases.

1. **api/internal/auth/application**

   * Added `InstallationLoginClassifier` and related domain models to classify login attempts as `known`, `first`, `new`, `revoked`, or `limit_reached`.
   * Introduced installation-aware login decisions and repository contracts for installation lookup and counting.
   * Extended `LoginUserUseCase` to receive and process `InstallationID`.
   * Added support for first-installation bootstrap orchestration using transactional execution.
   * Added bootstrap race-condition protection through reclassification inside a database transaction.
   * Added configuration hooks for installation classifier, bootstrapper, and transaction manager.
   * Added comprehensive unit test coverage for installation classification and first-installation bootstrap scenarios.

2. **api/internal/auth/delivery**

   * Added mandatory `X-Installation-Id` validation during `POST /auth/login`.
   * Implemented canonical UUID v4 parsing and validation.
   * Propagated validated installation identifiers to the login use case.
   * Added HTTP contract validation for missing or malformed installation identifiers.
   * Extended handler and integration tests to cover installation header propagation and validation behavior.

3. **api/internal/shared/http**

   * Introduced centralized shared header definitions through the new `headers` package.
   * Consolidated `X-App-Token`, `X-Step-Up-Token`, and `X-Installation-Id` constants.
   * Refactored middleware and transaction handlers to consume shared header definitions instead of local constants.
   * Updated middleware tests to use the centralized header package.

4. **api/internal/shared/errors**

   * Added `INVALID_INSTALLATION_ID` error code.
   * Introduced installation-specific error registration and HTTP mapping.
   * Standardized validation failures for missing, malformed, or non-canonical installation identifiers.

5. **api/internal/account/transaction**

   * Refactored step-up token retrieval to use the centralized shared header package.
   * Removed duplicated header constant definitions from transaction delivery handlers and tests.

6. **docs/backlogs**

   * Reorganized the Installation Identity MVP backlog structure around implementation dependencies.
   * Replaced endpoint-oriented backlog decomposition with architecture-oriented phases.
   * Added new backlog documents covering:

     * Entry contract.
     * Persistence foundation.
     * Domain and repositories.
     * Session and token infrastructure.
     * Login flow.
     * Registration flow.
     * Installation management.
     * Refresh enforcement.
     * Audit and retention.
   * Removed superseded backlog documents and task breakdowns that no longer matched the revised implementation strategy.
   * Updated backlog documentation to distinguish active, completed, and archived discussions.

7. **docs/.gitignore**

   * Added support for ignoring the `olds/` backlog archive directory.

### Conclusion

This change establishes the first operational layer of the Installation Identity MVP by enforcing installation identifiers during login and introducing the application-level classification model required for future installation-aware authentication flows.

The authentication architecture is now prepared for transactional first-installation registration, installation lifecycle management, restricted authorization flows, and session binding based on installation context.

Documentation was reorganized to align implementation work with architectural dependencies, creating a clearer roadmap for the remaining phases of the Installation Identity initiative.


## 2026/06/16 - api/installation-identity-01

This change expands the Installation Identity MVP planning for the API and aligns the mobile backlog with the newly consolidated backend contracts.

The umbrella backlog was refined into endpoint-specific and infrastructure-specific planning documents, covering login, restricted authorization, step-up, installation registration, listing, revocation, refresh, session enforcement, and shared persistence/token infrastructure.

1. **docs/backlogs/README.md**

   * Expanded the API backlog index for the Installation Identity MVP.
   * Added references to the new API backlogs from `011` to `018`.
   * Replaced the previous broad description of backlog `010` with a clearer umbrella backlog role.

2. **docs/backlogs/api/010 - installation-identity-mvp.md**

   * Reworked the Installation Identity MVP backlog into a consolidated umbrella document.
   * Added `app_build` to the installation metadata model.
   * Clarified that revoked installations are denied without recovery flow in the MVP.
   * Finalized the `installation_limit_reached` response contract.
   * Replaced the generic restricted access token naming with `restricted_access_token`.
   * Defined the restricted authorization model, JWT claims, persistence table, statuses, consistency rules, and middleware expectations.
   * Clarified that transactional password must be active before issuing step-up for installation registration.
   * Defined that `POST /security/installations` completes the operational session bootstrap without requiring an intermediate login.
   * Documented the operational `access_token` claim `installation_id` and its relationship with `X-Installation-Id`.
   * Defined revocation behavior for installations, including no step-up requirement in the MVP and prevention of revoking the current installation.
   * Added explicit error contracts for invalid installation headers and installation/session mismatches.
   * Replaced the previous open-decision checklist with derived endpoint maps and shared infrastructure responsibilities.

3. **docs/backlogs/api/011 - installation-identity-auth-login.md**

   * Added the API backlog for handling installation identity during `POST /auth/login`.
   * Defined login behavior for known installations, first installation bootstrap, new installation registration, revoked installations, and installation limit enforcement.
   * Documented response contracts for `installation_registration_required` and `installation_limit_reached`.

4. **docs/backlogs/api/011 - installation-identity-auth-login_tasks.md**

   * Added detailed implementation tasks for login handling.
   * Covered header validation, installation classification, first-installation bootstrap, operational session emission, restricted authorization, revoked installation denial, limit handling, and tests.

5. **docs/backlogs/api/012 - installation-identity-step-up-authorize.md**

   * Added the API backlog for authorizing `POST /security/installations` through the existing step-up flow.
   * Defined the use of `restricted_access_token` for this flow.
   * Documented the requirement for an active transactional password before issuing the step-up token.

6. **docs/backlogs/api/012 - installation-identity-step-up-authorize_tasks.md**

   * Added implementation tasks for step-up support in installation registration.
   * Covered endpoint registration in the step-up policy, restricted access token acceptance, transactional password state validation, scoped step-up token issuance, and tests.

7. **docs/backlogs/api/013 - installation-identity-register-installation.md**

   * Added the API backlog for explicit installation registration through `POST /security/installations`.
   * Defined required headers, restricted authorization validation, step-up token consumption, installation matching, limit validation, installation creation, grant invalidation, and operational token issuance.

8. **docs/backlogs/api/013 - installation-identity-register-installation_tasks.md**

   * Added implementation tasks for the installation registration endpoint.
   * Covered use case orchestration, grant validation, step-up consumption, atomic installation creation, operational session creation, HTTP handler implementation, and test coverage.

9. **docs/backlogs/api/014 - installation-identity-list-installations.md**

   * Added the API backlog for listing user installations through `GET /security/installations`.
   * Defined authenticated access, public management identifiers, supported installation states, MVP metadata, and session/header enforcement.

10. **docs/backlogs/api/014 - installation-identity-list-installations_tasks.md**

* Added implementation tasks for installation listing.
* Covered use case creation, public response modeling, HTTP endpoint implementation, and tests for state visibility and user isolation.

11. **docs/backlogs/api/015 - installation-identity-revoke-installation.md**

* Added the API backlog for logical installation revocation through `DELETE /security/installations/{installation_resource_id}`.
* Defined revocation rules, preservation of history, prevention of current-installation revocation, and immediate access cutoff.

12. **docs/backlogs/api/015 - installation-identity-revoke-installation_tasks.md**

* Added implementation tasks for installation revocation.
* Covered logical revocation, current-installation protection, immediate session/token invalidation, endpoint exposure, and tests.

13. **docs/backlogs/api/016 - installation-identity-auth-refresh.md**

* Added the API backlog for binding `POST /auth/refresh` to the installation that originated the session.
* Defined header requirements, session-installation matching, revoked installation denial, and refreshed access token behavior.

14. **docs/backlogs/api/016 - installation-identity-auth-refresh_tasks.md**

* Added implementation tasks for refresh enforcement.
* Covered header validation, session and installation matching, revoked installation denial, refresh rotation behavior, and tests.

15. **docs/backlogs/api/017 - installation-identity-session-enforcement.md**

* Added the API backlog for enforcing `X-Installation-Id` on authenticated requests.
* Defined operational token claims, middleware validation, mismatch errors, and MVP scope boundaries.

16. **docs/backlogs/api/017 - installation-identity-session-enforcement_tasks.md**

* Added implementation tasks for session enforcement.
* Covered token claim extension, middleware adaptation, authenticated route enforcement, revoked installation denial, and tests.

17. **docs/backlogs/api/018 - installation-identity-shared-infrastructure.md**

* Added the shared infrastructure backlog for Installation Identity.
* Defined the `app_installations` and `installation_registration_authorizations` data models.
* Consolidated shared rules for operational sessions, restricted access tokens, grant uniqueness, header enforcement, and MVP recovery limitations.

18. **docs/backlogs/api/018 - installation-identity-shared-infrastructure_tasks.md**

* Added implementation tasks for shared infrastructure.
* Covered database migrations, domain and repository contracts, Postgres persistence, JWT and context support, restricted middleware, and metadata retention/auditing decisions.

19. **docs/backlogs/mobile/013 - installation-identity-mvp.md**

* Aligned the mobile Installation Identity backlog with the finalized API planning.
* Clarified that `X-Installation-Id` is mandatory from the first release of the feature.
* Added `platform`, `app_version`, and `app_build` as the minimum metadata set expected by the API.
* Updated the mobile flow to use `restricted_access_token` terminology.
* Clarified that successful installation registration returns operational tokens directly.
* Documented mobile behavior for missing or locked transactional password.
* Refined `installation_limit_reached` handling with `known_installations_count`, `max_installations`, and `next_action`.
* Added MVP revocation expectations, including no step-up requirement, disabling removal of the current installation, and immediate session loss after revocation.

### Conclusion

This change turns the Installation Identity MVP from a broad concept into a structured implementation plan split by endpoint and shared infrastructure.

The API planning now defines the main contracts, token flows, persistence models, revocation behavior, session binding, refresh enforcement, and middleware responsibilities required to support installation-aware authentication.

The mobile backlog was updated to reflect the backend decisions, keeping both sides aligned around mandatory `X-Installation-Id`, restricted registration, operational token bootstrap, installation limits, and revocation behavior.


## 2026/06/15 - api/backlog-device-identity-01

This change refines the BankLab Zero Trust backlog by replacing the previous device-oriented terminology with an installation-oriented model. The update clarifies that the MVP identifies an app installation, not the physical device, and defines `X-Installation-Id` as the shared contract between API and mobile.

The affected areas are the roadmap, backlog index, API research backlog, mobile research backlog, and the broader transactional security discussion.

1. **.gitmodules**

   * Added `templates/pandoc-latex-template` as a Git submodule.
   * Registered the external Pandoc LaTeX template repository under the `templates` folder.

2. **docs/ROADMAP.md**

   * Replaced device registration terminology with app installation terminology.
   * Updated the Zero Trust roadmap to focus on known app installations instead of physical device control.
   * Added `X-Installation-Id` as the planned installation identifier.
   * Deferred physical device identification and trust evaluation to a future evolution.

3. **docs/backlogs/README.md**

   * Replaced the API backlog reference from `device-identity-mvp` to `installation-identity-mvp`.
   * Added the new mobile backlog for installation identity.
   * Documented that API 010 and Mobile 013 share the installation identity contract.
   * Clarified that physical device identification is outside the MVP.

4. **docs/backlogs/api/010 - device-identity-mvp.md**

   * Removed the previous device identity research backlog.
   * Replaced the broader device model with a more precise installation identity model.

5. **docs/backlogs/api/010 - installation-identity-mvp.md**

   * Added the API research backlog for Installation Identity MVP.
   * Defined API responsibilities for validating `X-Installation-Id`, associating installations with users and sessions, listing installations, and revoking associations.
   * Established UUID v4 canonical format as the accepted installation identifier.
   * Documented first-installation bootstrap during login.
   * Added the restricted access flow for registering later installations.
   * Defined the step-up flow for `POST /security/installations`.
   * Introduced the maximum limit of three `known` installations per user.
   * Documented logical revocation, audit preservation, session binding, rollout concerns, and open decisions.

6. **docs/backlogs/mobile/013 - installation-identity-mvp.md**

   * Added the mobile research backlog for Installation Identity MVP.
   * Defined mobile responsibilities for generating, persisting, and sending a UUID v4 installation identifier.
   * Documented the installation lifecycle across first execution, app update, logout, user switch, reinstall, and data cleanup.
   * Introduced proposed mobile components such as `InstallationIdentityService` and `InstallationInterceptor`.
   * Defined mobile behavior for registration-required, limit-reached, restricted-token, and step-up flows.
   * Clarified that the identifier is not a secret, authentication factor, physical device ID, or proof of possession.

7. **docs/backlogs/discussion.md**

   * Renamed the second Zero Trust capability from Device Registration to Installation Registration.
   * Updated the conceptual model from `devices` to `app_installations`.
   * Reworked login and request flows around `X-Installation-Id`.
   * Added the restricted registration flow, step-up requirement, installation limit, and session-installation binding.
   * Reclassified installation context as a weak signal rather than a strong factor.
   * Deferred physical device identification, fingerprinting, attestation, and device trust to future discussions.

### Conclusion

This change narrows the Zero Trust MVP from device identity to app installation identity, making the security model more precise and less dependent on unreliable physical device assumptions.

The API and mobile responsibilities are now separated into dedicated backlogs while sharing a clear `X-Installation-Id` contract.

The resulting documentation provides a stronger basis for future implementation tasks around installation registration, session binding, revocation, and contextual policy evaluation.


## 2026/06/15 - chore/backlog-ci-transfer-fix

This change introduces a continuous integration workflow for both API and mobile applications, updates project visibility through repository status badges, and reorganizes backlog artifacts to reflect the completion of recent transactional security initiatives.

The documentation backlog has been advanced to the next Zero Trust Architecture discussion topic, focusing on device identity as a contextual security signal. At the same time, completed mobile security backlogs were moved into the historical archive, keeping the active backlog aligned with the current roadmap.

The mobile transfer flow also received a small cleanup by removing temporary diagnostic logging from the transfer confirmation process.

1. **.github/workflows/ci.yml**

   * Added a new GitHub Actions CI pipeline.
   * Configured automated execution for pull requests and pushes to `main`.
   * Added API validation through Go test execution.
   * Added mobile validation through dependency installation, formatting checks, static analysis, and Flutter test execution.
   * Enabled workflow concurrency control to prevent overlapping CI runs.

2. **README.md**

   * Added a GitHub Actions status badge.
   * Exposed CI execution status directly in the repository landing page.
   * Improved project visibility for contributors and reviewers.

3. **docs/backlogs/README.md**

   * Updated the active backlog index.
   * Removed references to completed session bootstrap and transactional transfer initiatives.
   * Added the new API research backlog for device identity.
   * Clarified that there are currently no active mobile backlogs.
   * Documented the migration of recent security-related work to the historical archive.

4. **docs/backlogs/api/010 - device-identity-mvp.md**

   * Added a new research backlog focused on device identity within the Zero Trust Architecture roadmap.
   * Defined objectives, security principles, lifecycle considerations, candidate data models, and API contracts for discussion.
   * Documented open architectural decisions regarding registration, revocation, session association, and policy enforcement.
   * Established acceptance criteria required before implementation tasks can be derived.

5. **docs/backlogs/mobile/done/**

   * Moved the internal transfer step-up authorization backlog and its task list into the historical backlog archive.
   * Preserved completed design decisions and implementation guidance for future reference.

6. **mobile/lib/ui/pages/transfer/transfer_confirmation_page.dart**

   * Removed temporary transfer diagnostic logging.
   * Eliminated the unused logging dependency and logger instance.
   * Simplified transfer completion and error-handling paths without changing user-facing behavior.

### Conclusion

This change establishes automated quality verification for both backend and mobile applications through GitHub Actions, providing continuous feedback on code quality and test execution.

The backlog structure was advanced to the next Zero Trust Architecture discussion phase by introducing device identity research while archiving recently completed transactional security work.

The mobile transfer confirmation flow was also cleaned up by removing temporary diagnostic instrumentation, keeping the implementation focused on production behavior.


## 2026/06/15 - mobile/transactional-password-10

This update refines the transfer completion experience by improving the navigation options presented after a transaction result is displayed. The changes focus on providing clearer user actions for both successful and failed transfers, aligning the flow with expected banking application behavior.

The update also includes a documentation asset refresh related to the database diagram used in the API documentation.

1. **mobile/lib/ui/pages/transfer/transfer_status_page.dart**

   * Refactored the transfer result page action area to use a consistent dual-button layout for both success and failure scenarios.
   * Updated the success flow to provide direct navigation back to the home screen through a "Cancel" action while preserving receipt access.
   * Replaced the previous single-button failure state with a dual-action interface, allowing users to either return to the home screen or immediately retry the transfer flow.
   * Added dedicated home navigation handling through a new `_navHome` method.
   * Simplified back navigation by converting `_navBack` into a concise route pop operation.
   * Removed the dependency on `BigButton` since all transfer status outcomes now use the shared `DoubleBottomButton` component.

2. **api/docs/images/database.png**

   * Updated the database diagram image used in the API documentation.
   * Refreshed the visual representation of the database structure to reflect the current project documentation state.

### Conclusion

This change improves the transfer status user experience by standardizing action controls and providing clearer navigation paths after transaction completion or failure.

The transfer flow now offers more consistent behavior across outcomes while reducing UI component variation within the status page.

Documentation assets were also refreshed to keep the database reference material aligned with the current system structure.



## 2026/06/13 - mobile/transactional-password-09

Implements the guarded step-up behavior for internal transfers, refining how transactional password errors are handled after confirmation and strengthening route safety for sensitive transfer and password setup screens.

This change also updates project documentation, mobile guidance, and tests to align the implemented behavior with the intended security model: credentials and setup data must stay typed, local, and transient.

1. **docs/backlogs/mobile/012 - step-up-transferencia-interna.md**

   * Updated the defensive `TRANSACTION_PASSWORD_NOT_SET` behavior during step-up.
   * Replaced the previous remote session refresh expectation with a local `AppSection` update to `not_set`.
   * Clarified that no new session call is required after this backend-authoritative error.
   * Preserved the rule that password setup must only be reached through a new transfer entry check.

2. **docs/backlogs/mobile/012 - step-up-transferencia-interna_tasks.md**

   * Updated task descriptions and acceptance criteria for `TRANSACTION_PASSWORD_NOT_SET`.
   * Replaced remote session refresh coverage with local `AppSection` update coverage.
   * Clarified that the defensive state change must happen without a new session query.

3. **mobile/AGENT.md**

   * Updated architectural guidance for use case orchestration.
   * Clarified that workflows spanning multiple repositories should remain in `lib/domain/usecases`.
   * Updated dependency registration order to include repositories, use cases, and view models explicitly.
   * Documented the correct module entrypoints for new repositories, use cases, and view models.

4. **mobile/Changelog.md**

   * Added the 2026/06/13 changelog entry for step-up internal transfer.
   * Documented transaction password authorization, protected transfer execution, use case orchestration, UI PIN input, and security constraints.
   * Recorded the test coverage status for the feature.

5. **mobile/README.md**

   * Expanded the transfer feature description to document step-up authorization.
   * Clarified the flow from transfer confirmation to PIN collection, step-up authorization, protected transfer execution, and retry behavior.
   * Documented that the step-up token is single-use and never persisted.

6. **mobile/docs/01-implemented-features.md**

   * Updated implemented transfer documentation to include step-up authorization.
   * Corrected relevant transfer file references to the current folder structure.
   * Added `TransferUsecase`, `TransactionPasswordRepository`, `TransactionPasswordApi`, and `TransactionPasswordInputPage` to the documented implementation surface.
   * Clarified token lifecycle, idempotency behavior, and backend error-code preservation.

7. **mobile/lib/core/routing/extra_codec.dart**

   * Removed support for serializing arbitrary maps and lists through route extras.
   * Preserved primitive, typed DTO, and typed origin serialization.
   * Added safe fallback decoding for unknown `TransactionPasswordSetupOrigin` values, defaulting to `postLogin`.
   * Reduced the risk of passing credentials or sensitive setup payloads through encoded route extras.

8. **mobile/lib/core/routing/routes/transaction_password_routes.dart**

   * Added defensive redirects for transaction password introduction and creation routes when the required typed origin is missing.
   * Redirected direct access to the confirm route back to Home.
   * Injected `TransactionPasswordViewModel` directly into the creation page.
   * Removed routed confirmation page construction based on map extras.

9. **mobile/lib/core/routing/routes/transfer_routes.dart**

   * Added route guards for payment and confirmation routes.
   * Redirected missing or invalid extras to Home.
   * Prevented direct navigation into transfer steps without the required typed data.

10. **mobile/lib/data/services/apis/contact_verification/enums/contact_verification_channel.dart**

* Normalized enum formatting for `ContactVerificationChannel`.

11. **mobile/lib/domain/common/receipt/enums/transfer_receipt_status.dart**

* Normalized enum formatting for `TransferReceiptStatus`.

12. **mobile/lib/domain/common/user/enums/user_role.dart**

* Normalized enum formatting for `UserRole`.

13. **mobile/lib/ui/pages/transaction_password/setup/create_transaction_password_page.dart**

* Added direct `TransactionPasswordViewModel` injection.
* Replaced route-extra based confirmation navigation with in-memory `MaterialPageRoute` navigation.
* Passed the created PIN and setup origin directly to `ConfirmTransactionPasswordPage`.
* Avoided serializing the PIN through router extras.

14. **mobile/lib/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart**

* Simplified post-setup navigation for the transfer origin.
* Replaced manual stack manipulation with direct navigation to the transfer recipient route.
* Removed the no longer needed context extension import.

15. **mobile/lib/ui/pages/transfer/transfer_confirmation_page.dart**

* Refactored transfer submission to build a stable `TransferDraft` before step-up.
* Added protected transfer execution with PIN collection through the transaction password route.
* Preserved the same `idempotency_key` across retryable step-up attempts.
* Added explicit handling for step-up backend error codes:

  * invalid transaction password requests a new PIN before transfer execution;
  * expired or consumed step-up token requests a new PIN and retries with the same draft;
  * locked transaction password blocks the current attempt;
  * missing transaction password returns the user to Home;
  * unknown or non-retryable authorization failures go to the failure flow.
* Ensured the transfer failure screen is not shown before offering retry for retryable step-up errors.

16. **mobile/test/core/routing/extra_codec_test.dart**

* Added coverage to reject arbitrary map serialization.
* Added coverage for typed transaction password setup origin serialization.
* Added coverage for safe fallback decoding of unknown setup origins.

17. **mobile/test/core/routing/protected_route_fallback_test.dart**

* Added route fallback tests for transaction password and transfer routes.
* Verified that sensitive routes redirect to Home when required extras are missing.

18. **mobile/test/ui/pages/transaction_password/setup/create_transaction_password_page_test.dart**

* Updated the creation page test to use a real `TransactionPasswordViewModel`.
* Verified that PIN and origin are propagated only in memory to the confirmation page.
* Added a fake transaction password repository for the test boundary.

19. **mobile/test/ui/pages/transfer/transfer_confirmation_page_test.dart**

* Added coverage for invalid transaction password retry without transfer execution.
* Added coverage for expired and consumed step-up tokens using a stable idempotency key.
* Added coverage for locked transaction password behavior.
* Added coverage for `TRANSACTION_PASSWORD_NOT_SET` returning to Home without transfer execution.
* Added coverage for non-retryable authorization failures reaching the failure flow.
* Added coverage to ensure rejected transfer tokens are not reused.
* Extended test fakes to record PINs, transfer tokens, requests, and sequential configured results.

### Conclusion

This change set tightens the internal transfer step-up flow by keeping sensitive PIN data out of router serialization, protecting direct access to sensitive routes, and handling backend step-up error codes with explicit user-flow decisions.

The resulting behavior is more secure and predictable: retryable step-up failures request a new PIN while preserving transfer idempotency, non-retryable failures terminate correctly, and `TRANSACTION_PASSWORD_NOT_SET` now updates local application state before returning to Home.

Documentation, architectural guidance, changelog entries, and tests were updated to keep the implemented behavior aligned with the BankLab mobile security model.


## 2026/06/13 - mobile/transactional-password-08

This change refines the transactional password flow used during internal transfers, consolidating route usage, improving submit-state handling, and adding widget coverage for the PIN prompt and transfer confirmation behavior.

The update also introduces a reusable loading state for `BigButton`, preventing duplicate transfer submissions while the protected operation is running.

1. **mobile/lib/core/routing/routes.dart**

   * Removed the transfer-specific `verifyTransactionPassword` route entry.
   * Kept transactional password verification under the dedicated transaction password route module.

2. **mobile/lib/core/routing/routes/transfer_routes.dart**

   * Updated the transfer flow to open `TransactionPasswordRoutes.transactionPassword`.
   * Reused the centralized transaction password route name and path instead of defining the verification route inside `TransferRoutes`.

3. **mobile/lib/ui/components/buttons/big_button.dart**

   * Added the `isRunning` property to represent an in-progress button action.
   * Disabled button presses while `isRunning` is active.
   * Rendered a compact `CircularProgressIndicator` in place of the configured left or right icon during loading.

4. **mobile/lib/ui/pages/transaction_password/verification/transaction_password_input_page.dart**

   * Simplified the right action icon to always show the confirmation icon.
   * Removed the local loading indicator logic from the transaction password prompt button.

5. **mobile/lib/ui/pages/transfer/transfer_confirmation_page.dart**

   * Added `_hasSubmitted` state tracking to prevent duplicate submissions.
   * Disabled the transfer button after the first submit attempt until the flow completes or is canceled.
   * Connected the button loading state to the transfer command execution state.
   * Reset the submission guard when the PIN prompt is canceled or when an unknown authorization error occurs.
   * Updated the cancellation feedback message to clarify that the transaction password was not provided.

6. **mobile/test/ui/pages/transaction_password/verification/transaction_password_input_page_test.dart**

   * Added widget tests for the transactional password input page.
   * Verified that submission remains disabled until all six digits are entered.
   * Verified that the PIN is masked in the UI while the original six-digit value is returned to the caller.
   * Verified that canceling returns `null` and that reopening the prompt starts with an empty field.

7. **mobile/test/ui/pages/transfer/transfer_confirmation_page_test.dart**

   * Added widget tests for the transfer confirmation page.
   * Verified that the PIN input page opens before authorization or transfer execution.
   * Verified that canceling the PIN prompt does not call authorization or transfer and allows reopening the flow.
   * Verified duplicate-tap protection, loading feedback, authorization token usage, and success navigation.
   * Verified that transfer failure navigates to the existing failure status flow.

### Conclusion

The transfer confirmation flow now uses the centralized transactional password route and has stronger protection against duplicate submissions.

The UI behavior is more consistent through `BigButton.isRunning`, and the new widget tests cover the critical PIN, authorization, transfer, cancelation, success, and failure paths.


## 2026/06/13 - mobile/transactional-password-07

This change integrates the transactional password step-up flow into the internal transfer execution path. The transfer confirmation flow now collects the transactional PIN, delegates authorization to the transfer use case, and sends the resulting step-up token through the transfer API header.

The implementation also simplifies the transactional password verification page, moves authorization responsibility out of the UI layer, strengthens idempotency handling, and expands tests around protected transfers, step-up token propagation, and retry behavior.

1. **mobile/lib/core/routing/routes.dart**

   * Added a dedicated transactional password verification route at `/transaction-password/verify`.
   * Exposed the route through `TransactionPasswordRoutes.transactionPassword`.

2. **mobile/lib/core/routing/routes/transfer_routes.dart**

   * Replaced the old `VerifyTansactionPasswordPage` route target with `TransactionPasswordInputPage`.
   * Removed the direct dependency on `VerifyTansactionPasswordViewmodel`.
   * Kept the transfer verification route responsible only for collecting the PIN and returning it to the caller.

3. **mobile/lib/data/repositories/transfer/transfer_repository.dart**

   * Updated the transfer contract to require both a step-up token and a transfer request DTO.
   * Made token usage explicit at repository boundary.

4. **mobile/lib/data/repositories/transfer/transfer_repository_impl.dart**

   * Added validation for blank step-up tokens before calling the API.
   * Forwarded the token and DTO to `ApiTransfer.transfer`.
   * Preserved the existing `lastTransfer` cache behavior for success and failure responses.

5. **mobile/lib/data/services/apis/transfer/api_transfer.dart**

   * Updated the internal transfer request to send the step-up token in the `X-Step-Up-Token` header.
   * Kept transfer payload data restricted to transfer fields, avoiding sensitive authorization data in the request body.

6. **mobile/lib/data/services/apis/transfer/dtos/transfer_request_dto.dart**

   * Added `TransferRequestDto.fromTransferDraft`.
   * Centralized conversion from domain transfer draft to API transfer request.
   * Prevented invalid DTO creation when the destination account is empty or equal to the source account.

7. **mobile/lib/domain/usecases/transfer/inputs/protected_transfer_input.dart**

   * Added `ProtectedTransferInput` to group the transfer draft and transactional PIN.
   * Introduced a domain input model for protected transfer execution.

8. **mobile/lib/domain/usecases/transfer/inputs/transfer_draft.dart**

   * Made `idempotencyKey` required.
   * Removed the empty default value to force explicit idempotency handling by the caller.

9. **mobile/lib/domain/usecases/transfer/transfer_usecase.dart**

   * Injected `TransactionPasswordRepository`.
   * Changed `transfer` to receive `ProtectedTransferInput`.
   * Validated selected account, idempotency key, and destination account before authorization.
   * Authorized the transactional password before executing the transfer.
   * Passed the returned step-up token to the transfer repository.
   * Returned authorization failures without attempting the transfer.

10. **mobile/lib/ui/pages/transaction_password/verification/transaction_password_input_page.dart**

* Added a lightweight transactional password input page.
* Used `TokenInput` to collect a six-digit hidden PIN.
* Returned the entered PIN through navigation instead of authorizing directly.
* Supported cancel and submit actions through the bottom button layout.

11. **mobile/lib/ui/pages/transaction_password/verification/verify_tansaction_password_page.dart**

* Removed the previous verification page implementation.
* Eliminated UI-level step-up authorization and related error handling from this screen.

12. **mobile/lib/ui/pages/transaction_password/verification/viewmodel/verify_tansaction_password_viewmodel.dart**

* Removed the verification view model.
* Moved transactional password authorization responsibility into the transfer use case.

13. **mobile/lib/ui/pages/transfer/transfer_confirmation_page.dart**

* Added UUID v7 idempotency key generation per confirmation page instance.
* Collected the transactional PIN before submitting the transfer.
* Built `ProtectedTransferInput` with transfer draft and PIN.
* Displayed an informational snackbar when the operation is cancelled.
* Fixed success navigation to use `TransferRoutes.statusSuccess.routeName`.

14. **mobile/lib/ui/pages/transfer/transfer_payment_page.dart**

* Removed direct disposal of the shared transfer view model.

15. **mobile/lib/ui/pages/transfer/transfer_recipient_page.dart**

* Removed direct disposal of the shared transfer view model.

16. **mobile/lib/ui/pages/transfer/viewmodel/transfer_viewmodel.dart**

* Changed the transfer command input from `TransferDraft` to `ProtectedTransferInput`.
* Delegated transfer execution directly to `TransferUsecase.transfer`.
* Removed view-model-level UUID generation and unused disposal logic.
* Cleaned up commented recipient selection state.

17. **mobile/lib/ui/viewmodels.dart**

* Removed registration of `VerifyTansactionPasswordViewmodel`.
* Kept transactional password setup and transfer view models registered.

18. **mobile/test/data/repositories/transaction/transaction_repository_impl_test.dart**

* Updated transfer repository tests to pass step-up tokens.
* Added assertion that the token reaches the fake API.
* Added validation coverage for blank step-up tokens.
* Updated fake API transfer implementation to capture token and DTO separately.

19. **mobile/test/data/services/apis/transfer/api_transfer_test.dart**

* Updated transfer API tests to pass the step-up token.
* Added assertions for the `X-Step-Up-Token` request header.
* Verified that step-up token and transactional password are not sent in the request body.
* Added coverage for preserved backend step-up error codes.
* Confirmed recipient lookup requests do not send headers.

20. **mobile/test/domain/usecases/transfer/transfer_usecase_test.dart**

* Added fake transactional password repository support.
* Updated transfer use case tests to use `ProtectedTransferInput`.
* Verified authorization happens before transfer execution.
* Added coverage for failed authorization without transfer execution.
* Tested one-token-per-transfer-attempt behavior.
* Tested retry behavior with the same idempotency key and a new PIN/token.
* Tested distinct idempotency keys across distinct transfer attempts.
* Verified validation failures avoid unnecessary authorization and transfer calls.

21. **mobile/test/ui/pages/transaction_password/verification/viewmodel/verify_tansaction_password_viewmodel_test.dart**

* Removed tests for the deleted verification view model.
* Authorization behavior is now covered through transfer use case tests.

### Conclusion

This change completes the protected internal transfer flow by moving transactional password authorization into the domain use case and sending the resulting step-up token through the transfer API header.

The UI now only collects the transactional PIN, while the use case coordinates validation, authorization, idempotency, and transfer execution. This improves separation of concerns and makes protected transfer behavior easier to test.

The test suite was expanded to cover token propagation, authorization failures, retry behavior, idempotency preservation, and API request structure.


## 2026/06/13 - mobile/transactional-password-06

This change refactors the mobile transfer module to use clearer domain naming and a flatter UI page structure. The previous `transaction` repository naming was replaced by `transfer`, aligning the repository contract, implementation, dependency injection, use cases, and tests with the actual transfer flow.

It also moves transfer UI pages out of the `home/transfer` subtree into a dedicated `transfer` page module, updating routing, codec imports, view model registration, and page documentation accordingly.

1. **mobile/lib/data/repositories/transfer**

   * Renamed `TransactionRepository` to `TransferRepository`.
   * Renamed `TransactionRepositoryImpl` to `TransferRepositoryImpl`.
   * Moved the repository files from `data/repositories/transaction` to `data/repositories/transfer`.
   * Updated the implementation import to reference the new repository contract.

2. **mobile/lib/data/repositories.dart**

   * Replaced the transaction repository registration with `TransferRepository`.
   * Updated dependency injection to bind `TransferRepository` to `TransferRepositoryImpl`.

3. **mobile/lib/domain/usecases/transfer/transfer_usecase.dart**

   * Updated the use case dependency from `TransactionRepository` to `TransferRepository`.
   * Renamed the injected field and constructor parameter from transaction-oriented naming to transfer-oriented naming.
   * Updated transfer, receipt, and recipient lookup calls to use the new repository dependency.

4. **mobile/lib/domain/usecases/details/details_usecase.dart**

   * Replaced the transaction repository dependency with `TransferRepository`.
   * Updated the receipt retrieval flow to use the renamed transfer repository.
   * Normalized the account summary DTO import path.

5. **mobile/lib/ui/pages/transfer**

   * Moved transfer pages, models, view model, and widgets from `ui/pages/home/transfer` to `ui/pages/transfer`.
   * Preserved the existing transfer page implementations while promoting the transfer flow to its own page module.

6. **mobile/lib/core/routing**

   * Updated transfer route imports to point to the new `ui/pages/transfer` module.
   * Updated `ExtraCodec` to import `TransferConfirmationData` from the new transfer page path.
   * Kept transaction password verification routing integrated with the transfer route flow.

7. **mobile/lib/ui/viewmodels.dart**

   * Updated the transfer view model import to use the new `pages/transfer/viewmodel` location.
   * Preserved the existing view model registration behavior.

8. **mobile/lib/ui/pages/AGENT.md**

   * Updated page layout documentation to reflect the new `transfer/models` location.
   * Updated the current example for `transfer_confirmation_data.dart` to match the new module structure.

9. **mobile/test/data/repositories/transaction/transaction_repository_impl_test.dart**

   * Updated repository implementation imports from `transaction` to `transfer`.
   * Replaced test construction of `TransactionRepositoryImpl` with `TransferRepositoryImpl`.
   * Preserved the existing transfer, receipt, and recipient lookup test coverage.

10. **mobile/test/domain/usecases/transfer/transfer_usecase_test.dart**

* Updated the fake repository contract from `TransactionRepository` to `TransferRepository`.
* Renamed the fake test implementation from `_FakeTransactionRepository` to `_FakeTransferRepository`.
* Updated test variables and expectations to use transfer-oriented naming while preserving existing use case behavior coverage.

### Conclusion

This change consolidates transfer-related naming and structure across the mobile application. The repository layer now reflects the transfer domain explicitly, and the transfer UI flow is organized as its own page module instead of being nested under `home`.

The refactor improves architectural clarity without changing the existing transfer behavior covered by the current tests.


## 2026/06/13 - mobile/transactional-password-05

This change refactors the transaction password verification flow by moving the internal transfer authorization context into the repository layer. The UI and ViewModel now work with a simplified contract that only requires the transaction password, while the repository becomes responsible for building the step-up authorization request.

The update reduces coupling between presentation and API-specific details, centralizing the internal transfer operation mapping within the transaction password repository implementation. The affected areas include repository contracts, ViewModel commands, verification UI, and automated tests.

1. **mobile/lib/data/repositories/transaction_password**

   * Refactored the repository contract by replacing `stepUpAuthorize(StepUpAuthorizeRequestDto)` with `authorizeInternalTransfer(String transactionPassword)`.
   * Removed the need for callers to construct step-up authorization request DTOs.
   * Centralized internal transfer authorization as an explicit repository operation.

2. **mobile/lib/data/repositories/transaction_password/transaction_password_repository_impl.dart**

   * Added responsibility for creating `StepUpAuthorizeRequestDto` inside the repository implementation.
   * Automatically maps internal transfer authorization requests to `StepUpOperation.internalTransfer`.
   * Encapsulated API-specific request construction details within the data layer.
   * Preserved existing API integration and error propagation behavior.

3. **mobile/lib/ui/pages/transaction_password/verification**

   * Simplified the verification page by removing direct dependencies on step-up authorization DTOs and operation enums.
   * Updated command execution to submit only the transaction password.
   * Renamed command references from `stepUpAuthorize` to `authorizeInternalTransfer`.
   * Kept existing success, failure, and navigation behavior unchanged.

4. **mobile/lib/ui/pages/transaction_password/verification/viewmodel**

   * Updated the ViewModel command definition to use `String` as input instead of `StepUpAuthorizeRequestDto`.
   * Renamed the exposed command to better represent the business action being performed.
   * Reduced presentation-layer knowledge of backend authorization request structures.

5. **mobile/test/data/repositories/transaction_password**

   * Updated repository tests to validate automatic request creation inside the repository.
   * Added assertions verifying the generated operation, transaction password, and serialized request payload.
   * Adjusted backend error validation to cover updated authorization error codes.
   * Removed helper methods that are no longer required by the new contract.

6. **mobile/test/ui/pages/auth/transaction_password and mobile/test/ui/pages/transaction_password/setup**

   * Updated repository test doubles to match the new repository interface.
   * Replaced DTO-based authorization methods with transaction-password-based authorization methods.

7. **mobile/test/ui/pages/transaction_password/verification/viewmodel**

   * Added dedicated ViewModel tests for `authorizeInternalTransfer`.
   * Verified delegation of the transaction password to the repository.
   * Verified successful result propagation.
   * Verified backend errors are exposed unchanged to the presentation layer.

### Conclusion

This refactoring simplifies the transaction password verification flow by removing API request construction responsibilities from the UI and ViewModel layers.

The repository now acts as the boundary responsible for translating business actions into step-up authorization requests, improving encapsulation and reducing coupling with API-specific DTOs and operation identifiers.

The accompanying test updates ensure the new contract is validated across repository and presentation layers while preserving existing authorization behavior.


## 2026/06/12 - mobile/transactional-password-04

This change consolidates the mobile step-up flow for internal transfers by making the Home screen route users according to the transaction password status already stored in the authenticated session.

It also improves step-up authorization handling, preserves backend error codes across API layers, redacts sensitive HTTP headers in logs, and simplifies authentication view model reuse between full login and short login.

1. **docs/backlogs/mobile/012 - step-up-transferencia-interna_tasks.md**

   * Updated the task scope for transfer step-up navigation.
   * Replaced the previous verification-operation approach with direct exposure of `TransactionPasswordStatus` through `HomeViewmodel`.
   * Adjusted acceptance criteria to focus on UI routing behavior for `active`, `not_set`, `locked`, and `unknown`.

2. **mobile/lib/core/services/app_section/app_section.dart**

   * Added accessors for `hasActiveTransactionPassword` and `transactionPasswordStatus`.
   * Centralized transaction password readiness lookup from the current authenticated session.
   * Defaulted missing session status to `TransactionPasswordStatus.notSet`.

3. **mobile/lib/ui/pages/home/home_page.dart** and **mobile/lib/ui/pages/home/viewmodel/home_viewmodel.dart**

   * Exposed `transactionPasswordStatus` from `HomeViewmodel`.
   * Updated the **Transferir** action to:

     * open recipient selection when the status is `active`;
     * open transaction password setup with transfer origin when the status is `notSet`;
     * keep the user on Home and show an error message when the status is `locked` or `unknown`.
   * Kept session interpretation outside the Home page by delegating status access to the view model.

4. **mobile/lib/core/routing/extensions/context_extencions.dart**

   * Replaced the commented routing helper with a concrete `BuildContext.popUntil(String routeName)` extension.
   * Added support for popping the navigation stack until a named route is reached, preserving the first route when the target is not found.

5. **mobile/lib/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart**

   * Updated the post-setup transfer flow.
   * After creating a transaction password from the transfer origin, the app now returns to Home and then pushes the transfer recipient route.

6. **mobile/lib/data/services/apis/transaction_password/enums/transaction_password_status.dart**

   * Expanded transaction password statuses to `active`, `notSet`, `locked`, and `unknown`.
   * Changed unsupported backend status values to map to `unknown` instead of throwing.
   * Removed the duplicated domain enum and re-exported the API enum through the auth session model.

7. **mobile/lib/domain/common/auth/models/auth_session/**

   * Removed the duplicated `TransactionPasswordStatus` enum from the domain auth session folder.
   * Updated `ReadinessSession` and auth session exports to use the shared transaction password status enum from the API layer.

8. **mobile/lib/data/services/apis/transaction_password/dtos/**

   * Renamed `SetUpAuthorizeRequestDto` to `StepUpAuthorizeRequestDto`.
   * Renamed `SetUpAuthorizeResponseDto` to `StepUpAuthorizeResponseDto`.
   * Updated request and response usage across repositories, APIs, view models, pages, and tests.

9. **mobile/lib/data/services/apis/transaction_password/transaction_password_api.dart**

   * Updated step-up authorization to use the renamed DTOs.
   * Added error-envelope handling for step-up authorization responses.
   * Preserved backend error details through `ApiError.toAppErrorDetails()`.
   * Kept parsing and missing-data failures explicit for malformed or incomplete success envelopes.

10. **mobile/lib/data/repositories/transaction_password/transaction_password_repository.dart** and **mobile/lib/data/repositories/transaction_password/transaction_password_repository_impl.dart**

* Updated repository contracts to use the renamed step-up DTOs.
* Ensured step-up authorization failures are returned explicitly.
* Preserved backend error codes while still updating `AppSection` when the backend reports that the transaction password is not set.

11. **mobile/lib/data/services/apis/core/api_error.dart**

* Added `toAppErrorDetails()` to normalize backend error codes and details into `AppError.details`.
* Reused this mapping from transaction password and transfer API flows.

12. **mobile/lib/core/services/client_http/dio/dio_error_mapper.dart**

* Normalized backend error details so backend error codes remain available even when the backend also sends nested details.
* Updated mappings for account approval, contact verification, transaction password locked, transaction password not set, and generic HTTP errors.

13. **mobile/lib/core/services/client_http/dio/dio_rest_client.dart**

* Added request logging for all HTTP methods.
* Redacted sensitive headers, including `Authorization`, `X-App-Token`, and `X-Step-Up-Token`.
* Avoided logging request bodies, preventing transaction passwords and step-up credentials from appearing in logs.

14. **mobile/lib/data/services/apis/transfer/api_transfer.dart**

* Preserved backend error codes and details for transfer and recipient API error envelopes.
* Replaced direct envelope error access with safer pattern matching.

15. **mobile/lib/ui/pages/auth/** and **mobile/lib/ui/viewmodels.dart**

* Moved `LoginViewModel` to a shared auth view model folder.
* Reused `LoginViewModel` for both login and short-login flows.
* Removed the dedicated `ShortLoginViewModel`.
* Updated route registration and dependency injection accordingly.

16. **mobile/lib/ui/pages/transaction_password/verification/**

* Updated transaction password verification to use `StepUpAuthorizeRequestDto` and `StepUpAuthorizeResponseDto`.
* Simplified request creation during verification submit.

17. **mobile/test/core/routing/context_extensions_test.dart**

* Added widget tests for `popUntil`.
* Covered route removal up to a named route and fallback behavior when the route name is not found.

18. **mobile/test/ui/pages/home/home_transaction_password_navigation_test.dart**

* Added Home navigation tests for all transaction password statuses.
* Verified recipient navigation, setup navigation with transfer origin, locked state messaging, and unknown state error handling.

19. **mobile/test/core/services/client_http/dio/**

* Added coverage for backend code preservation in mapped HTTP errors.
* Added tests ensuring sensitive headers and transaction password bodies are not exposed in logs.

20. **mobile/test/data/services/apis/transaction_password/**

* Added step-up authorization API tests for successful token parsing, backend error envelopes, missing data, malformed data, and propagated RestClient failures.
* Added DTO tests for step-up request serialization and response parsing.
* Updated status parsing tests to validate fallback to `unknown`.

21. **mobile/test/data/repositories/transaction_password/transaction_password_repository_impl_test.dart**

* Added success-path coverage for step-up authorization.
* Updated failure tests to use renamed DTOs.
* Verified backend code preservation for transaction password and step-up authorization errors.

22. **mobile/test/data/services/apis/transfer/api_transfer_test.dart**

* Added assertions to ensure transfer and recipient backend error codes remain available through `backendErrorCode`.

23. **mobile/test/ui/pages/auth/post_login_destination_test.dart** and **mobile/test/ui/pages/transaction_password/setup/confirm_transaction_password_page_test.dart**

* Updated tests to use the shared `LoginViewModel`.
* Updated transaction password test doubles to use the renamed step-up DTOs.

### Conclusion

This change completes the Home-driven decision flow for internal transfer access based on the transaction password status already loaded in the session.

It also strengthens the step-up authorization path by preserving backend error codes, improving API error normalization, and protecting sensitive credentials from logs.

The refactor reduces duplicated authentication view model code, removes duplicated transaction password status definitions, and adds test coverage for routing, error handling, DTOs, repositories, APIs, and logging behavior.


## 2026/06/12 - mobile/transactional-password-03

Implemented origin-aware transaction password setup navigation, allowing the same setup flow to be reused after login and during internal transfer step-up requirements.

This change updates routing extras, setup pages, confirmation behavior, and tests so the transaction password flow can return to the correct destination after completion or cancellation.

1. **docs/backlogs/mobile/012 - step-up-transferencia-interna_tasks.md**

   * Removed the acceptance note requiring PIN and confirmation cleanup on completion, cancellation, or failure.
   * Aligned the backlog with the updated setup behavior covered by the implementation and tests.

2. **mobile/lib/core/routing/extra_codec.dart**

   * Added serialization and deserialization support for `TransactionPasswordSetupOrigin`.
   * Updated imports to use project-root absolute paths.
   * Enabled setup origin data to survive route extra encoding.

3. **mobile/lib/core/routing/models/transaction_password_setup_origin.dart**

   * Added `TransactionPasswordSetupOrigin` with `postLogin` and `transfer` origins.
   * Added `fromName` factory validation for restoring origins from route extras.

4. **mobile/lib/core/routing/routes/transaction_password_routes.dart**

   * Updated transaction password routes to require and propagate setup origin.
   * Renamed the introduction page import to `IntroductionTransactionPasswordPage`.
   * Changed confirmation route extra handling from a raw PIN string to a map containing token and origin.
   * Added fallback behavior to return to the introduction page when confirmation route data is invalid.

5. **mobile/lib/ui/pages/auth/login/login_page.dart**

   * Passed `TransactionPasswordSetupOrigin.postLogin` when navigating to transaction password setup after login.
   * Normalized imports to project-root absolute paths.

6. **mobile/lib/ui/pages/auth/short_login/short_login_page.dart**

   * Passed `TransactionPasswordSetupOrigin.postLogin` when short login requires transaction password setup.
   * Normalized input component imports.

7. **mobile/lib/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart**

   * Replaced `pin` with `token` and added setup `origin`.
   * Used the origin to decide the destination after successful creation.
   * Navigates to Home for post-login setup and to transfer recipient flow for transfer setup.
   * Preserved failure behavior by keeping the setup screen open when creation fails.

8. **mobile/lib/ui/pages/transaction_password/setup/create_transaction_password_page.dart**

   * Added setup origin as a required page parameter.
   * Propagated both token and origin to the confirmation route.
   * Replaced introduction navigation with stack pop behavior for the back action.
   * Renamed internal PIN state to token.

9. **mobile/lib/ui/pages/transaction_password/setup/introduction_transaction_password_page.dart**

   * Renamed the introduction page class and file.
   * Added setup origin as a required constructor parameter.
   * Updated navigation so cancellation pops back to the caller.
   * Propagated origin when opening the create password page.
   * Extracted the information item widget from the page.

10. **mobile/lib/ui/pages/transaction_password/setup/widgets/information_item.dart**

* Added reusable `InformationItem` widget for transaction password setup guidance.
* Moved presentation details out of the introduction page.

11. **mobile/lib/ui/pages/transaction_password/verification/verify_tansaction_password_page.dart**

* Normalized routing imports to project-root absolute paths.

12. **mobile/test/ui/pages/auth/transaction_password/transaction_password_introduction_page_test.dart**

* Updated tests for the renamed introduction page.
* Added origin propagation validation when starting password creation.
* Added cancellation behavior coverage to ensure the flow returns to the caller.
* Added validation for unsupported setup origins.

13. **mobile/test/ui/pages/transaction_password/setup/create_transaction_password_page_test.dart**

* Added test coverage confirming that the create page forwards both token and origin to confirmation.

14. **mobile/test/ui/pages/transaction_password/setup/confirm_transaction_password_page_test.dart**

* Added tests for successful post-login setup navigation to Home.
* Added tests for transfer setup navigation to the recipient flow.
* Added coverage for inactive local status after creation.
* Added failure retry behavior validation.
* Added handling validation for the transaction password already set error.

### Conclusion

The transaction password setup flow is now reusable across post-login onboarding and transfer step-up scenarios.

Routing now carries an explicit setup origin, allowing each entry point to complete or cancel back to the correct destination while preserving validation and failure behavior.

The added tests strengthen confidence around origin propagation, navigation outcomes, retry behavior, and setup error handling.


## 2026/06/11 - mobile/transactional-password-02

Implements the mobile step-up authorization flow for internal transfers using the transactional password. The change adds the authorization API contract, repository integration, verification screen, route registration, error mapping, and local session snapshot updates based on backend responses.

It also refactors route enums to expose explicit `routePath` and `routeName` properties, reducing ambiguity between enum names and URL paths across navigation calls.

1. **mobile/lib/core/routing/routes.dart**

   * Added the `AppRoute` interface with `routePath` and `routeName`.
   * Updated route enums to implement `AppRoute`.
   * Replaced the generic `path` field with `routePath`.
   * Added `TransferRoutes.verifyTransactionPassword`.
   * Renamed `RegisterRoutes.name` to `RegisterRoutes.fullName` to avoid semantic ambiguity.

2. **mobile/lib/core/routing/**

   * Updated router initialization and route declarations to use `routePath` and `routeName`.
   * Registered the transactional password verification route in the transfer flow.
   * Updated auth, base, register, shared, transfer, and transactional password route files.

3. **mobile/lib/ui/pages/transaction_password/setup/**

   * Renamed the previous `creation_flow` folder to `setup`.
   * Updated imports and navigation calls to match the new route contract.
   * Preserved the existing transactional password creation behavior.

4. **mobile/lib/ui/pages/transaction_password/verification/**

   * Added `VerifyTansactionPasswordPage` for step-up password confirmation.
   * Added `VerifyTansactionPasswordViewmodel`.
   * Integrated `TokenInput`, loading state, cancel behavior, success pop with response, and error handling.
   * Handles locked transactional password and missing transactional password scenarios.

5. **mobile/lib/data/services/apis/transaction_password/**

   * Added step-up authorization request and response DTOs.
   * Added `HttpMethod` and `StepUpOperation` enums.
   * Implemented `TransactionPasswordApi.stepUpAuthorize()`.
   * Added structured logging for create and step-up authorization failures.

6. **mobile/lib/data/repositories/transaction_password/**

   * Extended `TransactionPasswordRepository` with `stepUpAuthorize`.
   * Implemented repository forwarding to the API.
   * Updates `AppSection` locally to `notSet` when the backend returns `TRANSACTION_PASSWORD_NOT_SET`.

7. **mobile/lib/core/services/app_section/app_section.dart**

   * Added `markTransactionPasswordAsNotSet()`.
   * Allows the mobile snapshot to be corrected locally without calling `GET /auth/session` again.

8. **mobile/lib/core/services/client_http/dio/dio_error_mapper.dart**

   * Added mapping for `TRANSACTION_PASSWORD_LOCKED`.
   * Added mapping for `TRANSACTION_PASSWORD_NOT_SET`.
   * Exposes these backend errors as typed `AppErrorCode` values.

9. **mobile/lib/core/result/errors/app_error_code.dart**

   * Added `transactionPasswordLocked`.
   * Added `transactionPasswordNotSet`.

10. **mobile/lib/data/repositories/auth/auth_repository_impl.dart**

   * Renamed local session-loading variables for clarity.
   * Changed logout to clear `AppSection` before checking login state.
   * Ensures stale session snapshots are removed even when the repository is already not logged in.

11. **mobile/lib/ui/pages/auth, home, register, splash, transfer**

   * Updated navigation calls to use `routeName`.
   * Replaced direct enum `name` usage in `goNamed` and `pushNamed`.
   * Updated registration navigation to use `RegisterRoutes.fullName`.

12. **mobile/lib/ui/viewmodels.dart**

   * Updated transactional password setup import path.
   * Registered `VerifyTansactionPasswordViewmodel` in dependency injection.

13. **docs/backlogs/mobile/012 - step-up-transferencia-interna_tasks.md**

   * Updated the step-up consistency strategy.
   * Replaced the defensive remote refresh approach with local snapshot correction on `TRANSACTION_PASSWORD_NOT_SET`.
   * Clarified that the backend error is treated as authoritative.

14. **CHANGELOG.md**

   * Corrected login navigation documentation from `HomeRoutes.home.name` to `HomeRoutes.home.routeName`.

15. **mobile/test/**

   * Added error mapper coverage for `TRANSACTION_PASSWORD_LOCKED`.
   * Added auth repository tests for cached session retrieval and logout snapshot clearing.
   * Added transaction password repository tests for local status updates after creation and after step-up failures.
   * Updated transactional password setup tests to use the renamed folder and new route path API.

16. **mobile/android/build.gradle.kts**

   * Updated subproject build directory resolution from `project.name` to `project.routeName`.
   * Aligns the Gradle configuration with the route naming refactor introduced in the mobile code.

### Conclusion

This update completes the client-side foundation for transactional password step-up authorization in internal transfers.

The implementation introduces a dedicated authorization flow, typed backend error handling, local session snapshot reconciliation, and a reusable verification screen integrated into the routing system. The application now reacts immediately to authoritative backend responses such as `TRANSACTION_PASSWORD_NOT_SET` and `TRANSACTION_PASSWORD_LOCKED`, avoiding unnecessary session refreshes and keeping the local authentication state consistent.

Additionally, the routing layer was standardized through the introduction of `AppRoute`, improving navigation consistency and reducing coupling to enum implementation details. The transactional password creation flow was reorganized under the new `setup` structure, and the dependency injection, tests, and documentation were updated accordingly.

Overall, this change strengthens transfer security, improves session state reliability, and prepares the mobile application for the next phase of protected financial operations based on step-up authentication.


## 2026/06/11 - mobile/transactional-password-01

Implementa fluxo guiado de criação de senha transacional e reforça política de renovação de sessão

### Senha transacional

Foi adicionada uma etapa introdutória antes da criação da senha transacional, explicando ao usuário:

* o propósito da senha transacional;
* a diferença entre senha de acesso e senha transacional;
* os cenários em que ela será utilizada;
* recomendações de segurança para definição do código.

O fluxo de navegação pós-login foi atualizado para direcionar usuários sem senha transacional para esta nova etapa de apresentação.

Também foram realizados ajustes na UX do processo de criação:

* alteração dos textos orientativos;
* centralização visual da etapa de definição do PIN;
* revisão da navegação de retorno para evitar saídas inesperadas do fluxo;
* proteção da rota de confirmação, redirecionando para a introdução quando o PIN informado estiver ausente ou inválido.

Além disso, os arquivos relacionados ao processo de criação da senha transacional foram reorganizados em uma estrutura própria de `creation_flow`.

### Renovação de token

O `AuthInterceptor` passou a renovar sessões apenas para erros de autenticação relacionados à expiração ou invalidação do token de acesso.

Foram adicionadas validações para impedir tentativas de refresh em erros de negócio que também retornam HTTP 401, como:

* credenciais inválidas;
* senha transacional inválida;
* erros de step-up authentication;
* token de aplicação inválido.

Também foi implementada uma proteção contra ciclos de renovação infinitos, garantindo que uma requisição já reprocessada não tente executar novo refresh caso continue recebendo erro de autenticação.

### Cache de último acesso

O cache de último login foi ajustado para armazenar o e-mail utilizado na autenticação em vez do CPF.

Com isso:

* o short login passa a reutilizar um identificador compatível com a API atual;
* registros legados contendo CPF são rejeitados automaticamente;
* foi adicionada validação explícita de e-mail durante a recuperação dos dados armazenados.

### Testes

Foram adicionados testes cobrindo:

* exibição e navegação da tela introdutória da senha transacional;
* regras de renovação seletiva de tokens;
* prevenção de múltiplos refreshes consecutivos;
* validação do cache de último login utilizando e-mail;
* rejeição de identificadores legados baseados em CPF.


## 2026/06/10 - refactor/centralize-auth-session

Refactor authentication session handling by centralizing session state in `AppSection`, removing API-driven home access decisions, and aligning mobile navigation with readiness-derived rules.

### API

* Removed `can_access_home` from `GET /auth/session` contract.
* Simplified session readiness payload to expose only objective readiness state.
* Moved responsibility for post-login navigation decisions to clients.
* Updated session use case, handlers, DTO mappings, and tests accordingly.
* Clarified API documentation to state that authorization remains enforced by protected endpoints independently of client navigation.

### Mobile

#### Session Centralization

* Introduced `AppSection` as the single source of truth for the authenticated session.
* Registered `AppSection` as an application singleton.
* Moved cached session ownership from `AuthRepositoryImpl` to `AppSection`.
* Updated login flow to populate `AppSection` after loading the authenticated session.
* Updated logout flow to clear centralized session state.
* Renamed `profile()` to `getAuthSession()` to better reflect behavior.
* Removed repository-level `userProfile` cache access.

#### Readiness Model

* Removed persisted `canAccessHome` field from `ReadinessSession`.
* Added derived readiness rules inside the domain model:

  * `hasActiveTransactionPassword`
  * `canAccessHome`
* Added `copyWith()` support for:

  * `AuthSession`
  * `ReadinessSession`
* Corrected naming typo:

  * `CustommerSession` → `CustomerSession`

#### Transaction Password Synchronization

* Added local session synchronization after successful transaction password creation.
* Implemented `AppSection.markTransactionPasswordAsActive()`.
* Updated `TransactionPasswordRepositoryImpl` to immediately reflect active status in memory without reloading the session.
* Updated transaction password screens and view models to rely on centralized session readiness.

#### Login and Navigation

* Updated login and short-login view models to resolve destinations using `AppSection`.
* Removed dependency on repository-cached session state.
* Changed transaction password confirmation flow to validate readiness through the centralized session before navigating home.
* Prevented automatic navigation when transaction password creation succeeds but readiness does not allow home access.

### Documentation

* Documented removal of `can_access_home` from the API contract.
* Added migration notes to completed auth-session and transaction-password backlogs.
* Expanded internal-transfer step-up backlog with:

  * session centralization strategy;
  * transfer-entry readiness checks;
  * transaction password reuse flow;
  * AppSection synchronization rules;
  * step-up retry behavior;
  * interceptor requirements;
  * logging and security constraints;
  * acceptance criteria and implementation tasks.
* Added a complete execution task breakdown for internal transfer step-up implementation.

### Tests

* Updated API session handler and use case tests for the new contract.
* Added validation that `can_access_home` is no longer present in responses.
* Updated repository, view model, and destination resolution tests to use `AppSection`.
* Updated session fixtures and readiness builders to rely on computed readiness behavior instead of serialized API fields.


## 2026/06/10 - api/zta-mvp-transactional-password-16

This update completes the transactional password step-up flow by strengthening token lifecycle management, hardening request validation, expanding integration coverage, and consolidating the public API contract around operation-based authorization.

1. `step-up token lifecycle`

   * Added migration `000012_step_up_tokens_cleanup`.
   * Introduced `cleanup_step_up_tokens()` database function.
   * Added dedicated indexes for active and consumed token cleanup paths.
   * Configured a daily `pg_cron` job to remove expired operational records.
   * Implemented a 24-hour retention window for both expired and consumed tokens.
   * Preserved short-term auditability while preventing indefinite growth of the `step_up_tokens` table.

2. `step-up authorization flow`

   * Changed token generation flow to sign the JWT before persisting the authorization record.
   * Prevented persistence of unusable step-up authorizations when token signing fails.
   * Added explicit tests covering signing failures and persistence failures.
   * Preserved the guarantee that only valid signed tokens can be stored for future enforcement.

3. `strict JSON request validation`

   * Introduced shared `DecodeJSON()` helper under `internal/shared/http`.

   * Centralized strict request decoding logic.

   * Continued rejecting unknown fields.

   * Added rejection of multiple JSON values in a single request body.

   * Prevented payloads such as:

     ```json
     {"field":"value"}{"extra":"payload"}
     ```

   * Applied the shared decoder to:

     * transaction password creation
     * step-up authorization
     * internal transfer requests

4. `internal transfer protection`

   * Added integration coverage for the complete protected transfer flow.
   * Validated transaction password creation.
   * Validated step-up authorization issuance.
   * Validated protected transfer execution using `X-Step-Up-Token`.
   * Verified single-use token consumption behavior.
   * Verified rejection of reused step-up tokens.
   * Verified compatibility with transfer idempotency semantics.

5. `public step-up contract`

   * Standardized authorization requests around:

     * HTTP method
     * HTTP path
   * Removed public exposure of internal endpoint policy identifiers.
   * Documented canonical uppercase method usage.
   * Added normalization rules for method handling.
   * Clarified that paths are trimmed but never rewritten or normalized by the API.
   * Reinforced the separation between public operation contracts and internal policy resolution.

6. `API documentation`

   * Documented step-up token retention and cleanup behavior.
   * Added `step_up_tokens` table to database documentation.
   * Updated implementation documentation to include operational retention rules.
   * Clarified authorization error behavior for protected transfers.
   * Documented that internal transfers return `STEP_UP_TOKEN_REQUIRED` when the authorization header is missing.
   * Expanded step-up implementation notes to reflect the final operational model.

7. `security delivery layer`

   * Added handler documentation and constructor documentation.
   * Improved endpoint self-documentation around transaction password creation and step-up authorization responsibilities.
   * Increased maintainability of the security module entry points.

8. `mobile and tooling`

   * Added `AppEnv.isDev` and `AppEnv.isStaging`.
   * Enabled contact verification debug token handling in staging environments.
   * Updated the Bruno environment to target the public staging endpoint.

This commit finalizes the operational foundations of transactional-password-based step-up authentication. Authorization tokens are now single-use, retention-managed, strictly validated at the HTTP boundary, and fully covered by end-to-end transfer integration tests, providing a safer and more predictable security model for sensitive banking operations.


## 2026/06/09 - banklab/docker-environment-isolation-02

This update simplifies the BankLab runtime environment by isolating PostgreSQL in Docker while running the Go API directly on the host.

1. `Makefile`

   * Removed API container execution from the main workflow.
   * Changed `docker-up` to start only the selected PostgreSQL container.
   * Updated `setup`, `run`, `reset`, and `staging` to use the PostgreSQL-only Docker flow.
   * Removed production-related Make targets.
   * Changed `api-run` to execute the Go API directly on the host using the selected environment file.
   * Replaced API container stopping with a port-based `api-stop` command for port `8080`.
   * Simplified database connection variables by using `DB_HOST` and `DB_PORT` directly.

2. `docker-compose.yml`

   * Removed the API service definition.
   * Kept only the PostgreSQL service.
   * Changed PostgreSQL port publishing to use `DB_PORT` directly.

3. `infra/docker/api/Dockerfile`

   * Removed the obsolete API Dockerfile, since the API is no longer built or executed as a Docker container.

4. `infra/scripts/ensure-env-files.sh`

   * Removed production environment generation.
   * Updated API environment defaults for host-based execution.
   * Set `DB_HOST=127.0.0.1` and environment-specific `DB_PORT` values.
   * Removed obsolete Docker publishing variables from generated environment files.
   * Added `sha256sum` fallbacks for random value generation on Linux environments.

5. `infra/scripts/update-mobile-env-ip.sh`

   * Added Linux-compatible IP detection using the `ip` command.
   * Preserved macOS fallback behavior through `route` and `ifconfig`.

6. `README.md`

   * Updated setup documentation to describe the new `dev.env` and `staging.env` environment model.
   * Clarified that PostgreSQL runs in Docker and the API runs on the host.

7. `api/README.md`

   * Updated local setup instructions to remove API container references.
   * Removed production startup references.
   * Clarified the staging flow using the Linux Mint host and Cloudflare Tunnel.

8. `api/docs/00-getting_started.md`

   * Updated environment documentation to reflect the PostgreSQL-only Docker setup.
   * Removed production environment references.
   * Replaced Docker-network database settings with host-based PostgreSQL connection settings.
   * Updated development and staging instructions for the current host API workflow.

This commit consolidates the environment strategy around a simpler and more predictable development/staging model: PostgreSQL remains isolated in Docker, while the API runs directly on the host for easier debugging, tunnel exposure, and operational control.


## 2026/06/09 - banklab/docker-environment-isolation-01

This commit introduces Docker-based API execution with explicit environment isolation for development, staging, and production.

1. `.dockerignore`

   * Added Docker build exclusions for Git metadata, generated Flutter artifacts, local environment files, documentation, templates, and tools.
   * Reduced Docker build context size and avoided leaking local configuration files into images.

2. `.gitignore`

   * Added `api/*.env` to prevent API environment files from being committed.
   * Preserved the existing `.pdf` ignore rule with a proper trailing newline.

3. `Makefile`

   * Replaced the single `api/.env` flow with environment-specific files selected through `API_ENV`.
   * Added support for `dev`, `staging`, and `prod` targets.
   * Updated Docker Compose project names to isolate containers by environment.
   * Added `docker-db-up` to start only PostgreSQL.
   * Changed `api-run` to build and run the API as a Docker container.
   * Added `api-run-dev`, `api-run-staging`, and `api-run-prod`.
   * Updated setup, run, reset, and bootstrap flows to use the Dockerized API and selected database environment.
   * Limited mobile IP synchronization to the development environment.

4. `api/.gitignore`

   * Ignored `dev.env`, `staging.env`, and `prod.env`.

5. `api/README.md`

   * Updated execution instructions for Dockerized API containers.
   * Documented explicit environment selection through `ENV_FILE`.
   * Added development and staging URLs.
   * Added staging execution notes for Cloudflare Tunnel.

6. `api/cmd/api/main.go`

   * Passed the debug verification token exposure flag into the contact verification use case.
   * Updated server binding to use both `SERVER_HOST` and `SERVER_PORT`.
   * Improved startup logging to show the resolved server address.

7. `api/docs/00-getting_started.md`

   * Replaced the legacy single API `.env` documentation with isolated `dev`, `staging`, and `prod` environment files.
   * Documented Docker networking, published ports, environment-specific databases, and Cloudflare Tunnel usage.
   * Added documentation for `APP_ENV`, `PUBLIC_BASE_URL`, `EXPOSE_DEBUG_VERIFICATION_TOKEN`, `API_PUBLISHED_HOST`, and `API_PUBLISHED_PORT`.

8. API test setup files

   * Replaced full bootstrap initialization with `bootstrap.RegisterErrors()` in delivery tests.
   * Avoided requiring runtime environment files for tests that only need error registration.

9. `api/internal/auth/application/contact_verification.go`

   * Added configurable exposure of debug verification tokens.
   * Kept verification tokens available internally while omitting them from responses when disabled.

10. `api/internal/auth/application/contact_verification_test.go`

    * Updated tests for the new constructor signature.
    * Added coverage to ensure debug verification tokens are hidden when exposure is disabled.

11. `api/internal/auth/delivery/auth_authorization_integration_test.go`

    * Updated integration test wiring to pass the debug token exposure flag explicitly.

12. `api/internal/bootstrap/bootstrap.go`

    * Added explicit environment loading through `ENV_FILE`.
    * Removed implicit `.env` discovery to prevent accidental cross-environment configuration.
    * Added `APP_ENV`, `PUBLIC_BASE_URL`, `SERVER_HOST`, and debug token exposure configuration.
    * Added validation for accepted environments.
    * Added production safety validation to reject `EXPOSE_DEBUG_VERIFICATION_TOKEN=true`.

13. `api/internal/bootstrap/bootstrap_test.go`

    * Added tests for environment parsing.
    * Added tests for production debug token exposure validation.

14. `docker-compose.yml`

    * Added environment-specific Compose support for PostgreSQL and API services.
    * Added PostgreSQL health checks.
    * Published PostgreSQL only on loopback.
    * Added the Dockerized API service with explicit runtime env file mounting.
    * Used environment-specific published API host and port configuration.

15. `infra/docker/api/Dockerfile`

    * Added a multi-stage Docker build for the Go API.
    * Built a stripped Linux binary.
    * Added a minimal Alpine runtime image with certificates, timezone data, and a non-root user.

16. `infra/scripts/ensure-env-files.sh`

    * Reworked environment generation to create isolated API files for development, staging, and production.
    * Preserved existing secrets and appended only missing values.
    * Migrated legacy `api/.env` values into `api/dev.env` when available.
    * Added environment-specific database names, ports, public base URLs, and API publication settings.
    * Generated mobile environment files using the corresponding API application token.

17. `infra/scripts/update-mobile-env-ip.sh`

    * Restricted host LAN IP synchronization to `mobile/dev.env`.
    * Prevented staging and production mobile environments from being overwritten with local development URLs.

This commit establishes a safer and more reproducible runtime model by isolating development, staging, and production configuration, containerizing the API execution flow, and preventing accidental exposure of development-only verification tokens in production.


## 2026/06/08 - mobile/create_transaction_password-06

Implemented the transaction password registration flow completion and finalized the related backlog items.

### Changes

#### 1. Backlog Organization

**docs/backlogs/mobile/done/011 - cadastro-senha-transacional.md**

* Moved the transaction password registration backlog to the `done` directory.
* Marked the feature backlog as completed.

**docs/backlogs/mobile/done/011 - cadastro-senha-transacional_tasks.md**

* Moved the associated task list to the `done` directory.
* Consolidated completion status for all tasks related to the transaction password registration feature.

#### 2. OTP Registration Flow Improvement

**mobile/lib/ui/pages/register/register_token_page.dart**

* Enabled automatic token submission behavior when the OTP input is fully completed.
* Connected the `OtpInput.onCompleted` callback to the page logic.
* Restored the `_tokenCompleted` method implementation.
* Reused the existing `_tokenChanged` flow to avoid duplication and maintain a single source of token validation logic.
* Improved the registration experience by eliminating the need for additional user interaction after entering the final verification digit.

#### 3. Bruno Environment Update

**tools/bruno/banklab/collections/BankLab API/environments/BankLab Env.yml**

* Updated the default test identifier used in the Bruno environment configuration.
* Replaced the previous UUID with a new reference value for API testing scenarios.

### Conclusion

This commit completes the transaction password registration backlog, improves the OTP verification user experience through automatic completion handling, and updates the Bruno testing environment to align with the current development workflow.


## 2026/06/05 - mobile/create_transaction_password-05

Refactor environment configuration and migrate API testing assets from Postman to Bruno while introducing configurable JWT lifetimes and database connectivity settings.

### Changes

#### 1. Environment Configuration and Bootstrap

* Centralized runtime configuration around `api/.env` as the single source of truth.
* Added support for:

  * `SERVER_PORT`
  * `DB_HOST`
  * `DB_PORT`
  * `DB_NAME`
  * `DB_USER`
  * `DB_PASSWORD`
  * `JWT_ACCESS_TOKEN_DURATION`
  * `JWT_REFRESH_TOKEN_DURATION`
* Introduced typed database configuration structures in bootstrap and database layers.
* Added validation helpers for required environment variables.
* Added duration parsing and validation for JWT lifetime configuration.
* Added default values:

  * Access token: `15m`
  * Refresh session: `168h`
* Improved startup logging with configurable server URL information.

#### 2. Makefile Modernization

* Updated Makefile to load values directly from `api/.env`.
* Replaced hardcoded PostgreSQL connection details with configurable variables.
* Unified Docker Compose, migrations, database reset, schema export, and API startup around the same environment source.
* Added `env-init` dependency to critical targets to guarantee environment readiness.
* Parameterized:

  * database connection URL
  * database reset commands
  * readiness checks
  * schema export operations
* Updated Docker Compose invocation to use `--env-file api/.env`.

#### 3. Database Initialization Refactor

* Refactored PostgreSQL pool creation to receive an explicit configuration object.
* Removed hardcoded connection strings from application code.
* Added URL-safe PostgreSQL connection string builder.
* Improved separation between runtime configuration and infrastructure concerns.

#### 4. Configurable Session Lifetimes

* Removed hardcoded refresh session expiration values from authentication use cases.
* Added configurable refresh session TTL support for:

  * Login flow
  * Refresh token rotation flow
* Introduced fluent configuration methods:

  * `WithRefreshSessionTTL(...)`
* Connected session expiration to environment-driven JWT configuration.

#### 5. API Startup and Routing Cleanup

* Reorganized route registration into dedicated helper functions.
* Improved startup configuration flow.
* Replaced hardcoded port binding with configurable server port.
* Updated server logging to reflect runtime configuration.
* Preserved existing route structure while improving composition readability.

#### 6. Docker and Local Development Improvements

* Updated Docker Compose PostgreSQL configuration to consume database settings from environment variables.
* Enabled configurable database port mapping.
* Improved local environment portability across machines and contributors.

#### 7. Environment Automation Enhancements

* Expanded `ensure-env-files.sh` to automatically provision:

  * JWT duration variables
  * database configuration variables
  * server port configuration
* Added support for appending missing values into existing environment files without overwriting user configuration.
* Improved bootstrap experience for new contributors.

#### 8. Migration from Postman to Bruno

* Replaced all references to Postman across:

  * README files
  * API documentation
  * contributor documentation
  * roadmap
  * architecture presentations
* Added dedicated Bruno documentation.
* Introduced versionable Bruno workspace structure.
* Added:

  * workspace configuration
  * collections
  * environments
  * authentication requests
  * admin approval request examples
* Removed legacy Postman collections and environment files.
* Removed automatic Postman environment synchronization from development scripts.

#### 9. Documentation Updates

* Updated repository documentation in both Portuguese and English.
* Documented:

  * environment configuration strategy
  * JWT lifetime configuration
  * database configuration variables
  * Docker Compose integration
  * Bruno usage and setup
* Updated onboarding, authentication, infrastructure, and API reference documentation to reflect the new runtime model.

### Result

The project now uses a unified environment-driven configuration model, supports configurable authentication lifetimes and database connectivity, removes hardcoded infrastructure values, and standardizes API testing workflows around Bruno instead of Postman. The setup process is more portable, reproducible, and contributor-friendly while reducing configuration drift between the API, Docker environment, and development tooling.


## 2026/06/03 - mobile/create_transaction_password-04

Implemented the complete transactional password onboarding flow, including post-login gating, route integration, backend error mapping, dependency registration, and comprehensive test coverage.

### Documentation Updates

1. **docs/backlogs/mobile/011 - cadastro-senha-transacional.md**

   * Clarified the security model for transactional password handling.
   * Updated guidance to allow transient PIN transport via `GoRouter.extra` between creation and confirmation pages.
   * Refined restrictions around route state, storage, cache, analytics, logs, and long-lived memory.
   * Standardized UI decision rules around typed `AppErrorCode`.
   * Added explicit mapping of `TRANSACTION_PASSWORD_ALREADY_SET` to `AppErrorCode.transactionPasswordAlreadySet`.

2. **docs/backlogs/mobile/011 - cadastro-senha-transacional_tasks.md**

   * Added implementation results for dependency registration and route integration.
   * Documented post-login gate behavior.
   * Documented transient PIN transport strategy.
   * Added completion notes for tests, error mapping, and updated security documentation.

### Error Handling Improvements

3. **mobile/lib/core/result/errors/**

   * Added `BackendErrorCode` helper for extracting backend error codes from nested error payloads.
   * Exported the helper through `app_error.dart`.
   * Added new `AppErrorCode.transactionPasswordAlreadySet`.

### Routing and Navigation

4. **mobile/lib/core/routing/**

   * Added `TransactionPasswordRoutes` enum.
   * Registered transactional password routes in the application router.
   * Implemented dedicated route definitions for:

     * Transaction password creation.
     * Transaction password confirmation.
   * Added route validation to ensure only valid six-digit PINs reach the confirmation page.

### Transaction Password Feature

5. **mobile/lib/data/repositories/transaction_password/**

   * Added repository contract.
   * Added repository implementation.
   * Registered repository in dependency injection.

6. **mobile/lib/data/services/apis/transaction_password/**

   * Added backend-to-app error code mapping.
   * Preserved backend error metadata inside `AppError.details`.
   * Mapped `TRANSACTION_PASSWORD_ALREADY_SET` to a dedicated typed error.

7. **mobile/lib/ui/pages/auth/transaction_password/**

   * Added transactional password creation page.
   * Added transactional password confirmation page.
   * Added transactional password view model.
   * Implemented PIN confirmation validation.
   * Added handling for already-configured transactional passwords.
   * Cleared sensitive in-memory state after successful completion and special handling flows.
   * Registered the view model in dependency injection.

### Post-Login Readiness Gate

8. **mobile/lib/ui/pages/auth/models/post_login_destination.dart**

   * Added post-login destination resolution model.
   * Introduced navigation outcomes:

     * `home`
     * `transactionPassword`
     * `blocked`
     * `sessionError`
   * Centralized readiness evaluation based on session state.

9. **mobile/lib/ui/pages/auth/login/**

   * Added post-login destination evaluation after successful authentication.
   * Redirects users without a transactional password to the creation flow.
   * Added dedicated blocked/session error handling messages.
   * Replaced local email validation regex with shared string extension validation.
   * Improved Portuguese accentuation in validation messages.

10. **mobile/lib/ui/pages/auth/short_login/**

    * Implemented the same post-login readiness gate behavior used by the standard login flow.
    * Added destination resolution through the shared resolver.

### Dependency Registration

11. **mobile/lib/data/repositories.dart**

    * Registered `TransactionPasswordRepository`.

12. **mobile/lib/ui/viewmodels.dart**

    * Registered `TransactionPasswordViewModel`.

13. **mobile/lib/core/routing/router.dart**

    * Registered transactional password routes.

### Test Coverage

14. **mobile/test/**

    * Added tests for backend error code extraction.
    * Added tests for route extra encoding/decoding.
    * Added repository tests for transactional password creation.
    * Added API tests validating backend error mapping.
    * Added post-login destination resolution tests.
    * Added transactional password view model tests.
    * Verified preservation of backend error codes across layers.
    * Verified mapping of `TRANSACTION_PASSWORD_ALREADY_SET` to typed application errors.

### Conclusion

This commit completes the first version of the transactional password onboarding flow, integrating secure PIN creation, confirmation, post-login readiness gating, typed backend error handling, dependency registration, route orchestration, and automated test coverage while maintaining the requirement that transactional passwords remain transient and are never persisted in application storage.


## 2026/06/03 - mobile/create_transaction_password-03

Implemented the complete transactional password onboarding flow, including post-login gating, dedicated creation and confirmation screens, backend error mapping, dependency registration, routing integration, and comprehensive test coverage.

### Documentation

1. `docs/backlogs/mobile/011 - cadastro-senha-transacional.md`

   * Updated the security guidance to allow transient PIN transport through `GoRouter.extra`.
   * Clarified that transactional passwords must not be persisted in storage, cache, analytics, logs, or long-lived application state.
   * Documented the use of typed `AppErrorCode` values for flow decisions.
   * Added the mapping requirement from `TRANSACTION_PASSWORD_ALREADY_SET` to `AppErrorCode.transactionPasswordAlreadySet`.

2. `docs/backlogs/mobile/011 - cadastro-senha-transacional_tasks.md`

   * Documented the expected results for dependency registration, routing, and post-login gate integration.
   * Added implementation outcomes for testing, backend error handling, and secure PIN transport.

### Error Handling

3. `mobile/lib/core/result/errors/app_error.dart`

   * Exported backend error helper utilities.

4. `mobile/lib/core/result/errors/app_error_code.dart`

   * Added `transactionPasswordAlreadySet` to the application error catalog.

5. `mobile/lib/core/result/errors/backend_error_code.dart`

   * Introduced a reusable helper to extract backend error codes from nested API payload structures.
   * Supports direct, nested `error`, and nested `details` error formats.

### Routing

6. `mobile/lib/core/routing/router.dart`

   * Registered transactional password routes in the application router.

7. `mobile/lib/core/routing/routes.dart`

   * Added `TransactionPasswordRoutes`.
   * Applied minor enum formatting cleanup.

8. `mobile/lib/core/routing/routes/transaction_password_routes.dart`

   * Added creation and confirmation routes.
   * Validated route extras before rendering the confirmation page.
   * Redirected invalid navigation attempts back to the creation screen.
   * Injected the transactional password view model through DI.

### Dependency Registration

9. `mobile/lib/data/repositories.dart`

   * Registered `TransactionPasswordRepository`.

10. `mobile/lib/ui/viewmodels.dart`

    * Registered `TransactionPasswordViewModel`.

### Transaction Password Data Layer

11. `mobile/lib/data/repositories/transaction_password/transaction_password_repository.dart`

    * Added repository contract for transactional password creation.

12. `mobile/lib/data/repositories/transaction_password/transaction_password_repository_impl.dart`

    * Implemented repository delegation to the API layer.

13. `mobile/lib/data/services/apis/transaction_password/transaction_password_api.dart`

    * Added backend error code preservation.
    * Mapped `TRANSACTION_PASSWORD_ALREADY_SET` to a dedicated application error code.
    * Preserved backend error metadata through `details`.

### Post-Login Gate

14. `mobile/lib/ui/pages/auth/models/post_login_destination.dart`

    * Added centralized post-login destination resolution.
    * Implemented routing decisions based on readiness and transactional password status.

15. `mobile/lib/ui/pages/auth/login/viewmodel/login_viewmodel.dart`

    * Added destination resolution support.

16. `mobile/lib/ui/pages/auth/short_login/viewmodel/short_login_viewmodel.dart`

    * Added destination resolution support.
    * Stored repository reference for post-login evaluation.

17. `mobile/lib/ui/pages/auth/login/login_page.dart`

    * Replaced direct home navigation with post-login gate evaluation.
    * Added transactional password onboarding redirection.
    * Added blocked and session error handling.
    * Replaced local email regex validation with shared string extension validation.
    * Improved validation messages.

18. `mobile/lib/ui/pages/auth/short_login/short_login_page.dart`

    * Applied the same post-login gate behavior used in the full login flow.
    * Added transactional password onboarding routing and failure handling.

### Transaction Password UI

19. `mobile/lib/ui/pages/auth/transaction_password/create_transaction_password_page.dart`

    * Added PIN creation screen.
    * Implemented secure local PIN handling.
    * Navigated to confirmation using transient `GoRouter.extra`.

20. `mobile/lib/ui/pages/auth/transaction_password/confirm_transaction_password_page.dart`

    * Added PIN confirmation screen.
    * Implemented confirmation validation and mismatch handling.
    * Executed transactional password creation through the view model.
    * Handled already-configured password scenarios.
    * Cleared sensitive state after successful completion.

21. `mobile/lib/ui/pages/auth/transaction_password/viewmodel/transaction_password_viewmodel.dart`

    * Added command-based transactional password creation workflow.

### Tests

22. `mobile/test/core/result/errors/backend_error_code_test.dart`

    * Added coverage for backend error extraction logic.

23. `mobile/test/core/routing/extra_codec_test.dart`

    * Added validation for string extra encoding and decoding.

24. `mobile/test/data/repositories/transaction_password/transaction_password_repository_impl_test.dart`

    * Added repository success and failure propagation tests.
    * Verified backend error code preservation.

25. `mobile/test/data/services/apis/transaction_password/transaction_password_api_test.dart`

    * Added validation for mapping backend transaction password errors into typed application errors.

26. `mobile/test/ui/pages/auth/post_login_destination_test.dart`

    * Added coverage for all post-login destination scenarios.
    * Verified identical behavior between login and short-login flows.

27. `mobile/test/ui/pages/auth/transaction_password/viewmodel/transaction_password_viewmodel_test.dart`

    * Added view model success and failure tests.

### Conclusion

This update completes the first version of the transactional password onboarding flow, introducing secure PIN creation and confirmation screens, centralized post-login gating, typed backend error handling, dependency and routing integration, and a dedicated test suite to validate both navigation and error-processing behavior.


## 2026/06/03 - mobile/create_transaction_password-02

Start the mobile transaction password creation flow by documenting the implementation plan and adding the initial API integration layer.

1. `api/docs/07-api-rest.md`

   * Clarified that successful transaction password creation currently returns `status = active`.
   * Documented that `blocked` exists in the domain for future validation and step-up flows.

2. `docs/backlogs/README.md`

   * Added the task document for the mobile transaction password registration backlog.

3. `docs/backlogs/mobile/011 - cadastro-senha-transacional.md`

   * Refined the backlog to use the existing authenticated session snapshot from `GET /auth/session`.
   * Removed the need for a dedicated transaction password status endpoint.
   * Standardized naming around `transaction_password`.
   * Defined the initial implementation as view model + repository, without a dedicated use case.
   * Added decisions for post-login gating, two-page PIN flow, API error-code extraction, and test coverage.

4. `docs/backlogs/mobile/011 - cadastro-senha-transacional_tasks.md`

   * Added a complete 9-task implementation plan for the transaction password creation flow.
   * Covered DTOs, API service, repository, post-login gate helper, UI flow, dependency registration, tests, and documentation.

5. `mobile/lib/data/repositories/auth/auth_repository_impl.dart`

   * Updated the auth session loading call from `getProfile()` to `getAuthSession()`.
   * Adjusted imports to use the project absolute import style.

6. `mobile/lib/data/services/apis/auth/auth_api.dart`

   * Renamed `getProfile()` to `getAuthSession()` to better represent the `/auth/session` contract.
   * Updated log labels accordingly.
   * Adjusted imports to the project absolute import style.

7. `mobile/lib/data/services/apis/account/statement_api.dart`

   * Removed redundant documentation comments from the API service.

8. `mobile/lib/data/services/apis/transaction_password/`

   * Added `CreateTransactionPasswordRequestDto`.
   * Added `TransactionPasswordStatusResponseDto`.
   * Added `TransactionPasswordStatus` enum.
   * Added `TransactionPasswordApi` with support for `POST /security/transaction-password`.

9. `mobile/lib/data/services/services.dart`

   * Registered `TransactionPasswordApi` in the service container.

10. `mobile/test/data/repositories/auth/auth_repository_impl_test.dart`

    * Updated the fake auth API to match the new `getAuthSession()` method name.

11. `mobile/test/data/services/apis/auth/auth_api_get_profile_test.dart`

    * Updated the auth API test to call `getAuthSession()`.

12. `mobile/test/data/services/apis/transaction_password/`

    * Added DTO tests for request serialization and response parsing.
    * Added API service tests for success, envelope error, HTTP error, and client failure scenarios.

This commit establishes the first technical base for mobile transaction password creation, aligning documentation, backlog planning, API contracts, service registration, and tests.


## 2026/06/03 - mobile/create_transaction_password-01

Prepare mobile authentication session model and token input foundation.

1. `docs/backlogs/api`

   * Moved the completed auth session bootstrap backlog and task documents to the `done` folder.

2. `mobile/ios`

   * Removed CocoaPods integration references from Flutter iOS build configuration.
   * Deleted the iOS `Podfile`.
   * Cleaned Pods-related framework, group, build phase, and xcconfig references from the Xcode project.

3. `mobile/lib/data/repositories/auth`

   * Replaced `UserProfile` usage with the new `AuthSession` domain model.
   * Updated cached last login identity to use customer name and CPF from the session payload.
   * Adjusted profile caching to store the full authentication session.

4. `mobile/lib/data/services/apis/auth`

   * Changed profile loading to consume the new `GET /auth/session` endpoint.
   * Removed the previous two-step `/customers/me` and `/auth/me` composition.
   * Added direct parsing of the canonical API envelope into `AuthSession`.

5. `mobile/lib/domain/common/auth/models/auth_session`

   * Added `AuthSession`, `UserSession`, `CustommerSession`, and `ReadinessSession`.
   * Added `TransactionPasswordStatus` enum to represent transaction password readiness.
   * Exported the session-related models through the main auth session file.

6. `mobile/lib/domain/common/auth/models`

   * Removed the old `UserProfile` model.

7. `mobile/lib/ui/components/input_text`

   * Added `TokenInput`, a reusable numeric token/password input component with fixed-length cells, autofocus, visibility control, initial value support, and completion callbacks.

8. `mobile/test/data/repositories/auth`

   * Updated auth repository tests to validate the new `AuthSession` model and CPF-based cache identity.

9. `mobile/test/data/services/apis/auth`

   * Added tests covering the `/auth/session` contract parsing and `AuthApi.getProfile()` behavior.

This commit aligns the mobile client with the new auth session contract and introduces the first UI foundation required for transaction password flows.


## 2026/06/02 - api/auth-session-bootstrap-01

Implement authenticated session bootstrap endpoint for post-login API readiness.

1. `api/cmd/api/main.go`

   * Wired `GetSessionUseCase` into the auth module dependencies.
   * Registered `GET /auth/session` as a JWT-protected route.
   * Passed the session use case into the auth handler constructor.

2. `api/cmd/api/routes_test.go`

   * Added route coverage to ensure `/auth/session` is wrapped by the authentication middleware.

3. `api/internal/auth/application/get_session.go`

   * Added `GetSessionUseCase`.
   * Implemented authenticated user resolution from context.
   * Loaded user, linked customer, operational accounts, and transaction password state.
   * Calculated readiness fields for post-login routing.
   * Added transaction password session states: `active`, `not_set`, `locked`, and `unknown`.

4. `api/internal/auth/application/get_session_test.go`

   * Added application tests for ready sessions, missing transaction password, missing auth context, invalid customer link, and missing customer profile.

5. `api/internal/auth/delivery/handler.go`

   * Added `Session` handler.
   * Added response DTOs for `user`, `customer`, and `readiness`.
   * Ensured sensitive fields such as `user.customer_id`, `customer.email`, and credential material are not returned.
   * Added date formatting for customer birth date.

6. `api/internal/auth/delivery/handler_test.go`

   * Added delivery tests for successful session response.
   * Added unauthorized session response coverage.
   * Validated response shape and excluded fields.

7. `api/docs/07-api-rest.md`

   * Documented `GET /auth/session`.
   * Updated authentication requirements and endpoint index.
   * Added success and error response examples.
   * Clarified that `/auth/me` remains available for compatibility.

8. `tools/postman/Banklab_API.postman_collection.json`

   * Added Postman request for `GET /auth/session` using bearer authentication.

9. `docs/backlogs/api/009 - auth-session-bootstrap_tasks.md`

   * Marked all API session bootstrap tasks as completed.
   * Added final alignment notes and validation checklist.

10. `docs/backlogs/mobile/011 - cadastro-senha-transacional.md`

* Updated the mobile backlog to use `GET /auth/session` as the post-login gate source.
* Removed pending API-contract assumptions and obsolete `next_required_step` references.

This change completes the API-side session bootstrap contract and provides the mobile client with a single canonical endpoint for post-login readiness decisions.


## 2026/06/02 - backlog/api-009-auth-session-bootstrap

Add backlog documentation for the authenticated post-login session bootstrap flow.

1. `docs/backlogs/README.md`

   * Replaced completed API backlog entries with the new `009 - auth-session-bootstrap` backlog and tasks.
   * Split the previous mobile transactional password and step-up backlog into two independent entries:

     * transactional password registration;
     * internal transfer step-up authorization.

2. `docs/backlogs/api/009 - auth-session-bootstrap.md`

   * Added the parent backlog for `GET /auth/session`.
   * Defined the endpoint objective, response contract, readiness fields, compatibility strategy, errors, out-of-scope items, acceptance criteria, and implementation references.
   * Established `GET /auth/session` as the canonical authenticated bootstrap endpoint for post-login clients.

3. `docs/backlogs/api/009 - auth-session-bootstrap_tasks.md`

   * Added seven implementation tasks covering DTOs, use case, readiness calculation, HTTP handler, tests, REST documentation, and final API validation.

4. `docs/backlogs/api/done/007 - public-step-up-endpoint-contract.md`

   * Updated references to point to the renamed mobile step-up backlog.

5. `docs/backlogs/api/done/007 - public-step-up-endpoint-contract_tasks.md`

   * Updated documentation task references to the new mobile step-up backlog name.

6. `docs/backlogs/api/done/008 - transaction-password-pepper.md`

   * Updated references to include both the transactional password registration backlog and the internal transfer step-up backlog.

7. `docs/backlogs/mobile/011 - cadastro-senha-transacional.md`

   * Added a dedicated mobile backlog for transactional password registration.
   * Defined the post-login gate, PIN registration flow, API contracts, security constraints, architecture impacts, UX requirements, and acceptance criteria.

8. `docs/backlogs/mobile/012 - step-up-transferencia-interna.md`

   * Renamed and refocused the previous combined backlog.
   * Removed transactional password registration from its scope.
   * Kept the backlog focused on step-up authorization for internal transfers.
   * Added retry handling for expired step-up tokens and routing behavior when the transactional password is not set.

9. `api/docs/images/database.png`

   * Updated the database diagram image.

10. `tools/postman/Environment.postman_environment.json`

* Updated the local `base_url` from `192.168.0.21` to `192.168.0.17`.

This commit separates transactional password registration from transfer step-up authorization and introduces the API session bootstrap backlog required to support a cleaner post-login mobile flow.


## 2026/05/30 — api/transaction-password-papper-01

Implements transaction password pepper hardening in the API, adding a required server-side secret to strengthen PIN hashing while preserving the public mobile contract.

### 1. Configuration

* Added `TRANSACTION_PASSWORD_PEPPER` to API bootstrap configuration.
* Added fail-fast validation for:

  * missing pepper;
  * pepper shorter than 32 characters;
  * pepper equal to `APP_TOKEN`;
  * pepper equal to `JWT_SECRET`.

### 2. Transaction Password Hashing

* Updated `BcryptTransactionPasswordHasher` to require a pepper.
* Added HMAC-SHA256 preprocessing before bcrypt.
* Encoded the HMAC output as base64 before hashing.
* Updated both `Hash` and `Compare` to use the pepper-derived value.
* Kept the persisted value as a valid bcrypt hash.

### 3. API Wiring

* Updated `cmd/api/main.go` to inject `config.TransactionPasswordPepper` into the transaction password hasher.

### 4. Tests

* Added unit tests for the transaction password hasher.
* Covered:

  * empty pepper validation;
  * successful hash and compare;
  * wrong password failure;
  * different pepper failure;
  * valid bcrypt hash generation.

### 5. Documentation and Setup

* Updated API README and getting started documentation with `TRANSACTION_PASSWORD_PEPPER`.
* Added guidance for generating the pepper with `openssl rand -base64 32`.
* Documented that pepper rotation invalidates existing transaction password hashes without a migration strategy.
* Updated `ensure-env-files.sh` to generate and persist a local pepper for new API `.env` files.

### 6. Backlog Organization

* Moved the public step-up endpoint contract backlog to `done`.
* Added the completed transaction password pepper backlog and task documentation under `docs/backlogs/api/done`.

The API now requires a dedicated transaction password pepper, applies it before bcrypt hashing, and is ready to support the mobile transaction password flow without changing public request contracts.


## 2026/05/29 - api/zta-mvp-transactional-password-15

This commit completes the migration of the step-up authorization public contract from internal endpoint keys to public HTTP operations (`method` + `path`), strengthening API encapsulation and preventing clients from depending on internal policy identifiers.

### 1. Public Step-Up Contract Refactoring

#### `api/internal/security/application/authorize_step_up.go`

* Replaced `EndpointKey` input with:

  * `Method`
  * `Path`
* Introduced public operation validation before authorization.
* Added resolution layer that converts public operations into internal policy keys.
* Kept internal endpoint keys exclusively inside the backend and signed step-up tokens.

#### `api/internal/security/delivery/handler.go`

* Updated request contract:

  * removed `endpoint_key`
  * added `method`
  * added `path`
* Updated request mapping and normalization logic.

#### `api/internal/security/domain/interfaces.go`

* Added `StepUpPublicOperationResolver` contract.

#### `api/cmd/api/main.go`

* Replaced endpoint policy registration with the new public operation resolver.

### 2. Public HTTP Operation Domain Model

#### `api/internal/security/domain/public_http_operation.go`

* Added `PublicHTTPOperation` value object.
* Implemented normalization and validation for:

  * HTTP methods
  * public paths
* Added canonical lookup key generation.

Validation rules include:

* method must be non-empty and uppercase
* path must start with `/`
* no scheme (`http://`, `https://`)
* no host
* no query string
* no fragment

#### `api/internal/security/domain/errors.go`

* Added:

  * `ErrInvalidStepUpPublicOperation`
  * `ErrInvalidStepUpPublicOperationMethod`
  * `ErrInvalidStepUpPublicOperationPath`

### 3. Public Operation Resolver

#### `api/internal/security/domain/step_up_public_operation_resolver.go`

* Added whitelist-based resolver implementation.
* Added operation-to-policy mapping model.
* Implemented default mapping:

```text
POST /accounts/internal-transfers
    ->
internal_transfer.create
```

Benefits:

* clients only know public API contracts
* backend remains free to rename internal policies
* future ZTA policies remain decoupled from client implementations

### 4. Error Mapping Improvements

#### `api/internal/security/application/errors_registry.go`

* Registered new validation errors.
* Mapped invalid public operation inputs to:

  * `INVALID_DATA`
  * HTTP 400

This preserves consistency with the API error model. 

### 5. Test Coverage Expansion

#### `api/internal/security/application/authorize_step_up_test.go`

* Migrated all tests from endpoint-key inputs to public operation inputs.
* Added coverage for:

  * invalid methods
  * invalid paths
  * whitelist violations
  * successful operation resolution

#### `api/internal/security/domain/public_http_operation_test.go`

* Added comprehensive validation coverage for:

  * valid operations
  * empty methods
  * invalid paths
  * URLs containing hosts
  * query strings
  * fragments
  * templated routes

#### `api/internal/security/domain/step_up_public_operation_resolver_test.go`

* Added resolver coverage for:

  * successful mappings
  * unsupported operations
  * invalid operations
  * templated route mappings
  * nil resolver handling

#### `api/internal/security/delivery/handler_test.go`

* Updated request payload tests.
* Added coverage for:

  * invalid method errors
  * invalid path errors
  * legacy payload rejection
  * successful public contract handling

### 6. Documentation Alignment

#### Updated files

* `README.md`
* `README_en.md`
* `api/README.md`
* `api/docs/05-error_and_response.md`
* `api/docs/07-api-rest.md`
* `api/docs/implementations/03-zta-step-up-transaction-password.md`
* `mobile/README.md`
* `mobile/docs/01-implemented-features.md`
* `docs/backlogs/mobile/011 - senha-transacional-e-step-up.md`

Key changes:

* Replaced references to:

  * `endpoint_key`
  * `internal_transfer.create`
* Standardized documentation around:

  * `method`
  * `path`
* Clarified that internal policy keys are backend implementation details.
* Updated error descriptions to refer to public operations instead of logical endpoints.
* Updated transfer documentation to require step-up authorization for:

```text
POST /accounts/internal-transfers
```

### 7. Backlog Completion

#### `docs/backlogs/api/007 - public-step-up-endpoint-contract_tasks.md`

* Marked all six tasks as completed.
* Added final validation record documenting:

  * implementation status
  * testing results
  * documentation synchronization
  * whitelist configuration
  * enforcement behavior

### Conclusion

This change finalizes the first public-contract version of the step-up authorization flow. Clients now request authorization using only public HTTP operations, while internal policy identifiers remain hidden inside the backend. The implementation introduces a dedicated operation model, whitelist resolver, expanded validation, complete test coverage, and synchronized documentation, providing a cleaner and more maintainable foundation for future Zero Trust and step-up authorization policies.


## 2026/05/29 — api/zta-mvp-transactional-password-14

This commit reorganizes the ZTA backlog documentation and defines the next API and mobile work around the public step-up contract and transactional password flow.

1. `docs/backlogs/README.md`

   * Updated the active API backlog list to replace the completed ZTA MVP backlog entries with the new `007` public step-up endpoint contract backlog.
   * Added the new mobile backlog for transactional password and step-up.
   * Moved completed backlog references out of the active list.

2. `docs/backlogs/api/007 - public-step-up-endpoint-contract.md`

   * Added a new API backlog defining the migration from public `endpoint_key` input to a public HTTP operation contract based on `method` and `path`.
   * Documented the decision that mobile clients must not know internal policy keys such as `internal_transfer.create`.
   * Defined validation rules for public paths, allowed operations, compatibility decisions, expected errors, acceptance criteria, and out-of-scope items.

3. `docs/backlogs/api/007 - public-step-up-endpoint-contract_tasks.md`

   * Added six implementation tasks for the public step-up contract.
   * Covered the public HTTP operation model, operation-to-policy resolver, delivery contract changes, tests, documentation updates, and final alignment checks.

4. `docs/backlogs/api/done/*`

   * Moved the completed API ZTA MVP backlog files from the active API backlog folder to `done`.
   * Preserved file contents while marking the previous ZTA foundation, transaction password, step-up token, enforcement, contracts, and generated migration backlog as completed work.

5. `docs/backlogs/mobile/011 - senha-transacional-e-step-up.md`

   * Added a mobile backlog for transactional password creation and step-up authorization.
   * Documented API contracts, mobile product flow, architecture impacts, repository/use case responsibilities, UI requirements, ZTA error handling, security constraints, and acceptance criteria.

6. `docs/backlogs/mobile/done/*`

   * Moved the completed multi-page registration mobile backlog and tasks to `done`.

This update separates completed ZTA MVP planning from the next contract refinement stage, clarifies that public clients should use HTTP method and path instead of internal policy keys, and prepares the mobile implementation backlog for transactional password and protected internal transfers.


## 2026/05/29 - api/zta-mvp-transactional-password-13

Finalized the ZTA MVP contract documentation, consolidating step-up token behavior, protected endpoint requirements, error semantics, and client integration guidance across API, implementation, backlog, and mobile documentation.

### Documentation Updates

#### 1. `api/docs/05-error_and_response.md`

* Expanded the ZTA domain error mapping table.
* Added `TRANSACTION_PASSWORD_ALREADY_SET` to the documented error catalog.
* Clarified that clients must depend on `error.code` instead of `error.message`.
* Distinguished issuance/authorization errors from protected endpoint enforcement errors.
* Documented the exact meaning of `STEP_UP_ENDPOINT_NOT_ALLOWED`.
* Clarified all conditions covered by `STEP_UP_TOKEN_INVALID`, including:

  * malformed tokens;
  * invalid signatures;
  * missing required claims;
  * invalid scope;
  * missing persisted `jti`.

#### 2. `api/docs/07-api-rest.md`

* Added explicit request header requirements for protected internal transfers:

  * `Authorization: Bearer <access_token>`
  * `X-Step-Up-Token: <step_up_token>`
* Clarified the distinction between:

  * step-up authorization failures (`STEP_UP_ENDPOINT_NOT_ALLOWED`);
  * protected endpoint enforcement failures.
* Refined the definition of `STEP_UP_TOKEN_INVALID` to match implementation behavior.

#### 3. `api/docs/implementations/03-zta-step-up-transaction-password.md`

* Updated the implementation guide with the finalized ZTA error contract.
* Added guidance that consumers must rely on `error.code`.
* Documented issuance versus enforcement responsibilities.
* Clarified all scenarios covered by `STEP_UP_TOKEN_INVALID`.
* Improved consistency between implementation documentation and runtime behavior.

### Backlog Closure

#### 4. `docs/backlogs/api/006d - zta-contracts-and-docs.md`

* Added final validation notes for the ZTA MVP.
* Recorded successful validation through `go test ./...`.
* Confirmed alignment between:

  * error contracts;
  * JWT signer/verifier;
  * enforcement layer;
  * internal transfer handler;
  * published documentation.
* Marked the backlog as complete for MVP scope.
* Explicitly documented items intentionally left outside MVP scope:

  * payload-bound step-up authorization;
  * broader endpoint coverage;
  * trusted devices;
  * local biometrics;
  * liveness verification;
  * risk signals;
  * Postman collection requirements.

#### 5. `docs/backlogs/api/006d - zta-contracts-and-docs_tasks.md`

* Marked Tasks 3 through 6 as completed:

  * Consolidate ZTA error contract.
  * Update REST documentation.
  * Update ZTA documentation and references.
  * Perform final contract validation.

### Mobile Documentation Alignment

#### 6. `mobile/README.md`

* Updated transfer feature documentation to reflect the current security contract.
* Documented that internal transfers now require:

  * `POST /security/step-up/authorize`
  * `endpoint_key=internal_transfer.create`
  * `X-Step-Up-Token` in transfer requests.

#### 7. `mobile/docs/01-implemented-features.md`

* Added API contract notes for transfer operations.
* Documented:

  * mandatory `X-Step-Up-Token`;
  * issuance through the step-up authorization endpoint;
  * single-use token behavior;
  * retry requirements after `STEP_UP_TOKEN_CONSUMED`;
  * interaction between idempotency retries and step-up token renewal.

### Conclusion

This commit completes the documentation phase of the ZTA transactional password MVP by establishing a stable and fully documented contract for step-up authorization, protected endpoint enforcement, error handling, transfer requirements, and client integration expectations across both API and mobile layers.


## 2026/05/29 — api/zta-mvp-transactional-password-12

This commit consolidates the ZTA MVP step-up enforcement contract for internal transfers, aligning documentation, error semantics, retry behavior, and automated test coverage.

1. `README.md`

   * Documented that `POST /accounts/internal-transfers` requires `X-Step-Up-Token`.
   * Clarified that the token must be issued by `POST /security/step-up/authorize` for `internal_transfer.create`.

2. `README_en.md`

   * Added the same internal transfer step-up requirement to the English README.
   * Kept the authentication summary aligned with the Portuguese README.

3. `api/README.md`

   * Updated the API authentication summary to mention the mandatory `X-Step-Up-Token` for protected internal transfers.

4. `api/docs/07-api-rest.md`

   * Updated the standard error envelope to include optional `error.details`.
   * Clarified that clients must depend on `error.code`, not `error.message`.
   * Documented that the step-up token is an `HS256` JWT.
   * Added internal transfer enforcement rules for `X-Step-Up-Token`.
   * Documented atomic token consumption before transfer execution.
   * Clarified retry behavior when reusing consumed step-up tokens.
   * Added ZTA-related error codes and descriptions.
   * Added error scenarios for missing, invalid, expired, consumed, and endpoint-mismatched step-up tokens.

5. `api/docs/implementations/03-zta-step-up-transaction-password.md`

   * Documented that the step-up token must be signed with `HS256`.
   * Added explicit rules for token consumption and retry behavior.
   * Documented expected automated coverage for verifier, enforcement use case, and integration tests.
   * Reformatted the error contract table for readability.

6. `api/internal/account/transaction/delivery/handler_test.go`

   * Added coverage ensuring consumed step-up tokens return `STEP_UP_TOKEN_CONSUMED`.
   * Verified that the transfer use case is not executed when step-up enforcement fails.

7. `api/internal/security/application/enforce_step_up_test.go`

   * Added coverage for retrying enforcement with the same step-up token.
   * Confirmed that the second use of the same persisted `jti` returns `ErrStepUpTokenConsumed`.

8. `api/internal/security/infrastructure/step_up_enforcement_integration_test.go`

   * Added an integration test covering signer, verifier, PostgreSQL repository, and enforcement use case.
   * Verified that a persisted step-up `jti` is consumed exactly once.
   * Confirmed that retrying the same signed token fails with `ErrStepUpTokenConsumed`.

9. `docs/backlogs/api/006d - zta-contracts-and-docs.md`

   * Clarified that the step-up JWT uses `HS256`.
   * Added the MVP enforcement alignment status.
   * Documented the protected endpoint, required header, endpoint key, architecture boundary, single-use behavior, and retry/idempotency relationship.

10. `docs/backlogs/api/006d - zta-contracts-and-docs_tasks.md`

* Added the task breakdown for consolidating ZTA contracts and documentation.
* Defined scope, acceptance criteria, dependencies, and current status for the documentation backlog.

This change closes the main contract gap around ZTA step-up enforcement by making internal transfer protection explicit, tested, and consistently documented across API references, implementation notes, and backlog planning.


## 2026/05/29 — api/zta-mvp-transactional-password-11

Implement step-up enforcement for internal transfers.

This update connects the transactional password step-up authorization flow to the transfer endpoint, ensuring that internal transfers require a valid, endpoint-bound, single-use step-up token before the transfer use case is executed.

1. `api/cmd/api/main.go`

   * Added `JWTStepUpTokenVerifier`.
   * Created `EnforceStepUpUseCase`.
   * Injected step-up enforcement into the transaction handler.

2. `api/internal/account/transaction/delivery/handler.go`

   * Added support for `X-Step-Up-Token`.
   * Added `enforceStepUpUseCase` dependency.
   * Enforced step-up validation before executing transfers.
   * Blocked transfer execution when the step-up token is missing, invalid, expired, consumed, or bound to another endpoint.

3. `api/internal/security/application/enforce_step_up.go`

   * Added `EnforceStepUpUseCase`.
   * Validates authenticated user context.
   * Verifies signed step-up token claims.
   * Checks user ownership, endpoint binding, expiration, and persisted token consistency.
   * Consumes the step-up token by JTI to enforce single-use semantics.

4. `api/internal/security/infrastructure/jwt_step_up_token_verifier.go`

   * Added JWT verifier for step-up tokens.
   * Validates HS256 signature, required claims, scope, user ID, endpoint key, JTI, issued time, and expiration.

5. `api/internal/security/domain/errors.go`

   * Added transaction password required error.
   * Added step-up endpoint mismatch error.

6. `api/internal/shared/errors/codes.go`

   * Added stable error codes for transaction password and step-up enforcement failures.

7. `api/internal/security/application/errors_registry.go`

   * Registered HTTP mappings for the new transaction password and step-up errors.

8. Tests

   * Updated handler constructors across route and integration tests.
   * Added transfer handler tests to ensure failed step-up enforcement prevents transfer execution.
   * Added use case tests for successful enforcement, missing tokens, invalid users, verifier failures, user mismatch, endpoint mismatch, consumed tokens, expired tokens, and persisted record divergence.
   * Added JWT verifier tests for valid tokens, invalid secrets, invalid scopes, missing claims, expired tokens, invalid algorithms, and malformed tokens.
   * Extended error registry tests for the new error mappings.

This change moves internal transfers closer to the ZTA MVP model by requiring explicit, short-lived, endpoint-scoped authorization before executing a sensitive financial operation.


## 2026/05/29 — api/zta-mvp-transactional-password-10

Introduce the initial step-up token verification contract and align the enforcement backlog with the transactional password/ZTA flow.

1. `api/internal/security/domain`

   * Added `ErrStepUpTokenRequired`.
   * Added `StepUpTokenVerifier` contract.
   * Added `VerifiedStepUpTokenClaims` with validation and normalization.
   * Centralized the step-up JWT scope as `StepUpTokenScope`.

2. `api/internal/security/application`

   * Registered `STEP_UP_TOKEN_REQUIRED` as a mapped domain error.
   * Added test coverage for the new error mapping.

3. `api/internal/security/infrastructure`

   * Updated the JWT step-up token signer to use the shared domain scope constant.
   * Adjusted signer tests to validate the shared scope definition.

4. `api/internal/shared/errors`

   * Added the shared `STEP_UP_TOKEN_REQUIRED` error code.

5. `docs/backlogs/api`

   * Expanded the internal transfer step-up enforcement backlog.
   * Added detailed enforcement decisions for JWT validation, persisted `jti` comparison, atomic consumption, retry behavior, and error semantics.
   * Added a dedicated task breakdown for the enforcement implementation.
   * Updated the ZTA contract documentation with the step-up JWT claims and error separation.

This prepares the API for implementing step-up enforcement on internal transfers while keeping JWT verification behind a clean domain contract.


## 2026/05/29 — api/zta-mvp-transactional-password-09

Consolidate transactional password and step-up authorization documentation, aligning API contracts, error mappings, implementation notes, and project references with the current Zero Trust MVP flow.

### Documentation Updates

#### 1. API Contract (`api/docs/07-api-rest.md`)

* Added documentation for `POST /security/step-up/authorize`.
* Defined request and response payloads for step-up authorization.
* Documented step-up token lifetime (`expires_in: 120`).
* Added `STEP_UP_ENDPOINT_NOT_ALLOWED` error contract.
* Clarified the transactional password validation flow used to obtain a step-up token.
* Documented that `POST /accounts/internal-transfers` requires the `X-Step-Up-Token` header.
* Added transfer-related step-up enforcement errors:

  * `STEP_UP_TOKEN_REQUIRED`
  * `STEP_UP_TOKEN_INVALID`
  * `STEP_UP_TOKEN_EXPIRED`
  * `STEP_UP_TOKEN_CONSUMED`
  * `STEP_UP_ENDPOINT_MISMATCH`
  * `TRANSACTION_PASSWORD_REQUIRED`

#### 2. Error and Response Standard (`api/docs/05-error_and_response.md`)

* Expanded client error guidance to include:

  * `401 Unauthorized`
  * `403 Forbidden`
* Added security-domain error mappings for:

  * Transaction password setup and validation
  * Step-up authorization
  * Step-up token enforcement
* Improved HTTP status classification for authentication and authorization failures.

#### 3. Step-Up Implementation Documentation (`api/docs/implementations/03-zta-step-up-transaction-password.md`)

* Added `STEP_UP_ENDPOINT_NOT_ALLOWED` to the documented error contract.
* Updated the conceptual `step_up_tokens` model to include persisted `jti` support.
* Improved alignment between implementation notes and runtime behavior.

#### 4. Backlog Documentation (`docs/backlogs/api/006d - zta-contracts-and-docs.md`)

* Extended the MVP error matrix with `STEP_UP_ENDPOINT_NOT_ALLOWED`.
* Kept backlog specifications synchronized with the implemented contract.

### Project Documentation Improvements

#### 5. Repository Readmes

**README.md**

* Added security endpoints:

  * `POST /security/transaction-password`
  * `POST /security/step-up/authorize`

**README_en.md**

* Added the new security endpoints to the public API list.
* Fixed backlog documentation links using URL-encoded paths.

**api/README.md**

* Added transactional password and step-up authorization to the feature list.
* Added security endpoints to the API route summary.
* Updated migration command examples:

  * `make api-migrate-up` → `make migrate-up`
* Updated test command examples:

  * `make api-test` → `make api-tests`

**mobile/README.md**

* Updated test command reference:

  * `make mobile-test` → `make mobile-tests`
* Removed outdated mention of deposit and withdraw operations from the feature summary.

### Changelog

#### 6. `CHANGELOG.md`

* Added the release entry for `api/zta-mvp-transactional-password-08`.
* Recorded all contract, documentation, and security-flow updates introduced in this iteration.

### Result

This commit completes the documentation consolidation for the transactional password MVP, ensuring that API contracts, error mappings, implementation references, backlog specifications, and project-level documentation consistently reflect the current step-up authorization flow and Zero Trust enforcement model.


## 2026/05/29 — api/zta-mvp-transactional-password-08

Consolidate the ZTA transactional-password MVP contracts and documentation with the implemented step-up flow.

1. `api/docs/07-api-rest.md`

  * Added `POST /security/step-up/authorize` endpoint documentation.
  * Documented step-up request/response payload and lifetime (`expires_in: 120`).
  * Added `STEP_UP_ENDPOINT_NOT_ALLOWED` error to step-up authorization.
  * Clarified that `POST /accounts/internal-transfers` requires `X-Step-Up-Token`.
  * Added step-up enforcement errors for internal transfers (`STEP_UP_TOKEN_REQUIRED`, `STEP_UP_TOKEN_INVALID`, `STEP_UP_TOKEN_EXPIRED`, `STEP_UP_TOKEN_CONSUMED`, `STEP_UP_ENDPOINT_MISMATCH`, `TRANSACTION_PASSWORD_REQUIRED`).

2. `api/docs/implementations/03-zta-step-up-transaction-password.md`

  * Added `STEP_UP_ENDPOINT_NOT_ALLOWED` to the initial error contract table.
  * Updated conceptual `step_up_tokens` model with persisted `jti`.

3. `docs/backlogs/api/006d - zta-contracts-and-docs.md`

  * Added `STEP_UP_ENDPOINT_NOT_ALLOWED` to the MVP error contract matrix.

4. `api/docs/05-error_and_response.md`

  * Expanded HTTP client error guidance with 401/403 semantics.
  * Added security and step-up domain error mappings for transactional-password and step-up token flows.

## 2026/05/29 — api/zta-mvp-transactional-password-07

Add step-up authorization flow for transactional password validation.

1. `api/cmd/api/main.go`

   * Wired the step-up token repository and JWT step-up token signer.
   * Registered `AuthorizeStepUpUseCase` with transaction password validation, token persistence, signing, and endpoint policy.
   * Injected the new use case into the security handler.
   * Added `POST /security/step-up/authorize` as an authenticated route.

2. `api/cmd/api/routes_test.go`

   * Added route coverage to ensure the step-up authorization endpoint is protected by the auth middleware.

3. `api/internal/security/application/authorize_step_up.go`

   * Added the step-up authorization use case.
   * Validates authenticated user, endpoint policy, transaction password format, user status, password state, lock state, and password hash.
   * Registers failed attempts and lock behavior.
   * Resets validation state after successful authentication.
   * Persists the step-up token before signing it.
   * Returns a signed step-up token with expiration metadata.

4. `api/internal/security/application/authorize_step_up_test.go`

   * Added coverage for successful authorization.
   * Covered missing user, endpoint not allowed, missing password, invalid password, lock behavior, blocked password, expired lock normalization, and token persistence failure.

5. `api/internal/security/application/create_transaction_password_test.go`

   * Extended test mocks to support validation-state persistence and password comparison assertions.

6. `api/internal/security/application/errors_registry.go`

   * Registered the `ErrStepUpEndpointNotAllowed` domain error.

7. `api/internal/security/application/errors_registry_test.go`

   * Added coverage for the `STEP_UP_ENDPOINT_NOT_ALLOWED` error mapping.

8. `api/internal/security/delivery/handler.go`

   * Added request and response DTOs for step-up authorization.
   * Implemented `AuthorizeStepUp` HTTP handler.
   * Added authenticated user extraction, strict JSON decoding, request trimming, use case delegation, and response mapping.

9. `api/internal/security/delivery/handler_test.go`

   * Added handler tests for success, unauthorized access, invalid payload, mapped domain error, and missing use case configuration.

10. `api/internal/shared/errors/codes.go`

* Added the `STEP_UP_ENDPOINT_NOT_ALLOWED` shared error code.

This commit introduces the first complete step-up authorization endpoint, enabling critical operations to request a short-lived authorization token based on transactional password validation.


## 2026/05/29 — api/zta-mvp-transactional-password-06

This commit adds the infrastructure foundation for step-up tokens used by transactional password flows.

1. `api/internal/security/domain/errors.go`

   * Added `ErrStepUpEndpointNotAllowed` to represent rejected step-up endpoint usage.

2. `api/internal/security/domain/interfaces.go`

   * Added contracts for:

     * `StepUpTokenRepository`
     * `StepUpTokenSigner`
     * `StepUpEndpointPolicy`

3. `api/internal/security/domain/step_up_endpoint_policy.go`

   * Added whitelist-based endpoint validation for step-up operations.
   * Added the default allowed endpoint key for internal transfer creation.
   * Added normalization and rejection of blank or unsupported endpoint keys.

4. `api/internal/security/domain/step_up_endpoint_policy_test.go`

   * Added tests for default endpoint allowance, unknown endpoint rejection, blank input rejection, input trimming, custom whitelist behavior, and nil policy rejection.

5. `api/internal/security/infrastructure/jwt_step_up_token_signer.go`

   * Added JWT signer for persisted step-up tokens.
   * Added minimal step-up claims containing user ID, endpoint key, scope, JTI, issued-at, and expiration.
   * Ensured sensitive data and operation payloads are not embedded in the signed token.

6. `api/internal/security/infrastructure/jwt_step_up_token_signer_test.go`

   * Added tests for JWT claim generation, signature validation, invalid secret rejection, consumed/invalid token rejection, and absence of sensitive payload fields.

7. `api/internal/security/infrastructure/postgres_step_up_token_repository.go`

   * Added PostgreSQL repository for creating, finding, and consuming step-up tokens.
   * Added transaction-aware executor support.
   * Added atomic token consumption with status and expiration checks.
   * Mapped database constraint failures and token state failures to domain errors.

8. `api/internal/security/infrastructure/postgres_step_up_token_repository_test.go`

   * Added integration coverage for step-up token creation, lookup, duplicate JTI handling, successful consumption, repeated consumption, expired token handling, and missing token handling.

9. `api/internal/security/infrastructure/postgres_transaction_password_repository_test.go`

   * Extended the security repository test schema with the `step_up_tokens` table, constraints, indexes, default UUID generation, and cleanup logic.

This establishes the persistence, signing, endpoint policy, and test coverage required for the next step of the transactional password step-up flow.


## 2026/05/29 — api/zta-mvp-transactional-password-05

Introduce the initial domain and database foundation for step-up tokens used by transactional password flows.

1. `Makefile`

   * Added a centralized `DOCKER_COMPOSE` variable with explicit compose file and project name.
   * Replaced direct `docker compose` calls with `$(DOCKER_COMPOSE)`.
   * Updated database reset, readiness check, and schema export commands to use the compose service name.
   * Improved `db-reset` reliability by using `DROP DATABASE ... WITH (FORCE)` and `ON_ERROR_STOP=1`.

2. `api/internal/security/domain/errors.go`

   * Added domain errors for invalid, expired, and already consumed step-up tokens.

3. `api/internal/security/domain/step_up_token.go`

   * Added the `StepUpToken` domain entity.
   * Added active and consumed token statuses.
   * Added default token duration of two minutes.
   * Implemented creation, restoration, validation, expiration checking, and consumption behavior.
   * Enforced UTC normalization and consistency rules for token timestamps and status transitions.

4. `api/internal/security/domain/step_up_token_test.go`

   * Added unit tests covering token creation, validation, restoration, expiration checks, successful consumption, consumed-token rejection, and expired-token rejection.

5. `api/migrations/000011_step_up_tokens.up.sql`

   * Added the `step_up_tokens` table.
   * Added constraints for status validity, non-blank JTI, non-blank endpoint key, expiration consistency, and consumed-at consistency.
   * Added unique and lookup indexes for JTI, user ID, endpoint key, and expiration time.

6. `api/migrations/000011_step_up_tokens.down.sql`

   * Added rollback logic for indexes and the `step_up_tokens` table.

This commit establishes the persistence and domain rules required for short-lived, single-use step-up authorization tokens in the ZTA transactional password MVP.


## 2026/05/29 — api/zta-mvp-transactional-password-04

Advance the Zero Trust Architecture MVP by implementing the transactional password domain and the first operational step-up authorization components. This commit moves the initiative from architectural planning into executable backend functionality, introducing persistence, business rules, API endpoints, and enforcement primitives required for sensitive financial operations.

### Security module implementation

1. Introduced the initial security module structure:

   * `internal/security/domain/*`

   * `internal/security/application/*`

   * `internal/security/infrastructure/*`

   * `internal/security/delivery/*`

   * Established the layered foundation for security-related capabilities.

   * Preserved dependency direction consistent with the existing modular monolith architecture.

   * Isolated transactional authorization concerns from authentication and account modules.

### Transactional password domain

2. Added transactional password domain modeling:

   * Implemented transactional password entities and value objects.
   * Defined business invariants for password lifecycle management.
   * Added validation rules for PIN creation and verification.
   * Introduced failure tracking and lockout-related domain behavior.

3. Added transactional password repository contracts:

   * Defined persistence abstractions for transactional password storage.
   * Decoupled business rules from PostgreSQL implementations.
   * Prepared infrastructure support for future security factors.

### Transactional password persistence

4. Added PostgreSQL persistence support:

   * Implemented repository layer for transactional password management.
   * Added secure password storage using hashed values.
   * Introduced lookup and verification support required by step-up authorization flows.
   * Preserved transactional consistency requirements across security operations.

5. Added database migrations:

   * Created transactional password persistence structures.
   * Added required indexes and constraints.
   * Defined schema required for transactional password lifecycle management.
   * Prepared the database for future security-related extensions.

### Transactional password use cases

6. Implemented transactional password creation flow:

   * Added use case for transactional password registration.
   * Enforced PIN validation rules.
   * Prevented invalid or inconsistent password creation scenarios.
   * Established the first customer-managed authorization factor.

7. Implemented transactional password verification flow:

   * Added password validation use case.
   * Integrated failure counting behavior.
   * Added support for temporary protection against repeated invalid attempts.
   * Returned consistent domain errors for authorization failures.

### HTTP delivery layer

8. Added transactional password endpoints:

   * Implemented request parsing and response mapping.
   * Integrated security use cases with the HTTP layer.
   * Preserved the standard API response envelope.
   * Added authorization requirements aligned with authenticated customer sessions.

9. Added endpoint contracts and DTOs:

   * Defined request and response payloads.
   * Standardized validation behavior.
   * Added error mappings for transactional password operations.
   * Kept API behavior consistent with existing account and authentication endpoints.

### Step-up authorization foundation

10. Introduced step-up authorization services:

    * Added initial step-up token issuance flow.
    * Implemented token validation primitives.
    * Prepared support for endpoint-scoped authorization checks.
    * Established the foundation for sensitive-operation enforcement.

11. Added single-use token behavior:

    * Introduced token consumption semantics.
    * Prevented token reuse after successful authorization.
    * Prepared the infrastructure required for future policy enforcement expansion.

### Tests

12. Added automated test coverage:

    * Created domain tests covering transactional password rules.
    * Added repository tests for persistence behavior.
    * Added use case tests for creation and verification flows.
    * Validated error scenarios, lockout behavior, and authorization failures.

### Documentation updates

13. Updated implementation and backlog documentation:

    * Recorded implemented portions of the ZTA MVP roadmap.
    * Updated progress tracking for transactional password tasks.
    * Refined endpoint contracts and security flow descriptions.
    * Aligned implementation documents with the current codebase state.

This commit delivers the first operational authorization factor of the BankLab Zero Trust Architecture MVP, establishing transactional password management, verification flows, and the foundational building blocks required for step-up authorization of sensitive financial operations.


## 2026/05/28 — api/zta-mvp-transactional-password-03

This commit introduces the first part of the transactional password implementation for the API.

1. `api/cmd/api/main.go`

   * Wired the new security module into the API bootstrap.
   * Added the transaction password repository, bcrypt hasher, use case, handler, and route registration.
   * Registered the protected endpoint:

     * `POST /security/transaction-password`

2. `api/internal/security/domain`

   * Added the transactional password domain model.
   * Defined PIN validation rules for a numeric 6-digit password.
   * Added status handling for active and blocked passwords.
   * Added failure tracking, temporary locking, lock normalization, and success reset behavior.
   * Added domain errors and repository/hasher contracts.

3. `api/internal/security/application`

   * Added `CreateTransactionPasswordUseCase`.
   * Validates authenticated user context, PIN format, confirmation, active user status, and duplicate password creation.
   * Hashes the PIN before persistence.
   * Added security error registration and mappings to API error codes.

4. `api/internal/security/delivery`

   * Added the HTTP handler for creating the initial transaction password.
   * Enforces authenticated access.
   * Rejects unknown JSON fields.
   * Maps domain/application errors to the standard response envelope.

5. `api/internal/security/infrastructure`

   * Added bcrypt-based transaction password hasher.
   * Added PostgreSQL repository for creating, finding, updating validation state, and changing password hashes.
   * Added support for contextual transactions through the existing database transaction context.

6. `api/internal/shared/errors`

   * Added transaction password error codes:

     * `TRANSACTION_PASSWORD_ALREADY_SET`
     * `TRANSACTION_PASSWORD_NOT_SET`
     * `TRANSACTION_PASSWORD_INVALID`
     * `TRANSACTION_PASSWORD_LOCKED`

7. `api/internal/bootstrap/errors.go`

   * Registered security-specific errors during application bootstrap.

8. Tests

   * Added domain tests for PIN validation, password creation, lock behavior, failure tracking, and success reset.
   * Added use case tests covering success, invalid input, inactive user, missing user, duplicate password, hash failure, and repository failure.
   * Added handler tests for success, unauthorized access, invalid payloads, mapped domain errors, nil use case, and unknown errors.
   * Added PostgreSQL integration tests for repository persistence, duplicate handling, validation state updates, and hash updates.
   * Updated route tests to include the new security handler dependency.

9. Documentation

   * Updated REST API documentation with the new transaction password endpoint.
   * Documented request/response payloads, authentication requirements, and related error scenarios.
   * Added transaction password error codes to the documented API error catalog.

This establishes the initial transactional password foundation, including domain rules, persistence, HTTP exposure, error mapping, and test coverage, preparing the API for later transaction-password validation and step-up authorization flows.


## 2026/05/28 — api/zta-mvp-transactional-password-02

This change introduces database-generated UUIDs as the default strategy for persisted entities whenever possible, and starts the transactional password persistence model for the ZTA MVP.

1. Updated UUID generation strategy

   * Moved primary key generation from domain constructors and application code to PostgreSQL defaults using `gen_random_uuid()`.
   * Added `CREATE EXTENSION IF NOT EXISTS pgcrypto` to the baseline schema and test schemas.
   * Updated tables such as `customers`, `accounts`, `users`, `user_sessions`, and `transactions` to use database-generated IDs.
   * Adjusted repositories to omit `id` from inserts and retrieve generated IDs with `RETURNING id`.
   * Updated mocks and tests to populate IDs when simulating successful persistence.

2. Refactored domain constructors

   * Removed direct UUID generation from `Customer`, `CustomerDocument`, `Account`, `User`, `ContactVerification`, and `Transaction` constructors.
   * Simplified `NewUser` by removing the explicit ID argument.
   * Preserved domain invariants while leaving persistence identity assignment to the repository/database layer.

3. Updated registration flow

   * Adjusted user registration so the customer is persisted before creating the CPF document.
   * Normalized and validated CPF before customer persistence.
   * Created the CPF document only after the database has assigned the customer ID.

4. Updated transaction persistence

   * Changed transaction inserts to rely on database-generated IDs.
   * Preserved idempotency behavior with `ON CONFLICT`.
   * Mapped `pgx.ErrNoRows` to duplicate transaction handling when an idempotency conflict occurs.

5. Added transactional password migration

   * Added the `transaction_passwords` table.
   * Added fields for password hash, status, failed attempts, lock control, and timestamps.
   * Added constraints for valid status, non-negative failed attempts, and blocked-state consistency.
   * Added indexes for user uniqueness and lock expiration lookup.

6. Added migration planning documentation

   * Added a backlog document describing the staged migration toward database-generated IDs.
   * Documented the recommended order for moving each entity to the new persistence identity model.

This change intentionally breaks compatibility with the previous development database schema, because the project is still in active development. The practical impact is limited to recreating the current test data, especially the existing test user.


## 2026/05/28 — api/zta-mvp-transactional-password-01

Introduce the first architectural foundation for the Zero Trust Architecture (ZTA) MVP by defining transactional password, step-up authorization flow, enforcement boundaries, and related API contracts. This commit also reorganizes the API security backlog structure and includes supporting mobile/iOS adjustments.

### API ZTA MVP documentation and architecture

1. Added the new implementation document:

   * `api/docs/implementations/03-zta-step-up-transaction-password.md`

     * Defined the first ZTA MVP increment for BankLab.
     * Introduced the transactional password concept as an additional authorization factor.
     * Defined the short-lived step-up token model with single-use semantics.
     * Documented the complete authorization flow using:

       * JWT session authentication
       * transactional password validation
       * endpoint-scoped step-up token enforcement
     * Added Mermaid sequence diagrams and a visual flow image.
     * Defined:

       * endpoint names
       * logical endpoint keys
       * JSON payload fields
       * response envelopes
       * initial error contracts
     * Documented:

       * Policy Enforcement Point (PEP)
       * Policy Engine
       * Transaction Password Factor
       * architectural positioning between delivery and application layers
     * Introduced conceptual data models for:

       * `transaction_passwords`
       * `step_up_tokens`

2. Added ZTA MVP backlog foundation:

   * `docs/backlogs/api/006 - zta-mvp-foundation.md`

     * Established the architectural direction for the new `internal/security` module.
     * Defined module responsibilities and future evolution goals.
     * Documented the initial ZTA enforcement flow and dependency boundaries.
     * Listed MVP scope exclusions and future security evolution points.

3. Added transactional password backlog:

   * `docs/backlogs/api/006a - transaction-password.md`

     * Defined:

       * transactional password creation flow
       * PIN rules
       * temporary blocking behavior
       * failure counting
       * conceptual persistence model
     * Documented all initial business constraints and expected error scenarios.

4. Added transactional password implementation task breakdown:

   * `docs/backlogs/api/006a - transaction-password_tasks.md`

     * Created a complete 8-task implementation roadmap covering:

       * migrations
       * domain modeling
       * repository contracts
       * Postgres persistence
       * use cases
       * HTTP delivery
       * error mapping
       * automated tests
     * Added detailed objectives, scope, acceptance criteria, and dependency mapping for each task.

5. Added step-up token backlog:

   * `docs/backlogs/api/006b - step-up-token.md`

     * Defined the hybrid JWT + persisted `jti` model.
     * Specified:

       * token claims
       * expiration behavior
       * atomic consumption requirements
       * endpoint scoping
       * issuance flow
     * Documented conceptual persistence model and related errors.

6. Added internal transfer enforcement backlog:

   * `docs/backlogs/api/006c - internal-transfer-step-up-enforcement.md`

     * Defined protection rules for:

       * `POST /accounts/internal-transfers`
     * Documented:

       * `X-Step-Up-Token` validation
       * endpoint-key validation
       * single-use enforcement
       * user/token matching
       * atomic token consumption behavior
     * Clarified delivery-to-security interaction boundaries.

7. Added ZTA contracts and documentation backlog:

   * `docs/backlogs/api/006d - zta-contracts-and-docs.md`

     * Consolidated:

       * response envelope standards
       * endpoint naming
       * JSON fields
       * step-up response contract
       * HTTP error mapping
     * Formalized the initial public error-code contract for ZTA operations.

8. Updated backlog index:

   * `docs/backlogs/README.md`

     * Added references to all ZTA MVP backlogs.
     * Reorganized API backlog listing around the new security initiative.

9. Reorganized completed onboarding backlog:

   * Moved:

     * `docs/backlogs/api/001 - onboarding.md`
     * to:
     * `docs/backlogs/api/done/001 - onboarding.md`

### API documentation assets

10. Added new architectural flow image:

    * `api/docs/images/fluxo_senha-trans.png`

      * Added visual representation of the transactional password and step-up authorization flow.

### Mobile Flutter improvements

11. Refactored birthdate selection state handling:

    * `mobile/lib/ui/pages/register/register_birthdate_page.dart`

      * Replaced mutable `DateTime?` state with:

        * `ValueNotifier<DateTime?>`
      * Added `ValueListenableBuilder` for granular UI updates.
      * Simplified reactive rendering for selected date display.
      * Updated initialization and validation flow to use notifier values consistently.
      * Preserved age validation behavior for account creation.

### iOS / Flutter build system updates

12. Updated iOS Flutter integration:

    * `mobile/ios/Runner.xcodeproj/project.pbxproj`
    * `mobile/ios/Runner.xcodeproj/xcshareddata/xcschemes/Runner.xcscheme`
    * `mobile/ios/Podfile.lock`

      * Migrated plugin integration toward Swift Package Manager-based Flutter plugin handling.
      * Added `FlutterGeneratedPluginSwiftPackage`.
      * Added Xcode pre-build prepare script for Flutter framework generation.
      * Removed obsolete CocoaPods embed framework script entries.
      * Updated iOS dependency metadata accordingly.

### Tooling and local environment updates

13. Updated Postman environment:

    * `tools/postman/Environment.postman_environment.json`

      * Adjusted local `base_url` IP address for current development environment.

This commit establishes the first concrete architectural and contractual foundation for ZTA inside BankLab, introducing a layered security model centered on transactional authorization, short-lived step-up tokens, and explicit policy enforcement before sensitive financial operations.



## 2026/05/22 - milestone/basic-banking-core

This commit establishes an important architectural milestone for the BankLab project.

At this stage, the project consolidates a functional and coherent foundation for a basic banking platform, including transactional consistency, authentication flows, onboarding evolution, ledger operations, mobile integration, and progressive registration flows. Although the next natural step would be the implementation of KYC and stronger contextual verification mechanisms, the current state already fulfills the original objective of building a robust experimentation base for Zero Trust Architecture studies in mobile and backend environments.

This milestone also becomes strategically relevant because it creates a reusable foundation for future parallel projects that may evolve independently from the onboarding and regulatory layers.

### Documentation and Historical Consolidation

* Added a comprehensive technical evolution report:

  * `docs/relatorios/2026-05-18 - 05-22.md`

    * Consolidates the implementation history between:

      * `2026/05/18 - api/pre-onboarding-01`
      * `2026/05/22 - api/pg_cron-01`
    * Documents the transition from a simple CPF-centric registration model into a progressive onboarding architecture.
    * Summarizes:

      * identity normalization;
      * contact verification flows;
      * onboarding hardening;
      * mobile registration restructuring;
      * draft persistence architecture;
      * pg_cron lifecycle maintenance;
      * modularization efforts;
      * testing evolution;
      * architectural preparation for future ZTA/KYC flows.
    * Establishes a historical and architectural reference point for the current maturity of the project.

### Documentation Structure Reorganization

* Reorganized technical reports into a dedicated reports directory:

  * moved:

    * `docs/relatorio-api-implementada-2026-05-12.md`
  * to:

    * `docs/relatorios/relatorio-api-implementada-2026-05-12.md`

* Reorganized mobile implementation reports:

  * moved:

    * `docs/relatorio-mobile-implementado-2026-05-12.md`
  * to:

    * `docs/relatorios/relatorio-mobile-implementado-2026-05-12.md`

### Architectural Significance

This commit represents the closing of the project's first major operational cycle.

The system now contains:

* transactional banking core;
* account and ledger operations;
* JWT authentication model;
* onboarding orchestration;
* progressive mobile registration;
* contact verification lifecycle;
* secure onboarding persistence;
* modular layered architecture;
* PostgreSQL transactional consistency strategies;
* documentation sufficiently mature to support future contributors and derived projects.

From this point forward, the project can evolve in multiple directions independently:

* KYC and regulatory onboarding;
* contextual Zero Trust evaluation;
* MFA and device binding;
* liveness verification;
* web clients;
* production-grade operational hardening;
* derived banking and fintech experiments.

This branch is intended to preserve the current state as a stable architectural milestone before the project transitions into deeper security and contextual trust experimentation.


## 2026/05/22 — api/pg_cron-01

Introduce PostgreSQL 17 with `pg_cron` support and migrate contact verification cleanup from trigger-based execution to scheduled database jobs.

### Infrastructure and PostgreSQL runtime

* Replaced the default PostgreSQL 16 container with a custom PostgreSQL 17 image
* Added `infra/docker/postgres/Dockerfile`

  * installs `postgresql-17-cron`
  * keeps the image minimal by removing cached package metadata
* Updated `docker-compose.yml`

  * switched from direct `postgres:16` image usage to local image build
  * enabled `shared_preload_libraries=pg_cron`
  * configured:

    * `cron.database_name=bank`
    * `cron.timezone=America/Sao_Paulo`
* Standardized local infrastructure documentation around the new PostgreSQL runtime

### Contact verification cleanup redesign

* Reworked the contact verification cleanup strategy
* Removed the previous trigger-driven cleanup execution model

  * deleted:

    * `cleanup_contact_verifications_if_due()`
    * trigger-based cleanup scheduling
    * `contact_verification_cleanup_runs`
* Added `pg_cron`-based scheduled cleanup execution

  * creates the extension automatically with:

    * `CREATE EXTENSION IF NOT EXISTS pg_cron`
  * schedules:

    * `cleanup-contact-verifications`
    * daily execution at `03:00`
* Added migration safety logic

  * unschedules existing jobs before recreating them
  * safely removes jobs during rollback migrations
* Preserved immediate cleanup execution after migration through:

  * `SELECT cleanup_contact_verifications();`

### Database consistency and verification lifecycle

* Enforced a stronger verification lifecycle model
* Added uniqueness guarantees for `(target, channel)`

  * replaces previous pending verification attempts automatically
* Removed the obsolete `(target, channel)` non-unique index
* Added cleanup retention rules:

  * expired unverified records removed after 24 hours
  * verified records retained for 7 days
* Updated cleanup-related indexes and migration behavior accordingly

### Documentation updates

Updated multiple project documents to reflect the new PostgreSQL and scheduling architecture:

* `README.md`
* `README_en.md`
* `api/README.md`
* `api/docs/00-getting_started.md`
* `api/docs/06-implementation.md`
* `api/docs/09-database.md`
* `api/docs/infra.md`

Main documentation additions include:

* PostgreSQL 17 adoption
* `pg_cron` usage rationale
* custom Docker image explanation
* scheduled maintenance behavior
* onboarding verification lifecycle updates
* cleanup policy documentation

### Mobile documentation adjustments

Refined mobile onboarding documentation to reflect the current multi-page registration flow:

* updated onboarding references from the old single `RegisterPage`
* documented the new step-by-step registration structure
* updated route and architecture references
* clarified onboarding persistence and verification flow behavior

### Visual assets

* Updated database diagram image:

  * `api/docs/images/database.png`

This commit consolidates scheduled maintenance responsibilities inside PostgreSQL itself, simplifies application-side cleanup orchestration, and establishes a cleaner operational foundation for onboarding verification lifecycle management.


## 2026/05/22 — mobile/pre-onboarding-16

Refined the pre-onboarding registration flow by strengthening contact verification consistency on the API side and improving password validation and UX behavior on the Flutter client.

### API

1. Improved contact verification uniqueness handling

   * Replaced the non-unique `idx_contact_verifications_target_channel` index with the unique index `contact_verifications_unique_target_channel`.
   * Enforced a single active verification attempt per `(target, channel)` pair.
   * Updated integration test schema setup to reflect the new database constraint.

2. Added contact verification replacement semantics

   * Updated `PostgresContactVerificationRepository.CreateContactVerification` to support `ON CONFLICT (target, channel)`.
   * Existing verification attempts are now replaced atomically when a new request is created for the same target/channel.
   * Reset verification state on replacement:

     * `verification_token`
     * `verified_at`
     * expiration metadata
     * creation timestamp
   * Preserved deterministic behavior during repeated onboarding attempts.

3. Added automated cleanup migration for contact verifications

   * Created migration `000009_contact_verifications_cleanup.up.sql`.
   * Added deduplication cleanup before applying the unique index.
   * Added partial indexes for:

     * unverified expired records
     * verified records
   * Introduced `contact_verification_cleanup_runs` control table.
   * Added PostgreSQL cleanup functions:

     * `cleanup_contact_verifications`
     * `cleanup_contact_verifications_if_due`
   * Added trigger-based periodic cleanup execution.
   * Used advisory locks to prevent concurrent cleanup execution.
   * Established retention policies:

     * expired unverified records older than 24h
     * verified records older than 7 days
   * Added matching rollback migration.

### Mobile

1. Refactored password domain model

   * Converted `PasswordModel` constructor to named parameters.
   * Added centralized password validation rules:

     * minimum length
     * uppercase requirement
     * lowercase requirement
     * numeric requirement
   * Added semantic getters:

     * `hasNumber`
     * `hasUppercase`
     * `hasLowercase`
     * `hasMinLength`
     * `isValidPassword`
     * `hasEquals`
     * `isValid`
   * Centralized minimum length through `PasswordModel.minLength`.

2. Introduced `PasswordDraft`

   * Added mutable UI-oriented password draft model.
   * Encapsulated UI validation state around the domain `PasswordModel`.
   * Reduced duplicated password validation logic inside the page layer.
   * Improved separation between UI state and domain validation rules.

3. Simplified password page validation flow

   * Removed duplicated regex and validation logic from `RegisterPasswordPage`.
   * Replaced manual validation with `PasswordDraft` state accessors.
   * Centralized enable/disable button logic using `PasswordDraft.isValid`.
   * Updated password criteria label to use the centralized minimum length constant.
   * Removed direct submit on keyboard action to avoid premature submissions.

4. Improved registration diagnostics

   * Added structured registration request logging inside `RegistrationApi`.
   * Improved onboarding request traceability during integration/debug sessions.

5. Updated registration tests

   * Adjusted registration use case tests to reflect the new password policy.
   * Replaced weak password fixtures with compliant values (`Secret123`).

### Dependencies and tooling

1. Updated Flutter dependencies

   * Refreshed multiple package versions including:

     * `flutter_secure_storage`
     * `go_router`
     * `google_fonts`
     * `objective_c`
     * `meta`
     * `vm_service`
     * others from transitive dependency resolution.

2. Updated Postman environment

   * Adjusted local API base URL for the current development environment.

This commit consolidates an important pre-onboarding foundation by improving deterministic verification handling on the backend while moving password validation rules into a cleaner and more reusable domain-oriented structure on the Flutter client.


## 2026/05/21 — mobile/pre-onboarding-14

This commit advances the mobile onboarding and authentication experience by introducing a complete password registration flow, registration status feedback pages, and a broader UI component reorganization focused on reusable input abstractions.

### Registration Flow Enhancements

1. Added dedicated registration status routes and screens

   * Introduced success and failure routes in `register_routes.dart`
   * Added `RegisterStatusPage` to present onboarding completion feedback
   * Added navigation paths from the password submission flow to success/failure states
   * Included contextual actions:

     * success → navigate to login
     * failure → retry registration flow

2. Implemented the complete password registration page

   * Replaced the placeholder implementation in `register_password_page.dart`
   * Added:

     * password field
     * confirmation field
     * visibility toggle controls
     * focus management
     * validation feedback
     * reactive enable/disable logic
   * Added live password criteria validation:

     * minimum length
     * uppercase characters
     * lowercase characters
     * numeric characters
     * password confirmation matching
   * Integrated submit + register execution flow with proper error handling and navigation transitions

3. Added reusable password model abstraction

   * Introduced `PasswordModel`
   * Centralized password normalization and validation
   * Replaced tuple-based `(String, String)` transport with a semantic domain model
   * Improved readability and future extensibility of the onboarding flow

4. Updated register use case and viewmodel contracts

   * `RegisterUsecase.submitPassword` now receives `PasswordModel`
   * `RegisterViewmodel.submitPassword` updated accordingly
   * Added re-export support for password-related models directly from the usecase module
   * Updated tests to reflect the new API contract

### UI Component Refactor and Reorganization

1. Reorganized input components into dedicated namespaces

   * Moved:

     * `basic_text_form_field.dart` → `input_text/basic_input_text.dart`
     * `verification_code_field.dart` → `input_text/otp_input.dart`
     * `money_input_formatter.dart` → `input_formatters/money_input_formatter.dart`
   * Improved naming consistency and architectural clarity for reusable UI primitives

2. Renamed and generalized reusable input widgets

   * `BasicTextFormField` → `BasicInputText`
   * `VerificationCodeField` → `OtpInput`
   * Added `focusNode` support to the base text component

3. Migrated all onboarding/authentication pages to the new input components

   * Login page
   * Short login page
   * Register CPF page
   * Register email page
   * Register name page
   * Register phone page
   * Register token page
   * Transfer recipient page
   * Transfer payment page

### UX and Visual Improvements

1. Added password criteria helper widget

   * Introduced `CriterialItemRow`
   * Displays dynamic validation state using visual indicators

2. Improved transfer status screen readability

   * Cached `colorScheme` and `textTheme`
   * Reduced repeated theme access
   * Improved readability and maintainability of the widget tree

3. Minor route formatting cleanup

   * Reformatted transfer failure route builder for consistency

### Architectural Direction

This commit reinforces an important architectural transition in the mobile onboarding flow:

* replacing primitive transport structures with semantic models
* consolidating reusable input infrastructure
* isolating onboarding states into explicit routes
* improving UI composition around focused reusable widgets

The resulting structure makes the registration pipeline more maintainable, more expressive, and significantly better prepared for future onboarding extensions such as:

* transactional passwords
* device registration
* MFA/TOTP
* Zero Trust contextual validation
* biometric/liveness verification flows


## 2026/05/21 — mobile/pre-onboarding-14

This commit advances the mobile onboarding flow by introducing phone registration and token confirmation screens, while significantly improving validation, formatting, and verification code input behavior across the registration experience.

### Main Improvements

1. Added Brazilian phone validation support

   * Introduced `String.isValidPhone` extension for validating Brazilian landline and mobile numbers.
   * Added DDD validation rules.
   * Differentiated validation logic between mobile (`9XXXXXXXX`) and landline (`2-5XXXXXXX`) formats.
   * Centralized phone validation behavior inside the core string extensions layer.

2. Implemented `PhoneInputFormatter`

   * Added automatic phone formatting while typing.
   * Supports:

     * `(XX) XXXXX-XXXX`
     * `(XX) XXXX-XXXX`
   * Automatically strips non-numeric characters before formatting.
   * Keeps formatting logic isolated in reusable input formatter infrastructure.

3. Refactored `VerificationCodeField`

   * Improved keyboard navigation and editing experience.
   * Added focus-aware selection behavior.
   * Implemented proper backspace navigation between cells.
   * Fixed incorrect first-empty-cell focus condition.
   * Added clipboard paste support centralized through `_onTap`.
   * Simplified layout structure and removed unnecessary `Padding` wrappers.
   * Improved visual appearance:

     * outlined borders
     * rounded corners
     * bolder typography
     * space distribution between fields
   * Added internal helpers:

     * `_focusAndSelectCellContent`
     * `_selectCellContent`

### Registration Flow Enhancements

4. Implemented `RegisterPhonePage`

   * Added complete phone onboarding screen.
   * Integrated:

     * `BasicTextFormField`
     * `PhoneInputFormatter`
     * validation state management
     * async submission flow
     * navigation to phone token confirmation
   * Added bottom action bar with:

     * back action
     * continue action
     * loading/disabled state integration
   * Added snackbar error feedback support.
   * Added initialization logic for persisted onboarding state.

5. Implemented token confirmation UX in `RegisterTokenPage`

   * Added verification code UI using `VerificationCodeField`.
   * Added support for:

     * email token confirmation
     * phone token confirmation
   * Added dynamic header messages per verification channel.
   * Added async command integration for token confirmation.
   * Added validation and loading state management.
   * Added snackbar-based error handling.
   * Added continue/back navigation actions.

### UI Consistency Improvements

6. Standardized onboarding page spacing

   * Added `spacing: 12` to onboarding forms:

     * `register_birthdate_page`
     * `register_cpf_page`
     * `register_email_page`
     * `register_name_page`
   * Removed redundant `SizedBox` spacing where applicable.
   * Improved visual consistency across onboarding steps.

### Architectural Notes

This commit also reinforces some important architectural directions already present in the project:

* reusable input formatting components
* validation encapsulated in extensions
* UI state isolated with `ValueNotifier`
* async command orchestration through `Command`
* reusable onboarding navigation patterns
* clearer separation between formatting, validation, and UI concerns

The onboarding flow now has a substantially more complete pre-verification experience, especially around phone registration and token confirmation, while also improving the overall interaction quality of the verification code input system.


## 2026/05/21 — mobile/pre-onboarding-13

Refactor the contact verification flow to use strongly typed verification channels and introduce a reusable verification code input component for onboarding flows.

### Main Improvements

* Replaced raw string-based contact verification channels with the new `ContactVerificationChannel` enum across routing, DTOs, use cases, pages, repositories, APIs, and tests.
* Added a reusable `VerificationCodeField` widget focused on OTP/token entry UX improvements.
* Consolidated onboarding token handling around a single domain representation for verification channels.
* Improved typing consistency and reduced risks related to invalid string values during onboarding flows.

### Contact Verification Refactor

1. Added `ContactVerificationChannel` enum

   * Introduced:

     * `email`
     * `phone`
   * Added `fromString()` parser with validation and explicit error handling.
   * Centralized channel conversion logic.

2. Updated DTOs to use typed channels

   * Refactored:

     * `ContactVerificationRequestDto`
     * `ContactVerificationRequestResponseDto`
     * `ContactVerificationConfirmResponseDto`
   * Replaced `String channel` with `ContactVerificationChannel channel`.
   * Updated serialization/deserialization logic:

     * `.name` for outbound payloads
     * `fromString()` for inbound payloads

3. Updated onboarding use cases

   * Refactored `RegisterUsecase` to remove hardcoded channel strings.
   * Email and phone verification requests now use enum values directly.

4. Updated route and UI integration

   * Refactored `register_routes.dart`
   * Replaced legacy `TokenType` usage with `ContactVerificationChannel`.
   * Simplified `RegisterTokenPage` token flow handling.

5. Removed duplicated token channel abstraction

   * Eliminated the local `TokenType` enum from `RegisterTokenPage`.
   * Unified the onboarding verification flow around the API/domain enum.

### Verification Code Field Component

1. Added `VerificationCodeField`

   * Created reusable OTP/token input widget with:

     * one digit per field
     * automatic focus progression
     * backspace navigation behavior
     * clipboard paste support
     * autofill integration (`oneTimeCode`)
     * configurable dimensions and spacing
     * initial value support
     * completion callback support

2. Implemented improved keyboard handling

   * Backspace clears current field.
   * When empty, backspace navigates to the previous field.
   * Automatically advances focus after digit insertion.

3. Added paste handling

   * Detects clipboard numeric content.
   * Automatically distributes pasted digits across fields.
   * Correctly updates focus and completion state.

### Test Updates

1. Refactored repository tests

   * Updated `contact_verification_repository_impl_test.dart`
   * Replaced string assertions with enum assertions.

2. Refactored API tests

   * Updated `contact_verification_api_test.dart`
   * Adjusted request/response expectations for typed channels.

3. Refactored DTO tests

   * Updated serialization/deserialization validations.
   * Ensured enum parsing consistency.

4. Refactored use case tests

   * Updated fake repositories and verification flow assertions.
   * Removed string comparisons in verification channel checks.

### Architectural Impact

This refactor improves onboarding consistency by removing stringly-typed channel handling from the registration flow. The verification system now has stronger compile-time guarantees, clearer intent across layers, and reduced risk of invalid state propagation during email and phone verification operations.

The new verification input component also establishes a reusable foundation for future onboarding and security flows, including:

* transactional password confirmation
* MFA/TOTP verification
* device registration confirmation
* password recovery tokens
* Zero Trust challenge flows


## 2026/05/21 - mobile/pre-onboarding-12

Refined the pre-onboarding contact verification flow across API, mobile, tests, and documentation, with emphasis on separating temporary development/debug behavior from the stable verification contract.

### API and verification contract adjustments

1. Updated contact verification response semantics

   * Replaced `token` with `debug_token` in the contact verification request response.
   * Changed the API contract to explicitly treat the verification token as a temporary debug-only artifact while no e-mail/SMS provider exists.
   * Kept the stable contract centered on:

     * `verification_id`
     * `channel`
     * `target`
     * `expires_at`

2. Improved backend response modeling

   * Replaced:

     * `Token string`
   * With:

     * `DebugToken *string`
   * Added `omitempty` semantics to avoid coupling clients to temporary debug data.

3. Updated integration and handler tests

   * Adjusted request/confirmation integration flow to consume `debug_token`.
   * Added explicit assertions validating the presence of `debug_token` during local/dev execution.
   * Updated API handler tests and use case tests to reflect the new contract semantics.

4. Clarified operational documentation

   * Updated REST API documentation and onboarding backlog documents to explain:

     * why `debug_token` currently exists
     * why clients must not depend on it
     * future removal expectations once notification providers are integrated

### Mobile onboarding and DTO refactoring

1. Removed verification token from the main DTO model

   * Simplified `ContactVerificationRequestResponseDto`.
   * Removed `token` from the registration flow DTO structure.
   * Preserved the verification flow contract independently from temporary debug data.

2. Moved debug token handling to the API boundary

   * `ContactVerificationApi` now reads `debug_token` directly from the raw response envelope only in development mode.
   * Prevented propagation of debug-only fields into domain/application layers.
   * Kept debug logging available for local testing while preserving clean architecture boundaries.

3. Refactored `RegisterDraftState`

   * Replaced imperative `updateX()` methods with property setters:

     * `cpf =`
     * `name =`
     * `birthDate =`
     * `email =`
     * `phone =`
     * verification-related setters
   * Improved readability and reduced noise in `RegisterUsecase`.

4. Expanded model documentation

   * Added inline documentation to:

     * dirty tracking behavior
     * persistence lifecycle
     * snapshot hydration
     * timestamp semantics
     * mutation internals

5. Updated register use case orchestration

   * Migrated all draft mutations to setter-based syntax.
   * Preserved dirty-state persistence behavior and verification orchestration flow.

### Mobile tests and cleanup

1. Updated repository and DTO tests

   * Removed token expectations from DTO parsing tests.
   * Adjusted repository/use case fixtures to align with the new contract.

2. Updated draft state tests

   * Migrated all state mutation tests to the setter-based API.
   * Preserved validation for:

     * normalization
     * dirty tracking
     * persistence timestamps
     * snapshot hydration

3. Import normalization and cleanup

   * Standardized extension imports in register pages.
   * Removed outdated relative imports.

### Infrastructure and development workflow

1. Improved Docker/Colima startup behavior

   * Removed implicit `colima start` execution from `docker-up`.
   * Moved Colima initialization to `docker-check`.
   * Reduced unnecessary VM startup attempts during regular Docker commands.
   * Improved separation between:

     * environment validation
     * container lifecycle operations

2. Updated local Postman environment

   * Adjusted local `base_url` IP for current development environment.

This commit reinforces an important architectural distinction in the onboarding flow: temporary operational/debug behavior must remain isolated from stable application contracts and domain models.


## 2026/05/20 - mobile/pre-onboarding-11

Refactor and expand the mobile onboarding flow with persistent draft recovery, new registration screens, improved validation helpers, and a more complete pre-onboarding user experience.

### Core validation and utility improvements

* Added `DateTime.age` extension to centralize age calculation logic from birth dates.
* Added `String.isValidEmail` extension for reusable email validation across the onboarding flow.
* Updated `BasicTextFormField` to expose `textCapitalization` directly through the component API.
* Standardized several imports to use absolute project paths for consistency and maintainability.

### Register draft persistence redesign

Refactored the onboarding draft persistence layer to simplify the recovery model and make draft handling deterministic.

#### Repository layer

* Simplified `RegisterDraftRepository.getByCPF()` to return `RegisterDraftSnapshot` directly instead of `RegisterDraftLoadResult`.
* Removed the `RegisterDraftLoadResult` sealed hierarchy entirely.
* Added automatic draft recreation behavior when:

  * the draft does not exist
  * the draft has expired
* Added `_createASnapshotForCPF()` helper to centralize empty snapshot creation.
* Added `_isOld()` helper to isolate TTL expiration logic.
* Adjusted repository behavior to always return a valid snapshot when possible instead of exposing "not found" states to upper layers.

#### Store layer

* Simplified `RegisterDraftStore.getByCPF()` return type.
* Changed missing or invalid cache entries to return:

  * `Failure(AppErrorCode.storageNotFound)`
    instead of synthetic success states.
* Added cleanup of corrupted JSON cache entries before returning failures.
* Removed obsolete exports related to `RegisterDraftLoadResult`.

#### Domain model

* Added `RegisterDraftSnapshot.empty(String cpf)` factory constructor.
* Centralized initialization of empty onboarding state snapshots.

### Register use case refactor

Refactored onboarding state initialization and CPF recovery behavior.

* Removed the old `initialize()` flow.
* Introduced `startEmptyRegisterState()` as the explicit onboarding bootstrap entrypoint.
* Changed CPF submission flow to:

  * recover cached onboarding state automatically
  * recreate snapshots transparently when necessary
  * rebuild `RegisterDraftState` directly from persisted snapshots
* Removed legacy commented initialization logic.

### Registration UI expansion

Implemented and refined major portions of the onboarding user interface.

#### Register CPF page

* Updated onboarding copy:

  * `Qual o seu CPF?` → `Informe o CPF`
  * improved CPF hint message
* Improved user guidance for numeric CPF input.

#### Register Name page

Implemented the complete name step UI.

* Added:

  * `TextHeader`
  * `BasicTextFormField`
  * snackbar error handling
  * bottom navigation controls
* Added full name validation requiring at least two words.
* Added onboarding state restoration from cached draft.
* Added forward/back navigation handling.
* Added loading state integration with async commands.

#### Register Birthdate page

Implemented the complete birth date onboarding step.

* Added date picker interaction.
* Added minimum age validation (18+).
* Added onboarding state restoration from cached draft.
* Added disabled state handling for invalid dates.
* Added snackbar-based error feedback.
* Added navigation controls and async submit handling.
* Added formatted date presentation using `intl`.

One important focus of this commit was centralizing and formalizing birth date validation through the new `DateTime.age` extension, allowing the onboarding flow to consistently enforce minimum age requirements.

#### Register Email page

Implemented the complete email onboarding step.

* Added:

  * email input field
  * email validation flow
  * async token request handling
  * snackbar-based error feedback
  * navigation controls
* Integrated reusable `String.isValidEmail` validation helper.
* Added onboarding state restoration from persisted draft state.
* Added autofill support for email input.

#### Other onboarding pages

Standardized onboarding titles across:

* password page
* phone page
* token page
* birthdate page
* email page
* name page

Updated:

* `Criar conta`
  to:
* `Registro de Conta`

for a more consistent onboarding identity.

### Test suite updates

Updated repository, store, and use case tests to reflect the new onboarding persistence model.

#### Repository tests

* Reworked tests to validate:

  * automatic draft recreation
  * expired draft replacement
  * empty snapshot generation
  * storage failure propagation
* Added tracking of persisted snapshots in fake store implementations.

#### Store tests

* Updated expectations to use failure-based not-found semantics.
* Added validation for corrupted cache cleanup behavior.

#### Use case tests

* Replaced old initialization tests with:

  * `startEmptyRegisterState()` coverage
* Removed obsolete reset/reinitialize scenarios.
* Simplified fake repository behavior around snapshot recovery.
* Updated all onboarding setup helpers to use the new onboarding bootstrap flow.

### Tooling

* Updated Postman environment local API IP:

  * `192.168.0.16`
  * → `192.168.0.14`


## 2026/05/20 — mobile/pre-onboarding-10

This commit advances the mobile onboarding foundation by introducing the first functional CPF registration flow, refining registration state initialization, improving domain-specific error handling, and consolidating reusable bottom action components across the UI.

### Registration and CPF onboarding flow

* Implemented the first interactive CPF registration screen flow.
* Added CPF validation and formatting support directly in the UI layer.
* Added CPF availability verification before continuing registration.
* Added navigation integration between registration steps using GoRouter.
* Introduced snackbar-based feedback for onboarding validation failures.
* Added automatic focus dismissal after valid CPF input.
* Added login shortcut from the onboarding screen.

#### `mobile/lib/ui/pages/register/register_cpf_page.dart`

* Reworked the page from a placeholder screen into a functional onboarding step.
* Added:

  * `TextEditingController`
  * CPF validation state
  * formatted CPF input
  * navigation handlers
  * asynchronous CPF availability verification
  * onboarding action buttons
  * snackbar-based error presentation
* Added integration with:

  * `CpfInputFormatter`
  * `String.onlyNumbers`
  * `String.isValidCpf`
  * `DoubleBottomButton`
  * `AppSnackbar`

### Registration state lifecycle refactor

A major focus of this commit was simplifying the onboarding initialization flow by separating empty registration startup from future draft restoration logic.

#### `mobile/lib/domain/usecases/register/register_usecase.dart`

* Introduced `startEmptyRegisterState()`.
* Removed the dependency on asynchronous initialization for empty registration flows.
* Temporarily isolated the draft restoration logic through commented preservation for future recovery work.
* Improved CPF validation behavior:

  * replaced generic `invalidData`
  * introduced specific `cpfAlreadyRegistered` handling

#### `mobile/lib/ui/pages/register/viewmodel/register_viewmodel.dart`

* Removed `initialize` command abstraction.
* Added direct `startEmptyRegisterState()` forwarding to the use case.

### Domain-specific onboarding errors

#### `mobile/lib/core/result/errors/app_error_code.dart`

* Added:

  * `cpfAlreadyRegistered`

This improves semantic error handling and prepares the onboarding pipeline for more granular validation flows.

### DTO consistency and API alignment

#### `mobile/lib/data/services/apis/registration/dtos/cpf_check_response_dto.dart`

* Fixed typo:

  * `avaliable` → `available`
* Updated JSON parsing accordingly.
* Improved naming consistency between backend payloads and frontend DTOs.

### Registration navigation groundwork

Several onboarding pages received route navigation preparation methods for the future multi-step registration flow.

#### Updated pages

* `register_name_page.dart`
* `register_birthdate_page.dart`
* `register_email_page.dart`
* `register_token_page.dart`

Changes include:

* GoRouter integration
* next-step navigation methods
* cleanup of duplicated `super.initState()` calls

### Reusable bottom action component

#### `mobile/lib/ui/components/buttons/double_bottom_buttons.dart`

* Added reusable dual-action bottom button component.
* Encapsulates:

  * primary action button
  * secondary text button
  * enable/disable logic
  * optional icon support

### Transfer status page cleanup

#### `mobile/lib/ui/pages/home/transfer/transfer_status_page.dart`

* Replaced duplicated bottom action layout with `DoubleBottomButton`.
* Reduced UI duplication and improved consistency with onboarding screens.

### General cleanup

#### Registration pages

* Removed duplicated `super.initState()` calls from:

  * `register_birthdate_page.dart`
  * `register_email_page.dart`
  * `register_name_page.dart`
  * `register_password_page.dart`
  * `register_phone_page.dart`
  * `register_token_page.dart`

### Test coverage

#### `mobile/test/data/services/apis/registration/dtos/cpf_check_response_dto_test.dart`

* Added DTO parsing coverage for:

  * `available`

#### `mobile/test/domain/usecases/register/register_usecase_test.dart`

* Updated tests to reflect DTO field rename:

  * `avaliable` → `available`

This commit establishes the first concrete onboarding interaction flow in the mobile application while also improving domain error semantics, reducing UI duplication, and preparing the registration pipeline for future persistence and recovery capabilities.


## 2026/05/20 — mobile/pre-onboarding-09

Refactor and expand the mobile onboarding flow structure with dedicated register routes, draft persistence improvements, and the first implementation layer for the multi-step registration experience.

### Routing and onboarding flow restructuring

Refactored the registration route organization to support a fully segmented onboarding flow with explicit route paths.

Files updated:

* `mobile/lib/core/routing/routes.dart`
* `mobile/lib/core/routing/routes/register_routes.dart`

Changes:

* Reworked `RegisterRoutes` to use explicit `/register/*` paths.
* Renamed `fullName` route to `name` for consistency.
* Removed `passwordConfirmation` route from the enum.
* Added new `failure` route placeholder.
* Added route definitions for:

  * name
  * birth date
  * email
  * email token
  * phone
  * phone token
  * password
* Centralized all onboarding pages under `ui/pages/register`.
* Added token-type-aware route handling using `TokenType.email` and `TokenType.phone`.

This restructuring prepares the onboarding module for:

* deep linking
* route isolation
* independent validation stages
* resumable onboarding sessions
* future risk-analysis checkpoints during onboarding

### Registration UI module extraction

Moved the onboarding UI structure out of the auth module into its own dedicated registration module.

Files updated:

* `mobile/lib/ui/pages/register/register_cpf_page.dart`
* `mobile/lib/ui/pages/register/viewmodel/register_viewmodel.dart`
* `mobile/lib/ui/viewmodels.dart`

Changes:

* Migrated registration pages from:

  * `ui/pages/auth/register/*`
    to:
  * `ui/pages/register/*`
* Updated dependency imports accordingly.
* Adjusted page padding from `24` to `12` for a denser onboarding layout baseline.

This separation improves architectural clarity between:

* authentication/session flows
* onboarding/account creation flows

### New onboarding pages scaffold

Added the initial structure for the remaining onboarding screens.

New files:

* `mobile/lib/ui/pages/register/register_birthdate_page.dart`
* `mobile/lib/ui/pages/register/register_email_page.dart`
* `mobile/lib/ui/pages/register/register_name_page.dart`
* `mobile/lib/ui/pages/register/register_password_page.dart`
* `mobile/lib/ui/pages/register/register_phone_page.dart`
* `mobile/lib/ui/pages/register/register_token_page.dart`

Changes:

* Added `SafeScaffold`-based page templates.
* Added keyboard dismissal handling with `GestureDetector`.
* Added dedicated `TokenType` enum for verification flows.
* Established isolated stateful page structure for future form logic integration.

These pages currently serve as structural placeholders for:

* form implementation
* validation
* verification UX
* onboarding orchestration

### Register draft persistence refactor

Reorganized the register draft cache layer and improved dependency injection setup.

Files updated:

* `mobile/lib/data/repositories/register_draft/register_draft_repository.dart`
* `mobile/lib/data/repositories/register_draft/register_draft_repository_impl.dart`
* `mobile/lib/data/services/services.dart`

Files moved:

* `mobile/lib/data/services/cache/register_draft/register_draft_load_result.dart`
* `mobile/lib/data/services/cache/register_draft/register_draft_store.dart`

Changes:

* Removed `register_draft` from the `last_login` namespace.
* Promoted register draft persistence into its own cache domain.
* Added `RegisterDraftStore` as a lazy singleton.
* Updated all imports across repositories, services, use cases, and tests.

This better reflects the actual responsibility of the draft cache:

* onboarding state persistence
* onboarding recovery
* future interrupted-session restoration

instead of coupling it to login behavior.

### Register use case persistence optimization

Simplified and centralized draft persistence behavior inside `RegisterUsecase`.

File updated:

* `mobile/lib/domain/usecases/register/register_usecase.dart`

Changes:

* Replaced duplicated draft-save blocks with `_saveDirty()`.
* Added centralized dirty-state persistence logic.
* Added:

  * state existence validation
  * dirty-check optimization
  * clean-state marking after successful persistence
* Updated all onboarding submit/confirm methods to use the centralized save flow.

Benefits:

* reduced duplicated logic
* safer persistence lifecycle
* cleaner onboarding orchestration
* improved maintainability for future onboarding stages

### Tests and maintenance updates

Updated all affected test imports after cache-layer refactor.

Files updated:

* `mobile/test/data/repositories/register_draft/register_draft_repository_impl_test.dart`
* `mobile/test/data/services/cache/register_draft/register_draft_store_test.dart`
* `mobile/test/domain/usecases/register/register_usecase_test.dart`

### Environment update

File updated:

* `tools/postman/Environment.postman_environment.json`

Changes:

* Updated local API base URL from:

  * `192.168.0.14`
    to:
  * `192.168.0.16`

This commit establishes the structural foundation for the new onboarding architecture, separating onboarding concerns from authentication flows while preparing the application for resumable multi-step registration, verification stages, and future Zero Trust onboarding evolutions.


## 2026/05/20 — mobile/pre-onboarding-08

Refactor pre-onboarding architecture by splitting authentication responsibilities into dedicated APIs and repositories, reorganizing cache modules, and simplifying dependency injection wiring.

This commit restructures the mobile authentication and onboarding flow into clearer bounded contexts, separating login/session responsibilities from registration and contact verification concerns. The result is a cleaner dependency graph, improved modularity, and a more maintainable onboarding foundation for future Zero Trust and identity validation flows.

### Main architectural changes

1. Split `AuthApi` responsibilities into dedicated modules:

   * `AuthApi`
   * `RegistrationApi`
   * `ContactVerificationApi`

2. Split repository responsibilities previously concentrated in `AuthRepository`:

   * `RegistrationRepository`
   * `ContactVerificationRepository`

3. Reorganized cache structure:

   * moved `auth/cache/*` to `cache/last_login/*`
   * decoupled cache concerns from authentication domain semantics

4. Simplified dependency injection:

   * replaced verbose factory closures with constructor tear-offs
   * adopted cleaner `AutoInjector` registrations
   * introduced lazy singletons where appropriate

5. Refactored onboarding use case dependencies:

   * `RegisterUsecase` now depends explicitly on:

     * `RegistrationRepository`
     * `ContactVerificationRepository`
     * `RegisterDraftRepository`

### Repository layer refactor

#### `mobile/lib/data/repositories/auth/*`

* Removed onboarding responsibilities from `AuthRepository`
* Kept authentication/session-related operations only:

  * login
  * logout
  * profile
  * session state
  * last login identity

#### `mobile/lib/data/repositories/registration/*`

* Added dedicated repository for:

  * CPF validation
  * registration submission

#### `mobile/lib/data/repositories/contact_verification/*`

* Added dedicated repository for:

  * requesting verification tokens
  * confirming verification tokens

#### `mobile/lib/data/repositories/register_draft/*`

* Updated imports to new cache module structure

### API layer refactor

#### `mobile/lib/data/services/apis/auth/*`

* Reduced `AuthApi` scope to:

  * login
  * profile retrieval

#### `mobile/lib/data/services/apis/registration/*`

* Added isolated registration API implementation
* Added CPF validation endpoint integration

#### `mobile/lib/data/services/apis/contact_verification/*`

* Added isolated contact verification API implementation
* Added request/confirm verification flows

### Cache module reorganization

#### Renamed:

* `data/services/auth/cache/*`
  → `data/services/cache/last_login/*`

This change removes the incorrect conceptual coupling between:

* authentication/session management
* local onboarding persistence/cache

The new structure better reflects the actual responsibility of the module.

### Dependency injection cleanup

#### `mobile/lib/data/services/services.dart`

#### `mobile/lib/data/repositories.dart`

#### `mobile/lib/domain/usecases/usecases.dart`

* Simplified injector registrations using constructor tear-offs
* Reduced boilerplate factory wiring
* Added lazy singleton registrations for:

  * `RegisterUsecase`
  * `RegistrationRepository`
  * `ContactVerificationRepository`
  * `RegistrationApi`
  * `ContactVerificationApi`

### Register flow improvements

#### `mobile/lib/domain/usecases/register/register_usecase.dart`

Refactored onboarding flow to use explicit specialized repositories:

* CPF validation now uses `RegistrationRepository`
* contact token request/confirmation now uses `ContactVerificationRepository`
* final registration submission now uses `RegistrationRepository`

This significantly improves:

* separation of concerns
* testability
* onboarding flow clarity
* future extensibility for MFA/ZTA-related onboarding validations

### DTO and package reorganization

Moved DTOs into explicit feature-oriented namespaces:

* `apis/auth/dtos/*`
* `apis/registration/dtos/*`
* `apis/contact_verification/dtos/*`

This removes the previous overloaded `auth/api/dtos` structure and aligns DTO ownership with their actual business capability.

### UI and routing updates

Updated imports and dependencies across:

* login
* short login
* splash
* register flow
* route extra codec

to reflect the new modular package organization.

### Register draft behavior updates

#### `mobile/test/domain/common/auth/models/register_draft_test.dart`

Adjusted serialization expectations:

* `current_step` is no longer serialized when unnecessary
* draft restoration behavior now accepts persisted snapshots more flexibly

### Test suite refactor

Added dedicated test coverage for:

* `ContactVerificationRepositoryImpl`
* `ContactVerificationApi`

Refactored:

* onboarding use case tests
* repository tests
* DTO tests
* register draft tests

Removed outdated tests that validated onboarding responsibilities inside `AuthRepository`.

### Overall result

This refactor significantly improves the architectural consistency of the mobile onboarding flow by:

* reducing responsibility concentration
* aligning APIs/repositories with business capabilities
* clarifying dependency boundaries
* simplifying dependency injection
* improving test isolation
* preparing the codebase for future onboarding evolution and Zero Trust-related identity validation flows.


## 2026/05/20 - mobile/pre-onboarding-07

Refactored the registration onboarding flow around a dedicated `RegisterUsecase`, introducing secure draft persistence with TTL-aware repository orchestration and preparing the application for the upcoming multi-page onboarding experience.

### Main Changes

1. Registration onboarding architecture redesign

   * Introduced `RegisterUsecase` as the central orchestration layer for the onboarding flow.
   * Moved registration business logic away from `RegisterViewmodel`.
   * Added explicit flow methods for:

     * CPF validation;
     * name submission;
     * birth date submission;
     * e-mail verification request/confirmation;
     * phone verification request/confirmation;
     * password submission;
     * final registration execution;
     * onboarding reset.
   * Added in-memory handling for verification tokens and password to avoid persistence of sensitive transient data.
   * Added initialization flow capable of restoring onboarding drafts from storage.

2. Secure onboarding draft persistence redesign

   * Reworked the onboarding persistence architecture separating:

     * storage responsibilities;
     * expiration rules;
     * orchestration concerns.
   * Added `RegisterDraftRepository`.
   * Added `RegisterDraftRepositoryImpl`.
   * Kept `RegisterDraftStore` focused exclusively on:

     * secure storage;
     * hashing;
     * JSON serialization/deserialization.
   * Moved TTL responsibility entirely to the repository layer.
   * Added 24-hour expiration logic with automatic cleanup of expired drafts.
   * Added injectable clock support (`DateTime Function()`) to improve deterministic testing.

3. Draft load result formalization

   * Standardized onboarding draft loading behavior through sealed-style results:

     * `RegisterDraftFound`;
     * `RegisterDraftNotFound`.
   * Simplified draft recovery semantics across the application.
   * Improved separation between storage absence and infrastructure failures.

4. Register draft model simplification

   * Removed `RegisterDraftStep` from persistence and state models.
   * Removed `currentStep` serialization/deserialization logic.
   * Simplified `RegisterDraftSnapshot`.
   * Simplified `RegisterDraftState`.
   * Added `RegisterDraftState.fromSnapshot`.
   * Reduced coupling between persisted onboarding state and UI navigation structure.

5. CPF validation integration

   * Added `cpfCheck` to `AuthRepository`.
   * Added implementation support in `AuthRepositoryImpl`.
   * Integrated CPF availability validation into onboarding initialization flow.
   * Prevented draft persistence for already registered CPFs.

6. Dependency injection and application wiring

   * Registered `RegisterDraftRepository` in dependency injection.
   * Registered `RegisterUsecase` in the usecase injector.
   * Connected `RegisterViewmodel` to the new usecase layer.
   * Added command-based execution wrappers to the viewmodel.

7. Routing and onboarding preparation

   * Renamed onboarding entry page:

     * `register_page.dart` → `register_cpf_page.dart`
   * Updated routing configuration to use the new onboarding entry structure.
   * Prepared the route organization for future multi-page onboarding screens.

8. Documentation and backlog updates

   * Updated onboarding backlog tasks to reflect the architectural redesign.
   * Clarified separation between:

     * persistence store;
     * repository orchestration;
     * TTL management.
   * Refined acceptance criteria and responsibilities for onboarding draft persistence.

9. Testing improvements

   * Added complete test coverage for `RegisterDraftRepositoryImpl`.
   * Added extensive `RegisterUsecase` tests covering:

     * onboarding initialization;
     * CPF validation;
     * verification flows;
     * registration execution;
     * reset behavior;
     * persistence cleanup;
     * failure propagation.
   * Updated existing onboarding draft tests after removal of `RegisterDraftStep`.
   * Removed obsolete onboarding UI flow tests tied to the previous architecture.

### Technical Highlights

* The onboarding persistence layer now follows a cleaner responsibility split:

  * Store → secure persistence only;
  * Repository → TTL and orchestration;
  * Usecase → business flow coordination.
* Sensitive onboarding runtime data such as verification tokens and passwords are intentionally kept only in memory.
* The registration flow is now significantly more testable due to:

  * explicit orchestration boundaries;
  * injectable time provider;
  * repository abstraction;
  * isolated state management.
* The onboarding architecture is now structurally prepared for:

  * multi-page navigation;
  * onboarding recovery;
  * future Zero Trust validations;
  * device binding flows;
  * advanced onboarding checkpoints.


## 2026/05/19 — mobile/pre-onboarding-06

Refactor and restructure the mobile onboarding flow into a new pre-onboarding architecture focused on incremental navigation, CPF validation, local persistence, and future-proof step orchestration.

### Main Changes

1. Introduced the new pre-onboarding routing structure

   * Added `RegisterRoutes` with dedicated routes for:

     * CPF
     * full name
     * birth date
     * email
     * email token
     * phone
     * phone token
     * password
     * password confirmation
     * success
   * Added `register_routes.dart`
   * Removed onboarding registration route from `AuthRoutes`
   * Updated login navigation to redirect to the new CPF-first onboarding entry point
   * Registered onboarding routes in the global router

2. Started the onboarding flow decomposition

   * Renamed `RegisterPage` to `RegisterCpfPage`
   * Replaced the monolithic onboarding implementation with an isolated initial page
   * Removed the previous multi-step widget orchestration from the UI layer
   * Prepared the project for fully decoupled onboarding screens and state transitions

3. Added CPF pre-validation support

   * Added `cpfCheck()` to `AuthRepositoryImpl`
   * Implemented `/auth/cpf-check` integration inside `AuthApi`
   * Added `CpfCheckResponseDto`
   * Added standardized HTTP and envelope parsing flow for CPF validation
   * Centralized API error handling and parsing behavior for CPF checks

4. Implemented onboarding draft persistence infrastructure

   * Added `RegisterDraftStore`
   * Added secure local onboarding persistence using `LocalSecureStorage`
   * Added SHA-256 CPF hashing to generate non-plain-text storage keys
   * Added automatic cleanup for corrupted onboarding snapshots
   * Added lookup abstraction with `RegisterDraftLookup`

5. Added onboarding domain models

   * Added:

     * `RegisterDraftField`
     * `RegisterDraftStep`
     * `RegisterDraftSnapshot`
     * `RegisterDraftState`
   * Implemented:

     * dirty field tracking
     * snapshot hydration
     * persistable state serialization
     * onboarding step persistence
     * verification metadata persistence
     * timestamp tracking
   * Explicitly excluded sensitive transient runtime data from persistence:

     * passwords
     * verification tokens
     * session tokens

6. Improved core utility extensions

   * Added:

     * `DateTime.dateOnly`
     * `DateParser.parseDateOnly`
     * `String.trimToNull()`
   * Simplified normalization and persistence serialization flows

7. Simplified the registration viewmodel

   * Removed the old monolithic registration orchestration
   * Removed:

     * internal step machine
     * verification orchestration
     * inline onboarding validations
     * direct register execution flow
   * Reduced the viewmodel to a lightweight injectable dependency base for future incremental onboarding modules

8. Added comprehensive onboarding draft tests

   * Added tests for:

     * CPF hash generation
     * secure storage persistence
     * snapshot serialization
     * invalid payload cleanup
     * dirty field tracking
     * hydration behavior
     * persisted timestamp updates
   * Added fake secure storage implementation for isolated testing

9. Removed obsolete onboarding tests

   * Deleted legacy `register_viewmodel_test.dart`
   * Removed tests tied to the previous monolithic onboarding implementation

10. Updated dependencies

* Added direct dependency:

  * `crypto: ^3.0.7`

11. Updated backlog organization

* Moved completed onboarding backlog documents to:

  * `docs/backlogs/mobile/done/`
* Archived:

  * `009 - pre_onboarding_contact_verification.md`
  * `009 - pre_onboarding_contact_verification_tasks.md`

12. Updated Postman environment

* Adjusted local `base_url` IP address for current development environment

This commit establishes the foundation for a modular onboarding architecture with persistent draft recovery, CPF-first flow validation, and isolated onboarding stages, reducing coupling between UI, orchestration, and verification logic while preparing the system for more advanced onboarding and Zero Trust validation flows.


## 2026/05/19 — mobile/pre-onboarding-05

This commit introduces the first structural foundation for the new multi-page onboarding flow, centered around CPF pre-validation, onboarding hardening, and local draft persistence planning.

The backend onboarding flow was expanded to support CPF availability checks before user registration, while the mobile backlog and onboarding architecture were redesigned around resumable multi-step registration.

### API — Introduce CPF pre-check onboarding endpoint

Implemented a new onboarding endpoint:

* `POST /auth/cpf-check`

Main goals:

* validate CPF format before registration;
* normalize CPF input;
* verify CPF availability before collecting the remaining onboarding data;
* block duplicated registrations early in the flow.

The endpoint now integrates directly into the AppToken-protected onboarding surface.

Main changes:

* Added `CheckCPFUseCase`
* Added `CustomerDocumentRepository.ExistsCPF`
* Added PostgreSQL CPF existence query
* Added HTTP handler `CheckCPF`
* Added route registration in `main.go`
* Added onboarding route protection tests
* Added error mapping for:

  * `ErrCPFRequired`
  * `ErrCPFInvalid`
* Exported `NormalizeCPF` from domain layer
* Added complete use case and handler test coverage

The onboarding entry flow is now:

1. CPF check
2. Contact verification
3. Register
4. Login

This establishes a cleaner onboarding boundary and avoids unnecessary verification flows for already registered users.

### API — Strengthen contact verification uniqueness rules

`RequestContactVerificationUseCase` was expanded to validate uniqueness before issuing verification challenges.

New behavior:

* email normalization now lowercases input before lookup;
* phone input is trimmed before lookup;
* duplicated email now returns `ErrEmailAlreadyExists`;
* duplicated phone now returns `ErrPhoneAlreadyExists`.

This prevents generating verification challenges for identities already associated with existing users.

Additional improvements:

* added normalization helper for verification targets;
* extended mocks and tests for uniqueness validation;
* improved verification test coverage for:

  * normalized email handling;
  * duplicated email;
  * duplicated phone;
  * invalid input;
  * repository invocation guarantees.

### API — Documentation expansion for onboarding surface

The REST and authentication documentation were heavily expanded to reflect the new onboarding sequence.

Main documentation additions:

* full `POST /auth/cpf-check` contract;
* onboarding flow update;
* onboarding AppToken clarification;
* CPF normalization behavior;
* new error scenarios;
* onboarding Postman flow update;
* onboarding security explanations;
* verification uniqueness rules;
* updated onboarding sequence references across architecture documents.

Files updated include:

* `07-api-rest.md`
* `08-auth_implementation.md`
* onboarding overview chapters
* onboarding flow references and error sections

The onboarding documentation now better reflects the intended progressive registration model.

### Mobile backlog — Multi-page onboarding architecture definition

The mobile onboarding backlog was substantially expanded and formalized.

Key architectural decisions documented:

* onboarding starts with CPF validation;
* onboarding becomes multi-page instead of a monolithic form;
* CPF availability must be checked before progressing;
* onboarding state becomes resumable;
* onboarding drafts are persisted in secure storage;
* onboarding uses CPF-hash-based storage keys;
* onboarding drafts expire after 24 hours;
* passwords and verification tokens are never persisted;
* dirty tracking is introduced for onboarding snapshots;
* onboarding remains specific to user registration for now;
* no generic onboarding engine will be created yet.

A complete execution breakdown was added through a new task file containing 12 implementation stages.

Topics covered:

* draft state modeling;
* secure storage persistence;
* TTL management;
* `RegisterViewmodel` orchestration redesign;
* multi-route onboarding navigation;
* CPF API integration;
* e-mail verification pages;
* phone verification pages;
* password flow;
* onboarding recovery;
* onboarding cleanup;
* migration away from the monolithic `RegisterPage`.

### Architectural impact

This commit significantly improves onboarding separation and prepares the project for:

* resumable onboarding;
* future KYC expansion;
* progressive onboarding UX;
* device-aware onboarding flows;
* Zero Trust evidence collection during onboarding;
* safer onboarding retries and recovery.

It also reinforces the onboarding boundary protected by `X-App-Token`, keeping the public surface explicit and controlled.

## 2026/05/18 - mobile/pre-onboarding-04

This commit advances the pre-onboarding foundation by restructuring authentication route registration in the API, strengthening middleware coverage through focused router tests, and defining the complete multi-page mobile onboarding strategy for the future registration flow.

### API

1. Refactored authentication route registration in `api/cmd/api/main.go`

   * Extracted auth route configuration into a dedicated `newAuthRouter` function.
   * Reduced bootstrap noise inside `main()`.
   * Improved separation between runtime wiring and route composition.
   * Centralized onboarding and authentication endpoint registration.
   * Prepared the auth router for future onboarding expansion and isolated testing.

2. Preserved onboarding security boundaries

   * Maintained `X-App-Token` protection for:

     * `POST /auth/contact-verifications`
     * `POST /auth/contact-verifications/confirm`
     * `POST /auth/register`
     * `POST /auth/login`
   * Preserved JWT protection for authenticated endpoints such as `/auth/me`.
   * Reinforced the current multi-stage authentication model already documented in the API architecture and auth documentation. 

### API Tests

3. Added onboarding middleware coverage tests in `api/cmd/api/routes_test.go`

   * Introduced focused tests validating AppToken enforcement in onboarding endpoints.
   * Added coverage for:

     * missing app token
     * invalid app token
     * valid app token pass-through behavior
   * Added reusable `assertInvalidAppToken` helper.
   * Validated:

     * HTTP 401 responses
     * standardized error envelope structure
     * `INVALID_APP_TOKEN` error code consistency
   * Reinforced contract stability for the onboarding security boundary.

4. Improved router-level validation strategy

   * Tests now validate middleware composition directly from router registration.
   * Reduced coupling between middleware expectations and handler implementation details.
   * Increased confidence in future onboarding route expansion.

### Mobile

5. Added onboarding backlog specification in `docs/backlogs/mobile/010 - cadastro_multi_paginas.md`

   * Defined the complete migration from a monolithic `RegisterPage` into a multi-step onboarding journey.
   * Formalized a 10-step registration flow:

     1. CPF
     2. Full name
     3. Birth date
     4. Email
     5. Email confirmation
     6. Phone
     7. Phone confirmation
     8. Password
     9. Password confirmation
     10. Account creation

6. Defined onboarding persistence strategy

   * Added secure local onboarding draft persistence.
   * Introduced CPF-derived storage keys using SHA-256 hashing.
   * Explicitly prohibited persistence of:

     * passwords
     * password confirmation
     * verification tokens
     * session tokens
   * Added onboarding recovery and resume strategy based on CPF lookup.
   * Defined draft expiration and verification invalidation behavior.

7. Defined mobile onboarding architecture decisions

   * Preserved `RegisterViewmodel` as the registration orchestrator.
   * Allowed `RegisterViewmodel` to remain a `lazySingleton`.
   * Defined shared state behavior across onboarding pages.
   * Established explicit separation between:

     * onboarding flow
     * verification flow
     * account creation
     * authentication/login flow

8. Defined onboarding/API integration contract

   * Specified use of:

     * `POST /auth/contact-verifications`
     * `POST /auth/contact-verifications/confirm`
     * `POST /auth/register`
   * Explicitly documented AppToken usage during onboarding requests.
   * Preserved development-mode debug token logging behavior in `AuthApi`.

9. Added onboarding acceptance criteria and scope boundaries

   * Documented:

     * persistence rules
     * validation expectations
     * navigation behavior
     * onboarding recovery rules
     * success flow requirements
   * Clearly separated:

     * current scope
     * future scope
     * non-goals
     * pending architectural decisions

### Mobile Tests

10. Updated route builder signatures in `mobile/test/ui/pages/auth/register/register_page_flow_test.dart`

* Adjusted `GoRoute` builders to use explicit `(context, state)` parameters.
* Improved consistency with current `go_router` conventions.
* Reduced ambiguity in future route evolution.

This commit establishes the architectural and testing foundation for the upcoming pre-onboarding implementation while reinforcing the security model around onboarding endpoints and preparing the mobile application for a stateful, resumable, multi-step registration experience.


## 2026/05/18 — mobile/pre-onboarding-03

This commit introduces the first complete pre-onboarding and contact verification flow for the mobile application, evolving the authentication experience from a simple registration form into a staged onboarding process with explicit e-mail and phone verification steps.

The implementation also improves authentication feedback behavior by making login errors context-aware, especially for partially verified accounts.

### Authentication Feedback Improvements

#### `mobile/lib/ui/pages/auth/login/login_page.dart`

* Added structured login feedback resolution through `_resolveLoginErrorMessage`.
* Introduced differentiated feedback for:

  * pending account approval
  * missing e-mail verification
  * missing phone verification
  * both contact channels pending verification
* Added support for reading `error.details` to interpret backend verification state.
* Centralized login failure message resolution instead of directly exposing raw backend messages.
* Improved UX consistency for authentication failures.

#### `mobile/lib/ui/pages/auth/short_login/short_login_page.dart`

* Mirrored the same verification-aware login feedback logic implemented in the full login flow.
* Preserved remembered identity context while presenting verification-related authentication errors.
* Improved short-login usability for partially onboarded users.

### Multi-Step Registration Flow

#### `mobile/lib/ui/pages/auth/register/register_page.dart`

* Reworked the registration screen into a multi-step onboarding flow:

  * Personal data
  * Contact data
  * E-mail verification
  * Phone verification
  * Final review
* Added visual onboarding step indicators.
* Introduced step-aware navigation and validation.
* Added support for:

  * birth date selection
  * Brazilian phone formatting
  * phone normalization for API communication
  * review/confirmation step before submission
* Added command orchestration for:

  * requesting e-mail verification codes
  * confirming e-mail verification codes
  * requesting phone verification codes
  * confirming phone verification codes
* Added reusable snackbar feedback using `AppSnackbar`.
* Added state persistence between steps.
* Added conditional primary action behavior based on onboarding stage.
* Introduced command aggregation through `Listenable.merge`.
* Added UI locking while asynchronous operations are running.
* Added dynamic CTA labels according to current onboarding step.
* Added verification-aware enable/disable behavior for onboarding progression.
* Added custom `_BrazilPhoneInputFormatter`.
* Added birth date validation and date picker integration.
* Added review screen summarizing collected onboarding data before final registration.

### Registration ViewModel Refactor

#### `mobile/lib/ui/pages/auth/register/viewmodel/register_viewmodel.dart`

* Refactored the viewmodel into a full `ChangeNotifier` state machine.
* Introduced `RegisterStep` enum for explicit onboarding progression.
* Added internal onboarding state management:

  * personal data
  * contact data
  * verification identifiers
  * verification tokens
  * step errors
* Added onboarding progression methods:

  * `nextStep`
  * `previousStep`
  * `goToStep`
* Added validation-aware onboarding transitions.
* Added explicit onboarding state guards.
* Added e-mail verification orchestration.
* Added phone verification orchestration.
* Added final registration orchestration dependent on verification tokens.
* Added command-based async flow integration:

  * `Command0`
  * `Command1`
* Added validation helpers for:

  * e-mail
  * phone
  * CPF
  * onboarding state
* Added step-specific error propagation and recovery behavior.
* Added support for preserving entered data across onboarding navigation.
* Added registration payload enrichment with:

  * `birthDate`
  * `phone`
  * `emailVerificationToken`
  * `phoneVerificationToken`

### Authentication Feedback Tests

#### `mobile/test/ui/pages/auth/login_feedback_behavior_test.dart`

* Added tests covering:

  * generic contact-not-verified login feedback
  * e-mail-only pending verification feedback
  * phone-only pending verification feedback
  * short-login verification feedback behavior
* Added verification-state-aware error fixtures using `error.details`.

### Registration Flow Widget Tests

#### `mobile/test/ui/pages/auth/register/register_page_flow_test.dart`

* Added end-to-end widget test covering the entire onboarding process.
* Validated:

  * multi-step navigation
  * date picker integration
  * phone formatting
  * e-mail verification flow
  * phone verification flow
  * review step
  * final registration submission
* Added verification token assertions in final payload validation.
* Added navigation validation after successful registration.

### Registration ViewModel Tests

#### `mobile/test/ui/pages/auth/register/viewmodel/register_viewmodel_test.dart`

* Added unit tests validating:

  * initial onboarding state
  * invalid step progression blocking
  * onboarding data persistence
  * e-mail verification transitions
  * phone verification transitions
  * verification token requirements
  * final registration payload generation
* Added coverage for onboarding state machine behavior.

### Architectural Notes

This commit significantly advances the onboarding architecture toward a more realistic fintech onboarding model aligned with the project's Zero Trust and evidence-based authentication direction.

The new flow establishes the foundation for:

* progressive onboarding
* contextual trust evaluation
* stronger identity validation
* future device registration flows
* transactional security expansion
* adaptive authentication mechanisms

It also moves the mobile application closer to the backend authentication strategy already documented in the API architecture and authentication specifications.  


## 2026/05/18 — mobile/pre-onboarding-02

This commit introduces the first structural layer of the pre-onboarding flow for the mobile application, focusing on contact verification before account registration and improving the error propagation model between API and Flutter client.

The implementation establishes the initial foundation for a more robust onboarding process, where e-mail and phone validation become explicit prerequisites before authentication and account activation flows.

### Backlog and Project Organization

1. Reorganized onboarding backlog numbering

   * Renamed:

     * `docs/backlogs/api/006 - onboarding.md`
     * → `docs/backlogs/api/001 - onboarding.md`
   * Normalized onboarding backlog priority within the API planning structure.
   * Reinforced onboarding as one of the primary architectural flows of the platform.

### Contact Verification Flow Foundation

1. Added contact verification support to the authentication repository

   * Introduced:

     * `requestContactVerification(...)`
     * `confirmContactVerification(...)`
   * Added DTO integration for:

     * verification request
     * verification confirmation
     * verification tokens
     * verification responses

2. Extended repository implementation

   * Added delegation to `AuthApi` for verification operations.
   * Preserved authentication session isolation:

     * verification operations do not create sessions
     * verification operations do not persist login state
     * verification operations do not modify cached profile state

3. Added support for verification tokens during registration

   * `RegisterRequestDto` now includes:

     * `phone`
     * `birthDate`
     * `emailVerificationToken`
     * `phoneVerificationToken`
   * Added CPF normalization before serialization.
   * Added birth date serialization using ISO date-only format.
   * Added parsing validation for `birth_date`.

### HTTP Error Mapping Evolution

1. Introduced new application error code

   * Added:

     * `AppErrorCode.contactNotVerified`

2. Extended Dio error mapper

   * Added backend mapping for:

     * `CONTACT_NOT_VERIFIED`
   * Implemented contextual user feedback based on verification state:

     * both channels pending
     * only e-mail pending
     * only phone pending

3. Improved backend error details propagation

   * Added support for:

     * `error.details`
   * Preserved generic compatibility for legacy error payloads.
   * Improved HTTP error fallback handling.

4. Extended API envelope parsing

   * `ApiError` now supports structured `details`.
   * Added safe parsing for dynamic map payloads.

### Register Flow Transition

1. Adjusted register page behavior

   * Removed direct registration execution from the page.
   * Temporarily redirected flow toward pre-verification guidance.
   * Added onboarding feedback snackbar:

     * “Confirme seu e-mail e telefone antes de concluir o cadastro.”

2. Prepared the UI for the upcoming staged onboarding flow

   * The register page now behaves as an entry point for future:

     * e-mail confirmation
     * SMS confirmation
     * pre-auth onboarding orchestration

### Test Coverage Expansion

1. Added exhaustive tests for contact verification error mapping

   * Covered scenarios:

     * both channels pending
     * e-mail only pending
     * phone only pending
     * missing details fallback
     * generic 403 isolation

2. Added repository tests for verification flows

   * Verified:

     * API delegation
     * session isolation
     * failure propagation
     * success propagation

3. Added DTO serialization and parsing tests

   * Verified:

     * CPF normalization
     * token persistence
     * birth date formatting
     * payload compatibility

4. Added API envelope compatibility tests

   * Covered:

     * envelopes with details
     * envelopes without details
     * backward compatibility behavior

This commit establishes the first concrete layer of a staged onboarding architecture in the mobile application, moving the project away from immediate account creation and toward a verification-first model aligned with stronger identity validation and future Zero Trust onboarding flows.


## 2026/05/18 — mobile/pre-onboarding-01

This commit introduces the first mobile foundation for the new pre-onboarding flow with contact verification support. The work aligns the Flutter client with the updated authentication contract from the API, where users must verify both e-mail and phone before completing registration.

A major focus of this commit was the backlog reorganization and the creation of a dedicated mobile onboarding specification covering the new multi-step registration journey.

### Documentation and backlog restructuring

1. Reorganized onboarding backlog numbering and lifecycle

   * Renamed API onboarding backlog identifiers to maintain sequence consistency.
   * Moved completed pre-onboarding backlog items into the `done/` directory.
   * Preserved historical references and backlog traceability.

2. Added the mobile pre-onboarding specification

   * Created:

     * `docs/backlogs/mobile/009 - pre_onboarding_contact_verification.md`
     * `docs/backlogs/mobile/009 - pre_onboarding_contact_verification_tasks.md`
   * Documented:

     * new API contracts;
     * contact verification flow;
     * UI and ViewModel responsibilities;
     * login behavior for `CONTACT_NOT_VERIFIED`;
     * DTO and repository expectations;
     * acceptance criteria for all implementation phases.
   * Defined the complete multi-step registration journey:

     * e-mail verification;
     * phone verification;
     * birth date collection;
     * final registration with verification tokens.
   * Added a structured task breakdown with dependency mapping for incremental implementation.

### Mobile auth API integration

3. Added contact verification support to `AuthApi`

   * Implemented:

     * `requestContactVerification()`
     * `confirmContactVerification()`
   * Added integration with:

     * `POST /auth/contact-verifications`
     * `POST /auth/contact-verifications/confirm`
   * Added `X-App-Token` propagation for onboarding endpoints.
   * Added HTTP status validation and envelope parsing.
   * Preserved the existing `Result<AppError>` flow.
   * Added structured logging for:

     * request failures;
     * API envelope failures;
     * parsing failures.
   * Added development-mode debug logging for returned verification tokens.

### Contact verification DTOs

4. Added request and response DTOs for contact verification

   * Created:

     * `ContactVerificationRequestDto`
     * `ContactVerificationRequestResponseDto`
     * `ContactVerificationConfirmRequestDto`
     * `ContactVerificationConfirmResponseDto`
   * Added:

     * request serialization;
     * response parsing;
     * verification token mapping;
     * ISO datetime parsing for expiration and verification timestamps.

### Automated tests

5. Added DTO unit tests

   * Added serialization and parsing coverage for all new DTOs.
   * Validated:

     * request payload generation;
     * verification response parsing;
     * datetime conversion behavior.

6. Added `AuthApi` integration-oriented tests

   * Added coverage for:

     * request verification endpoint;
     * confirmation endpoint;
     * API error envelope mapping;
     * `X-App-Token` propagation;
     * request body generation.
   * Added fake REST client infrastructure for isolated testing.

### Architectural alignment

7. Advanced the mobile onboarding architecture toward the new API contract

   * Established the initial infrastructure required for:

     * multi-step onboarding;
     * pre-registration verification;
     * channel-specific verification handling;
     * future onboarding orchestration in `RegisterViewmodel`.
   * Prepared the mobile layer for upcoming support of:

     * `CONTACT_NOT_VERIFIED`;
     * `error.details`;
     * redesigned registration UI;
     * onboarding state management.

This commit represents the first operational mobile step toward the new onboarding model and establishes the API integration base required for the future multi-stage registration experience.


## 2026/05/18 — api/pre-onboarding-06

This commit advances the pre-onboarding and authentication flow by introducing contact verification enforcement during login, removing the remaining direct CPF dependency from the `customers` table, and consolidating document handling through `customer_documents`.

### Authentication and Contact Verification

Implemented mandatory contact verification validation before session creation during login.

#### `api/internal/auth/application/login_user.go`

* Added `validateContactVerification` to enforce verified email and phone before authentication completes.
* Introduced early interruption of the login flow when contact verification requirements are not satisfied.
* Ensured validation happens before account provisioning checks and token generation.

#### `api/internal/auth/domain/errors.go`

* Added `ErrContactNotVerified`.
* Introduced `ContactNotVerifiedError` carrying structured verification state:

  * `EmailVerified`
  * `PhoneVerified`
* Implemented `Unwrap()` support for compatibility with `errors.Is` and `errors.As`.

#### `api/internal/auth/application/errors_registry.go`

* Added structured registration for `CONTACT_NOT_VERIFIED`.
* Introduced dynamic error detail mapping for verification state exposure.

#### `api/internal/shared/errors/*`

* Extended `AppError` with `Details`.
* Added `RegisterDomainErrorWithDetails`.
* Updated mapper logic to dynamically populate structured error payloads.
* Added new error code:

  * `CONTACT_NOT_VERIFIED`

#### `api/internal/shared/http/response.go`

* Extended HTTP error responses to include:

  * `error.details`

### Login Flow Tests

Expanded login coverage to validate all verification scenarios.

#### `api/internal/auth/application/login_user_test.go`

* Updated all successful authentication scenarios to include:

  * `EmailVerifiedAt`
  * `PhoneVerifiedAt`
* Added tests for:

  * email not verified
  * phone not verified
* Verified that failed verification:

  * blocks token generation
  * blocks session persistence
  * skips provisioning validation

#### `api/internal/auth/delivery/handler_test.go`

* Added HTTP integration coverage for:

  * `CONTACT_NOT_VERIFIED`
* Validated:

  * HTTP 403 status
  * structured error response
  * verification detail payload

### Customer CPF Decoupling

One of the main goals of this commit was removing the remaining runtime dependency on the legacy `customers.cpf` column and consolidating CPF access through `customer_documents`.

#### `api/internal/account/bankaccount/infrastructure/repository.go`

* Replaced direct `customers.cpf` usage with joins against `customer_documents`.
* Updated transfer recipient lookup queries to:

  * resolve CPF from `customer_documents`
  * enforce:

    * `type = 'cpf'`
    * `country = 'BR'`
    * `is_primary = true`

#### `api/internal/account/bankaccount/infrastructure/repository_test.go`

* Added guards ensuring:

  * no legacy `customers.cpf` references remain
  * queries use `customer_documents`
  * document filtering uses `cd.value`

### Database Migration

Introduced migration removing the direct CPF column from `customers`.

#### `api/migrations/000008_remove_customers_direct_cpf.up.sql`

* Removed:

  * `customers.cpf`
  * CPF format constraint
  * CPF unique constraint

#### `api/migrations/000008_remove_customers_direct_cpf.down.sql`

* Recreated:

  * `customers.cpf`
  * CPF constraints
* Added restoration logic using `customer_documents`.

### Integration and Repository Test Schema Updates

Adjusted integration schemas and seed helpers to align with the new document model.

#### Updated test schemas

* `deposit_integration_test.go`
* `auth_authorization_integration_test.go`
* `postgres_user_repository_test.go`

Changes include:

* removal of direct CPF column usage
* creation of `customer_documents`
* primary document uniqueness index
* document seeding helpers
* explicit CPF document insertion

### User Seed and Verification Improvements

#### `api/internal/auth/delivery/auth_authorization_integration_test.go`

* Updated seeded users to include:

  * phone number
  * `email_verified_at`
  * `phone_verified_at`
* Ensured integration login flows reflect real verification requirements.

### Architectural Impact

This commit significantly advances the onboarding and identity model by:

* separating customer identity documents from the customer aggregate itself
* enabling multi-document support in future onboarding flows
* removing direct CPF coupling from transfer operations
* introducing structured authentication failure responses
* enforcing verified communication channels before authenticated access
* improving API consistency around domain-driven error propagation

The changes also prepare the platform for future onboarding expansion involving:

* multiple document types
* internationalization of identity records
* staged onboarding and verification workflows
* richer Zero Trust and pre-authentication validation strategies.


## 2026/05/18 — api/pre-onboarding-05

This commit advances the pre-onboarding and authentication flow by consolidating customer document ownership into `customer_documents`, introducing mandatory contact verification before login, and extending the shared error infrastructure to support structured error details.

### Customer document normalization and CPF decoupling

The persistence layer was refactored to remove direct CPF dependencies from the `customers` table and fully adopt the `customer_documents` model for identity resolution.

#### Account repository updates

* Refactored transfer recipient queries to load CPF values through `customer_documents`
* Removed direct references to `customers.cpf`
* Added filtering rules for:

  * `type = 'cpf'`
  * `country = 'BR'`
  * `is_primary = true`
* Preserved active account filtering semantics

Files:

* `api/internal/account/bankaccount/infrastructure/repository.go`

#### Repository test hardening

Added assertions to guarantee:

* no legacy CPF column usage
* mandatory `customer_documents` joins
* correct filtering through `cd.value`

Files:

* `api/internal/account/bankaccount/infrastructure/repository_test.go`

### Integration test schema migration alignment

Updated multiple integration schemas and seed helpers to reflect the new normalized document structure.

#### Deposit integration schema

* Removed direct CPF column from `customers`
* Added full `customer_documents` table
* Added:

  * uniqueness constraints
  * primary document partial index
  * timestamps
* Refactored customer seeding helpers to insert CPF documents separately

Files:

* `api/internal/account/transaction/delivery/deposit_integration_test.go`

#### Auth integration schema

* Removed CPF coupling from customer schema
* Added explicit CPF document seeding
* Updated seeded users to include:

  * `phone`
  * `email_verified_at`
  * `phone_verified_at`

Files:

* `api/internal/auth/delivery/auth_authorization_integration_test.go`

#### Auth repository schema cleanup

Removed obsolete CPF column assumptions from auth repository tests.

Files:

* `api/internal/auth/infrastructure/postgres_user_repository_test.go`

### Login contact verification enforcement

Introduced a mandatory verification gate before login token issuance.

#### New login validation flow

Added:

* `validateContactVerification(user)`

The login process now blocks authentication when:

* email is not verified
* phone is not verified

Validation occurs before:

* account provisioning checks
* token generation
* session persistence

Files:

* `api/internal/auth/application/login_user.go`

### Structured domain error for contact verification

Introduced a dedicated domain error type capable of exposing verification state metadata.

#### New domain error model

Added:

* `ErrContactNotVerified`
* `ContactNotVerifiedError`
* structured verification state:

  * `EmailVerified`
  * `PhoneVerified`

Implemented:

* `Error()`
* `Unwrap()`
* constructor helper

Files:

* `api/internal/auth/domain/errors.go`

### Shared error infrastructure improvements

Extended the shared error system to support dynamic structured payloads.

#### Shared error model updates

Added:

* `AppError.Details`

Implemented:

* `RegisterDomainErrorWithDetails(...)`
* dynamic detail extraction callbacks in mapper registry

Files:

* `api/internal/shared/errors/error.go`
* `api/internal/shared/errors/mapper.go`

#### HTTP error serialization

Extended HTTP responses to serialize:

* `error.details`

Files:

* `api/internal/shared/http/response.go`

#### New shared error code

Added:

* `CONTACT_NOT_VERIFIED`

Files:

* `api/internal/shared/errors/codes.go`

### Auth error registry enhancements

Registered structured mapping for contact verification failures.

Behavior:

* returns HTTP `403 Forbidden`
* includes:

  * `email_verified`
  * `phone_verified`

Files:

* `api/internal/auth/application/errors_registry.go`

### Login use case test coverage expansion

Refactored existing tests to seed verified users explicitly and added dedicated verification failure scenarios.

#### Added test coverage for

* email not verified
* phone not verified
* early interruption of login flow
* prevention of:

  * token generation
  * session creation
  * provisioning checks

Files:

* `api/internal/auth/application/login_user_test.go`

### HTTP handler verification error coverage

Added explicit response validation for:

* `CONTACT_NOT_VERIFIED`
* structured `details` payload

Files:

* `api/internal/auth/delivery/handler_test.go`

### Database migration for CPF removal

Added migration `000008_remove_customers_direct_cpf`.

#### Up migration

* removes:

  * `customers.cpf`
  * CPF constraints
  * CPF unique index

#### Down migration

* recreates `customers.cpf`
* restores CPF values from `customer_documents`
* recreates:

  * format validation constraint
  * uniqueness constraint

Files:

* `api/migrations/000008_remove_customers_direct_cpf.up.sql`
* `api/migrations/000008_remove_customers_direct_cpf.down.sql`

### Architectural impact

This commit significantly advances the separation between:

* customer identity
* customer documents
* authentication lifecycle

It also introduces the first structured authentication gating mechanism tied to onboarding state progression, creating a cleaner foundation for:

* multi-document support
* country-specific identification strategies
* progressive onboarding
* pre-account verification workflows
* future Zero Trust and KYC enforcement layers

The new shared error infrastructure also establishes a reusable mechanism for exposing structured domain context in API responses without leaking internal implementation details.


## 2026/05/18 — api/pre-onboarding-04

This commit evolves the onboarding and registration flow with a new pre-onboarding contact verification stage, introducing verified e-mail and phone validation before user creation. The implementation expands the authentication domain, persistence model, HTTP contracts, and integration coverage to support a more realistic onboarding pipeline aligned with future security and Zero Trust directions.

### Main Changes

1. Registration flow now requires verified e-mail and phone tokens

   * Extended `RegisterUserUseCase` to receive a `ContactVerificationRepository`
   * Added support for:

     * `phone`
     * `email_verification_token`
     * `phone_verification_token`
   * Added validation logic for:

     * verified e-mail token ownership
     * verified phone token ownership
     * channel consistency
     * verified state existence
   * Added phone uniqueness validation through `ExistsByPhone`
   * Persisted:

     * `phone`
     * `email_verified_at`
     * `phone_verified_at`
   * Improved registration transaction flow to validate contact verification state before user creation

2. Contact verification infrastructure expanded

   * Added:

     * `FindContactVerificationByVerificationToken`
   * Implemented PostgreSQL lookup by verification token
   * Added integration between:

     * verification confirmation
     * registration flow
   * Improved repository documentation and internal scanning helpers

3. User domain and persistence model expanded

   * `User` entity now supports:

     * `Phone`
     * `EmailVerifiedAt`
     * `PhoneVerifiedAt`
   * PostgreSQL repository updated to:

     * insert optional phone and verification timestamps
     * load nullable verification fields
     * support phone existence checks
   * Added helper utilities for nullable string/time persistence

4. Database schema updated for pre-onboarding verification support

   * Added columns to `users`:

     * `phone`
     * `email_verified_at`
     * `phone_verified_at`
   * Added new `contact_verifications` table
   * Added:

     * verification token uniqueness index
     * target/channel index
   * Updated integration and repository test schemas accordingly

5. HTTP authentication contract updated

   * `/auth/register` now requires:

     * phone
     * e-mail verification token
     * phone verification token
   * Registration request validation expanded
   * Added explicit handler documentation for:

     * contact verification request
     * contact verification confirmation

6. Integration tests now execute the full pre-onboarding flow

   * Added helper:

     * `requestAndConfirmContactVerification`
   * Integration flow now performs:

     * verification request
     * verification confirmation
     * token retrieval
     * final registration using verified tokens
   * Registration integration tests now simulate realistic onboarding behavior

7. Test suite updated across authentication and account modules

   * Added mock support for:

     * `ExistsByPhone`
     * verification token lookups
   * Added registration tests for:

     * duplicate phone
     * missing verification tokens
     * verified timestamp persistence
     * invalid phone scenarios
   * Updated all registration-related tests to use the new onboarding contract

8. Error handling expanded

   * Added:

     * `ErrPhoneAlreadyExists`
   * Registered HTTP conflict mapping for duplicate phone scenarios

### Architectural Impact

This commit introduces an important transition in the onboarding model:

* onboarding becomes stateful before account creation
* contact ownership is now validated independently from user persistence
* user registration becomes dependent on previously validated evidence
* authentication flow starts evolving toward evidence-based onboarding

The implementation also strengthens separation between:

* verification lifecycle
* onboarding orchestration
* persistence concerns
* HTTP delivery contracts

This lays the groundwork for future onboarding stages such as:

* device registration
* transactional password setup
* liveness verification
* contextual risk analysis
* Zero Trust evidence aggregation

Relevant architecture and authentication references remain aligned with the current project documentation.  


## 2026/05/18 - api/pre-onboarding-03

Implemented the first pre-onboarding contact verification flow, introducing a dedicated verification lifecycle for email and phone channels before user registration completion.

### Added contact verification domain model and repository contracts

* Introduced the `ContactVerification` entity with:

  * verification channel normalization
  * expiration handling
  * verification token support
  * validation rules for email and phone channels
* Added `ContactVerificationRepository` interface to the auth domain layer.
* Expanded auth domain errors with:

  * `ErrContactVerificationNotFound`
  * `ErrInvalidVerificationToken`
  * `ErrContactVerificationExpired`

### Implemented contact verification use cases

Added two new application use cases:

1. `RequestContactVerificationUseCase`
2. `ConfirmContactVerificationUseCase`

Main behaviors implemented:

* generation of secure numeric verification tokens
* verification expiration control with TTL
* verification token issuance after successful confirmation
* validation of verification channel and target
* verification state validation before confirmation

The implementation also introduced:

* cryptographically secure numeric token generation
* UTC-based deterministic time handling
* normalized channel processing (`email` and `phone`)

### Added PostgreSQL contact verification repository

Implemented `PostgresContactVerificationRepository` with support for:

* verification creation
* lookup by verification ID
* verification confirmation persistence
* transaction-aware executor support
* nullable verification metadata mapping
* PostgreSQL constraint error translation to domain errors

### Added database migration for contact verifications

Created migration `000007_contact_verifications` with:

* `contact_verifications` table
* channel constraint validation
* target/channel lookup index
* unique verification token index
* expiration and verification metadata fields

### Integrated pre-onboarding verification into API bootstrap

Updated `cmd/api/main.go` to:

* instantiate the new repository
* wire the new use cases
* inject verification flows into the auth handler
* expose new onboarding endpoints

New endpoints added:

* `POST /auth/contact-verifications`
* `POST /auth/contact-verifications/confirm`

Both endpoints are protected with the existing `AppToken` middleware, preserving the onboarding boundary architecture already established in the authentication model. 

### Extended auth delivery layer

Updated the auth handler to support:

* request payload parsing
* verification confirmation parsing
* UUID validation
* response serialization
* centralized error mapping
* structured error logging

Added dedicated request DTOs for:

* contact verification request
* contact verification confirmation

### Expanded automated test coverage

Added comprehensive unit coverage for:

* successful verification request flow
* invalid verification input
* successful verification confirmation
* invalid verification token handling

Expanded HTTP handler tests for:

* request endpoint success path
* confirmation endpoint success path
* response payload validation
* request parsing validation

Updated existing integration wiring tests to support the expanded auth handler constructor.

### Architectural impact

This commit establishes the first reusable onboarding verification primitive inside the authentication module.

The implementation keeps the architecture aligned with the current layered modular monolith approach:

* delivery → application → domain
* infrastructure → domain

while preserving explicit runtime composition through `cmd/api/main.go`. 

This also creates the foundation for future onboarding flows such as:

* email ownership validation
* phone validation
* transactional onboarding checkpoints
* device registration
* Zero Trust onboarding signals
* multi-step identity verification workflows

The solution was intentionally implemented as an isolated verification lifecycle instead of coupling verification state directly into the user entity, allowing future reuse across independent onboarding and security flows.


## 2026/05/18 — api/pre-onboarding-02

This commit advances the pre-onboarding cadastral redesign by introducing the first structural migrations for contact verification, customer documents, and customer addresses, while also formalizing the migration strategy and implementation backlog for the new onboarding model.

### Database migrations

1. Added `000004_pre_onboarding_cadastral` migration

   * Extended `users` with:

     * `phone`
     * `email_verified_at`
     * `phone_verified_at`
   * Added `birth_date` to `customers`
   * Introduced the new `customer_documents` table
   * Introduced the new `customer_addresses` table
   * Added unique document constraint:

     * `(type, value, country)`
   * Added partial unique indexes enforcing:

     * a single primary document per customer
     * a single primary address per customer
   * Standardized country handling using:

     * `CHAR(2) DEFAULT 'BR'`
   * Added rollback support for all new structures

2. Added `000005_migrate_customer_cpf_documents` migration

   * Migrated existing `customers.cpf` data into `customer_documents`
   * Preserved customer creation timestamps during migration
   * Marked migrated CPF entries as primary documents
   * Added defensive `NOT EXISTS` filtering for safer local retries
   * Added rollback migration removing migrated CPF document entries

### Pre-onboarding architecture and migration strategy

3. Expanded the pre-onboarding backlog documentation

   * Added an explicit migration order strategy to avoid:

     * intermediate schema breakage
     * data loss during CPF extraction
   * Clarified that `customers.cpf` removal must happen only after:

     * schema preparation
     * data migration
     * repository/query migration
     * test adaptation
   * Documented the transition from direct CPF ownership in `Customer`
     to document-based identity modeling

4. Added impact analysis checklist for the current codebase

   * Identified affected layers and repositories:

     * customer domain
     * auth registration flow
     * account lookup queries
     * transfer recipient queries
     * integration tests
     * Postman collection
     * technical documentation
   * Documented required query migrations from:

     * `customers.cpf`
       to:
     * `customer_documents`

### Backlog decomposition and implementation planning

5. Added `000 - pre-onboarding_tasks.md`

   * Created a detailed 12-task execution plan for the onboarding redesign
   * Structured the migration into isolated phases covering:

     * schema preparation
     * CPF migration
     * customer document modeling
     * contact verification flows
     * auth updates
     * login verification requirements
     * repository migration
     * final CPF column removal
     * documentation/test updates
   * Added:

     * acceptance criteria
     * dependency chains
     * suggested execution order
   * Explicitly separated:

     * domain responsibilities
     * persistence migration
     * onboarding/session evolution
     * verification workflows

### Architectural direction reinforced

6. Reinforced the transition toward a more extensible identity model

   * Customer identity is no longer tied exclusively to CPF
   * The system now evolves toward:

     * multi-document support
     * country-aware document modeling
     * onboarding verification checkpoints
     * future KYC extensibility
   * Prepared the foundation for:

     * email verification
     * phone verification
     * onboarding-scoped sessions
     * future Zero Trust onboarding flows

This commit establishes the structural base for the new onboarding architecture while preserving migration safety and preparing the codebase for the gradual removal of direct CPF coupling from the `Customer` model.


## 2026/05/18 - api/pre-onboarding-01

Introduced the initial pre-onboarding cadastral foundation for the Bank API, focusing on evolving the customer identity model from a CPF-centric structure to a flexible document-based architecture. This commit establishes the database preparation, migration strategy, and execution roadmap required to support future onboarding checkpoints, contact verification, and expanded customer identity flows.

### Database Migrations

1. Added `000004_pre_onboarding_cadastral` migration

   * Extended `users` table with:

     * `phone`
     * `email_verified_at`
     * `phone_verified_at`
   * Added `birth_date` to `customers`
   * Introduced `customer_documents` table

     * support for CPF, RG, CNH, and future document types
     * unique document constraint using `(type, value, country)`
     * optional issuer metadata
     * partial unique index enforcing a single primary document per customer
   * Introduced `customer_addresses` table

     * normalized address structure
     * partial unique index enforcing a single primary address per customer
   * Added reversible down migration for the entire cadastral preparation layer

2. Added `000005_migrate_customer_cpf_documents` migration

   * Migrated existing `customers.cpf` values into `customer_documents`
   * Preserved customer creation timestamps during migration
   * Marked migrated CPF documents as primary
   * Added rollback migration removing migrated CPF documents safely
   * Designed migration flow to avoid intermediate data loss during transition

### Backlog and Architectural Planning

3. Expanded `docs/backlogs/api/000 - pre-onboarding.md`

   * Documented the mandatory migration order to safely remove `customers.cpf`
   * Added impact checklist covering:

     * customer domain
     * account repository
     * auth registration flow
     * transfer recipient queries
     * tests and Postman collection
   * Clarified that `customer_documents` must become the authoritative identity source before removing legacy CPF fields
   * Reinforced transactional and migration safety concerns even for local development environments

4. Added `docs/backlogs/api/000 - pre-onboarding_tasks.md`

   * Introduced a complete execution roadmap with 12 detailed tasks
   * Defined:

     * objectives
     * scope
     * acceptance criteria
     * dependencies
     * suggested execution order
   * Covered:

     * domain refactoring
     * customer document contracts
     * repository migration
     * onboarding verification flow
     * login verification rules
     * transfer lookup updates
     * final CPF field removal
     * documentation and testing updates

### Architectural Direction

5. Established the foundation for a normalized customer identity model

   * Decouples `Customer` from direct CPF ownership
   * Enables multi-document support per customer
   * Creates a clear path for future onboarding checkpoints and KYC evolution
   * Preserves transactional consistency during migration and future onboarding flows
   * Prepares the authentication model for verified email and phone requirements

### Migration Strategy Highlights

6. Formalized the transitional migration approach

   * Create new structures before removing legacy fields
   * Copy existing CPF data into normalized tables
   * Refactor application code and queries incrementally
   * Remove `customers.cpf` only after all dependencies are eliminated
   * Preserve operational continuity and rollback safety throughout the transition

This commit represents the first structural step toward a richer onboarding and identity verification model, while maintaining compatibility with the current transactional architecture and existing authentication flow.


## 2026/05/18 — docs/update-17

Expanded the project presentation and documentation structure to better position BankLab as an open-source financial systems engineering laboratory, while improving onboarding material for contributors and international readers.

### Main Changes

1. `README.md`

   * Reworked the Portuguese project introduction with a stronger positioning around:

     * financial systems engineering;
     * transactional consistency;
     * onboarding flows;
     * future Zero Trust Architecture evolution.
   * Added a more narrative explanation describing the motivation behind the project and its relationship with research, architecture decisions, and practical experimentation.
   * Improved the explanation of the ledger-based financial model and the role of the `transactions` table as the source of truth.
   * Expanded the repository description to reinforce the educational and collaborative goals of the project.
   * Reformatted and improved the “Current Scope” table for better readability.
   * Updated the repository structure section to reflect the current monorepo organization, including:

     * `CHANGELOG.md`
     * `CONTRIBUTING.md`
     * `README_en.md`
     * `packages/`
     * `templates/`
   * Added the new LinkedIn/project presentation cover image to the README.
   * Improved wording and organization throughout the document to make the repository more approachable for external contributors.

2. `README_en.md`

   * Added a complete English version of the main project documentation.
   * Included:

     * project overview;
     * architectural principles;
     * current scope;
     * stack description;
     * repository structure;
     * local setup instructions;
     * endpoint overview;
     * contribution guidelines;
     * documentation index;
     * project status and license sections.
   * Standardized the public-facing presentation of the project for international audiences and potential collaborators.
   * Positioned the project as a practical study environment for backend, mobile, financial systems, and Zero Trust-related experimentation.

3. `docs/images/Capa_BankLab_LinkedIn_v3.png`

   * Added the new project visual identity image used in the README and public presentation materials.

4. `mobile/assets/images/banklab.svg`

   * Refactored and reorganized the SVG structure generated by Inkscape.
   * Grouped logo elements into a more structured hierarchy.
   * Preserved the existing visual identity while improving organization of the vector asset.

### Documentation and Positioning Improvements

This commit significantly improves how BankLab is presented externally by shifting the repository from a purely technical prototype into a more mature engineering laboratory focused on:

* financial systems;
* transactional correctness;
* applied architecture;
* backend/mobile integration;
* onboarding flows;
* future Zero Trust experimentation;
* collaborative open-source development.

The addition of a complete English README also prepares the repository for broader visibility on GitHub and LinkedIn, making it easier for international contributors to understand the project structure, goals, and technical direction.


## 2026/05/15 — docs/update-15

Refactored the Flutter presentation layer naming from `uis` to `ui`, consolidating the shared widget area into `ui/components` and aligning architecture, documentation, agent guides, imports, and tests with the updated structure.

### Architectural and Naming Consolidation

* Renamed the presentation root directory from `mobile/lib/uis` to `mobile/lib/ui`
* Renamed the shared presentation primitives area from `uis/core` to `ui/components`
* Standardized terminology across the project:

  * `uis` → `ui`
  * `shared UI primitives` → `shared UI components`
* Clarified that `ui/components` acts as:

  * a reusable presentation layer
  * a staging area for a future internal mobile widget/feature package
* Reinforced separation rules between:

  * `core`
  * `data`
  * `domain`
  * `ui`

### GitHub Copilot / Agent Instruction Updates

Updated all `.github/instructions/*` guides to reflect the new structure and terminology.

#### Updated instruction files

* `mobile-core.instructions.md`
* `mobile-data.instructions.md`
* `mobile-overview.instructions.md`
* `mobile-pages.instructions.md`
* `mobile-repositories.instructions.md`

#### Renamed instruction guides

* `mobile-uis.instructions.md`
  → `mobile-ui.instructions.md`

* `mobile-uis-core.instructions.md`
  → `mobile-ui-components.instructions.md`

#### Improvements added

* Updated all `applyTo` paths
* Updated dependency rules
* Updated architectural flow references
* Updated shared widget promotion guidance
* Added explicit guidance for future reusable widget packaging
* Clarified component promotion expectations:

  * reusable
  * presentation-only
  * feature-independent

### Mobile Documentation Refactor

Updated mobile documentation to reflect the new folder structure.

#### Updated files

* `mobile/README.md`
* `mobile/AGENT.md`
* `mobile/docs/00-getting_started.md`
* `mobile/docs/01-implemented-features.md`
* `mobile/docs/ARCHITECTURE.md`
* `mobile/Changelog.md`

#### Documentation adjustments

* Replaced all references from:

  * `lib/uis`
  * `uis/core`
  * `shared UI primitives`
* With:

  * `lib/ui`
  * `ui/components`
  * `shared UI components`

#### Added architectural clarification

Expanded documentation around `ui/components`:

* reusable presentation elements
* visual consistency patterns
* future extraction potential into a dedicated package
* distinction between:

  * page-local widgets
  * reusable shared components

### Contribution and Backlog Documentation

Updated project-wide documentation references.

#### Updated files

* `CONTRIBUTING.md`
* `docs/relatorio-mobile-implementado-2026-05-12.md`
* backlog documents under:

  * `docs/backlogs/mobile/done/...`

#### Adjustments

* Updated architecture references to `ui`
* Updated snackbar and component paths
* Updated test references
* Preserved the existing architectural intent while aligning terminology

### Flutter Source Refactor

Performed the physical migration of the Flutter UI layer.

#### Renamed root structure

* `mobile/lib/uis`
  → `mobile/lib/ui`

#### Renamed shared UI structure

* `mobile/lib/uis/core`
  → `mobile/lib/ui/components`

### Shared Component Migration

Migrated all reusable widgets and presentation helpers into `ui/components`.

#### Moved component categories

* base
* buttons
* cards
* input_formatters
* messages
* text
* text_form_field
* themes
* transaction

#### Examples migrated

* `SafeScaffold`
* `AppSnackbar`
* `BigButton`
* `BalanceCard`
* `RecipientCard`
* `BasicTextFormField`
* `CpfInputFormatter`
* `MaterialTheme`
* `TransactionMovement`

### Page and ViewModel Migration

Migrated all page modules into the new `ui/pages` structure.

#### Migrated areas

* auth
* home
* transfer
* shared/details
* splash
* statement

#### Included migrations

* pages
* view models
* local widgets
* presentation models
* exceptions
* route extras

### Routing and Dependency Injection Updates

Updated all imports and references across:

* route registration
* dependency injection
* router extras
* page builders
* typed route payloads

#### Updated files

* `dependencies.dart`
* `main.dart`
* all route registration files
* `extra_codec.dart`

### Test Suite Alignment

Updated Flutter test imports and paths to match the new structure.

#### Updated test files

* `mobile/test/ui/pages/auth/login_feedback_behavior_test.dart`

### Cleanup and Formatting

Performed small consistency improvements during the refactor.

#### Included

* fixed missing trailing newlines
* normalized exports
* aligned import organization
* cleaned obsolete references to `uis`
* updated comments and architecture wording

### Result

This refactor consolidates the Flutter presentation layer into a cleaner and more scalable structure:

* `ui` becomes the single presentation boundary
* `ui/components` becomes the canonical reusable presentation layer
* architectural documentation and agent guidance now align with the actual codebase
* future extraction of reusable mobile UI packages becomes significantly easier
* naming consistency improves readability, onboarding, and long-term maintenance across the entire mobile stack.


## 2026/05/15 — main

Refine collaboration and communication language across project documentation.

### Documentation updates

1. `CONTRIBUTING.md`

   * Reworked the collaboration section to remove the explicit focus on Portuguese-only communication.
   * Clarified that documentation may exist in either Portuguese or English depending on the context.
   * Reinforced the importance of clarity, respect, and collaborative improvement over language preference.

2. `README.md`

   * Adjusted the contribution invitation text to make the project more broadly welcoming.
   * Replaced the “local project” emphasis with a more technical and engineering-oriented positioning.

3. `docs/ROADMAP.md`

   * Updated the roadmap introduction to emphasize technical collaboration instead of local collaboration.
   * Simplified references to “local collaborators” and “local open source” to a more inclusive collaboration model.
   * Replaced the “Collaboration in Portuguese” principle with a clearer documentation-oriented communication principle.
   * Refined backlog and documentation guidance to prioritize clarity and accessibility rather than language-specific requirements.
   * Adjusted wording in contributor-oriented sections to better support mixed-language technical documentation.

### Result

This commit improves the project’s positioning for broader technical collaboration while preserving the original educational and engineering-focused identity of BankLab.


## 2026/05/15 — docs/update-14

Expanded and reorganized the public project documentation structure, transforming the former internal discussion notes into a documented and versioned public decision process for the project.

This update formalizes `docs/` and `docs/backlogs/` as part of the collaborative surface of the repository, exposing architectural discussions, onboarding evolution plans, and implementation debates that were previously treated as internal working material.

The documentation now reflects the reality of the project:

* a constantly evolving codebase;
* iterative architectural refinement;
* ongoing redesigns and trade-off evaluations;
* historical preservation of technical decisions and discarded approaches.

### Documentation Structure Reorganization

1. `README.md`

   * Updated the repository structure description.
   * Repositioned `docs/` as a central area for:

     * roadmap
     * backlogs
     * architectural decisions
     * reports
   * Added direct references to:

     * `docs/README.md`
     * `docs/backlogs/README.md`
     * active onboarding-related backlogs.

2. `docs/README.md`

   * Added a new root-level documentation guide for the `docs/` directory.
   * Documented:

     * overall documentation structure
     * purpose of each directory
     * public communication goals
     * relationship between roadmap, backlogs, and disclosure material.
   * Explicitly established backlogs as part of the project documentation process instead of temporary private notes.

3. `docs/backlogs/README.md`

   * Added a dedicated guide for backlog organization and decision tracking.
   * Defined:

     * active vs resolved backlog separation
     * `done/` directory conventions
     * historical preservation policy
     * contributor guidance for consulting previous discussions before implementation.

4. `docs/ROADMAP.md`

   * Added roadmap guidance reinforcing:

     * public documentation of decision processes;
     * separation between active and completed backlogs;
     * long-term maintainability of collaborative discussions.

5. `CONTRIBUTING.md`

   * Expanded contributor onboarding instructions.
   * Added references to:

     * documentation organization;
     * backlog tracking;
     * architectural discussion history.
   * Established backlog files as first-class project documentation artifacts.
   * Defined preservation rules for completed discussions and implementation history.

6. `docs/.gitignore`

   * Removed backlog exclusion from git tracking.
   * Officially versioned the backlog structure as public repository content.

### Public Backlog Integration

1. `docs/backlogs/api/000 - pre-onboarding.md`

   * Added a large architectural backlog discussing the pre-onboarding restructuring phase.
   * Documented extensive discussions regarding:

     * separation between authentication identity and business identity;
     * restructuring of `users` and `customers`;
     * introduction of:

       * `customer_documents`
       * `customer_addresses`
     * migration strategy for CPF handling;
     * onboarding preparation before full onboarding implementation.
   * Preserved iterative architectural reasoning and evolving decisions instead of reducing the document to a finalized specification.

### Project Philosophy Reinforcement

This update also reinforces an important characteristic of the repository:

* architectural maturity is treated as an incremental process;
* discussions, reversals, refinements, and discarded approaches are intentionally preserved;
* documentation is used not only to describe the final state, but also to expose the reasoning path behind the evolution of the system.

The result is a repository that documents both:

* the implemented software;
* and the engineering decision process that shaped it.


## 2026/05/14 — docs/update-13

Expanded the project documentation structure to position BankLab as a collaborative engineering-oriented financial systems laboratory, improving contributor onboarding, public narrative, roadmap visibility, and repository guidance.

### Main changes

1. `CONTRIBUTING.md`

   * Replaced the previous lightweight GitHub project organization note with a complete contribution guide in Portuguese.
   * Added:

     * project introduction and contribution philosophy;
     * onboarding guidance for new contributors;
     * issue classification standards (`Type`, `Area`, `Priority`);
     * suggested issue template;
     * recommended development workflow;
     * branch naming examples;
     * commit message guidelines;
     * local setup and test commands;
     * API architectural responsibilities and layering guidance;
     * Flutter/mobile contribution conventions;
     * pull request expectations and suggested template;
     * documentation update rules;
     * collaboration and communication principles;
     * explicit out-of-scope contribution boundaries.
   * Established the repository as a learning-oriented but engineering-focused environment for backend, mobile, financial systems, and architecture discussions.

2. `README.md`

   * Completely rewrote and expanded the repository root documentation.
   * Migrated the README from a concise English technical summary to a broader Portuguese presentation focused on:

     * project vision;
     * financial engineering rationale;
     * architectural principles;
     * collaboration narrative;
     * onboarding clarity.
   * Added:

     * detailed project objectives;
     * explicit architectural and financial consistency principles;
     * current scope and future exclusions;
     * stack overview for API, mobile, and infrastructure;
     * repository structure explanation;
     * implemented feature summary;
     * local bootstrap instructions;
     * endpoint overview;
     * development command reference;
     * contribution guidance;
     * suggested first contributions;
     * expanded documentation index;
     * current project status section.
   * Updated endpoint references to reflect newer flows:

     * `/admin/customers/{customer_id}/accounts`
     * `/accounts/internal-transfers`
     * `/accounts/internal-transfers/recipients`
     * transfer receipt endpoint.
   * Improved the positioning of the project as:

     * an engineering laboratory;
     * a transactional consistency study project;
     * a collaboration-oriented fintech sandbox.

3. `docs/ROADMAP.md`

   * Added a new long-form roadmap document.
   * Structured the roadmap around:

     * project identity;
     * engineering motivations;
     * long-term evolution goals;
     * architectural and product directions.
   * Added strategic sections for:

     * public narrative and collaboration;
     * progressive onboarding;
     * Zero Trust Architecture evolution;
     * backend-mediated external integrations;
     * local payment experimentation;
     * web/admin channels;
     * architectural evolution strategy.
   * Added phased planning sections:

     * `Agora`
     * `Próximo`
     * `Futuro`
   * Introduced roadmap items covering:

     * onboarding checkpoints;
     * transactional password modeling;
     * device registration;
     * TOTP and step-up authentication;
     * Pix/DOC/TED mock flows;
     * CI/testing improvements;
     * observability;
     * educational and technical communication goals.
   * Explicitly documented:

     * evolution philosophy;
     * consistency-first mindset;
     * modular monolith rationale;
     * future microservice extraction criteria.

4. `tools/postman/Banklab_API.postman_collection.json`

   * Added a new authenticated request:

     * `GET /customers/me`
   * Configured bearer token authentication using `{{access_token}}`.
   * Expanded the Postman collection coverage for authenticated customer flows.

5. `docs/.gitignore`

   * Added `disclosure` directory to ignored documentation artifacts.

### Documentation direction improvements

* Strengthened the repository positioning as a serious engineering study project rather than a generic CRUD banking demo.
* Improved the public-facing narrative for:

  * potential collaborators;
  * recruiters;
  * technical readers;
  * contributors interested in fintech architecture.
* Consolidated roadmap, contribution, and repository guidance into a more coherent collaboration model.
* Expanded documentation around:

  * transactional consistency;
  * ledger-based modeling;
  * onboarding evolution;
  * Zero Trust Architecture ambitions;
  * engineering rationale and trade-offs.

### Architectural alignment

The new documentation reinforces and contextualizes previously documented concepts such as:

* modular monolith architecture; 
* ledger-centered financial modeling; 
* transactional consistency guarantees; 
* application-layer ownership enforcement; 
* JWT + AppToken authentication model; 
* standardized API response contracts. 

### Result

This commit significantly improves the project’s discoverability, contributor onboarding experience, architectural communication quality, and long-term roadmap clarity, while also preparing the repository for broader technical collaboration and public presentation.


## 2026/05/13 — mobile/login-approved-account-01

Implemented approval-aware authentication feedback for the mobile client, introducing semantic handling for pending-account login attempts and expanding documentation around the current mobile architecture and implemented features.

### Documentation

1. `mobile/README.md`

   * Expanded the authentication flow description to explicitly mention approval-pending login states before account access is granted.

2. `mobile/docs/01-implemented-features.md`

   * Added a new implementation-oriented mobile documentation file.
   * Documented:

     * layered mobile architecture
     * authentication flows
     * account and statement features
     * transfer flows
     * routing organization
     * state and error handling
     * persistence/session handling
     * current testing coverage
   * Included references to repositories, APIs, pages, view models, and storage services currently implemented in the project.

3. `mobile/docs/ARCHITECTURE.md`

   * Standardized `AppSnackbar` as the preferred transient user-feedback mechanism.
   * Added architectural guidance for semantic authentication error mapping.
   * Documented the distinction between invalid credentials and approval-required login states.

4. `mobile/lib/uis/AGENT.md`

   * Reinforced the UI guideline that transient user-facing feedback should use `AppSnackbar.show(...)`.

### Authentication Error Modeling

1. `mobile/lib/core/result/errors/app_error_code.dart`

   * Added the new semantic app error:

     * `AppErrorCode.accountApprovalRequired`

2. `mobile/lib/core/services/client_http/dio/dio_error_mapper.dart`

   * Added explicit backend error mapping for:

     * `ACCOUNT_APPROVAL_REQUIRED`
   * Introduced a stable user-facing approval-pending message.
   * Preserved generic HTTP error behavior for unrelated failures.
   * Improved `details` propagation by preserving the entire backend error payload when `details` is absent.

### Login UI Behavior

1. `mobile/lib/uis/pages/auth/login/login_page.dart`

   * Added semantic handling for `accountApprovalRequired`.
   * Implemented approval-pending snackbar feedback for the full login flow.
   * Preserved existing invalid-credential behavior.

2. `mobile/lib/uis/pages/auth/short_login/short_login_page.dart`

   * Added approval-pending handling for the remembered-account login flow.
   * Preserved short-login identity state after failure.
   * Kept generic failure handling unchanged for non-approval scenarios.

### Test Coverage

1. `mobile/test/core/services/client_http/dio/dio_error_mapper_test.dart`

   * Added tests validating:

     * backend approval-required mapping
     * unchanged invalid-credentials behavior
     * fallback forbidden handling

2. `mobile/test/data/repositories/auth/auth_repository_impl_test.dart`

   * Added repository-level tests validating:

     * approval-required login failures do not persist tokens
     * profile loading is skipped on approval failures
     * remembered-login cache is not updated on rejected login attempts
     * successful login still persists tokens and updates remembered identity correctly

3. `mobile/test/uis/pages/auth/login_feedback_behavior_test.dart`

   * Added widget tests covering:

     * full login approval-pending feedback
     * short-login approval-pending feedback
     * preservation of generic invalid credential behavior
     * preservation of generic server failure behavior

This commit improves the authentication UX by distinguishing operational account state from credential failure, aligning the mobile client with the backend semantic error model while strengthening architectural documentation and automated test coverage.


## 2026/05/13 — api/login-approved-account-01

Implemented login eligibility enforcement for customer users, requiring completed approval and account provisioning before session creation.

### Authentication Flow

* Added account provisioning validation to `LoginUserUseCase`
* Injected `AccountRepository` as an `AccountProvisioningChecker` dependency during application wiring
* Introduced `validateLoginEligibility()` to centralize operational login checks
* Restricted customer login based on:

  * `pending` status → `ACCOUNT_APPROVAL_REQUIRED`
  * missing `customer_id` → `ACCOUNT_APPROVAL_REQUIRED`
  * no provisioned account → `ACCOUNT_APPROVAL_REQUIRED`
  * `blocked` or invalid lifecycle states → `FORBIDDEN`
* Preserved admin login behavior without requiring account provisioning
* Prevented token generation and session persistence when approval requirements are not satisfied

### Domain and Error Handling

* Added new domain error:

  * `ErrAccountApprovalRequired`
* Added shared stable error code:

  * `ACCOUNT_APPROVAL_REQUIRED`
* Registered HTTP mapping:

  * HTTP 403 Forbidden
  * message: `Account approval required`
* Extended centralized error registry with the new authorization/business-state error

### Application Layer

* Added `AccountProvisioningChecker` interface abstraction to the auth application layer
* Improved constructor dependency documentation and responsibility description
* Added defensive validation for missing provisioning checker configuration
* Wrapped infrastructure failures during account provisioning verification with contextual errors

### Tests

* Expanded `login_user_test.go` with approval/account provisioning coverage:

  * successful active customer login
  * pending customer rejection
  * active customer without account rejection
  * active customer without customer_id rejection
  * admin login without account success
  * provisioning repository error wrapping
* Added provisioning checker mock implementation
* Added assertions ensuring:

  * token generation is skipped on rejection
  * sessions are not persisted on rejection
  * provisioning checks are invoked only when required
* Updated existing login tests to support new constructor signature
* Added delivery-layer handler test validating:

  * HTTP 403 response
  * `ACCOUNT_APPROVAL_REQUIRED` payload contract
* Updated auth integration test bootstrap wiring with account repository injection

### Documentation

* Updated REST API documentation:

  * login eligibility behavior
  * new `ACCOUNT_APPROVAL_REQUIRED` error
  * explicit 403 response example
  * onboarding/account provisioning dependency explanation
* Updated implementation documentation describing:

  * login blocking for unapproved customers
* Updated architectural and conceptual docs:

  * authentication lifecycle
  * onboarding boundary
  * operational eligibility semantics
  * approval/account provisioning relationship
* Clarified that customer login now depends on:

  * admin approval
  * atomic account provisioning via `/admin/users/{id}/approve`

### Result

The authentication flow now enforces a stricter operational lifecycle boundary, ensuring that customer users can only establish authenticated sessions after administrative approval and successful account provisioning. This aligns login behavior with the system’s lifecycle and operational consistency model.


## 2026/05/13 — api-mobile/routes-02-statement-04

Refactored account provisioning boundaries, removed customer self-service account creation routes, and aligned the API surface with the intended operational model for admin-controlled account provisioning and future terminal channels.

### API and Route Architecture

1. `api/cmd/api/main.go`

   * Extracted API route registration into `newAPIRouter(...)` for clearer composition and isolated router testing.
   * Added admin route:

     * `POST /admin/customers/{customer_id}/accounts`
   * Removed customer-facing:

     * `POST /accounts`
   * Kept terminal cash routes commented and intentionally unregistered:

     * `/terminal/accounts/{id}/deposit`
     * `/terminal/accounts/{id}/withdraw`
   * Simplified `main()` wiring responsibilities.

2. `api/cmd/api/routes_test.go`

   * Added router-level tests validating that:

     * `POST /accounts` is no longer registered.
     * terminal deposit route is not exposed.
     * terminal withdraw route is not exposed.
   * Ensured operational-only endpoints remain outside the public REST surface.

### Account Provisioning Refactor

3. `api/internal/account/bankaccount/application/create_account.go`

   * Refactored account creation flow to support explicit admin provisioning.
   * Added `CustomerID` to `CreateAccountInput`.
   * Changed authorization model:

     * now requires authenticated admin role.
   * Removed dependency on authenticated customer ownership context.
   * Updated documentation/comments to reflect provisioning semantics.

4. `api/internal/account/bankaccount/delivery/account_handler.go`

   * Replaced `CreateAccount` with:

     * `CreateAccountForCustomer`
   * Added:

     * admin role validation
     * path-based `customer_id` parsing
     * invalid UUID handling
   * Updated operational logging to:

     * `event=admin_create_account`

5. `api/internal/account/bankaccount/delivery/auth_test.go`

   * Added `testAdminRequest(...)` helper for admin-authenticated request contexts.

### Account Provisioning Tests

6. `api/internal/account/bankaccount/application/create_account_test.go`

   * Migrated all account creation tests to admin provisioning semantics.
   * Added explicit `CustomerID` usage throughout tests.
   * Updated use case expectations for admin execution context.

7. `api/internal/account/bankaccount/delivery/account_handler_test.go`

   * Reworked handler tests for:

     * admin-only provisioning
     * invalid customer ID handling
     * unknown field rejection
     * empty-body success flow
     * non-admin rejection
   * Removed obsolete customer self-service account creation tests.

### Documentation Updates

8. `api/README.md`

   * Updated exposed API surface.
   * Clarified that:

     * account creation is an admin provisioning capability.
     * terminal deposit/withdraw channels are intentionally disabled.

9. `api/docs/06-implementation.md`

   * Updated runtime route registration section.
   * Documented disabled terminal operations.

10. `api/docs/07-api-rest.md`

    * Added:

      * `3.6 Create Customer Account (Admin Only)`
    * Removed customer-facing `POST /accounts` documentation.
    * Updated error scenarios and route references.
    * Clarified provisioning responsibilities and onboarding behavior.

11. `api/docs/ARCHITECTURE.md`

    * Updated registered route inventory.
    * Documented inactive terminal routes.

12. `api/docs/presentation/presentation-api-architecture.md`

    * Refined presentation narrative around:

      * admin provisioning
      * disabled terminal operations
      * active REST surface
      * future terminal channel direction

13. `api/docs/visao_geral/chapters/*`

    * Updated architectural overview chapters to reflect:

      * removal of `POST /accounts`
      * admin provisioning semantics
      * disabled terminal routes
      * ownership/context derivation adjustments
      * revised operational flows

### Mobile Statement Improvements

14. `mobile/lib/core/extensions/datetime_extension.dart`

    * Added:

      * `formatMonthLabel`
      * `formatDayLabel`
      * `formatHour`

* Introduced `DateParser` utility class.
* Added `parseOrNow(...)`.

15. `mobile/lib/data/services/apis/account/dtos/statement_response_dto.dart`

    * Migrated statement date fields from `String` to `DateTime`.
    * Improved parsing consistency for:

      * statement items
      * pagination cursors

16. `mobile/lib/data/services/auth/api/dtos/customer_me_response_dto.dart`

    * Replaced legacy parsing with `DateParser.parseOrNow(...)`.

17. `mobile/lib/domain/common/auth/models/user_profile.dart`

    * Replaced legacy parsing helpers with `DateParser`.

18. `mobile/lib/data/repositories/account/account_repository_impl.dart`

    * Reset statement cache before loading a new statement request.

### Statement UI Refactor

19. `mobile/lib/uis/pages/statement/statement_page.dart`

    * Refactored statement grouping to use typed `DateTime` keys.
    * Removed legacy string/date parsing helpers.
    * Simplified sorting and grouping logic.
    * Improved formatting reuse through extension methods.
    * Extracted reusable widgets for:

      * load error state
      * empty transaction state

20. `mobile/lib/uis/pages/statement/widgets/load_statement_error.dart`

    * Added reusable error-state widget with retry support.

21. `mobile/lib/uis/pages/statement/widgets/no_transactions_card.dart`

    * Added reusable empty-state widget.

22. `mobile/lib/uis/pages/statement/widgets/statement_item_card.dart`

    * Simplified layout structure.
    * Reduced visual nesting complexity.
    * Improved transaction amount visibility.
    * Reworked description/hour alignment.
    * Removed obsolete commented balance display code.

### Mobile Tests

23. `mobile/test/core/extensions/datetime_extension_test.dart`

    * Added coverage for:

      * `DateParser.parseOrNull`
      * `DateParser.parseOrNow`
      * locale-aware formatting
      * month/day/hour formatting extensions

### Tooling

24. `tools/postman/Environment.postman_environment.json`

    * Updated local `base_url` environment IP.

This commit consolidates the transition from customer self-service account creation to explicit administrative provisioning while preserving the future architectural direction for terminal-based cash operations. It also significantly improves statement rendering consistency on mobile through typed date handling, reusable UI components, and cleaner formatting utilities.


## 2026/05/13 — api/routes-01

Refined the API route surface to better distinguish operational ledger endpoints from customer-facing product flows, while aligning runtime wiring, documentation, and tests around the new terminal-oriented route structure.

### Main Changes

1. Route surface restructuring and operational boundary clarification

   * Repositioned deposit and withdraw endpoints from:

     * `/accounts/{id}/deposit`
     * `/accounts/{id}/withdraw`
   * To:

     * `/terminal/accounts/{id}/deposit`
     * `/terminal/accounts/{id}/withdraw`
   * Explicitly documented these operations as terminal/operational ledger flows rather than customer-facing product capabilities.
   * Added architectural guidance clarifying that:

     * account creation is still a provisioning-oriented operation
     * deposit/withdraw should eventually move behind protected operational/admin surfaces
     * future onboarding and cash-in/cash-out flows should replace direct ledger mutation endpoints

2. Runtime API wiring hardening

   * Removed active registration of deposit and withdraw routes from `cmd/api/main.go`.
   * Left the route definitions commented with explanatory notes indicating intentional disabling until a real terminal channel exists.
   * Preserved transfer, balance, statement, and recipient routes unchanged.

3. Documentation consistency updates

   * Updated REST documentation references, route indexes, endpoint sections, examples, and error scenario sections to use `/terminal/accounts/...`.
   * Added operational notes to:

     * account creation
     * deposit
     * withdraw
   * Clarified that deposit/withdraw routes are intentionally disabled in the current runtime wiring.
   * Updated architecture and implementation documents to reflect the revised route semantics and operational positioning.
   * Updated presentation and “visão geral” chapters to maintain consistency across:

     * flow descriptions
     * examples
     * diagrams
     * route listings
     * testing references

4. Test suite alignment

   * Updated integration and handler tests to use the new terminal route namespace.
   * Adjusted:

     * deposit integration tests
     * authorization integration tests
     * handler tests for deposit and withdraw
   * Preserved existing authorization and ownership validation behavior while aligning with the renamed routes.

5. Added implementation snapshot reports

   * Added `docs/relatorio-api-implementada-2026-05-12.md`

     * comprehensive implementation inventory of the API
     * architecture, auth model, persistence, transactional guarantees, routes, flows, invariants, migrations, and operational notes
   * Added `docs/relatorio-mobile-implementado-2026-05-12.md`

     * detailed overview of the Flutter mobile implementation
     * architecture, navigation, DI, auth/session handling, repositories, use cases, and implemented user journeys

### Architectural Impact

This change improves conceptual separation between:

* customer-facing banking flows
* operational ledger mutation endpoints
* future onboarding/admin provisioning surfaces

The current API now communicates more explicitly that:

* transfer is a real customer operation
* deposit/withdraw are infrastructure or terminal-oriented operations
* account provisioning is transitional and expected to evolve into administrative onboarding flows

This reduces ambiguity in the public API contract and strengthens the long-term architectural direction of the project.


## 2026/05/12 — mobile/statement-02

Implemented the first complete statement flow for the mobile application, including backend support for transaction descriptions, statement navigation, grouped UI rendering, cached state handling, and receipt behavior improvements.

### API

1. Statement response enrichment

   * Added optional `description` support to statement items returned by the API.
   * Updated statement application model, delivery DTOs, handlers, and infrastructure repository mappings.
   * Extended SQL queries to fetch transaction descriptions from the `transactions` table.
   * Added `omitempty` behavior for empty descriptions in JSON responses.
   * Updated REST API documentation with statement description examples and optional field notes.

2. Statement repository test coverage

   * Added `repository_test.go` for statement infrastructure.
   * Validated:

     * SQL query selection includes `description`
     * transaction description mapping
     * timestamp propagation
     * nullable field handling

3. Transfer receipt behavior correction

   * Adjusted transfer receipt operation type resolution for destination customers.
   * Destination users now correctly see `transfer_in` instead of the original transfer direction.
   * Added dedicated unit test validation for the corrected behavior.

### Mobile — Routing & Navigation

1. Statement route integration

   * Added `BaseRoutes.statement`.
   * Registered `StatementPage` and `StatementViewmodel` in GoRouter and dependency injection.

2. Home page navigation

   * Replaced the previous “coming soon” placeholder action with real statement navigation.
   * Added `_navToStatement()` navigation handler.

3. Details page navigation improvements

   * Replaced forced navigation to home with `context.pop()`.
   * Updated button label from `Fechar` to `Voltar`.

### Mobile — Statement Feature

1. Statement page implementation

   * Added full `StatementPage`.
   * Implemented:

     * loading state
     * retry flow
     * refresh indicator
     * empty state
     * grouped rendering by month/day
     * operation detail navigation
     * snackbar-based error handling

2. Statement grouping and formatting

   * Added:

     * `MonthHeader`
     * `DayHeader`
     * `StatementItemCard`
   * Implemented:

     * month grouping
     * daily grouping
     * daily consolidated balance display
     * localized date/hour formatting

3. Statement interaction flow

   * Transactions with valid references now navigate to details screen.
   * Transactions without references show informational feedback.

4. Statement caching support

   * Added `lastStatement` support to `AccountRepository`.
   * Added internal `_statementCache`.
   * Cache is invalidated when switching accounts.
   * Cached statement is reused during refresh/error scenarios.

5. Statement DTO improvements

   * Added `description` field to `StatementItemDto`.
   * Added safe parsing fallback for nullable descriptions.

### Mobile — Shared Transaction Presentation

1. Transaction movement abstraction

   * Added `TransactionMovement`.
   * Centralized:

     * operation labels
     * debit/credit semantics
     * amount sign formatting

2. Details page integration

   * Receipt operation labels now use `TransactionMovement`.
   * Improved semantic consistency between statement and receipt views.

### Mobile — UI & Refactoring

1. Statement card component

   * Added reusable visual representation for statement entries.
   * Included:

     * debit/credit coloring
     * optional descriptions
     * detail navigation indicators
     * formatted timestamps

2. Splash screen adjustment

   * Wrapped splash logo in a card container with padding for improved visual contrast.

3. Import normalization

   * Standardized several imports to use absolute project paths.

4. Exception extraction

   * Moved `ReceiptImageException` into its own dedicated exception file.

5. BigButton flexibility improvements

   * `leftIcon` and `rightIcon` now accept `Widget` instead of `IconData`.
   * Increased flexibility for future UI composition.

### Tests

1. Statement handler tests

   * Added assertions validating description propagation through HTTP responses.

2. Transfer receipt tests

   * Added validation for destination-side `transfer_in` operation type behavior.

3. Fake repository updates

   * Extended fake repositories with `lastStatement` support for compatibility with the new repository contract.


## 2026/05/12 — main

Refactor authentication and transfer action buttons to support fully customizable icon widgets and improve UI consistency across the mobile application.

### Button Component Improvements

Updated the reusable button components to support generic widget injection instead of limiting icons to `IconData`.

#### `mobile/lib/uis/core/buttons/big_button.dart`

* Replaced `IconData?` with `Widget?` for:

  * `leftIcon`
  * `rightIcon`
* Removed internal `Icon(...)` wrapping
* Allowed direct rendering of arbitrary widgets inside buttons
* Preserved spacing and alignment behavior

This change enables:

* loading indicators
* animated widgets
* custom icon compositions
* future extensibility without modifying the button component again

#### `mobile/lib/uis/core/buttons/big_text_button.dart`

* Added support for:

  * `leftIcon`
  * `rightIcon`
* Refactored button content into a centered `Row`
* Added conditional rendering and spacing logic for optional widgets

This aligns `BigTextButton` behavior with `BigButton`, improving component consistency across the UI layer.

### Authentication UI Refactor

Refactored authentication pages to use the shared `BigButton` component instead of raw `FilledButton`.

#### `mobile/lib/uis/pages/auth/login/login_page.dart`

* Replaced `FilledButton` with `BigButton`
* Added:

  * dynamic label (`Entrando...` / `Entrar`)
  * loading spinner as `rightIcon`
  * login icon when idle

Improves:

* visual consistency
* loading feedback
* button standardization

#### `mobile/lib/uis/pages/auth/register/register_page.dart`

* Replaced `FilledButton` with `BigButton`
* Added:

  * dynamic loading label (`Cadastrando...`)
  * loading spinner support
  * registration icon

Also normalized imports during refactor.

#### `mobile/lib/uis/pages/auth/short_login/short_login_page.dart`

* Applied the same refactor pattern used in login page
* Added loading indicator support through widget-based icons
* Unified interaction behavior with the standard login flow

### Transfer Flow UI Adjustments

Updated all transfer-related screens to comply with the new widget-based icon API.

#### Updated Pages

* `transfer_confirmation_page.dart`
* `transfer_payment_page.dart`
* `transfer_recipient_page.dart`
* `transfer_status_page.dart`
* `details_page.dart`

Changes:

* replaced raw `IconData` values with explicit `Icon(...)` widgets
* standardized icon sizing (`size: 24`)
* preserved existing button semantics and flow behavior

### UX and Architectural Impact

This refactor improves the UI component architecture by:

* reducing coupling between button components and icon implementation details
* enabling richer interactive widgets inside buttons
* improving visual consistency across authentication and transfer flows
* centralizing button styling and behavior
* simplifying future evolution of loading and action states

The resulting API for button components is now significantly more flexible while remaining backward-compatible with standard Flutter widgets.


## 2026/05/12 — mobile/short-login-02

Implemented the first version of the short login flow with cached user identity recovery, splash bootstrap routing, and authentication module restructuring.

### Routing and navigation refactor

1. `mobile/lib/core/routing/routes.dart`

   * Added `AuthRoutes.shortLogin`
   * Replaced `HomeRoutes` with `BaseRoutes`
   * Added `BaseRoutes.splash`

2. `mobile/lib/core/routing/router.dart`

   * Changed application startup route from `/login` to `/splash`
   * Replaced `homeRoutes()` with `baseRoutes()`

3. `mobile/lib/core/routing/routes/base_routes.dart`

   * Added centralized base route configuration
   * Registered:

     * splash route
     * home route

4. `mobile/lib/core/routing/routes/auth_routes.dart`

   * Added short login route support
   * Added `LastLoginIdentity` route payload validation
   * Added fallback to regular login when route payload is invalid
   * Integrated:

     * `ShortLoginPage`
     * `ShortLoginViewModel`

5. `mobile/lib/core/routing/routes/home_routes.dart`

   * Removed obsolete home route file after migration to `base_routes.dart`

6. `mobile/lib/core/routing/routes/tranfer_routes.dart`

   * Renamed file to `transfer_routes.dart`
   * Fixed routing naming typo

### Splash bootstrap flow

7. `mobile/lib/uis/pages/splash/viewmodel/splash_viewmodel.dart`

   * Added splash initialization flow
   * Introduced startup identity recovery command
   * Integrated cached login identity retrieval through `AuthRepository`

### Short login infrastructure

8. `mobile/lib/data/services/auth/cache/`

   * Added last login cache service infrastructure
   * Introduced secure storage based login identity persistence

9. `mobile/lib/data/repositories/auth/auth_repository_impl.dart`

   * Injected `LastLoginCacheService`
   * Added integration between authentication lifecycle and cached identity persistence

10. `mobile/lib/data/repositories.dart`

    * Updated repository dependency graph
    * Registered new authentication cache dependencies

11. `mobile/lib/data/services.dart`

    * Registered:

      * `LastLoginCacheService`
      * `LastLoginCacheServiceImpl`
    * Added secure storage dependency wiring

### Authentication module restructuring

12. Authentication DTO imports and service organization

    * Migrated auth service imports from:

      * `data/services/apis/auth`
    * To:

      * `data/services/auth/api`
    * Improved module separation between:

      * authentication APIs
      * authentication cache services

13. Updated multiple files with normalized import paths

    * Reduced inconsistent relative imports
    * Improved package organization consistency

### ViewModel registration updates

14. `mobile/lib/uis/viewmodels.dart`

    * Registered:

      * `SplashViewmodel`
      * `ShortLoginViewModel`

### UI and asset updates

15. `mobile/pubspec.yaml`

    * Registered `assets/images/brand.png`

16. `mobile/assets/images/banklab.svg`

    * Added new SVG branding asset for BankLab visual identity

### Environment updates

17. `tools/postman/Environment.postman_environment.json`

    * Updated local API host IP for development environment synchronization

### General impact

This commit introduces the foundation for a faster authentication experience through cached login identity recovery while also restructuring the authentication module for future expansion.

The new splash-driven startup flow centralizes bootstrapping responsibilities and prepares the mobile application for:

* persistent login experiences
* device-aware authentication flows
* future transactional authentication layers
* Zero Trust oriented identity evolution
* biometric and device registration extensions in future iterations


## 2026/05/11 - mobile/short-login-01

Refactor authentication profile loading flow and introduce support structures for short login persistence on the mobile application.

### Authentication and profile integration

This commit restructures the mobile authentication profile retrieval flow to align the client with the current API contract by separating identity and customer profile responsibilities between `/auth/me` and `/customers/me`.

#### Updated `AuthApi.getProfile()`

* Replaced the old `profile/me` request with:

  * `GET /customers/me`
  * `GET /auth/me`
* Added explicit parsing and validation for both response envelopes.
* Merged both endpoint responses into a unified `UserProfile` model through a dedicated factory constructor.
* Improved failure propagation and parsing isolation for each endpoint.
* Preserved compatibility with the API response envelope strategy already established in the backend documentation.

### New DTOs for authentication/profile separation

Added dedicated transport DTOs:

#### `AuthMeResponseDto`

Represents authenticated identity/session information:

* user id
* role
* email
* customer id

#### `CustomerMeResponseDto`

Represents customer profile information:

* customer id
* name
* cpf
* email
* createdAt

This separation better reflects the backend architecture and authorization model currently implemented in the API.

### User profile model improvements

Extended `UserProfile` with:

```dart
UserProfile.fromMe(...)
```

This constructor consolidates:

* authenticated user context
* customer profile data

into a single application-facing representation.

The new approach reduces coupling between API transport contracts and UI/domain consumption.

### Short login persistence support

Added new secure/local storage keys:

```dart
lastLoginName
lastLoginIdentifier
```

These keys prepare the application for:

* cached login identity
* short login flows
* reduced friction during authentication
* future biometric or transactional authentication extensions

This is an important step toward a banking-style login experience where:

* the user identity remains cached locally
* only the credential confirmation step is required on subsequent accesses

### Import and path normalization

Adjusted several imports to:

* use absolute project-based imports
* reduce relative path traversal
* improve consistency across the mobile module

This includes updates in:

* `details_page.dart`
* authentication DTO/model integration files

### Architectural alignment

The changes reinforce the separation already defined in the backend architecture:

* identity/session context from `/auth/me`
* customer/business profile from `/customers/me`

This keeps the mobile client aligned with:

* JWT ownership rules
* customer-bound authorization
* modular backend boundaries
* future Zero Trust evolution paths

The result is a cleaner authentication/profile pipeline with improved transport separation, reduced coupling, and a stronger foundation for future secure login and contextual authentication features.


## 2026/05/11 - mobile/internal-transfer-11

Refactor the internal transfer flow and introduce a complete transfer receipt/details experience in the Flutter mobile client.

### Documentation

1. `api/docs/00-getting_started.md`

   * Added the bootstrap flow for the first administrator user.
   * Documented the manual PostgreSQL promotion process for the initial admin account.
   * Clarified the approval lifecycle for newly registered users.
   * Added operational guidance for calling `POST /admin/users/{id}/approve`.
   * Improved onboarding clarity for local environment setup and first system access.

### Routing and Navigation

1. `mobile/lib/core/routing/router.dart`

   * Added shared route registration.

2. `mobile/lib/core/routing/routes.dart`

   * Introduced `SharedRoutes` enum with the transfer details route.

3. `mobile/lib/core/routing/routes/shared_routes.dart`

   * Added shared navigation route for transfer details/receipt visualization.
   * Injected `DetailsViewmodel` through dependency injection.
   * Enabled route argument passing through `state.extra`.

4. `mobile/lib/core/routing/routes/tranfer_routes.dart`

   * Renamed transfer pages to explicit `Transfer*` naming convention.
   * Updated route builders to use the renamed widgets.
   * Simplified imports and normalized route organization.
   * Added integration with the new transfer status and receipt navigation flow.

### Domain and Dependency Injection

1. `mobile/lib/domain/usecases/details/details_usecase.dart`

   * Added `DetailsUsecase` for transfer receipt retrieval.
   * Integrated account and transaction repositories.
   * Exposed selected account access for shared details flows.

2. `mobile/lib/domain/usecases/usecases.dart`

   * Registered `DetailsUsecase` in dependency injection.

3. `mobile/lib/uis/viewmodels.dart`

   * Registered `DetailsViewmodel` in dependency injection.

### Shared Transfer Receipt Experience

1. `mobile/lib/uis/pages/shared/details/details_page.dart`

   * Added a dedicated transfer receipt/details page.
   * Implemented asynchronous receipt loading flow.
   * Added receipt rendering with transfer metadata presentation.
   * Added transaction reference copy support.
   * Added receipt sharing support using generated PNG images.
   * Implemented `RepaintBoundary` image capture flow.
   * Added temporary image generation and persistence.
   * Added integration with native share APIs using `share_plus`.
   * Added detailed snackbar feedback handling.
   * Added loading, retry, and failure states.
   * Added responsive receipt card rendering.
   * Added transaction status visualization.
   * Added share state protection to avoid duplicated actions.
   * Added defensive rendering and capture validation logic.

2. `mobile/lib/uis/pages/shared/details/viewmodel/details_viewmodel.dart`

   * Added command-based receipt retrieval state management.
   * Connected details flow to `DetailsUsecase`.

3. `mobile/lib/uis/pages/shared/details/widgets/detail_line.dart`

   * Added reusable detail line widget for receipt rendering.

### Transfer Flow Improvements

1. `mobile/lib/uis/pages/home/transfer/transfer_status_page.dart`

   * Replaced modal receipt visualization with dedicated details navigation.
   * Added new action layout using `BigTextButton`.
   * Simplified successful transfer presentation.
   * Improved CTA hierarchy and spacing.
   * Added receipt navigation through shared routes.
   * Updated success message wording.

2. `mobile/lib/uis/pages/home/transfer/transfer_confirmation_page.dart`

   * Renamed confirmation page class for consistency.
   * Migrated snackbar usage to `AppSnackbar`.
   * Improved transfer failure feedback handling.

3. `mobile/lib/uis/pages/home/transfer/transfer_payment_page.dart`

   * Renamed payment page class for consistency.
   * Replaced deprecated snackbar helper usage.
   * Improved invalid amount feedback handling.

4. `mobile/lib/uis/pages/home/transfer/transfer_recipient_page.dart`

   * Renamed recipient page class for consistency.
   * Added forward navigation icon to primary CTA.

5. `mobile/lib/uis/pages/home/transfer/transfer_page.dart`

   * Removed obsolete transfer page implementation.

### UI Components and Shared Widgets

1. `mobile/lib/uis/core/buttons/big_button.dart`

   * Made `onPressed` nullable.
   * Added default enabled state.
   * Improved disabled button compatibility.

2. `mobile/lib/uis/core/buttons/big_text_button.dart`

   * Added reusable large text button component.

3. `mobile/lib/uis/core/text/card_text_row.dart`

   * Improved multiline layout handling.
   * Added expanded text rendering.
   * Improved alignment for long values.

### Feedback and Messaging

1. `mobile/lib/uis/core/messages/app_snackbar.dart`

   * Moved snackbar utilities into a dedicated `messages` namespace.

2. `mobile/lib/uis/core/messages/scaffold_snackbar_message.dart`

   * Removed legacy snackbar helper implementation.

3. Updated snackbar imports across:

   * `login_page.dart`
   * `home_page.dart`
   * transfer pages

### Dependencies and Platform Support

1. `mobile/pubspec.yaml`

   * Added `share_plus`.
   * Added `path_provider`.

2. `mobile/pubspec.lock`

   * Updated dependency graph for sharing and filesystem support.
   * Added transitive URL launcher and sharing platform packages.

3. `mobile/ios/Podfile.lock`

   * Added iOS CocoaPods integration for `share_plus`.

This commit consolidates the internal transfer UX around a dedicated receipt/details flow, improves route consistency, modernizes shared feedback components, and establishes the foundation for richer financial operation visualization inside the mobile client.


## 2026/05/11 - mobile/internal-transfer-10

Implemented the complete internal transfer confirmation and status flow in the Flutter mobile application, introducing presentation-layer models, typed route serialization, money parsing/formatting helpers, and improved transfer UX/navigation behavior.

### Transfer flow implementation

* Added the new transfer confirmation step before executing the transfer request.
* Implemented success and failure transfer status screens with transaction feedback and receipt preview support.
* Introduced dedicated success and failure routes:

  * `TransferRoutes.statusSuccess`
  * `TransferRoutes.statusFailure`
* Added typed route serialization support for transfer confirmation payloads through `ExtraCodec`.

### Transfer presentation model architecture

* Introduced `TransferConfirmationData` as a feature-local presentation model under:

  * `uis/pages/home/transfer/models`
* Added serialization/deserialization support:

  * `toMap`
  * `fromMap`
* Encapsulated immutable transfer confirmation state independently from DTOs and domain use case inputs.
* Improved architectural separation between:

  * API DTOs
  * domain use case contracts
  * UI workflow state

### Payment page improvements

* Reworked `PaymentPage` to support:

  * amount input
  * transfer description input
  * recipient review
  * navigation to confirmation flow
* Added money validation before confirmation.
* Added formatted currency input using a custom formatter.
* Added invalid transfer feedback through snackbar messages.
* Improved keyboard dismissal behavior using `GestureDetector`.

### Confirmation page implementation

* Added `ConfirmationPage` with:

  * transfer summary visualization
  * balance visualization
  * recipient/account review
  * transfer confirmation execution
* Integrated transfer execution with `TransferUsecase`.
* Added navigation handling for:

  * success state
  * failure state
  * unexpected execution state

### Transfer status experience

* Added `TransferStatusPage` supporting:

  * success visualization
  * failure visualization
  * transaction reference display
  * receipt modal preview
* Added contextual UI feedback using:

  * icons
  * themed status colors
  * receipt action buttons

### Money formatting and parsing utilities

* Added `MoneyInputFormatter` for currency-formatted numeric input.
* Added `String.parseToMoney()` extension helper using `money2`.
* Added reusable numeric sanitization and currency conversion support.
* Centralized invalid amount error handling through `Result<AppError>`.

### Shared UI improvements

* Enhanced `BigButton` with optional:

  * `leftIcon`
  * `rightIcon`
* Added reusable snackbar helper:

  * `showScaffoldSnackBarMessage`
* Moved `RecipientCard` into shared UI components:

  * `uis/core/cards`

### Recipient flow refactor

* Removed recipient selection state ownership from `TransferViewmodel`.
* Moved recipient selection state management into `RecipientPage`.
* Improved recipient/account auto-fill behavior:

  * CPF lookup now updates branch/account
  * account lookup now updates CPF
* Fixed branch method naming typo:

  * `_onBranchChenged` → `_onBranchChanged`
* Improved recipient error handling and invalid-state reset behavior.

### ViewModel cleanup

* Simplified `TransferViewmodel` responsibilities by removing UI-local selection state.
* Added lightweight `dispose()` method placeholder for future lifecycle expansion.
* Preserved transfer idempotency generation behavior through UUID generation.

### Documentation and architecture guidance

* Expanded `AGENT.md` documentation under:

  * `uis/AGENT.md`
  * `uis/pages/AGENT.md`
* Documented the concept of feature-local presentation models.
* Clarified architectural boundaries between:

  * DTOs
  * use case inputs
  * UI flow models
* Added guidance for route extras and immutable presentation snapshots.

This commit significantly advances the mobile banking transfer workflow, improving UX continuity, route safety, UI architecture consistency, and transfer execution feedback while reinforcing the separation between domain contracts, transport DTOs, and presentation-layer state.


## 2026/05/10 - mobile/internal-transfer-09

Refactor the internal transfer flow structure by introducing route serialization support, reusable financial UI components, and balance stream propagation across the transfer workflow.

### Routing and transfer flow evolution

This commit expands the internal transfer navigation pipeline, preparing the application for a complete multi-step transfer experience.

#### Routing improvements

1. Added new route groups and navigation stages:

   * `GeneralRoutes`
   * `TransferRoutes.payment`
   * `TransferRoutes.confirmation`
   * `TransferRoutes.status`
   * `GeneralRoutes.receipt`

2. Added the `PaymentPage` route to the transfer flow:

   * integrated with `go_router`
   * receives typed `RecipientInfoDto` through `state.extra`

3. Extended `ExtraCodec` to support typed serialization/deserialization of:

   * `RecipientInfoDto`

This change enables safe route persistence and structured navigation state handling during app lifecycle restoration and deep navigation scenarios.

#### Recipient transfer flow

4. Updated `RecipientPage`:

   * extracted button implementation into reusable `BigButton`
   * replaced direct pop navigation with forward navigation to payment flow
   * added `context.pushNamed(...)`
   * passes selected recipient through route extra serialization
   * renamed internal validation helper to `_isValidBranchAndAccount`

### Transfer balance propagation

One of the main focuses of this commit was restructuring how balance information propagates through transfer-related screens.

#### Transfer use case and viewmodel

5. Added balance stream exposure to:

   * `TransferUsecase`
   * `TransferViewmodel`

6. Introduced:

   * `Stream<BalanceResponseDto> balance`

This allows transfer screens to consume live balance updates directly from the account repository layer.

#### Repository balance behavior

7. Refactored `AccountRepository`:

   * `selectAccount` now returns `AsyncResult<Unit>`
   * removed `clearSelectedAccount`

8. Updated `AccountRepositoryImpl`:

   * `balance()` now emits cached balance immediately before streaming updates
   * `selectAccount(...)` now:

     * validates cache existence
     * validates account existence
     * loads account balance automatically
     * returns explicit success/failure results

This significantly improves transfer screen initialization consistency and reduces duplicated orchestration responsibilities in the UI layer.

### Reusable UI components

This commit also introduces reusable financial UI building blocks.

#### Added reusable widgets

9. Added `BigButton`

   * centralized primary large CTA button styling
   * standardized enabled/disabled states

10. Added `BalanceCard`

    * reusable balance visualization component
    * supports:

      * live balance streams
      * account visibility toggling
      * account selection actions
      * selected account display

11. Added `AccountCard`

    * reusable account information visualization

12. Added `CardTextRow`

    * reusable structured row presentation for cards

### Home and transfer UI integration

13. Refactored `HomePage`

* replaced legacy `BalanceTile`
* integrated new reusable `BalanceCard`

14. Added initial `PaymentPage`

* displays:

  * current balance
  * selected account
* integrates transfer flow navigation structure
* prepared for future transfer confirmation implementation

### DTO and serialization support

15. Added `RecipientInfoDto.toMap()`

* required for route serialization through `ExtraCodec`

This creates a proper serialization boundary for typed navigation payloads.

### Tests and repository contract updates

16. Updated transfer use case tests:

* adapted fake repositories to new async `selectAccount` contract
* removed obsolete `clearSelectedAccount`

### Development environment

17. Updated Postman environment:

* changed `base_url`
* `192.168.0.16`
* to `192.168.0.20`

### Architectural impact

This commit moves the mobile application closer to a complete transactional flow architecture by:

* separating navigation concerns from UI state
* centralizing reusable financial widgets
* standardizing balance propagation through streams
* introducing typed route serialization
* reducing orchestration leakage into presentation widgets

The resulting structure is significantly more prepared for:

* transfer confirmation
* receipt visualization
* transaction replay
* transactional authentication
* future Zero Trust contextual flows
* state restoration during complex navigation stacks


## 2026/05/09 - mobile/internal-transfer-08

Refactor authentication refresh flow and improve internal transfer recipient UX.

This commit introduces a major update to the mobile authentication interceptor by implementing serialized refresh handling for concurrent `401` responses, while also aligning backend documentation and routing behavior with the actual refresh-token authentication model. Additionally, the internal transfer flow received UI refinements and component extraction to improve readability and interaction consistency.

### API

1. `api/cmd/api/main.go`

   * Removed JWT middleware from `POST /auth/refresh`
   * Clarified that refresh authentication is handled by the refresh token payload and session validation itself

2. Documentation updates

   * Updated authentication semantics across:

     * `api/docs/06-implementation.md`
     * `api/docs/07-api-rest.md`
     * `api/docs/08-auth_implementation.md`
     * `api/docs/ARCHITECTURE.md`
   * Replaced “JWT required” references for `/auth/refresh` with refresh-token-based authentication
   * Clarified refresh token lifecycle, validation behavior, and error semantics
   * Improved consistency between runtime behavior and architectural documentation
   * Refined wording around `INVALID_TOKEN` responses for refresh sessions

### Mobile Core Infrastructure

1. `mobile/lib/core/services/client_http/interceptors/auth/auth_interceptor.dart`

   * Implemented serialized refresh execution using a shared in-flight `Future`
   * Prevented duplicate refresh requests during concurrent `401` failures
   * Added retry optimization when another request already refreshed the access token
   * Extracted refresh execution into dedicated private methods:

     * `_refreshAccessToken`
     * `_performRefresh`
     * `_retryWithToken`
     * `_accessTokenUpdatedAfter`
     * `_bearerToken`
   * Added optional injected `refreshDio` for testability
   * Centralized session cleanup behavior
   * Improved error logging around refresh/retry failures
   * Reduced duplicated retry composition logic
   * Removed outdated concurrency-risk warning comments now that the locking strategy is implemented

2. `mobile/test/core/services/client_http/interceptors/auth/auth_interceptor_test.dart`

   * Added concurrency test coverage for refresh serialization
   * Validated:

     * single refresh execution under concurrent failures
     * successful retry behavior for all pending requests
     * token persistence updates
     * absence of unnecessary session cleanup
   * Added lightweight in-memory secure storage mock
   * Added fake HTTP adapter infrastructure for interceptor-level testing

3. `mobile/docs/ARCHITECTURE.md`

   * Documented shared in-flight refresh behavior
   * Removed obsolete “known concurrency risk” section
   * Added transfer routes reference

4. `mobile/lib/core/AGENT.md`

   * Updated interceptor guidance to reflect implemented concurrency protection
   * Reinforced importance of interceptor-specific tests

### Internal Transfer Flow

1. `mobile/lib/uis/pages/home/transfer/recipient_page.dart`

   * Refactored recipient page structure
   * Replaced repeated section labels with reusable `TextHeader`
   * Extracted dropdown and recipient summary into dedicated widgets
   * Improved selected-recipient visualization
   * Added dynamic “Prosseguir” button enable/disable behavior
   * Added custom enabled-state button styling
   * Improved account formatting consistency
   * Simplified recipient state handling with `ValueListenableBuilder`
   * Removed redundant `return` statements after snackbar handling

2. `mobile/lib/uis/pages/home/transfer/viewmodel/transfer_viewmodel.dart`

   * Replaced mutable selected recipient field with `ValueNotifier`
   * Improved UI reactivity around recipient selection
   * Adjusted auto-selection behavior when a single account is returned

3. Added reusable transfer widgets

   * `mobile/lib/uis/pages/home/transfer/widgets/dropdown_recipient.dart`

     * Extracted account selection dropdown
     * Added ellipsis handling and selected item builder
     * Standardized account rendering
   * `mobile/lib/uis/pages/home/transfer/widgets/recipient_card.dart`

     * Added summarized recipient preview card
     * Introduced reusable labeled information rows

### Shared UI Components

1. `mobile/lib/uis/core/text/text_header.dart`

   * Added reusable bold section header component

2. `mobile/lib/uis/core/text_form_field/basic_text_form_field.dart`

   * Added default bold text style for form input content
   * Added dedicated hint text styling
   * Improved visual consistency for transfer forms

This commit strengthens authentication robustness under concurrent network conditions while also improving transfer flow usability and maintainability across the mobile application.


## 2026/05/09 — mobile/internal-transfer-07

Refactor transfer recipient flow, strengthen CPF domain validation, and improve local bootstrap workflow consistency across API and mobile layers.

### Infrastructure and developer workflow

1. Makefile

   * Added `bootstrap` to `.PHONY`
   * Reworked `make bootstrap` into a complete first-run environment setup flow
   * Integrated:

     * `.env` initialization
     * Docker startup
     * database readiness wait
     * migrations execution
     * API startup
   * Updated `docker compose up` to use `--no-recreate` for safer repeated executions
   * Improved local development reproducibility and reduced onboarding friction

2. API Getting Started documentation

   * Expanded onboarding documentation with:

     * index/table of contents
     * Docker engine clarification
     * Colima support explanation
     * first-run bootstrap flow
     * troubleshooting guidance
   * Documented the new `make bootstrap` execution sequence
   * Clarified relationship between `bootstrap`, `setup`, and `docker-up`

### Customer domain and CPF validation

3. Customer domain

   * Introduced dedicated CPF validation module
   * Added:

     * CPF normalization
     * digit verification algorithm
     * repeated-digit rejection
     * formatting support (`123.456.789-09`)
   * Moved CPF validation responsibility fully into the domain layer through `NewCustomer`
   * Normalized stored CPF values before persistence

4. Customer domain errors

   * Added:

     * `ErrCPFInvalid`

5. Customer application error registry

   * Registered CPF invalid errors into shared HTTP/domain error mapping
   * Added standardized API response support for invalid CPF format

6. Register user use case

   * Refactored customer creation flow to use `customerdomain.NewCustomer`
   * Centralized customer invariant validation in the domain factory
   * Removed duplicated normalization logic from application layer

### Tests and validation coverage

7. CPF tests

   * Added complete CPF validation test suite covering:

     * formatted CPFs
     * invalid lengths
     * invalid verification digits
     * repeated digits
     * invalid characters
     * normalization behavior

8. Customer creation tests

   * Added validation tests for:

     * invalid CPF
     * empty CPF
     * empty customer name
     * formatted CPF normalization

9. Register user tests

   * Updated existing tests to use valid CPF values
   * Added invalid CPF registration scenario
   * Added assertions ensuring:

     * customer creation is aborted
     * hashing is not executed
     * user persistence is not triggered after CPF validation failure

### Mobile transfer flow restructuring

10. Routing

* Introduced dedicated transfer routing module:

  * `transfer_routes.dart`
* Added:

  * `TransferRoutes.recipient`
* Registered transfer routes separately from home routes
* Removed direct transfer route coupling from `home_routes.dart`

11. Home page

* Updated transfer navigation flow
* Replaced temporary placeholder navigation with recipient flow entrypoint

12. Recipient page

* Added new `RecipientPage`
* Implemented recipient lookup flow by:

  * CPF
  * branch + account number
* Added:

  * CPF input formatter
  * recipient dropdown selection
  * reactive account loading
  * recipient preview card
  * validation-triggered lookup behavior
  * focus management improvements
  * snackbar-based failure feedback

13. Transfer viewmodel

* Added recipient account state management
* Added:

  * `receipientAccounts`
  * `selectedRecipient`
* Refactored recipient retrieval into dedicated internal handler
* Added automatic recipient selection for single-account results

14. Transfer page

* Refactored ViewModel access using local getter
* Improved consistency and readability

### Mobile shared UI and repository cleanup

15. BasicTextFormField

* Made `labelText` optional
* Added `onChanged` support
* Increased flexibility for lightweight form compositions

16. Command result abstraction

* Renamed:

  * `data` → `value`
* Improved semantic consistency with `Result<T>`

17. Transaction repositories

* Refactored imports to absolute package-style imports
* Improved consistency across repository implementation files

18. Postman environment

* Updated local API base URL for current development environment

This commit consolidates three major directions:

* stronger domain invariants around CPF handling
* a more realistic internal transfer recipient flow in Flutter
* a cleaner and more reproducible local bootstrap/development workflow.


## 2026/05/08 - mobile/internal-transfer-06

Refactor internal transfer flow to use account IDs as the primary identity model across API and mobile layers, while introducing recipient lookup support and CPF/CNPJ validation utilities.

This commit aligns the mobile transfer implementation with the backend transfer contract evolution, simplifying recipient identification and removing legacy branch/account-based transfer semantics from the core transaction flow.

1. API contract and transfer domain adjustments

   * Removed `account_type` from internal transfer recipient responses and domain models.
   * Simplified recipient payloads to expose only:

     * `account_id`
     * `holder_name`
     * `document`
     * `branch`
     * `account_number`
   * Updated REST documentation for internal transfer recipient lookup responses.
   * Adjusted delivery and domain mapping layers to stop propagating account type information.

2. Internal transfer identity migration

   * Reworked transfer DTOs to use:

     * `from_account_id`
     * `to_account_id`
   * Removed legacy transfer fields:

     * `from_branch`
     * `from_account_number`
     * `to_branch`
     * `to_account_number`
   * Updated:

     * `TransferRequestDto`
     * `TransferResponseDto`
     * `TransferDraft`
     * `TransferUsecase`
     * repository validation logic
   * Simplified transfer execution flow by using account UUIDs as the canonical transaction identity.

3. API route evolution

   * Updated mobile transfer integration to consume:

     * `POST /accounts/internal-transfers`
   * Replaced legacy route:

     * `POST /accounts/transfer`
   * Renamed receipt retrieval API method:

     * `getTransferReceipt()` → `getReceipt()`

4. Internal recipient lookup implementation

   * Added support for:

     * `GET /accounts/internal-transfers/recipients`
   * Implemented:

     * `RecipientRequestDto`
     * `RecipientResponseDto`
     * `RecipientInfoDto`
   * Added query support for:

     * CPF/CNPJ document lookup
     * branch + account number lookup
   * Added robust envelope parsing and HTTP failure handling.
   * Added parsing safeguards for malformed backend payloads.

5. CPF/CNPJ validation utilities

   * Added `StringExtension` utilities:

     * `onlyNumbers`
     * `isValidCpf`
     * `isValidCnpj`
   * Implemented:

     * CPF verification digit validation
     * CNPJ verification digit validation
     * repeated digit rejection
     * formatted/unformatted document support
   * Centralized numeric sanitization logic for document processing.

6. Test coverage expansion

   * Updated all transfer-related tests to the new account ID model.
   * Added extensive coverage for:

     * internal recipient lookup
     * network failures
     * malformed payload handling
     * HTTP error mapping
     * empty recipient lists
     * DTO serialization changes
   * Added CPF/CNPJ validation tests covering:

     * formatted/unformatted values
     * invalid lengths
     * invalid check digits
     * repeated digits
   * Updated fake REST client support for GET operations and query parameter assertions.

7. Architectural consistency improvements

   * Reduced leakage of presentation-oriented banking concepts into transaction execution.
   * Consolidated account identity semantics around immutable UUIDs.
   * Improved separation between:

     * transfer execution identity
     * recipient display information
   * Kept branch/account number as presentation and lookup data instead of transactional identifiers.

This refactor prepares the mobile application for a more deterministic and scalable internal transfer model while improving API consistency, DTO clarity, and validation reliability across the transfer flow.


## 2026/05/08 - mobile/internal-transfer-05

Refactor internal transfer flow to use account IDs as the canonical identity model and implement recipient lookup support across API and Flutter layers.

This commit aligns the mobile application with the new internal transfer contract exposed by the API, replacing branch/account-number-based transfer execution with UUID-based account references while preserving recipient discovery through dedicated lookup endpoints.

### API

1. Simplified internal transfer recipient contract

   * Removed `account_type` from:

     * internal transfer recipient domain model
     * delivery DTOs
     * handler response mapping
     * REST API documentation
   * Reduced exposure of unnecessary banking metadata in recipient lookup responses.

2. Updated REST documentation

   * Adjusted internal transfer recipient payload examples.
   * Kept transfer semantics centered on:

     * `from_account_id`
     * `to_account_id`
   * Reinforced consistency with the current transfer domain model and ledger-oriented architecture.

### Mobile — Core

1. Added `StringExtension`

   * Introduced:

     * `onlyNumbers`
     * `isValidCpf`
     * `isValidCnpj`
   * Implemented CPF/CNPJ check digit validation algorithms.
   * Added complete unit coverage for:

     * formatted/unformatted values
     * invalid lengths
     * repeated digits
     * invalid check digits

### Mobile — Transfer API Integration

1. Refactored transfer requests to UUID-based account identity

   * Replaced:

     * `from_branch`
     * `from_account_number`
     * `to_branch`
     * `to_account_number`
   * With:

     * `from_account_id`
     * `to_account_id`

2. Updated transfer endpoint path

   * Changed:

     * `/accounts/transfer`
   * To:

     * `/accounts/internal-transfers`

3. Updated transfer DTOs

   * `TransferRequestDto`

     * now serializes only account IDs
   * `TransferResponseDto`

     * now exposes account IDs instead of branch/account-number identity

4. Updated transfer use case

   * Transfer execution now derives:

     * source account ID directly from authenticated account context
     * destination account ID from recipient selection

5. Simplified transfer validation

   * Repository validation now checks only:

     * `fromAccountId`
     * `toAccountId`

### Mobile — Recipient Lookup Flow

1. Added recipient lookup API support

   * Introduced:

     * `RecipientRequestDto`
     * `RecipientResponseDto`
     * `RecipientInfoDto`

2. Added lookup endpoint integration

   * Implemented:

     * `GET /accounts/internal-transfers/recipients`

3. Added recipient query strategies

   * Lookup by:

     * CPF/CNPJ document
     * branch + account number

4. Added envelope parsing and error handling

   * Implemented:

     * backend envelope validation
     * parsing error mapping
     * HTTP status validation
     * network failure propagation

### Mobile — Receipt API

1. Renamed receipt retrieval method

   * `getTransferReceipt`
   * → `getReceipt`

2. Updated repository integration and tests accordingly.

### Mobile — Tests

1. Expanded transfer API coverage

   * Added tests for:

     * endpoint path validation
     * request serialization
     * account-ID-based payloads
     * parsing failures
     * HTTP failures
     * network failures
     * malformed envelopes

2. Added recipient lookup test suite

   * Covered:

     * multiple recipients
     * empty responses
     * malformed payloads
     * backend envelope errors
     * query parameter generation
     * network failures

3. Updated DTO tests

   * Removed legacy branch/account-number assertions.
   * Added UUID-based transfer identity assertions.

4. Updated repository and receipt tests

   * Adjusted mocks and fixtures to reflect the new transfer contract.

This refactor consolidates internal transfer identity around stable account UUIDs while separating operational transfer execution from recipient discovery. The result is a cleaner transfer contract, lower exposure of transport-specific banking fields, and better alignment between the Flutter client and the current API/domain architecture.


## 2026/05/08 — api/internal-transfer-01

Refactor the internal transfer flow to operate with account UUIDs instead of branch/account-number pairs, and introduce a dedicated recipient lookup endpoint for internal transfers. This update also expands the transfer contract, improves authorization consistency, and aligns the API/documentation surface with the new transfer model.

### API and Route Changes

1. Updated internal transfer routes and naming conventions

   * Replaced `POST /accounts/transfer` with `POST /accounts/internal-transfers`
   * Added `GET /accounts/internal-transfers/recipients`
   * Updated route registration in `cmd/api/main.go`
   * Wired the new recipient lookup use case into the account handler

2. Refactored transfer payloads to use account UUIDs

   * Removed:

     * `from_branch`
     * `from_account_number`
     * `to_branch`
     * `to_account_number`
   * Added:

     * `from_account_id`
     * `to_account_id`
   * Updated request/response DTOs in delivery layer
   * Updated handler parsing and validation logic to parse UUIDs explicitly

3. Added internal transfer recipient lookup flow

   * Introduced lookup endpoint supporting:

     * branch + account number
     * CPF document lookup
   * Added normalization and validation rules
   * Added masked document support
   * Restricted exposed fields to confirmation-safe transfer metadata only

### Account Module

#### Application Layer

1. Added `LookupInternalTransferRecipients` use case

   * Supports:

     * lookup by branch/account number
     * lookup by CPF
   * Rejects:

     * mixed lookup modes
     * incomplete query combinations
     * unsupported document formats
   * Normalizes:

     * branch
     * account number
     * CPF document

2. Added helper normalization utilities

   * `normalizeBranch`
   * `normalizeAccountNumber`
   * `onlyDigits`
   * Reused shared document normalization logic

3. Added authorization enforcement

   * Lookup flow now validates authenticated customer access before querying recipients

#### Domain Layer

1. Added `TransferRecipient` domain structure

   * Includes:

     * account ID
     * holder name
     * masked document
     * branch
     * account number
     * optional account type

2. Added document normalization and masking utilities

   * `NormalizeDocument`
   * `MaskDocument`

#### Infrastructure Layer

1. Added recipient lookup repository operations

   * `FindTransferRecipientsByBranchAndNumber`
   * `FindTransferRecipientsByDocument`

2. Added active-account filtering to lookup queries

   * Queries now explicitly restrict recipient accounts to active accounts only

3. Added centralized recipient query mapper

   * Shared recipient scan/mapping logic
   * Automatic document masking before returning data

#### Delivery Layer

1. Added recipient lookup HTTP handler

   * Parses query params
   * Validates authenticated user
   * Maps domain recipients to response DTOs
   * Returns stable API envelope responses

2. Added recipient response DTOs

   * `InternalTransferRecipientsData`
   * `InternalTransferRecipientData`

### Transfer Use Case Refactor

1. Refactored transfer orchestration to use account IDs directly

   * Removed source/destination lookup by branch/number
   * Uses:

     * `GetByIDForUpdate`
     * deterministic UUID lock ordering

2. Improved authorization semantics

   * Ownership validation now occurs immediately after source account resolution
   * Replay/idempotency authorization now depends on validated ownership

3. Simplified transaction flow

   * Reduced duplicate account lookup logic
   * Unified balance updates around UUID-based operations

4. Updated idempotency scope

   * Explicitly scoped to:

     * `from_account_id`
     * `idempotency_key`

5. Updated ledger relationship handling

   * Related account references now use UUIDs directly
   * Transfer replay reconstruction preserved

### Tests

1. Added complete unit coverage for recipient lookup use case

   * Successful account lookup
   * Successful CPF lookup
   * Multiple-account CPF responses
   * Invalid query combinations
   * Forbidden access
   * Repository failure propagation
   * Empty-result behavior

2. Added delivery tests for recipient lookup endpoint

   * Authentication validation
   * Success responses
   * Response field restrictions
   * Forbidden scenarios
   * Invalid parameter handling
   * Empty result responses

3. Added repository tests for recipient lookup queries

   * Active account filtering validation
   * Document masking validation
   * Query argument validation

4. Refactored transfer tests to UUID-based payloads

   * Updated application tests
   * Updated delivery tests
   * Updated integration wiring tests
   * Added legacy payload rejection coverage

### Documentation

1. Updated API documentation

   * Added:

     * internal transfer recipient lookup section
     * new internal transfer contract
   * Replaced all branch/account-number transfer examples with UUID-based examples
   * Added lookup examples and response semantics
   * Expanded Postman environment documentation

2. Updated architecture and implementation documents

   * Route listings
   * transfer flow descriptions
   * startup wiring references
   * REST surface explanations
   * authentication-protected route lists

3. Updated presentation and overview documentation

   * Adjusted transfer endpoint references
   * Added recipient lookup flow descriptions
   * Clarified internal-transfer-only scope

### Mobile Layer

1. Added internal transfer recipient DTOs

   * `InternalTransferRecipientDto`
   * `InternalTransferRecipientLookupResponseDto`
   * `InternalTransferRecipientLookupQueryDto`

2. Added query serialization support

   * Account-based lookup serialization
   * CPF-based lookup serialization

3. Added DTO parsing and validation tests

   * Optional account type handling
   * Multiple account parsing
   * Prohibited-field exposure protection
   * Empty-result handling

This commit consolidates the transition from branch/account-number transfer execution to a UUID-driven internal transfer architecture, introduces a dedicated recipient discovery flow, strengthens authorization boundaries, and aligns the API, tests, documentation, and mobile DTO contracts around the new transfer model.


## 2026/05/08 - mobile/internal-transfer-04

Refactor the internal transfer page structure and introduce reusable UI components for account selection and section rendering.

This update focuses on improving the organization and maintainability of the transfer flow screen by extracting repeated UI structures into dedicated widgets and simplifying controller lifecycle management.

### Mobile

1. Refactored `transfer_page.dart`

   * Replaced local helper widgets with reusable components:

     * `AccountDropdown`
     * `SectionTitle`
   * Migrated relative imports to absolute project imports for better consistency.
   * Simplified `TextEditingController` initialization using direct field initialization.
   * Added `viewModel.initialize()` execution during `initState`.
   * Removed obsolete private widget builders:

     * `_buildSectionTitle`
     * `_buildOriginAccountDropdown`
   * Replaced hardcoded account dropdown items with dynamic account rendering from `viewModel.accounts`.
   * Improved screen readability by separating transfer sections into explicit reusable widgets.
   * Added TODO note indicating future reactive integration for account selection state.

2. Updated `transfer_viewmodel.dart`

   * Simplified `initialize()` from asynchronous to synchronous execution.
   * Kept UUID v7 idempotency generation isolated in initialization flow.

3. Added `widgets/account_dropdown.dart`

   * Introduced reusable account selection widget.
   * Added support for dynamic account rendering using `AccountSummaryResponseDto`.
   * Centralized dropdown styling and selection behavior.
   * Standardized account display format using:

     * branch
     * account number

4. Added `widgets/section_title.dart`

   * Extracted section title rendering into a reusable stateless widget.
   * Centralized typography styling for transfer form sections.

5. General Improvements

   * Reduced UI duplication inside transfer page implementation.
   * Improved component isolation and future extensibility.
   * Prepared the transfer flow for reactive state evolution and integration with notifier-based widgets.
   * Improved readability and separation of responsibilities within the transfer feature module.

This refactor establishes a cleaner foundation for the internal transfer flow, making the UI structure more modular and easier to evolve as the banking operations and reactive state management continue to grow.


## 2026/05/07 — mobile/internal-transfer-04

Refined the Flutter mobile architecture organization by formalizing dependency injection entrypoints, consolidating use case orchestration patterns, and improving transfer workflow coordination across repositories, use cases, and view models.

Updated the dependency injection structure and naming conventions across the mobile architecture.

1. Dependency injection and module organization

   * Renamed `data/data.dart` to `data/repositories.dart`
   * Renamed `uis/uis.dart` to `uis/viewmodels.dart`
   * Introduced `domain/usecases/usecases.dart` as the dedicated registration entrypoint for domain workflows
   * Updated `core/config/dependencies.dart` to use the new registration pipeline:

     * `CoreServices`
     * `Services`
     * `Repositories`
     * `Usecases`
     * `Viewmodels`
   * Added explicit `Usecases.add(injector)` bootstrap integration
   * Standardized dependency registration terminology across the project

2. Use case architecture consolidation

   * Expanded the architectural guidance for `domain/usecases`
   * Formalized the intended orchestration flow:

     * `UI -> ViewModel -> UseCase -> Repository -> API/Service -> RestClient -> Dio`
   * Clarified when workflows should remain in repositories versus when they should become dedicated use cases
   * Added implementation guidance for:

     * reusable orchestration
     * multi-repository coordination
     * app-facing workflow inputs
     * request mapping ownership
   * Added architectural restrictions preventing:

     * UI state leakage into use cases
     * repository duplication
     * transport-layer dependencies inside use cases

3. Domain organization improvements

   * Reorganized the conceptual structure of `domain/common`
   * Updated references from:

     * `domain/auth/...`
     * `domain/enums/...`
   * To:

     * `domain/common/auth/...`
     * `domain/common/user/...`
     * `domain/common/receipt/...`
   * Clarified separation between:

     * stable app-facing domain models
     * workflow orchestration use cases
   * Added stronger framework isolation rules for `domain/common`

4. Transfer workflow improvements

   * Improved `TransferUsecase` validation flow
   * Added defensive validation for missing selected accounts
   * Replaced unsafe nullable access on `selectedAccount`
   * Added transfer receipt retrieval orchestration:

     * `getTransferReceipt`
   * Added account selection workflow orchestration:

     * `selectAccount`
   * Added validation and error handling for invalid account selection
   * Delegated balance loading after account selection through the use case layer

5. Transfer view model enhancements

   * Added transfer receipt command support
   * Added account selection command support
   * Connected new use case operations into the view model command system
   * Extended transfer workflow state orchestration

6. Repository API refinements

   * Refactored `AccountRepository.selectAccount`

     * changed from object-based selection to ID-based selection
   * Improved account cache lookup behavior
   * Prevented unnecessary balance loading during account selection
   * Added cache cleanup when clearing selected account
   * Reset balance cache during account deselection

7. Repository documentation improvements

   * Added detailed documentation comments to:

     * `AccountRepository`
     * `AuthRepository`
     * `TransactionRepository`
   * Clarified:

     * cache semantics
     * session behavior
     * authentication invariants
     * transfer constraints
     * statement retrieval behavior
   * Improved repository contract readability for future contributors

8. Mobile architecture documentation updates

   * Updated:

     * `.github/instructions/*`
     * `mobile/lib/*/AGENT.md`
     * `mobile/docs/ARCHITECTURE.md`
     * `mobile/README.md`
   * Standardized references to:

     * `repositories.dart`
     * `viewmodels.dart`
     * `usecases/usecases.dart`
   * Documented the new dependency registration flow and architectural responsibilities
   * Added clearer guidance about use case placement and orchestration responsibilities

9. Import and structure cleanup

   * Normalized import ordering in repositories
   * Fixed relative import inconsistencies
   * Improved readability and architectural consistency across modules

This commit consolidates the mobile architecture around explicit dependency registration boundaries, establishes a clearer use case orchestration model, and strengthens the internal transfer workflow foundation for future transactional and security-related features.


## 2026/05/07 - mobile/internal-transfer-03

Implemented the first complete internal transfer flow structure in the Flutter mobile application, including routing, UI foundation, transfer orchestration, account selection support, and idempotency preparation aligned with the backend transactional model.

Also improved local development ergonomics for macOS environments using Colima.

### Infrastructure and Development Environment

1. Updated `Makefile`

   * Added automatic Colima detection before `docker compose up`
   * Ensured Docker daemon startup in macOS environments using Colima
   * Improved local developer experience and reduced manual environment setup friction

### Routing and Navigation

2. Updated `mobile/lib/core/routing/routes.dart`

   * Added `HomeRoutes.transfer`
   * Introduced dedicated navigation path for transfer operations

3. Updated `mobile/lib/core/routing/routes/home_routes.dart`

   * Registered `TransferPage`
   * Added dependency injection wiring for `TransferViewmodel`
   * Extended GoRouter configuration with transfer navigation support

4. Updated `mobile/lib/uis/pages/home/home_page.dart`

   * Replaced placeholder transfer action with actual navigation flow
   * Added initialization execution during `initState`
   * Removed redundant `didPush` initialization logic
   * Integrated navigation using `context.pushNamed`

### Transfer Domain and Use Case Layer

5. Added `mobile/lib/domain/usecases/transfer/inputs/transfer_draft.dart`

   * Introduced immutable transfer draft structure
   * Added support for:

     * origin/destination data
     * amount handling with `money2`
     * optional description
     * idempotency key propagation
   * Implemented `copyWith` pattern for immutable state evolution

6. Added `mobile/lib/domain/usecases/transfer/transfer_usecase.dart`

   * Implemented transfer orchestration use case
   * Connected account and transaction repositories
   * Added DTO conversion from domain draft to API request
   * Centralized transfer execution logic
   * Exposed selected account and available accounts from repository layer

### Repository Improvements

7. Updated `mobile/lib/data/repositories/account/account_repository.dart`

   * Added cached accounts exposure through `accounts` getter

8. Updated `mobile/lib/data/repositories/account/account_repository_impl.dart`

   * Added in-memory account cache
   * Persisted loaded accounts for reuse across flows
   * Prepared repository layer for account-origin selection in transfers

### Transfer UI Foundation

9. Added `mobile/lib/uis/pages/home/transfer/transfer_page.dart`

   * Created initial transfer screen
   * Added structured sections for:

     * origin account
     * beneficiary data
     * transfer amount
   * Implemented:

     * dropdown account selection
     * branch/account input fields
     * amount input
     * beneficiary input
   * Added confirmation action placeholder
   * Used reusable `BasicTextFormField` components
   * Structured layout for future validation and execution integration

10. Added `mobile/lib/uis/pages/home/transfer/viewmodel/transfer_viewmodel.dart`

    * Introduced transfer state orchestration layer
    * Integrated `Command1` async execution pattern
    * Added UUID v7 idempotency generation
    * Prepared deterministic retry behavior aligned with backend transfer guarantees
    * Centralized transfer command execution logic

### Dependency Injection

11. Updated `mobile/lib/uis/uis.dart`

    * Registered `TransferViewmodel` in dependency injection container

### UI and Theme Refinements

12. Updated `mobile/lib/uis/app_widget.dart`

    * Replaced `EB Garamond` with `Google Sans`
    * Improved overall visual consistency for application typography

13. Updated `mobile/lib/uis/core/text_form_field/basic_text_form_field.dart`

    * Reduced border radius from `24` to `8`
    * Improved visual alignment with banking-style UI patterns
    * Standardized input appearance for future financial flows

### Tooling and API Testing

14. Updated `tools/postman/Environment.postman_environment.json`

    * Updated local API base URL for current development environment

This commit establishes the initial mobile transfer architecture and prepares the application for full transactional integration with the backend transfer pipeline, including idempotent execution semantics and reusable account context handling.


## 2026/05/07 — mobile/internal-transfer-02

Refine mobile transfer architecture, repository boundaries, and DTO usage strategy while introducing the first transaction repository implementation for internal transfers and receipts.

### Architectural and Documentation Updates

* Updated mobile architecture and agent instruction documents to clarify when DTOs are acceptable across repository and view model boundaries.
* Standardized the guidance that domain models should only exist when they add semantic meaning, behavior, aggregation, or decoupling from unstable contracts.
* Explicitly documented that curated app-facing DTOs using idiomatic Dart types (`Money`, `DateTime`, enums) are valid application-facing contracts.
* Reinforced the restriction against leaking low-level transport concerns such as:

  * raw JSON maps
  * backend envelopes
  * Dio types
  * HTTP handling details
  * snake_case payload structures
* Updated:

  * `.github/instructions/*`
  * `mobile/lib/data/AGENT.md`
  * `mobile/lib/domain/AGENT.md`
  * `mobile/lib/data/repositories/AGENT.md`
  * `mobile/lib/data/services/apis/AGENT.md`
  * `mobile/docs/ARCHITECTURE.md`

### Result API Improvements

#### `mobile/lib/core/result/result.dart`

* Added nullable `value` getter documentation explaining its intended use for lightweight success access patterns such as caching successful responses without explicit folding.

### Transaction Repository Layer

#### `mobile/lib/data/repositories/transaction/transaction_repository.dart`

Added the new transaction repository contract responsible for:

* transfer execution
* transfer receipt retrieval
* exposing cached transfer state

Introduced:

* `lastTransfer`
* `lastReceipt`
* `transfer()`
* `getTransferReceipt()`

The repository intentionally exposes DTOs directly as curated app-facing contracts aligned with the updated architectural direction.

### Transaction Repository Implementation

#### `mobile/lib/data/repositories/transaction/transaction_repository_impl.dart`

Implemented the first transaction repository orchestration layer.

Features include:

* integration with:

  * `ApiTransfer`
  * `ApiReceipt`
* lightweight application-level validation before API execution
* transfer success caching
* receipt success caching
* automatic cache invalidation on failures

Added defensive validations for:

* missing source account
* missing destination account
* zero or negative amounts

Implemented consistent failure propagation using `Result` and `AppError`.

### Dependency Injection Wiring

#### `mobile/lib/data/data.dart`

Registered:

* `TransactionRepository`
* `TransactionRepositoryImpl`
* `ApiTransfer`
* `ApiReceipt`

into the mobile dependency injection graph.

This completes the first transaction orchestration vertical slice for the Flutter client.

### Automated Tests

#### `mobile/test/data/repositories/transaction/transaction_repository_impl_test.dart`

Added extensive repository test coverage for:

#### Transfer Flow

* successful transfer execution
* cache persistence after success
* cache clearing after failure
* validation failures before API execution
* source account validation
* destination account validation
* amount validation

#### Receipt Flow

* successful receipt retrieval
* receipt cache persistence
* cache invalidation after backend failures
* propagation of:

  * 404 not found
  * 403 forbidden
  * generic backend failures

#### Test Infrastructure

Added:

* fake API implementations
* noop HTTP client
* reusable DTO builders
* helper money factory methods

The tests validate both repository orchestration behavior and architectural boundary expectations.

### Architectural Direction

This commit formalizes an important mobile architecture decision:

* DTOs are now treated as first-class application-facing contracts when:

  * the backend contract is intentionally designed for the mobile app
  * fields are already idiomatic in Dart
  * no additional semantic abstraction is necessary

This avoids redundant domain model duplication while preserving clear transport isolation boundaries.

The result is a leaner and more pragmatic mobile architecture with lower mapping overhead and clearer repository/view-model contracts.


## 2026/05/07 - api/transfer-description-01

Expanded transfer and receipt support across the mobile layer while refining transfer receipt semantics and standardizing money transport conversions.

1. Transfer and receipt API services

   * Added `ApiTransfer` service for `POST /accounts/transfer`
   * Added `ApiReceipt` service for `GET /accounts/transfer/{transaction_reference}/receipt`
   * Registered both services in dependency injection
   * Standardized API envelope parsing and HTTP error handling for transfer flows
   * Added parsing failure handling with explicit `AppErrorCode.parsingError`

2. Transfer DTO contracts

   * Added `TransferRequestDto` with branch/account-number-based transfer identity
   * Added `TransferResponseDto` for transfer result parsing
   * Added `TransferReceiptResponseDto` for transfer receipt parsing
   * Standardized money scalar serialization/deserialization using:

     * `ApiParse.toInt`
     * `ApiParse.toMoney`
   * Preserved transport identity through branch + account number instead of internal UUID exposure
   * Added optional `idempotency_key` serialization support

3. Transfer receipt domain modeling

   * Added `TransferReceiptStatus` enum with:

     * `completed`
     * `pending`
     * `failed`
     * `cancelled`
     * `rejected`
   * Added semantic helpers:

     * `isSuccess`
     * `isPending`
     * `isFailed`
   * Added strict parsing through `TransferReceiptStatus.fromString`
   * Documented future-compatible receipt status evolution

4. Mobile domain reorganization

   * Migrated domain models into `domain/common/...`
   * Added:

     * `domain/common/auth/models`
     * `domain/common/user/enums`
     * `domain/common/receipt/enums`
   * Moved `UserRole` into `domain/common/user/enums/user_role.dart`
   * Updated imports across repositories, APIs, view models, and domain models
   * Refined architecture documentation for:

     * `domain/common`
     * `domain/usecases`
     * future domain growth organization

5. API parsing standardization

   * Replaced manual money serialization helpers with:

     * `ApiParse.toInt(Money)`
   * Standardized guidance across:

     * AGENT instructions
     * architecture docs
     * mobile data layer instructions
     * API service instructions
   * Explicitly documented that DTOs must not hand-roll money scalar conversions

6. Transfer and receipt DTO test coverage

   * Added `TransferRequestDto` serialization tests
   * Added `TransferResponseDto` parsing tests
   * Added `TransferReceiptResponseDto` parsing tests
   * Validated:

     * money parsing semantics
     * enum parsing behavior
     * UTC date parsing
     * idempotency key serialization
     * absence of internal account/customer identifiers
     * required field failures
     * runtime protection against accidental DTO field leakage

7. API documentation updates

   * Expanded transfer receipt status documentation in `api/docs/07-api-rest.md`
   * Added explicit status semantics for:

     * `completed`
     * `pending`
     * `failed`
     * `cancelled`
     * `rejected`
   * Documented current backend behavior returning `completed`
   * Added forward-compatibility notes for future receipt state evolution
   * Fixed transfer receipt route anchor escaping in the table of contents

8. Repository and documentation organization

   * Moved Mermaid-generated assets into `docs/mermaid-images`
   * Updated mobile architecture documentation to reflect the new domain structure
   * Refined agent guidance for domain placement and application organization

This commit establishes the first complete mobile-side transfer and receipt integration baseline, introduces explicit receipt lifecycle semantics, and standardizes financial scalar handling across the API contract and Flutter data layer.


## 2026/05/07 - mobile/internal-transfer-01

Introduced the first complete mobile-side internal transfer API integration layer, including transfer execution, transfer receipt retrieval, DTO validation coverage, and domain structure normalization for future growth.

### Mobile API and DTO Enhancements

1. Added transfer API service infrastructure:

   * Created `ApiTransfer` for `POST /accounts/transfer`
   * Implemented envelope parsing and HTTP error handling
   * Added parsing failure protection with structured `AppError`

2. Added transfer receipt API service:

   * Created `ApiReceipt` for `GET /accounts/transfer/{transaction_reference}/receipt`
   * Implemented envelope parsing and status validation
   * Added standardized HTTP/parsing error handling

3. Added transfer request DTO:

   * Created `TransferRequestDto`
   * Added serialization via `toMap`
   * Added optional `idempotency_key` support
   * Standardized monetary serialization using `ApiParse.toInt`

4. Added transfer response DTO:

   * Created `TransferResponseDto`
   * Added parsing for balances and transferred amount
   * Standardized monetary parsing using `ApiParse.toMoney`

5. Added transfer receipt response DTO:

   * Created `TransferReceiptResponseDto`
   * Added parsing for:

     * transfer status
     * operation metadata
     * source/destination account presentation data
     * transaction reference
     * operation timestamp
   * Added conversion to `TransferReceiptStatus`

6. Updated API parsing utilities:

   * Replaced `moneyToBigInt` with `ApiParse.toInt`
   * Standardized transport scalar conversion strategy for Money types

### Domain Structure Refactor

1. Reorganized domain root structure for scalability:

   * Introduced:

     * `domain/common`
     * `domain/usecases`
   * Documented the new structure across architecture and agent files

2. Moved auth models into contextual domain folders:

   * `domain/auth/models/auth_user.dart`
     → `domain/common/auth/models/auth_user.dart`
   * `domain/auth/models/user_profile.dart`
     → `domain/common/auth/models/user_profile.dart`

3. Moved and normalized user role enum:

   * `domain/enums/user_role.dart`
     → `domain/common/user/enums/user_role.dart`

4. Added transfer receipt domain enum:

   * Created `TransferReceiptStatus`
   * Added:

     * parsing helper
     * semantic helpers (`isSuccess`, `isPending`, `isFailed`)
     * documented operational semantics

5. Updated imports throughout repositories, APIs, and UI layers to follow the new domain structure.

### Dependency Injection and Service Registration

1. Updated service registration:

   * Added `ApiTransfer`
   * Added `ApiReceipt`
   * Registered both in `mobile/lib/data/services/services.dart`

### Documentation and Architecture Updates

1. Updated mobile architecture documentation:

   * Documented new domain folder strategy
   * Clarified separation between:

     * stable app-facing models
     * workflow orchestration use cases

2. Updated AGENT instructions:

   * Added standardized Money conversion rules
   * Enforced usage of:

     * `ApiParse.toInt`
     * `ApiParse.toMoney`
   * Added domain placement conventions for:

     * `common/<area>/models`
     * `common/<area>/enums`
     * `usecases`

3. Updated REST API documentation:

   * Added transfer receipt status semantics
   * Documented current and future-compatible status values
   * Clarified backend behavior for persisted receipts

4. Fixed markdown anchor escaping for:

   * `/accounts/transfer/{transaction_reference}/receipt`

5. Relocated generated Mermaid assets:

   * moved `mermaid-images/*`
     → `docs/mermaid-images/*`

### Test Coverage

1. Added `TransferRequestDto` tests:

   * serialization validation
   * Money conversion validation
   * idempotency serialization behavior
   * prevention of internal ID exposure

2. Added `TransferResponseDto` tests:

   * Money parsing validation
   * payload validation
   * protection against leaking internal account IDs

3. Added `TransferReceiptResponseDto` tests:

   * status parsing
   * Money parsing
   * timestamp parsing
   * required field validation
   * unknown status rejection
   * protection against internal account/customer ID leakage

This commit establishes the first complete mobile transfer transport layer, standardizes Money transport serialization rules, and reorganizes the mobile domain structure into a scalable context-oriented layout aligned with future workflow orchestration growth.


## 2026/05/06 - api/transfer-by-account-number-04

Refine transfer-by-account-number behavior, strengthen idempotency validation, and expand transfer integration coverage.

1. Updated transfer API documentation

   * Clarified idempotency semantics for transfer operations:

     * idempotency is now explicitly scoped to the resolved source account plus `idempotency_key`
     * different source accounts may reuse the same key independently
   * Added detailed error scenarios for:

     * invalid request payloads
     * invalid amount values
     * malformed branch/account number data
     * account not found
     * insufficient funds
     * inactive account
     * malformed transfer receipt reference
   * Updated Postman environment documentation:

     * replaced transfer account UUID variables with public branch/account number variables
     * documented `transaction_reference` usage for receipt retrieval
   * Updated onboarding/testing flow examples to reflect transfer-by-account-number behavior.

2. Expanded transfer application test coverage

   * Added support for custom idempotency lookup behavior in `transferTxMock` through `getTransactionByKeyFn`
   * Added regression coverage validating:

     * identical idempotency keys are allowed across different source accounts
     * idempotency replay remains scoped to the resolved source account only
     * financial effects execute exactly once for valid independent transfers
     * new ledger references are generated for distinct source accounts.

3. Added transfer integration test suite

   * Implemented full integration coverage for:

     * transfer execution
     * idempotent replay behavior
     * persisted balance validation
     * transfer receipt retrieval
   * Added end-to-end validation for:

     * transfer requests using branch/account number resolution
     * replay preservation of historical balances
     * transaction reference propagation
     * receipt endpoint correctness
     * recipient name and account number mapping
   * Added response safety validation ensuring internal account IDs are not exposed in transfer responses.

4. Added transfer integration test infrastructure helpers

   * Added helpers for:

     * transfer test customer/account seeding
     * account insertion
     * customer insertion
     * transfer cleanup
     * transfer request execution
     * unique CPF/account number generation
   * Added cleanup handling for:

     * transfer ledger rows
     * accounts
     * customers.

5. Expanded transfer delivery handler tests

   * Added authentication validation coverage:

     * transfer endpoint now explicitly tested for missing authenticated user context
   * Added error mapping validation for:

     * forbidden access
     * account not found
     * insufficient funds
     * inactive account
   * Verified proper HTTP status mapping and stable API error codes.

6. Improved transfer API robustness and contract validation

   * Reinforced public transfer flow centered on:

     * branch + account number resolution
     * ledger-backed deterministic replay
     * account-scoped idempotency
     * receipt retrieval through transaction reference
   * Strengthened validation around transport-layer contract behavior and protected internal identifiers from API consumers.

This update consolidates the transfer-by-account-number flow as the public transfer contract while improving idempotency correctness, integration reliability, and API-level validation coverage.


## 2026/05/06 - api/transfer-by-account-number-03

Implemented transfer receipt retrieval flow and exposed transaction references across transfer operations.

1. Added transfer receipt use case and delivery endpoint

   * Created `GetTransferReceipt` application use case
   * Added ownership validation for transfer receipts
   * Implemented `GET /accounts/transfer/{transaction_reference}/receipt`
   * Added transfer receipt response DTOs
   * Registered new route in API bootstrap
   * Added HTTP handler for receipt retrieval
   * Added support for `TRANSACTION_NOT_FOUND`

2. Extended transfer result contract with transaction references

   * Added `TransactionReference` to `TransferResult`
   * Returned transfer reference from successful operations
   * Preserved transaction reference during idempotency replay
   * Exposed `transaction_reference` in transfer HTTP responses

3. Implemented transfer receipt persistence query

   * Added `GetTransferReceiptByReference` repository contract
   * Implemented PostgreSQL query joining:

     * transfer_out
     * transfer_in
     * accounts
     * customers
   * Reconstructed persisted receipt data from ledger entries
   * Added transaction not found mapping from `pgx.ErrNoRows`

4. Expanded domain model for receipt support

   * Added `TransferReceipt` domain structure
   * Added `ErrTransactionNotFound`
   * Updated repository interfaces
   * Improved domain documentation comments for:

     * account validation methods
     * transaction constructors

5. Improved delivery layer and response contracts

   * Added `TransferReceiptData`
   * Added RFC3339 formatting for operation dates
   * Added authorization enforcement for receipt access
   * Added invalid UUID validation for transaction references

6. Expanded test coverage

   * Added application tests for:

     * source customer access
     * destination customer access
     * forbidden access
     * invalid reference
     * transaction not found
   * Added delivery tests for:

     * successful receipt retrieval
     * unauthorized access
     * invalid references
     * forbidden access
     * not found responses
   * Updated transfer tests to validate transaction references
   * Updated mocks to support receipt repository contract

7. Updated API documentation

   * Added Transfer Receipt endpoint documentation
   * Added new error scenarios and payload examples
   * Added `TRANSACTION_NOT_FOUND` to error catalog
   * Updated table of contents and endpoint indexes
   * Updated transfer success payload examples with `transaction_reference`

8. Improved repository documentation and transaction comments

   * Added descriptive comments for repository methods
   * Added transaction lifecycle documentation
   * Added commit/rollback behavior documentation
   * Added repository delegation comments

9. Updated Postman environment

   * Updated `base_url` to local network IP for external device testing
   * Reformatted exported environment JSON

This change completes the first version of persisted transfer receipt retrieval, enabling clients to recover transfer metadata from immutable ledger data using the public transaction reference identifier.


## 2026/05/06 - api/transfer-by-account-number-02

Refactor internal transfer identification to use `(branch, account_number)` instead of UUID-based account references, while aligning API documentation, repository contracts, database constraints, and tooling behavior with the new transfer model.

### Documentation

1. Updated `api/docs/07-api-rest.md`

   * Replaced transfer request and response payloads from:

     * `from_account_id`
     * `to_account_id`
       to:
     * `from_branch`
     * `from_account_number`
     * `to_branch`
     * `to_account_number`
   * Clarified that transfers are currently internal-only and bank resolution is implicit.
   * Updated idempotency scope from:

     * `(from_account_id, idempotency_key)`
       to:
     * `(from_branch, from_account_number, idempotency_key)`
   * Refined ownership documentation for admin users and `GET /accounts` customer-context behavior.
   * Standardized invalid query parameter handling from `INVALID_DATA` to `INVALID_REQUEST`.
   * Updated transfer validation error descriptions to reflect branch/account-number semantics instead of UUID validation.

### Transaction Delivery Layer

2. Updated `api/internal/account/transaction/delivery/data.go`

   * Added `NewTransferData` DTO using:

     * branch
     * account number
     * amount
     * resulting balances

3. Updated `api/internal/account/transaction/delivery/request.go`

   * Added `NewTransferRequest` with transfer identification based on:

     * branch
     * account number
     * idempotency key

4. Updated `api/internal/account/transaction/delivery/handler.go`

   * Added documentation comment to `requireUser`.

### Domain and Repository Contracts

5. Updated `api/internal/account/transaction/domain/domain.go`

   * Extended repository contract with:

     * `GetByBranchAndNumberForUpdate(ctx, branch, number)`

### Infrastructure Layer

6. Updated `api/internal/account/transaction/infrastructure/base_repository.go`

   * Implemented:

     * `GetByBranchAndNumberForUpdate`
   * Added row-locking lookup using:

     * `WHERE branch = $1 AND number = $2 FOR UPDATE`
   * Preserved transactional locking semantics for transfer consistency.
   * Added repository-level documentation comments for:

     * account locking
     * balance updates
     * transaction creation
     * idempotency replay lookups
     * reference-based transaction retrieval

7. Updated `api/internal/account/transaction/infrastructure/repository.go`

   * Exposed `GetByBranchAndNumberForUpdate` through:

     * `Repository`
     * `txRepository`

### Tests

8. Updated `deposit_test.go`

   * Added mock implementation for:

     * `GetByBranchAndNumberForUpdate`

9. Updated `transfer_test.go`

   * Added mock implementation for:

     * `GetByBranchAndNumberForUpdate`

### Database Migrations

10. Added `api/migrations/000002_account_number_key.up.sql`

    * Replaced unique constraint:

      * `accounts_number_key`
        with:
      * `accounts_branch_number_key`
        enforcing uniqueness on:
      * `(branch, number)`

11. Added `api/migrations/000002_account_number_key.down.sql`

    * Restored previous unique constraint on:

      * `number`

### Tooling and Environment Automation

12. Updated `infra/scripts/update-mobile-env-ip.sh`

* Added automatic Postman environment synchronization.
* Added support for updating:

  * `base_url`
  * `app_token`
    inside:
  * `tools/postman/Environment.postman_environment.json`
* Added `jq`-based JSON update flow.
* Added automatic extraction of `APP_ACCESS_TOKEN` from `mobile/dev.env`.
* Improved developer workflow consistency between:

  * mobile
  * API
  * Postman environments

### Architectural Impact

This change shifts transfer addressing from opaque UUID references toward banking-oriented identifiers, preparing the API for more realistic transfer semantics while preserving:

* transactional integrity
* deterministic locking
* ledger consistency
* idempotent replay behavior

The new `(branch, account_number)` lookup model also creates a cleaner path for future expansion into:

* inter-bank routing
* PIX-style abstractions
* external transfer gateways
* customer-facing account operations.


## 2026/05/06 - api/transfer-by-account-number-01

Refactor transfer identification to use `(branch, account_number)` instead of account UUIDs and align repository, database, and API contracts with the new transfer model.

### Documentation Updates

1. Updated REST API documentation (`api/docs/07-api-rest.md`)

   * Replaced transfer payload fields from `from_account_id` / `to_account_id` to:

     * `from_branch`
     * `from_account_number`
     * `to_branch`
     * `to_account_number`
   * Documented that transfers are currently internal-only.
   * Clarified that account resolution is performed through `(branch, account_number)`.
   * Updated idempotency scope to:

     * `(from_branch, from_account_number, idempotency_key)`
   * Adjusted transfer response examples to expose branch and account number instead of UUIDs.
   * Refined `INVALID_DATA` semantics for malformed branch/account information.
   * Improved ownership and authorization documentation for admin behavior and customer-scoped account listing.
   * Standardized unexpected query parameter errors from `INVALID_DATA` to `INVALID_REQUEST`.

### Domain and Repository Changes

2. Extended transaction repository contract (`api/internal/account/transaction/domain/domain.go`)

   * Added:

     * `GetByBranchAndNumberForUpdate(ctx, branch, number)`
   * Enables transactional account resolution using business identifiers instead of UUIDs.

3. Updated PostgreSQL repository implementation

   * Added `GetByBranchAndNumberForUpdate` to:

     * `base_repository.go`
     * `repository.go`
   * Implemented row locking with:

     * `SELECT ... FOR UPDATE`
   * Added proper error mapping to `ErrAccountNotFound`.

4. Improved infrastructure documentation/comments

   * Added explanatory comments for:

     * account locking
     * balance updates
     * transaction persistence
     * idempotency replay lookups
     * helper methods

### Transfer Delivery Layer Preparation

5. Added new transfer request/response transport models

   * `NewTransferRequest`
   * `NewTransferData`
   * Introduced branch/account-number-based payload structures while preserving existing transfer DTOs during migration.

6. Improved delivery handler readability

   * Added documentation comment to `requireUser`.

### Database Migration

7. Added migration to support branch-scoped account uniqueness

   * New migration:

     * `000002_account_number_key.up.sql`
     * `000002_account_number_key.down.sql`
   * Replaced:

     * `UNIQUE(number)`
   * With:

     * `UNIQUE(branch, number)`

This enables reuse of account numbers across different branches and aligns the schema with the new transfer lookup strategy.

### Test Adjustments

8. Updated transaction-related mocks/tests

   * Added `GetByBranchAndNumberForUpdate` support to:

     * deposit mocks
     * transfer mocks
     * transactional repository mocks
   * Keeps repository interfaces compatible after contract expansion.

### Tooling Improvements

9. Enhanced mobile/Postman environment synchronization script

   * Updated `infra/scripts/update-mobile-env-ip.sh`
   * Added automatic Postman environment synchronization for:

     * `base_url`
     * `app_token`
   * Added safe handling for:

     * missing environment files
     * missing `jq`
   * Extracts token values automatically from `mobile/dev.env`.

### Result

This change begins the transition from internal UUID-based transfers to bank-style account addressing using branch and account number identifiers. It also prepares the persistence and API layers for future inter-bank evolution while preserving transactional consistency and idempotent replay semantics.


## 2026/05/06 - api/split-account-04

Refactor the account module by separating bank account and transaction responsibilities into distinct bounded domains and infrastructure layers.

This change removes the old shared `internal/account/domain` abstraction and establishes clearer ownership boundaries between:

* `bankaccount`
* `transaction`
* `statement`

The refactor also reduces interface coupling, simplifies repository contracts, and aligns the implementation more closely with the modular monolith architecture.

### Main Changes

1. Account Domain Separation

   * Moved account entity ownership to `internal/account/bankaccount/domain`
   * Moved transaction entity ownership to `internal/account/transaction/domain`
   * Removed the old shared `internal/account/domain` package entirely
   * Split domain errors between bank account and transaction contexts
   * Preserved business invariants while reducing cross-module coupling

2. Transaction Module Isolation

   * Introduced dedicated transaction repository contracts:

     * `Repository`
     * `Tx`
   * Created standalone transaction infrastructure layer:

     * `transaction/infrastructure/base_repository.go`
     * `transaction/infrastructure/repository.go`
   * Migrated:

     * balance updates
     * row locking
     * idempotency handling
     * transaction persistence
     * transaction orchestration
   * Isolated transfer consistency logic inside the transaction module

3. Infrastructure Refactor

   * Removed old generic account infrastructure:

     * `base_repository.go`
     * `repository.go`
     * `tx_repository.go`
   * Updated runtime wiring in `cmd/api/main.go`
   * Replaced shared repository initialization with dedicated transaction repository initialization
   * Simplified repository responsibilities per module

4. Error Registry Reorganization

   * Moved account error registry to:

     * `internal/account/errors/registry.go`
   * Added independent registration for:

     * bank account errors
     * transaction errors
   * Updated bootstrap initialization to use the new registry package

5. Statement Module Decoupling

   * Updated statement application and infrastructure layers to:

     * consume bank account entities from `bankaccount/domain`
     * consume ledger transactions from `transaction/domain`
   * Adjusted repository contracts accordingly
   * Preserved cursor pagination and statement semantics

6. Test Cleanup and Simplification

   * Removed obsolete mock implementations from:

     * account tests
     * statement tests
     * admin tests
   * Eliminated unnecessary transaction-related methods from bank account mocks
   * Reduced unused imports and dead test code
   * Updated tests to use the new domain package boundaries

7. Runtime and Dependency Wiring

   * Updated all transaction use cases:

     * Deposit
     * Withdraw
     * Transfer
   * Rewired use cases to depend on:

     * `transaction/domain.Repository`
   * Updated integration tests and delivery handlers to use the new infrastructure path

8. Architectural Improvements

   * Reinforced modular monolith boundaries
   * Reduced accidental coupling between read/write concerns
   * Clarified ownership of:

     * account lifecycle
     * transaction ledger
     * statement queries
   * Improved long-term maintainability for future evolution of:

     * ledger consistency
     * transactional security
     * financial operations

### Removed Components

The following legacy structures were removed:

* `internal/account/domain`
* old shared account repository contracts
* old generic account transaction infrastructure
* duplicated repository responsibilities
* obsolete transaction-related test scaffolding

### Result

The account subsystem is now structured around explicit module boundaries:

* `bankaccount` handles account lifecycle and ownership
* `transaction` handles financial mutations and ledger consistency
* `statement` handles read/query behavior

This significantly improves separation of concerns while preserving the existing transactional guarantees and consistency model already defined by the system architecture.


## 2026/05/06 - api/split-account-03

Refactor the account module structure by introducing the `bankaccount` bounded package separation and isolating transactional repository responsibilities from account lifecycle operations.

This change reorganizes the account-related codebase to better distinguish:

* account lifecycle management
* transactional balance operations
* delivery/application/domain boundaries

The refactor reinforces the modular monolith direction already established in the architecture and prepares the project for stronger feature isolation inside the `account` module.

1. Updated runtime wiring and dependency composition

   * Refactored imports in `api/cmd/api/main.go`
   * Migrated account lifecycle dependencies to:

     * `internal/account/bankaccount/application`
     * `internal/account/bankaccount/delivery`
     * `internal/account/bankaccount/infrastructure`
   * Introduced a dedicated transaction repository instance:

     * `transactionAccountRepo := accountTransactionInfrastructure.New(db)`
   * Updated transaction use cases (`Deposit`, `Withdraw`, `Transfer`) to use the transaction-oriented repository instead of the generic account repository

2. Removed obsolete duplicated application-layer files

   * Removed old duplicated files from:

     * `internal/account/application`
   * Deleted:

     * `access_policy.go`
     * `access_policy_test.go`
     * `auth_test.go`
   * Consolidated responsibility into the new `bankaccount` structure

3. Introduced the new `bankaccount` module structure

   * Migrated account lifecycle components to:

     * `internal/account/bankaccount/application`
     * `internal/account/bankaccount/delivery`
     * `internal/account/bankaccount/domain`
     * `internal/account/bankaccount/infrastructure`
   * Renamed and relocated:

     * access policies
     * branch policy
     * create account use case
     * balance use case
     * account listing use case
     * HTTP handlers
     * request/response DTOs
     * tests

4. Updated imports to reflect the new domain boundary

   * Replaced imports from:

     * `internal/account/domain`
   * With:

     * `internal/account/bankaccount/domain`
   * Applied across:

     * application layer
     * delivery layer
     * admin approval flow
     * auth integration tests

5. Added a dedicated `bankaccount/domain` compatibility layer

   * Created:

     * `bankaccount/domain/account.go`
     * `bankaccount/domain/repository.go`
   * Introduced type aliases for:

     * `Account`
     * `Transaction`
     * `AccountStatus`
     * `TransactionType`
     * `Tx`
   * Re-exported:

     * account status constants
     * domain errors
     * `NewAccount`
   * Added the new repository contract abstraction:

     * `Repository`
     * `AccountRepository`

6. Added a dedicated bank account repository implementation

   * Implemented:

     * `internal/account/bankaccount/infrastructure/repository.go`
   * Added support for:

     * account creation
     * account lookup
     * customer account listing
     * existence checks
     * account number generation
   * Introduced a generic `executor` abstraction supporting:

     * pooled connections
     * transactional execution compatibility
   * Improved repository separation between:

     * account state persistence
     * financial transaction orchestration

7. Updated admin approval flow dependencies

   * Refactored `ApproveUser` use case imports
   * Updated tests to use the new `bankaccount` package hierarchy
   * Preserved atomic account creation behavior during approval flow

8. Updated auth integration setup

   * Adjusted integration server wiring to use:

     * `bankaccount` repository for account lifecycle
     * transaction repository for balance-changing operations
   * Improved dependency clarity in integration tests

9. Architectural impact

   * Strengthens feature-oriented modularization inside `internal/account`
   * Separates:

     * account aggregate responsibilities
     * transaction processing responsibilities
   * Reduces accidental coupling between:

     * account lifecycle operations
     * ledger/balance mutation logic
   * Moves the project closer to:

     * clearer bounded contexts
     * more explicit repository contracts
     * safer future extraction paths if the module evolves further

10. Test and compatibility considerations

* Preserved existing behavior through type aliasing compatibility
* Updated affected tests and imports
* Maintained runtime compatibility while reorganizing package ownership

This refactor significantly improves the internal organization of the account module while preserving the current runtime behavior and transactional guarantees of the banking operations.


## 2026/05/06 - api/split-account-02

Refactor statement feature into an isolated account statement module.

This change separates statement responsibilities from the core account module, introducing dedicated `statement` application, delivery, domain, and infrastructure layers. The refactor strengthens module boundaries and aligns the implementation with the modular monolith architecture adopted by the project.

Implemented changes:

1. Refactored statement package structure

   * Moved statement application use cases from:

     * `internal/account/application/statement`
     * to `internal/account/statement/application`
   * Added dedicated statement layers:

     * `statement/domain`
     * `statement/delivery`
     * `statement/infrastructure`
   * Isolated statement responsibilities from account delivery and repository concerns

2. Introduced dedicated statement repository contract

   * Added `statement/domain.Repository`
   * Extracted statement-specific operations from `domain.AccountRepository`
   * Removed transaction query responsibilities from the account repository abstraction
   * Reduced coupling between account balance operations and statement retrieval

3. Added dedicated statement infrastructure repository

   * Implemented `statement/infrastructure.Repository`
   * Moved transaction query logic from account infrastructure into the new statement repository
   * Preserved cursor pagination, date filtering, and transaction mapping behavior
   * Kept statement query semantics unchanged

4. Introduced dedicated statement delivery handler

   * Created `statement/delivery.Handler`
   * Moved:

     * statement endpoint logic
     * query parsing
     * cursor validation
     * response mapping
     * authenticated user extraction
   * Added statement-specific DTOs:

     * `StatementData`
     * `StatementItemData`
     * `StatementCursorData`

5. Simplified account delivery layer

   * Removed statement endpoint handling from:

     * `account_handler.go`
     * `handler.go`
   * Removed statement DTOs from `account_data.go`
   * Removed statement parsing helpers from account delivery
   * Simplified `accountdelivery.New(...)` constructor signature

6. Updated runtime composition and route registration

   * Added statement repository initialization in `cmd/api/main.go`
   * Injected statement repository into `NewGetStatement`
   * Added dedicated `statementHandler`
   * Redirected:

     * `GET /accounts/{id}/statement`
     * to the new statement delivery module

7. Migrated and reorganized tests

   * Moved statement handler tests into:

     * `internal/account/statement/delivery`
   * Added isolated statement delivery test setup
   * Added local auth helper for statement delivery tests
   * Removed statement-related mocks and tests from account handler tests
   * Preserved existing statement behavior validation

8. Reduced account module responsibilities

   * Removed `GetTransactions(...)` from:

     * `domain.AccountRepository`
     * account infrastructure repositories
     * transactional repositories
   * Reinforced separation between:

     * balance mutation operations
     * ledger query operations

9. Architectural impact

   * Reinforces feature-oriented modularization inside the account bounded context
   * Improves cohesion by grouping statement-specific concerns together
   * Reduces accidental coupling between transactional account operations and ledger querying
   * Creates a cleaner evolution path for future statement-specific features:

     * exports
     * filters
     * analytics
     * caching
     * asynchronous projections

The resulting structure better reflects the architectural direction of the project: a layered modular monolith with explicit responsibility boundaries and domain-oriented organization.


## 2026/05/06 - api/split-account-01

Refactor transaction operations into a dedicated transaction module and delivery layer.

This change separates deposit, withdraw, and transfer responsibilities from the account delivery layer, introducing a clearer modular boundary for transaction-related operations. The refactor improves package organization, aligns the codebase with the modular monolith architecture, and prepares the project for future expansion of transaction-specific behaviors.

### Main Changes

1. Transaction module extraction

   * Moved transaction application use cases from:

     * `internal/account/application/transaction`
   * To:

     * `internal/account/transaction/application`
   * Preserved all existing business logic and tests during the move.
   * Updated imports throughout the project to reflect the new module structure.

2. Dedicated transaction delivery layer

   * Added:

     * `internal/account/transaction/delivery`
   * Implemented a dedicated `Handler` responsible for:

     * `Deposit`
     * `Withdraw`
     * `Transfer`
   * Moved request DTOs into the transaction module:

     * `DepositRequest`
     * `WithdrawRequest`
     * `TransferRequest`
   * Moved transaction response DTOs into the transaction module:

     * `TransferData`

3. Account handler simplification

   * Removed transaction-related responsibilities from:

     * `internal/account/delivery/account_handler.go`
   * The account handler now focuses only on:

     * account creation
     * account listing
     * balance retrieval
     * statement retrieval
   * Removed transaction use case dependencies from:

     * `internal/account/delivery/handler.go`
   * Simplified `accountDelivery.New(...)` constructor signature.

4. Application wiring update

   * Updated `cmd/api/main.go`:

     * introduced `transactionHandler`
     * split route ownership between account and transaction handlers
   * Transaction routes now use:

     * `transactionHandler.Deposit`
     * `transactionHandler.Withdraw`
     * `transactionHandler.Transfer`

5. Test migration and isolation

   * Moved transaction handler unit tests into the new transaction delivery package.
   * Added:

     * transaction-specific mocks
     * auth test helpers
     * test bootstrap setup
   * Migrated deposit integration tests into:

     * `internal/account/transaction/delivery`
   * Updated integration server wiring in auth authorization integration tests.

6. Cleanup and modular consistency

   * Removed obsolete transaction DTOs and interfaces from the account delivery package.
   * Reduced coupling between account delivery and transaction orchestration.
   * Improved alignment with the documented layered architecture and module boundaries.

### Architectural Impact

This refactor strengthens the modular monolith organization by distinguishing:

* account lifecycle responsibilities
* transaction execution responsibilities

The new structure better reflects the financial domain model, where transactions are a distinct operational concern with their own request flows, validation rules, and delivery behavior.

The result is a cleaner separation of concerns, improved maintainability, and a more scalable foundation for future transaction features such as:

* transaction-specific middleware
* auditing
* fraud analysis
* transactional security policies
* asynchronous integrations
* richer transfer orchestration

### Result

The API now exposes transaction operations through a dedicated transaction delivery module while preserving the existing HTTP contract and transactional behavior. The refactor reduces accidental coupling inside the account module and improves long-term maintainability of the financial operation flows.


## 2026/05/06 - api/docs-04

Add comprehensive GoDoc documentation across account, auth, admin, customer, infrastructure, and shared layers

This commit improves the internal documentation quality of the Bank API by adding extensive GoDoc comments throughout the codebase. The changes focus on clarifying responsibilities, execution flow, transactional guarantees, authorization behavior, repository contracts, and HTTP handler semantics.

Key improvements include:

1. Account module documentation

   * Added documentation for access policies (`CanAccessCustomer`, `CanAccessAccount`)
   * Documented account creation, balance retrieval, statement retrieval, deposit, withdraw, and transfer use cases
   * Added explanations for concurrency helpers such as deterministic UUID ordering and transfer ledger reconstruction
   * Documented branch policy responsibilities
   * Added repository interface documentation covering:

     * transaction lookup
     * idempotency behavior
     * transactional execution
     * balance mutation semantics
     * statement pagination behavior
   * Added GoDoc to domain entities and business rule methods:

     * `NewAccount`
     * `CanDeposit`
     * `CanWithdraw`
     * `CanTransfer`
     * transaction constructors
   * Documented HTTP handlers and parsing helpers for:

     * deposits
     * withdrawals
     * transfers
     * statements
     * balances
     * query parameter parsing
   * Added infrastructure repository documentation for:

     * SQL persistence behavior
     * row locking semantics
     * transaction creation
     * account existence checks
     * atomic balance updates

2. Transaction and consistency clarification

   * Added comments explaining:

     * ledger replay strategy
     * transfer pair reconstruction
     * deadlock reduction through deterministic lock ordering
     * transactional guarantees
   * Reinforced append-only ledger concepts already described in architecture documentation 
   * Clarified consistency strategy aligned with the concurrency documentation 

3. Auth module documentation

   * Added detailed GoDoc for:

     * login flow
     * registration flow
     * refresh token rotation
     * current user retrieval
     * JWT middleware behavior
     * authenticated context helpers
   * Documented repository behavior for:

     * users
     * sessions
     * transactions
   * Clarified executor resolution logic for transactional contexts
   * Added explanations for:

     * email normalization
     * password validation
     * session revocation
     * refresh token persistence
   * Improved middleware documentation for bearer token extraction and authentication context propagation

4. Admin module documentation

   * Added GoDoc for:

     * user approval use case
     * approval transaction semantics
     * admin handler behavior
     * error registry registration
   * Clarified atomic relationship between:

     * user activation
     * account creation
   * Reinforced guarantees already defined in the lifecycle and consistency documentation 

5. Customer module documentation

   * Added documentation for:

     * customer creation flow
     * authenticated customer retrieval
     * customer repository behavior
     * customer entity construction
     * HTTP handlers
   * Clarified PostgreSQL error handling and validation mapping behavior

6. Shared infrastructure and error handling documentation

   * Added GoDoc for:

     * shared error registry
     * error mapping strategy
     * internal error fallback behavior
     * JSON response helpers
   * Clarified response standardization and mapping rules aligned with the API error specification 

7. Repository and transaction infrastructure improvements

   * Added extensive documentation for:

     * `Repository`
     * `txRepository`
     * transactional execution helpers
     * commit and rollback semantics
     * nested transaction protection
   * Clarified responsibilities of:

     * `BeginTx`
     * `WithTransaction`
     * `runInTransaction`

8. General architectural alignment

   * The added documentation now better reflects the implemented architecture and layer responsibilities described in the architecture and implementation documents:

     * modular monolith structure 
     * application orchestration model
     * transactional use case execution

The result is a substantially more self-documented codebase, improving:

* onboarding for contributors
* maintainability
* IDE-assisted navigation
* architectural clarity
* long-term documentation consistency between implementation and project docs.


## 2026/05/06 - mobile/docs-03

Add mobile architecture instruction system and Postman integration documentation

This commit introduces a complete instruction-based guidance structure for the Flutter mobile application, expands the mobile architecture documentation, and adds a reusable Postman setup for API validation and onboarding workflows.

1. Added layered mobile instruction system under `.github/instructions`

   * Created top-level mobile architecture guidance with:

     * dependency direction rules
     * architectural boundaries
     * routing and state management conventions
     * dependency injection guidance
     * Result/AppError usage patterns
   * Added specialized instruction files for:

     * core infrastructure layer
     * data layer
     * repositories
     * API services
     * domain layer
     * use cases
     * UI layer
     * pages/view models
     * shared UI core components
   * Standardized architecture expectations around:

     * `UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`
     * constructor injection
     * Result-based error handling
     * separation of UI, orchestration, and transport concerns
   * Documented:

     * folder placement rules
     * DTO responsibilities
     * AppError mapping
     * Command usage
     * route lifecycle handling
     * repository ownership rules
     * caching and session ownership
     * shared UI promotion rules
   * Added `applyTo` patterns for coding-agent-aware contextual guidance.

2. Expanded `mobile/docs/ARCHITECTURE.md`

   * Added a full table of contents for navigation.
   * Added the new “Copilot Instructions Mirror (.github/instructions)” section.
   * Documented:

     * relationship between `AGENT.md` files and `.github/instructions`
     * maintenance synchronization expectations
     * layer-specific instruction file references
   * Improved discoverability of architecture and coding conventions.

3. Added Postman collection and environment support

   * Added `tools/postman/Banklab_API.postman_collection.json`.
   * Added `tools/postman/Environment.postman_environment.json`.
   * Added `tools/postman/README.md`.
   * Included ready-to-use requests for:

     * register
     * login
     * refresh
     * auth/me
     * admin approval
   * Added environment variables for:

     * `base_url`
     * `app_token`
     * `access_token`
     * `refresh_token`
     * `account_id`
     * `account_id_2`
     * `id`
   * Standardized onboarding/testing flow for local API validation.

4. Updated API REST documentation

   * Added a new “Postman Setup” section to `api/docs/07-api-rest.md`.
   * Documented:

     * repository Postman files
     * environment variables
     * import/configuration flow
     * recommended request execution order
   * Added operational notes for:

     * `X-App-Token`
     * Bearer authentication
     * token refresh behavior

5. Documentation alignment improvements

   * Reinforced consistency between:

     * architecture documentation
     * coding-agent instructions
     * mobile layer responsibilities
     * API onboarding workflow
   * Improved onboarding support for:

     * contributors
     * coding agents
     * API consumers
     * mobile/frontend development workflows

This update significantly improves project maintainability, architectural clarity, and contributor onboarding while establishing a consistent instruction-driven workflow for the Flutter mobile application and API integration ecosystem.


## 2026/05/04 — mobile/docs-01

Introduces a **comprehensive documentation and guidance system for the mobile layer**, consolidating architectural definitions, development rules, and AI-assisted coding guidelines across all layers of the Flutter application.

### 1. Documentation Structure — Expansion and Consolidation

* Expanded `mobile/docs/ARCHITECTURE.md` into a **complete architectural reference**:

  * clarified layer responsibilities (UI, Domain, Data, Core)
  * introduced explicit flow models for simple and complex workflows
  * documented dependency graph and injection order
  * formalized rules for use cases, repositories, and API services
* Improved readability and structure with:

  * consistent sectioning
  * clearer separation of responsibilities
  * explicit conventions for evolution

### 2. Getting Started — Developer Guidance Integration

* Updated `mobile/docs/00-getting_started.md` to include:

  * mandatory reference to `AGENT.md` guides before code changes
  * direct links to all layer-specific agent documents
* Added **Architecture Notes section**:

  * defines canonical execution flows:

    * `UI -> ViewModel -> Repository -> API -> RestClient -> Dio`
    * optional use case insertion for complex workflows
  * reinforces separation of concerns at runtime

### 3. Agent-Based Development Model (Major Addition)

* Introduced a **full set of `AGENT.md` guides** across the project:

  * `core`, `data`, `domain`, `uis`
  * sub-guides for repositories, APIs, pages, and shared UI
* These guides define:

  * strict dependency boundaries
  * responsibilities per layer
  * rules for error handling (`Result`, `AppError`)
  * constraints for Flutter, Dio, and storage usage
* Establishes a **controlled environment for AI-assisted development**, reducing architectural drift

### 4. Core Layer — Formalization

* Rewrote `mobile/lib/core/AGENT.md` with:

  * explicit role as cross-cutting infrastructure layer
  * detailed breakdown of:

    * HTTP abstraction (`RestClient`)
    * interceptors (auth flow, refresh logic)
    * dependency injection rules
  * strict prohibition of dependencies on higher layers
* Formalizes core as **infrastructure boundary**, not a feature layer

### 5. Data Layer — Clear Separation of Concerns

* Reworked `mobile/lib/data/AGENT.md` and added:

  * dedicated repository guide
  * dedicated API service guide
* Key improvements:

  * repositories defined as **app-facing orchestration boundary**
  * APIs defined as **transport-level boundary**
  * DTO placement and ownership clarified
* Enforces rule:

  * UI never calls APIs directly
  * repositories abstract all data operations

### 6. Domain Layer — Stabilization Rules

* Expanded `mobile/lib/domain/AGENT.md`:

  * defines domain as **framework-agnostic layer**
  * separates domain models from DTOs
  * clarifies when parsing belongs in domain vs API
* Added `domain/usecases/AGENT.md`:

  * introduces use cases as **controlled orchestration layer**
  * defines when to introduce them vs keeping logic in view models

### 7. UI Layer — Strict Presentation Boundaries

* Rewrote `mobile/lib/uis/AGENT.md` and added:

  * `pages/AGENT.md`
  * `core/AGENT.md` (shared UI primitives)
* Key rules enforced:

  * UI owns Flutter concerns only
  * view models coordinate commands, not infrastructure
  * repositories/use cases are the only data entry points
* Formalizes:

  * command-based async handling
  * navigation responsibilities
  * widget lifecycle ownership

### 8. Architectural Consistency Model

* The documentation now defines a **cohesive end-to-end model**:

  * layered architecture with strict boundaries
  * predictable data flow
  * explicit orchestration points (ViewModel vs UseCase)
* Aligns mobile architecture with backend principles such as:

  * separation of concerns
  * deterministic flows
  * controlled evolution paths 

### 9. Development Impact

* Establishes:

  * a **single source of truth** for mobile architecture
  * enforceable rules for contributors and AI tools
  * reduced ambiguity in where logic should live
* Improves:

  * maintainability
  * onboarding clarity
  * consistency across features

### Conclusion

This commit transitions the mobile project from **implicit architectural conventions to an explicit, enforceable system**, introducing a structured documentation and agent-guided model that significantly reduces ambiguity and increases long-term scalability.


## 2026/05/03 — docs/update-12

Introduces a formal definition for project organization in GitHub and refines the system overview diagram to better represent architectural boundaries and composition.

### 1. Contribution Guidelines

* Added `CONTRIBUTING.md` defining the project board structure
* Established three classification dimensions:

  * **Type**: Feature, Improvement, Bug, Research
  * **Area**: aligned with system modules (Auth, Account, Ledger, Customer, Mobile, Web, Security, Infra)
  * **Priority**: High, Medium, Low
* Clarified that:

  * every issue must include all three dimensions
  * Research issues precede implementation
  * Area represents primary ownership, not full impact scope
* This aligns backlog organization with the modular architecture of the system 

### 2. Documentation — System Overview

* Refactored the main architecture diagram in `00-visao_geral.md`
* Introduced explicit grouping:

  * **BankLab** as the top-level system boundary
  * **BankAPI** as an internal subgraph
* Adjusted flow layout to:

  * emphasize separation between client, API, and database
  * reflect the modular monolith structure more clearly
* Improved visual readability by enforcing horizontal flow (`direction LR`)

### 3. Architectural Consistency

* The updated diagram better communicates:

  * the API as a centralized runtime component
  * the database as a consistency boundary
  * the client as an external consumer
* This representation is now more consistent with the documented architecture and execution model, particularly the layered and modular structure described in the implementation and architecture documents 

### Conclusion

This commit strengthens both collaboration and communication aspects of the project by formalizing backlog organization and improving architectural visualization. It reduces ambiguity in issue management while making the system structure clearer for contributors and stakeholders.


## 2026/05/01 — docs/update-11

Refactors application bootstrap, improves infrastructure organization, and enhances internal documentation across authentication, security, and middleware layers. The changes focus on **centralization, clarity, and maintainability**, without altering core business behavior.

### 1. Bootstrap Refactor — Configuration Centralization

* Removed `Config` and `LoadConfig` from `main.go`
* Introduced `bootstrap.LoadConfig()` as the single source of configuration
* Enforced **fail-fast validation** for:

  * `APP_TOKEN`
  * `JWT_SECRET`
* Consolidates startup concerns into the bootstrap module, aligning with the architectural boundary defined for initialization logic 

### 2. Application Initialization Improvements

* `main.go` now delegates:

  * environment loading
  * error registration
  * configuration retrieval
* Reduces coupling between runtime wiring and infrastructure concerns
* Improves readability of the application entrypoint

### 3. Documentation — Presentation and Positioning

* Updated API architecture presentation:

  * improved system diagram with clearer module boundaries
  * explicit separation between **BankLab** and **BankAPI**
* Added contextual reference to ZTA work (Aurix – Feb 2026)
* Expanded and refined MIT license justification:

  * clearer positioning for fintech context
  * explicit trade-off between permissiveness and future flexibility

### 4. Middleware Documentation Enhancements

* Added detailed comments to JWT middleware:

  * `RequireAuth`
  * `OptionalAuth`
  * token extraction behavior
* Clarifies runtime behavior:

  * strict authentication vs opportunistic authentication
* Improves maintainability and onboarding for future contributors

### 5. Security Infrastructure Documentation

* Added documentation to:

  * `JWTTokenService`

    * access token generation
    * refresh token generation
    * parsing logic
  * `BcryptPasswordHasher`

    * cost handling
    * hashing and comparison behavior
* Makes explicit the responsibilities of the infrastructure layer, which is consistent with the architectural separation of concerns 

### 6. Shared Middleware — App Token

* Documented `AppToken` middleware:

  * explains constant-time comparison
  * clarifies role in onboarding protection
* Reinforces the boundary between:

  * system entry control (AppToken)
  * user authentication (JWT) 

### 7. Bootstrap Enhancements

* Added documentation for:

  * `Init()` lifecycle
  * `RegisterErrors()` centralized mapping
  * `.env` resolution strategy across multiple paths
* Strengthens consistency in error handling across the system, aligning with the defined response standard 

### 8. Database Layer Documentation

* Documented `database.NewPool()`:

  * clarifies its role as connection factory
  * reinforces database as a **first-class architectural component**
* Aligns with the system’s design principle where the database acts as a consistency boundary 

### 9. Overall Architectural Impact

* No changes to domain rules, use cases, or API contracts
* Improvements are strictly:

  * structural (bootstrap centralization)
  * documentary (clarity and intent)
* Reinforces the modular monolith discipline:

  * clearer boundaries
  * better separation of responsibilities
  * improved long-term maintainability

### Conclusion

This commit is a **structural and documentary consolidation step**. It reduces duplication, clarifies responsibilities, and strengthens architectural consistency without introducing behavioral changes.

From an engineering perspective, this is a necessary refinement: the system evolves from “working code” to **explicitly documented and maintainable architecture**, which is critical for scaling both the codebase and collaboration.



## 2026/04/30 — docs/update-10

Refactors the account infrastructure layer, introduces structured transaction handling, and expands documentation tooling with PDF generation support. Also improves context propagation and architectural clarity across the codebase.

### 1. Infrastructure Layer — Repository Refactor

* Extracted core persistence logic into `baseRepository`

  * isolates SQL operations and shared behavior
  * removes duplication between transactional and non-transactional flows
* Introduced new `Repository` implementation:

  * delegates all operations to `baseRepository`
  * centralizes database pool (`pgxpool.Pool`) usage
* Introduced `txRepository`:

  * encapsulates transaction-scoped execution (`pgx.Tx`)
  * enforces boundary: no nested transactions allowed
* Reorganized transaction control:

  * `BeginTx` returns a domain-level `Tx`
  * `WithTransaction` executes use cases with commit/rollback guarantees
* This refactor aligns with the architectural direction of **explicit transaction boundaries at the application layer**, preserving consistency guarantees 

### 2. Removal of Legacy Coupling

* Eliminated duplicated method implementations from previous `Repository`
* Removed implicit mixing of:

  * connection pool logic
  * transaction logic
* Result:

  * clearer separation between execution context (pool vs tx)
  * improved maintainability and testability

### 3. Context Propagation Enhancements

* Improved `authctx` package:

  * simplified `GetAuthenticatedUser` with early-return pattern
  * supports both value and pointer storage safely
  * introduced cleaner `RequireAuthenticatedUser` flow
* Added comprehensive documentation:

  * clarifies role of authenticated user in request lifecycle
* These changes reinforce the pattern where **authentication context is explicit and propagated through layers**, not hidden dependencies 

### 4. Database Context Utilities

* Expanded `tx_context.go`:

  * added documentation and usage examples
  * formalized pattern for embedding `pgx.Tx` into `context.Context`
* Enables transaction-aware operations without polluting method signatures

### 5. API Documentation Updates

* Updated `README.md`:

  * added new endpoints:

    * `POST /admin/users/{id}/approve`
    * `GET /accounts`
  * included architecture presentation material reference
* Keeps documentation aligned with current API surface and flows 

### 6. Documentation Tooling — Book Generation

* Added Makefile target:

  * `make book-pt` to generate PDF documentation via `pandoc`
* Introduced:

  * `metadata.yaml` (document metadata)
  * Eisvogel LaTeX template and supporting files
* Enables:

  * structured documentation pipeline
  * reproducible generation of technical material (book format)
* Added Mermaid diagram assets for architectural visualization

### 7. Project Configuration Improvements

* Updated `.gitignore`:

  * ignores generated `.pdf` artifacts
* Introduced template directory and LaTeX customization:

  * improves control over formatting, typography, and layout

### 8. Compile-Time Guarantees

* Added explicit interface assertions:

  * ensures repository implementations satisfy domain contracts at compile time
* Improves reliability and prevents silent contract drift

### Conclusion

This commit significantly improves the **structural integrity of the infrastructure layer** while introducing a robust documentation pipeline.

From an architectural standpoint, the most relevant evolution is the **clear separation between execution context (pool vs transaction)** combined with explicit transaction orchestration. This aligns the implementation more closely with the system’s consistency model and reinforces the database as the central coordination boundary for financial operations.


## 2026/04/22 — mobile/balance-04

Refines the mobile balance flow by introducing typed money handling, centralized feedback components, and route-aware home refresh behavior. This update improves monetary consistency, simplifies navigation contracts, and makes the home balance experience more reliable after login and route transitions. 

### 1. Monetary handling and parsing

* Added `money2` as a direct dependency and introduced a centralized currency registry in `app_currencies.dart`
* Defined `BRL` as the default application currency, while keeping `USD`, `EUR`, and `BTC` available for future compatibility
* Added `ApiParse` helpers to convert API numeric payloads into `BigInt` and `Money`
* Updated `BalanceResponseDto` and `StatementItemDto` to use `Money` instead of raw integers
* Centralized conversion from API values to money objects, reducing formatting duplication and improving type safety for balance and statement data

### 2. Home route and navigation behavior

* Added a global `RouteObserver` and registered it in the router
* Simplified the home route by removing `accountId` extraction from route `extra` and query parameters
* Adjusted `HomePage` to no longer depend on a route-supplied account id
* This change makes the home screen less coupled to navigation payload details and more aligned with the authenticated app state

### 3. Home page lifecycle and automatic balance refresh

* Converted `HomePage` state to `RouteAware`
* Subscribed the page to route lifecycle events through the shared route observer
* Triggered initialization on `didPush` and `didPopNext`
* Stopped periodic refresh on `didPop` and `didPushNext`
* Added a timer in `HomeViewmodel` to reload the balance every 20 seconds after initialization
* Added explicit `startTimer`, `stopTimer`, and `dispose` methods to manage polling safely

### 4. Balance tile refactor

* Refactored `BalanceTile` to receive typed DTOs instead of preformatted strings
* Replaced manual currency formatting with `Money.format()`
* Updated the supporting label to show agency and account information directly from the account DTO
* Slightly adjusted the card gradient to soften the visual presentation

### 5. Snackbar standardization

* Added reusable `AppSnackbar` with support for `success`, `error`, and `info` variants
* Centralized snackbar styling, duration, and presentation behavior
* Replaced direct `ScaffoldMessenger` usage in `LoginPage`
* Updated pending-feature feedback on the home screen to use the same component
* This creates a more consistent feedback pattern across the app

### 6. Login flow adjustment

* Updated login result handling to use `AppSnackbar`
* Moved navigation to `context.goNamed(HomeRoutes.home.routeName)` after successful command completion handling
* Changed `_submit()` to trigger the command without duplicating result processing logic afterward
* Added mounted checks before navigation

### 7. General outcome

* The balance feature now uses a stronger monetary representation end to end
* The home screen refresh behavior is better synchronized with navigation lifecycle
* UI feedback is more consistent and reusable
* The codebase becomes cleaner by removing manual money formatting and reducing route-state complexity

### Conclusion

This commit strengthens the mobile balance experience by combining typed money support, cleaner route handling, periodic balance refresh, and standardized user feedback. The result is a more robust and maintainable implementation for the authenticated home flow.


## 2026/04/22 — mobile/balance-03

Introduces **account balance and statement support in the mobile layer**, including repository abstraction, API integration, caching strategy, and alignment with backend semantics. Also refines documentation to clarify user lifecycle versus account operational status across the system.

### 1. Data Layer — Account Repository

* Added `AccountRepository` interface with:

  * `getBalance`
  * `watchBalance`
  * `getStatement`
  * `getCachedBalance`
* Implemented `AccountRepositoryImpl`:

  * integrates `BalanceApi` and `StatementApi`
  * introduces in-memory cache for balance
  * exposes reactive updates via `StreamController.broadcast`
* Design choice:

  * repository acts as **single source of truth for account state in the mobile layer**
  * balances are treated as **event-driven updates**, not just request/response

### 2. API Integration — Balance and Statement

* Added `BalanceApi`:

  * GET `/accounts/{id}/balance`
  * parses standard API envelope
  * validates HTTP status and maps errors to `AppError`
* Added `StatementApi`:

  * GET `/accounts/{id}/statement`
  * supports query parameters via DTO
  * handles pagination (`cursor`, `cursor_id`, `limit`, `from`, `to`)
* Both APIs:

  * follow existing `RestClient` abstraction
  * enforce explicit parsing and failure handling
  * align with API contract 

### 3. DTO Layer

* Added balance DTO:

  * `BalanceResponseDto`
* Added statement DTOs:

  * `StatementResponseDto`
  * `StatementItemDto`
  * `StatementNextCursorDto`
* Added query DTO:

  * `StatementQueryParamsDto`
* Design observation:

  * DTOs mirror backend payloads closely, preserving **contract fidelity**
  * avoids premature mapping into domain models

### 4. Dependency Injection

* Registered new services:

  * `BalanceApi`
  * `StatementApi`
* Registered new repository:

  * `AccountRepository`
* Maintains existing injection order and modular structure

### 5. Reactive Balance Flow

* Implemented `watchBalance(accountId)`:

  * emits cached value immediately (if available)
  * streams future updates filtered by `accountId`
* Practical implication:

  * UI can subscribe once and react to balance changes
  * avoids redundant polling
* Opinion:

  * this is a **clean and scalable approach**, especially aligned with fintech UX expectations

### 6. Documentation — Status Semantics Clarification

* Updated application and auth documentation to explicitly separate:

  * **User lifecycle status (`users.status`)**
  * **Account operational status (`accounts.status`)**
* Key clarification:

  * user status controls onboarding and eligibility
  * account status controls financial operability
* This aligns:

  * domain model invariants 
  * application behavior 
  * database semantics 

### 7. Database Documentation Enhancements

* Added schema diagram reference
* Documented:

  * `users.status` as lifecycle control
  * `accounts.status` as operational constraint
* Reinforces database as **active consistency boundary**, not passive storage 

### 8. Architectural Guides (AGENT.md)

* Introduced structured guidance files for:

  * root mobile layer
  * core
  * data
  * domain
  * UI
* Defines:

  * responsibilities per layer
  * dependency boundaries
  * coding conventions
* Opinion:

  * this is a **high-leverage addition** — it reduces architectural drift and enforces consistency, especially when scaling or using AI-assisted development

### Conclusion

This commit establishes the **foundation for account visibility in the mobile client**, combining:

* strong API contract alignment
* explicit error handling
* reactive state propagation
* clear architectural boundaries

From a design standpoint, the separation between **user lifecycle and account operability** is particularly important. It eliminates a common ambiguity in financial systems and aligns the entire stack — API, domain, database, and mobile — under a consistent mental model.


## 2026/04/22 — api/list_accounts-01

Refines the **conceptual separation between user lifecycle and account operational status**, updates documentation to reflect this distinction across layers, and introduces the initial mobile-side infrastructure for account balance and statement consumption.

### 1. Documentation — Status Semantics Clarification

* Explicitly separates:

  * `users.status` as **user lifecycle status**
  * `accounts.status` as **account operational status**
* Clarifies responsibilities:

  * user status governs onboarding and account creation eligibility
  * account status governs financial operation execution
* Updates reinforce domain alignment with existing invariants and rules 

### 2. Auth Documentation — Model Refinement

* Reframes the operational model into three distinct concerns:

  * authentication (identity)
  * authorization (role)
  * user lifecycle (onboarding eligibility)
  * account operational state (financial capability)
* Adjusts invariants:

  * removes ambiguity between “active user” and “operational account”
  * introduces clearer responsibility boundaries between application and domain layers
* Aligns documentation with real execution flows already implemented in use cases 

### 3. Database Documentation — Semantic Enrichment

* Adds schema visualization for improved system comprehension
* Documents meaning of:

  * `users.status` as lifecycle control
  * `accounts.status` as operational constraint
* Reinforces database role as consistency boundary and rule enforcer 

### 4. Mobile Layer — Account Data Integration

* Introduces account module in data layer:

  * `AccountRepository` abstraction
  * `AccountRepositoryImpl` with:

    * balance caching
    * reactive balance stream
    * statement retrieval support
* Adds API clients:

  * `BalanceApi` → `/accounts/{id}/balance`
  * `StatementApi` → `/accounts/{id}/statement`
* Implements DTOs:

  * `BalanceResponseDto`
  * `StatementResponseDto` (+ cursor pagination)
  * `StatementQueryParamsDto`
* Ensures adherence to API envelope contract and error handling standard 

### 5. Dependency Injection — Data & Services

* Registers new services and repositories:

  * `BalanceApi`
  * `StatementApi`
  * `AccountRepository`
* Maintains layered dependency flow:

  * UI → Repository → API → RestClient

### 6. Architectural Guidance — Mobile AGENT Files

* Adds structured guidelines for:

  * core, data, domain, and UI layers
* Defines:

  * responsibilities per layer
  * dependency rules
  * error handling conventions
  * extension boundaries
* Establishes a consistent development contract for future iterations

### Conclusion

This commit is primarily **semantic and architectural**, not behavioral.

It improves the system by:

* eliminating ambiguity between user and account states
* aligning documentation with domain and execution reality
* preparing the mobile layer for real API consumption with a clean repository/API split

From a design perspective, this is a critical step toward **maintaining conceptual integrity across API, database, and client**, which is essential in financial systems where subtle semantic confusion often leads to incorrect business rules.


## 2026/04/22 — mobile/balance-01

Introduces the first mobile integration for **account balance and statement retrieval**, while also updating API documentation to clarify status semantics and database representation. The mobile layer now has the minimum data flow needed to consume account read operations from the backend. 

### 1. Mobile architecture guidance

* Added agent guides for:

  * `mobile/`
  * `mobile/lib/core`
  * `mobile/lib/data`
  * `mobile/lib/domain`
  * `mobile/lib/uis`
* Documents the current Flutter layering, dependency direction, DI rules, routing conventions, async handling, and boundaries between UI, repositories, and API services
* Helps preserve the intended structure as the mobile app grows

### 2. Data layer wiring

* Extended dependency registration in `mobile/lib/data/data.dart`
* Registered new mobile account dependencies:

  * `AccountRepository`
  * `AccountRepositoryImpl`
  * `BalanceApi`
  * `StatementApi`
* Keeps account read features aligned with the existing injector-based composition model

### 3. Account repository for mobile

* Added `AccountRepository` contract with support for:

  * cached balance access
  * balance streaming
  * remote balance retrieval
  * statement retrieval with query parameters
* Added `AccountRepositoryImpl` implementing:

  * in-memory balance cache
  * broadcast stream for balance updates
  * API delegation for balance and statement endpoints
* This is a solid first step because it keeps API access out of the UI and establishes the repository as the app-facing entry point

### 4. New account APIs

* Added `BalanceApi`
* Added `StatementApi`
* Both APIs:

  * use the existing `RestClient`
  * parse the standard API envelope
  * convert transport and parsing failures into `AppError`
  * return `AsyncResult<T>`
* This keeps the transport layer consistent with the existing auth implementation and avoids leaking raw HTTP concerns upward

### 5. DTOs for account reading

* Added `BalanceResponseDto`
* Added `StatementQueryParamsDto`
* Added `StatementResponseDto`
* Added statement support types:

  * `StatementItemDto`
  * `StatementNextCursorDto`
* The statement DTO structure matches the backend contract for paginated ledger/history reads, including cursor-based pagination

### 6. Service registration

* Updated `mobile/lib/data/services/services.dart`
* Registered:

  * `BalanceApi`
  * `StatementApi`
* Keeps API construction centralized and consistent with the existing `AuthApi` pattern

### 7. API documentation refinements

* Clarified the distinction between:

  * `users.status` as **user lifecycle status**
  * `accounts.status` as **account operational status**
* Updated auth/application/database documentation to reflect:

  * onboarding/account-opening eligibility tied to user lifecycle
  * financial operation eligibility tied to account status
* This is an important clarification because it removes conceptual ambiguity that could easily leak into both backend rules and mobile assumptions

### 8. Database documentation updates

* Added a schema image reference in `api/docs/09-database.md`
* Added notes documenting:

  * `users.status` semantics
  * `accounts.status` semantics
* Added the new schema image asset under `api/docs/images/databese.png`

### Conclusion

This commit establishes the **mobile read path for account balance and statement**, with proper layering, dependency injection, DTOs, transport handling, and repository orchestration. At the same time, it improves backend documentation by making status responsibilities explicit, which is particularly valuable for keeping mobile and API behavior aligned.


## 2026/04/22 — refactor/api-application-structure-04

Refines the **API and application structure** to better align documentation, use case flows, and architectural responsibilities with the current implementation, particularly around account lifecycle, admin operations, and branch resolution strategy. 

### 1. API Surface Adjustments

* Added new endpoint:

  * `GET /accounts/{id}/balance`
* Introduced explicit access control note:

  * `/admin/users/{id}/approve` requires authenticated user with `admin` role
* Improves clarity of exposed capabilities and security model at the API boundary

### 2. Use Case Flow Corrections

* Refined **Approve User flow**:

  * added explicit validation of associated customer existence
  * introduced account number generation and branch resolution as distinct steps
  * clarified transactional sequence and responsibilities
* Updated **Create Account flow**:

  * replaced implicit branch generation with explicit **branch policy resolution**
* These changes align documentation with actual application orchestration and domain expectations 

### 3. Implementation Documentation Alignment

* Updated approval flow to reflect:

  * dependency on `CustomerRepository`
  * validation of `customer_id` before account creation
  * explicit ownership of branch resolution by application layer
* Adjusted create account flow:

  * branch is no longer a fixed value, but derived from a policy
* Expanded test coverage section:

  * includes admin handler tests
* Clarified implementation note:

  * branch policy now belongs to `account/application/account` instead of being hardcoded
* This resolves inconsistencies between documentation and actual system behavior 

### 4. REST Contract Simplification

* Simplified response of **Approve User endpoint**:

  * removed `email` and `customer_id` from response payload
* Adjusted error mapping:

  * added `CUSTOMER_NOT_FOUND`
  * changed `USER_ALREADY_ACTIVE` from 400 → 409
* Improves contract minimalism and correctness, reducing redundant data exposure

### 5. Architectural Evolution

* Introduced **admin module** as first-class component:

  * added to module list and architectural description
* Updated architecture documentation:

  * reflects expanded module set (`account`, `admin`, `auth`, `customer`)
  * includes new use case: `approve user`
  * registers new route `/admin/users/{id}/approve`
* Reinforces modular monolith structure and clearer separation of responsibilities 

### 6. Conceptual Improvement — Branch Policy

* Replaced hardcoded branch generation with **explicit branch policy abstraction**
* Defined ownership at application layer instead of infrastructure or domain
* This is a subtle but important design improvement:

  * removes hidden assumptions
  * enables future extensibility without refactoring core flows

### Conclusion

This refactor is primarily **structural and semantic**, not functional. It eliminates inconsistencies between documentation, API contracts, and implementation while introducing clearer boundaries and responsibilities.

The most relevant improvement is the **formalization of application-level policies (such as branch resolution) and the elevation of admin operations into the architecture**, resulting in a more coherent and evolvable system design.


## 2026/04/22 — refactor/api-application-structure-03

Refactors the application structure to improve **module separation, dependency clarity, and cross-cutting concerns reuse**, with emphasis on isolating admin responsibilities, centralizing authentication context, and formalizing branch generation as a policy.

### 1. Admin Module Extraction

* Moved `ApproveUserUseCase` from auth module to new `admin` module
* Introduced:

  * `internal/admin/application`
  * `internal/admin/delivery`
* Removed approval responsibility from auth handler
* Created dedicated `adminHandler` and route wiring
* Result:

  * clear separation between **authentication concerns** and **administrative workflows**
  * aligns with use case flow definition for user lifecycle 

### 2. Branch Policy Abstraction

* Introduced `BranchPolicy` interface in account application layer
* Implemented `DefaultBranchPolicy` (currently fixed as "0001")
* Injected policy into:

  * `CreateAccount`
  * `ApproveUserUseCase`
* Removed hardcoded branch generation function
* Added explicit validation for missing policy
* Architectural impact:

  * removes hidden business rule
  * prepares system for future variability (multi-branch, configuration-driven)
* Aligns with application layer responsibility of operational rules 

### 3. Auth Context Centralization

* Created shared package: `internal/shared/authctx`
* Migrated all context-related helpers:

  * `WithAuthenticatedUser`
  * `GetAuthenticatedUser`
  * `RequireAuthenticatedUser`
* Updated:

  * middleware
  * handlers
  * tests
* Eliminated duplication across modules
* Improves consistency of authentication flow handling 

### 4. Auth Handler Simplification

* Removed approval logic from auth handler
* Reduced constructor dependencies
* Eliminated approve-related DTOs and tests
* Auth module now strictly focused on:

  * register
  * login
  * token refresh
  * current user
* Reinforces boundary between **identity management** and **domain operations**

### 5. Customer Repository Contract Unification

* Removed legacy `account/domain.CustomerRepository`
* Consolidated on `customer/domain.CustomerRepository`
* Standardized interface:

  * `Create`
  * `Exists`
  * `GetByID`
* Updated all consumers and mocks
* Improves consistency across modules and reduces duplication

### 6. Delivery Layer Adjustments

* Introduced dedicated admin handler with:

  * role validation (admin-only)
  * UUID parsing
  * error mapping via shared error system
* Updated routing:

  * `/admin/users/{id}/approve` now handled by admin module
* Updated account and customer handlers to use shared auth context

### 7. Error Registration Extension

* Added admin-specific error mappings via `admin/application/errors_registry.go`
* Integrated into bootstrap
* Ensures consistency with global error contract 

### 8. Main Wiring Refinement

* Updated `main.go`:

  * introduced `branchPolicy` as shared dependency
  * separated auth and admin handlers
  * aligned use case construction with new module boundaries
* Improves clarity of runtime composition 

### 9. Test Suite Updates

* Updated all tests to:

  * inject `BranchPolicy`
  * use shared auth context
  * reflect new module boundaries
* Added coverage for:

  * missing branch policy scenarios
  * admin handler behavior
* Maintains strong validation of transactional and lifecycle flows

### 10. Architectural Impact

This refactor reinforces key architectural principles:

* **Separation of concerns**

  * auth vs admin vs account responsibilities clearly isolated
* **Explicit dependencies**

  * no hidden rules (branch now injected)
* **Shared cross-cutting infrastructure**

  * authentication context unified
* **Consistency across layers**

  * repository contracts and error handling standardized

From a design perspective, this is a meaningful evolution toward a more **cohesive modular monolith**, improving maintainability without introducing unnecessary abstraction.


## 2026/04/22 — refactor/api-application-structure-02

Refactors the **application layer structure of the account module**, introducing a clear separation by subdomains (account, transaction, statement) and aligning imports, use cases, and delivery contracts with this new organization. This change improves modularity, reduces implicit coupling, and better reflects the system’s use case boundaries as described in the architecture and application model. 

### 1. Application Layer Restructuring

* Split `internal/account/application` into three focused sub-packages:

  * `account` → account lifecycle (create account, get balance)
  * `transaction` → financial operations (deposit, withdraw, transfer)
  * `statement` → read models (account statement)
* Moved all use cases to their respective domains without changing behavior
* Renamed imports across the project to reflect the new structure

This change reinforces the idea that the application layer is the orchestration boundary of use cases, not a generic container. 

### 2. Main Wiring Adjustments

* Updated `cmd/api/main.go`:

  * `Deposit`, `Withdraw`, `Transfer` now originate from `transaction`
  * `GetStatement` now originates from `statement`
  * `CreateAccount` and `GetAccountBalance` remain under `account`
* Maintains explicit composition while improving semantic clarity of dependencies

### 3. Access Policy Consolidation per Subdomain

* Introduced `access_policy.go` in:

  * `application/account`
  * `application/transaction`
  * `application/statement`

* Provides:

  * `CanAccessCustomer`
  * `CanAccessAccount`

* Enforces authorization at the application layer, consistent with the system’s model where ownership is derived from JWT context and not from client input 

* Although duplicated across packages, this decision keeps each subdomain **self-contained**, avoiding cross-package coupling at the cost of controlled redundancy

### 4. Delivery Layer Alignment

* Updated handlers to consume the new package structure:

  * `accountapp`, `transactionapp`, `statementapp`
* Adjusted all DTO mappings to match new package locations
* Preserved handler responsibilities:

  * request parsing
  * use case invocation
  * error mapping to HTTP responses 

### 5. Interface and Contract Updates

* Updated handler interfaces to depend on:

  * `accountapp` inputs/outputs
  * `transactionapp` inputs/outputs
  * `statementapp` inputs/outputs
* Maintains strict separation between use cases and delivery contracts

### 6. Test Suite Adaptation

* Refactored all unit and integration tests to match new package boundaries:

  * mocks updated to use new DTO types
  * transfer, deposit, withdraw tests now reference `transaction` package
  * statement tests now reference `statement` package
* Added shared test helpers (`testCustomerUser`, `testAdminUser`) in each subdomain
* No behavioral changes, only structural alignment

### 7. Cross-Module Impact

* Updated auth module to use:

  * `account/application/account` instead of the previous flat application package
* Updated integration tests to reference `transaction` where applicable

### 8. Architectural Implications

This refactor is not cosmetic. It aligns the codebase with a more precise interpretation of:

* **use case boundaries** (transaction vs. account vs. read models)
* **application layer responsibilities** (orchestration and authorization)
* **modular monolith principles**, reducing implicit dependencies 

From a design perspective, this is a significant improvement. The previous flat structure was already functional, but it blurred responsibilities. The new structure introduces:

* clearer mental model
* improved scalability for future features
* safer evolution of independent use case groups

### Conclusion

This commit restructures the application layer into **cohesive subdomains**, improving clarity, modularity, and alignment with the system’s architectural principles. While introducing some duplication (notably in access policies), the trade-off favors **low coupling and explicit boundaries**, which is a sound decision for a financial system that prioritizes correctness and maintainability.


## 2026/04/22 — refactor/api-application-structure-01

Refactors the **application layer structure for the account module**, introducing a clearer separation of concerns between account management, transactions, and statement operations. This change aligns the codebase with a more explicit use-case organization and improves modularity across layers.

### 1. Application Layer — Structural Decomposition

* Split `account/application` into three focused submodules:

  * `account` → account lifecycle (create, balance)
  * `transaction` → deposit, withdraw, transfer
  * `statement` → statement retrieval
* Moved existing use cases into their respective domains without altering behavior
* Updated all imports to reflect the new structure
* This refactor reinforces the role of the application layer as an orchestration boundary, as defined in the architecture 

### 2. Access Control — Explicit Policy Duplication per Module

* Introduced `access_policy.go` in:

  * `account`
  * `transaction`
  * `statement`
* Each module now contains:

  * `CanAccessCustomer`
  * `CanAccessAccount`
* These policies enforce ownership rules based on authenticated user context
* Although duplicated, this design favors **locality and independence of use cases**, avoiding cross-module coupling

### 3. Authorization Test Helpers

* Added `auth_test.go` in each application submodule
* Introduced helper builders:

  * `testCustomerUser`
  * `testAdminUser`
* Standardizes test setup for authorization scenarios across modules

### 4. Main Composition Root (cmd/api/main.go)

* Updated use case instantiation to reflect new modular structure:

  * transaction use cases (`Deposit`, `Withdraw`, `Transfer`)
  * statement use case (`GetStatement`)
  * account use cases remain in `account`
* Improves clarity of responsibilities during application wiring
* Makes the composition root more expressive and aligned with business capabilities

### 5. Delivery Layer — Handler Adaptation

* Updated handler imports to use new module paths:

  * `accountapp`
  * `transactionapp`
  * `statementapp`
* Adjusted all DTO mappings to reference the correct module
* Ensures that HTTP layer remains a thin adapter over application use cases, consistent with the architecture 

### 6. Delivery Tests — Type Alignment

* Updated all mocks and test cases to use new input/output types per module
* Adjusted:

  * `CreateAccountInput`
  * `DepositInput`, `WithdrawInput`, `TransferInput`
  * `GetStatementInput`
  * `AccountBalance`, `Statement`, `TransferResult`
* Maintains full test coverage with the new structure

### 7. Integration Tests

* Updated integration tests to reference:

  * `transaction` module for deposit use case
* Confirms that runtime wiring remains consistent after refactor

### 8. Auth Module Integration

* Updated `approve_user` use case to import account use cases from the new `account` module
* Adjusted authorization integration tests to use `transaction` module where applicable
* Keeps cross-module dependencies aligned with new boundaries

### 9. Architectural Impact

* This refactor makes the application layer explicitly reflect **use case categories**, rather than grouping everything under a single package
* Improves:

  * readability
  * maintainability
  * scalability of the module
* Reinforces the separation between:

  * account state management
  * financial operations (transactions)
  * historical data access (statements)

This structure better matches the defined use case flows and system behavior  and prepares the codebase for future evolution without introducing unnecessary complexity.

### Conclusion

This change significantly improves the **semantic organization of the application layer**, making the system easier to reason about and extend. While it introduces some controlled duplication (access policies), it results in a cleaner, more modular design that aligns closely with the domain and use case boundaries.


## 2026/04/22 — api/balane-01

Adds support for **account balance retrieval** as a first-class use case, exposing a new protected endpoint and ensuring consistency with the existing authorization, domain, and repository contracts.

### 1. Application Layer — GetAccountBalance Use Case

* Introduced `GetAccountBalance` use case

  * Input includes authenticated user and account ID
  * Validates:

    * non-nil UUID
    * account existence
    * ownership via `CanAccessAccount`
* Returns a minimal projection:

  * `AccountID`
  * `Balance`
* Aligns with the existing principle that balance is a **snapshot derived from the ledger** 
* Keeps the operation read-only and free from transactional overhead, which is appropriate given that no state mutation occurs 

### 2. Domain Contract Clarification

* Strengthened `AccountRepository.GetByID` contract:

  * must never return `(nil, nil)`
  * explicitly returns `ErrAccountNotFound`
* This removes ambiguity and prevents invalid states propagating into the application layer
* Consistent with domain invariant that **account must exist for any operation** 

### 3. Infrastructure Layer — PostgreSQL

* Updated repository implementations (`Repository` and `txRepository`):

  * enforce `ErrAccountNotFound` when no row is returned
* Eliminates implicit null handling and aligns persistence behavior with domain expectations
* Improves correctness and reduces defensive checks in upper layers

### 4. Delivery Layer — HTTP Endpoint

* Added new endpoint:

  * `GET /accounts/{id}/balance`
* Implemented handler `GetBalance`:

  * requires authenticated user (JWT)
  * rejects any query parameters (`400 INVALID_DATA`)
  * validates UUID parsing
  * delegates to application use case
* Response structure:

  * `{ account_id, balance }`
  * balance returned in cents, consistent with system-wide monetary representation 
* Error mapping aligned with global response standard 

### 5. Handler Wiring and Routing

* Extended handler constructor to include `balance` use case
* Updated dependency injection in `main.go`
* Registered new route:

  * `GET /accounts/{id}/balance`
* Maintains consistency with existing route protection via JWT middleware 

### 6. Delivery DTOs

* Added `AccountBalanceData` for HTTP response mapping
* Keeps delivery layer isolated from application structs, preserving layer boundaries

### 7. Test Coverage

#### 7.1 Application Tests

* Covered scenarios:

  * invalid account ID
  * account not found
  * forbidden access (ownership violation)
  * successful retrieval
* Ensures no repository call on invalid input (early validation pattern)

#### 7.2 Delivery Tests

* Added handler tests for:

  * query params rejection
  * unauthorized access
  * invalid UUID
  * account not found
  * forbidden access
  * success response validation

#### 7.3 Integration Adjustments

* Updated handler constructor usage across tests to include new dependency
* Maintains compatibility with existing integration setup

### 8. Documentation Updates

* Extended REST API documentation with:

  * endpoint definition
  * request constraints
  * response format
  * error scenarios
* Adjusted section numbering to accommodate the new endpoint
* Reinforces API contract consistency and discoverability 

### 9. Customer Repository Fix

* Corrected `GetByID` behavior:

  * now returns `ErrNotFound` instead of `(nil, nil)`
* Aligns customer module with the same repository contract expectations applied to accounts

### Conclusion

This change introduces a **clean and well-isolated read use case** that fits naturally into the existing architecture.

From a design perspective, the implementation is correct and disciplined:

* no leakage of persistence concerns into application logic
* strict enforcement of ownership and domain invariants
* consistent error semantics across layers

Although technically simple, this endpoint is important because it formalizes **balance as an explicit query operation**, reinforcing the separation between:

* ledger (source of truth)
* account snapshot (read model)

This is a necessary step toward a more complete and coherent financial API surface.


## 2026/04/22 — api/migrate-update-01

Refactors the database migration strategy by consolidating all historical migrations into a single baseline schema and aligning the entire codebase with the `transactions` table as the authoritative ledger.

### 1. Migration Strategy Simplification

* Removed all legacy incremental migration files (000000 → 000010 series)
* Introduced a new baseline migration:

  * `000001_init_schema.up.sql`
  * `000001_init_schema.down.sql`
* The new migration represents the **full current database state**, eliminating unnecessary historical transitions
* Establishes a clean starting point for future migrations with reduced cognitive overhead and improved reproducibility

### 2. Ledger Standardization

* Definitively replaces `account_transactions` with `transactions` across the entire system
* Reinforces the architectural decision that:

  * the ledger is **append-only**
  * the ledger is the **single source of truth for financial state**
* Aligns with the domain model where all balance changes must be backed by transactions 

### 3. Documentation Alignment

* Updated all documentation references:

  * README
  * domain model
  * use case flows
  * consistency and concurrency strategy
  * implementation documentation
  * database documentation
* Ensures consistency between:

  * conceptual model
  * implementation
  * persistence layer
* All flows now explicitly reference `transactions` as the ledger, including deposit, withdraw, and transfer operations 

### 4. Infrastructure Layer Updates

* Updated PostgreSQL repository queries:

  * `INSERT`, `SELECT`, and query operations now target `transactions`
* Ensures full compatibility with the new schema without introducing behavioral changes

### 5. Test Suite Adjustments

* Updated integration tests and schema setup:

  * replaced `account_transactions` with `transactions`
  * adjusted index names and constraints
* Maintained backward compatibility logic where necessary:

  * conditional column creation for evolving schemas
* Updated cleanup routines and assertions to reflect the new ledger table

### 6. Schema Evolution Improvements

* New baseline schema includes:

  * enums (`account_status`, `transaction_type`)
  * sequences (`account_number_seq`)
  * immutable ledger enforcement via trigger (`prevent_transactions_mutation`)
  * transfer integrity constraints (`reference_id + type`)
  * idempotency guarantees via unique index
* Consolidates previously scattered concerns into a **cohesive and explicit schema definition**

### 7. Removal of Redundant Artifacts

* Deleted legacy `schema.sql` dump file
* Eliminated outdated migration logic related to:

  * ledger consolidation
  * user/session evolution steps
  * incremental schema adjustments
* Reduces ambiguity about which structure represents the current truth

### Conclusion

This change is a **structural simplification and consistency correction** of the persistence layer.

It eliminates migration noise, enforces a single authoritative ledger (`transactions`), and aligns documentation, code, and database schema under a unified model.

From an architectural standpoint, this is a decisive improvement: it reduces ambiguity, strengthens the ledger model, and makes the system significantly easier to reason about and evolve.


## 2026/04/19 — docs/update-02

Updates documentation, developer workflow, and project structure to reflect the current state of the system, with emphasis on bootstrap automation, ledger-based modeling, and clearer onboarding for both API and mobile applications.

### 1. Makefile — Developer Workflow and Environment Management

* Introduced a complete **Docker orchestration layer**:

  * `docker-up`, `docker-down`, `docker-logs`, `docker-clean`, `docker-check`
* Added **environment lifecycle commands**:

  * `setup`: full bootstrap (Docker + DB + migrations)
  * `run`: full system startup including API
  * `reset`: deterministic full reset (containers + database + migrations)
  * `db-reset`, `db-wait`
* Added semantic aliases:

  * `bootstrap` → setup
  * `dev` → run
* Removed duplicated Docker section and consolidated structure

This significantly improves reproducibility and reduces manual setup errors.

### 2. Root README — Conceptual Clarification and Onboarding

* Reframed system model:

  * emphasizes **ledger as source of truth**
  * explicitly states `account_transactions` as authoritative
* Removed ambiguity around legacy `transactions` table
* Expanded API capabilities:

  * refresh token support
  * admin approval flow
  * customer self-access (`/customers/me`)
* Updated authentication model:

  * AppToken for onboarding
  * JWT for all other routes
* Replaced manual setup steps with:

  * `make run` as the primary entrypoint
* Added direct references to:

  * API getting started guide
  * Mobile getting started guide

This aligns the documentation with the actual implementation and reduces cognitive friction during onboarding.

### 3. API README — Documentation Reorganization

* Replaced outdated relative paths with **local docs structure**
* Introduced a structured documentation index:

  * getting started
  * architecture
  * domain, application, and persistence models
  * database and infrastructure
* Standardized naming conventions across documents

Improves discoverability and reinforces the API as a well-documented subsystem.

### 4. API Bootstrap — Configuration Refactor

* Introduced `Config` struct and `LoadConfig()` function:

  * centralizes environment variable validation
  * enforces fail-fast behavior
* Replaced direct environment access with structured configuration usage
* Updated dependency wiring to use config values

This is a subtle but important step toward better configuration management and testability.

### 5. API Documentation — Getting Started

* Added new guide: `api/docs/00-getting_started.md`
* Defines:

  * environment variables (`APP_TOKEN`, `JWT_SECRET`)
  * Docker initialization
  * bootstrap via `make setup`
  * execution via `make run`
  * reset strategy
* Reinforces database-first approach and deterministic environment setup

Complements the implementation described in  by making the runtime flow explicit.

### 6. Mobile Documentation — Getting Started

* Added new guide: `mobile/docs/00-getting_started.md`
* Defines:

  * environment configuration via `dart-define`
  * API dependency requirement
  * execution with environment files (`dev.env`, `staging.env`, `prod.env`)
* Documents networking constraints (emulator vs device)
* Integrates mobile startup with backend lifecycle (`make run`)

This removes ambiguity in mobile setup and enforces alignment with backend availability.

### 7. Mobile README — Documentation Alignment

* Updated references to:

  * API getting started
  * mobile getting started
  * local architecture docs
* Removed outdated paths and centralized documentation access

### 8. Documentation Cleanup and Standardization

* Removed obsolete file:

  * `api/docs/changes.md`
* Simplified formatting in `errors.md`:

  * removed visual noise and emoji markers
  * improved readability and long-term maintainability
* Standardized document naming and structure across the repository

### 9. Conceptual Consolidation — Ledger Model

* Reinforced across documentation:

  * ledger entries are append-only
  * balance is derived, not authoritative
* Aligns documentation with:

  * domain model invariants 
  * consistency strategy 
  * use case execution model 

### Conclusion

This commit transforms the project documentation from a partial reference into a **cohesive, executable guide** for the system. It also aligns the conceptual model, developer workflow, and runtime behavior, resulting in a more robust and maintainable foundation for further evolution.


## 2026/04/18 — docs/update-01

Refactors and consolidates the entire project documentation to align with the current implementation, architectural decisions, and system positioning. The update establishes a clear narrative centered on **transactional consistency**, **ledger-driven design**, and **explicit system guarantees**, while significantly improving documentation discoverability and structure.

### 1. Root Documentation — README.md

* Rewrote the project introduction to reflect the system’s core premise:

  * balance as a consequence of transactions
  * transaction as the central element of the model
* Introduced a formal **System Scope** section:

  * defined goals, boundaries, and responsibilities
  * clarified in-scope vs out-of-scope features
* Added explicit **system guarantees**:

  * financial integrity, atomicity, traceability, consistency, synchronous model, single source of truth
* Improved API and Mobile descriptions to better reflect actual capabilities
* Reorganized documentation references into structured sections:

  * API documentation
  * Mobile documentation
* Elevated the README from a simple entry point to a **high-level architectural contract**

### 2. Documentation Consolidation

* Centralized system definition previously described in:

  * `api/docs/objetivos.md`
* Replaced duplicated content with a reference to the root README:

  * avoids divergence between documents
  * enforces a single source of truth for system scope

### 3. Full Documentation Update (api/docs)

* All documents under `api/docs` were reviewed and updated to reflect the current state of the system
* The documentation now consistently describes:

  * domain model centered on transactions and invariants 
  * data model aligned with financial consistency and traceability 
  * strong consistency and concurrency strategy using ACID transactions and pessimistic locking 
  * use case flows with explicit transactional behavior and validation order 
  * REST contract with standardized response envelope and error mapping 
  * error handling model with stable codes and predictable structure 
  * authentication and authorization model with staged token strategy 
  * implementation details reflecting actual runtime behavior and layering 

### 4. Structural Improvements

* Eliminated conceptual duplication across documents
* Strengthened separation between:

  * conceptual documentation (domain, flows, guarantees)
  * implementation documentation (runtime behavior, wiring, persistence)
* Established a clearer documentation hierarchy:

  * README as entry point and system definition
  * api/docs as deep technical reference

### 5. Documentation Positioning

* The documentation now functions as:

  * a **technical specification of the system**
  * a **learning artifact for architectural decisions**
  * a **foundation for future evolution (e.g., Zero Trust, distributed systems)**

### Conclusion

This commit transforms the documentation from a fragmented set of notes into a **cohesive, authoritative, and implementation-aligned body of knowledge**. The system is now clearly defined not only by its code, but by explicit guarantees, constraints, and architectural intent, which significantly increases its value as both a technical asset and a reference model.


## 2026/04/18 — mobile/login-02

Implements a **refined authentication UI layer** with standardized input components, improved routing transitions, and consistent theming across login and registration flows. Also aligns the mobile layer with the backend contract and documentation updates 

### 1. Linting and Import Consistency

* Enabled `prefer_relative_imports` in `analysis_options.yaml`
* Refactored imports in `dio_rest_client.dart` to use relative paths
* Improves modular isolation and reduces dependency coupling across layers

### 2. Routing Layer — Custom Transitions

* Introduced `AppCustomTransactionPage`

  * Combines `FadeTransition` with `ScaleTransition`
  * Uses `easeOutCubic` curve for smoother perception
* Migrated auth routes (`login`, `register`) from `builder` to `pageBuilder`
* Centralizes navigation behavior at routing level rather than UI layer

Opinion: this is a correct architectural move. Transition concerns belong to routing, not to pages. It avoids UI fragmentation and implicit navigation behavior.

### 3. UI Abstraction — Form Components

* Introduced `BasicTextFormField`

  * Encapsulates:

    * decoration
    * border styling
    * radius standardization
    * icon handling
  * Reduces duplication across forms
* Applied to:

  * `LoginPage`
  * `RegisterPage`

Impact:

* Eliminates repeated `InputDecoration`
* Enforces visual consistency
* Simplifies future global UI changes

### 4. Input Formatting — CPF

* Added `CpfInputFormatter`

  * Normalizes input to digits only
  * Applies mask: `000.000.000-00`
* Integrated into registration form with:

  * `FilteringTextInputFormatter.digitsOnly`
  * custom formatter

Opinion: this is a critical UX improvement. It reduces invalid payloads before they reach the API and aligns well with backend expectations for CPF validation.

### 5. Login Page Refinements

* Replaced raw `TextFormField` with `BasicTextFormField`
* Added:

  * `textInputAction` flow (next/done)
  * explicit hints
  * submit via keyboard (`onFieldSubmitted`)
* Simplified `ValueListenableBuilder` usage for password visibility

### 6. Register Page Refinements

* Fully migrated all inputs to `BasicTextFormField`
* Improved UX flow:

  * sequential navigation via keyboard actions
  * consistent hints and labels
* Integrated CPF formatting and validation pipeline

### 7. Theme Standardization

* Centralized theme adjustments in `AppWidget`:

  * `AppBarTheme` now uses `colorScheme` explicitly
  * Introduced `InputDecorationTheme`:

    * filled inputs
    * semi-transparent background
    * rounded borders (24px)
* Aligns visual identity across all form fields without per-widget duplication

Opinion: this is a strong step toward a design system. The combination of theme + base component is significantly more maintainable than scattered styling.

### 8. Documentation Update (API Layer)

* Updated entire `api/docs` structure (not included in diff due to size)
* Ensures alignment between:

  * mobile DTOs
  * authentication flow (login/register/refresh)
  * response envelope and error handling patterns

This is particularly relevant given the API contract:

* standardized `data/error` envelope
* authentication flow using AppToken + JWT
* consistent error codes and payloads 

### Conclusion

This commit represents a **UI architecture consolidation for authentication flows**, focusing on:

* component reuse
* visual consistency
* controlled navigation behavior
* input validation at the edge

From a design perspective, the most important gain is the **emergence of a coherent UI foundation**. The introduction of a base input component combined with theme-level control significantly reduces long-term maintenance cost and prevents divergence across screens.


## 2026/04/18 — mobile/login-01

Implements the **initial mobile authentication integration aligned with the backend contract**, introducing structured logging, AppToken support, stricter HTTP handling, and internal refactoring for environment configuration. Also includes database migration adjustments and a full documentation update across the API layer. 

### 1. Environment Configuration Refactor

* Moved configuration files from `core/config` to `core/resources`:

  * `app_env.dart`
  * `storage_keys.dart`
* Refactored `AppEnv`:

  * introduced `APP_ACCESS_TOKEN` support (`appToken`)
  * encapsulated timeouts and mode as private constants with getters
  * added `isProd` helper
* Improved validation for `BASE_URL`
* Overall effect:

  * clearer separation between configuration and infrastructure
  * safer access to environment variables

### 2. HTTP Layer Enhancements (Dio + RestClient)

* Added structured logging via `ConsoleLog`:

  * request logging for `POST` and `PUT`
  * centralized error logging with stack trace
* Improved error handling in `DioRestClient`:

  * ensures all exceptions are mapped consistently to `Result.failure`
* This aligns the client with the backend response contract:

  * strict envelope parsing
  * predictable failure paths 

### 3. Auth Interceptor Improvements

* Replaced `dart:developer.log` with `ConsoleLog`
* Added detailed error diagnostics:

  * status code
  * error type
  * message
  * full stack trace
* Improved observability for:

  * token refresh failures
  * unexpected network errors
* Maintains current refresh flow while making failure states explicit

### 4. Secure Storage Logging

* Integrated `ConsoleLog` into `FlutterSecureStorageLocalStorage`
* All operations now include structured error logging:

  * read
  * write
  * delete
  * deleteAll
* Converts silent failures into traceable events without changing behavior

### 5. Auth API Integration (Login/Register)

* Added required header:

  * `X-App-Token: AppEnv.appToken`
* Enforced HTTP validation before parsing:

  * rejects non-2xx responses explicitly
* Refactored response handling:

  * now uses `RestClientResponse` instead of raw maps
  * validates `statusCode` before envelope parsing
* Improved error semantics:

  * HTTP errors mapped explicitly
  * parsing errors logged with stack trace
* Aligns mobile client with backend authentication model:

  * AppToken required for onboarding endpoints
  * JWT used after login 

### 6. Logging Infrastructure Introduction

* Introduced `ConsoleLog` as a reusable logging utility:

  * supports `error`, `warn`, `info`, and raw `log`
  * context-aware logging (`Class.method`)
  * debug-only execution via `kDebugMode`
* Establishes a consistent logging standard across:

  * HTTP layer
  * interceptors
  * storage
  * APIs

### 7. Dependency Wiring Adjustments

* Updated imports across core services to reflect new structure
* Ensured `DioFactory` and `CoreServices` use the new `AppEnv` location
* Maintains existing DI structure without introducing new abstractions

### 8. Database and Backend Alignment

* Added migrations:

  * `customer.email` becomes nullable
  * full removal of `customer.email`
* Introduced `schema.sql` dump for full database snapshot
* Updated Makefile:

  * simplified schema export command (`dbschema`)
* These changes reflect the architectural shift:

  * email responsibility moved from Customer to Auth layer

### 9. Documentation Update (API)

* Entire `api/docs` directory updated and normalized:

  * heading hierarchy fixes
  * consistency improvements across endpoints
* Documentation now reflects:

  * AppToken requirement for `/auth/*`
  * JWT-based access control
  * standardized response envelope
* This ensures the mobile client is aligned with the real API contract 

### Conclusion

This commit establishes a **cohesive foundation for mobile authentication**, ensuring strict alignment with backend contracts while improving observability and error handling.

From an architectural perspective, the most relevant gains are:

* explicit separation of configuration concerns
* deterministic HTTP behavior
* consistent logging strategy
* formalization of AppToken-based onboarding

The result is a **much more predictable and debuggable client**, which is essential before advancing to higher-level flows such as session management and UI integration.


## 2026/04/17 — api/customer-email-removal-01

Refactors the **customer domain boundary by removing email from the Customer aggregate**, repositioning it as a responsibility of the authentication layer, and aligning persistence, application flows, and delivery contracts accordingly. Also introduces a utility for database schema extraction and updates all related tests and documentation. 

### 1. Domain Layer — Customer Simplification

* Removed `email` from `Customer` entity
* Updated `NewCustomer` factory:

  * now validates only `name` and `cpf`
* Removed domain errors:

  * `ErrEmailRequired`
  * `ErrEmailAlreadyExists`
* Reinforces a cleaner domain model where:

  * Customer represents **identity (CPF-based)**
  * Email belongs to **authentication/user context**

This is a **structural correction**, not merely a field removal, improving separation of concerns.

### 2. Repository Contract Redesign

* Updated `CustomerRepository.GetByID` signature:

  * now returns `(customer, email, error)`
* Email is retrieved via join with `users` table at infrastructure level
* Makes explicit that:

  * email is **not part of the aggregate**
  * but still required for **read models / API responses**

### 3. Infrastructure Layer — PostgreSQL

* Updated `customers` persistence:

  * removed `email` from INSERT operations
* Adjusted constraint handling:

  * removed email unique constraint mapping
* Refactored `GetByID` query:

  * now performs:

    ```sql
    JOIN users u ON u.customer_id = c.id
    ```
  * returns email as a separate value
* Removed unused helper:

  * `nullableStringValue`

This aligns persistence with the **new aggregate boundary** while preserving API requirements.

### 4. Application Layer Adjustments

#### 4.1 Create Customer

* Input no longer includes email
* Delegates email responsibility entirely to auth flow

#### 4.2 Get Customer Me

* Updated return signature to include email separately
* Propagates `(customer, email)` across layers
* Ensures:

  * domain purity
  * complete response composition

### 5. Auth Module Alignment

* `RegisterUser` no longer injects email into Customer entity
* Maintains email strictly within `User` context
* Strengthens consistency between:

  * user identity (auth)
  * customer identity (domain)

### 6. Delivery Layer — HTTP Contract

#### 6.1 Create Customer Endpoint

* Introduced explicit request DTO:

  * `createCustomerRequest`
* Email remains in request payload but:

  * is not persisted in Customer
  * is returned as part of response
* Adjusted validation:

  * removed email-related domain errors

#### 6.2 Get /customers/me

* Now composes response using:

  * `customer` (domain)
  * `email` (auth layer)
* Preserves API contract expected by clients 

This is a **read model composition at delivery level**, which is the correct architectural placement.

### 7. Test Suite Updates

* Updated all mocks and signatures to reflect new repository contract
* Adjusted integration tests:

  * removed email from `customers` schema
  * removed email from seed data
* Extended assertions:

  * validate email returned separately
* Ensures full coverage of new behavior across:

  * application
  * delivery
  * infrastructure

### 8. Makefile Enhancement

* Added `database-schema` target:

  * exports current DB schema using `pg_dump --schema-only`
* Improves:

  * observability
  * documentation workflow
  * version tracking of schema

### 9. Documentation Update

* Entire `api/docs` directory reviewed and updated to reflect:

  * new domain boundary (Customer without email)
  * updated repository contracts
  * adjusted API behavior
* This aligns implementation with documented architecture and avoids conceptual drift

### Conclusion

This commit introduces an important architectural refinement by **decoupling customer identity from authentication data**, resulting in:

* clearer aggregate boundaries
* improved domain purity
* better alignment with layered architecture
* explicit read model composition at delivery layer

From a design standpoint, this is a **high-quality correction**, reducing leakage of concerns into the domain and preparing the system for future evolution in authentication and user management.


## 2026/04/17 — api/user_status-05

Consolidates the **ledger model as the single source of truth for transactions**, removes the legacy operation abstraction, and introduces a robust idempotency and replay mechanism based entirely on persisted ledger data. This change also aligns repository contracts, database schema, and tests with the new model, while updating the entire documentation set to reflect the evolved architecture.

### 1. Ledger Consolidation and Domain Simplification

* Removed `Operation` entity and its entire persistence flow
* Eliminated `transactions` table as a separate idempotency store via migration
* Promoted `account_transactions` as the **canonical ledger**
* Extended `Transaction` entity:

  * added `related_account_id`
  * added `idempotency_key`
* Introduced `NewTransactionWithIdempotency` to support origin-side idempotent writes
* Replaced `ErrOperationAlreadyProcessed` with `ErrTransferDuplicate`, aligning error semantics with the new model

This change reinforces the domain principle that **all financial truth must be derivable from the ledger itself**, as already outlined in the domain model .

### 2. Transfer Use Case — Idempotency Redesign

* Replaced operation-based idempotency with **ledger-based replay**
* Introduced early idempotency check using:

  * `GetTransactionByIdempotencyKey`
* Implemented replay via `transferResultFromLedger`:

  * reconstructs result using `transfer_out` and paired `transfer_in`
  * validates ledger consistency (reference_id, related_account_id, type pairing)
* Added conflict handling:

  * DB unique constraint triggers `ErrTransferDuplicate`
  * reloads committed transaction and returns replay result
  * forces rollback while preserving response

This approach eliminates parallel state (operations table) and ensures **idempotency is enforced at the same layer that guarantees financial consistency**.

### 3. Repository Contract Evolution

* Removed:

  * `CreateOperation`
  * `GetOperationByIdempotencyKey`
* Added:

  * `GetTransactionByIdempotencyKey`
  * `GetTransactionByReference`
* Updated all mocks and tests accordingly

The repository now exposes only **ledger-native queries**, reducing conceptual duplication and improving cohesion.

### 4. PostgreSQL Layer Enhancements

* Updated `CreateTransaction`:

  * includes `related_account_id` and `idempotency_key`
  * uses `ON CONFLICT (account_id, idempotency_key) DO NOTHING`
  * detects duplicates via `RowsAffected == 0`
* Implemented:

  * `GetTransactionByIdempotencyKey`
  * `GetTransactionByReference`
* Extended query projections across the repository to include new fields

These changes move idempotency enforcement fully into the database, consistent with the system’s **strong consistency strategy** .

### 5. Database Migrations

* Added migration `000007_consolidate_ledger`:

  * moves idempotency data into `account_transactions`
  * converts `type` to enum
  * creates unique index for idempotency
  * drops legacy `transactions` table
* Added migration `000008_transfer_pair_integrity`:

  * enforces uniqueness of `(reference_id, type)` for transfer pairs
  * introduces index for efficient lookup

These migrations formalize:

* **idempotency at the ledger level**
* **structural integrity of transfer pairs**

### 6. Transfer Integrity Guarantees

* Enforced bidirectional linkage:

  * `transfer_out.related_account_id → destination`
  * `transfer_in.related_account_id → origin`
* Ensured:

  * both sides share the same `reference_id`
  * only one row per `(reference_id, type)`
* Added defensive checks in replay logic:

  * missing reference_id
  * missing related_account_id
  * mismatched pairing

This elevates the ledger from a log to a **verifiable structure**, not merely a record of events.

### 7. Test Suite Refactor

* Updated all mocks to align with new repository interface
* Removed operation-related expectations
* Added:

  * validation of `related_account_id` on both transfer legs
  * replay correctness using ledger data
  * duplicate conflict handling via DB constraint simulation
* Adjusted expectations:

  * single ledger write attempt before conflict
  * rollback behavior preserved
* Updated integration schema setup:

  * ensures backward compatibility with pre-existing databases
  * adds missing columns and indexes dynamically

### 8. Error Handling Adjustments

* Standardized duplicate semantics:

  * `ErrTransferDuplicate` mapped to idempotent replay
* Simplified authorization error message:

  * "Access denied to account" → "Access denied"
* Maintains consistency with API error contract 

### 9. Makefile Refinement

* Renamed migration commands:

  * `api-migrate-up` → `migrate-up`
  * `api-migrate-down` → `migrate-down`
* Grouped under a dedicated "Database Migrations" section

Improves clarity and aligns CLI with project scope.

### 10. Documentation Update (api/docs)

* Entire documentation set updated to reflect:

  * ledger as single source of truth
  * removal of `transactions` table usage
  * idempotency embedded in ledger
  * updated transfer flow and consistency guarantees
* Affects architecture, domain, data model, flows, and API contract
* Ensures alignment between implementation and documented behavior 

### Conclusion

This commit represents a **critical architectural refinement**, removing an artificial abstraction (`Operation`) and converging the system toward a **pure ledger-driven model**.

The resulting system is:

* more coherent (single source of truth)
* more robust (DB-enforced idempotency and integrity)
* easier to reason about (no dual-write model)
* better aligned with financial system principles

From a design standpoint, this is a substantial improvement. The previous model introduced unnecessary indirection; the current one correctly treats the ledger as both **execution log and verification mechanism**, which is the right direction for any system that aims at financial correctness.


## 2026/04/17 — api/user_status-04

Introduces **user status enforcement across account creation and admin approval flow**, strengthening authorization guarantees and aligning onboarding with an explicit lifecycle (pending → active).

### 1. Application Layer — CreateAccount Hardening

* Extended `CreateAccount` use case to depend on `UserRepository`
* Added validation pipeline before account creation:

  * user must exist
  * user must have valid `UserID`
  * user must be in `active` status
* Enforced strict access control:

  * any non-active user (pending, blocked) is rejected with `ErrForbidden`
* Prevents invalid states where accounts could be created for non-approved users
* This change aligns account creation with the authentication model where identity alone is insufficient without valid state 

### 2. New Admin Capability — Approve User

* Integrated `ApproveUserUseCase` into application wiring
* Added new protected endpoint:

  * `POST /admin/users/{id}/approve`
* Approval flow:

  * transitions user status to `active`
  * creates associated account atomically
* Extended output contract:

  * now includes `status` field alongside `user_id` and `account_id`
* Establishes explicit onboarding lifecycle:

  * register → pending → approved → active → operational
* This is a **structural improvement**, not just a feature addition

### 3. Delivery Layer — Authorization Enforcement

* Implemented `ApproveUser` handler with strict guards:

  * requires authenticated user
  * enforces `admin` role
  * validates UUID path parameter
* Maps domain errors consistently:

  * `FORBIDDEN`, `INVALID_DATA`, `USER_NOT_FOUND`, etc.
* Response includes:

  * `user_id`
  * `status`
  * `account_id`
* Maintains API contract consistency with existing envelope pattern 

### 4. Error Handling Standardization

* Added missing domain error mappings in auth layer:

  * `ErrForbidden`
  * `ErrInvalidData`
* Unified error message:

  * `"Access denied"` replaces inconsistent variants
* Removed duplicate error registration guard in shared mapper:

  * simplifies registry behavior
  * shifts responsibility to developer discipline
* Keeps alignment with global error strategy and response model 

### 5. Dependency Wiring (Composition Root)

* Updated `main.go`:

  * `CreateAccount` now receives `userRepo`
  * `ApproveUserUseCase` added to auth handler
* Ensures correct dependency flow:

  * Delivery → Application → Domain
* Reinforces modular monolith structure and explicit wiring rules 

### 6. Test Coverage Expansion

#### 6.1 CreateAccount Tests

* Added scenarios for:

  * user repository not configured
  * user not found
  * user lookup failure
  * pending user rejection
  * blocked user rejection
* Verified:

  * no repository side-effects on failure paths
  * correct interaction counts

#### 6.2 ApproveUser Tests

* Added full coverage:

  * success case (admin)
  * unauthorized request
  * non-admin rejection
  * invalid UUID
  * domain error mapping (not found, conflict, forbidden, internal)
* Validates both:

  * authorization boundary
  * HTTP contract correctness

#### 6.3 Integration Adjustments

* Updated handler constructors to include new dependency
* Ensured compatibility with existing integration tests

### 7. Domain Alignment

* Reinforces the concept that:

  * **user status is part of authorization context**, not just authentication
* Prevents illegal transitions such as:

  * financial operations executed by pending users
* Aligns with domain invariants where operations depend on valid state, not only identity 

### Conclusion

This commit introduces a **critical correction in the authorization model** by incorporating user status into the decision process.

The system now eliminates an important invalid state:
users could previously authenticate and operate without being formally approved.

From an architectural perspective, this is a **consistency and correctness fix**, ensuring that onboarding, authorization, and financial operations are coherently integrated.


## 2026/04/16 — api/user_status-03

Implements the **user approval flow with automatic account creation**, introducing transactional consistency across auth and account modules, strengthening user lifecycle control, and expanding domain and infrastructure support for status transitions.

### 1. Application Layer — ApproveUser Use Case

* Added `ApproveUserUseCase` to handle transition from `pending` to `active`
* Full transactional orchestration via `Transactor.RunInTx`
* Execution flow:

  * load user with `FindByIDForUpdate` (row-level lock)
  * validate user existence and current status
  * update user status to `active`
  * validate `customer_id` presence and existence
  * generate account number
  * create and persist new account
* Ensures **atomic activation + account creation**, preventing partial states
* Reuses `accountapplication.GenerateBranch()` for consistency with account module

### 2. Domain Layer Enhancements

* Introduced new domain error:

  * `ErrUserAlreadyActive`
* Reinforces user lifecycle invariants:

  * only `pending` users can be approved
  * active users cannot be reprocessed
* Aligns with invariant enforcement strategy described in the domain model 

### 3. Repository Contract Evolution

* Extended `UserRepository`:

  * added `FindByIDForUpdate` for pessimistic locking
* Enables safe concurrent approval handling, consistent with system-wide locking strategy 

### 4. Infrastructure Layer — PostgreSQL

* Implemented `FindByIDForUpdate` using:

  * `SELECT ... FOR UPDATE`
* Ensures:

  * row-level locking during approval
  * protection against concurrent status transitions
* Behavior:

  * returns `nil` when user not found (mapped at application level)

### 5. Error Handling Standardization

* Registered new domain errors in error registry:

  * `USER_NOT_FOUND → 404`
  * `USER_ALREADY_ACTIVE → 409`
* Added corresponding error codes in shared layer:

  * `ErrCodeUserNotFound`
  * `ErrCodeUserAlreadyActive`
* Maintains consistency with global API error contract 

### 6. Account Module Adjustment

* Refactored branch generation:

  * `generateBranch` → `GenerateBranch`
* Promotes reuse across modules and avoids duplication in account creation logic

### 7. Test Coverage

#### 7.1 Application Tests — ApproveUser

* Added comprehensive test suite:

  * success scenario (user activation + account creation)
  * user not found
  * user already active
  * missing customer_id
  * customer not found
  * account creation failure
* Validates transactional integrity and invariant enforcement

#### 7.2 Test Infrastructure Updates

* Updated mocks across auth tests to support:

  * `FindByIDForUpdate`
* Ensures compatibility with new repository contract without breaking existing tests

### 8. Architectural Impact

* Introduces a **controlled onboarding progression**:

  * register → pending user
  * approve → active user + account creation
* Eliminates invalid intermediate states:

  * active user without account
  * account created for unapproved user
* Strengthens consistency guarantees across modules, aligned with system transaction model 

### Conclusion

This commit introduces a **critical lifecycle transition for users**, integrating authentication and financial domains under a single transactional boundary.

The implementation is technically robust, particularly due to:

* explicit use of pessimistic locking
* strict invariant enforcement
* elimination of partial states

From an architectural perspective, this significantly improves the correctness and consistency of the onboarding flow, aligning it with the system’s financial integrity requirements.


## 2026/04/16 — api/user_status-02

Refactors user registration transaction handling, introduces **user status as a first-class domain attribute**, and standardizes persistence behavior across repository and application layers.

### 1. Application Layer — Transaction Handling Refactor

* Replaced repository-specific transaction coupling (`WithTransaction`) with a **generic `Transactor` abstraction**
* Updated `RegisterUserUseCase` to depend on `domain.Transactor` instead of casting the repository
* Execution now uses:

  * `transactor.RunInTx(ctx, fn)`
* Removes implicit assumptions about repository capabilities and enforces **explicit transaction orchestration at the application layer**
* Improves architectural consistency with the layered design already used in account operations 

### 2. Domain Layer — User Status Introduction

* Added `status` attribute to `User` entity lifecycle
* Introduced new domain error:

  * `ErrUserNotFound`
* Extended `UserRepository` contract:

  * added `UpdateStatus(userID, status)`
* Clarified repository responsibilities with explicit method semantics (e.g. `ExistsByEmail`)
* Aligns user lifecycle with a **state-driven model**, enabling future approval/activation flows

### 3. Persistence Layer — PostgreSQL Updates

* Added `status` column to `users` table:

  * `VARCHAR(20) NOT NULL DEFAULT 'pending'`
* Updated repository behavior:

  * `Create` now persists `status`
  * `FindByEmail` and `FindByID` now map `status`
  * implemented `UpdateStatus` with:

    * `UPDATE users SET status = $1, updated_at = NOW()`
    * returns `ErrUserNotFound` when no rows affected
* Removed legacy `WithTransaction` implementation from repository
* Consolidates responsibility: **repository handles persistence, application handles transactions**

### 4. Migration Layer

* Added migration:

  * `000006_user_status.up.sql`
  * `000006_user_status.down.sql`
* Ensures schema evolution is:

  * incremental
  * reversible
  * aligned with existing migration strategy

### 5. Test Infrastructure Refactor

* Updated all mocks to support new contract:

  * added `UpdateStatus` to `UserRepository` mocks
* Replaced transaction mocking:

  * removed `WithTransaction`
  * introduced `registerTransactorMock` with `RunInTx`
* Adjusted assertions:

  * validate `RunInTx` invocation instead of repository transaction calls

### 6. Integration Tests

* Updated integration setup:

  * ensures `users.status` column exists via `ALTER TABLE IF NOT EXISTS`
* Added validation for:

  * default status (`pending`) on creation
  * status transition via `UpdateStatus`
* Strengthens alignment between:

  * schema
  * repository
  * domain behavior

### 7. Wiring Adjustments

* Updated `main.go`:

  * `RegisterUserUseCase` now receives `transactor`
* Keeps dependency graph explicit and consistent with other use cases

### 8. Design Considerations

This change is particularly relevant from an architectural standpoint:

* eliminates hidden coupling between repository and transaction control
* enforces **application-layer ownership of transactional boundaries**
* introduces a **stateful user lifecycle**, which is essential for:

  * approval flows
  * onboarding pipelines
  * access control evolution

The decision to move away from `WithTransaction` is correct and aligns the auth module with the same rigor already present in financial operations.

### Conclusion

This commit is a structural improvement rather than a feature addition.

It establishes:

* a **clean transaction boundary model**
* a **state-driven user lifecycle**
* a **more consistent repository contract**

These changes prepare the system for more advanced flows such as user approval, activation, and policy-based authorization without requiring further architectural refactoring.


## 2026/04/16 — api/user_status-01

Introduces **user status management** into the authentication domain and restructures the HTTP layer to support a clearer multi-stage authentication model, including AppToken-based onboarding and JWT-protected routes. 

### 1. Domain Layer — User Status

* Added `UserStatus` type with explicit states:

  * `pending`
  * `active`
  * `blocked`
* Extended `User` entity to include `Status` field
* Updated `NewUser` factory:

  * initializes all new users with `UserStatusPending`
* This change establishes a **foundation for lifecycle control** (approval, activation, blocking), which was previously absent in the model

### 2. Bootstrap — Environment Configuration

* Introduced automatic `.env` loading using `github.com/joho/godotenv`
* Implemented flexible resolution strategy:

  * local `.env`
  * `api/.env`
  * executable-relative paths
* Ensures configuration is available regardless of execution context
* Aligns with **fail-fast configuration validation** already present in `main.go`

### 3. Main Wiring Refactor (cmd/api/main.go)

* Reorganized startup into explicit sections:

  * Config
  * Repositories
  * Services
  * Use Cases
  * Handlers
  * Middlewares
  * Routers
* Enforced validation of critical environment variables:

  * `APP_TOKEN`
  * `JWT_SECRET`
* Improved readability and maintainability of composition root
* This is a **structural improvement**, not just cosmetic; it clarifies dependency boundaries

### 4. Routing Architecture — Separation of Concerns

* Split routing into three layers:

  * `authRouter` (authentication endpoints)
  * `apiRouter` (business endpoints)
  * `mainRouter` (composition)
* Introduced explicit middleware application per route group:

  * AppToken for onboarding
  * JWT for authenticated access
* Eliminated global middleware wrapping, replacing it with **route-level control**, which is more precise and safer

### 5. Authentication Model — AppToken + JWT

* Applied AppToken middleware to:

  * `POST /auth/register`
  * `POST /auth/login`
* Applied JWT middleware to:

  * `POST /auth/refresh`
  * `GET /auth/me`
  * all `/accounts/*`
  * `/customers/me`
* This formalizes a **two-phase authentication model**:

  * controlled entry (AppToken)
  * authenticated session (JWT)
* Matches the intended design described in the authentication documentation 

### 6. API Contract Documentation Updates

* Updated REST documentation to reflect:

  * AppToken requirement for onboarding endpoints
  * JWT requirement for all protected endpoints
* Added new error code:

  * `INVALID_APP_TOKEN` (HTTP 401)
* Expanded error scenarios with concrete payload examples
* Clarified access control rules and authentication flows
* Improves alignment between implementation and public contract 

### 7. Dependency Updates

* Added `godotenv` dependency to `go.mod` and `go.sum`
* Enables environment-based configuration without external tooling

### 8. Architectural Impact

* Introduces the first step toward **user lifecycle governance** via status
* Establishes a clearer boundary between:

  * onboarding security
  * session-based authentication
  * resource authorization
* Prepares the system for future features such as:

  * user approval workflows
  * account activation
  * access blocking

### Conclusion

This commit is a **strategic evolution of the authentication layer**, not merely a feature addition.

It introduces user lifecycle semantics and formalizes a multi-stage authentication model, improving both **security posture and architectural clarity**, while keeping the system aligned with its current simplicity goals and ready for future extensions.


## 2026/04/16 — api/app_token-01

Introduces **application-level request validation via App Token middleware**, enforces stricter environment configuration, and refactors HTTP server initialization to support middleware composition and improved security boundaries.

### 1. Security — App Token Middleware

* Implemented `AppToken` middleware to enforce presence and validity of `X-App-Token` header
* Uses `crypto/subtle.ConstantTimeCompare` to prevent timing attacks during token comparison
* Rejects unauthorized requests early in the pipeline with standardized error response
* Integrates with existing error contract via `ErrInvalidAppToken`
* Establishes a clear separation between:

  * client identity (JWT)
  * client application validation (App Token)

This aligns with a layered security model, where multiple independent signals are validated before request execution, consistent with the system’s architectural direction 

### 2. Error Standardization

* Added `ErrInvalidAppToken` to shared sentinel errors
* Ensures consistency with API response envelope (`data` / `error`) 
* Avoids inline error construction, reinforcing centralized error definitions

### 3. HTTP Pipeline Refactor

* Replaced global `http.Handle` usage with explicit `http.ServeMux`
* Introduced handler composition:

  * `router → auth middleware → app token middleware`
* Final pipeline:

  * `handler := AppToken(...) (router)`
* Improves:

  * composability
  * testability
  * visibility of request flow

### 4. Environment Hardening

* Enforced mandatory environment variables:

  * `APP_TOKEN` now required (fail-fast with `log.Fatal`)
  * `JWT_SECRET` no longer has fallback value
* Eliminates insecure default configurations
* Guarantees that authentication and application validation cannot run in an invalid state

### 5. Routing Organization

* Centralized all routes into `ServeMux`
* Maintains explicit registration of:

  * auth endpoints
  * customer endpoints
  * account endpoints (including transfer)
* Preserves existing authorization behavior via JWT middleware

### 6. Test Coverage — Middleware

* Added comprehensive tests for `AppToken`:

  * missing header
  * invalid token
  * valid token (happy path)
* Validates:

  * correct HTTP status (`401`)
  * response envelope structure
  * prevention of downstream handler execution on failure

### 7. Developer Experience

* Added `api-run` command to `Makefile` for local server execution
* Introduced `.gitignore` for API build artifacts

### 8. Documentation Reorganization

* Moved API documentation into `api/docs`
* Improves cohesion between code and documentation
* Maintains consistency with modular project structure

### Conclusion

This commit introduces a **critical security boundary at the application level**, ensuring that requests are validated not only by user identity (JWT) but also by client context (App Token).

From an architectural standpoint, this is a meaningful step toward a **multi-signal validation model**, where authentication alone is no longer treated as sufficient.


## 2026/04/16 — api/refresh_token-02

Refines authentication flow and project documentation, introducing **refresh token persistence on the client side** and restructuring the repository documentation to reflect the current architecture and usage model.

### 1. Mobile — Refresh Token Support

* Extended `LoggedUser` model:

  * added `refreshToken` field
  * updated `fromMap` to deserialize `refresh_token`
* Updated authentication repository:

  * persist `refreshToken` alongside `accessToken` on login
  * ensure removal of both tokens on logout
* This aligns the mobile client with the backend session model, enabling **token rotation and session continuity**

### 2. Authentication Flow Consistency

* Ensures that client-side storage reflects server expectations:

  * access + refresh token pair becomes the canonical session representation
* Prepares the mobile layer for:

  * automatic token refresh via interceptor
  * retry logic on `401` responses
* Reinforces contract implied by auth endpoints and JWT usage 

### 3. Monorepo Documentation Simplification

* Rewrote root `README.md`:

  * removed narrative-heavy content
  * introduced concise structure and quick-start flow
  * clarified dual-app nature (API + mobile)
* Focus shifted from conceptual description to **operational clarity**

### 4. API Documentation Restructuring

* Simplified `api/README.md`:

  * clearer separation of stack, architecture, and features
  * explicit route listing
  * streamlined setup instructions
* Removed legacy architecture document:

  * replaced with new `ARCHITECTURE.md`
* New architecture document:

  * formalizes modular monolith structure
  * clarifies layer responsibilities and dependency direction
  * documents authentication and refresh flow behavior
* Maintains alignment with layered architecture principles 

### 5. Documentation Standardization (Mobile)

* Rewrote `mobile/README.md`:

  * emphasizes role as integration client
  * adds environment configuration guidance
  * documents test and build workflows
* Introduced `docs/mobile/ARCHITECTURE.md`:

  * defines layered structure (UI, Data, Domain, Core)
  * formalizes request flow and interceptor behavior

### 6. Licensing

* Added MIT license to API module:

  * clarifies usage and distribution terms
  * aligns repository with open-source conventions

### 7. Structural Improvements

* Normalized directory descriptions across README files
* Improved onboarding flow:

  * Docker → migrations → API → mobile
* Reduced redundancy across documentation layers

### Conclusion

This commit is primarily a **consistency and alignment step** between client, API, and documentation.

It establishes:

* a **complete token lifecycle on the client (access + refresh)**
* a **clearer and more operational documentation structure**
* a **more explicit architectural baseline for future evolution**

From a design standpoint, this is a necessary consolidation step before advancing into more complex authentication concerns such as concurrent refresh handling and session control.


## 2026/04/15 — api/refresh_token-02

Refactors and hardens the refresh token flow to guarantee **atomic rotation, consistency, and correctness of session management**, while simplifying dependency contracts and eliminating invalid execution paths.

---

### 1. Application Layer — Refresh Token Atomicity

* Introduced `Transactor` as a first-class dependency in `RefreshAccessTokenUseCase`
* Implemented atomic rotation using `RunInTx`:

  * `Revoke(old_token)` + `Create(new_token)` executed in a single transaction
* Removed non-transactional revoke operation
* Guarantees:

  * no partial state (no “revoked without replacement” or “duplicated sessions”)
  * rollback preserves original token validity on failure
* Aligns transaction control with Application layer responsibilities 

---

### 2. Infrastructure Layer — Transaction Support

* Added `PostgresTransactor`:

  * wraps `pgx` transaction lifecycle (`Begin → Commit / Rollback`)
  * injects transaction into context (`ContextWithTx`)
* Enables multiple repositories to share the same transaction transparently
* Strengthens infrastructure compliance with domain contracts (`Transactor` interface)

---

### 3. Domain Layer — Contract Expansion

* Introduced `Transactor` interface:

  * explicit control of transactional boundaries at use case level
* Added `ErrSessionNotFound`:

  * ensures revoke failures are explicit and not silently ignored
* Reinforces domain-driven consistency for session lifecycle

---

### 4. Session Repository — Correctness Enforcement

* Updated `Revoke` implementation:

  * now validates `RowsAffected()`
  * returns `ErrSessionNotFound` when token does not exist
* Eliminates silent inconsistencies in session state

---

### 5. Login Use Case — Contract Tightening

* Removed optional `SessionRepository` dependency (variadic → required)
* Enforced invariant:

  * **no refresh token is issued without persisted session**
* Simplified logic:

  * always hashes and persists refresh token
  * failure in session creation aborts login
* This removes previously possible invalid states

---

### 6. Delivery Layer — Dependency Integrity

* Refactored `Handler` constructor:

  * now requires all use cases upfront (including refresh)
  * removed setter injection (`SetRefreshAccessTokenUseCase`)
* Ensures:

  * no partially initialized handlers
  * no runtime mutation of dependencies
* Added defensive check for `nil` output in login flow

---

### 7. Wiring (main.go & integration)

* Registered `PostgresTransactor` in composition root
* Injected into `RefreshAccessTokenUseCase`
* Updated handler initialization to reflect new constructor contract
* Ensures consistent dependency graph across application and tests

---

### 8. Test Suite — Coverage Expansion

#### 8.1 Refresh Flow Tests

* Updated all tests to include `Transactor` dependency
* Added `transactorMock` for transactional execution

#### 8.2 Rotation Integrity Test (Stateful)

* Introduced `statefulSessionMock` to simulate real session lifecycle
* Validates full rotation behavior:

  * old token becomes unusable after refresh
  * new token is immediately valid
  * reuse of revoked token fails correctly
* This is a critical validation of **system invariants**

#### 8.3 Login Tests

* Updated to reflect mandatory `SessionRepository`
* Ensures session persistence is always exercised

#### 8.4 Infrastructure Test Fix

* Simplified JWT error assertion (removed invalid `errors.Is` usage)

---

### 9. API & Documentation Updates

* Login now explicitly returns `refresh_token`
* Introduced `/auth/refresh` endpoint with:

  * token rotation semantics
  * single-use refresh tokens
  * atomic revoke + create behavior
* Documented error scenarios:

  * invalid, expired, revoked, or missing sessions
* Clarifies contract for clients and aligns behavior with implementation 

---

### Conclusion

This commit is a **consistency and correctness milestone** for authentication flow.

It eliminates entire classes of invalid states by enforcing:

* **atomic token rotation**
* **mandatory session persistence**
* **explicit failure handling**
* **constructor-level dependency integrity**

From an architectural standpoint, this is a decisive improvement:
transaction boundaries are now correctly owned by the Application layer, and the authentication model becomes **predictable, verifiable, and resilient under failure conditions**.


## 2026/04/11 — api/refresh_token-01

Implements a **complete refresh token flow with session management**, evolving the authentication model from stateless JWT-only to a **stateful, revocable, and rotating session-based approach**. This change aligns the system with a more robust security posture while preserving the layered architecture principles 

---

### 1. Domain Layer — Contracts Expansion

* Extended `TokenService`:

  * added `GenerateRefreshToken(userID)`
  * added `ParseRefreshToken(token)`
* Introduced `SessionRepository`:

  * `Create`
  * `FindByTokenHash`
  * `Revoke`
* Establishes **explicit session lifecycle control** at the domain boundary

---

### 2. Application Layer — Login Flow Evolution

* `LoginUserUseCase` updated to:

  * generate **access token + refresh token**
  * hash refresh token using `SHA-256`
  * persist session with expiration (`30 days TTL`)
* Output now includes:

  * `AccessToken`
  * `RefreshToken`
* Optional session dependency supported (backward-safe injection)

**Key observation (architectural):**
This is the first point where authentication becomes **state-aware**, breaking the purely stateless JWT model intentionally.

---

### 3. Application Layer — Refresh Token Use Case

* Introduced `RefreshAccessTokenUseCase`:

  * validates refresh token integrity
  * validates session (existence, expiration, revocation, ownership)
  * loads user from repository
  * generates new access token
  * performs **refresh token rotation**:

    * revoke old token
    * generate new refresh token
    * persist new session

**Security properties introduced:**

* replay protection (rotation)
* server-side invalidation
* binding between token and stored session

---

### 4. Infrastructure Layer — Token Service

* Extended `JWTTokenService`:

  * added **opaque refresh token generation**

    * payload: `userID + nonce`
    * signature: `HMAC-SHA256`
    * encoding: `base64url`
  * implemented `ParseRefreshToken`

    * signature validation using constant-time comparison
    * strict payload validation

* Access token improvements:

  * ensured `exp` claim correctness with TTL enforcement

**Technical decision:**
Refresh token is **not JWT**, which is a correct choice to:

* reduce attack surface
* simplify validation
* avoid overloading JWT semantics

---

### 5. Infrastructure Layer — Session Persistence

* Added `PostgresSessionRepository`

  * `Create`: inserts hashed token
  * `FindByTokenHash`: retrieves session state
  * `Revoke`: soft-revokes via `revoked_at`

* Supports transaction-aware execution via context

* Migration `000005_user_sessions`:

  * new table `user_sessions`
  * indexed by `user_id` and `expires_at`
  * unique constraint on `token_hash`

**Critical design choice:**

* only **hashed tokens are stored**
* prevents token leakage from DB compromise

---

### 6. Delivery Layer — HTTP Contract Updates

* `/auth/login`:

  * now returns `refresh_token`

* New endpoint:

  * `POST /auth/refresh`

* Handler additions:

  * request validation (`refresh_token`)
  * consistent error mapping
  * response envelope preserved

* Introduced DTOs:

  * `refreshAccessTokenRequest`
  * `refreshAccessTokenData`

**Important:**
This extends the API contract beyond what is currently documented  and requires documentation update.

---

### 7. Dependency Wiring (main.go)

* Registered:

  * `SessionRepository`
  * `RefreshAccessTokenUseCase`
* Injected into:

  * `LoginUserUseCase`
  * handler via setter
* Exposed route:

  * `POST /auth/refresh`

---

### 8. Test Coverage

#### 8.1 Application Tests

* Login:

  * validates access + refresh generation
  * verifies session persistence (hash + expiration)
  * covers failure scenarios:

    * access token failure
    * refresh token failure
    * session persistence failure

* Refresh flow:

  * success path
  * invalid token
  * session not found
  * revoked session
  * expired session
  * user mismatch
  * repository failures
  * rotation integrity

#### 8.2 Infrastructure Tests

* Refresh token:

  * generation + parsing
  * entropy validation
  * tampering detection
  * malformed token handling
* Access token:

  * expiration correctness

#### 8.3 Delivery Tests

* Login:

  * validates response now includes `refresh_token`
* Refresh:

  * success case
  * invalid token → `401 INVALID_TOKEN`

#### 8.4 Integration Tests

* End-to-end validation:

  * login returns both tokens
  * refresh endpoint wired correctly
* Test DB isolation:

  * switched to `bank_test`
* CPF constraint repair added for test consistency

---

### 9. Test & Environment Adjustments

* Updated default test database:

  * `bank → bank_test`
* Added defensive SQL for constraint repair:

  * prevents flaky test runs due to regex mismatch

---

### 10. Behavioral Changes Summary

* Access tokens are now **short-lived**
* Refresh tokens are:

  * generated securely
  * persisted as hashed values
  * validated against DB
  * rotated on use
* Authentication becomes:

  * **stateful**
  * **revocable**
  * **traceable**

---

### Conclusion

This commit represents a **major security and architectural milestone**:

* transitions authentication from **stateless JWT** to **session-backed model**
* introduces **refresh token rotation**, a critical protection against replay attacks
* enforces **server-side control over sessions**, enabling future features such as:

  * logout
  * device/session listing
  * anomaly detection

From a technical standpoint, the implementation is **well-aligned with Clean Architecture principles**, keeping:

* domain contracts pure
* application responsible for orchestration
* infrastructure isolated

The only architectural caveat is the **optional session dependency in login**, which introduces a potential inconsistency. In a production-grade system, this should be mandatory to avoid issuing unusable refresh tokens.

Overall, this is a **production-grade foundation for authentication**, suitable for fintech-level requirements.


## 2026/04/10 — infra/layout-01

Introduces a **UI layout standardization layer** for the Flutter application, centralizing structural concerns and improving consistency across authentication screens, while also refining routing behavior and state handling patterns.

### 1. Routing Adjustment

* Updated initial route:

  * from `HomeRoutes.home` to `AuthRoutes.login`
* Aligns application startup with authentication flow, enforcing a more realistic entry point for protected systems
* This change is consistent with the backend contract where authentication precedes access to account resources 

---

### 2. Introduction of SafeScaffold

* Added new base component: `SafeScaffold`
* Encapsulates:

  * `SafeArea` handling for body and bottom navigation
  * consistent horizontal constraints (`maxWidth: 460`)
  * standardized padding for bottom actions
* Provides a **reusable layout abstraction**, reducing duplication and enforcing UI consistency
* Conceptually aligns with separation of responsibilities seen in the backend architecture, isolating structural concerns from business/UI logic 

---

### 3. Login Page Refactor

* Migrated from `Scaffold` to `SafeScaffold`
* Introduced `AppBar` for clearer navigation structure
* Refactored state handling:

  * replaced `setState` with `ValueNotifier<bool>` for password visibility
* Improved layout:

  * consistent spacing using `Column.spacing`
  * moved primary action to `bottomNavigationBar`
  * added `GestureDetector` to dismiss keyboard
* Decoupled navigation logic into dedicated methods (`_navToRegister`)
* Replaced direct widget access with local `_viewModel` reference for better readability and lifecycle control

---

### 4. Register Page Refactor

* Applied same structural pattern as Login:

  * `SafeScaffold`
  * `AppBar`
  * bottom action bar for primary CTA
* Introduced local `_viewmodel` reference
* Improved layout consistency:

  * removed redundant spacing widgets
  * standardized vertical rhythm using `spacing`
* Added explicit navigation method (`_navToLogin`)
* Ensures both auth screens follow the same **visual and interaction contract**

---

### 5. UI Behavior Improvements

* Centralized primary actions (Entrar / Cadastrar) in bottom area:

  * improves ergonomics on mobile devices
  * creates a consistent interaction pattern
* Added loading state handling directly in action buttons
* Improved keyboard UX with tap-to-dismiss behavior

---

### 6. Architectural Considerations

This change is subtle but important from a design perspective:

* Introduces a **UI composition layer**, analogous to how backend layers isolate responsibilities
* Reduces duplication while preserving flexibility
* Moves toward a **design system mindset**, even without formalizing one yet

A critical observation:
this abstraction is well-scoped. It does not attempt to generalize business logic or navigation, only layout concerns. This is a good boundary and avoids premature over-engineering.

---

### Conclusion

This commit establishes a **foundation for consistent UI composition**, improving maintainability, readability, and user experience.

The introduction of `SafeScaffold` combined with the refactoring of authentication screens represents a **clear step toward a scalable UI architecture**, mirroring the layered discipline already present in the backend.


## 2026/04/10 — infra/routing-01

Introduces a **structured routing architecture using GoRouter**, along with UI composition, dependency injection integration, and initial authentication flows. This commit establishes a clear separation of routing concerns aligned with a modular layered approach 

### 1. Routing Architecture Refactor

* Replaced monolithic route definition with **modular route groups**:

  * `authRoutes()`
  * `homeRoutes()`
* Router now composes routes using spread operators, improving scalability and readability
* Updated `initialLocation` to use `HomeRoutes.home.path`, removing reliance on generic enums

### 2. Route Definition Strategy

* Replaced generic `Routes` enum with **domain-oriented route enums**:

  * `AuthRoutes`
  * `HomeRoutes`
* Each enum encapsulates its own path, improving cohesion and reducing accidental coupling
* Introduced dedicated route files:

  * `routes/auth_routes.dart`
  * `routes/home_routes.dart`

Opinion: This is a strong architectural move. It prevents the typical “god enum” anti-pattern and aligns routing with feature boundaries.

### 3. GoRouter Integration

* Migrated from `MaterialApp` to `MaterialApp.router`
* Centralized router creation via `router()` factory
* Added `ExtraCodec` support for serialization:

  * now explicitly supports `null` values
  * prevents runtime failures when passing optional navigation data

### 4. Dependency Injection Integration

* Introduced `Uis.add(injector)` into dependency setup
* ViewModels are now resolved directly in route builders via injector:

  * `LoginViewModel`
  * `RegisterViewmodel`
  * `HomeViewmodel`
* Removed redundant LocalSecureStorage registration from `Data` layer, keeping DI responsibilities better distributed

Opinion: Injecting ViewModels at the routing boundary is a pragmatic choice. It keeps UI decoupled while avoiding premature abstraction layers.

### 5. Application Entry Point Refactor

* Renamed `MainApp` to `AppWidget`
* Moved it into `/uis`, reinforcing UI ownership
* Introduced internal router instance (`GoRouter`) inside the widget
* Replaced `home:` with `routerConfig`, aligning app initialization with navigation system

### 6. Authentication UI Implementation

#### Login Flow

* Implemented full `LoginPage`:

  * form validation (email/password)
  * loading state via `Command`
  * success/failure feedback using `SnackBar`
* Navigation:

  * success → `HomeRoutes.home`
  * register link → `AuthRoutes.register`

#### Register Flow

* Replaced placeholder with full implementation:

  * fields: name, email, cpf, password
  * validation rules for each field
  * command-based execution
* Navigation:

  * success → `AuthRoutes.login`

### 7. ViewModel Layer Introduction

* Added ViewModels:

  * `LoginViewModel`
  * `RegisterViewmodel`
  * `HomeViewmodel`
* Standardized usage of `Command1` for async actions
* Established consistent interaction pattern:

  * UI observes command state
  * ViewModel delegates to repository

### 8. UI Composition Adjustments

* `HomePage` now receives `HomeViewmodel` via constructor
* Ensures consistency with DI-driven UI pattern
* Created centralized `uis.dart` for ViewModel registration

### 9. Codebase Cleanup and Direction

* Removed unused imports and redundant DI registrations
* Added note to relocate `getProfile` from `AuthApi` to a future profile service
* Introduced (commented) navigation extension for future evaluation

### Conclusion

This commit represents a **foundational shift in navigation and UI architecture**, achieving:

* modular routing aligned with feature boundaries
* clean integration between routing and dependency injection
* consistent ViewModel-driven UI pattern
* scalable structure for future expansion (auth, home, and beyond)

From an architectural standpoint, this is a well-directed evolution. The system moves closer to a **feature-oriented modular design**, reducing global coupling and improving long-term maintainability.


## 2026/04/10 — infra/http-client-setup-01

Establishes a **centralized and environment-driven HTTP client configuration**, removing runtime mutation patterns and aligning the mobile client with a more deterministic and infrastructure-oriented design.

### 1. Environment Configuration Refactor

* Introduced `AppEnv` as the single source of truth for runtime configuration:

  * `baseUrl` with strict validation (non-empty and valid URI)
  * `connectTimeout` and `receiveTimeout` via compile-time environment variables
  * `AppMode` enum with explicit parsing and validation
* Removed legacy `EnviromentKey`, eliminating loosely validated configuration access
* This change enforces **fail-fast behavior**, which is a critical improvement for reliability in distributed systems

### 2. HTTP Client Design Simplification

* Removed `setBaseUrl` from `RestClient` interface and its implementation
* Eliminated runtime base URL mutation across the application layer
* All configuration is now resolved at instantiation time via `DioFactory`
* This is a **significant architectural improvement**, as it:

  * removes hidden side effects
  * avoids per-request configuration inconsistencies
  * enforces immutability of infrastructure concerns

### 3. DioFactory Redesign

* Refactored `DioFactory` to return a configured `Dio` instance instead of `RestClient`
* Integrated `AppEnv` directly into `BaseOptions`:

  * `baseUrl`
  * timeouts
  * default headers
* Added support for optional `defaultHeaders`
* Improved interceptor registration:

  * avoids duplicate interceptor instances using type comparison
* This aligns the HTTP client with an **infrastructure-first responsibility model**, consistent with layered architecture principles 

### 4. Dependency Injection Restructuring

* Reorganized `CoreServices` with explicit layering:

  1. `FlutterSecureStorage`
  2. `LocalSecureStorage` abstraction
  3. base `Dio` instance
  4. `AuthInterceptor` with isolated configuration
  5. `RestClient` composed from `Dio`
* Notable design decision:

  * `AuthInterceptor` uses a dedicated `Dio` instance to avoid recursive interception
* This setup improves:

  * testability
  * separation of concerns
  * predictability of request flow

### 5. API Layer Cleanup

* Removed manual base URL overrides from `AuthApi`
* All endpoints now rely on centralized configuration
* This eliminates duplication and prevents divergence across API calls
* Aligns the client with a **contract-driven API consumption model** 

### 6. Interceptor Behavior Clarification

* Updated `AuthInterceptor` comment to explicitly document behavior:

  * skips token injection when `Authorization` header is already present
* Improves readability and reduces ambiguity in request handling

### 7. Test Adjustments

* Updated `DioRestClient` tests:

  * removed dependency on `setBaseUrl`
  * now validate behavior based on `Dio` configuration
* Ensures tests reflect the new immutable configuration model

### Conclusion

This commit represents a **structural upgrade of the HTTP client layer**, shifting from mutable, scattered configuration to a **centralized, deterministic, and environment-driven approach**.

From an architectural standpoint, the most relevant gain is the clear separation between **application logic and infrastructure concerns**, reinforcing the principles of layered architecture and significantly reducing the risk of inconsistent network behavior across the application.


## 2026/04/09 — infra/di-and-env-setup-01

Establishes the **foundational infrastructure layer for dependency injection and environment configuration** in the Flutter client, aligning the mobile architecture with a modular, scalable structure and enabling controlled environment-based execution.

### 1. Development Environment Configuration

* Added `.vscode/launch.json` with predefined run configurations:

  * Dev, Staging, Prod
  * Integration test profile (Dev)
* Each configuration uses `--dart-define-from-file`, enabling **externalized environment configuration**
* Introduced `.env` file strategy (`dev.env`, `staging.env`, `prod.env`) and ensured they are ignored via `.gitignore`
* This approach is technically sound and aligns with production-grade practices for **environment isolation and reproducibility**

### 2. Dependency Injection Setup

* Introduced centralized DI configuration via `dependencies.dart`
* Adopted `AutoInjector` as DI container
* Implemented idempotent initialization (`_initialized` guard)
* Structured registration into modular layers:

  * `CoreServices`
  * `Services`
  * `Data`
* This is a **critical architectural improvement**, bringing the mobile project closer to the same separation principles already present in the backend 

### 3. Core Services Layer

* Added `CoreServices` module:

  * Registers `FlutterSecureStorage`
  * Configures `RestClient` via `DioFactory`
* Environment-driven configuration:

  * `baseUrl` via `EnviromentKey`
  * timeouts defined explicitly
* This enforces **centralized HTTP client configuration**, avoiding scattered setup across the codebase

### 4. Environment Abstraction

* Introduced `EnviromentKey`:

  * Maps compile-time variables using `String.fromEnvironment` and `int.fromEnvironment`
* Supports:

  * base URL
  * timeouts
  * app mode
  * access token (for internal usage)
* This design is particularly robust, as it avoids runtime parsing and ensures **compile-time guarantees**

### 5. Data Layer Composition

* Introduced `Data` module for DI registration:

  * `LocalSecureStorage` abstraction
  * `AuthRepository` implementation
* Proper dependency chaining:

  * Repository depends on API + storage
* This reinforces the **Repository as SSOT pattern**, consistent with your architectural direction

### 6. Services Layer Refactor

* Introduced `Services` module:

  * Registers `AuthApi` with injected `RestClient`
* Removed legacy empty `services.dart`
* Clean separation between:

  * core infrastructure (HTTP, storage)
  * feature services (API layer)

### 7. Authentication Repository Implementation

* Added `AuthRepository` contract and `AuthRepositoryImpl`
* Responsibilities:

  * manage authentication state (`currentUser`, `isLoggedIn`)
  * persist access token
  * handle login, logout, register, and profile
* Introduced explicit unauthenticated handling:

  * new `AppErrorCode.unauthenticated`
* This is a **well-structured implementation**, with clear boundaries between:

  * API (remote)
  * storage (local)
  * state (in-memory)

### 8. Storage and Auth Adjustments

* Renamed `authToken` → `accessToken` for semantic clarity
* Updated `AuthInterceptor` to use new key consistently
* Improved session lifecycle:

  * proper token write on login
  * cleanup on logout and refresh failure
* These changes reduce ambiguity and improve long-term maintainability

### 9. Application Bootstrap

* Updated `main.dart`:

  * introduced `setupDependencies()` before `runApp`
* Ensures all dependencies are resolved prior to UI initialization
* Aligns with proper application lifecycle control

### 10. Minor Improvements

* Adjusted imports in `AuthApi`
* Improved test launch configuration for integration tests
* Small consistency fixes across modules

### Conclusion

This commit introduces a **structural turning point in the mobile application architecture**.

Key gains:

* centralized dependency management
* environment-driven configuration
* clear separation of layers (core, services, data)
* improved authentication flow consistency

From an architectural perspective, this is a **necessary and well-executed foundation**, enabling the project to scale without accumulating coupling or configuration debt.


## 2026/04/09 — theme/composition-01

Introduces a structured **theme composition system** for the Flutter application, including dynamic theme resolution, Material 3 integration, custom typography, and improvements in developer tooling via Makefile refinements.

### 1. Theme Composition Architecture

* Refactored `MainApp` from `StatelessWidget` to `StatefulWidget` to support context-dependent initialization
* Introduced controlled theme composition flow:

  * resolve system brightness (`platformBrightness`)
  * select base theme (`light` / `dark`)
  * apply app-level overrides via `_buildAppTheme`
* Encapsulates theme creation logic, improving cohesion and avoiding scattered configuration across widgets
* This approach is conceptually aligned with layered responsibility principles, where configuration is centralized and isolated 

### 2. Material Theme Abstraction

* Added `MaterialTheme` class:

  * centralizes all `ColorScheme` definitions
  * supports multiple variants:

    * light / dark
    * medium contrast
    * high contrast
* Provides factory methods:

  * `light()`, `dark()`, and contrast variations
* Uses Material 3 (`useMaterial3: true`)
* Ensures consistency and scalability of design tokens across the application
* This is a **notable improvement in design maturity**, replacing ad-hoc theming with a reusable and extensible system

### 3. Typography System with Google Fonts

* Introduced `createTextTheme` helper:

  * composes two font families:

    * body font (Quicksand)
    * display font (EB Garamond)
* Uses `google_fonts` package for runtime font resolution
* Merges text styles to preserve semantic roles (`body`, `label`, etc.)
* Enables consistent typography without coupling UI components to font configuration

### 4. Dynamic Theme Initialization

* Theme is initialized in `didChangeDependencies`:

  * ensures access to `BuildContext`
  * avoids unnecessary recomputation
* Separation between:

  * theme construction (`MaterialTheme`)
  * runtime selection (`brightness`)
  * UI overrides (`AppBarTheme`)
* Improves maintainability and testability of UI configuration

### 5. UI Adjustments

* Updated `AppBar` styling:

  * uses `primaryContainer` and `onPrimaryContainer`
  * enforces semi-bold title (`FontWeight.w600`)
* Minor text change in HomePage:

  * "Home Page" → "Type Home Page"

### 6. Dependency Updates

* Added `google_fonts` dependency for typography support
* Introduced transitive dependency `http` (via ecosystem resolution)

### 7. Makefile Improvements

* Added `tests` target:

  * aggregates `api-test` and `mobile-test`
* Renamed Flutter commands for consistency and ergonomics:

  * `flutter-clean` → `fclean`
  * `flutter-build` → `fbuild`
* Added new utility:

  * `fadd pkg=<name>` to simplify dependency installation
* Improves developer experience and standardizes command usage across environments

### Conclusion

This commit establishes a **robust and scalable theming foundation**, transitioning from a basic configuration to a **composable design system** with clear separation of concerns.

From a technical standpoint, the introduction of a dedicated theme layer combined with dynamic resolution and Material 3 alignment significantly improves maintainability, consistency, and long-term extensibility of the UI layer.


## 2026/04/08 — main

Restructures the repository into a cohesive **monorepo architecture**, consolidating backend, mobile, infrastructure, and documentation while improving developer experience, build orchestration, and project clarity.

### 1. Monorepo Consolidation

* Introduced unified repository structure:

  * `api/` (Go backend)
  * `mobile/` (Flutter client)
  * `infra/` (Docker/infrastructure)
  * `docs/` (centralized documentation)
* Promoted project to a **full-stack system workspace**, aligning backend and mobile under a single lifecycle
* Reinforces the modular monolith approach described in the architecture documentation 

### 2. Documentation Reorganization

* Moved all API documentation from `api/docs/` → `docs/api/`
* Updated all internal references to reflect new structure
* Centralized architectural and API design artifacts:

  * architecture
  * domain model
  * use cases
  * API contract
* Improves discoverability and enforces documentation as a **first-class artifact of the system design**

### 3. Root-Level README Overhaul

* Replaced minimal README with comprehensive project documentation:

  * system purpose and engineering goals
  * architectural overview (layered modular monolith)
  * API capabilities and guarantees
  * mobile role as integration validator
  * local development workflow
* Explicitly documents:

  * transactional consistency strategy
  * concurrency handling (row-level locking)
  * API contract conventions
* Aligns with the REST contract and system behavior expectations 

### 4. Build and Tooling Unification

* Introduced root-level `Makefile` as a **monorepo task runner**
* Added commands:

  * Docker lifecycle (`docker-up`, `docker-down`, `docker-logs`)
  * Flutter utilities (`flutter-clean`, `flutter-build`)
* Removed duplicated Makefiles from:

  * `api/`
  * `mobile/`
* Establishes a **single entry point for all development workflows**, reducing operational fragmentation

### 5. Infrastructure Standardization

* Moved `docker-compose.yml` to repository root
* Simplifies environment setup and aligns with monorepo conventions
* Enables consistent orchestration across backend and mobile dependencies

### 6. Dependency Management Improvements (Go)

* Promoted key dependencies from indirect to direct:

  * `jwt`
  * `uuid`
  * `pgx`
  * `crypto`
* Updated `go.sum` with explicit versions and additional test dependencies (`testify`, `difflib`)
* Improves dependency clarity and reproducibility of builds

### 7. Repository Hygiene

* Added `.gitignore` covering:

  * Go build artifacts
  * Flutter build/cache directories
  * environment files and OS artifacts
* Introduced MIT `LICENSE`, formalizing project usage and distribution rights

### 8. API Project Adjustments

* Updated `api/README.md`:

  * aligned commands with new root Makefile
  * corrected build paths (`api/build/`)
  * updated documentation links to `docs/api/`
* Ensures consistency between documentation and actual project structure

### Conclusion

This commit represents a **structural milestone** rather than a feature addition.

Key impacts:

* Establishes a **clean monorepo foundation**
* Improves **developer ergonomics and workflow consistency**
* Elevates documentation to a **core part of the system design**
* Aligns project organization with its architectural principles

From an engineering perspective, this is a highly valuable refactor that reduces cognitive load, eliminates duplication, and prepares the codebase for scalable evolution across both backend and mobile layers.
