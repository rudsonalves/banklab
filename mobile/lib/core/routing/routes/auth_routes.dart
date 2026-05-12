import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/data/services/auth/cache/models/last_login_identity.dart';
import '/uis/pages/auth/login/login_page.dart';
import '/uis/pages/auth/login/viewmodel/login_viewmodel.dart';
import '/uis/pages/auth/register/register_page.dart';
import '/uis/pages/auth/register/viewmodel/register_viewmodel.dart';
import '/uis/pages/auth/short_login/short_login_page.dart';
import '/uis/pages/auth/short_login/viewmodel/short_login_viewmodel.dart';
import '../animations_page/app_custom_transaction.dart';

List<RouteBase> authRoutes() => [
  GoRoute(
    path: AuthRoutes.login.path,
    name: AuthRoutes.login.name,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: LoginPage(viewModel: injector.get<LoginViewModel>()),
    ),
  ),
  GoRoute(
    path: AuthRoutes.register.path,
    name: AuthRoutes.register.name,
    pageBuilder: (context, state) => AppCustomTransactionPage(
      key: state.pageKey,
      child: RegisterPage(viewmodel: injector.get<RegisterViewmodel>()),
    ),
  ),
  GoRoute(
    path: AuthRoutes.shortLogin.path,
    name: AuthRoutes.shortLogin.name,
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
          viewModel: injector.get<ShortLoginViewModel>(),
          identity: identity,
        ),
      );
    },
  ),
];
