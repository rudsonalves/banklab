class RecipientInfoDto {
  final String accountId;
  final String holderName;
  final String document;
  final String branch;
  final String accountNumber;

  RecipientInfoDto({
    required this.accountId,
    required this.holderName,
    required this.document,
    required this.branch,
    required this.accountNumber,
  });

  factory RecipientInfoDto.fromMap(Map<String, dynamic> map) {
    return RecipientInfoDto(
      accountId: map['account_id'],
      holderName: map['holder_name'],
      document: map['document'],
      branch: map['branch'],
      accountNumber: map['account_number'],
    );
  }
}
