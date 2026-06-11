import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/animations_page/app_custom_transaction.dart';
import '/core/routing/routes.dart';
import '../../../ui/pages/transaction_password/creation_flow/confirm_transaction_password_page.dart';
import '../../../ui/pages/transaction_password/creation_flow/create_transaction_password_page.dart';
import '../../../ui/pages/transaction_password/creation_flow/transaction_password_introduction_page.dart';
import '../../../ui/pages/transaction_password/creation_flow/viewmodel/transaction_password_viewmodel.dart';

List<RouteBase> transactionPasswordRoutes() => [
  GoRoute(
    path: TransactionPasswordRoutes.introduction.path,
    name: TransactionPasswordRoutes.introduction.name,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: const TransactionPasswordIntroductionPage(),
    ),
  ),
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
