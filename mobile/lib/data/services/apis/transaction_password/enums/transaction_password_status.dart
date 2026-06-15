enum TransactionPasswordStatus {
  active,
  notSet,
  locked,
  unknown;

  factory TransactionPasswordStatus.byName(String name) {
    switch (name) {
      case 'active':
        return TransactionPasswordStatus.active;
      case 'not_set':
        return TransactionPasswordStatus.notSet;
      case 'locked':
        return TransactionPasswordStatus.locked;
      default:
        return TransactionPasswordStatus.unknown;
    }
  }
}
