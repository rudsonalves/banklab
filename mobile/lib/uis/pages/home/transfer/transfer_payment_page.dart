import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '/uis/core/cards/balance_card.dart';
import '/uis/core/cards/recipient_card.dart';
import '/uis/core/text/text_header.dart';
import '/uis/core/text_form_field/basic_text_form_field.dart';
import '/uis/core/text_form_field/money_input_formatter.dart';
import '../../../core/messages/app_snackbar.dart';
import 'models/transfer_confirmation_data.dart';
import 'viewmodel/transfer_viewmodel.dart';

class TransferPaymentPage extends StatefulWidget {
  final TransferViewmodel viewModel;
  final RecipientInfoDto recipientInfo;

  const TransferPaymentPage({
    super.key,
    required this.viewModel,
    required this.recipientInfo,
  });

  @override
  State<TransferPaymentPage> createState() => _TransferPaymentPageState();
}

class _TransferPaymentPageState extends State<TransferPaymentPage> {
  TransferViewmodel get _viewModel => widget.viewModel;
  RecipientInfoDto get _recipientInfo => widget.recipientInfo;

  final _amountController = TextEditingController();
  final _descriptionController = TextEditingController();
  final ValueNotifier<bool> _amountIsValid = ValueNotifier(false);

  @override
  void dispose() {
    _amountController.dispose();
    _descriptionController.dispose();
    _amountIsValid.dispose();
    _viewModel.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Transferência'),
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        child: SingleChildScrollView(
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

              SizedBox(height: 12),
              RecipientCard(selectedRecipient: _recipientInfo),

              SizedBox(height: 12),
              TextHeader('Valor'),
              BasicTextFormField(
                controller: _amountController,
                hintText: 'R\$ 0,00',
                inputFormatters: [MoneyInputFormatter()],
                keyboardType: TextInputType.number,
                onChanged: _amountChanged,
              ),

              SizedBox(height: 12),
              TextHeader('Descrição'),
              BasicTextFormField(
                controller: _descriptionController,
                hintText: 'Digite a descrição',
                keyboardType: TextInputType.text,
                textCapitalization: TextCapitalization.sentences,
              ),
            ],
          ),
        ),
      ),

      bottomNavigationBar: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16),
        child: ListenableBuilder(
          listenable: _amountIsValid,
          builder: (context, _) => BigButton(
            label: 'Prosseguir',
            onPressed: _onConfirmTransfer,
            enabled: _amountIsValid.value,
            rightIcon: Icons.arrow_forward_ios_rounded,
          ),
        ),
      ),
    );
  }

  void _amountChanged(String value) {
    final result = value.parseToMoney();

    if (result.isSuccess && result.value!.isPositive) {
      _amountIsValid.value = true;
    } else {
      _amountIsValid.value = false;
    }
  }

  void _onConfirmTransfer() {
    final amountResult = _amountController.text.parseToMoney();
    if (amountResult.isFailure) {
      AppSnackbar.show(
        context,
        message:
            'Valor inválido. Por favor, insira um valor maior que zero'
            ' para a transferência.',
        type: SnackbarType.error,
      );
      return;
    }

    final amount = amountResult.value!;
    final description = _descriptionController.text.trim();

    final transferData = TransferConfirmationData.fromRecipientInfo(
      recipientInfo: _recipientInfo,
      fromAccountId: _viewModel.selectedAccount!.id,
      amount: amount,
      description: description,
    );

    context.pushNamed(
      TransferRoutes.confirmation.name,
      extra: transferData,
    );
  }
}
