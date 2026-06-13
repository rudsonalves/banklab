import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/core/result/result.dart';
import '/core/routing/models/transaction_password_setup_origin.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/big_button.dart';
import '/ui/components/input_text/basic_input_text.dart';
import '/ui/components/messages/app_snackbar.dart';
import '../models/post_login_destination.dart';
import '../viewmodel/login_viewmodel.dart';

const _accountApprovalRequiredMessage =
    'Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você poderá acessar sua conta.';
const _contactNotVerifiedGenericMessage =
    'Confirme seu e-mail e telefone antes de entrar.';
const _contactNotVerifiedEmailOnlyMessage =
    'Confirme seu e-mail antes de entrar.';
const _contactNotVerifiedPhoneOnlyMessage =
    'Confirme seu telefone antes de entrar.';
const _postLoginBlockedMessage =
    'Não foi possível liberar seu acesso agora. Tente novamente.';
const _postLoginSessionErrorMessage =
    'Não foi possível carregar sua sessão. Entre novamente.';

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

                        BasicInputText(
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
                          builder: (context, value, child) => BasicInputText(
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
                : Icon(Icons.login_rounded, size: 24),
          );
        },
      ),
    );
  }

  void _navToRegister() {
    context.goNamed(RegisterRoutes.cpf.routeName);
  }

  void _obscurePasswordListener() {
    _obscurePassword.value = !_obscurePassword.value;
  }

  String? _emailValidator(String? value) {
    final email = (value ?? '').trim();
    if (email.isEmpty) return 'Informe o e-mail.';

    if (!email.isValidEmail) return 'Informe um e-mail válido.';

    return null;
  }

  String? _passwordValidator(String? value) {
    final password = value ?? '';
    if (password.isEmpty) return 'Informe a senha.';
    if (password.length < 6) {
      return 'A senha deve ter no mínimo 6 caracteres.';
    }
    return null;
  }

  void _onLoginCommandChanged() {
    final loginCommand = _viewModel.login;
    if (!mounted || loginCommand.isRunning) return;

    if (loginCommand.isFailure) {
      final error = loginCommand.error;
      final message = _resolveLoginErrorMessage(error);

      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message: message,
      );

      return;
    }

    if (loginCommand.isSuccess) _handlePostLoginDestination();
  }

  void _handlePostLoginDestination() {
    final destination = _viewModel.resolvePostLoginDestination();

    switch (destination) {
      case PostLoginDestination.home:
        context.goNamed(BaseRoutes.home.routeName);
        return;

      case PostLoginDestination.transactionPassword:
        context.goNamed(
          TransactionPasswordRoutes.introduction.routeName,
          extra: TransactionPasswordSetupOrigin.postLogin,
        );
        return;

      case PostLoginDestination.blocked:
        AppSnackbar.show(
          context,
          type: SnackbarType.error,
          title: 'Acesso indisponível',
          message: _postLoginBlockedMessage,
        );
        return;

      case PostLoginDestination.sessionError:
        AppSnackbar.show(
          context,
          type: SnackbarType.error,
          title: 'Sessão indisponível',
          message: _postLoginSessionErrorMessage,
        );
        return;
    }
  }

  String _resolveLoginErrorMessage(AppError? error) {
    if (error == null) return 'Falha ao autenticar.';

    if (error.code == AppErrorCode.accountApprovalRequired) {
      return _accountApprovalRequiredMessage;
    }

    if (error.code == AppErrorCode.contactNotVerified) {
      final details = error.details;
      if (details is Map<String, dynamic>) {
        final emailVerified = details['email_verified'];
        final phoneVerified = details['phone_verified'];

        if (emailVerified == false && phoneVerified == true) {
          return _contactNotVerifiedEmailOnlyMessage;
        }

        if (emailVerified == true && phoneVerified == false) {
          return _contactNotVerifiedPhoneOnlyMessage;
        }

        return _contactNotVerifiedGenericMessage;
      }

      return error.message.isNotEmpty
          ? error.message
          : _contactNotVerifiedGenericMessage;
    }

    return error.message.isNotEmpty ? error.message : 'Falha ao autenticar.';
  }

  Future<void> _submit() async {
    final form = _formKey.currentState;
    if (form == null || !form.validate()) return;

    FocusScope.of(context).unfocus();

    _viewModel.login.execute(
      LoginRequestDto(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      ),
    );
  }
}
