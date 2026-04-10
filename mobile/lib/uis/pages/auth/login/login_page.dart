import 'package:bankflow/core/routing/routes.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/data/services/apis/auth/dtos/login_request_dto.dart';
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
  final _formKey = GlobalKey<FormState>();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();

  bool _obscurePassword = true;

  @override
  void initState() {
    super.initState();
    widget.viewModel.login.addListener(_onLoginCommandChanged);
  }

  @override
  void dispose() {
    widget.viewModel.login.removeListener(_onLoginCommandChanged);
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 460),
              child: AnimatedBuilder(
                animation: widget.viewModel.login,
                builder: (context, _) {
                  final isRunning = widget.viewModel.login.isRunning;

                  return Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      mainAxisSize: MainAxisSize.min,
                      children: [
                        Text(
                          'Entrar',
                          style: Theme.of(context).textTheme.headlineMedium,
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 8),
                        Text(
                          'Acesse sua conta para continuar no BankFlow.',
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                              ),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 28),
                        TextFormField(
                          controller: _emailController,
                          keyboardType: TextInputType.emailAddress,
                          autofillHints: const [AutofillHints.email],
                          enabled: !isRunning,
                          decoration: const InputDecoration(
                            labelText: 'E-mail',
                            hintText: 'voce@exemplo.com',
                            prefixIcon: Icon(Icons.email_outlined),
                          ),
                          validator: _emailValidator,
                        ),
                        const SizedBox(height: 16),
                        TextFormField(
                          controller: _passwordController,
                          obscureText: _obscurePassword,
                          enabled: !isRunning,
                          autofillHints: const [AutofillHints.password],
                          decoration: InputDecoration(
                            labelText: 'Senha',
                            prefixIcon: const Icon(Icons.lock_outline),
                            suffixIcon: IconButton(
                              onPressed: isRunning
                                  ? null
                                  : () {
                                      setState(() {
                                        _obscurePassword = !_obscurePassword;
                                      });
                                    },
                              icon: Icon(
                                _obscurePassword
                                    ? Icons.visibility_outlined
                                    : Icons.visibility_off_outlined,
                              ),
                            ),
                          ),
                          validator: _passwordValidator,
                          onFieldSubmitted: (_) => _submit(),
                        ),
                        const SizedBox(height: 24),
                        FilledButton(
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
                        ),
                        const SizedBox(height: 12),
                        TextButton(
                          onPressed: isRunning
                              ? null
                              : () => context.goNamed(AuthRoutes.register.name),
                          child: const Text('Nao tem conta? Cadastre-se'),
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
    );
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
    final loginCommand = widget.viewModel.login;
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

    await widget.viewModel.login.execute(
      LoginRequestDto(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      ),
    );

    final result = widget.viewModel.login.result!;
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
