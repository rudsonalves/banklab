import 'package:go_router/go_router.dart';

import '/uis/pages/shared/details/details_page.dart';
import '/uis/pages/shared/details/viewmodel/details_viewmodel.dart';
import '../../config/dependencies.dart';
import '../routes.dart';

List<RouteBase> sharedRoutes() => [
  GoRoute(
    path: SharedRoutes.details.path,
    name: SharedRoutes.details.name,
    builder: (context, state) => DetailsPage(
      viewModel: injector.get<DetailsViewmodel>(),
      reference: state.extra as String,
    ),
  ),
];
