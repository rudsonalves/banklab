class RecipientRequestDto {
  final String? branch;
  final String? accountNumber;
  final String? document;

  const RecipientRequestDto({
    this.branch,
    this.accountNumber,
    this.document,
  });

  bool get isEmpty =>
      document == null && (branch == null || accountNumber == null);

  Map<String, dynamic> toMap() {
    if (document != null) {
      return {'document': document};
    }

    if (branch != null && accountNumber != null) {
      return {
        'branch': branch,
        'account_number': accountNumber,
      };
    }

    return {};
  }
}
