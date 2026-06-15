import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:uuid/uuid.dart';

import '/core/result/errors/backend_error_code.dart';
import '/core/routing/routes.dart';
import '/core/services/logging/console_log.dart';
import '/domain/usecases/transfer/inputs/protected_transfer_input.dart';
import '/domain/usecases/transfer/transfer_usecase.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/big_button.dart';
import '/ui/components/cards/balance_card.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/card_text_row.dart';
import '/ui/components/text/text_header.dart';
import 'models/transfer_confirmation_data.dart';
import 'viewmodel/transfer_viewmodel.dart';

class TransferConfirmationPage extends StatefulWidget {
  final TransferViewmodel viewModel;
  final TransferConfirmationData transferData;

  const TransferConfirmationPage({
    super.key,
    required this.viewModel,
    required this.transferData,
  });

  @override
  State<TransferConfirmationPage> createState() =>
      _TransferConfirmationPageState();
}

class _TransferConfirmationPageState extends State<TransferConfirmationPage> {
  TransferViewmodel get _viewModel => widget.viewModel;
  TransferConfirmationData get _transferData => widget.transferData;

  final _log = ConsoleLog('TransferConfirmationPage');
  late final String _idempotencyKey;
  final _hasSubmitted = ValueNotifier<bool>(false);

  @override
  void initState() {
    super.initState();
    _idempotencyKey = const Uuid().v7();
  }

  @override
  void dispose() {
    _hasSubmitted.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Confirmação'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            BalanceCard(
              balance: _viewModel.balance,
              isVisible: true,
              selectedAccount: _viewModel.selectedAccount,
            ),

            SizedBox(height: 16),
            TextHeader('Detalhes da transferência'),
            Card(
              color: Theme.of(context).colorScheme.onPrimary,
              margin: EdgeInsets.zero,
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  spacing: 12,
                  crossAxisAlignment: CrossAxisAlignment.stretch,
                  children: [
                    CardTextRow(
                      label: 'Destinatário',
                      value: _transferData.toHolderName,
                    ),
                    // CardTextRow(label: 'Banco:', value: transferData.toBankName),
                    CardTextRow(
                      label: 'Conta',
                      value: '0001 - ${_transferData.toNumber}',
                    ),
                    CardTextRow(
                      label: 'Valor',
                      value: _transferData.amount.format(),
                    ),
                    if (_transferData.description.isNotEmpty)
                      CardTextRow(
                        label: 'Descrição',
                        value: _transferData.description,
                      ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),

      bottomNavigationBar: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16),
        child: ListenableBuilder(
          listenable: Listenable.merge([
            _viewModel.transfer,
            _hasSubmitted,
          ]),
          builder: (context, _) => BigButton(
            label: 'Transferir',
            onPressed: _submit,
            leftIcon: Icon(Icons.check_rounded, size: 24),
            enabled: !_hasSubmitted.value,
            isRunning: _viewModel.transfer.isRunning,
          ),
        ),
      ),
    );
  }

  Future<void> _submit() async {
    if (_hasSubmitted.value) return;
    _hasSubmitted.value = true;

    final draft = TransferDraft(
      toAccountId: _transferData.toAccountId,
      description: _transferData.description,
      amount: _transferData.amount,
      idempotencyKey: _idempotencyKey,
    );

    await _executeProtectedTransfer(draft);
  }

  Future<void> _executeProtectedTransfer(TransferDraft draft) async {
    final transferCommand = _viewModel.transfer;

    final pin = await context.pushNamed<String?>(
      TransactionPasswordRoutes.transactionPassword.routeName,
    );

    if (!mounted) return;
    if (pin == null) {
      AppSnackbar.show(
        context,
        message: 'Operação cancelada. Senha de transação não fornecida.',
        type: SnackbarType.info,
      );
      _finishSubmit();
      return;
    }

    await transferCommand.execute(
      ProtectedTransferInput(draft: draft, pin: pin),
    );

    final rawError = transferCommand.error;
    final mappedErrorCode = backendErrorCode(rawError);
    _log.info(
      'Transfer command completed: '
      'state=${transferCommand.state.name}, '
      'isSuccess=${transferCommand.isSuccess}, '
      'isFailure=${transferCommand.isFailure}, '
      'errorCode=${rawError?.code.name}, '
      'backendErrorCode=$mappedErrorCode, '
      'errorMessage=${rawError?.message}, '
      'hasValue=${transferCommand.result?.value != null}',
    );

    if (!mounted) return;
    if (transferCommand.isSuccess) {
      final transferResponse = transferCommand.result?.value;
      if (transferResponse == null) {
        AppSnackbar.show(
          context,
          message: 'Erro desconhecido. Por favor, tente novamente mais tarde.',
          type: SnackbarType.error,
        );
        _finishSubmit();
        return;
      }

      context.pushNamed(
        TransferRoutes.statusSuccess.routeName,
        extra: transferResponse.transactionReference,
      );
      _finishSubmit();
      return;
    }

    final errorCode = mappedErrorCode;

    switch (errorCode) {
      case 'TRANSACTION_PASSWORD_INVALID':
        _showError('Senha de transação inválida. Tente novamente.');
        return _executeProtectedTransfer(draft);
      case 'STEP_UP_TOKEN_EXPIRED':
      case 'STEP_UP_TOKEN_CONSUMED':
        AppSnackbar.show(
          context,
          message:
              'A autorização expirou. Informe novamente sua senha de transação.',
          type: SnackbarType.info,
        );
        return _executeProtectedTransfer(draft);
      case 'TRANSACTION_PASSWORD_LOCKED':
        _showError(
          'Senha de transação bloqueada. Tente novamente mais tarde.',
        );
        return;
      case 'TRANSACTION_PASSWORD_NOT_SET':
        _showError(
          'Sua senha de transação precisa ser configurada novamente.',
        );
        context.goNamed(BaseRoutes.home.routeName);
        _finishSubmit();
        return;
      default:
        _log.warn(
          'Routing to failure status page. '
          'Unhandled backendErrorCode=$errorCode, '
          'appErrorCode=${transferCommand.error?.code.name}, '
          'message=${transferCommand.error?.message}',
        );
        context.pushNamed(TransferRoutes.statusFailure.routeName);
        _finishSubmit();
        return;
    }
  }

  void _finishSubmit() {
    if (!mounted) return;
    _hasSubmitted.value = false;
  }

  void _showError(String message) {
    AppSnackbar.show(
      context,
      message: message,
      type: SnackbarType.error,
    );
  }
}
