# Changelog

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
