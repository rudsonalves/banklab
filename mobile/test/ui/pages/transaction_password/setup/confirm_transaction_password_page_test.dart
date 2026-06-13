import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/routing/models/transaction_password_setup_origin.dart';
import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/enums/transaction_password_status.dart'
    as api;
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/user/enums/user_role.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/viewmodel/transaction_password_viewmodel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('navigates to Home after post-login setup', (tester) async {
    final harness = await _pumpPage(
      tester,
      origin: TransactionPasswordSetupOrigin.postLogin,
      result: Success(_activeResponse()),
    );

    await _submit(tester);

    expect(find.text('Home destination'), findsOneWidget);
    expect(harness.repository.createCalls, 1);
    expect(
      harness.appSection.readiness?.transactionPasswordStatus,
      TransactionPasswordStatus.active,
    );
  });

  testWidgets('navigates to recipient after transfer setup', (tester) async {
    final harness = await _pumpPage(
      tester,
      origin: TransactionPasswordSetupOrigin.transfer,
      result: Success(_activeResponse()),
    );

    await _submit(tester);

    expect(find.text('Recipient destination'), findsOneWidget);
    expect(harness.repository.createCalls, 1);
  });

  testWidgets('does not open transfer when local status is not active', (
    tester,
  ) async {
    final harness = await _pumpPage(
      tester,
      origin: TransactionPasswordSetupOrigin.transfer,
      result: Success(_activeResponse()),
      updateStatusOnSuccess: false,
    );

    await _submit(tester);

    expect(find.text('Recipient destination'), findsNothing);
    expect(find.text('Confirme sua senha transacional'), findsOneWidget);
    expect(
      find.text('Não foi possível ativar a senha transacional.'),
      findsOneWidget,
    );
    expect(
      harness.appSection.readiness?.transactionPasswordStatus,
      TransactionPasswordStatus.notSet,
    );
  });

  testWidgets(
    'keeps the setup open after failure and permits another attempt',
    (
      tester,
    ) async {
      final harness = await _pumpPage(
        tester,
        origin: TransactionPasswordSetupOrigin.postLogin,
        result: const Failure(
          AppError(
            code: AppErrorCode.networkError,
            message: 'Não foi possível criar a senha.',
          ),
        ),
      );

      await _submit(tester);

      expect(find.text('Confirme sua senha transacional'), findsOneWidget);
      expect(find.text('Não foi possível criar a senha.'), findsOneWidget);

      harness.repository.result = Success(_activeResponse());
      await tester.tap(find.text('Concluir'));
      await tester.pumpAndSettle();

      expect(find.text('Home destination'), findsOneWidget);
      expect(harness.repository.createCalls, 2);
    },
  );

  testWidgets(
    'keeps setup and session unchanged when password already exists',
    (
      tester,
    ) async {
      final harness = await _pumpPage(
        tester,
        origin: TransactionPasswordSetupOrigin.postLogin,
        result: const Failure(
          AppError(
            code: AppErrorCode.transactionPasswordAlreadySet,
            message: 'Backend message',
          ),
        ),
      );

      await _submit(tester);

      expect(find.text('Confirme sua senha transacional'), findsOneWidget);
      expect(
        find.text('Sua senha transacional já está cadastrada.'),
        findsOneWidget,
      );
      expect(
        harness.appSection.readiness?.transactionPasswordStatus,
        TransactionPasswordStatus.notSet,
      );
    },
  );
}

Future<_Harness> _pumpPage(
  WidgetTester tester, {
  required TransactionPasswordSetupOrigin origin,
  required Result<TransactionPasswordStatusResponseDto> result,
  bool updateStatusOnSuccess = true,
}) async {
  final appSection = AppSection()
    ..setAuthSession(_session(TransactionPasswordStatus.notSet));
  final repository = _FakeTransactionPasswordRepository(
    appSection: appSection,
    result: result,
    updateStatusOnSuccess: updateStatusOnSuccess,
  );
  final viewModel = TransactionPasswordViewModel(
    repository: repository,
    appSection: appSection,
  );
  final router = GoRouter(
    initialLocation: TransactionPasswordRoutes.confirm.routePath,
    routes: [
      GoRoute(
        path: TransactionPasswordRoutes.confirm.routePath,
        name: TransactionPasswordRoutes.confirm.routeName,
        builder: (context, state) => ConfirmTransactionPasswordPage(
          token: '123456',
          origin: origin,
          viewModel: viewModel,
        ),
      ),
      GoRoute(
        path: BaseRoutes.home.routePath,
        name: BaseRoutes.home.routeName,
        builder: (context, state) =>
            const Scaffold(body: Text('Home destination')),
      ),
      GoRoute(
        path: TransferRoutes.recipient.routePath,
        name: TransferRoutes.recipient.routeName,
        builder: (context, state) =>
            const Scaffold(body: Text('Recipient destination')),
      ),
    ],
  );
  addTearDown(router.dispose);

  await tester.pumpWidget(MaterialApp.router(routerConfig: router));
  await tester.pumpAndSettle();

  return _Harness(appSection: appSection, repository: repository);
}

Future<void> _submit(WidgetTester tester) async {
  await tester.enterText(
    find.byType(TextField, skipOffstage: false),
    '123456',
  );
  await tester.pump();
  await tester.tap(find.text('Concluir'));
  await tester.pumpAndSettle();
}

TransactionPasswordStatusResponseDto _activeResponse() {
  return TransactionPasswordStatusResponseDto(
    userId: 'user-1',
    status: api.TransactionPasswordStatus.active,
    createdAt: DateTime(2026, 6, 12),
  );
}

AuthSession _session(TransactionPasswordStatus status) {
  return AuthSession(
    user: UserSession(
      userId: 'user-1',
      email: 'customer@example.com',
      role: UserRole.customer,
    ),
    customer: CustomerSession(
      id: 'customer-1',
      name: 'Maria Silva',
      cpf: '12345678901',
      birthDate: DateTime(1990),
      createdAt: DateTime(2026, 6, 12),
    ),
    readiness: ReadinessSession(
      onboardingCompleted: true,
      approved: true,
      hasOperationalAccount: true,
      transactionPasswordStatus: status,
    ),
  );
}

class _Harness {
  final AppSection appSection;
  final _FakeTransactionPasswordRepository repository;

  const _Harness({
    required this.appSection,
    required this.repository,
  });
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  final AppSection appSection;
  final bool updateStatusOnSuccess;
  Result<TransactionPasswordStatusResponseDto> result;
  int createCalls = 0;

  _FakeTransactionPasswordRepository({
    required this.appSection,
    required this.result,
    required this.updateStatusOnSuccess,
  });

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) async {
    createCalls++;
    if (result.isSuccess && updateStatusOnSuccess) {
      appSection.markTransactionPasswordAsActive();
    }
    return result;
  }

  @override
  AsyncResult<StepUpAuthorizeResponseDto> stepUpAuthorize(
    StepUpAuthorizeRequestDto dto,
  ) {
    throw UnimplementedError();
  }
}
