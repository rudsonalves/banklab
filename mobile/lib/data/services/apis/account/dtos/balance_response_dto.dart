import 'package:money2/money2.dart';

import '../../core/api_parse.dart';

class BalanceResponseDto {
  final String accountId;
  final Money balance;

  BalanceResponseDto({
    required this.accountId,
    required this.balance,
  });

  factory BalanceResponseDto.fromMap(Map<String, dynamic> map) {
    return BalanceResponseDto(
      accountId: map['account_id'] as String,
      balance: ApiParse.toMoney(map['balance']),
    );
  }
}
