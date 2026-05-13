import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/result.dart';
import '/core/routing/routes.dart';
import '/data/services/auth/api/dtos/login_request_dto.dart';
import '/data/services/auth/cache/models/last_login_identity.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';
import '/uis/core/messages/app_snackbar.dart';
import '/uis/core/text_form_field/basic_text_form_field.dart';
import 'viewmodel/short_login_viewmodel.dart';

const _accountApprovalRequiredMessage =
    'Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você poderá acessar sua conta.';

class ShortLoginPage extends StatefulWidget {
  final ShortLoginViewModel viewModel;
  final LastLoginIdentity identity;

  const ShortLoginPage({
    super.key,
    required this.viewModel,
    required this.identity,
  });

  @override
  State<ShortLoginPage> createState() => _ShortLoginPageState();
}

class _ShortLoginPageState extends State<ShortLoginPage> {
  late final ShortLoginViewModel _viewModel;

  final _formKey = GlobalKey<FormState>();
  final _passwordController = TextEditingController();

  final ValueNotifier<bool> _obscurePassword = ValueNotifier<bool>(true);

  @override
  void initState() {
    _viewModel = widget.viewModel;
    _viewModel.login.addListener(_onLoginCommandChanged);

    super.initState();
  }

  @override
  void dispose() {
    _viewModel.login.removeListener(_onLoginCommandChanged);

    _passwordController.dispose();
    _obscurePassword.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Entrar'),
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 460),
              child: AnimatedBuilder(
                animation: _viewModel.login,
                builder: (context, _) {
                  final isRunning = _viewModel.login.isRunning;
                  final userName = widget.identity.name.trim();

                  return Form(
                    key: _formKey,
                    child: Column(
                      spacing: 16,
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          'Bem-vindo de volta',
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                              ),
                          textAlign: TextAlign.center,
                        ),
                        Text(
                          userName,
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                                fontWeight: FontWeight.bold,
                              ),
                          textAlign: TextAlign.center,
                        ),

                        const SizedBox(height: 12),

                        ValueListenableBuilder<bool>(
                          valueListenable: _obscurePassword,
                          builder: (context, value, child) =>
                              BasicTextFormField(
                                controller: _passwordController,
                                obscureText: value,
                                enabled: !isRunning,
                                autofillHints: const [AutofillHints.password],
                                textInputAction: TextInputAction.done,
                                labelText: 'Senha',
                                hintText: '********',
                                prefixIcon: const Icon(Icons.lock_outline),
                                suffixIcon: IconButton(
                                  onPressed: isRunning
                                      ? null
                                      : _obscurePasswordListener,
                                  icon: Icon(
                                    value
                                        ? Icons.visibility_outlined
                                        : Icons.visibility_off_outlined,
                                  ),
                                ),
                                validator: _passwordValidator,
                                onFieldSubmitted: (_) => _submit(),
                              ),
                        ),

                        const SizedBox(height: 6),
                        Align(
                          alignment: Alignment.centerRight,
                          child: TextButton(
                            onPressed: isRunning ? null : _navToLogin,
                            child: const Text('Entrar com outra conta'),
                          ),
                        ),
                      ],
                    ),
                  );
                },
              ),
            ),
          ),
        ),
      ),

      bottomNavigationBar: AnimatedBuilder(
        animation: _viewModel.login,
        builder: (context, _) {
          final isRunning = _viewModel.login.isRunning;

          return BigButton(
            onPressed: isRunning ? null : _submit,
            label: isRunning ? 'Entrando...' : 'Entrar',
            rightIcon: isRunning
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                    ),
                  )
                : const Icon(Icons.login_rounded, size: 24),
          );
        },
      ),
    );
  }

  void _navToLogin() {
    context.goNamed(AuthRoutes.login.name);
  }

  void _obscurePasswordListener() {
    _obscurePassword.value = !_obscurePassword.value;
  }

  String? _passwordValidator(String? value) {
    final password = value ?? '';
    if (password.isEmpty) return 'Informe a senha.';
    if (password.length < 6) {
      return 'A senha deve ter no minimo 6 caracteres.';
    }
    return null;
  }

  void _onLoginCommandChanged() {
    final loginCommand = _viewModel.login;
    if (!mounted || loginCommand.isRunning) return;

    if (loginCommand.isFailure) {
      final error = loginCommand.error;
      final message = error?.code == AppErrorCode.accountApprovalRequired
          ? _accountApprovalRequiredMessage
          : error?.message ?? 'Falha ao autenticar.';

      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message: message,
      );

      return;
    }

    if (loginCommand.isSuccess) {
      AppSnackbar.show(
        context,
        type: SnackbarType.success,
        title: 'Sucesso',
        message: 'Login realizado com sucesso.',
      );
    }

    if (!mounted) return;
    context.goNamed(BaseRoutes.home.name);
  }

  Future<void> _submit() async {
    final form = _formKey.currentState;
    if (form == null || !form.validate()) return;

    FocusScope.of(context).unfocus();

    _viewModel.login.execute(
      LoginRequestDto(
        email: widget.identity.identifier.trim(),
        password: _passwordController.text,
      ),
    );
  }
}
