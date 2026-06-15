enum TransferReceiptStatus {
  /// Transfer operation completed successfully.
  /// - Money has been debited from source account
  /// - Money has been credited to destination account
  /// - Transaction is immutable and persisted in ledger
  completed('completed'),

  /// Transfer operation is pending processing.
  /// - Request has been received and validated
  /// - Awaiting backend processing or external settlement
  /// - State may transition to completed, failed, or rejected
  pending('pending'),

  /// Transfer operation failed during processing.
  /// - Technical error occurred during execution
  /// - Example: database connectivity, service unavailable
  /// - No funds were transferred; both accounts remain unchanged
  failed('failed'),

  /// Transfer operation was cancelled.
  /// - User or system cancelled the transfer before completion
  /// - No funds were transferred; source account balance unchanged
  cancelled('cancelled'),

  /// Transfer operation was rejected.
  /// - Validation failed or business rules prevented execution
  /// - Examples: insufficient funds, account inactive, same account transfer
  /// - No funds were transferred
  rejected('rejected');

  final String value;

  const TransferReceiptStatus(this.value);

  /// Parses a string value into a TransferReceiptStatus enum.
  ///
  /// Returns the corresponding enum value if found, otherwise throws
  /// a [FormatException] for invalid status values.
  factory TransferReceiptStatus.fromString(String value) {
    return TransferReceiptStatus.values.firstWhere(
      (status) => status.value == value,
      orElse: () => throw FormatException(
        'Unknown transfer receipt status: $value',
      ),
    );
  }

  /// Returns true if the transfer operation succeeded.
  bool get isSuccess => this == completed;

  /// Returns true if the transfer operation is still processing.
  bool get isPending => this == pending;

  /// Returns true if the transfer operation did not succeed.
  bool get isFailed => this == failed || this == rejected || this == cancelled;
}
