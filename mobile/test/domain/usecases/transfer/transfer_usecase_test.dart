import 'dart:async';

import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/data/repositories/account/account_repository.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository.dart';
import 'package:bankflow/data/repositories/transfer/transfer_repository.dart';
import 'package:bankflow/data/services/apis/account/dtos/account_summary_response_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/balance_response_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/statement_query_params_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/statement_response_dto.dart';
import 'package:bankflow/data/services/apis/receipt/dtos/transfer_receipt_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_request_dto.dart';
import 'package:bankflow/data/services/apis/transfer/dtos/transfer_response_dto.dart';
import 'package:bankflow/domain/usecases/transfer/inputs/protected_transfer_input.dart';
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
        transactionPasswordRepo: _FakeTransactionPasswordRepository(),
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
        transactionPasswordRepo: _FakeTransactionPasswordRepository(),
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
        transactionPasswordRepo: _FakeTransactionPasswordRepository(),
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
        transactionPasswordRepo: _FakeTransactionPasswordRepository(),
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
      'authorizes before building the protected transfer request',
      () async {
        final events = <String>[];
        final transferRepo = _FakeTransferRepository(
          transferResult: Success(_transferResponse()),
          events: events,
        );
        final transactionPasswordRepo = _FakeTransactionPasswordRepository(
          events: events,
        );
        final usecase = TransferUsecase(
          accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
          transferRepo: transferRepo,
          transactionPasswordRepo: transactionPasswordRepo,
        );

        final result = await usecase.transfer(
          _protectedInput(
            draft: TransferDraft(
              toAccountId: 'acc-recipient-001',
              amount: brl(2500),
              idempotencyKey: 'client-key-001',
              description: 'Aluguel de maio',
            ),
          ),
        );

        expect(result, isA<Success<TransferResponseDto>>());
        expect(events, ['authorize', 'transfer']);
        expect(transactionPasswordRepo.authorizeCalls, 1);
        expect(transactionPasswordRepo.lastTransactionPassword, '123456');
        expect(transferRepo.transferCalls, 1);
        expect(transferRepo.lastTransferToken, 'step-up-token');
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

    test('does not transfer when step-up authorization fails', () async {
      final events = <String>[];
      final transferRepo = _FakeTransferRepository(events: events);
      final transactionPasswordRepo = _FakeTransactionPasswordRepository(
        events: events,
        authorizeResults: const [
          Failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Invalid transaction password',
              details: {'code': 'TRANSACTION_PASSWORD_INVALID'},
            ),
          ),
        ],
      );
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
        transactionPasswordRepo: transactionPasswordRepo,
      );

      final result = await usecase.transfer(_protectedInput());

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(
        backendErrorCode(result.error),
        'TRANSACTION_PASSWORD_INVALID',
      );
      expect(events, ['authorize']);
      expect(transferRepo.transferCalls, 0);
    });

    test('uses the authorization token in only one transfer attempt', () async {
      final events = <String>[];
      final transferRepo = _FakeTransferRepository(
        events: events,
        transferResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'source account has insufficient funds',
          ),
        ),
      );
      final transactionPasswordRepo = _FakeTransactionPasswordRepository(
        events: events,
      );
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
        transactionPasswordRepo: transactionPasswordRepo,
      );

      final result = await usecase.transfer(_protectedInput());

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(events, ['authorize', 'transfer']);
      expect(transactionPasswordRepo.authorizeCalls, 1);
      expect(transferRepo.transferCalls, 1);
      expect(transferRepo.transferTokens, ['step-up-token']);
    });

    test(
      'retrying step-up keeps idempotency key and uses new PIN and token',
      () async {
        final events = <String>[];
        final transferRepo = _FakeTransferRepository(events: events);
        final transactionPasswordRepo = _FakeTransactionPasswordRepository(
          events: events,
          authorizeResults: [
            const Failure(
              AppError(
                code: AppErrorCode.httpError,
                message: 'Invalid transaction password',
                details: {'code': 'TRANSACTION_PASSWORD_INVALID'},
              ),
            ),
            Success(
              StepUpAuthorizeResponseDto(
                stepUpToken: 'step-up-token-2',
                expiresIn: 120,
              ),
            ),
          ],
        );
        final usecase = TransferUsecase(
          accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
          transferRepo: transferRepo,
          transactionPasswordRepo: transactionPasswordRepo,
        );
        final draft = _transferDraft(idempotencyKey: 'same-attempt-key');

        final firstResult = await usecase.transfer(
          ProtectedTransferInput(draft: draft, pin: '111111'),
        );
        final secondResult = await usecase.transfer(
          ProtectedTransferInput(draft: draft, pin: '222222'),
        );

        expect(firstResult, isA<Failure<TransferResponseDto>>());
        expect(secondResult, isA<Success<TransferResponseDto>>());
        expect(events, ['authorize', 'authorize', 'transfer']);
        expect(
          transactionPasswordRepo.transactionPasswords,
          ['111111', '222222'],
        );
        expect(transferRepo.transferTokens, ['step-up-token-2']);
        expect(
          transferRepo.transferRequests.single.idempotencyKey,
          'same-attempt-key',
        );
      },
    );

    test('different transfer attempts preserve their distinct keys', () async {
      final transferRepo = _FakeTransferRepository();
      final transactionPasswordRepo = _FakeTransactionPasswordRepository(
        authorizeResults: [
          Success(
            StepUpAuthorizeResponseDto(
              stepUpToken: 'step-up-token-1',
              expiresIn: 120,
            ),
          ),
          Success(
            StepUpAuthorizeResponseDto(
              stepUpToken: 'step-up-token-2',
              expiresIn: 120,
            ),
          ),
        ],
      );
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
        transactionPasswordRepo: transactionPasswordRepo,
      );

      await usecase.transfer(
        _protectedInput(
          draft: _transferDraft(idempotencyKey: 'attempt-key-1'),
        ),
      );
      await usecase.transfer(
        _protectedInput(
          draft: _transferDraft(idempotencyKey: 'attempt-key-2'),
        ),
      );

      expect(
        transferRepo.transferRequests.map((dto) => dto.idempotencyKey),
        ['attempt-key-1', 'attempt-key-2'],
      );
      expect(
        transferRepo.transferTokens,
        ['step-up-token-1', 'step-up-token-2'],
      );
    });

    test('fails when selected recipient is missing', () async {
      final transferRepo = _FakeTransferRepository();
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
        transactionPasswordRepo: _FakeTransactionPasswordRepository(),
      );

      final result = await usecase.transfer(
        _protectedInput(
          draft: TransferDraft(
            toAccountId: '   ',
            amount: brl(2500),
            idempotencyKey: 'client-key-001',
          ),
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
          transactionPasswordRepo: _FakeTransactionPasswordRepository(),
        );

        final result = await usecase.transfer(
          _protectedInput(
            draft: TransferDraft(
              toAccountId: 'acc-src-001',
              amount: brl(2500),
              idempotencyKey: 'client-key-001',
            ),
          ),
        );

        expect(result, isA<Failure<TransferResponseDto>>());
        expect(result.error?.code, AppErrorCode.invalidData);
        expect(transferRepo.transferCalls, 0);
      },
    );

    test('fails when idempotency key is missing', () async {
      final transferRepo = _FakeTransferRepository();
      final transactionPasswordRepo = _FakeTransactionPasswordRepository();
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
        transactionPasswordRepo: transactionPasswordRepo,
      );

      final result = await usecase.transfer(
        _protectedInput(
          draft: TransferDraft(
            toAccountId: 'acc-recipient-001',
            amount: brl(2500),
            idempotencyKey: '   ',
          ),
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.invalidData);
      expect(
        result.error?.message,
        'Idempotency key is required for transfer.',
      );
      expect(transactionPasswordRepo.authorizeCalls, 0);
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
      final transactionPasswordRepo = _FakeTransactionPasswordRepository();
      final usecase = TransferUsecase(
        accountRepo: _FakeAccountRepository(selected: _selectedAccount()),
        transferRepo: transferRepo,
        transactionPasswordRepo: transactionPasswordRepo,
      );

      final result = await usecase.transfer(
        _protectedInput(
          draft: TransferDraft(
            toAccountId: 'acc-recipient-001',
            amount: brl(2500),
            idempotencyKey: 'client-key-001',
          ),
        ),
      );

      expect(result, isA<Failure<TransferResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'source account has insufficient funds');
      expect(transactionPasswordRepo.authorizeCalls, 1);
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

ProtectedTransferInput _protectedInput({
  TransferDraft? draft,
  String pin = '123456',
}) {
  return ProtectedTransferInput(
    draft: draft ?? _transferDraft(),
    pin: pin,
  );
}

TransferDraft _transferDraft({
  String idempotencyKey = 'client-key-001',
}) {
  return TransferDraft(
    toAccountId: 'acc-recipient-001',
    amount: brl(2500),
    idempotencyKey: idempotencyKey,
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
    this.events,
  });

  Result<TransferResponseDto>? transferResult;
  Result<List<RecipientInfoDto>> recipientResult;
  final List<String>? events;
  int transferCalls = 0;
  int recipientCalls = 0;
  String? lastTransferToken;
  TransferRequestDto? lastTransferRequest;
  final List<String> transferTokens = [];
  final List<TransferRequestDto> transferRequests = [];
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
  AsyncResult<TransferResponseDto> transfer({
    required String token,
    required TransferRequestDto dto,
  }) async {
    transferCalls++;
    events?.add('transfer');
    lastTransferToken = token;
    lastTransferRequest = dto;
    transferTokens.add(token);
    transferRequests.add(dto);
    return transferResult ?? Success(_transferResponse());
  }
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  _FakeTransactionPasswordRepository({
    List<Result<StepUpAuthorizeResponseDto>>? authorizeResults,
    this.events,
  }) : authorizeResults =
           authorizeResults ??
           [
             Success(
               StepUpAuthorizeResponseDto(
                 stepUpToken: 'step-up-token',
                 expiresIn: 120,
               ),
             ),
           ];

  final List<Result<StepUpAuthorizeResponseDto>> authorizeResults;
  final List<String>? events;
  final List<String> transactionPasswords = [];
  int authorizeCalls = 0;

  String? get lastTransactionPassword =>
      transactionPasswords.isEmpty ? null : transactionPasswords.last;

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInternalTransfer(
    String transactionPassword,
  ) async {
    events?.add('authorize');
    transactionPasswords.add(transactionPassword);
    final resultIndex = authorizeCalls;
    authorizeCalls++;
    return authorizeResults[resultIndex < authorizeResults.length
        ? resultIndex
        : authorizeResults.length - 1];
  }

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInstallationRegistration(
    String transactionPassword,
  ) {
    throw UnimplementedError();
  }

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) {
    throw UnimplementedError();
  }
}
