import 'dart:async';

import 'package:bankflow/core/resources/app_currencies.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/routing/routes.dart';
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
import 'package:bankflow/domain/usecases/transfer/transfer_usecase.dart';
import 'package:bankflow/ui/pages/transaction_password/verification/transaction_password_input_page.dart';
import 'package:bankflow/ui/pages/transfer/models/transfer_confirmation_data.dart';
import 'package:bankflow/ui/pages/transfer/transfer_confirmation_page.dart';
import 'package:bankflow/ui/pages/transfer/viewmodel/transfer_viewmodel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:money2/money2.dart';

void main() {
  testWidgets('opens PIN input before starting the protected operation', (
    tester,
  ) async {
    final harness = await _pumpConfirmation(tester);

    await tester.tap(find.text('Transferir'));
    await tester.pumpAndSettle();

    expect(find.byType(TransactionPasswordInputPage), findsOneWidget);
    expect(harness.passwordRepository.authorizeCalls, 0);
    expect(harness.transferRepository.transferCalls, 0);
  });

  testWidgets('cancel does not authorize or transfer and permits reopening', (
    tester,
  ) async {
    final harness = await _pumpConfirmation(tester);

    await tester.tap(find.text('Transferir'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Cancelar'));
    await tester.pumpAndSettle();

    expect(find.text('Confirmação'), findsOneWidget);
    expect(harness.passwordRepository.authorizeCalls, 0);
    expect(harness.transferRepository.transferCalls, 0);

    final button = tester.widget<ElevatedButton>(
      find.widgetWithText(ElevatedButton, 'Transferir'),
    );
    expect(button.onPressed, isNotNull);

    await tester.tap(find.text('Transferir'));
    await tester.pumpAndSettle();
    expect(find.byType(TransactionPasswordInputPage), findsOneWidget);
  });

  testWidgets(
    'blocks duplicate taps, shows loading, and navigates only after transfer',
    (tester) async {
      final authorization = Completer<Result<StepUpAuthorizeResponseDto>>();
      final transfer = Completer<Result<TransferResponseDto>>();
      final harness = await _pumpConfirmation(
        tester,
        authorization: authorization,
        transfer: transfer,
      );

      await tester.tap(find.text('Transferir'));
      await tester.pumpAndSettle();
      await tester.enterText(
        find.byType(TextField, skipOffstage: false),
        '123456',
      );
      await tester.pump();
      await tester.tap(find.text('Concluir'));
      await tester.pump();

      expect(harness.passwordRepository.authorizeCalls, 1);
      expect(harness.passwordRepository.lastPin, '123456');
      expect(harness.transferRepository.transferCalls, 0);
      expect(find.byType(CircularProgressIndicator), findsWidgets);
      expect(find.text('Success destination'), findsNothing);

      await tester.tap(find.text('Transferir'), warnIfMissed: false);
      await tester.pump();
      expect(harness.passwordRepository.authorizeCalls, 1);

      authorization.complete(
        Success(
          StepUpAuthorizeResponseDto(
            stepUpToken: 'single-use-token',
            expiresIn: 120,
          ),
        ),
      );
      await tester.pump();

      expect(harness.transferRepository.transferCalls, 1);
      expect(harness.transferRepository.lastToken, 'single-use-token');
      expect(find.byType(CircularProgressIndicator), findsWidgets);
      expect(find.text('Success destination'), findsNothing);

      transfer.complete(Success(_transferResponse()));
      await tester.pumpAndSettle();

      expect(find.text('Success destination'), findsOneWidget);
    },
  );

  testWidgets('failure navigates to the existing failure status flow', (
    tester,
  ) async {
    final harness = await _pumpConfirmation(
      tester,
      authorizationResult: Success(
        StepUpAuthorizeResponseDto(
          stepUpToken: 'single-use-token',
          expiresIn: 120,
        ),
      ),
      transferResult: const Failure(
        AppError(
          code: AppErrorCode.httpError,
          message: 'Transfer failed',
        ),
      ),
    );

    await tester.tap(find.text('Transferir'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '123456',
    );
    await tester.pump();
    await tester.tap(find.text('Concluir'));
    await tester.pumpAndSettle();

    expect(find.text('Failure destination'), findsOneWidget);
    expect(harness.passwordRepository.authorizeCalls, 1);
    expect(harness.transferRepository.transferCalls, 1);
  });

  testWidgets('invalid password requests a new PIN without transferring', (
    tester,
  ) async {
    final harness = await _pumpConfirmation(
      tester,
      authorizationResults: [
        _failure('TRANSACTION_PASSWORD_INVALID'),
        Success(
          StepUpAuthorizeResponseDto(
            stepUpToken: 'second-token',
            expiresIn: 120,
          ),
        ),
      ],
    );

    await _startAndSubmitPin(tester, '111111');
    expect(find.byType(TransactionPasswordInputPage), findsOneWidget);
    expect(harness.transferRepository.transferCalls, 0);

    await _submitVisiblePin(tester, '222222');

    expect(find.text('Success destination'), findsOneWidget);
    expect(harness.passwordRepository.pins, ['111111', '222222']);
    expect(harness.transferRepository.transferTokens, ['second-token']);
  });

  for (final errorCode in const [
    'STEP_UP_TOKEN_EXPIRED',
    'STEP_UP_TOKEN_CONSUMED',
  ]) {
    testWidgets(
      '$errorCode requests a new token with the same idempotency key',
      (
        tester,
      ) async {
        final harness = await _pumpConfirmation(
          tester,
          authorizationResults: [
            Success(
              StepUpAuthorizeResponseDto(
                stepUpToken: 'expired-token',
                expiresIn: 120,
              ),
            ),
            Success(
              StepUpAuthorizeResponseDto(
                stepUpToken: 'fresh-token',
                expiresIn: 120,
              ),
            ),
          ],
          transferResults: [
            _failure(errorCode),
            Success(_transferResponse()),
          ],
        );

        await _startAndSubmitPin(tester, '111111');
        expect(find.byType(TransactionPasswordInputPage), findsOneWidget);

        await _submitVisiblePin(tester, '222222');

        expect(find.text('Success destination'), findsOneWidget);
        expect(
          harness.transferRepository.transferTokens,
          ['expired-token', 'fresh-token'],
        );
        expect(
          harness.transferRepository.transferRequests
              .map((request) => request.idempotencyKey)
              .toSet(),
          hasLength(1),
        );
      },
    );
  }

  testWidgets('locked password blocks another immediate attempt', (
    tester,
  ) async {
    final harness = await _pumpConfirmation(
      tester,
      authorizationResult: _failure('TRANSACTION_PASSWORD_LOCKED'),
    );

    await _startAndSubmitPin(tester, '111111');

    expect(find.text('Confirmação'), findsOneWidget);
    expect(harness.transferRepository.transferCalls, 0);
    final button = tester.widget<ElevatedButton>(
      find.widgetWithText(ElevatedButton, 'Transferir'),
    );
    expect(button.onPressed, isNull);
  });

  testWidgets('password not set returns home without transferring', (
    tester,
  ) async {
    final harness = await _pumpConfirmation(
      tester,
      authorizationResult: _failure('TRANSACTION_PASSWORD_NOT_SET'),
    );

    await _startAndSubmitPin(tester, '111111');

    expect(harness.transferRepository.transferCalls, 0);
    expect(find.text('Home destination'), findsOneWidget);
  });

  for (final errorCode in const [
    'STEP_UP_ENDPOINT_NOT_ALLOWED',
    'UNAUTHORIZED',
    'INVALID_TOKEN',
    'FORBIDDEN',
    'INVALID_DATA',
    'INVALID_REQUEST',
  ]) {
    testWidgets('$errorCode ends authorization through the failure flow', (
      tester,
    ) async {
      final harness = await _pumpConfirmation(
        tester,
        authorizationResult: _failure(errorCode),
      );

      await _startAndSubmitPin(tester, '111111');

      expect(find.text('Failure destination'), findsOneWidget);
      expect(harness.transferRepository.transferCalls, 0);
    });
  }

  for (final errorCode in const [
    'STEP_UP_TOKEN_REQUIRED',
    'STEP_UP_TOKEN_INVALID',
    'STEP_UP_ENDPOINT_MISMATCH',
  ]) {
    testWidgets('$errorCode does not reuse the rejected token', (tester) async {
      final harness = await _pumpConfirmation(
        tester,
        transferResult: _failure(errorCode),
      );

      await _startAndSubmitPin(tester, '111111');

      expect(find.text('Failure destination'), findsOneWidget);
      expect(harness.passwordRepository.authorizeCalls, 1);
      expect(harness.transferRepository.transferCalls, 1);
      expect(harness.transferRepository.transferTokens, ['single-use-token']);
    });
  }
}

Future<void> _startAndSubmitPin(WidgetTester tester, String pin) async {
  await tester.tap(find.text('Transferir'));
  await tester.pumpAndSettle();
  await _submitVisiblePin(tester, pin);
}

Future<void> _submitVisiblePin(WidgetTester tester, String pin) async {
  await tester.enterText(
    find.byType(TextField, skipOffstage: false),
    pin,
  );
  await tester.pump();
  await tester.tap(find.text('Concluir'));
  await tester.pumpAndSettle();
}

Failure<T> _failure<T extends Object>(String code) {
  return Failure(
    AppError(
      code: AppErrorCode.httpError,
      message: code,
      details: {'code': code},
    ),
  );
}

Future<_Harness> _pumpConfirmation(
  WidgetTester tester, {
  Completer<Result<StepUpAuthorizeResponseDto>>? authorization,
  Completer<Result<TransferResponseDto>>? transfer,
  Result<StepUpAuthorizeResponseDto>? authorizationResult,
  List<Result<StepUpAuthorizeResponseDto>>? authorizationResults,
  Result<TransferResponseDto>? transferResult,
  List<Result<TransferResponseDto>>? transferResults,
}) async {
  final accountRepository = _FakeAccountRepository();
  final passwordRepository = _FakeTransactionPasswordRepository(
    completer: authorization,
    result:
        authorizationResult ??
        Success(
          StepUpAuthorizeResponseDto(
            stepUpToken: 'single-use-token',
            expiresIn: 120,
          ),
        ),
    results: authorizationResults,
  );
  final transferRepository = _FakeTransferRepository(
    completer: transfer,
    result: transferResult ?? Success(_transferResponse()),
    results: transferResults,
  );
  final viewModel = TransferViewmodel(
    TransferUsecase(
      accountRepo: accountRepository,
      transferRepo: transferRepository,
      transactionPasswordRepo: passwordRepository,
    ),
  );
  final router = GoRouter(
    initialLocation: TransferRoutes.confirmation.routePath,
    routes: [
      GoRoute(
        path: TransferRoutes.confirmation.routePath,
        name: TransferRoutes.confirmation.routeName,
        builder: (context, state) => TransferConfirmationPage(
          viewModel: viewModel,
          transferData: _confirmationData(),
        ),
      ),
      GoRoute(
        path: BaseRoutes.home.routePath,
        name: BaseRoutes.home.routeName,
        builder: (context, state) =>
            const Scaffold(body: Text('Home destination')),
      ),
      GoRoute(
        path: TransactionPasswordRoutes.transactionPassword.routePath,
        name: TransactionPasswordRoutes.transactionPassword.routeName,
        builder: (context, state) => const TransactionPasswordInputPage(),
      ),
      GoRoute(
        path: TransferRoutes.statusSuccess.routePath,
        name: TransferRoutes.statusSuccess.routeName,
        builder: (context, state) =>
            const Scaffold(body: Text('Success destination')),
      ),
      GoRoute(
        path: TransferRoutes.statusFailure.routePath,
        name: TransferRoutes.statusFailure.routeName,
        builder: (context, state) =>
            const Scaffold(body: Text('Failure destination')),
      ),
    ],
  );
  addTearDown(router.dispose);

  await tester.pumpWidget(MaterialApp.router(routerConfig: router));
  await tester.pumpAndSettle();

  return _Harness(
    passwordRepository: passwordRepository,
    transferRepository: transferRepository,
  );
}

TransferConfirmationData _confirmationData() {
  return TransferConfirmationData(
    toAccountId: 'acc-dst-001',
    toHolderName: 'Maria Silva',
    toDocument: '***.456.789-**',
    toNumber: '00067890',
    fromAccountId: 'acc-src-001',
    amount: _brl(2500),
    description: 'Aluguel',
  );
}

TransferResponseDto _transferResponse() {
  return TransferResponseDto(
    fromAccountId: 'acc-src-001',
    toAccountId: 'acc-dst-001',
    transactionReference: 'tx-ref-001',
    amount: _brl(2500),
    fromBalance: _brl(97500),
    toBalance: _brl(32500),
  );
}

Money _brl(int cents) =>
    Money.fromBigIntWithCurrency(BigInt.from(cents), AppCurrencies.brl);

class _Harness {
  final _FakeTransactionPasswordRepository passwordRepository;
  final _FakeTransferRepository transferRepository;

  const _Harness({
    required this.passwordRepository,
    required this.transferRepository,
  });
}

class _FakeAccountRepository implements AccountRepository {
  final AccountSummaryResponseDto _selected = AccountSummaryResponseDto(
    id: 'acc-src-001',
    customerId: 'customer-1',
    number: '00012345',
    branch: '0001',
    status: 'active',
  );

  @override
  List<AccountSummaryResponseDto> get accounts => [_selected];

  @override
  BalanceResponseDto? get lastBalance => null;

  @override
  StatementResponseDto? get lastStatement => null;

  @override
  AccountSummaryResponseDto get selectedAccount => _selected;

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
      Success(accounts);

  @override
  AsyncResult<Unit> selectAccount(String accountId) async =>
      const Success(unit);
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  _FakeTransactionPasswordRepository({
    required this.result,
    this.completer,
    this.results,
  });

  final Result<StepUpAuthorizeResponseDto> result;
  final Completer<Result<StepUpAuthorizeResponseDto>>? completer;
  final List<Result<StepUpAuthorizeResponseDto>>? results;
  int authorizeCalls = 0;
  String? lastPin;
  final pins = <String>[];

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInternalTransfer(
    String transactionPassword,
  ) async {
    authorizeCalls++;
    lastPin = transactionPassword;
    pins.add(transactionPassword);
    if (completer != null) return completer!.future;
    final configuredResults = results;
    if (configuredResults == null || configuredResults.isEmpty) return result;
    final index = authorizeCalls - 1;
    return configuredResults[index < configuredResults.length
        ? index
        : configuredResults.length - 1];
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

class _FakeTransferRepository implements TransferRepository {
  _FakeTransferRepository({
    required this.result,
    this.completer,
    this.results,
  });

  final Result<TransferResponseDto> result;
  final Completer<Result<TransferResponseDto>>? completer;
  final List<Result<TransferResponseDto>>? results;
  int transferCalls = 0;
  String? lastToken;
  final transferTokens = <String>[];
  final transferRequests = <TransferRequestDto>[];

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
  ) {
    throw UnimplementedError();
  }

  @override
  AsyncResult<TransferResponseDto> transfer({
    required String token,
    required TransferRequestDto dto,
  }) async {
    transferCalls++;
    lastToken = token;
    transferTokens.add(token);
    transferRequests.add(dto);
    if (completer != null) return completer!.future;
    final configuredResults = results;
    if (configuredResults == null || configuredResults.isEmpty) return result;
    final index = transferCalls - 1;
    return configuredResults[index < configuredResults.length
        ? index
        : configuredResults.length - 1];
  }
}
