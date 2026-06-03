/// Domain-oriented error aligned with API contract
enum AppErrorCode {
  // HTTP
  httpError,
  accountApprovalRequired,
  contactNotVerified,
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

  // Registration-specific
  cpfAlreadyRegistered,

  // Transaction password-specific
  transactionPasswordAlreadySet,
}
