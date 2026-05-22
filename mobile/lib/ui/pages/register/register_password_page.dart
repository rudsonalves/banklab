import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/errors/app_error.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/basic_input_text.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import 'model/password_draft.dart';
import 'viewmodel/register_viewmodel.dart';
import 'widgets/criterial_item_row.dart';

class RegisterPasswordPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterPasswordPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterPasswordPage> createState() => _RegisterPasswordPageState();
}

class _RegisterPasswordPageState extends State<RegisterPasswordPage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  final _passwordFocusNode = FocusNode();
  final _confirmPasswordFocusNode = FocusNode();
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  final ValueNotifier<bool> _isPasswordObscured = ValueNotifier(true);
  final ValueNotifier<bool> _isConfirmPasswordObscured = ValueNotifier(true);
  final ValueNotifier<bool> _hasNumber = ValueNotifier(false);
  final ValueNotifier<bool> _hasUppercase = ValueNotifier(false);
  final ValueNotifier<bool> _hasLowercase = ValueNotifier(false);
  final ValueNotifier<bool> _hasMinLength = ValueNotifier(false);
  final ValueNotifier<bool> _passwordsMatch = ValueNotifier(false);

  final PasswordDraft _password = PasswordDraft();

  @override
  void dispose() {
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    _passwordFocusNode.dispose();
    _confirmPasswordFocusNode.dispose();
    _isDisabled.dispose();
    _isPasswordObscured.dispose();
    _isConfirmPasswordObscured.dispose();
    _hasNumber.dispose();
    _hasUppercase.dispose();
    _hasLowercase.dispose();
    _hasMinLength.dispose();
    _passwordsMatch.dispose();

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
              TextHeader('Crie sua senha'),
              ValueListenableBuilder<bool>(
                valueListenable: _isPasswordObscured,
                builder: (context, isObscured, _) => BasicInputText(
                  controller: _passwordController,
                  focusNode: _passwordFocusNode,
                  hintText: 'Digite sua senha',
                  obscureText: isObscured,
                  onChanged: _passwordChanged,
                  textInputAction: TextInputAction.next,
                  onFieldSubmitted: (_) =>
                      _confirmPasswordFocusNode.requestFocus(),
                  suffixIcon: Focus(
                    canRequestFocus: false,
                    skipTraversal: true,
                    descendantsAreFocusable: false,
                    child: IconButton(
                      onPressed: _togglePasswordVisibility,
                      icon: Icon(
                        isObscured ? Icons.visibility : Icons.visibility_off,
                      ),
                    ),
                  ),
                ),
              ),
              ValueListenableBuilder<bool>(
                valueListenable: _hasNumber,
                builder: (context, checked, _) => CriterialItemRow(
                  checked: checked,
                  label: 'tem números',
                ),
              ),
              ValueListenableBuilder<bool>(
                valueListenable: _hasUppercase,
                builder: (context, checked, _) => CriterialItemRow(
                  checked: checked,
                  label: 'tem caracteres maiúsculos',
                ),
              ),
              ValueListenableBuilder<bool>(
                valueListenable: _hasLowercase,
                builder: (context, checked, _) => CriterialItemRow(
                  checked: checked,
                  label: 'tem caracteres minúsculos',
                ),
              ),
              ValueListenableBuilder<bool>(
                valueListenable: _hasMinLength,
                builder: (context, checked, _) => CriterialItemRow(
                  checked: checked,
                  label:
                      'tem comprimento maior ou igual a ${PasswordDraft.minLength}'
                      ' caracteres',
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'Confirme sua senha',
                style: Theme.of(context).textTheme.titleMedium,
              ),
              ValueListenableBuilder<bool>(
                valueListenable: _isConfirmPasswordObscured,
                builder: (context, isObscured, _) => BasicInputText(
                  controller: _confirmPasswordController,
                  focusNode: _confirmPasswordFocusNode,
                  hintText: 'Confirme sua senha',
                  obscureText: isObscured,
                  onChanged: _passwordChanged,
                  textInputAction: TextInputAction.done,
                  suffixIcon: Focus(
                    canRequestFocus: false,
                    skipTraversal: true,
                    descendantsAreFocusable: false,
                    child: IconButton(
                      onPressed: _toggleConfirmPasswordVisibility,
                      icon: Icon(
                        isObscured ? Icons.visibility : Icons.visibility_off,
                      ),
                    ),
                  ),
                ),
              ),
              ValueListenableBuilder<bool>(
                valueListenable: _passwordsMatch,
                builder: (context, checked, _) => CriterialItemRow(
                  checked: checked,
                  label: 'senhas são iguais',
                ),
              ),
            ],
          ),
        ),
      ),

      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([
          _isDisabled,
          _viewmodel.submitPassword,
          _viewmodel.register,
        ]),
        builder: (context, _) => DoubleBottomButton(
          leftButtonLabel: 'Voltar',
          rightButtonLabel: 'Concluir',
          leftOnPressed: _navBack,
          rightOnPressed: _submitPassword,
          isRightEnabled:
              !_isDisabled.value &&
              !_viewmodel.submitPassword.isRunning &&
              !_viewmodel.register.isRunning,
          rightButtonIcon: const Icon(Icons.check),
        ),
      ),
    );
  }

  void _navBack() => context.pop();

  void _navToSuccessStatus() => context.pushNamed(RegisterRoutes.success.name);

  void _navToFailureStatus() => context.pushNamed(RegisterRoutes.failure.name);

  void _togglePasswordVisibility() {
    _isPasswordObscured.value = !_isPasswordObscured.value;
  }

  void _toggleConfirmPasswordVisibility() {
    _isConfirmPasswordObscured.value = !_isConfirmPasswordObscured.value;
  }

  void _passwordChanged(String _) {
    _password.value = _passwordController.text.trim();
    _password.confirm = _confirmPasswordController.text.trim();

    _hasNumber.value = _password.hasNumber;
    _hasUppercase.value = _password.hasUppercase;
    _hasLowercase.value = _password.hasLowercase;
    _hasMinLength.value = _password.hasMinLength;
    _passwordsMatch.value = _password.hasEquals;

    _isDisabled.value = !_password.isValid;
  }

  Future<void> _submitPassword() async {
    if (_isDisabled.value) return;

    _password.value = _passwordController.text.trim();
    _password.confirm = _confirmPasswordController.text.trim();

    await _viewmodel.submitPassword.execute(_password.model);

    final passwordResult = _viewmodel.submitPassword.result;
    if (passwordResult == null) return;
    if (passwordResult.isFailure) {
      _showErrorMessage(passwordResult.error!);
      return;
    }

    await _viewmodel.register.execute();

    final registerResult = _viewmodel.register.result;
    if (registerResult == null) return;
    if (registerResult.isFailure) {
      _navToFailureStatus();
      return;
    }

    _navToSuccessStatus();
  }

  void _showErrorMessage(AppError error) {
    if (!mounted) return;

    AppSnackbar.show(
      context,
      message: error.message,
      type: SnackbarType.error,
    );
  }
}
