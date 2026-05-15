class TransactionMovement {
  final String label;
  final String sign;
  final bool isDebit;

  const TransactionMovement({
    required this.label,
    required this.sign,
    required this.isDebit,
  });

  factory TransactionMovement.fromType(String type) {
    switch (type.toLowerCase()) {
      case 'transfer_out':
        return const TransactionMovement(
          label: 'Débito',
          sign: '-',
          isDebit: true,
        );
      case 'transfer_in':
        return const TransactionMovement(
          label: 'Crédito',
          sign: '+',
          isDebit: false,
        );
      case 'withdraw':
        return const TransactionMovement(
          label: 'Saque',
          sign: '-',
          isDebit: true,
        );
      case 'deposit':
        return const TransactionMovement(
          label: 'Depósito',
          sign: '+',
          isDebit: false,
        );
      default:
        return TransactionMovement(
          label: type,
          sign: '',
          isDebit: false,
        );
    }
  }
}
