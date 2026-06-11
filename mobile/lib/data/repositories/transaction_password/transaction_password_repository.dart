import '/core/result/result.dart';
import '/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/set_up_authorize_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/set_up_authorize_response_dto.dart';
import '/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';

abstract class TransactionPasswordRepository {
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  );

  AsyncResult<SetUpAuthorizeResponseDto> stepUpAuthorize(
    SetUpAuthorizeRequestDto dto,
  );
}
