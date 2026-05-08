import 'package:money2/money2.dart';

import '../../core/api_parse.dart';

class TransferResponseDto {
  final String fromAccountId;
  final String toAccountId;
  final String transactionReference;
  final String toBranch;
  final Money amount;
  final Money fromBalance;
  final Money toBalance;

  TransferResponseDto({
    required this.fromAccountId,
    required this.toAccountId,
    required this.transactionReference,
    required this.toBranch,
    required this.amount,
    required this.fromBalance,
    required this.toBalance,
  });

  factory TransferResponseDto.fromMap(Map<String, dynamic> map) {
    return TransferResponseDto(
      fromAccountId: map['from_account_id'],
      toAccountId: map['to_account_id'],
      transactionReference: map['transaction_reference'],
      toBranch: map['to_branch'],
      amount: ApiParse.toMoney(map['amount']),
      fromBalance: ApiParse.toMoney(map['from_balance']),
      toBalance: ApiParse.toMoney(map['to_balance']),
    );
  }
}
