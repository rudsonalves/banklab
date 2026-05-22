import 'package:go_router/go_router.dart';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/ui/pages/home/transfer/models/transfer_confirmation_data.dart';
import '/ui/pages/home/transfer/transfer_confirmation_page.dart';
import '/ui/pages/home/transfer/transfer_payment_page.dart';
import '/ui/pages/home/transfer/transfer_recipient_page.dart';
import '/ui/pages/home/transfer/transfer_status_page.dart';
import '/ui/pages/home/transfer/viewmodel/transfer_viewmodel.dart';
import '../../config/dependencies.dart';
import '../routes.dart';

List<RouteBase> transferRoutes() => [
  GoRoute(
    path: TransferRoutes.recipient.path,
    name: TransferRoutes.recipient.name,
    builder: (context, state) => TransferRecipientPage(
      viewModel: injector.get<TransferViewmodel>(),
    ),
  ),

  GoRoute(
    path: TransferRoutes.payment.path,
    name: TransferRoutes.payment.name,
    builder: (context, state) => TransferPaymentPage(
      viewModel: injector.get<TransferViewmodel>(),
      recipientInfo: state.extra as RecipientInfoDto,
    ),
  ),

  GoRoute(
    path: TransferRoutes.confirmation.path,
    name: TransferRoutes.confirmation.name,
    builder: (context, state) => TransferConfirmationPage(
      viewModel: injector.get<TransferViewmodel>(),
      transferData: state.extra as TransferConfirmationData,
    ),
  ),

  GoRoute(
    path: TransferRoutes.statusSuccess.path,
    name: TransferRoutes.statusSuccess.name,
    builder: (context, state) => TransferStatusPage(
      isSuccess: true,
      transactionReference: state.extra as String,
    ),
  ),

  GoRoute(
    path: TransferRoutes.statusFailure.path,
    name: TransferRoutes.statusFailure.name,
    builder: (context, state) => TransferStatusPage(
      isSuccess: false,
    ),
  ),
];
