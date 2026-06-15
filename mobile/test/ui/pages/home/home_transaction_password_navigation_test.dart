import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/routing/models/transaction_password_setup_origin.dart';
import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/data/repositories/account/account_repository.dart';
import 'package:bankflow/data/services/apis/account/dtos/account_summary_response_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/balance_response_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/statement_query_params_dto.dart';
import 'package:bankflow/data/services/apis/account/dtos/statement_response_dto.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/user/enums/user_role.dart';
import 'package:bankflow/ui/pages/home/home_page.dart';
import 'package:bankflow/ui/pages/home/viewmodel/home_viewmodel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('active opens the recipient selection', (tester) async {
    await _pumpHome(tester, TransactionPasswordStatus.active);

    await tester.tap(find.text('Transferir'));
    await tester.pumpAndSettle();

    expect(find.text('Recipient page'), findsOneWidget);
  });

  testWidgets('notSet opens setup with transfer origin', (tester) async {
    await _pumpHome(tester, TransactionPasswordStatus.notSet);

    await tester.tap(find.text('Transferir'));
    await tester.pumpAndSettle();

    expect(find.text('Setup origin: transfer'), findsOneWidget);
  });

  testWidgets('locked stays on Home and presents the blocked message', (
    tester,
  ) async {
    await _pumpHome(tester, TransactionPasswordStatus.locked);

    await tester.tap(find.text('Transferir'));
    await tester.pump();

    expect(find.text('Minha conta'), findsOneWidget);
    expect(find.text('Recipient page'), findsNothing);
    expect(find.text('Setup origin: transfer'), findsNothing);
    expect(find.text('Acesso bloqueado'), findsOneWidget);
  });

  testWidgets('unknown stays on Home and presents an error', (tester) async {
    await _pumpHome(tester, TransactionPasswordStatus.unknown);

    await tester.tap(find.text('Transferir'));
    await tester.pump();

    expect(find.text('Minha conta'), findsOneWidget);
    expect(find.text('Recipient page'), findsNothing);
    expect(find.text('Setup origin: transfer'), findsNothing);
    expect(
      find.text(
        'Não foi possível verificar o status da sua senha transacional.'
        ' Por favor, tente novamente mais tarde.',
      ),
      findsOneWidget,
    );
  });
}

Future<void> _pumpHome(
  WidgetTester tester,
  TransactionPasswordStatus status,
) async {
  final appSection = AppSection()..setAuthSession(_session(status));
  final viewModel = HomeViewmodel(
    accountRepository: _FakeAccountRepository(),
    appSection: appSection,
  );
  final router = GoRouter(
    initialLocation: BaseRoutes.home.routePath,
    routes: [
      GoRoute(
        path: BaseRoutes.home.routePath,
        name: BaseRoutes.home.routeName,
        builder: (context, state) => HomePage(viewModel: viewModel),
      ),
      GoRoute(
        path: TransferRoutes.recipient.routePath,
        name: TransferRoutes.recipient.routeName,
        builder: (context, state) =>
            const Scaffold(body: Text('Recipient page')),
      ),
      GoRoute(
        path: TransactionPasswordRoutes.introduction.routePath,
        name: TransactionPasswordRoutes.introduction.routeName,
        builder: (context, state) {
          final origin = state.extra as TransactionPasswordSetupOrigin;
          return Scaffold(body: Text('Setup origin: ${origin.name}'));
        },
      ),
    ],
  );
  addTearDown(router.dispose);

  await tester.pumpWidget(MaterialApp.router(routerConfig: router));
  await tester.pumpAndSettle();
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

class _FakeAccountRepository implements AccountRepository {
  @override
  List<AccountSummaryResponseDto>? get accounts => const [];

  @override
  BalanceResponseDto? get lastBalance => null;

  @override
  AccountSummaryResponseDto? get selectedAccount => null;

  @override
  StatementResponseDto? get lastStatement => null;

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
      const Success([]);

  @override
  AsyncResult<Unit> selectAccount(String accountId) async =>
      const Success(unit);
}
