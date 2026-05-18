import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/datetime_extension.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/big_button.dart';
import '/ui/components/input_formatters/cpf_input_formatter.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text_form_field/basic_text_form_field.dart';
import '/ui/pages/auth/register/viewmodel/register_viewmodel.dart';

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
  final _birthDateController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _cpfController = TextEditingController();
  final _passwordController = TextEditingController();
  final _emailCodeController = TextEditingController();
  final _phoneCodeController = TextEditingController();

  bool _obscurePassword = true;

  @override
  void initState() {
    _viewmodel = widget.viewmodel;

    _nameController.text = _viewmodel.name;
    _emailController.text = _viewmodel.email;
    _cpfController.text = _viewmodel.cpf;
    _phoneController.text = _formatPhoneForDisplay(_viewmodel.phone);
    _passwordController.text = _viewmodel.password;
    _birthDateController.text = _viewmodel.birthDate?.formatDayLabel ?? '';

    _viewmodel.register.addListener(_onRegisterCommandChanged);
    _viewmodel.requestEmailCode.addListener(_onRequestEmailCodeChanged);
    _viewmodel.confirmEmailCode.addListener(_onConfirmEmailCodeChanged);
    _viewmodel.requestPhoneCode.addListener(_onRequestPhoneCodeChanged);
    _viewmodel.confirmPhoneCode.addListener(_onConfirmPhoneCodeChanged);

    super.initState();
  }

  @override
  void dispose() {
    _viewmodel.register.removeListener(_onRegisterCommandChanged);
    _viewmodel.requestEmailCode.removeListener(_onRequestEmailCodeChanged);
    _viewmodel.confirmEmailCode.removeListener(_onConfirmEmailCodeChanged);
    _viewmodel.requestPhoneCode.removeListener(_onRequestPhoneCodeChanged);
    _viewmodel.confirmPhoneCode.removeListener(_onConfirmPhoneCodeChanged);

    _nameController.dispose();
    _birthDateController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _cpfController.dispose();
    _passwordController.dispose();
    _emailCodeController.dispose();
    _phoneCodeController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final mergedListenable = Listenable.merge([
      _viewmodel,
      _viewmodel.register,
      _viewmodel.requestEmailCode,
      _viewmodel.confirmEmailCode,
      _viewmodel.requestPhoneCode,
      _viewmodel.confirmPhoneCode,
    ]);

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
                animation: mergedListenable,
                builder: (context, _) {
                  final isRunning = _isAnyCommandRunning;

                  return Form(
                    key: _formKey,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.stretch,
                      mainAxisSize: MainAxisSize.min,
                      spacing: 16,
                      children: [
                        Text(
                          'Cadastro em etapas para proteger sua conta.',
                          style: Theme.of(context).textTheme.bodyLarge
                              ?.copyWith(
                                color: colorScheme.onSurfaceVariant,
                              ),
                          textAlign: TextAlign.center,
                        ),
                        _buildStepHeader(context),
                        _buildCurrentStepContent(context, isRunning),
                        if (_viewmodel.stepError != null &&
                            _viewmodel.stepErrorStep == _viewmodel.currentStep)
                          Container(
                            padding: const EdgeInsets.all(12),
                            decoration: BoxDecoration(
                              color: colorScheme.errorContainer,
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              _viewmodel.stepError!.message,
                              style: Theme.of(context).textTheme.bodyMedium
                                  ?.copyWith(
                                    color: colorScheme.onErrorContainer,
                                  ),
                            ),
                          ),

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
        animation: mergedListenable,
        builder: (context, _) {
          final isRunning = _isAnyCommandRunning;

          return Padding(
            padding: const EdgeInsets.fromLTRB(16, 8, 16, 16),
            child: Row(
              spacing: 12,
              children: [
                if (_viewmodel.currentStep != RegisterStep.personalData)
                  Expanded(
                    child: BigButton(
                      onPressed: isRunning ? null : _previousStep,
                      label: 'Voltar',
                      leftIcon: const Icon(Icons.arrow_back_ios_new_rounded),
                    ),
                  ),
                Expanded(
                  flex: 2,
                  child: BigButton(
                    onPressed: isRunning ? null : _submit,
                    enabled: _isPrimaryEnabled,
                    label: _primaryLabel,
                    rightIcon: isRunning
                        ? const SizedBox(
                            width: 20,
                            height: 20,
                            child: CircularProgressIndicator(strokeWidth: 2),
                          )
                        : const Icon(Icons.arrow_forward_rounded, size: 24),
                  ),
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  bool get _isAnyCommandRunning {
    return _viewmodel.register.isRunning ||
        _viewmodel.requestEmailCode.isRunning ||
        _viewmodel.confirmEmailCode.isRunning ||
        _viewmodel.requestPhoneCode.isRunning ||
        _viewmodel.confirmPhoneCode.isRunning;
  }

  bool get _isPrimaryEnabled {
    return switch (_viewmodel.currentStep) {
      RegisterStep.personalData => true,
      RegisterStep.contactData => true,
      RegisterStep.emailVerification => _viewmodel.isEmailVerified,
      RegisterStep.phoneVerification => _viewmodel.isPhoneVerified,
      RegisterStep.review => _viewmodel.canRegister,
    };
  }

  String get _primaryLabel {
    if (_isAnyCommandRunning) return 'Aguarde...';

    return switch (_viewmodel.currentStep) {
      RegisterStep.personalData => 'Continuar',
      RegisterStep.contactData => 'Continuar',
      RegisterStep.emailVerification => 'Ir para telefone',
      RegisterStep.phoneVerification => 'Revisar cadastro',
      RegisterStep.review => 'Concluir cadastro',
    };
  }

  Widget _buildStepHeader(BuildContext context) {
    final steps = [
      (RegisterStep.personalData, 'Dados'),
      (RegisterStep.contactData, 'Contato'),
      (RegisterStep.emailVerification, 'E-mail'),
      (RegisterStep.phoneVerification, 'Telefone'),
      (RegisterStep.review, 'Revisao'),
    ];

    return Wrap(
      spacing: 8,
      runSpacing: 8,
      alignment: WrapAlignment.center,
      children: steps.map((entry) {
        final isCurrent = _viewmodel.currentStep == entry.$1;
        final isDone =
            _stepIndex(_viewmodel.currentStep) > _stepIndex(entry.$1);

        final color = isCurrent
            ? Theme.of(context).colorScheme.primary
            : isDone
            ? Theme.of(context).colorScheme.secondary
            : Theme.of(context).colorScheme.surfaceContainerHighest;
        final textColor = isCurrent || isDone
            ? Theme.of(context).colorScheme.onPrimary
            : Theme.of(context).colorScheme.onSurfaceVariant;

        return Container(
          padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
          decoration: BoxDecoration(
            color: color,
            borderRadius: BorderRadius.circular(999),
          ),
          child: Text(
            entry.$2,
            style: TextStyle(
              color: textColor,
              fontWeight: FontWeight.w700,
            ),
          ),
        );
      }).toList(),
    );
  }

  Widget _buildCurrentStepContent(BuildContext context, bool isRunning) {
    return switch (_viewmodel.currentStep) {
      RegisterStep.personalData => _buildPersonalDataStep(isRunning),
      RegisterStep.contactData => _buildContactDataStep(isRunning),
      RegisterStep.emailVerification => _buildEmailVerificationStep(isRunning),
      RegisterStep.phoneVerification => _buildPhoneVerificationStep(isRunning),
      RegisterStep.review => _buildReviewStep(context),
    };
  }

  Widget _buildPersonalDataStep(bool isRunning) {
    return Column(
      spacing: 16,
      children: [
        BasicTextFormField(
          controller: _nameController,
          textCapitalization: TextCapitalization.words,
          enabled: !isRunning,
          textInputAction: TextInputAction.next,
          labelText: 'Nome completo',
          hintText: 'Seu nome completo',
          prefixIcon: const Icon(Icons.person_outline),
          validator: _nameValidator,
        ),
        BasicTextFormField(
          controller: _cpfController,
          keyboardType: TextInputType.number,
          enabled: !isRunning,
          textInputAction: TextInputAction.next,
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
            CpfInputFormatter(),
          ],
          labelText: 'CPF',
          hintText: '000.000.000-00',
          prefixIcon: const Icon(Icons.badge_outlined),
          validator: _cpfValidator,
        ),
        BasicTextFormField(
          controller: _birthDateController,
          enabled: false,
          labelText: 'Data de nascimento',
          hintText: 'Selecione sua data de nascimento',
          prefixIcon: const Icon(Icons.calendar_month_outlined),
          validator: (_) => _birthDateValidator(),
        ),
        Align(
          alignment: Alignment.centerLeft,
          child: TextButton.icon(
            onPressed: isRunning ? null : _pickBirthDate,
            icon: const Icon(Icons.calendar_today_outlined),
            label: const Text('Selecionar data de nascimento'),
          ),
        ),
        BasicTextFormField(
          controller: _passwordController,
          obscureText: _obscurePassword,
          enabled: !isRunning,
          autofillHints: const [AutofillHints.newPassword],
          textInputAction: TextInputAction.done,
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
          validator: _passwordValidator,
          onFieldSubmitted: (_) => _submit(),
        ),
      ],
    );
  }

  Widget _buildContactDataStep(bool isRunning) {
    return Column(
      spacing: 16,
      children: [
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
        BasicTextFormField(
          controller: _phoneController,
          keyboardType: TextInputType.phone,
          enabled: !isRunning,
          textInputAction: TextInputAction.done,
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
            _BrazilPhoneInputFormatter(),
          ],
          labelText: 'Telefone',
          hintText: '(27) 99999-9999',
          prefixIcon: const Icon(Icons.phone_outlined),
          validator: _phoneValidator,
        ),
      ],
    );
  }

  Widget _buildEmailVerificationStep(bool isRunning) {
    return Column(
      spacing: 12,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text('Verifique o e-mail ${_viewmodel.email}.'),
        Row(
          spacing: 12,
          children: [
            Expanded(
              child: BigButton(
                onPressed: isRunning ? null : _requestEmailCode,
                label: 'Enviar codigo e-mail',
              ),
            ),
          ],
        ),
        BasicTextFormField(
          controller: _emailCodeController,
          keyboardType: TextInputType.number,
          enabled: !isRunning,
          textInputAction: TextInputAction.done,
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
            LengthLimitingTextInputFormatter(6),
          ],
          labelText: 'Codigo de e-mail',
          hintText: '000000',
          prefixIcon: const Icon(Icons.mark_email_read_outlined),
        ),
        BigButton(
          onPressed: isRunning ? null : _confirmEmailCode,
          label: _viewmodel.isEmailVerified
              ? 'E-mail confirmado'
              : 'Confirmar e-mail',
          enabled: !_viewmodel.isEmailVerified,
        ),
      ],
    );
  }

  Widget _buildPhoneVerificationStep(bool isRunning) {
    return Column(
      spacing: 12,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        Text('Verifique o telefone ${_phoneController.text}.'),
        BigButton(
          onPressed: isRunning ? null : _requestPhoneCode,
          label: 'Enviar codigo telefone',
        ),
        BasicTextFormField(
          controller: _phoneCodeController,
          keyboardType: TextInputType.number,
          enabled: !isRunning,
          textInputAction: TextInputAction.done,
          inputFormatters: [
            FilteringTextInputFormatter.digitsOnly,
            LengthLimitingTextInputFormatter(6),
          ],
          labelText: 'Codigo de telefone',
          hintText: '000000',
          prefixIcon: const Icon(Icons.sms_outlined),
        ),
        BigButton(
          onPressed: isRunning ? null : _confirmPhoneCode,
          label: _viewmodel.isPhoneVerified
              ? 'Telefone confirmado'
              : 'Confirmar telefone',
          enabled: !_viewmodel.isPhoneVerified,
        ),
      ],
    );
  }

  Widget _buildReviewStep(BuildContext context) {
    final lines = [
      'Nome: ${_viewmodel.name}',
      'CPF: ${_cpfController.text}',
      'Nascimento: ${_birthDateController.text}',
      'E-mail: ${_viewmodel.email}',
      'Telefone: ${_phoneController.text}',
      'E-mail verificado: ${_viewmodel.isEmailVerified ? 'sim' : 'nao'}',
      'Telefone verificado: ${_viewmodel.isPhoneVerified ? 'sim' : 'nao'}',
    ];

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Theme.of(context).colorScheme.surfaceContainerHighest,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        spacing: 8,
        children: [
          Text(
            'Revise seus dados antes de concluir.',
            style: Theme.of(context).textTheme.titleMedium,
          ),
          for (final line in lines)
            Text(
              line,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
        ],
      ),
    );
  }

  void _navToLogin() {
    context.goNamed(AuthRoutes.login.name);
  }

  void _previousStep() {
    _viewmodel.previousStep();
  }

  String? _nameValidator(String? value) {
    final name = (value ?? '').trim();
    if (name.isEmpty) return 'Informe o nome completo.';
    if (name.length < 3) {
      return 'Informe um nome valido.';
    }
    return null;
  }

  String? _birthDateValidator() {
    if (_viewmodel.birthDate == null) {
      return 'Selecione a data de nascimento.';
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

  String? _phoneValidator(String? value) {
    final phoneDigits = (value ?? '').replaceAll(RegExp(r'\D'), '');
    if (phoneDigits.isEmpty) return 'Informe o telefone.';
    if (phoneDigits.length < 10 || phoneDigits.length > 11) {
      return 'Informe um telefone valido.';
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
      AppSnackbar.show(
        context,
        type: SnackbarType.success,
        title: 'Sucesso',
        message: 'Cadastro realizado com sucesso.',
      );

      context.goNamed(AuthRoutes.login.name);
    }
  }

  void _onRequestEmailCodeChanged() {
    final command = _viewmodel.requestEmailCode;
    if (!mounted || command.isRunning) return;

    if (command.isFailure) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message:
            command.error?.message ?? 'Falha ao solicitar codigo de e-mail.',
      );
      return;
    }

    if (command.isSuccess) {
      AppSnackbar.show(
        context,
        type: SnackbarType.info,
        title: 'Codigo enviado',
        message: 'Codigo de verificacao enviado para o e-mail.',
      );
    }
  }

  void _onConfirmEmailCodeChanged() {
    final command = _viewmodel.confirmEmailCode;
    if (!mounted || command.isRunning) return;

    if (command.isFailure) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message:
            command.error?.message ?? 'Falha ao confirmar codigo de e-mail.',
      );
      return;
    }

    if (command.isSuccess) {
      AppSnackbar.show(
        context,
        type: SnackbarType.success,
        title: 'E-mail confirmado',
        message: 'Seu e-mail foi confirmado com sucesso.',
      );
    }
  }

  void _onRequestPhoneCodeChanged() {
    final command = _viewmodel.requestPhoneCode;
    if (!mounted || command.isRunning) return;

    if (command.isFailure) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message:
            command.error?.message ?? 'Falha ao solicitar codigo de telefone.',
      );
      return;
    }

    if (command.isSuccess) {
      AppSnackbar.show(
        context,
        type: SnackbarType.info,
        title: 'Codigo enviado',
        message: 'Codigo de verificacao enviado para o telefone.',
      );
    }
  }

  void _onConfirmPhoneCodeChanged() {
    final command = _viewmodel.confirmPhoneCode;
    if (!mounted || command.isRunning) return;

    if (command.isFailure) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Erro',
        message:
            command.error?.message ?? 'Falha ao confirmar codigo de telefone.',
      );
      return;
    }

    if (command.isSuccess) {
      AppSnackbar.show(
        context,
        type: SnackbarType.success,
        title: 'Telefone confirmado',
        message: 'Seu telefone foi confirmado com sucesso.',
      );
    }
  }

  Future<void> _submit() async {
    if (_isAnyCommandRunning) return;

    switch (_viewmodel.currentStep) {
      case RegisterStep.personalData:
        final form = _formKey.currentState;
        if (form == null || !form.validate()) return;

        _viewmodel.updatePersonalData(
          name: _nameController.text,
          cpf: _cpfController.text,
          birthDate: _viewmodel.birthDate,
          password: _passwordController.text,
        );
        _viewmodel.nextStep();
        return;

      case RegisterStep.contactData:
        final form = _formKey.currentState;
        if (form == null || !form.validate()) return;

        _viewmodel.updateContactData(
          email: _emailController.text,
          phone: _toApiPhone(_phoneController.text),
        );
        _viewmodel.nextStep();
        return;

      case RegisterStep.emailVerification:
      case RegisterStep.phoneVerification:
        _viewmodel.nextStep();
        return;

      case RegisterStep.review:
        await _viewmodel.register.execute();
        return;
    }
  }

  Future<void> _pickBirthDate() async {
    final now = DateTime.now();
    final initial = _viewmodel.birthDate ?? DateTime(now.year - 18, 1, 1);

    final selected = await showDatePicker(
      context: context,
      initialDate: initial,
      firstDate: DateTime(1900, 1, 1),
      lastDate: now,
      helpText: 'Selecione a data de nascimento',
    );

    if (selected == null) return;

    _viewmodel.updatePersonalData(birthDate: selected);
    _birthDateController.text = selected.formatDayLabel;
  }

  Future<void> _requestEmailCode() async {
    if (_viewmodel.email.trim().isEmpty) {
      _viewmodel.updateContactData(email: _emailController.text);
    }
    await _viewmodel.requestEmailCode.execute();
  }

  Future<void> _confirmEmailCode() async {
    await _viewmodel.confirmEmailCode.execute(_emailCodeController.text);
  }

  Future<void> _requestPhoneCode() async {
    if (_viewmodel.phone.trim().isEmpty) {
      _viewmodel.updateContactData(phone: _toApiPhone(_phoneController.text));
    }
    await _viewmodel.requestPhoneCode.execute();
  }

  Future<void> _confirmPhoneCode() async {
    await _viewmodel.confirmPhoneCode.execute(_phoneCodeController.text);
  }

  int _stepIndex(RegisterStep step) {
    return switch (step) {
      RegisterStep.personalData => 0,
      RegisterStep.contactData => 1,
      RegisterStep.emailVerification => 2,
      RegisterStep.phoneVerification => 3,
      RegisterStep.review => 4,
    };
  }

  String _toApiPhone(String value) {
    final digits = value.replaceAll(RegExp(r'\D'), '');
    if (digits.isEmpty) return '';
    if (digits.startsWith('55')) return '+$digits';
    return '+55$digits';
  }

  String _formatPhoneForDisplay(String apiPhone) {
    final digits = apiPhone.replaceAll(RegExp(r'\D'), '');
    final localDigits = digits.startsWith('55') ? digits.substring(2) : digits;
    return _BrazilPhoneInputFormatter.format(localDigits);
  }
}

class _BrazilPhoneInputFormatter extends TextInputFormatter {
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    final digits = newValue.text.replaceAll(RegExp(r'\D'), '');
    final clipped = digits.length > 11 ? digits.substring(0, 11) : digits;
    final formatted = format(clipped);

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(offset: formatted.length),
    );
  }

  static String format(String digits) {
    if (digits.isEmpty) return '';

    final d = digits.length > 11 ? digits.substring(0, 11) : digits;
    if (d.length <= 2) return '($d';
    if (d.length <= 6) return '(${d.substring(0, 2)}) ${d.substring(2)}';
    if (d.length <= 10) {
      return '(${d.substring(0, 2)}) ${d.substring(2, 6)}-${d.substring(6)}';
    }
    return '(${d.substring(0, 2)}) ${d.substring(2, 7)}-${d.substring(7)}';
  }
}
