import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository_impl.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/step_up_operation.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/transaction_password_status.dart';
import 'package:bankflow/data/services/apis/transaction_password/transaction_password_api.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart'
    as auth;
import 'package:bankflow/domain/common/user/enums/user_role.dart';
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

    test('updates transaction password status from notSet to active', () async {
      final appSection = AppSection()
        ..setAuthSession(
          _authSession(auth.TransactionPasswordStatus.notSet),
        );
      final repository = TransactionPasswordRepositoryImpl(
        api: _FakeTransactionPasswordApi(
          result: Success(_response()),
        ),
        appSection: appSection,
      );

      final result = await repository.create(_request());

      expect(result, isA<Success<TransactionPasswordStatusResponseDto>>());
      expect(
        appSection.readiness?.transactionPasswordStatus,
        auth.TransactionPasswordStatus.active,
      );
    });

    test('preserves TRANSACTION_PASSWORD_ALREADY_SET backend code', () async {
      final appSection = AppSection()
        ..setAuthSession(
          _authSession(auth.TransactionPasswordStatus.notSet),
        );
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
        appSection: appSection,
      );

      final result = await repository.create(_request());

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(
        backendErrorCode(result.error),
        'TRANSACTION_PASSWORD_ALREADY_SET',
      );
      expect(
        appSection.readiness?.transactionPasswordStatus,
        auth.TransactionPasswordStatus.notSet,
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

  group('TransactionPasswordRepositoryImpl.stepUpAuthorize', () {
    test('delegates to API and returns its authorization response', () async {
      final api = _FakeStepUpTransactionPasswordApi(
        result: Success(
          StepUpAuthorizeResponseDto(
            stepUpToken: 'opaque-step-up-token',
            expiresIn: 120,
          ),
        ),
      );
      final repository = TransactionPasswordRepositoryImpl(
        api: api,
        appSection: AppSection(),
      );
      final request = _stepUpRequest();

      final result = await repository.stepUpAuthorize(request);

      expect(result, isA<Success<StepUpAuthorizeResponseDto>>());
      expect(result.value?.stepUpToken, 'opaque-step-up-token');
      expect(result.value?.expiresIn, 120);
      expect(api.calls, 1);
      expect(api.lastRequest, same(request));
    });

    test(
      'updates transaction password status from active to notSet when missing',
      () async {
        final appSection = AppSection()
          ..setAuthSession(
            _authSession(auth.TransactionPasswordStatus.active),
          );
        final repository = TransactionPasswordRepositoryImpl(
          api: _FakeStepUpTransactionPasswordApi(
            result: const Failure(
              AppError(
                code: AppErrorCode.transactionPasswordNotSet,
                message: 'Transaction password not set',
                details: {'code': 'TRANSACTION_PASSWORD_NOT_SET'},
              ),
            ),
          ),
          appSection: appSection,
        );

        final result = await repository.stepUpAuthorize(_stepUpRequest());

        expect(result, isA<Failure<StepUpAuthorizeResponseDto>>());
        expect(
          backendErrorCode(result.error),
          'TRANSACTION_PASSWORD_NOT_SET',
        );
        expect(
          appSection.readiness?.transactionPasswordStatus,
          auth.TransactionPasswordStatus.notSet,
        );
      },
    );

    test('preserves active status for failures other than notSet', () async {
      final appSection = AppSection()
        ..setAuthSession(
          _authSession(auth.TransactionPasswordStatus.active),
        );
      final repository = TransactionPasswordRepositoryImpl(
        api: _FakeStepUpTransactionPasswordApi(
          result: const Failure(
            AppError(
              code: AppErrorCode.transactionPasswordLocked,
              message: 'Transaction password locked',
              details: {'code': 'TRANSACTION_PASSWORD_LOCKED'},
            ),
          ),
        ),
        appSection: appSection,
      );

      final result = await repository.stepUpAuthorize(_stepUpRequest());

      expect(result, isA<Failure<StepUpAuthorizeResponseDto>>());
      expect(
        backendErrorCode(result.error),
        'TRANSACTION_PASSWORD_LOCKED',
      );
      expect(
        appSection.readiness?.transactionPasswordStatus,
        auth.TransactionPasswordStatus.active,
      );
    });

    for (final errorCode in [
      'TRANSACTION_PASSWORD_INVALID',
      'INVALID_TOKEN',
      'STEP_UP_POLICY_DENIED',
    ]) {
      test('preserves $errorCode', () async {
        final repository = TransactionPasswordRepositoryImpl(
          api: _FakeStepUpTransactionPasswordApi(
            result: Failure(
              AppError(
                code: AppErrorCode.httpError,
                message: 'Authorization failed',
                details: {'code': errorCode},
              ),
            ),
          ),
          appSection: AppSection(),
        );

        final result = await repository.stepUpAuthorize(_stepUpRequest());

        expect(result, isA<Failure<StepUpAuthorizeResponseDto>>());
        expect(backendErrorCode(result.error), errorCode);
      });
    }
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

StepUpAuthorizeRequestDto _stepUpRequest() {
  return StepUpAuthorizeRequestDto(
    operation: StepUpOperation.internalTransfer,
    transactionPassword: '123456',
  );
}

auth.AuthSession _authSession(auth.TransactionPasswordStatus status) {
  return auth.AuthSession(
    user: auth.UserSession(
      userId: 'user-1',
      email: 'customer@example.com',
      role: UserRole.customer,
    ),
    customer: auth.CustomerSession(
      id: 'customer-1',
      name: 'Maria Silva',
      cpf: '12345678901',
      birthDate: DateTime(1990, 1, 1),
      createdAt: DateTime(2026, 5, 13),
    ),
    readiness: auth.ReadinessSession(
      onboardingCompleted: true,
      approved: true,
      hasOperationalAccount: true,
      transactionPasswordStatus: status,
    ),
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

class _FakeStepUpTransactionPasswordApi extends TransactionPasswordApi {
  _FakeStepUpTransactionPasswordApi({
    required this.result,
  }) : super(_NoopRestClient());

  final Result<StepUpAuthorizeResponseDto> result;
  StepUpAuthorizeRequestDto? lastRequest;
  int calls = 0;

  @override
  AsyncResult<StepUpAuthorizeResponseDto> stepUpAuthorize(
    StepUpAuthorizeRequestDto dto,
  ) async {
    calls++;
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
