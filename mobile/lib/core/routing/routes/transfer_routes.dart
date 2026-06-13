import 'package:go_router/go_router.dart';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/ui/pages/transaction_password/verification/transaction_password_input_page.dart';
import '/ui/pages/transfer/models/transfer_confirmation_data.dart';
import '/ui/pages/transfer/transfer_confirmation_page.dart';
import '/ui/pages/transfer/transfer_payment_page.dart';
import '/ui/pages/transfer/transfer_recipient_page.dart';
import '/ui/pages/transfer/transfer_status_page.dart';
import '/ui/pages/transfer/viewmodel/transfer_viewmodel.dart';
import '../../config/dependencies.dart';
import '../routes.dart';

List<RouteBase> transferRoutes() => [
  GoRoute(
    path: TransferRoutes.recipient.routePath,
    name: TransferRoutes.recipient.routeName,
    builder: (context, state) => TransferRecipientPage(
      viewModel: injector.get<TransferViewmodel>(),
    ),
  ),

  GoRoute(
    path: TransferRoutes.payment.routePath,
    name: TransferRoutes.payment.routeName,
    builder: (context, state) => TransferPaymentPage(
      viewModel: injector.get<TransferViewmodel>(),
      recipientInfo: state.extra as RecipientInfoDto,
    ),
  ),

  GoRoute(
    path: TransferRoutes.confirmation.routePath,
    name: TransferRoutes.confirmation.routeName,
    builder: (context, state) => TransferConfirmationPage(
      viewModel: injector.get<TransferViewmodel>(),
      transferData: state.extra as TransferConfirmationData,
    ),
  ),

  GoRoute(
    path: TransactionPasswordRoutes.transactionPassword.routePath,
    name: TransactionPasswordRoutes.transactionPassword.routeName,
    builder: (context, state) => const TransactionPasswordInputPage(),
  ),

  GoRoute(
    path: TransferRoutes.statusSuccess.routePath,
    name: TransferRoutes.statusSuccess.routeName,
    builder: (context, state) => TransferStatusPage(
      isSuccess: true,
      transactionReference: state.extra as String,
    ),
  ),

  GoRoute(
    path: TransferRoutes.statusFailure.routePath,
    name: TransferRoutes.statusFailure.routeName,
    builder: (context, state) => TransferStatusPage(
      isSuccess: false,
    ),
  ),
];
