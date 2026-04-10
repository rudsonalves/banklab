# Changelog

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


## 2026/04/08 — doc/adjustments-03

Introduces the **initial authentication and profile integration on the client side**, along with API standardization, domain modeling, and utility improvements. The changes align the Flutter application with the backend REST contract and reinforce consistency in error handling and data parsing.

### 1. Core Utilities and Error Handling

* Added `DateTimeExtensions`:

  * localized date formatting using `intl`
  * safe parsing via `parseOrNull`
* Introduced new error code:

  * `parsingError` to explicitly represent serialization/deserialization failures
* This is a relevant improvement, as it separates **transport errors from data integrity issues**, increasing observability and debuggability

### 2. API Layer — Authentication Module

* Implemented `AuthApi` with endpoints:

  * `POST /auth/register`
  * `POST /auth/login`
  * `GET /profile/me`
* Standardized API consumption using `ApiEnvelope<T>`:

  * consistent handling of `data` and `error` fields
  * explicit validation of null payloads
* Introduced structured error mapping:

  * HTTP-level errors → `AppErrorCode.httpError`
  * parsing failures → `AppErrorCode.parsingError`
* This design is strongly aligned with the backend contract defined in , ensuring consistency between client and server

### 3. DTOs and Serialization

* Added request/response DTOs:

  * `LoginRequestDto`
  * `RegisterRequestDto`
  * `RegisterResponseDto`
* Clear separation between transport models and domain models
* Mapping strategy ensures:

  * strong typing
  * controlled transformation boundaries

### 4. API Envelope Standardization

* Introduced:

  * `ApiEnvelope<T>`
  * `ApiError`
* Centralizes response parsing logic according to backend specification
* Eliminates duplication across API calls and enforces a single response contract

### 5. Domain Layer — Authentication Models

* Introduced domain entities:

  * `AuthUser` (sealed abstraction)
  * `LoggedUser`
  * `NotLoggedUser`
  * `UserProfile`
* Added `UserRole` enum with safe parsing (`byName`)
* Domain modeling reflects backend payload structure and invariants
* Notably, `UserProfile` integrates date parsing via shared extension, improving consistency

### 6. UI Layer — Registration Scaffold

* Added initial structure for:

  * `RegisterPage`
  * `RegisterViewmodel`
* Current implementation is a placeholder, but establishes:

  * separation between view and state logic
  * preparation for MVVM-style composition

### 7. Dependency Management

* Promoted `uuid` to a direct dependency
* Suggests upcoming usage for:

  * client-side identifiers
  * correlation or request tracing

### 8. Test Adjustments

* Minor refactor in Dio test adapter signatures:

  * removed unused parameters (`__`)
* Improves code clarity and consistency with Dart conventions

### 9. Architectural Alignment

The changes respect the layered architecture principles described in :

* API layer acts as **infrastructure integration**
* domain models remain isolated from transport concerns
* UI layer depends on abstractions rather than concrete implementations

This reinforces:

* low coupling
* clear separation of concerns
* improved testability

### Conclusion

This commit establishes a **solid foundation for authentication on the client**, with emphasis on:

* standardized API communication via envelope pattern
* explicit and granular error handling
* clear separation between DTOs and domain models
* preparation for scalable UI architecture

From a technical perspective, the most valuable aspect is the **formalization of the API contract consumption**, which significantly reduces ambiguity and future integration errors.


## 2026/04/08 — feat/http-core-module-02

Refines the initial Flutter client shell to better support the project’s role as a controlled integration surface for the banking API, while also bringing the iOS workspace into a consistent CocoaPods-managed state. This change improves the project presentation, removes demo-oriented UI leftovers, and prepares the mobile layer for a cleaner HTTP core evolution aligned with the backend architecture and REST contract. 

### 1. Documentation and project positioning

* Reworked `README.md` to present **BankFlow** as a validation client for a custom banking API rather than a generic Flutter demo
* Clarified the project purpose around:

  * end-to-end validation of financial workflows
  * client/backend contract consistency
  * integration boundaries
  * transaction safety
* Added a more precise architectural framing for the mobile app as part of the broader system design effort
* Improved the motivation and scope sections to better reflect the engineering focus of the project
* Moved the license section to the end and simplified its wording

### 2. Flutter app bootstrap cleanup

* Simplified `MainApp` by removing the obsolete demo title parameter from the home screen instantiation
* Kept the application entry point leaner and closer to the real project intent, instead of Flutter template defaults

### 3. Home page simplification

* Removed the default counter-based demo behavior from `HomePage`
* Eliminated:

  * `title` parameter
  * internal counter state
  * increment action
  * floating action button
  * demo-specific text rendering
* Replaced the old template structure with a cleaner and more neutral home screen scaffold
* Corrected enum usage from shorthand syntax to explicit `MainAxisAlignment.center`, which is technically clearer and more idiomatic for maintainable code

### 4. iOS CocoaPods and workspace integration

* Added `ios/Podfile.lock` to capture the current CocoaPods dependency state
* Updated `Runner.xcodeproj` to include:

  * Pods framework references
  * Pods xcconfig references
  * manifest lock check phases
  * embed frameworks phase
  * test target Pods integration
* Updated `Runner.xcworkspace` to include `Pods/Pods.xcodeproj`
* These changes indicate the iOS side is now aligned with a proper CocoaPods-managed workspace structure, which is important for plugin stability and deterministic local builds

### 5. Architectural significance

* This commit is small in visible feature scope, but strategically relevant
* It removes noise inherited from the Flutter template and makes the mobile app more compatible with its intended role inside a layered system, where the client should act as a predictable integration surface rather than a sandbox demo
* That direction is coherent with the backend architectural model centered on layered boundaries and controlled execution flow, as documented in the project architecture and REST contract references  

### Conclusion

This commit cleans the client foundation, improves project documentation, and stabilizes the iOS workspace configuration. The result is a more intentional mobile base, better aligned with the project’s real objective: serving as a structured frontend for validating HTTP flows and backend behavior rather than remaining tied to Flutter’s default starter template.


## 2026/04/08 — feat/http-core-module-01

Introduces a **core HTTP module standardization** with unified error modeling, improved result handling, and a more robust command execution pattern. This refactor significantly improves consistency, testability, and alignment with layered architecture principles 

### 1. Unified Error Model (AppError)

* Introduced `AppError` and `AppErrorCode` as the **single error abstraction across the core layer**
* Replaced legacy exception hierarchy:

  * removed `BaseException`, `GenericException`, `StorageException`, and `RestClientException`
* Standardized error structure:

  * `code` (domain-oriented)
  * `message` (human-readable)
  * `details` (raw context)
  * optional `statusCode`
* Aligns error handling with API contract (`data` / `error` envelope) 

### 2. Result Pattern Refinement

* Updated `Result` to use `AppError` instead of generic `Exception`
* Improved type safety and consistency across all layers
* Added exports and modularization (`unit.dart`, error modules)
* Ensures **explicit success/failure modeling**, avoiding exception-driven flow

### 3. Command Pattern Evolution

* Introduced `CommandState` enum:

  * `idle`, `running`, `success`, `failure`
* Added execution control:

  * `_executionId` to prevent race conditions
  * guards against concurrent execution
* Improved API:

  * `data` and `error` accessors
  * state-driven UI integration
* Added **infra-agnostic fallback error handling** using `AppError.unexpected`
* This is a **notable design improvement**, especially for Flutter state management

### 4. HTTP Layer Improvements (Dio Integration)

* Introduced `dio_error_mapper.dart`:

  * maps `DioException` → `AppError`
  * handles:

    * timeout
    * network errors
    * structured API errors
    * fallback cases
* Refactored `DioRestClient`:

  * removed exception-based flow
  * now returns `Result.success` / `Result.failure`
* Ensures HTTP layer is:

  * predictable
  * strongly typed
  * aligned with core error model

### 5. Storage Layer Refactor

* Replaced exception-based storage errors with `AppError` via extension:

  * `StorageAppError.storage`
  * `StorageAppError.notFound`
* Updated `FlutterSecureStorageLocalStorage`:

  * all operations now return `Result`
  * consistent error mapping and logging
* Improves reliability and removes hidden exception propagation

### 6. REST Client Enhancements

* Added `parse<T>` helper in `RestClientResponse`

  * simplifies DTO transformation
* Cleaned exports:

  * removed deprecated exception exposure
  * added error mapper export
* Reinforces clean contract for HTTP consumers

### 7. Test Infrastructure and Coverage

* Added comprehensive unit tests for:

  * `AppError` behavior and formatting
  * `RestClientRequest` (copy semantics)
  * `RestClientResponse` (status + copy)
  * `DioRestClient` (success + error mapping)
  * `FlutterSecureStorageLocalStorage` (success + failure scenarios)
* Introduced HTTP adapter mock for deterministic testing
* Strengthens confidence in core module behavior

### 8. Makefile Improvements

* Added test commands:

  * `make test`
  * `make test-unit`
* Encourages standardized execution of test suites

### Conclusion

This commit establishes a **solid foundation for the HTTP core module**, with emphasis on:

* **unified error handling (AppError)**
* **predictable result-based flow**
* **elimination of exception-driven control paths**
* **improved concurrency safety in command execution**

From an architectural standpoint, this is a **high-impact refactor**, bringing the core layer closer to a clean, deterministic, and framework-agnostic design, fully aligned with modular layered principles.


## 2026/04/07 — feat/http-core-module-01

Introduces a **foundational HTTP core module** for the Flutter client, including a structured REST client abstraction, secure storage integration, interceptor pipeline, and unified error/result handling. This establishes a solid infrastructure layer aligned with a modular and scalable architecture.

### 1. HTTP Client Abstraction

* Introduced `RestClient` contract defining standard HTTP operations (GET, POST, PUT, PATCH, DELETE)
* Added request/response models:

  * `RestClientRequest`
  * `RestClientResponse`
* Implemented `RestClientException` with semantic helpers (`isUnauthorized`, `isForbidden`, etc.)
* Designed for **framework-agnostic usage**, decoupling higher layers from HTTP implementation details
* Aligns with layered architecture principles where infrastructure details are isolated 

### 2. Dio-Based Implementation

* Implemented `DioRestClient` as concrete adapter for `RestClient`
* Centralized request execution with:

  * unified success mapping → `Success<Result>`
  * error normalization → `Failure<RestClientException>`
* Added `DioFactory`:

  * configurable base URL, timeouts, and interceptors
  * prevents duplicate interceptor registration
* This design provides **high flexibility and testability**, enabling easy replacement of HTTP engine if needed

### 3. Result Pattern Refactor

* Refactored `Result<T>`:

  * introduced `sealed class` with `Success` and `Failure`
  * added pattern matching via `switch`
  * improved `fold` implementation
* Eliminated nullable state ambiguity (`value`/`error`)
* This is a **significant improvement in API ergonomics and correctness**, especially for async flows

### 4. Secure Storage Layer

* Introduced `LocalSecureStorage` contract
* Implemented `FlutterSecureStorageLocalStorage`:

  * supports read, write, delete, deleteAll, and prefix filtering
  * returns `Result` instead of throwing exceptions
* Added structured storage keys (`StorageKeys`)
* Introduced storage-specific exceptions:

  * `StorageNotFoundException`
  * `StorageConflictException`
  * `StorageCorruptedException`, etc.
* Provides a **robust and explicit error model**, avoiding silent failures in critical flows

### 5. Interceptor Pipeline

* Implemented `AuthInterceptor`:

  * automatically injects `Authorization: Bearer <token>`
  * handles `401` responses with token refresh flow
  * retries original request after successful refresh
  * clears session on failure
* Explicitly documents **race condition risk** during concurrent refresh attempts (design trade-off acknowledged)
* Prepared extensibility for additional interceptors (device, security), currently scaffolded
* This is a **production-grade foundation**, though future synchronization for refresh is recommended

### 6. Exception Hierarchy

* Introduced base exception models:

  * `BaseException`
  * `GenericException`
* Standardized exception structure across modules
* Enables consistent error propagation and logging strategy

### 7. Project Tooling (Makefile)

* Added developer-oriented Makefile with commands:

  * `commit` (uses predefined message file)
  * `diff` (staged diff + line count)
  * `push` / `pull` with dynamic branch resolution
  * `gitlog`
* Improves workflow consistency and reduces manual errors

### 8. iOS Configuration

* Added `Podfile` and updated `.xcconfig` files
* Ensures proper integration of CocoaPods and Flutter plugins
* Required for `flutter_secure_storage` compatibility on iOS

### 9. Dependencies

* Added `flutter_secure_storage`
* Updated transitive dependencies accordingly
* Enables secure persistence of sensitive data (tokens)

### 10. Architectural Impact

This commit establishes the **entire HTTP infrastructure layer** for the client, with clear separation of concerns:

* HTTP communication isolated behind `RestClient`
* Security concerns handled via interceptors
* Persistence handled via secure storage abstraction
* Error handling standardized via `Result` and exception hierarchy

This design is strongly aligned with the principles described in the system architecture, particularly:

* isolation of infrastructure concerns
* explicit contracts between layers
* improved testability and replaceability 

Additionally, the module is consistent with the REST contract expectations (JWT, error envelope, status handling), enabling seamless integration with the backend API 

### Conclusion

This is a **foundational and high-impact commit**. It does not deliver end-user features directly, but it defines the technical backbone required for all future network communication, authentication flows, and secure data handling.

The design is clean, extensible, and largely production-ready. The only notable gap is the lack of synchronization in token refresh, which should be addressed as concurrency increases.
