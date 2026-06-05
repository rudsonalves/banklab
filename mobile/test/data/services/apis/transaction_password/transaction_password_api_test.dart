import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/transaction_password_status.dart';
import 'package:bankflow/data/services/apis/transaction_password/transaction_password_api.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('TransactionPasswordApi.create', () {
    test(
      'returns success when backend responds with data envelope',
      () async {
        final client = _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 201,
              data: _successEnvelope(),
            ),
          ),
        );
        final api = TransactionPasswordApi(client);

        final result = await api.create(
          CreateTransactionPasswordRequestDto(
            password: '123456',
            confirmation: '123456',
          ),
        );

        expect(result, isA<Success<TransactionPasswordStatusResponseDto>>());
        expect(client.postCalls, 1);
        expect(client.lastPostRequest?.path, '/security/transaction-password');
        expect(client.lastPostRequest?.headers, isNull);
        expect(client.lastPostRequest?.body, {
          'transaction_password': '123456',
          'transaction_password_confirmation': '123456',
        });

        final dto = result.value!;
        expect(dto.userId, 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b');
        expect(dto.status, TransactionPasswordStatus.active);
      },
    );

    test('maps backend envelope error to AppError failure', () async {
      final api = TransactionPasswordApi(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'VALIDATION_ERROR',
                  'message': 'transaction password is invalid',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.create(
        CreateTransactionPasswordRequestDto(
          password: '123456',
          confirmation: '123456',
        ),
      );

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'transaction password is invalid');
      expect(backendErrorCode(result.error), 'VALIDATION_ERROR');
    });

    test(
      'maps already set backend envelope error to specific AppErrorCode',
      () async {
        final api = TransactionPasswordApi(
          _FakeRestClient(
            postResult: const Result.success(
              RestClientResponse(
                statusCode: 200,
                data: {
                  'data': null,
                  'error': {
                    'code': 'TRANSACTION_PASSWORD_ALREADY_SET',
                    'message': 'transaction password already set',
                  },
                },
              ),
            ),
          ),
        );

        final result = await api.create(
          CreateTransactionPasswordRequestDto(
            password: '123456',
            confirmation: '123456',
          ),
        );

        expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
        expect(result.error?.code, AppErrorCode.transactionPasswordAlreadySet);
        expect(result.error?.message, 'transaction password already set');
        expect(
          backendErrorCode(result.error),
          'TRANSACTION_PASSWORD_ALREADY_SET',
        );
      },
    );

    test('maps non-2xx response to HTTP error failure', () async {
      final api = TransactionPasswordApi(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 500,
              statusMessage: 'Internal Server Error',
              data: {'data': null, 'error': null},
            ),
          ),
        ),
      );

      final result = await api.create(
        CreateTransactionPasswordRequestDto(
          password: '123456',
          confirmation: '123456',
        ),
      );

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'HTTP error: 500 Internal Server Error');
    });

    test('returns client failure when RestClient.post fails', () async {
      final api = TransactionPasswordApi(
        _FakeRestClient(
          postResult: Result.failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Network timeout',
            ),
          ),
        ),
      );

      final result = await api.create(
        CreateTransactionPasswordRequestDto(
          password: '123456',
          confirmation: '123456',
        ),
      );

      expect(result, isA<Failure<TransactionPasswordStatusResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'Network timeout');
    });
  });
}

Map<String, dynamic> _successEnvelope() => {
  'data': {
    'user_id': 'fb3a1709-57a9-4c35-ba90-5a5dca6fdb4b',
    'status': 'active',
    'created_at': '2026-05-18T12:03:00Z',
  },
  'error': null,
};

class _FakeRestClient implements RestClient {
  _FakeRestClient({required this.postResult});

  final Result<RestClientResponse> postResult;
  RestClientRequest? lastPostRequest;
  int postCalls = 0;

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async {
    postCalls++;
    lastPostRequest = request;
    return postResult;
  }

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      throw UnimplementedError();
}
