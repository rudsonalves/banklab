class CreateTransactionPasswordRequestDto {
  final String password;
  final String confirmation;

  CreateTransactionPasswordRequestDto({
    required this.password,
    required this.confirmation,
  });

  Map<String, dynamic> toMap() => {
    'transaction_password': password,
    'transaction_password_confirmation': confirmation,
  };
}
