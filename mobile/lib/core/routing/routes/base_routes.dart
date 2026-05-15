import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/ui/pages/home/home_page.dart';
import '/ui/pages/home/viewmodel/home_viewmodel.dart';
import '/ui/pages/splash/splash_page.dart';
import '/ui/pages/splash/viewmodel/splash_viewmodel.dart';
import '/ui/pages/statement/statement_page.dart';
import '/ui/pages/statement/viewmodel/statement_viewmodel.dart';

List<RouteBase> baseRoutes() => [
  GoRoute(
    path: BaseRoutes.home.path,
    name: BaseRoutes.home.name,
    builder: (context, state) => HomePage(
      viewModel: injector.get<HomeViewmodel>(),
    ),
  ),

  GoRoute(
    path: BaseRoutes.splash.path,
    name: BaseRoutes.splash.name,
    builder: (context, state) => SplashPage(
      viewModel: injector.get<SplashViewmodel>(),
    ),
  ),

  GoRoute(
    path: BaseRoutes.statement.path,
    name: BaseRoutes.statement.name,
    builder: (context, state) =>
        StatementPage(viewModel: injector.get<StatementViewmodel>()),
  ),
];
