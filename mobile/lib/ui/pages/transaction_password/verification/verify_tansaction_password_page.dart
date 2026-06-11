import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/errors/app_error.dart';
import '/data/services/apis/transaction_password/dtos/set_up_authorize_request_dto.dart';
import '/data/services/apis/transaction_password/enums/step_up_operation.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/token_input.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import '../../../../core/routing/routes.dart';
import 'viewmodel/verify_tansaction_password_viewmodel.dart';

class VerifyTansactionPasswordPage extends StatefulWidget {
  final VerifyTansactionPasswordViewmodel viewModel;

  const VerifyTansactionPasswordPage({
    super.key,
    required this.viewModel,
  });

  @override
  State<VerifyTansactionPasswordPage> createState() =>
      _VerifyTansactionPasswordPageState();
}

class _VerifyTansactionPasswordPageState
    extends State<VerifyTansactionPasswordPage> {
  VerifyTansactionPasswordViewmodel get _viewModel => widget.viewModel;

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
        listenable: Listenable.merge([
          _isDisabled,
          _viewModel.stepUpAuthorize,
        ]),
        builder: (context, _) {
          final isRunning = _viewModel.stepUpAuthorize.isRunning;

          return DoubleBottomButton(
            leftButtonLabel: 'Cancelar',
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

    final transPasswdRequest = SetUpAuthorizeRequestDto(
      operation: StepUpOperation.internalTransfer,
      transactionPassword: _transPasswd,
    );
    await _viewModel.stepUpAuthorize.execute(transPasswdRequest);

    final result = _viewModel.stepUpAuthorize.result;
    if (result == null || !mounted) return;

    if (result.isFailure) {
      _errorSession(result.error!);
      return;
    }

    _clearSensitiveState();
    final response = result.value!;
    if (!mounted) return;
    context.pop(response);
    return;
  }

  void _errorSession(AppError error) {
    if (error.code == AppErrorCode.transactionPasswordLocked) {
      _clearSensitiveState();
      context.pop(null);
      return;
    }

    if (error.code == AppErrorCode.transactionPasswordNotSet) {
      _clearSensitiveState();
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        message:
            'Você precisa configurar uma senha transacional para realizar esta operação.',
      );
      context.goNamed(BaseRoutes.home.routeName);
      return;
    }

    AppSnackbar.show(
      context,
      type: SnackbarType.error,
      message: 'Falha ao autorizar a operação. Tente novamente.',
    );
  }

  void _clearSensitiveState() {
    _transPasswd = '';
  }
}
