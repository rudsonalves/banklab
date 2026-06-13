import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/data/services/cache/last_login/models/last_login_identity.dart';
import '/ui/pages/auth/login/login_page.dart';
import '/ui/pages/auth/short_login/short_login_page.dart';
import '/ui/pages/auth/viewmodel/login_viewmodel.dart';
import '../animations_page/app_custom_transaction.dart';

List<RouteBase> authRoutes() => [
  GoRoute(
    path: AuthRoutes.login.routePath,
    name: AuthRoutes.login.routeName,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: LoginPage(viewModel: injector.get<LoginViewModel>()),
    ),
  ),

  GoRoute(
    path: AuthRoutes.shortLogin.routePath,
    name: AuthRoutes.shortLogin.routeName,
    pageBuilder: (context, state) {
      final identity = state.extra;
      if (identity is! LastLoginIdentity) {
        return AppCustomTransactionPage(
          key: state.pageKey,
          child: LoginPage(viewModel: injector.get<LoginViewModel>()),
        );
      }

      return AppCustomTransactionPage(
        key: state.pageKey,
        child: ShortLoginPage(
          viewModel: injector.get<LoginViewModel>(),
          identity: identity,
        ),
      );
    },
  ),
];
