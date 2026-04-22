import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/uis/pages/home/home_page.dart';
import '/uis/pages/home/viewmodel/home_viewmodel.dart';

List<RouteBase> homeRoutes() => [
  GoRoute(
    path: HomeRoutes.home.path,
    name: HomeRoutes.home.name,
    builder: (context, state) {
      final extra = state.extra;
      final accountId = switch (extra) {
        String value when value.trim().isNotEmpty => value,
        Map<String, dynamic> value => value['accountId'] as String?,
        _ => state.uri.queryParameters['accountId'],
      };

      return HomePage(
        viewModel: injector.get<HomeViewmodel>(),
        accountId: accountId,
      );
    },
  ),
];
