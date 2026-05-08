import 'package:flutter/material.dart';

import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/text_form_field/basic_text_form_field.dart';
import 'viewmodel/transfer_viewmodel.dart';
import 'widgets/account_dropdown.dart';
import 'widgets/section_title.dart';

class TransferPage extends StatefulWidget {
  final TransferViewmodel viewModel;

  const TransferPage({super.key, required this.viewModel});

  @override
  State<TransferPage> createState() => _TransferPageState();
}

class _TransferPageState extends State<TransferPage> {
  final _beneficiaryNameController = TextEditingController();
  final TextEditingController _branchController = TextEditingController(
    text: "0001",
  );
  final TextEditingController _accountController = TextEditingController();
  final TextEditingController _amountController = TextEditingController();

  String? _selectedOriginAccount;

  @override
  void dispose() {
    _beneficiaryNameController.dispose();
    _branchController.dispose();
    _accountController.dispose();
    _amountController.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Transferência'),
      ),
      body: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            spacing: 12,
            children: [
              // Section: Source Account
              const SectionTitle('Conta de Origem'),
              // TODO: Put this in the command reactive Notifier Widget
              AccountDropdown(
                accounts: widget.viewModel.accounts!,
                selectedAccountId: _selectedOriginAccount!,
                onChanged: (value) {
                  setState(() {
                    _selectedOriginAccount = value;
                  });
                },
              ),
              const SizedBox(height: 12),

              // Section: Beneficiary Data
              const SectionTitle('Dados do Beneficiário'),
              BasicTextFormField(
                controller: _beneficiaryNameController,
                labelText: 'Nome do Beneficiário',
                hintText: 'Digite o nome completo',
              ),

              Row(
                children: [
                  Expanded(
                    flex: 1,
                    child: BasicTextFormField(
                      controller: _branchController,
                      labelText: 'Agência',
                      hintText: '0001',
                    ),
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    flex: 2,
                    child: BasicTextFormField(
                      controller: _accountController,
                      labelText: 'Conta',
                      hintText: 'Digite o número da conta',
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),
              // Section: Amount
              const SectionTitle('Valor'),
              BasicTextFormField(
                controller: _amountController,
                labelText: 'Valor da Transferência',
                hintText: 'R\$ 0,00',
                keyboardType: const TextInputType.numberWithOptions(
                  decimal: true,
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: Padding(
        padding: const EdgeInsets.all(16),
        child: ElevatedButton(
          onPressed: null,
          child: const Padding(
            padding: EdgeInsets.symmetric(vertical: 12),
            child: Text('Confirmar Transferência'),
          ),
        ),
      ),
    );
  }
}
