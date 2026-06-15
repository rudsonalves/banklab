import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/models/transaction_password_setup_origin.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/token_input.dart';
import '/ui/components/text/text_header.dart';
import 'confirm_transaction_password_page.dart';
import 'viewmodel/transaction_password_viewmodel.dart';

class CreateTransactionPasswordPage extends StatefulWidget {
  final TransactionPasswordSetupOrigin origin;
  final TransactionPasswordViewModel viewModel;

  const CreateTransactionPasswordPage({
    super.key,
    required this.origin,
    required this.viewModel,
  });

  @override
  State<CreateTransactionPasswordPage> createState() =>
      _CreateTransactionPasswordPageState();
}

class _CreateTransactionPasswordPageState
    extends State<CreateTransactionPasswordPage> {
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  TransactionPasswordSetupOrigin get _origin => widget.origin;
  TransactionPasswordViewModel get _viewModel => widget.viewModel;

  String _token = '';

  @override
  void dispose() {
    _token = '';
    _isDisabled.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Senha transacional'),
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 24, horizontal: 12),
          child: Column(
            spacing: 20,
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const TextHeader('Escolha uma sequência de 6 dígitos'),
              Center(
                child: TokenInput(
                  visible: true,
                  onChanged: _onPinChanged,
                  onCompleted: _onPinCompleted,
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ValueListenableBuilder<bool>(
        valueListenable: _isDisabled,
        builder: (context, isDisabled, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navBack,
            rightOnPressed: isDisabled ? null : _navToConfirmation,
            isRightEnabled: !isDisabled,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _onPinChanged(String value) {
    _token = value.trim();
    _isDisabled.value = _token.length != 6;
  }

  void _onPinCompleted(String value) {
    _onPinChanged(value);
    FocusScope.of(context).unfocus();
  }

  void _navBack() => context.pop();

  Future<void> _navToConfirmation() async {
    if (_token.length != 6) return;

    await Navigator.of(context).push(
      MaterialPageRoute<void>(
        builder: (context) => ConfirmTransactionPasswordPage(
          token: _token,
          origin: _origin,
          viewModel: _viewModel,
        ),
      ),
    );

    if (!mounted) return;
  }
}
