import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/transaction_password_introduction_page.dart';
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
              const TransactionPasswordIntroductionPage(),
        ),
        GoRoute(
          path: TransactionPasswordRoutes.create.routePath,
          name: TransactionPasswordRoutes.create.name,
          builder: (context, state) =>
              const Scaffold(body: Text('Create password page')),
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

    expect(find.text('Create password page'), findsOneWidget);
  });
}
