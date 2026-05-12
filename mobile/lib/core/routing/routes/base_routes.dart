import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/uis/pages/home/home_page.dart';
import '/uis/pages/home/viewmodel/home_viewmodel.dart';
import '/uis/pages/splash/splash_page.dart';
import '/uis/pages/splash/viewmodel/splash_viewmodel.dart';

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
];
