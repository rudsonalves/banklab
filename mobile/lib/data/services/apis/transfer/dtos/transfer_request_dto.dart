import 'package:money2/money2.dart';

import '../../core/api_parse.dart';

class TransferRequestDto {
  final String fromBranch;
  final String fromAccountNumber;
  final String toBranch;
  final String toAccountNumber;
  final Money amount;
  final String idempotencyKey;
  final String? description;

  TransferRequestDto({
    required this.fromBranch,
    required this.fromAccountNumber,
    required this.toBranch,
    required this.toAccountNumber,
    required this.amount,
    required this.idempotencyKey,
    this.description,
  });

  Map<String, dynamic> toMap() {
    return {
      'from_branch': fromBranch,
      'from_account_number': fromAccountNumber,
      'to_branch': toBranch,
      'to_account_number': toAccountNumber,
      'amount': ApiParse.toInt(amount),
      'idempotency_key': idempotencyKey,
      if (description != null) 'description': description,
    };
  }
}
