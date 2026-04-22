class AccountSummaryResponseDto {
  final String id;
  final String customerId;
  final String number;
  final String branch;
  final String status;

  AccountSummaryResponseDto({
    required this.id,
    required this.customerId,
    required this.number,
    required this.branch,
    required this.status,
  });

  factory AccountSummaryResponseDto.fromMap(Map<String, dynamic> map) {
    return AccountSummaryResponseDto(
      id: map['id'] as String,
      customerId: map['customer_id'] as String,
      number: map['number'] as String,
      branch: map['branch'] as String,
      status: map['status'] as String,
    );
  }
}
