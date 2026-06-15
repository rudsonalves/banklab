import 'package:go_router/go_router.dart';

import '/ui/pages/shared/details/details_page.dart';
import '/ui/pages/shared/details/viewmodel/details_viewmodel.dart';
import '../../config/dependencies.dart';
import '../routes.dart';

List<RouteBase> sharedRoutes() => [
  GoRoute(
    path: SharedRoutes.details.routePath,
    name: SharedRoutes.details.routeName,
    builder: (context, state) => DetailsPage(
      viewModel: injector.get<DetailsViewmodel>(),
      reference: state.extra as String,
    ),
  ),
];
