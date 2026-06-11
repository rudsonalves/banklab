import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/result.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/token_input.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import 'viewmodel/transaction_password_viewmodel.dart';

const _pinMismatchMessage = 'A confirmação deve ser igual à senha criada.';
const _alreadySetMessage = 'Sua senha transacional já está cadastrada.';

class ConfirmTransactionPasswordPage extends StatefulWidget {
  final TransactionPasswordViewModel viewModel;
  final String pin;

  const ConfirmTransactionPasswordPage({
    super.key,
    required this.viewModel,
    required this.pin,
  });

  @override
  State<ConfirmTransactionPasswordPage> createState() =>
      _ConfirmTransactionPasswordPageState();
}

class _ConfirmTransactionPasswordPageState
    extends State<ConfirmTransactionPasswordPage> {
  TransactionPasswordViewModel get _viewModel => widget.viewModel;

  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  String _confirmation = '';

  @override
  void dispose() {
    _confirmation = '';
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
              const TextHeader('Confirme sua senha transacional'),
              Center(
                child: TokenInput(
                  visible: false,
                  onChanged: _onConfirmationChanged,
                  onCompleted: _onConfirmationCompleted,
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([
          _isDisabled,
          _viewModel.create,
        ]),
        builder: (context, _) {
          final isRunning = _viewModel.create.isRunning;

          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: isRunning ? 'Criando...' : 'Concluir',
            leftOnPressed: isRunning ? null : _navBack,
            rightOnPressed: _isDisabled.value || isRunning ? null : _submit,
            isRightEnabled: !_isDisabled.value && !isRunning,
            rightButtonIcon: isRunning
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Icon(Icons.check),
          );
        },
      ),
    );
  }

  void _onConfirmationChanged(String value) {
    _confirmation = value.trim();
    _isDisabled.value = _confirmation.length != 6;
  }

  void _onConfirmationCompleted(String value) {
    _onConfirmationChanged(value);
    FocusScope.of(context).unfocus();
  }

  void _navBack() => context.pop();

  Future<void> _submit() async {
    if (_confirmation.length != 6) return;

    if (_confirmation != widget.pin) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        message: _pinMismatchMessage,
      );
      return;
    }

    final transPasswdRequest = CreateTransactionPasswordRequestDto(
      password: widget.pin,
      confirmation: _confirmation,
    );
    await _viewModel.create.execute(transPasswdRequest);

    final result = _viewModel.create.result;
    if (result == null || !mounted) return;

    if (result.isFailure) {
      _errorSession(result.error!);
      return;
    }

    if (!_viewModel.canAccessHome) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        message: 'Não foi possível ativar a senha transacional.',
      );
      return;
    }

    _clearSensitiveState();
    context.goNamed(BaseRoutes.home.routeName);
    return;
  }

  void _errorSession(AppError error) {
    if (error.code == AppErrorCode.transactionPasswordAlreadySet) {
      _clearSensitiveState();
      AppSnackbar.show(
        context,
        type: SnackbarType.info,
        message: _alreadySetMessage,
      );
      return;
    }

    AppSnackbar.show(
      context,
      type: SnackbarType.error,
      message: error.message,
    );
  }

  void _clearSensitiveState() {
    _confirmation = '';
  }
}
