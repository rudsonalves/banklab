# Mobile Implemented Features

This document summarizes what the BankFlow mobile app has implemented so far.
It is a living status view of the current mobile client, not a changelog.

## Overview

The mobile app is a Flutter client organized around a layered architecture:

- UI layer for screens, view models, and shared widgets
- Domain layer for stable app-facing models and enums
- Data layer for repositories, API services, DTOs, and local persistence
- Core layer for Result, Command, routing, HTTP, secure storage, and app config

Current flow model:

`UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio`

When a flow needs reusable orchestration across multiple repositories, the app is
prepared to use a domain use case. Current transfer and details flows already
use this pattern, while simpler flows remain repository-driven.

## Authentication

Authentication is implemented with the following flows:

- full login with e-mail and password
- short login using the remembered identity
- splash-based auth entry decision
- registration
- repository-level logout/session clearing
- profile loading after successful login

The current auth stack includes:

- [AuthRepository](../lib/data/repositories/auth/auth_repository.dart)
- [AuthRepositoryImpl](../lib/data/repositories/auth/auth_repository_impl.dart)
- [RegistrationRepository](../lib/data/repositories/registration/registration_repository.dart)
- [ContactVerificationRepository](../lib/data/repositories/contact_verification/contact_verification_repository.dart)
- [RegisterDraftRepository](../lib/data/repositories/register_draft/register_draft_repository.dart)
- [AuthApi](../lib/data/services/apis/auth/auth_api.dart)
- [RegistrationApi](../lib/data/services/apis/registration/registration_api.dart)
- [ContactVerificationApi](../lib/data/services/apis/contact_verification/contact_verification_api.dart)
- [RegisterUsecase](../lib/domain/usecases/register/register_usecase.dart)
- [LoginPage](../lib/ui/pages/auth/login/login_page.dart)
- [ShortLoginPage](../lib/ui/pages/auth/short_login/short_login_page.dart)
- [RegisterCpfPage](../lib/ui/pages/register/register_cpf_page.dart)
- [RegisterNamePage](../lib/ui/pages/register/register_name_page.dart)
- [RegisterBirthdatePage](../lib/ui/pages/register/register_birthdate_page.dart)
- [RegisterEmailPage](../lib/ui/pages/register/register_email_page.dart)
- [RegisterTokenPage](../lib/ui/pages/register/register_token_page.dart)
- [RegisterPhonePage](../lib/ui/pages/register/register_phone_page.dart)
- [RegisterPasswordPage](../lib/ui/pages/register/register_password_page.dart)
- [RegisterStatusPage](../lib/ui/pages/register/register_status_page.dart)
- [SplashPage](../lib/ui/pages/splash/splash_page.dart)
- [LoginViewModel](../lib/ui/pages/auth/login/viewmodel/login_viewmodel.dart)
- [ShortLoginViewModel](../lib/ui/pages/auth/short_login/viewmodel/short_login_viewmodel.dart)
- [RegisterViewmodel](../lib/ui/pages/register/viewmodel/register_viewmodel.dart)
- [SplashViewModel](../lib/ui/pages/splash/viewmodel/splash_viewmodel.dart)

Implemented auth behavior:

- access and refresh tokens are stored in secure storage after a successful login
- the user profile is loaded after login succeeds
- the last login identity is cached after profile loading succeeds
- splash uses the remembered identity cache to route to full login or short
  login
- short login reuses the remembered identity and lets the user switch accounts
- repository-level logout clears the persisted auth tokens
- approval-required login failures are handled as a distinct app error state
- invalid credentials remain a separate failure path
- registration now runs as a multi-page onboarding flow under
  [ui/pages/register](../lib/ui/pages/register)
- the onboarding persists draft progress between steps through the register
  draft repository
- e-mail and phone contact verification are requested and confirmed before the
  final registration submission
- the flow ends in a success or failure status page after the final submit

Current login feedback uses `AppSnackbar` for transient messages.
Approval-pending login now shows a specific user-facing message instead of a
wrong-password style error.

Registration flow entry points include:

- [RegisterCpfPage](../lib/ui/pages/register/register_cpf_page.dart)
- [RegisterNamePage](../lib/ui/pages/register/register_name_page.dart)
- [RegisterBirthdatePage](../lib/ui/pages/register/register_birthdate_page.dart)
- [RegisterEmailPage](../lib/ui/pages/register/register_email_page.dart)
- [RegisterTokenPage](../lib/ui/pages/register/register_token_page.dart)
- [RegisterPhonePage](../lib/ui/pages/register/register_phone_page.dart)
- [RegisterPasswordPage](../lib/ui/pages/register/register_password_page.dart)
- [RegisterStatusPage](../lib/ui/pages/register/register_status_page.dart)
- [RegisterViewmodel](../lib/ui/pages/register/viewmodel/register_viewmodel.dart)

## Account And Statement Features

The mobile app already implements the main account views and related data flows:

- account list and selected account handling
- account balance loading and refresh
- bank statement listing
- statement entry point from the home screen
- statement empty-state and error-state handling
- statement grouping by month and day
- transaction detail pages for statement items
- statement state preservation through `AccountRepository.lastStatement`

Relevant files include:

- [AccountRepositoryImpl](../lib/data/repositories/account/account_repository_impl.dart)
- [BalanceApi](../lib/data/services/apis/account/balance_api.dart)
- [ListAccountsApi](../lib/data/services/apis/account/list_accounts_api.dart)
- [StatementApi](../lib/data/services/apis/account/statement_api.dart)
- [HomePage](../lib/ui/pages/home/home_page.dart)
- [HomeViewModel](../lib/ui/pages/home/viewmodel/home_viewmodel.dart)
- [StatementPage](../lib/ui/pages/statement/statement_page.dart)
- [StatementViewModel](../lib/ui/pages/statement/viewmodel/statement_viewmodel.dart)
- [DetailsPage](../lib/ui/pages/shared/details/details_page.dart)
- [DetailsViewModel](../lib/ui/pages/shared/details/viewmodel/details_viewmodel.dart)

Implemented account behavior:

- the app can load and display the user account list
- the app can show the current account balance
- the app can load statement entries for the selected account
- statement views support loading, empty, and error states
- statement items can open the shared details page using `referenceId`
- statement labels and colors use the shared transaction movement mapping

## Transfer Features

Transfer-related functionality is implemented at the client level:

- internal transfer orchestration
- recipient lookup for internal transfers
- transfer confirmation flow
- transfer payment flow
- transfer result/status flow
- receipt visualization after transfer

Relevant files include:

- [TransactionRepositoryImpl](../lib/data/repositories/transaction/transaction_repository_impl.dart)
- [ApiTransfer](../lib/data/services/apis/transfer/api_transfer.dart)
- [ApiReceipt](../lib/data/services/apis/receipt/api_receipt.dart)
- [TransferPage](../lib/ui/pages/home/transfer/transfer_recipient_page.dart)
- [TransferConfirmationPage](../lib/ui/pages/home/transfer/transfer_confirmation_page.dart)
- [TransferPaymentPage](../lib/ui/pages/home/transfer/transfer_payment_page.dart)
- [TransferStatusPage](../lib/ui/pages/home/transfer/transfer_status_page.dart)
- [TransferViewModel](../lib/ui/pages/home/transfer/viewmodel/transfer_viewmodel.dart)

Implemented transfer behavior:

- the app can search and select an internal recipient
- the app can confirm a transfer before submitting it
- the app can submit the transfer request
- the app can show the final transfer status
- the app can retrieve and display transfer receipts
- transfer and receipt flows are coordinated through domain use cases where
  multi-repository orchestration is needed

API contract note:

- the backend now requires `POST /accounts/internal-transfers` to include
  `X-Step-Up-Token`
- the token must be issued by `POST /security/step-up/authorize` for
  `internal_transfer.create`
- the token is single-use; after `STEP_UP_TOKEN_CONSUMED`, the client must
  request a new step-up token before retrying the transfer
- retrying with the same `idempotency_key` may still require a new step-up token

## Routing And Navigation

Routing is handled with GoRouter and typed route groups.

Implemented navigation areas include:

- splash route
- auth routes
- statement route
- home routes
- transfer routes
- shared details routes

Relevant files:

- [router](../lib/core/routing/router.dart)
- [routes](../lib/core/routing/routes.dart)
- [auth routes](../lib/core/routing/routes/auth_routes.dart)
- [home routes](../lib/core/routing/routes/home_routes.dart)
- [transfer routes](../lib/core/routing/routes/transfer_routes.dart)

## State, Errors, And Feedback

The mobile app uses the shared core result model:

- `Result<T>` for success/failure flows
- `AppError` for typed application errors
- `Command` for async UI-triggered actions

Implemented error handling traits:

- transport failures are converted into `AppError`
- parsing failures become `AppErrorCode.parsingError`
- HTTP failures stay explicit through `AppErrorCode.httpError`
- login approval-required is represented by `AppErrorCode.accountApprovalRequired`

Implemented feedback traits:

- pages react to command state changes
- `AppSnackbar` is used for transient user messages
- loading states are driven by command execution state

Shared UI primitives already in use include:

- `SafeScaffold`
- `AppSnackbar`
- `BigButton`
- `BasicTextFormField`
- account, balance, and recipient cards
- `TransactionMovement` for transaction labels, signs, and colors

These primitives live under `mobile/lib/ui/components`, which is also the staging
area for a future internal mobile widget/feature package. Components should
only move there when they are reusable, presentation-only, and not tied to a
single page workflow.

## HTTP And Session Infrastructure

The mobile app has a reusable HTTP layer around Dio:

- `RestClient` abstraction
- `DioRestClient` implementation
- request/response wrappers
- Dio error mapping into `AppError`
- auth interceptor for JWT handling

Implemented session behavior:

- authenticated requests receive the stored access token automatically
- refresh is attempted after `401` responses when appropriate
- concurrent refresh attempts are coordinated by the interceptor
- failed refresh clears persisted auth tokens

## Persistence And Session State

The app currently persists and restores session-related data:

- access token
- refresh token
- last login identity
- current user session state in the auth repository

Implemented storage helpers include:

- [LocalSecureStorage](../lib/core/services/secure_storage/local_secure_storage.dart)
- [FlutterSecureStorageLocalStorage](../lib/core/services/secure_storage/flutter_secure_storage_local_storage.dart)
- [LastLoginCacheService](../lib/data/services/cache/last_login/last_login_cache_service.dart)

## Testing Coverage

The current mobile test suite covers:

- HTTP client behavior
- request/response/exception wrappers
- Dio REST client behavior
- auth interceptor behavior
- auth approval-required error mapping
- repository login side effects
- full login and short login feedback behavior
- secure storage adapter behavior
- string and date extensions
- transfer API behavior
- transfer request/response DTOs
- internal transfer recipient lookup DTOs
- receipt API behavior
- receipt response DTO parsing
- transfer repository behavior
- transfer use case behavior
- contact verification API behavior
- registration draft repository behavior
- registration use case orchestration and validation behavior

The newly added approval-required handling is covered by focused tests in:

- [dio_error_mapper_test.dart](../test/core/services/client_http/dio/dio_error_mapper_test.dart)
- [auth_repository_impl_test.dart](../test/data/repositories/auth/auth_repository_impl_test.dart)
- [login_feedback_behavior_test.dart](../test/ui/pages/auth/login_feedback_behavior_test.dart)

## Notes

This document reflects the current state of the mobile client. It should be
updated when new screens, repositories, or auth behaviors are introduced.
