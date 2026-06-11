import '/core/result/result.dart';
import '/core/services/app_section/app_section.dart';
import '/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/set_up_authorize_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/set_up_authorize_response_dto.dart';
import '/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import '/data/services/apis/transaction_password/enums/transaction_password_status.dart';
import '/data/services/apis/transaction_password/transaction_password_api.dart';
import 'transaction_password_repository.dart';

class TransactionPasswordRepositoryImpl
    implements TransactionPasswordRepository {
  final TransactionPasswordApi _api;
  final AppSection _appSection;

  TransactionPasswordRepositoryImpl({
    required TransactionPasswordApi api,
    required AppSection appSection,
  }) : _api = api,
       _appSection = appSection;

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) async {
    final result = await _api.create(dto);
    if (result.isFailure) return Result.failure(result.error!);

    final transPasswdResp = result.value!;
    if (transPasswdResp.status == TransactionPasswordStatus.active) {
      _appSection.markTransactionPasswordAsActive();
    }

    return Success(transPasswdResp);
  }

  @override
  AsyncResult<SetUpAuthorizeResponseDto> stepUpAuthorize(
    SetUpAuthorizeRequestDto dto,
  ) async {
    final result = await _api.stepUpAuthorize(dto);

    if (result.isFailure) {
      final error = result.error!;
      if (error.code == AppErrorCode.transactionPasswordNotSet) {
        _appSection.markTransactionPasswordAsNotSet();
      }
    }

    return result;
  }
}
