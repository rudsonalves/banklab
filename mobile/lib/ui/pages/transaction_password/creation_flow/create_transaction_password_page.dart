import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/token_input.dart';
import '/ui/components/text/text_header.dart';

class CreateTransactionPasswordPage extends StatefulWidget {
  const CreateTransactionPasswordPage({
    super.key,
  });

  @override
  State<CreateTransactionPasswordPage> createState() =>
      _CreateTransactionPasswordPageState();
}

class _CreateTransactionPasswordPageState
    extends State<CreateTransactionPasswordPage> {
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  String _pin = '';

  @override
  void dispose() {
    _pin = '';
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
            leftOnPressed: _navIntroduction,
            rightOnPressed: isDisabled ? null : _navToConfirmation,
            isRightEnabled: !isDisabled,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _onPinChanged(String value) {
    _pin = value.trim();
    _isDisabled.value = _pin.length != 6;
  }

  void _onPinCompleted(String value) {
    _onPinChanged(value);
    FocusScope.of(context).unfocus();
  }

  void _navIntroduction() =>
      context.goNamed(TransactionPasswordRoutes.introduction.name);

  void _navToConfirmation() {
    if (_pin.length != 6) return;

    context.pushNamed(
      TransactionPasswordRoutes.confirm.name,
      extra: _pin,
    );
  }
}
