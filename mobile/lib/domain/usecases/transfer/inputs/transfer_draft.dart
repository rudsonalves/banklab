import 'package:money2/money2.dart';

class TransferDraft {
  final String toAccountId;
  final String? description;
  final Money amount;

  /// A uuid.v7 string to ensure idempotency of transfer requests.
  final String idempotencyKey;

  TransferDraft({
    required this.toAccountId,
    this.description,
    required this.amount,
    required this.idempotencyKey,
  });

  TransferDraft copyWith({
    String? toAccountId,
    String? description,
    Money? amount,
    String? idempotencyKey,
  }) {
    return TransferDraft(
      toAccountId: toAccountId ?? this.toAccountId,
      description: description ?? this.description,
      amount: amount ?? this.amount,
      idempotencyKey: idempotencyKey ?? this.idempotencyKey,
    );
  }
}
