import 'package:bankflow/core/routing/models/transaction_password_setup_origin.dart';
import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/create_transaction_password_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('propagates the PIN and origin to confirmation', (tester) async {
    final router = GoRouter(
      initialLocation: TransactionPasswordRoutes.create.routePath,
      routes: [
        GoRoute(
          path: TransactionPasswordRoutes.create.routePath,
          name: TransactionPasswordRoutes.create.routeName,
          builder: (context, state) => const CreateTransactionPasswordPage(
            origin: TransactionPasswordSetupOrigin.transfer,
          ),
        ),
        GoRoute(
          path: TransactionPasswordRoutes.confirm.routePath,
          name: TransactionPasswordRoutes.confirm.routeName,
          builder: (context, state) {
            final extra = state.extra! as Map<String, dynamic>;
            return Scaffold(
              body: Text('${extra['token']}:${extra['origin']}'),
            );
          },
        ),
      ],
    );
    addTearDown(router.dispose);

    await tester.pumpWidget(MaterialApp.router(routerConfig: router));
    await tester.pumpAndSettle();

    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '123456',
    );
    await tester.pump();
    await tester.tap(find.text('Continuar'));
    await tester.pumpAndSettle();

    expect(find.text('123456:transfer'), findsOneWidget);
  });
}
