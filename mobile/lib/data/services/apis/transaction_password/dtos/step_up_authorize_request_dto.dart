import '../enums/step_up_operation.dart';

class StepUpAuthorizeRequestDto {
  final StepUpOperation operation;
  final String transactionPassword;

  StepUpAuthorizeRequestDto({
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
