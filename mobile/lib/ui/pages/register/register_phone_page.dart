import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/extensions/string.dart';
import '/core/result/errors/app_error.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import '../../components/input_formatters/phone_input_formatter.dart';
import '../../components/input_text/basic_input_text.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterPhonePage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterPhonePage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterPhonePage> createState() => _RegisterPhonePageState();
}

class _RegisterPhonePageState extends State<RegisterPhonePage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  final _phoneController = TextEditingController();
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);

  @override
  void initState() {
    super.initState();

    _initialize();
  }

  @override
  void dispose() {
    _phoneController.dispose();
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
              TextHeader('Informe o telefone'),
              BasicInputText(
                controller: _phoneController,
                hintText: 'Digite seu telefone',
                keyboardType: TextInputType.number,
                textInputAction: TextInputAction.next,
                inputFormatters: [PhoneInputFormatter()],
                onChanged: _phoneChanged,
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([
          _isDisabled,
          _viewmodel.submitAndRequestPhoneToken,
        ]),
        builder: (context, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navBack,
            rightOnPressed: _submitPhone,
            isRightEnabled:
                !_isDisabled.value &&
                !_viewmodel.submitAndRequestPhoneToken.isRunning,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _navBack() => context.pop();

  void _navToPhoneToken() => context.pushNamed(RegisterRoutes.phoneToken.name);

  void _phoneChanged(String value) {
    _isDisabled.value = !value.isValidPhone;

    if (value.isValidPhone) {
      FocusScope.of(context).nextFocus();
    }
  }

  Future<void> _submitPhone() async {
    final phone = _phoneController.text.trim();

    if (!phone.isValidPhone) {
      _isDisabled.value = true;
      return;
    }

    await _viewmodel.submitAndRequestPhoneToken.execute(phone);

    final result = _viewmodel.submitAndRequestPhoneToken.result!;
    if (result.isFailure) {
      _showErrorMessage(result.error!);
      return;
    }

    _navToPhoneToken();
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
    final initialPhone = _viewmodel.state?.phone?.trim().onlyNumbers;
    if (initialPhone != null) {
      _phoneController.text = initialPhone;
      _isDisabled.value = !initialPhone.isValidPhone;
    }
  }
}
