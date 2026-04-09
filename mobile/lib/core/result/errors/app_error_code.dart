/// Domain-oriented error aligned with API contract
enum AppErrorCode {
  // HTTP
  httpError,
  timeout,
  networkError,
  parsingError,
  unauthenticated,

  // Storage
  storageError,
  storageNotFound,
  storageConflict,
  storageCorrupted,
  storageExpired,

  // Generic
  unexpected,
}
