import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/data/services/apis/auth/dtos/register_request_dto.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/pages/auth/register/viewmodel/register_viewmodel.dart';

class RegisterPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  late final RegisterViewmodel _viewmodel;

  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _cpfController = TextEditingController();
  final _passwordController = TextEditingController();

  bool _obscurePassword = true;

  @override
  void initState() {
    _viewmodel = widget.viewmodel;
    _viewmodel.register.addListener(_onRegisterCommandChanged);

    super.initState();
  }

  @override
  void dispose() {
    _viewmodel.register.removeListener(_onRegisterCommandChanged);
    _nameController.dispose();
    _emailController.dispose();
    _cpfController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Criar conta'),
      ),
      body: GestureDetector(
        onTap: () => FocusScope.of(context).unfocus(),
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24),
            child: ConstrainedBox(
              constraints: const BoxConstraints(maxWidth: 460),
              child: AnimatedBuilder(
                animation: _viewmodel.register,
                builder: (context, _) {
                  final isRunning = _viewmodel.register.isRunning;

                  return Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      mainAxisSize: MainAxisSize.min,
                      spacing: 16,
                      children: [
                        Text(
                          'Cadastre-se para começar a usar o BankFlow.',
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                              ),
                          textAlign: TextAlign.center,
                        ),
                        const SizedBox(height: 12),
                        TextFormField(
                          controller: _nameController,
                          textCapitalization: TextCapitalization.words,
                          enabled: !isRunning,
                          decoration: const InputDecoration(
                            labelText: 'Nome completo',
                            hintText: 'Seu nome completo',
                            prefixIcon: Icon(Icons.person_outline),
                          ),
                          validator: _nameValidator,
                        ),

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

                        TextFormField(
                          controller: _cpfController,
                          keyboardType: TextInputType.number,
                          enabled: !isRunning,
                          decoration: const InputDecoration(
                            labelText: 'CPF',
                            hintText: '00000000000',
                            prefixIcon: Icon(Icons.badge_outlined),
                          ),
                          validator: _cpfValidator,
                        ),

                        TextFormField(
                          controller: _passwordController,
                          obscureText: _obscurePassword,
                          enabled: !isRunning,
                          autofillHints: const [AutofillHints.newPassword],
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

                        const SizedBox(height: 6),
                        Align(
                          alignment: Alignment.centerRight,
                          child: TextButton(
                            onPressed: isRunning ? null : _navToLogin,
                            child: const Text('Já tem conta? Faça login'),
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
        animation: _viewmodel.register,
        builder: (context, _) {
          final isRunning = _viewmodel.register.isRunning;

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
                : const Text('Cadastrar'),
          );
        },
      ),
    );
  }

  void _navToLogin() {
    context.goNamed(AuthRoutes.login.name);
  }

  String? _nameValidator(String? value) {
    final name = (value ?? '').trim();
    if (name.isEmpty) return 'Informe o nome completo.';
    if (name.length < 3) {
      return 'Informe um nome valido.';
    }
    return null;
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

  String? _cpfValidator(String? value) {
    final cpf = (value ?? '').replaceAll(RegExp(r'\D'), '');
    if (cpf.isEmpty) return 'Informe o CPF.';
    if (cpf.length != 11) {
      return 'O CPF deve ter 11 digitos.';
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

  void _onRegisterCommandChanged() {
    final registerCommand = _viewmodel.register;
    if (!mounted || registerCommand.isRunning) return;

    if (registerCommand.isFailure) {
      final message = registerCommand.error?.message ?? 'Falha ao cadastrar.';
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

    if (registerCommand.isSuccess) {
      ScaffoldMessenger.of(context)
        ..hideCurrentSnackBar()
        ..showSnackBar(
          const SnackBar(
            content: Text('Cadastro realizado com sucesso.'),
            behavior: SnackBarBehavior.floating,
          ),
        );
    }
  }

  Future<void> _submit() async {
    final form = _formKey.currentState;
    if (form == null || !form.validate()) return;

    FocusScope.of(context).unfocus();

    await _viewmodel.register.execute(
      RegisterRequestDto(
        name: _nameController.text.trim(),
        email: _emailController.text.trim(),
        password: _passwordController.text,
        cpf: _cpfController.text.replaceAll(RegExp(r'\D'), ''),
      ),
    );

    final result = _viewmodel.register.result!;
    if (result.isFailure) {
      final message = result.error?.message ?? 'Falha ao cadastrar.';
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
    context.goNamed(AuthRoutes.login.name);
  }
}
