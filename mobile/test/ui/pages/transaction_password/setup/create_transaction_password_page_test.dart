import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/routing/extra_codec.dart';
import 'package:bankflow/core/routing/models/transaction_password_setup_origin.dart';
import 'package:bankflow/core/routing/routes.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/confirm_transaction_password_page.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/create_transaction_password_page.dart';
import 'package:bankflow/ui/pages/transaction_password/setup/viewmodel/transaction_password_viewmodel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('propagates the PIN and origin only in memory', (tester) async {
    final viewModel = TransactionPasswordViewModel(
      repository: _FakeTransactionPasswordRepository(),
      appSection: AppSection(),
    );
    final router = GoRouter(
      initialLocation: TransactionPasswordRoutes.create.routePath,
      extraCodec: const ExtraCodec(),
      routes: [
        GoRoute(
          path: TransactionPasswordRoutes.create.routePath,
          name: TransactionPasswordRoutes.create.routeName,
          builder: (context, state) => CreateTransactionPasswordPage(
            origin: TransactionPasswordSetupOrigin.transfer,
            viewModel: viewModel,
          ),
        ),
      ],
    );
    addTearDown(router.dispose);

    await tester.pumpWidget(MaterialApp.router(routerConfig: router));
    await tester.pumpAndSettle();

    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '123456',
    );
    await tester.pump();
    await tester.tap(find.text('Continuar'));
    await tester.pumpAndSettle();

    final confirmation = tester.widget<ConfirmTransactionPasswordPage>(
      find.byType(ConfirmTransactionPasswordPage),
    );
    expect(confirmation.token, '123456');
    expect(confirmation.origin, TransactionPasswordSetupOrigin.transfer);
  });
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) {
    throw UnimplementedError();
  }

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInternalTransfer(
    String transactionPassword,
  ) {
    throw UnimplementedError();
  }
}
