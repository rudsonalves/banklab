enum TransactionPasswordStatus {
  active,
  blocked;

  factory TransactionPasswordStatus.byName(String value) {
    switch (value) {
      case 'active':
        return TransactionPasswordStatus.active;
      case 'blocked':
        return TransactionPasswordStatus.blocked;
      default:
        throw ArgumentError('Unknown transaction password status: $value');
    }
  }
}
