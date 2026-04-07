# Changelog

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
