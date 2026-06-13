import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:uuid/uuid.dart';

import '/core/routing/routes.dart';
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

  late final String _idempotencyKey;

  @override
  void initState() {
    super.initState();
    _idempotencyKey = const Uuid().v7();
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
        child: BigButton(
          label: 'Transferir',
          onPressed: _submit,
          leftIcon: Icon(Icons.check_rounded, size: 24),
          enabled: true,
        ),
      ),
    );
  }

  Future<void> _submit() async {
    final pin = await context.pushNamed<String?>(
      TransactionPasswordRoutes.transactionPassword.routeName,
    );

    if (pin == null) {
      if (!mounted) return;
      AppSnackbar.show(
        context,
        message: 'Operação cancelada.',
        type: SnackbarType.info,
      );
      return;
    }

    final transfer = ProtectedTransferInput(
      draft: TransferDraft(
        toAccountId: _transferData.toAccountId,
        description: _transferData.description,
        amount: _transferData.amount,
        idempotencyKey: _idempotencyKey,
      ),
      pin: pin,
    );

    await _viewModel.transfer.execute(transfer);

    if (!mounted) return;
    if (_viewModel.transfer.isFailure) {
      context.pushNamed(TransferRoutes.statusFailure.routeName);
    } else {
      final transferResponse = _viewModel.transfer.result?.value;
      if (transferResponse == null) {
        AppSnackbar.show(
          context,
          message: 'Erro desconhecido. Por favor, tente novamente mais tarde.',
          type: SnackbarType.error,
        );
        return;
      }

      context.pushNamed(
        TransferRoutes.statusSuccess.routeName,
        extra: transferResponse.transactionReference,
      );
    }
  }
}
