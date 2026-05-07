import 'package:money2/money2.dart';

import '../../core/api_parse.dart';

class TransferResponseDto {
  final String fromBranch;
  final String fromAccountNumber;
  final String transactionReference;
  final String toBranch;
  final String toAccountNumber;
  final Money amount;
  final Money fromBalance;
  final Money toBalance;

  TransferResponseDto({
    required this.fromBranch,
    required this.fromAccountNumber,
    required this.transactionReference,
    required this.toBranch,
    required this.toAccountNumber,
    required this.amount,
    required this.fromBalance,
    required this.toBalance,
  });

  factory TransferResponseDto.fromMap(Map<String, dynamic> map) {
    return TransferResponseDto(
      fromBranch: map['from_branch'],
      fromAccountNumber: map['from_account_number'],
      transactionReference: map['transaction_reference'],
      toBranch: map['to_branch'],
      toAccountNumber: map['to_account_number'],
      amount: ApiParse.toMoney(map['amount']),
      fromBalance: ApiParse.toMoney(map['from_balance']),
      toBalance: ApiParse.toMoney(map['to_balance']),
    );
  }
}
