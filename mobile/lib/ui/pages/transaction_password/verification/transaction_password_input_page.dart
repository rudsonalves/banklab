import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/token_input.dart';
import '/ui/components/text/text_header.dart';

class TransactionPasswordInputPage extends StatefulWidget {
  const TransactionPasswordInputPage({super.key});

  @override
  State<TransactionPasswordInputPage> createState() =>
      _TransactionPasswordInputPageState();
}

class _TransactionPasswordInputPageState
    extends State<TransactionPasswordInputPage> {
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  String _transPasswd = '';

  @override
  void dispose() {
    _transPasswd = '';
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
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(vertical: 24, horizontal: 12),
          child: Column(
            spacing: 20,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              const TextHeader(
                'Confirme a operação com sua senha transacional',
              ),
              Center(
                child: TokenInput(
                  visible: false,
                  onChanged: _onTokenChanged,
                  onCompleted: _onTokenCompleted,
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([_isDisabled]),
        builder: (context, _) => DoubleBottomButton(
          leftButtonLabel: 'Cancelar',
          rightButtonLabel: 'Concluir',
          leftOnPressed: _navBack,
          rightOnPressed: _isDisabled.value ? null : _submit,
          isRightEnabled: !_isDisabled.value,
          rightButtonIcon: const Icon(Icons.check),
        ),
      ),
    );
  }

  void _onTokenChanged(String value) {
    _transPasswd = value.trim();
    _isDisabled.value = _transPasswd.length != 6;
  }

  void _onTokenCompleted(String value) {
    _onTokenChanged(value);
    FocusScope.of(context).unfocus();
  }

  void _navBack() => context.pop(null);

  Future<void> _submit() async {
    if (_transPasswd.length != 6) return;

    if (!mounted) return;
    context.pop(_transPasswd);
    return;
  }
}
