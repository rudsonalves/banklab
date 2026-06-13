import 'dart:async';

import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/data/repositories/account/account_repository.dart';
import 'package:bankflow/data/repositories/transfer/transfer_repository.dart';
import 'package:bankflow/data/services/apis/account/dtos/account_summary_response_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/balance_response_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/statement_query_params_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/statement_response_dto.dart';
import 'package:bankflow/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'package:bankflow/domain/usecases/transfer/transfer_usecase.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';

Money brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

void main() {
  group('TransferUsecase.getInternalRecipient', () {
    test('searches recipients by CPF document through repository', () async {
      final transferRepo = _FakeTransferRepository(
        recipientResult: Success([_recipientInfo()]),
      );
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
      );

      const request = RecipientRequestDto(document: '12345678901');

      final result = await usecase.getInternalRecipient(request);

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      expect(result.value, hasLength(1));
      expect(transferRepo.recipientCalls, 1);
      expect(transferRepo.lastRecipientRequest, same(request));
    });

    test('searches recipients by branch and account number', () async {
      final transferRepo = _FakeTransferRepository(
        recipientResult: Success([_recipientInfo()]),
      );
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
      );

      const request = RecipientRequestDto(
        branch: '0001',
        accountNumber: '00067890',
      );

      final result = await usecase.getInternalRecipient(request);

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      expect(transferRepo.recipientCalls, 1);
      expect(transferRepo.lastRecipientRequest, same(request));
    });

    test('returns multiple recipient lookup results', () async {
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: _FakeTransferRepository(
          recipientResult: Success([
            _recipientInfo(accountId: 'acc-recipient-001'),
            _recipientInfo(accountId: 'acc-recipient-002'),
          ]),
        ),
      );

      final result = await usecase.getInternalRecipient(
        const RecipientRequestDto(document: '12345678901'),
      );

      expect(result, isA<Success<List<RecipientInfoDto>>>());
      expect(result.value, hasLength(2));
      expect(result.value?.last.accountId, 'acc-recipient-002');
    });

    test('fails before repository call when lookup query is empty', () async {
      final transferRepo = _FakeTransferRepository();
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
      );

      final result = await usecase.getInternalRecipient(
        const RecipientRequestDto(),
      );

      expect(result, isA<Failure<List<RecipientInfoDto>>>());
      expect(result.error?.code, AppErrorCode.invalidData);
      expect(result.error?.message, 'Recipient search query cannot be empty.');
      expect(transferRepo.recipientCalls, 0);
    });
  });

  group('TransferUsecase.transfer', () {
    test(
      'builds ID-based transfer request from selected account and recipient',
      () async {
        final transferRepo = _FakeTransferRepository(
          transferResult: Success(_transferResponse()),
        );
        final usecase = TransferUsecase(
          accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
          transferRepo: transferRepo,
        );

        final result = await usecase.transfer(
          TransferDraft(
            toAccountId: 'acc-recipient-001',
            amount: brl(2500),
            idempotencyKey: 'client-key-001',
            description: 'Aluguel de maio',
          ),
        );

        expect(result, isA<Success<TransferResponseDto>>());
        expect(transferRepo.transferCalls, 1);
        expect(
          transferRepo.lastTransferRequest?.fromAccountId,
          'acc-src-001',
        );
        expect(
          transferRepo.lastTransferRequest?.toAccountId,
          'acc-recipient-001',
        );
        expect(transferRepo.lastTransferRequest?.amount, brl(2500));
        expect(
          transferRepo.lastTransferRequest?.idempotencyKey,
          'client-key-001',
        );
        expect(
          transferRepo.lastTransferRequest?.description,
          'Aluguel de maio',
        );
      },
    );

    test('fails when selected recipient is missing', () async {
      final transferRepo = _FakeTransferRepository();
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
      );

      final result = await usecase.transfer(
        TransferDraft(
          toAccountId: '   ',
          amount: brl(2500),
          idempotencyKey: 'client-key-001',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.invalidData);
      expect(
        result.error?.message,
        'Destination account ID cannot be empty or the same as the source account.',
      );
      expect(transferRepo.transferCalls, 0);
    });

    test(
      'fails when selected recipient is the selected source account',
      () async {
        final transferRepo = _FakeTransferRepository();
        final usecase = TransferUsecase(
          accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
          transferRepo: transferRepo,
        );

        final result = await usecase.transfer(
          TransferDraft(
            toAccountId: 'acc-src-001',
            amount: brl(2500),
            idempotencyKey: 'client-key-001',
          ),
        );

        expect(result, isA<Failure<TransferResponseDto>>());
        expect(result.error?.code, AppErrorCode.invalidData);
        expect(transferRepo.transferCalls, 0);
      },
    );

    test('fails when idempotency key is missing', () async {
      final transferRepo = _FakeTransferRepository();
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
      );

      final result = await usecase.transfer(
        TransferDraft(toAccountId: 'acc-recipient-001', amount: brl(2500)),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.invalidData);
      expect(
        result.error?.message,
        'Idempotency key is required for transfer.',
      );
      expect(transferRepo.transferCalls, 0);
    });

    test('propagates backend transfer failure', () async {
      final transferRepo = _FakeTransferRepository(
        transferResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'source account has insufficient funds',
          ),
        ),
      );
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
      );

      final result = await usecase.transfer(
        TransferDraft(
          toAccountId: 'acc-recipient-001',
          amount: brl(2500),
          idempotencyKey: 'client-key-001',
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'source account has insufficient funds');
      expect(transferRepo.transferCalls, 1);
    });
  });
}

AccountSummaryResponseDto _selectedAccount() {
  return AccountSummaryResponseDto(
    id: 'acc-src-001',
    customerId: 'cus-001',
    number: '00012345',
    branch: '0001',
    status: 'active',
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

TransferResponseDto _transferResponse() {
  return TransferResponseDto(
    fromAccountId: 'acc-src-001',
    toAccountId: 'acc-recipient-001',
    transactionReference: 'tx-ref-001',
    amount: brl(2500),
    fromBalance: brl(97500),
    toBalance: brl(32500),
  );
}

class _FakeAccountRepository implements AccountRepository {
  _FakeAccountRepository({this.selected});

  AccountSummaryResponseDto? selected;

  @override
  List<AccountSummaryResponseDto>? get accounts =>
      selected == null ? null : [selected!];

  @override
  BalanceResponseDto? get lastBalance => null;

  @override
  StatementResponseDto? get lastStatement => null;

  @override
  AccountSummaryResponseDto? get selectedAccount => selected;

  @override
  Stream<BalanceResponseDto> balance() => const Stream.empty();

  @override
  AsyncResult<StatementResponseDto> getStatement(
    StatementQueryParamsDto queryParams,
  ) {
    throw UnimplementedError();
  }

  @override
  AsyncResult<Unit> loadBalance() async => const Success(unit);

  @override
  AsyncResult<List<AccountSummaryResponseDto>> listAccounts() async =>
      Success(accounts ?? <AccountSummaryResponseDto>[]);

  @override
  AsyncResult<Unit> selectAccount(String accountId) async {
    return const Success(unit);
  }
}

class _FakeTransferRepository implements TransferRepository {
  _FakeTransferRepository({
    this.transferResult,
    this.recipientResult = const Success(<RecipientInfoDto>[]),
  });

  Result<TransferResponseDto>? transferResult;
  Result<List<RecipientInfoDto>> recipientResult;
  int transferCalls = 0;
  int recipientCalls = 0;
  TransferRequestDto? lastTransferRequest;
  RecipientRequestDto? lastRecipientRequest;

  @override
  TransferReceiptResponseDto? get lastReceipt => null;

  @override
  TransferResponseDto? get lastTransfer => null;

  @override
  AsyncResult<TransferReceiptResponseDto> getTransferReceipt(
    String transactionReference,
  ) {
    throw UnimplementedError();
  }

  @override
  AsyncResult<List<RecipientInfoDto>> getInternalRecipient(
    RecipientRequestDto dto,
  ) async {
    recipientCalls++;
    lastRecipientRequest = dto;
    return recipientResult;
  }

  @override
  AsyncResult<TransferResponseDto> transfer(TransferRequestDto dto) async {
    transferCalls++;
    lastTransferRequest = dto;
    return transferResult ?? Success(_transferResponse());
  }
}
