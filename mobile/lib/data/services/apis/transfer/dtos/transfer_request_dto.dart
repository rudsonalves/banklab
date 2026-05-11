import 'package:money2/money2.dart';

import '../../core/api_parse.dart';

class TransferRequestDto {
  final String fromAccountId;
  final String toAccountId;
  final Money amount;
  final String idempotencyKey;
  final String? description;

  TransferRequestDto({
    required this.fromAccountId,
    required this.toAccountId,
    required this.amount,
    required this.idempotencyKey,
    this.description,
  });

  Map<String, dynamic> toMap() {
    return {
      'from_account_id': fromAccountId,
      'to_account_id': toAccountId,
      'amount': ApiParse.toInt(amount),
      'idempotency_key': idempotencyKey,
      if (description != null) 'description': description,
    };
  }
}
