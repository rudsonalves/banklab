# Changelog

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
