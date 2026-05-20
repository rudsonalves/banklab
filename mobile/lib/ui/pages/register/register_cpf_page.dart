import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/core/result/errors/app_error.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_formatters/cpf_input_formatter.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import '/ui/components/text_form_field/basic_text_form_field.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterCpfPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterCpfPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterCpfPage> createState() => _RegisterCpfPageState();
}

class _RegisterCpfPageState extends State<RegisterCpfPage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  final _cpfController = TextEditingController();
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);

  @override
  void initState() {
    super.initState();

    _viewmodel.startEmptyRegisterState();
  }

  @override
  void dispose() {
    _cpfController.dispose();
    _isDisabled.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Registro de Conta'),
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        child: SingleChildScrollView(
          padding: const EdgeInsets.symmetric(vertical: 24, horizontal: 12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextHeader('Informe o CPF'),
              BasicTextFormField(
                controller: _cpfController,
                hintText: 'Digite apenas os números do seu CPF',
                keyboardType: TextInputType.number,
                onChanged: _cpfChanged,
                inputFormatters: [CpfInputFormatter()],
              ),
            ],
          ),
        ),
      ),

      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([_isDisabled, _viewmodel.submitCPF]),
        builder: (context, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Login',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navLogin,
            rightOnPressed: _checkCpfIsAvailable,
            isRightEnabled:
                !_isDisabled.value && !_viewmodel.submitCPF.isRunning,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _navLogin() => context.goNamed(AuthRoutes.login.name);

  void _navToName() => context.pushNamed(RegisterRoutes.name.name);

  void _cpfChanged(String value) {
    _isDisabled.value = !value.isValidCpf;

    if (!_isDisabled.value) {
      FocusScope.of(context).unfocus();
    }
  }

  Future<void> _checkCpfIsAvailable() async {
    final cpf = _cpfController.text.onlyNumbers;

    await _viewmodel.submitCPF.execute(cpf);

    final result = _viewmodel.submitCPF.result!;
    if (result.isFailure) {
      _showErrorMessages(result.error!);
      return;
    }

    _navToName();
  }

  void _showErrorMessages(AppError error) {
    if (!mounted) return;
    if (error.code == AppErrorCode.cpfAlreadyRegistered) {
      AppSnackbar.show(
        context,
        message: 'Este CPF já se encontra cadastrado.',
        type: SnackbarType.error,
      );
      return;
    }

    AppSnackbar.show(
      context,
      message: 'Erro inesperado',
    );
  }
}
