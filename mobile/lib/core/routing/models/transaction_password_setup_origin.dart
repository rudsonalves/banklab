enum TransactionPasswordSetupOrigin {
  postLogin,
  transfer;

  factory TransactionPasswordSetupOrigin.fromName(String originName) {
    return TransactionPasswordSetupOrigin.values.firstWhere(
      (origin) => origin.name == originName,
      orElse: () => throw UnsupportedError(
        'Unsupported TransactionPasswordSetupOrigin: $originName',
      ),
    );
  }
}
