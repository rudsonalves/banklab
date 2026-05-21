import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/core/result/errors/app_error.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import '../../components/input_text/basic_input_text.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterEmailPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterEmailPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterEmailPage> createState() => _RegisterEmailPageState();
}

class _RegisterEmailPageState extends State<RegisterEmailPage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  final _emailController = TextEditingController();
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);

  @override
  void initState() {
    super.initState();

    _initialize();
  }

  @override
  void dispose() {
    _emailController.dispose();
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
            spacing: 12,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextHeader('Informe o email'),
              BasicInputText(
                controller: _emailController,
                hintText: 'Digite seu email',
                keyboardType: TextInputType.emailAddress,
                autofillHints: const [AutofillHints.email],
                textInputAction: TextInputAction.next,
                onChanged: _emailChanged,
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([
          _isDisabled,
          _viewmodel.submitAndRequestEmailToken,
        ]),
        builder: (context, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navBack,
            rightOnPressed: _submitEmail,
            isRightEnabled:
                !_isDisabled.value &&
                !_viewmodel.submitAndRequestEmailToken.isRunning,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _navBack() => context.pop();

  void _navToEmailToken() => context.pushNamed(RegisterRoutes.emailToken.name);

  void _emailChanged(String value) {
    _isDisabled.value = !value.isValidEmail;
  }

  Future<void> _submitEmail() async {
    final email = _emailController.text.trim();

    if (!email.isValidEmail) {
      _isDisabled.value = true;
      return;
    }

    await _viewmodel.submitAndRequestEmailToken.execute(email);

    final result = _viewmodel.submitAndRequestEmailToken.result!;
    if (result.isFailure) {
      _showErrorMessage(result.error!);
      return;
    }

    _navToEmailToken();
  }

  void _showErrorMessage(AppError error) {
    if (!mounted) return;

    AppSnackbar.show(
      context,
      message: error.message,
      type: SnackbarType.error,
    );
  }

  void _initialize() {
    final initialEmail = _viewmodel.state?.email;
    if (initialEmail != null) {
      _emailController.text = initialEmail;
      _isDisabled.value = !initialEmail.isValidEmail;
    }
  }
}
