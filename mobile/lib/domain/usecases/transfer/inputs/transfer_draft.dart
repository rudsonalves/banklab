import 'package:money2/money2.dart';

class TransferDraft {
  final String toAccountNumber;
  final String toBranch;
  final String? description;
  final Money amount;
  final String idempotencyKey;

  TransferDraft({
    required this.toAccountNumber,
    required this.toBranch,
    this.description,
    required this.amount,
    this.idempotencyKey = '',
  });

  TransferDraft copyWith({
    String? toAccountNumber,
    String? toBranch,
    String? description,
    Money? amount,
    String? idempotencyKey,
  }) {
    return TransferDraft(
      toAccountNumber: toAccountNumber ?? this.toAccountNumber,
      toBranch: toBranch ?? this.toBranch,
      description: description ?? this.description,
      amount: amount ?? this.amount,
      idempotencyKey: idempotencyKey ?? this.idempotencyKey,
    );
  }
}
