import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '../../../core/text_form_field/basic_text_form_field.dart';
import 'viewmodel/login_viewmodel.dart';

class LoginPage extends StatefulWidget {
  final LoginViewModel viewModel;

  const LoginPage({
    super.key,
    required this.viewModel,
  });

  @override
  State<LoginPage> createState() => _LoginPageState();
}

class _LoginPageState extends State<LoginPage> {
  late final LoginViewModel _viewModel;

  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
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

    _emailController.dispose();
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

                  return Form(
                    key: _formKey,
                    child: Column(
                      spacing: 16,
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          'Acesse sua conta para continuar no BankFlow.',
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                              ),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 12),

                        BasicTextFormField(
                          controller: _emailController,
                          keyboardType: TextInputType.emailAddress,
                          autofillHints: const [AutofillHints.email],
                          enabled: !isRunning,
                          textInputAction: TextInputAction.next,
                          labelText: 'E-mail',
                          hintText: 'voce@exemplo.com',
                          prefixIcon: const Icon(Icons.email_outlined),
                          validator: _emailValidator,
                        ),

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
                            onPressed: isRunning ? null : _navToRegister,
                            child: const Text('Não tem conta? Cadastre-se'),
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

          return FilledButton(
            onPressed: isRunning ? null : _submit,
            child: isRunning
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                    ),
                  )
                : const Text('Entrar'),
          );
        },
      ),
    );
  }

  void _navToRegister() {
    context.goNamed(AuthRoutes.register.name);
  }

  void _obscurePasswordListener() {
    _obscurePassword.value = !_obscurePassword.value;
  }

  String? _emailValidator(String? value) {
    final email = (value ?? '').trim();
    if (email.isEmpty) return 'Informe o e-mail.';

    final emailRegex = RegExp(
      r'^[^@\s]+@[^@\s]+\.[^@\s]+$',
    );
    if (!emailRegex.hasMatch(email)) {
      return 'Informe um e-mail valido.';
    }

    return null;
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
      final message = loginCommand.error?.message ?? 'Falha ao autenticar.';
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(
          SnackBar(
            content: Text(message),
            behavior: SnackBarBehavior.floating,
          ),
        );
      return;
    }

    if (loginCommand.isSuccess) {
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(
          const SnackBar(
            content: Text('Login realizado com sucesso.'),
            behavior: SnackBarBehavior.floating,
          ),
        );
    }
  }

  Future<void> _submit() async {
    final form = _formKey.currentState;
    if (form == null || !form.validate()) return;

    FocusScope.of(context).unfocus();

    await _viewModel.login.execute(
      LoginRequestDto(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      ),
    );

    final result = _viewModel.login.result!;
    if (result.isFailure) {
      final message = result.error?.message ?? 'Falha ao autenticar.';
      if (!mounted) return;
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(
          SnackBar(
            content: Text(message),
            behavior: SnackBarBehavior.floating,
          ),
        );
      return;
    }

    if (!mounted) return;
    context.goNamed(HomeRoutes.home.name);
  }
}
