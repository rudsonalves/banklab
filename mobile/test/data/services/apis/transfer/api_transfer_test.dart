import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/transfer/api_transfer.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Money brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('ApiTransfer.transfer', () {
    test(
      'calls POST /accounts/internal-transfers with serialized request body',
      () async {
        final client = _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successEnvelope(),
            ),
          ),
        );
        final api = ApiTransfer(client);

        final result = await api.transfer(
          TransferRequestDto(
            fromAccountId: 'acc-src-001',
            toAccountId: 'acc-dst-001',
            amount: brl(2500),
            idempotencyKey: 'client-key',
            description: 'Aluguel de maio',
          ),
        );

        expect(result, isA<Success<TransferResponseDto>>());
        expect(client.postCalls, 1);
        expect(client.lastPostRequest?.path, '/accounts/internal-transfers');
        expect(client.lastPostRequest?.body, {
          'from_account_id': 'acc-src-001',
          'to_account_id': 'acc-dst-001',
          'amount': 2500,
          'idempotency_key': 'client-key',
          'description': 'Aluguel de maio',
        });
      },
    );

    test('parses representative success envelope', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successEnvelope(),
            ),
          ),
        ),
      );

      final result = await api.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: 'acc-dst-001',
          amount: brl(2500),
          idempotencyKey: 'client-key',
        ),
      );

      expect(result, isA<Success<TransferResponseDto>>());
      final dto = result.value!;
      expect(dto.transactionReference, '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31');
      expect(dto.amount, brl(2500));
      expect(dto.fromBalance, brl(97500));
      expect(dto.toBalance, brl(32500));
      expect(dto.fromAccountId, 'acc-src-001');
      expect(dto.toAccountId, 'acc-dst-001');
    });

    test('maps backend error envelope to AppError failure', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'INSUFFICIENT_FUNDS',
                  'message': 'source account has insufficient funds',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: 'acc-dst-001',
          amount: brl(2500),
          idempotencyKey: 'client-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'source account has insufficient funds');
    });

    test('maps unknown HTTP failure to generic HTTP error behavior', () async {
      final api = ApiTransfer(
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

      final result = await api.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: 'acc-dst-001',
          amount: brl(2500),
          idempotencyKey: 'client-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'HTTP error: 500 Internal Server Error');
    });

    test('returns network failure when RestClient fails', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          postResult: Result.failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Network timeout',
            ),
          ),
        ),
      );

      final result = await api.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: 'acc-dst-001',
          amount: brl(2500),
          idempotencyKey: 'client-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'Network timeout');
    });

    test('returns parsing error when envelope data is malformed', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': {
                  'from_account_number': 'invalid', // missing required fields
                },
                'error': null,
              },
            ),
          ),
        ),
      );

      final result = await api.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: 'acc-dst-001',
          amount: brl(2500),
          idempotencyKey: 'client-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.parsingError);
    });
  });

  group('ApiTransfer.getInternalRecipient', () {
    test(
      'calls GET /accounts/internal-transfers/recipients with query parameters',
      () async {
        final client = _FakeRestClient(
          getResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successRecipientsEnvelope(),
            ),
          ),
        );
        final api = ApiTransfer(client);

        final result = await api.getInternalRecipient(
          const RecipientRequestDto(document: '12345678901'),
        );

        expect(result, isA<Success<List<RecipientInfoDto>>>());
        expect(client.getCalls, 1);
        expect(
          client.lastGetRequest?.path,
          '/accounts/internal-transfers/recipients',
        );
        expect(client.lastGetRequest?.queryParameters, {
          'document': '12345678901',
        });
      },
    );

    test('parses successful envelope with multiple recipients', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          getResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successRecipientsEnvelope(),
            ),
          ),
        ),
      );

      final result = await api.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      final recipients = result.value!;
      expect(recipients, hasLength(2));
      expect(recipients[0].accountId, 'acc-001');
      expect(recipients[0].holderName, 'John Doe');
      expect(recipients[0].document, '12345678901');
      expect(recipients[1].accountId, 'acc-002');
      expect(recipients[1].holderName, 'Jane Smith');
    });

    test(
      'returns empty list when envelope contains empty accounts array',
      () async {
        final api = ApiTransfer(
          _FakeRestClient(
            getResult: const Result.success(
              RestClientResponse(
                statusCode: 200,
                data: {
                  'data': {
                    'accounts': [],
                  },
                  'error': null,
                },
              ),
            ),
          ),
        );

        final result = await api.getInternalRecipient(
          const RecipientRequestDto(document: '12345678901'),
        );

        expect(result, isA<Success<List<RecipientInfoDto>>>());
        expect(result.value, isEmpty);
      },
    );

    test('maps backend error envelope to AppError failure', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          getResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'ACCOUNT_NOT_FOUND',
                  'message': 'Account not found',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'Account not found');
    });

    test('returns HTTP error when status code indicates failure', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          getResult: const Result.success(
            RestClientResponse(
              statusCode: 400,
              statusMessage: 'Bad Request',
              data: {'data': null, 'error': null},
            ),
          ),
        ),
      );

      final result = await api.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'HTTP error: 400 Bad Request');
    });

    test('returns parsing error when recipients data is malformed', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          getResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': {
                  'accounts': [
                    {
                      'account_id': 'acc-001',
                      // missing holder_name field
                    },
                  ],
                },
                'error': null,
              },
            ),
          ),
        ),
      );

      final result = await api.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.parsingError);
    });

    test('returns network failure when RestClient fails', () async {
      final api = ApiTransfer(
        _FakeRestClient(
          getResult: Result.failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Network error',
            ),
          ),
        ),
      );

      final result = await api.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'Network error');
    });

    test(
      'uses branch and account_number as query parameters when provided',
      () async {
        final client = _FakeRestClient(
          getResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successRecipientsEnvelope(),
            ),
          ),
        );
        final api = ApiTransfer(client);

        await api.getInternalRecipient(
          const RecipientRequestDto(
            branch: '0001',
            accountNumber: '12345',
          ),
        );

        expect(client.lastGetRequest?.queryParameters, {
          'branch': '0001',
          'account_number': '12345',
        });
      },
    );
  });
}

Map<String, dynamic> _successEnvelope() => {
  'data': {
    'from_account_id': 'acc-src-001',
    'to_account_id': 'acc-dst-001',
    'transaction_reference': '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
    'to_branch': '0001',
    'amount': 2500,
    'from_balance': 97500,
    'to_balance': 32500,
  },
  'error': null,
};

Map<String, dynamic> _successRecipientsEnvelope() => {
  'data': {
    'accounts': [
      {
        'account_id': 'acc-001',
        'holder_name': 'John Doe',
        'document': '12345678901',
        'branch': '0001',
        'account_number': '00012345',
      },
      {
        'account_id': 'acc-002',
        'holder_name': 'Jane Smith',
        'document': '12345678901',
        'branch': '0002',
        'account_number': '00067890',
      },
    ],
  },
  'error': null,
};

class _FakeRestClient implements RestClient {
  _FakeRestClient({this.postResult, this.getResult});

  final Result<RestClientResponse>? postResult;
  final Result<RestClientResponse>? getResult;
  RestClientRequest? lastPostRequest;
  RestClientRequest? lastGetRequest;
  int postCalls = 0;
  int getCalls = 0;

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async {
    postCalls++;
    lastPostRequest = request;
    final result = postResult;
    if (result == null) throw UnimplementedError();
    return result;
  }

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async {
    getCalls++;
    lastGetRequest = request;
    final result = getResult;
    if (result == null) throw UnimplementedError();
    return result;
  }

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
