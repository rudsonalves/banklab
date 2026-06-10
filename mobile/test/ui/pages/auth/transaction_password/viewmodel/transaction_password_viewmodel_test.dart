import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/transaction_password_status.dart';
import 'package:bankflow/ui/pages/auth/transaction_password/viewmodel/transaction_password_viewmodel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('TransactionPasswordViewModel.create', () {
    test('delegates to repository and returns success metadata', () async {
      final response = TransactionPasswordStatusResponseDto(
        userId: 'user-1',
        status: TransactionPasswordStatus.active,
        createdAt: DateTime.parse('2026-05-18T12:03:00Z'),
      );
      final repository = _FakeTransactionPasswordRepository(
        result: Success(response),
      );
      final viewModel = TransactionPasswordViewModel(
        repository: repository,
        appSection: AppSection(),
      );
      final request = CreateTransactionPasswordRequestDto(
        password: '123456',
        confirmation: '123456',
      );

      await viewModel.create.execute(request);

      expect(viewModel.create.isSuccess, isTrue);
      expect(viewModel.create.value, same(response));
      expect(repository.createCalls, 1);
      expect(repository.lastRequest, same(request));
    });

    test('exposes repository failure', () async {
      final repository = _FakeTransactionPasswordRepository(
        result: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'Invalid data',
            details: {'code': 'INVALID_DATA'},
          ),
        ),
      );
      final viewModel = TransactionPasswordViewModel(
        repository: repository,
        appSection: AppSection(),
      );

      await viewModel.create.execute(
        CreateTransactionPasswordRequestDto(
          password: '123456',
          confirmation: '654321',
        ),
      );

      expect(viewModel.create.isFailure, isTrue);
      expect(backendErrorCode(viewModel.create.error), 'INVALID_DATA');
      expect(repository.createCalls, 1);
    });
  });
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  _FakeTransactionPasswordRepository({
    required this.result,
  });

  Result<TransactionPasswordStatusResponseDto> result;
  CreateTransactionPasswordRequestDto? lastRequest;
  int createCalls = 0;

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) async {
    createCalls++;
    lastRequest = dto;
    return result;
  }
}
