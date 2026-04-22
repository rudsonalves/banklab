class StatementItemDto {
  final String transactionId;
  final String type;
  final int amount;
  final int balanceAfter;
  final String? referenceId;
  final String createdAt;

  StatementItemDto({
    required this.transactionId,
    required this.type,
    required this.amount,
    required this.balanceAfter,
    required this.referenceId,
    required this.createdAt,
  });

  factory StatementItemDto.fromMap(Map<String, dynamic> map) {
    return StatementItemDto(
      transactionId: map['transaction_id'] as String,
      type: map['type'] as String,
      amount: map['amount'] as int,
      balanceAfter: map['balance_after'] as int,
      referenceId: map['reference_id'] as String?,
      createdAt: map['created_at'] as String,
    );
  }
}

class StatementNextCursorDto {
  final String createdAt;
  final String id;

  StatementNextCursorDto({
    required this.createdAt,
    required this.id,
  });

  factory StatementNextCursorDto.fromMap(Map<String, dynamic> map) {
    return StatementNextCursorDto(
      createdAt: map['created_at'] as String,
      id: map['id'] as String,
    );
  }
}

class StatementResponseDto {
  final String accountId;
  final List<StatementItemDto> items;
  final StatementNextCursorDto? nextCursor;

  StatementResponseDto({
    required this.accountId,
    required this.items,
    required this.nextCursor,
  });

  factory StatementResponseDto.fromMap(Map<String, dynamic> map) {
    final rawItems = map['items'] as List<dynamic>? ?? const [];

    return StatementResponseDto(
      accountId: map['account_id'] as String,
      items: rawItems
          .map((item) => StatementItemDto.fromMap(item as Map<String, dynamic>))
          .toList(),
      nextCursor: map['next_cursor'] != null
          ? StatementNextCursorDto.fromMap(
              map['next_cursor'] as Map<String, dynamic>,
            )
          : null,
    );
  }
}
