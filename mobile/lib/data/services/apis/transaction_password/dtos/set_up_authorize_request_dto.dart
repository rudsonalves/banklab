import '../enums/step_up_operation.dart';

class SetUpAuthorizeRequestDto {
  final StepUpOperation operation;
  final String transactionPassword;

  SetUpAuthorizeRequestDto({
    required this.operation,
    required this.transactionPassword,
  });

  Map<String, dynamic> toMap() {
    return {
      'method': operation.method.value,
      'path': operation.path,
      'transaction_password': transactionPassword,
    };
  }
}
