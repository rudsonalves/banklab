import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/animations_page/app_custom_transaction.dart';
import '/core/routing/models/transaction_password_setup_origin.dart';
import '/core/routing/routes.dart';
import '/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart';
import '/ui/pages/transaction_password/setup/create_transaction_password_page.dart';
import '/ui/pages/transaction_password/setup/introduction_transaction_password_page.dart';
import '/ui/pages/transaction_password/setup/viewmodel/transaction_password_viewmodel.dart';

List<RouteBase> transactionPasswordRoutes() => [
  GoRoute(
    path: TransactionPasswordRoutes.introduction.routePath,
    name: TransactionPasswordRoutes.introduction.routeName,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: IntroductionTransactionPasswordPage(
        origin: state.extra as TransactionPasswordSetupOrigin,
      ),
    ),
  ),

  GoRoute(
    path: TransactionPasswordRoutes.create.routePath,
    name: TransactionPasswordRoutes.create.routeName,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: CreateTransactionPasswordPage(
        origin: state.extra as TransactionPasswordSetupOrigin,
      ),
    ),
  ),

  GoRoute(
    path: TransactionPasswordRoutes.confirm.routePath,
    name: TransactionPasswordRoutes.confirm.routeName,
    pageBuilder: (context, state) {
      final map = state.extra as Map<String, dynamic>?;

      final token = map?['token'] as String?;
      final origin = TransactionPasswordSetupOrigin.fromName(
        map?['origin'] as String? ?? 'postLogin',
      );

      if (token is! String || token.length != 6) {
        return AppCustomTransactionPage(
          key: state.pageKey,
          child: IntroductionTransactionPasswordPage(
            origin: TransactionPasswordSetupOrigin.postLogin,
          ),
        );
      }

      return AppCustomTransactionPage(
        key: state.pageKey,
        child: ConfirmTransactionPasswordPage(
          viewModel: injector.get<TransactionPasswordViewModel>(),
          token: token,
          origin: origin,
        ),
      );
    },
  ),
];
