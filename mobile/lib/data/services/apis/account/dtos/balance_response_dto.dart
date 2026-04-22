class BalanceResponseDto {
  final String accountId;
  final int balance;

  BalanceResponseDto({
    required this.accountId,
    required this.balance,
  });

  factory BalanceResponseDto.fromMap(Map<String, dynamic> map) {
    return BalanceResponseDto(
      accountId: map['account_id'] as String,
      balance: map['balance'] as int,
    );
  }
}
