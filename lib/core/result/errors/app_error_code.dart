/// Domain-oriented error aligned with API contract
enum AppErrorCode {
  // HTTP
  httpError,
  timeout,
  networkError,

  // Storage
  storageError,
  storageNotFound,
  storageConflict,
  storageCorrupted,
  storageExpired,

  // Generic
  unexpected,
}
