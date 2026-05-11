# Changelog

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
