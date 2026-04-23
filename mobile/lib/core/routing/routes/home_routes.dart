import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/uis/pages/home/home_page.dart';
import '/uis/pages/home/viewmodel/home_viewmodel.dart';

List<RouteBase> homeRoutes() => [
  GoRoute(
    path: HomeRoutes.home.path,
    name: HomeRoutes.home.name,
    builder: (context, state) => HomePage(
      viewModel: injector.get<HomeViewmodel>(),
    ),
  ),
];
