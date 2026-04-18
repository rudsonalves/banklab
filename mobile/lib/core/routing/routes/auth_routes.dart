import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/uis/pages/auth/login/login_page.dart';
import '/uis/pages/auth/login/viewmodel/login_viewmodel.dart';
import '/uis/pages/auth/register/register_page.dart';
import '/uis/pages/auth/register/viewmodel/register_viewmodel.dart';
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
];
