import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/receipt/api_receipt.dart';
import 'package:bankflow/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import 'package:bankflow/domain/common/receipt/enums/transfer_receipt_status.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Money brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('ApiReceipt.getTransferReceipt', () {
    test(
      'calls receipt endpoint with transaction_reference path parameter',
      () async {
        final client = _FakeRestClient(
          getResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successEnvelope(),
            ),
          ),
        );
        final api = ApiReceipt(client);

        final result = await api.getReceipt(
          '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
        );

        expect(result, isA<Success<TransferReceiptResponseDto>>());
        expect(client.getCalls, 1);
        expect(
          client.lastGetRequest?.path,
          '/accounts/transfer/2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31/receipt',
        );
      },
    );

    test('parses representative success envelope', () async {
      final api = ApiReceipt(
        _FakeRestClient(
          getResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _successEnvelope(),
            ),
          ),
        ),
      );

      final result = await api.getReceipt(
        '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
      );

      expect(result, isA<Success<TransferReceiptResponseDto>>());
      final dto = result.value!;
      expect(dto.transactionReference, '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31');
      expect(dto.amount, brl(2500));
      expect(dto.status, TransferReceiptStatus.completed);
      expect(dto.description, 'Aluguel de maio');
    });

    test('maps not found backend envelope to AppError failure', () async {
      final api = ApiReceipt(
        _FakeRestClient(
          getResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'TRANSACTION_NOT_FOUND',
                  'message': 'transfer receipt does not exist',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.getReceipt('missing-reference');

      expect(result, isA<Failure<TransferReceiptResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'transfer receipt does not exist');
    });

    test('maps forbidden backend envelope to AppError failure', () async {
      final api = ApiReceipt(
        _FakeRestClient(
          getResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'FORBIDDEN',
                  'message': 'authenticated user cannot access this receipt',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.getReceipt(
        '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
      );

      expect(result, isA<Failure<TransferReceiptResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(
        result.error?.message,
        'authenticated user cannot access this receipt',
      );
    });

    test('maps unknown HTTP failure to generic HTTP error behavior', () async {
      final api = ApiReceipt(
        _FakeRestClient(
          getResult: const Result.success(
            RestClientResponse(
              statusCode: 500,
              statusMessage: 'Internal Server Error',
              data: {'data': null, 'error': null},
            ),
          ),
        ),
      );

      final result = await api.getReceipt(
        '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
      );

      expect(result, isA<Failure<TransferReceiptResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'HTTP error: 500 Internal Server Error');
    });
  });
}

Map<String, dynamic> _successEnvelope() => {
  'data': {
    'operation_type': 'transfer_out',
    'amount': 2500,
    'status': 'completed',
    'transaction_reference': '2e3ef0c7-ef10-4f4e-a62b-56c71c3c5b31',
    'operation_date': '2026-05-06T12:30:00Z',
    'source_branch': '0001',
    'source_account_number': '00012345',
    'destination_branch': '0001',
    'destination_account_number': '00067890',
    'recipient_name': 'Maria Silva',
    'description': 'Aluguel de maio',
  },
  'error': null,
};

class _FakeRestClient implements RestClient {
  _FakeRestClient({required this.getResult});

  final Result<RestClientResponse> getResult;
  RestClientRequest? lastGetRequest;
  int getCalls = 0;

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async {
    getCalls++;
    lastGetRequest = request;
    return getResult;
  }

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async =>
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
