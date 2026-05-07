import 'package:flutter/material.dart';

import '../../../core/base/safe_scaffold.dart';
import '../../../core/text_form_field/basic_text_form_field.dart';
import 'viewmodel/transfer_viewmodel.dart';

class TransferPage extends StatefulWidget {
  final TransferViewmodel viewModel;

  const TransferPage({super.key, required this.viewModel});

  @override
  State<TransferPage> createState() => _TransferPageState();
}

class _TransferPageState extends State<TransferPage> {
  late final TextEditingController _beneficiaryNameController;
  late final TextEditingController _branchController;
  late final TextEditingController _accountController;
  late final TextEditingController _amountController;

  // String? _selectedTransferType;
  String? _selectedOriginAccount;

  @override
  void initState() {
    super.initState();
    _beneficiaryNameController = TextEditingController();
    _branchController = TextEditingController(text: "0001");
    _accountController = TextEditingController();
    _amountController = TextEditingController();
  }

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
              // Seção: Conta de Origem
              _buildSectionTitle(context, 'Conta de Origem'),
              _buildOriginAccountDropdown(),
              const SizedBox(height: 12),

              // Seção: Dados do Beneficiário
              _buildSectionTitle(context, 'Dados do Beneficiário'),
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
              // Seção: Valor
              _buildSectionTitle(context, 'Valor'),
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

  Widget _buildSectionTitle(BuildContext context, String title) {
    return Text(
      title,
      style: Theme.of(context).textTheme.titleMedium?.copyWith(
        fontWeight: FontWeight.w600,
      ),
    );
  }

  Widget _buildOriginAccountDropdown() {
    return DropdownButtonFormField<String>(
      initialValue: _selectedOriginAccount,
      hint: const Text('Selecione uma conta'),
      items: [
        DropdownMenuItem(
          value: 'account_1',
          child: const Text('Conta Corrente - 0001-2'),
        ),
        DropdownMenuItem(
          value: 'account_2',
          child: const Text('Conta Poupança - 0001-3'),
        ),
      ],
      onChanged: (value) {
        setState(() {
          _selectedOriginAccount = value;
        });
      },
      decoration: InputDecoration(
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
        ),
        contentPadding: const EdgeInsets.symmetric(
          horizontal: 12,
          vertical: 12,
        ),
      ),
    );
  }
}
