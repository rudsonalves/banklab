import 'package:go_router/go_router.dart';

import '/uis/pages/home/transfer/recipient_page.dart';
import '/uis/pages/home/transfer/viewmodel/transfer_viewmodel.dart';
import '../../config/dependencies.dart';
import '../routes.dart';

List<RouteBase> transferRoutes() => [
  GoRoute(
    path: TransferRoutes.recipient.path,
    name: TransferRoutes.recipient.name,
    builder: (context, state) => RecipientPage(
      viewModel: injector.get<TransferViewmodel>(),
    ),
  ),
];
