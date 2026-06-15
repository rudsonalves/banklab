import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/core/routing/models/transaction_password_setup_origin.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/introduction_transaction_password_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('explains the transaction password before creation', (
    tester,
  ) async {
    final router = GoRouter(
      initialLocation: TransactionPasswordRoutes.introduction.routePath,
      routes: [
        GoRoute(
          path: TransactionPasswordRoutes.introduction.routePath,
          name: TransactionPasswordRoutes.introduction.name,
          builder: (context, state) =>
              const IntroductionTransactionPasswordPage(
                origin: TransactionPasswordSetupOrigin.postLogin,
              ),
        ),
        GoRoute(
          path: TransactionPasswordRoutes.create.routePath,
          name: TransactionPasswordRoutes.create.name,
          builder: (context, state) {
            final origin = state.extra as TransactionPasswordSetupOrigin;
            return Scaffold(body: Text('Origin: ${origin.name}'));
          },
        ),
        GoRoute(
          path: AuthRoutes.login.routePath,
          name: AuthRoutes.login.name,
          builder: (context, state) => const Scaffold(body: Text('Login page')),
        ),
      ],
    );
    addTearDown(router.dispose);

    await tester.pumpWidget(
      MaterialApp.router(routerConfig: router),
    );
    await tester.pumpAndSettle();

    expect(find.text('Proteção extra para suas transações'), findsOneWidget);
    expect(find.text('É diferente da senha de acesso'), findsOneWidget);
    expect(find.text('Confirma operações financeiras'), findsOneWidget);
    expect(find.text('É pessoal e intransferível'), findsOneWidget);

    await tester.tap(find.text('Criar senha'));
    await tester.pumpAndSettle();

    expect(find.text('Origin: postLogin'), findsOneWidget);
  });

  testWidgets('returns to the page that started the setup when cancelled', (
    tester,
  ) async {
    final router = GoRouter(
      initialLocation: '/host',
      routes: [
        GoRoute(
          path: '/host',
          builder: (context, state) => Scaffold(
            body: Column(
              children: [
                const Text('Host page'),
                TextButton(
                  onPressed: () => context.pushNamed(
                    TransactionPasswordRoutes.introduction.routeName,
                    extra: TransactionPasswordSetupOrigin.transfer,
                  ),
                  child: const Text('Start setup'),
                ),
              ],
            ),
          ),
        ),
        GoRoute(
          path: TransactionPasswordRoutes.introduction.routePath,
          name: TransactionPasswordRoutes.introduction.routeName,
          builder: (context, state) => IntroductionTransactionPasswordPage(
            origin: state.extra as TransactionPasswordSetupOrigin,
          ),
        ),
      ],
    );
    addTearDown(router.dispose);

    await tester.pumpWidget(MaterialApp.router(routerConfig: router));
    await tester.tap(find.text('Start setup'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Agora não'));
    await tester.pumpAndSettle();

    expect(find.text('Host page'), findsOneWidget);
  });

  test('rejects an arbitrary setup origin', () {
    expect(
      () => TransactionPasswordSetupOrigin.fromName('/recipient'),
      throwsUnsupportedError,
    );
  });
}
