/// Domain-oriented error aligned with API contract
enum AppErrorCode {
  // HTTP
  httpError,
  accountApprovalRequired,
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

  // Validation
  invalidData,
}
