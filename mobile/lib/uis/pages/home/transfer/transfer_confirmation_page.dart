import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/domain/usecases/transfer/transfer_usecase.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '/uis/core/cards/balance_card.dart';
import '/uis/core/text/card_text_row.dart';
import '/uis/core/text/text_header.dart';
import '../../../core/messages/app_snackbar.dart';
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
              balance: widget.viewModel.balance,
              isVisible: true,
              selectedAccount: widget.viewModel.selectedAccount,
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
                      value: widget.transferData.toHolderName,
                    ),
                    // CardTextRow(label: 'Banco:', value: transferData.toBankName),
                    CardTextRow(
                      label: 'Conta',
                      value: '0001 - ${widget.transferData.toNumber}',
                    ),
                    CardTextRow(
                      label: 'Valor',
                      value: widget.transferData.amount.format(),
                    ),
                    if (widget.transferData.description.isNotEmpty)
                      CardTextRow(
                        label: 'Descrição',
                        value: widget.transferData.description,
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
          onPressed: _onConfirmTransfer,
          leftIcon: Icon(Icons.check_rounded, size: 24),
          enabled: true,
        ),
      ),
    );
  }

  Future<void> _onConfirmTransfer() async {
    final transfer = TransferDraft(
      toAccountId: widget.transferData.toAccountId,
      description: widget.transferData.description,
      amount: widget.transferData.amount,
    );

    await widget.viewModel.transfer.execute(transfer);

    if (!mounted) return;
    if (widget.viewModel.transfer.isFailure) {
      context.pushNamed(TransferRoutes.statusFailure.name);
    } else {
      final transferResponse = widget.viewModel.transfer.result?.value;
      if (transferResponse == null) {
        AppSnackbar.show(
          context,
          message: 'Erro desconhecido. Por favor, tente novamente mais tarde.',
          type: SnackbarType.error,
        );
        return;
      }

      context.pushNamed(
        TransferRoutes.statusSuccess.name,
        extra: transferResponse.transactionReference,
      );
    }
  }
}
