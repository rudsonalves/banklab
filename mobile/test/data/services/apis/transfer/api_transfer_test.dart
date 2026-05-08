import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/transfer/api_transfer.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Money brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('ApiTransfer.transfer', () {
    test(
      'calls POST /accounts/transfer with serialized request body',
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
            fromBranch: '0001',
            fromAccountNumber: '00012345',
            toBranch: '0001',
            toAccountNumber: '00067890',
            amount: brl(2500),
            idempotencyKey: 'client-key',
            description: 'Aluguel de maio',
          ),
        );

        expect(result, isA<Success<TransferResponseDto>>());
        expect(client.postCalls, 1);
        expect(client.lastPostRequest?.path, '/accounts/transfer');
        expect(client.lastPostRequest?.body, {
          'from_branch': '0001',
          'from_account_number': '00012345',
          'to_branch': '0001',
          'to_account_number': '00067890',
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
          fromBranch: '0001',
          fromAccountNumber: '00012345',
          toBranch: '0001',
          toAccountNumber: '00067890',
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
          fromBranch: '0001',
          fromAccountNumber: '00012345',
          toBranch: '0001',
          toAccountNumber: '00067890',
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
          fromBranch: '0001',
          fromAccountNumber: '00012345',
          toBranch: '0001',
          toAccountNumber: '00067890',
          amount: brl(2500),
          idempotencyKey: 'client-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'HTTP error: 500 Internal Server Error');
    });
  });
}

Map<String, dynamic> _successEnvelope() => {
  'data': {
    'from_branch': '0001',
    'from_account_number': '00012345',
    'transaction_reference': '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
    'to_branch': '0001',
    'to_account_number': '00067890',
    'amount': 2500,
    'from_balance': 97500,
    'to_balance': 32500,
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
