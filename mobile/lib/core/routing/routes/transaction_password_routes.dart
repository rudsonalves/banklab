import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/animations_page/app_custom_transaction.dart';
import '/core/routing/routes.dart';
import '/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart';
import '/ui/pages/transaction_password/setup/create_transaction_password_page.dart';
import '/ui/pages/transaction_password/setup/transaction_password_introduction_page.dart';
import '/ui/pages/transaction_password/setup/viewmodel/transaction_password_viewmodel.dart';

List<RouteBase> transactionPasswordRoutes() => [
  GoRoute(
    path: TransactionPasswordRoutes.introduction.routePath,
    name: TransactionPasswordRoutes.introduction.routeName,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: const TransactionPasswordIntroductionPage(),
    ),
  ),
  GoRoute(
    path: TransactionPasswordRoutes.create.routePath,
    name: TransactionPasswordRoutes.create.routeName,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: const CreateTransactionPasswordPage(),
    ),
  ),
  GoRoute(
    path: TransactionPasswordRoutes.confirm.routePath,
    name: TransactionPasswordRoutes.confirm.routeName,
    pageBuilder: (context, state) {
      final pin = state.extra;
      if (pin is! String || pin.length != 6) {
        return AppCustomTransactionPage(
          key: state.pageKey,
          child: const TransactionPasswordIntroductionPage(),
        );
      }

      return AppCustomTransactionPage(
        key: state.pageKey,
        child: ConfirmTransactionPasswordPage(
          viewModel: injector.get<TransactionPasswordViewModel>(),
          pin: pin,
        ),
      );
    },
  ),
];
