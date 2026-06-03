import '/core/result/result.dart';
import '/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import '/data/services/apis/transaction_password/transaction_password_api.dart';
import 'transaction_password_repository.dart';

class TransactionPasswordRepositoryImpl
    implements TransactionPasswordRepository {
  final TransactionPasswordApi _api;

  TransactionPasswordRepositoryImpl({
    required TransactionPasswordApi api,
  }) : _api = api;

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) async {
    final result = await _api.create(dto);
    if (result.isFailure) return Result.failure(result.error!);

    return Success(result.value!);
  }
}
