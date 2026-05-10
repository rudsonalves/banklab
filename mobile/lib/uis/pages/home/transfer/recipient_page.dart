import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '/uis/core/input_formatters/cpf_input_formatter.dart';
import '/uis/core/text/text_header.dart';
import '/uis/core/text_form_field/basic_text_form_field.dart';
import '../../../../core/routing/routes.dart';
import 'viewmodel/transfer_viewmodel.dart';
import 'widgets/dropdown_recipient.dart';
import 'widgets/recipient_card.dart';

class RecipientPage extends StatefulWidget {
  final TransferViewmodel viewModel;

  const RecipientPage({super.key, required this.viewModel});

  @override
  State<RecipientPage> createState() => _RecipientPageState();
}

class _RecipientPageState extends State<RecipientPage> {
  TransferViewmodel get _viewModel => widget.viewModel;

  final _documentController = TextEditingController();
  final _branchController = TextEditingController(text: "0001");
  final _accountController = TextEditingController();

  @override
  void dispose() {
    _documentController.dispose();
    _branchController.dispose();
    _accountController.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Conta do Destinatário'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          mainAxisAlignment: MainAxisAlignment.center,
          spacing: 12,
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextHeader('CPF'),
                BasicTextFormField(
                  controller: _documentController,
                  hintText: '000.000.000-00',
                  inputFormatters: [CpfInputFormatter()],
                  keyboardType: TextInputType.number,
                  onChanged: _onCpfChanged,
                ),
              ],
            ),

            Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextHeader('Agência e número da conta'),
                Row(
                  children: [
                    Expanded(
                      child: BasicTextFormField(
                        controller: _branchController,
                        hintText: '0001',
                        keyboardType: TextInputType.number,
                        onChanged: _onBranchChenged,
                      ),
                    ),
                    const SizedBox(width: 16),
                    Expanded(
                      child: BasicTextFormField(
                        controller: _accountController,
                        hintText: '0000000-0',
                        keyboardType: TextInputType.number,
                        onChanged: _onAccountChanged,
                      ),
                    ),
                  ],
                ),
              ],
            ),

            const SizedBox(height: 12),
            Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                TextHeader('Selecione a conta destino'),
                ListenableBuilder(
                  listenable: _viewModel.getInternalRecipient,
                  builder: (context, _) {
                    final inicialValue =
                        _viewModel.receipientAccounts.isNotEmpty
                        ? _viewModel.receipientAccounts.first.accountId
                        : null;

                    return Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      children: [
                        DropdownRecipient(
                          receipientAccounts: _viewModel.receipientAccounts,
                          inicialValue: inicialValue,
                          onChanged: _onRecipientChanged,
                        ),

                        if (_viewModel.selectedRecipient.value != null) ...[
                          SizedBox(height: 16),
                          TextHeader('Conta selecionada'),
                          ValueListenableBuilder(
                            valueListenable: _viewModel.selectedRecipient,
                            builder: (context, value, _) => RecipientCard(
                              selectedRecipient: value!,
                            ),
                          ),
                        ],
                      ],
                    );
                  },
                ),
              ],
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
              onPressed: _onConfigTransfer,
              enabled: isButtonEnabled,
            );
          },
        ),
      ),
    );
  }

  void _onConfigTransfer() {
    context.pushNamed(
      TransferRoutes.payment.name,
      extra: _viewModel.selectedRecipient.value,
    );
  }

  void _onRecipientChanged(String? value) {
    if (value != null) {
      FocusScope.of(context).unfocus();
    }
  }

  void _onCpfChanged(String value) {
    if (value.isValidCpf) {
      FocusScope.of(context).unfocus();
      _getAccountByCpf(value.onlyNumbers);
    }
  }

  void _onBranchChenged(String value) {
    if (value.onlyNumbers.length == 4) {
      FocusScope.of(context).nextFocus();
      if (_isValidBranchAndAccount()) {
        _getAccountByBranchAndAccount();
      }
    }
  }

  void _onAccountChanged(String value) {
    if (value.onlyNumbers.length == 8) {
      FocusScope.of(context).unfocus();
      if (_isValidBranchAndAccount()) {
        _getAccountByBranchAndAccount();
      }
    }
  }

  bool _isValidBranchAndAccount() {
    bool isValidBranch = _branchController.text.onlyNumbers.length == 4;
    bool isValidAccount = _accountController.text.onlyNumbers.length == 8;

    return isValidBranch && isValidAccount;
  }

  Future<void> _getAccountByCpf(String cpf) async {
    final recipientRequest = RecipientRequestDto(document: cpf);
    await _viewModel.getInternalRecipient.execute(recipientRequest);

    if (_viewModel.getInternalRecipient.isFailure) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Erro ao buscar conta do destinatário')),
      );
    }
  }

  Future<void> _getAccountByBranchAndAccount() async {
    final recipientRequest = RecipientRequestDto(
      branch: _branchController.text.onlyNumbers,
      accountNumber: _accountController.text.onlyNumbers,
    );
    await _viewModel.getInternalRecipient.execute(recipientRequest);

    if (_viewModel.getInternalRecipient.isFailure) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Erro ao buscar conta do destinatário')),
      );
    }
  }
}
