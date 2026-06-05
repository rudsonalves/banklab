import '/core/result/command.dart';
import '/data/repositories/transaction_password/transaction_password_repository.dart';
import '/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';

class TransactionPasswordViewModel {
  TransactionPasswordViewModel({
    required TransactionPasswordRepository repository,
  }) {
    create = Command1(repository.create);
  }

  late final Command1<
    TransactionPasswordStatusResponseDto,
    CreateTransactionPasswordRequestDto
  >
  create;
}
