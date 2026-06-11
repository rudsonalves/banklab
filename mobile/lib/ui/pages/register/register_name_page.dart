import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/result/errors/app_error.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import '../../components/input_text/basic_input_text.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterNamePage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterNamePage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterNamePage> createState() => _RegisterNamePageState();
}

class _RegisterNamePageState extends State<RegisterNamePage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  final _nameController = TextEditingController();
  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);

  @override
  void initState() {
    super.initState();

    _initialize();
  }

  @override
  void dispose() {
    _nameController.dispose();
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
              TextHeader('Informe o nome completo'),
              BasicInputText(
                controller: _nameController,
                hintText: 'Digite seu nome completo',
                textCapitalization: TextCapitalization.words,
                onChanged: _nameChanged,
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([_isDisabled, _viewmodel.submitName]),
        builder: (context, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navBack,
            rightOnPressed: _submitName,
            isRightEnabled:
                !_isDisabled.value && !_viewmodel.submitName.isRunning,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _navBack() => context.pop();

  void _navToBirthdate() =>
      context.pushNamed(RegisterRoutes.birthDate.routeName);

  void _nameChanged(String value) {
    final name = value.trim();
    _isDisabled.value = name.isEmpty || name.split(' ').length < 2;
  }

  Future<void> _submitName() async {
    final name = _nameController.text.trim();

    await _viewmodel.submitName.execute(name);

    final result = _viewmodel.submitName.result!;
    if (result.isFailure) {
      _showErrorMessage(result.error!);
      return;
    }

    _navToBirthdate();
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
    final initialName = _viewmodel.state?.name;
    if (initialName != null) {
      _nameController.text = initialName;
      _isDisabled.value = initialName.trim().isEmpty;
    }
  }
}
