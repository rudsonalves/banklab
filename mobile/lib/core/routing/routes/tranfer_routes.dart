import 'package:go_router/go_router.dart';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/uis/pages/home/transfer/payment_page.dart';
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

  GoRoute(
    path: TransferRoutes.payment.path,
    name: TransferRoutes.payment.name,
    builder: (context, state) => PaymentPage(
      viewModel: injector.get<TransferViewmodel>(),
      recipientInfo: state.extra as RecipientInfoDto,
    ),
  ),
];
