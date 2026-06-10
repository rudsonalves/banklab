import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository_impl.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/transaction_password_status.dart';
import 'package:bankflow/data/services/apis/transaction_password/transaction_password_api.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('TransactionPasswordRepositoryImpl.create', () {
    test('delegates to API and returns success metadata', () async {
      final response = _response();
      final api = _FakeTransactionPasswordApi(
        result: Success(response),
      );
      final repository = TransactionPasswordRepositoryImpl(
        api: api,
        appSection: AppSection(),
      );
      final request = _request();

      final result = await repository.create(request);

      expect(result, isA<Success<TransactionPasswordStatusResponseDto>>());
      expect(result.value, same(response));
      expect(api.createCalls, 1);
      expect(api.lastRequest, same(request));
    });

    test('preserves TRANSACTION_PASSWORD_ALREADY_SET backend code', () async {
      final repository = TransactionPasswordRepositoryImpl(
        api: _FakeTransactionPasswordApi(
          result: const Failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Transaction password already set',
              details: {'code': 'TRANSACTION_PASSWORD_ALREADY_SET'},
            ),
          ),
        ),
        appSection: AppSection(),
      );

      final result = await repository.create(_request());

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(
        backendErrorCode(result.error),
        'TRANSACTION_PASSWORD_ALREADY_SET',
      );
    });

    test('preserves INVALID_DATA backend code', () async {
      final repository = TransactionPasswordRepositoryImpl(
        api: _FakeTransactionPasswordApi(
          result: const Failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Invalid data',
              details: {'code': 'INVALID_DATA'},
            ),
          ),
        ),
        appSection: AppSection(),
      );

      final result = await repository.create(_request());

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(backendErrorCode(result.error), 'INVALID_DATA');
    });

    test('preserves invalid session backend code', () async {
      final repository = TransactionPasswordRepositoryImpl(
        api: _FakeTransactionPasswordApi(
          result: const Failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Invalid token',
              details: {'code': 'INVALID_TOKEN'},
            ),
          ),
        ),
        appSection: AppSection(),
      );

      final result = await repository.create(_request());

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(backendErrorCode(result.error), 'INVALID_TOKEN');
    });

    test('propagates unexpected errors', () async {
      final repository = TransactionPasswordRepositoryImpl(
        api: _FakeTransactionPasswordApi(
          result: const Failure(
            AppError(
              code: AppErrorCode.unexpected,
              message: 'Unexpected failure',
            ),
          ),
        ),
        appSection: AppSection(),
      );

      final result = await repository.create(_request());

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(result.error?.code, AppErrorCode.unexpected);
      expect(result.error?.message, 'Unexpected failure');
      expect(backendErrorCode(result.error), isNull);
    });
  });
}

CreateTransactionPasswordRequestDto _request() {
  return CreateTransactionPasswordRequestDto(
    password: '123456',
    confirmation: '123456',
  );
}

TransactionPasswordStatusResponseDto _response() {
  return TransactionPasswordStatusResponseDto(
    userId: 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b',
    status: TransactionPasswordStatus.active,
    createdAt: DateTime.parse('2026-05-18T12:03:00Z'),
  );
}

class _FakeTransactionPasswordApi extends TransactionPasswordApi {
  _FakeTransactionPasswordApi({
    required this.result,
  }) : super(_NoopRestClient());

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

class _NoopRestClient implements RestClient {
  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      throw UnimplementedError();
}
