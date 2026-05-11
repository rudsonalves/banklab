import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/repositories/transaction/transaction_repository_impl.dart';
import 'package:bankflow/data/services/apis/receipt/api_receipt.dart';
import 'package:bankflow/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import 'package:bankflow/data/services/apis/transfer/api_transfer.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'package:bankflow/domain/common/receipt/enums/transfer_receipt_status.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Money brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('TransactionRepositoryImpl.transfer', () {
    test('returns success and caches lastTransfer', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(
          receiptResult: Success(_receiptResponse()),
        ),
      );

      final request = _validTransferRequest();

      final result = await repository.transfer(request);

      expect(result, isA<Success<TransferResponseDto>>());
      expect(transferApi.transferCalls, 1);
      expect(transferApi.lastTransferRequest, same(request));
      expect(repository.lastTransfer?.transactionReference, 'tx-ref-001');
    });

    test('returns backend failure and clears previous lastTransfer', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(
          receiptResult: Success(_receiptResponse()),
        ),
      );

      final success = await repository.transfer(_validTransferRequest());
      expect(success, isA<Success<TransferResponseDto>>());
      expect(repository.lastTransfer?.transactionReference, 'tx-ref-001');

      transferApi.transferResult = const Failure(
        AppError(
          code: AppErrorCode.httpError,
          message: 'source account has insufficient funds',
        ),
      );

      final result = await repository.transfer(_validTransferRequest());

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'source account has insufficient funds');
      expect(transferApi.transferCalls, 2);
      expect(repository.lastTransfer, isNull);
    });

    test(
      'fails before API call when selected source account is missing',
      () async {
        final transferApi = _FakeApiTransfer(
          transferResult: Success(_transferResponse()),
        );
        final repository = TransactionRepositoryImpl(
          apiTransfer: transferApi,
          apiReceipt: _FakeApiReceipt(
            receiptResult: Success(_receiptResponse()),
          ),
        );

        final result = await repository.transfer(
          TransferRequestDto(
            fromAccountId: '',
            toAccountId: 'acc-dst-001',
            amount: brl(2500),
            idempotencyKey: 'idempotency-key',
          ),
        );

        expect(result, isA<Failure<TransferResponseDto>>());
        expect(result.error?.code, AppErrorCode.unexpected);
        expect(result.error?.message, 'No account selected.');
        expect(transferApi.transferCalls, 0);
        expect(repository.lastTransfer, isNull);
      },
    );

    test('fails before API call when destination account is missing', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(
          receiptResult: Success(_receiptResponse()),
        ),
      );

      final result = await repository.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: '',
          amount: brl(2500),
          idempotencyKey: 'idempotency-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.unexpected);
      expect(result.error?.message, 'Destination account is required.');
      expect(transferApi.transferCalls, 0);
      expect(repository.lastTransfer, isNull);
    });

    test('fails before API call when amount is not positive', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(
          receiptResult: Success(_receiptResponse()),
        ),
      );

      final result = await repository.transfer(
        TransferRequestDto(
          fromAccountId: 'acc-src-001',
          toAccountId: 'acc-dst-001',
          amount: brl(0),
          idempotencyKey: 'idempotency-key',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.unexpected);
      expect(result.error?.message, 'Amount must be greater than zero.');
      expect(transferApi.transferCalls, 0);
      expect(repository.lastTransfer, isNull);
    });
  });

  group('TransactionRepositoryImpl.getTransferReceipt', () {
    test('returns success and caches lastReceipt', () async {
      final receiptApi = _FakeApiReceipt(
        receiptResult: Success(_receiptResponse()),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: _FakeApiTransfer(
          transferResult: Success(_transferResponse()),
        ),
        apiReceipt: receiptApi,
      );

      final result = await repository.getTransferReceipt('tx-ref-001');

      expect(result, isA<Success<TransferReceiptResponseDto>>());
      expect(receiptApi.receiptCalls, 1);
      expect(receiptApi.lastTransactionReference, 'tx-ref-001');
      expect(repository.lastReceipt?.transactionReference, 'tx-ref-001');
    });

    test(
      'propagates not found failure and clears previous lastReceipt',
      () async {
        final receiptApi = _FakeApiReceipt(
          receiptResult: Success(_receiptResponse()),
        );
        final repository = TransactionRepositoryImpl(
          apiTransfer: _FakeApiTransfer(
            transferResult: Success(_transferResponse()),
          ),
          apiReceipt: receiptApi,
        );

        final success = await repository.getTransferReceipt('tx-ref-001');
        expect(success, isA<Success<TransferReceiptResponseDto>>());
        expect(repository.lastReceipt?.transactionReference, 'tx-ref-001');

        receiptApi.receiptResult = const Failure(
          AppError(
            code: AppErrorCode.httpError,
            statusCode: 404,
            message: 'transfer receipt does not exist',
            details: {'code': 'TRANSACTION_NOT_FOUND'},
          ),
        );

        final result = await repository.getTransferReceipt('missing-reference');

        expect(result, isA<Failure<TransferReceiptResponseDto>>());
        expect(result.error?.code, AppErrorCode.httpError);
        expect(result.error?.statusCode, 404);
        expect(result.error?.message, 'transfer receipt does not exist');
        expect(result.error?.details, {'code': 'TRANSACTION_NOT_FOUND'});
        expect(receiptApi.receiptCalls, 2);
        expect(receiptApi.lastTransactionReference, 'missing-reference');
        expect(repository.lastReceipt, isNull);
      },
    );

    test('propagates forbidden failure from receipt API', () async {
      final receiptApi = _FakeApiReceipt(
        receiptResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            statusCode: 403,
            message: 'authenticated user cannot access this receipt',
            details: {'code': 'FORBIDDEN'},
          ),
        ),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: _FakeApiTransfer(
          transferResult: Success(_transferResponse()),
        ),
        apiReceipt: receiptApi,
      );

      final result = await repository.getTransferReceipt('forbidden-reference');

      expect(result, isA<Failure<TransferReceiptResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.statusCode, 403);
      expect(
        result.error?.message,
        'authenticated user cannot access this receipt',
      );
      expect(result.error?.details, {'code': 'FORBIDDEN'});
      expect(receiptApi.receiptCalls, 1);
      expect(receiptApi.lastTransactionReference, 'forbidden-reference');
      expect(repository.lastReceipt, isNull);
    });

    test('propagates generic backend failure from receipt API', () async {
      final receiptApi = _FakeApiReceipt(
        receiptResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            statusCode: 500,
            message: 'HTTP error: 500 Internal Server Error',
          ),
        ),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: _FakeApiTransfer(
          transferResult: Success(_transferResponse()),
        ),
        apiReceipt: receiptApi,
      );

      final result = await repository.getTransferReceipt('tx-ref-001');

      expect(result, isA<Failure<TransferReceiptResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.statusCode, 500);
      expect(result.error?.message, 'HTTP error: 500 Internal Server Error');
      expect(receiptApi.receiptCalls, 1);
      expect(receiptApi.lastTransactionReference, 'tx-ref-001');
      expect(repository.lastReceipt, isNull);
    });
  });

  group('TransactionRepositoryImpl.getInternalRecipient', () {
    test('returns lookup success and forwards query to API', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
        recipientResult: Success([_recipientInfo()]),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(receiptResult: Success(_receiptResponse())),
      );

      const request = RecipientRequestDto(document: '12345678901');

      final result = await repository.getInternalRecipient(request);

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      expect(result.value, hasLength(1));
      expect(result.value?.single.accountId, 'acc-recipient-001');
      expect(transferApi.recipientCalls, 1);
      expect(transferApi.lastRecipientRequest, same(request));
    });

    test('returns multiple lookup results', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
        recipientResult: Success([
          _recipientInfo(accountId: 'acc-recipient-001'),
          _recipientInfo(accountId: 'acc-recipient-002'),
        ]),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(receiptResult: Success(_receiptResponse())),
      );

      final result = await repository.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      expect(result.value, hasLength(2));
      expect(result.value?.first.accountId, 'acc-recipient-001');
      expect(result.value?.last.accountId, 'acc-recipient-002');
      expect(transferApi.recipientCalls, 1);
    });

    test('returns empty lookup results', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
        recipientResult: const Success(<RecipientInfoDto>[]),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(receiptResult: Success(_receiptResponse())),
      );

      final result = await repository.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      expect(result.value, isEmpty);
      expect(transferApi.recipientCalls, 1);
    });

    test('fails before API call when lookup query is empty', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
        recipientResult: Success([_recipientInfo()]),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(receiptResult: Success(_receiptResponse())),
      );

      final result = await repository.getInternalRecipient(
        const RecipientRequestDto(),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.unexpected);
      expect(result.error?.message, 'Search query is required.');
      expect(transferApi.recipientCalls, 0);
    });

    test(
      'propagates invalid query, forbidden, and unauthorized lookup failures',
      () async {
        final cases = [
          (
            error: const AppError(
              code: AppErrorCode.httpError,
              statusCode: 400,
              message: 'invalid or unsupported query parameter combination',
              details: {'code': 'INVALID_DATA'},
            ),
          ),
          (
            error: const AppError(
              code: AppErrorCode.httpError,
              statusCode: 403,
              message: 'access denied',
              details: {'code': 'FORBIDDEN'},
            ),
          ),
          (
            error: const AppError(
              code: AppErrorCode.httpError,
              statusCode: 401,
              message: 'authentication required',
              details: {'code': 'UNAUTHORIZED'},
            ),
          ),
        ];

        for (final testCase in cases) {
          final transferApi = _FakeApiTransfer(
            transferResult: Success(_transferResponse()),
            recipientResult: Failure(testCase.error),
          );
          final repository = TransactionRepositoryImpl(
            apiTransfer: transferApi,
            apiReceipt: _FakeApiReceipt(
              receiptResult: Success(_receiptResponse()),
            ),
          );

          final result = await repository.getInternalRecipient(
            const RecipientRequestDto(document: '12345678901'),
          );

          expect(result, isA<Failure<List<RecipientInfoDto>>>());
          expect(result.error?.code, testCase.error.code);
          expect(result.error?.statusCode, testCase.error.statusCode);
          expect(result.error?.message, testCase.error.message);
          expect(result.error?.details, testCase.error.details);
          expect(transferApi.recipientCalls, 1);
        }
      },
    );

    test('propagates generic backend lookup failure', () async {
      final transferApi = _FakeApiTransfer(
        transferResult: Success(_transferResponse()),
        recipientResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            statusCode: 500,
            message: 'HTTP error: 500 Internal Server Error',
          ),
        ),
      );
      final repository = TransactionRepositoryImpl(
        apiTransfer: transferApi,
        apiReceipt: _FakeApiReceipt(receiptResult: Success(_receiptResponse())),
      );

      final result = await repository.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.statusCode, 500);
      expect(result.error?.message, 'HTTP error: 500 Internal Server Error');
      expect(transferApi.recipientCalls, 1);
    });
  });
}

TransferRequestDto _validTransferRequest() {
  return TransferRequestDto(
    fromAccountId: 'acc-src-001',
    toAccountId: 'acc-dst-001',
    amount: brl(2500),
    idempotencyKey: 'idempotency-key',
  );
}

TransferResponseDto _transferResponse() {
  return TransferResponseDto(
    fromAccountId: 'acc-src-001',
    toAccountId: 'acc-dst-001',
    transactionReference: 'tx-ref-001',
    amount: brl(2500),
    fromBalance: brl(97500),
    toBalance: brl(32500),
  );
}

TransferReceiptResponseDto _receiptResponse() {
  return TransferReceiptResponseDto(
    operationType: 'transfer_out',
    amount: brl(2500),
    status: TransferReceiptStatus.completed,
    transactionReference: 'tx-ref-001',
    operationDate: DateTime.utc(2026, 5, 7, 12, 0, 0),
    sourceBranch: '0001',
    sourceAccountNumber: '00012345',
    destinationBranch: '0001',
    destinationAccountNumber: '00067890',
    recipientName: 'Maria Silva',
    description: 'Aluguel de maio',
  );
}

RecipientInfoDto _recipientInfo({String accountId = 'acc-recipient-001'}) {
  return RecipientInfoDto(
    accountId: accountId,
    holderName: 'Maria Silva',
    document: '***.456.789-**',
    accountNumber: '00067890',
  );
}

class _FakeApiTransfer extends ApiTransfer {
  _FakeApiTransfer({
    required this.transferResult,
    this.recipientResult = const Success(<RecipientInfoDto>[]),
  }) : super(_NoopRestClient());

  Result<TransferResponseDto> transferResult;
  Result<List<RecipientInfoDto>> recipientResult;
  int transferCalls = 0;
  int recipientCalls = 0;
  TransferRequestDto? lastTransferRequest;
  RecipientRequestDto? lastRecipientRequest;

  @override
  AsyncResult<TransferResponseDto> transfer(TransferRequestDto dto) async {
    transferCalls++;
    lastTransferRequest = dto;
    return transferResult;
  }

  @override
  AsyncResult<List<RecipientInfoDto>> getInternalRecipient(
    RecipientRequestDto dto,
  ) async {
    recipientCalls++;
    lastRecipientRequest = dto;
    return recipientResult;
  }
}

class _FakeApiReceipt extends ApiReceipt {
  _FakeApiReceipt({required this.receiptResult}) : super(_NoopRestClient());

  Result<TransferReceiptResponseDto> receiptResult;
  int receiptCalls = 0;
  String? lastTransactionReference;

  @override
  AsyncResult<TransferReceiptResponseDto> getReceipt(
    String transactionReference,
  ) async {
    receiptCalls++;
    lastTransactionReference = transactionReference;
    return receiptResult;
  }
}

class _NoopRestClient implements RestClient {
  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async {
    throw UnimplementedError();
  }

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async {
    throw UnimplementedError();
  }

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async {
    throw UnimplementedError();
  }

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async {
    throw UnimplementedError();
  }

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async {
    throw UnimplementedError();
  }
}
