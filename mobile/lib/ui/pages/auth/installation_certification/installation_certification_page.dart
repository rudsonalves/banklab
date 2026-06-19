import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/result.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/input_text/token_input.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import 'viewmodel/installation_certification_viewmodel.dart';

class InstallationCertificationPage extends StatefulWidget {
  const InstallationCertificationPage({
    super.key,
    required this.viewModel,
  });

  final InstallationCertificationViewModel viewModel;

  @override
  State<InstallationCertificationPage> createState() =>
      _InstallationCertificationPageState();
}

class _InstallationCertificationPageState
    extends State<InstallationCertificationPage> {
  late final InstallationCertificationViewModel _viewModel;
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  String _transactionPassword = '';

  @override
  void initState() {
    super.initState();
    _viewModel = widget.viewModel;
    _viewModel.certify.addListener(_onCertifyChanged);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted || _viewModel.hasRestrictedInstallationAuth) return;
      context.goNamed(AuthRoutes.login.routeName);
    });
  }

  @override
  void dispose() {
    _transactionPassword = '';
    _isDisabled.dispose();
    _viewModel.certify.removeListener(_onCertifyChanged);
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Certificar instalação'),
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
                'Confirme esta instalação com sua senha transacional',
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
        listenable: Listenable.merge([_isDisabled, _viewModel.certify]),
        builder: (context, _) => DoubleBottomButton(
          leftButtonLabel: 'Cancelar',
          rightButtonLabel: 'Certificar',
          leftOnPressed: _viewModel.certify.isRunning ? null : _cancel,
          rightOnPressed: _isDisabled.value ? null : _submit,
          isRightEnabled: !_isDisabled.value,
          rightButtonIcon: _viewModel.certify.isRunning
              ? const SizedBox(
                  width: 16,
                  height: 16,
                  child: CircularProgressIndicator(strokeWidth: 2),
                )
              : const Icon(Icons.check),
        ),
      ),
    );
  }

  /// Listener for changes in the token input. Trims the input and updates the
  /// disabled state of the submit button.
  void _onTokenChanged(String value) {
    _transactionPassword = value.trim();
    _isDisabled.value = _transactionPassword.length != 6;
  }

  /// Listener for when the token input is completed. Trims the input and unfocuses the field.
  void _onTokenCompleted(String value) {
    _onTokenChanged(value);
    FocusScope.of(context).unfocus();
  }

  /// Submits the certification request using the current transaction password.
  /// Shows a loading state while processing.
  Future<void> _submit() async {
    if (_transactionPassword.length != 6 || _viewModel.certify.isRunning) {
      return;
    }

    await _viewModel.certify.execute(_transactionPassword);
    _transactionPassword = '';
  }

  /// Cancels the certification process and navigates back to the login page.
  Future<void> _cancel() async {
    await _viewModel.cancel();
    if (!mounted) return;
    context.goNamed(AuthRoutes.login.routeName);
  }

  /// Listener for changes in the certification process. Navigates or shows messages
  /// based on the outcome.
  void _onCertifyChanged() {
    final command = _viewModel.certify;
    if (!mounted || command.isRunning) return;

    if (command.isSuccess) {
      context.goNamed(BaseRoutes.home.routeName);
      return;
    }

    if (command.isFailure) {
      AppSnackbar.show(
        context,
        type: SnackbarType.error,
        title: 'Certificação não concluída',
        message: _resolveCertificationErrorMessage(command.error),
      );
      context.goNamed(AuthRoutes.login.routeName);
    }
  }

  /// Resolves a user-friendly error message based on the provided [AppError].
  String _resolveCertificationErrorMessage(AppError? error) {
    if (error?.code == AppErrorCode.transactionPasswordNotSet) {
      return 'Sua senha transacional ainda não está configurada. Entre por'
          ' uma instalação já autorizada para regularizar antes de tentar novamente.';
    }

    if (error?.code == AppErrorCode.transactionPasswordLocked) {
      return 'Sua senha transacional está bloqueada. Entre por'
          ' uma instalação já autorizada para regularizar antes de tentar novamente.';
    }

    final backendCode = backendErrorCode(error);
    if (backendCode == 'INVALID_TOKEN' || backendCode == 'UNAUTHORIZED') {
      return 'Sua autorização expirou. Entre novamente para continuar.';
    }

    return error?.message ?? 'Não foi possível certificar esta instalação.';
  }
}
