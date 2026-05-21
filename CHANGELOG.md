# Changelog

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
* Moved navigation to `context.goNamed(HomeRoutes.home.name)` after successful command completion handling
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
