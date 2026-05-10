import 'package:flutter/material.dart';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '../../../core/cards/balance_card.dart';
import 'viewmodel/transfer_viewmodel.dart';

class PaymentPage extends StatefulWidget {
  final TransferViewmodel viewModel;
  final RecipientInfoDto recipientInfo;

  const PaymentPage({
    super.key,
    required this.viewModel,
    required this.recipientInfo,
  });

  @override
  State<PaymentPage> createState() => _PaymentPageState();
}

class _PaymentPageState extends State<PaymentPage> {
  TransferViewmodel get _viewModel => widget.viewModel;
  RecipientInfoDto get _recipientInfo => widget.recipientInfo;

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Transferência'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            BalanceCard(
              balance: _viewModel.balance,
              isVisible: true,
              onToggleVisibility: () {},
              selectedAccount: _viewModel.selectedAccount,
            ),
          ],
        ),
      ),

      bottomNavigationBar: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16),
        child: ValueListenableBuilder(
          valueListenable: _viewModel.selectedRecipient,
          builder: (context, value, _) {
            final isButtonEnabled = value != null;

            return BigButton(
              label: 'Prosseguir',
              onPressed: _onConfirmTransfer,
              enabled: isButtonEnabled,
            );
          },
        ),
      ),
    );
  }

  void _onConfirmTransfer() {
    Navigator.of(context).pop();
  }
}
