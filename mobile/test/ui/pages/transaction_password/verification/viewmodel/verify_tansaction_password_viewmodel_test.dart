import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/ui/pages/transaction_password/verification/viewmodel/verify_tansaction_password_viewmodel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('VerifyTansactionPasswordViewmodel.authorizeInternalTransfer', () {
    test('delegates only the transaction password to the repository', () async {
      final response = StepUpAuthorizeResponseDto(
        stepUpToken: 'opaque-step-up-token',
        expiresIn: 120,
      );
      final repository = _FakeTransactionPasswordRepository(
        authorizeResult: Success(response),
      );
      final viewModel = VerifyTansactionPasswordViewmodel(
        repository: repository,
        appSection: AppSection(),
      );

      await viewModel.authorizeInternalTransfer.execute('123456');

      expect(viewModel.authorizeInternalTransfer.isSuccess, isTrue);
      expect(viewModel.authorizeInternalTransfer.value, same(response));
      expect(repository.authorizeCalls, 1);
      expect(repository.lastTransactionPassword, '123456');
    });

    test('exposes the repository backend error unchanged', () async {
      const error = AppError(
        code: AppErrorCode.httpError,
        message: 'Endpoint not allowed',
        details: {'code': 'STEP_UP_ENDPOINT_NOT_ALLOWED'},
      );
      final repository = _FakeTransactionPasswordRepository(
        authorizeResult: const Failure(error),
      );
      final viewModel = VerifyTansactionPasswordViewmodel(
        repository: repository,
        appSection: AppSection(),
      );

      await viewModel.authorizeInternalTransfer.execute('123456');

      expect(viewModel.authorizeInternalTransfer.isFailure, isTrue);
      expect(viewModel.authorizeInternalTransfer.error, same(error));
      expect(
        backendErrorCode(viewModel.authorizeInternalTransfer.error),
        'STEP_UP_ENDPOINT_NOT_ALLOWED',
      );
    });
  });
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  _FakeTransactionPasswordRepository({
    required this.authorizeResult,
  });

  final Result<StepUpAuthorizeResponseDto> authorizeResult;
  String? lastTransactionPassword;
  int authorizeCalls = 0;

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInternalTransfer(
    String transactionPassword,
  ) async {
    authorizeCalls++;
    lastTransactionPassword = transactionPassword;
    return authorizeResult;
  }

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) {
    throw UnimplementedError();
  }
}
