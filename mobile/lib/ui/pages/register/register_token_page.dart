import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/command.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/contact_verification/enums/contact_verification_channel.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import '../../components/input_text/otp_input.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterTokenPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;
  final ContactVerificationChannel channel;

  const RegisterTokenPage({
    super.key,
    required this.viewmodel,
    required this.channel,
  });

  @override
  State<RegisterTokenPage> createState() => _RegisterTokenPageState();
}

class _RegisterTokenPageState extends State<RegisterTokenPage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;
  ContactVerificationChannel get _tokenType => widget.channel;

  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);
  String _token = '';

  @override
  void dispose() {
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
          padding: const EdgeInsets.symmetric(
            vertical: 24,
            horizontal: 12,
          ),
          child: Column(
            spacing: 12,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              TextHeader(_headerText),
              OtpInput(
                onChanged: _tokenChanged,
                // onCompleted: _tokenCompleted,
              ),
            ],
          ),
        ),
      ),

      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([
          _isDisabled,
          _confirmTokenCommand,
        ]),
        builder: (context, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navBack,
            rightOnPressed: _submitToken,
            isRightEnabled:
                !_isDisabled.value && !_confirmTokenCommand.isRunning,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  Command1<Unit, String> get _confirmTokenCommand {
    switch (_tokenType) {
      case ContactVerificationChannel.email:
        return _viewmodel.confirmEmailToken;
      case ContactVerificationChannel.phone:
        return _viewmodel.confirmPhoneToken;
    }
  }

  String get _headerText {
    switch (_tokenType) {
      case ContactVerificationChannel.email:
        return 'Informe o código enviado para o seu e-mail';
      case ContactVerificationChannel.phone:
        return 'Informe o código enviado para o seu telefone';
    }
  }

  void _navBack() => context.pop();

  void _navToNext() {
    switch (_tokenType) {
      case ContactVerificationChannel.email:
        context.pushNamed(RegisterRoutes.phone.name);
        break;
      case ContactVerificationChannel.phone:
        context.pushNamed(RegisterRoutes.password.name);
        break;
    }
  }

  void _tokenChanged(String value) {
    _token = value;
    _isDisabled.value = value.trim().length < 6;

    if (!_isDisabled.value) {
      FocusScope.of(context).unfocus();
    }
  }

  // void _tokenCompleted(String value) {
  //   _tokenChanged(value);
  // }

  Future<void> _submitToken() async {
    final token = _token.trim();
    if (token.length < 6) {
      _isDisabled.value = true;
      return;
    }

    await _confirmTokenCommand.execute(token);

    final result = _confirmTokenCommand.result;
    if (result == null) return;

    if (result.isFailure) {
      _showErrorMessage(result.error!);
      return;
    }

    _navToNext();
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
