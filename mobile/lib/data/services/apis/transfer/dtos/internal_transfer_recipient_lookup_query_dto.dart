class InternalTransferRecipientLookupQueryDto {
  final String? branch;
  final String? accountNumber;
  final String? document;

  const InternalTransferRecipientLookupQueryDto._({
    this.branch,
    this.accountNumber,
    this.document,
  });

  factory InternalTransferRecipientLookupQueryDto.byAccount({
    required String branch,
    required String accountNumber,
  }) {
    return InternalTransferRecipientLookupQueryDto._(
      branch: branch,
      accountNumber: accountNumber,
    );
  }

  factory InternalTransferRecipientLookupQueryDto.byCpf({
    required String cpf,
  }) {
    return InternalTransferRecipientLookupQueryDto._(document: cpf);
  }

  Map<String, dynamic> toMap() {
    return {
      if (branch != null) 'branch': branch,
      if (accountNumber != null) 'account_number': accountNumber,
      if (document != null) 'document': document,
    };
  }
}
