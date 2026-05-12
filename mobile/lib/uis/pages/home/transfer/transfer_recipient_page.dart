import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/apis/transfer/dtos/recipient_request_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '/uis/core/cards/recipient_card.dart';
import '/uis/core/input_formatters/cpf_input_formatter.dart';
import '/uis/core/text/text_header.dart';
import '/uis/core/text_form_field/basic_text_form_field.dart';
import 'viewmodel/transfer_viewmodel.dart';
import 'widgets/dropdown_recipient.dart';

class TransferRecipientPage extends StatefulWidget {
  final TransferViewmodel viewModel;

  const TransferRecipientPage({super.key, required this.viewModel});

  @override
  State<TransferRecipientPage> createState() => _TransferRecipientPageState();
}

class _TransferRecipientPageState extends State<TransferRecipientPage> {
  TransferViewmodel get _viewModel => widget.viewModel;

  final _documentController = TextEditingController();
  final _branchController = TextEditingController(text: "0001");
  final _accountController = TextEditingController();
  final _selectedRecipient = ValueNotifier<RecipientInfoDto?>(null);

  @override
  void dispose() {
    _documentController.dispose();
    _branchController.dispose();
    _accountController.dispose();
    _selectedRecipient.dispose();
    _viewModel.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Conta do Destinatário'),
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        child: SingleChildScrollView(
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
                          onChanged: _onBranchChanged,
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

                          if (_selectedRecipient.value != null) ...[
                            SizedBox(height: 16),
                            TextHeader('Conta selecionada'),
                            ValueListenableBuilder(
                              valueListenable: _selectedRecipient,
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
      ),

      bottomNavigationBar: Padding(
        padding: const EdgeInsets.symmetric(vertical: 16),
        child: ValueListenableBuilder(
          valueListenable: _selectedRecipient,
          builder: (context, value, _) {
            final isButtonEnabled = value != null;

            return BigButton(
              label: 'Prosseguir',
              onPressed: _onConfigTransfer,
              enabled: isButtonEnabled,
              rightIcon: Icon(Icons.arrow_forward_ios_rounded, size: 24),
            );
          },
        ),
      ),
    );
  }

  void _onConfigTransfer() {
    context.pushNamed(
      TransferRoutes.payment.name,
      extra: _selectedRecipient.value,
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

  void _onBranchChanged(String value) {
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

      _selectedRecipient.value = null;
      return;
    }

    _selectedRecipient.value = _viewModel.receipientAccounts.first;

    _branchController.text = '0001';
    _accountController.text = _viewModel.receipientAccounts.first.accountNumber;
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

      _selectedRecipient.value = null;
      return;
    }

    _selectedRecipient.value = _viewModel.receipientAccounts.first;

    _documentController.text = _viewModel.receipientAccounts.first.document;
  }
}
