import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/core/routing/routes/transaction_password_routes.dart';
import 'package:bankflow/core/routing/routes/transfer_routes.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  for (final path in [
    TransactionPasswordRoutes.introduction.routePath,
    TransactionPasswordRoutes.create.routePath,
    TransactionPasswordRoutes.confirm.routePath,
    TransferRoutes.payment.routePath,
    TransferRoutes.confirmation.routePath,
  ]) {
    testWidgets('$path redirects missing extras to Home', (tester) async {
      final router = GoRouter(
        initialLocation: path,
        routes: [
          ...transactionPasswordRoutes(),
          ...transferRoutes(),
          _homeRoute(),
        ],
      );
      addTearDown(router.dispose);

      await tester.pumpWidget(MaterialApp.router(routerConfig: router));
      await tester.pumpAndSettle();

      expect(find.text('Safe Home'), findsOneWidget);
    });
  }
}

GoRoute _homeRoute() {
  return GoRoute(
    path: BaseRoutes.home.routePath,
    name: BaseRoutes.home.routeName,
    builder: (context, state) => const Scaffold(body: Text('Safe Home')),
  );
}
