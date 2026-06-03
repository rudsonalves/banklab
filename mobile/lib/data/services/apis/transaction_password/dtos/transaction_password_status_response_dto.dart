import '../enums/transaction_password_status.dart';

class TransactionPasswordStatusResponseDto {
  final String userId;
  final TransactionPasswordStatus status;
  final DateTime createdAt;

  TransactionPasswordStatusResponseDto({
    required this.userId,
    required this.status,
    required this.createdAt,
  });

  factory TransactionPasswordStatusResponseDto.fromApi(
    Map<String, dynamic> map,
  ) {
    return TransactionPasswordStatusResponseDto(
      userId: map['user_id'] as String,
      status: TransactionPasswordStatus.byName(map['status'] as String),
      createdAt: DateTime.parse(map['created_at'] as String),
    );
  }
}
