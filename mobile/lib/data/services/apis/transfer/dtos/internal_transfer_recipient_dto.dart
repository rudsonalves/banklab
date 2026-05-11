class InternalTransferRecipientDto {
  final String accountId;
  final String holderName;
  final String document;
  final String branch;
  final String accountNumber;

  const InternalTransferRecipientDto({
    required this.accountId,
    required this.holderName,
    required this.document,
    required this.branch,
    required this.accountNumber,
  });

  factory InternalTransferRecipientDto.fromMap(Map<String, dynamic> map) {
    return InternalTransferRecipientDto(
      accountId: map['account_id'] as String,
      holderName: map['holder_name'] as String,
      document: map['document'] as String,
      branch: map['branch'] as String,
      accountNumber: map['account_number'] as String,
    );
  }
}
