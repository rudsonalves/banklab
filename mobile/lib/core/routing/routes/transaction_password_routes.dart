import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/animations_page/app_custom_transaction.dart';
import '/core/routing/routes.dart';
import '/ui/pages/auth/transaction_password/confirm_transaction_password_page.dart';
import '/ui/pages/auth/transaction_password/create_transaction_password_page.dart';
import '/ui/pages/auth/transaction_password/viewmodel/transaction_password_viewmodel.dart';

List<RouteBase> transactionPasswordRoutes() => [
  GoRoute(
    path: TransactionPasswordRoutes.create.path,
    name: TransactionPasswordRoutes.create.name,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: const CreateTransactionPasswordPage(),
    ),
  ),
  GoRoute(
    path: TransactionPasswordRoutes.confirm.path,
    name: TransactionPasswordRoutes.confirm.name,
    pageBuilder: (context, state) {
      final pin = state.extra;
      if (pin is! String || pin.length != 6) {
        return AppCustomTransactionPage(
          key: state.pageKey,
          child: const CreateTransactionPasswordPage(),
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
