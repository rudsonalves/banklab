import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:intl/intl.dart';

import '/core/extensions/datetime_extension.dart';
import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/ui/components/messages/app_snackbar.dart';
import '/ui/components/text/text_header.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterBirthdatePage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterBirthdatePage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterBirthdatePage> createState() => _RegisterBirthdatePageState();
}

class _RegisterBirthdatePageState extends State<RegisterBirthdatePage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  final _selectedDate = ValueNotifier<DateTime?>(null);

  static final DateTime _maxDate = DateTime.now();
  static final DateTime _minDate = DateTime(1900);

  final ValueNotifier<bool> _isDisabled = ValueNotifier(true);

  @override
  void initState() {
    super.initState();

    _initialize();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

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
              TextHeader('Informe a data de nascimento'),
              InkWell(
                onTap: _pickDate,
                borderRadius: BorderRadius.circular(8),
                child: Container(
                  width: double.infinity,
                  padding: const EdgeInsets.symmetric(
                    vertical: 16,
                    horizontal: 12,
                  ),
                  decoration: BoxDecoration(
                    border: Border.all(color: colorScheme.outline),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: ValueListenableBuilder<DateTime?>(
                    valueListenable: _selectedDate,
                    builder: (context, value, child) {
                      return Row(
                        children: [
                          Icon(
                            Icons.calendar_today_outlined,
                            color: colorScheme.onSurfaceVariant,
                          ),
                          const SizedBox(width: 12),
                          Text(
                            value != null
                                ? DateFormat(
                                    'dd/MM/yyyy',
                                  ).format(value)
                                : 'Selecione a data',
                            style: textTheme.bodyLarge?.copyWith(
                              color: value != null
                                  ? colorScheme.onSurface
                                  : colorScheme.onSurfaceVariant,
                              fontWeight: FontWeight.w700,
                            ),
                          ),
                        ],
                      );
                    },
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
      bottomNavigationBar: ListenableBuilder(
        listenable: Listenable.merge([_isDisabled, _viewmodel.submitBirthDate]),
        builder: (context, _) {
          return DoubleBottomButton(
            leftButtonLabel: 'Voltar',
            rightButtonLabel: 'Continuar',
            leftOnPressed: _navBack,
            rightOnPressed: _submitBirthDate,
            isRightEnabled:
                !_isDisabled.value && !_viewmodel.submitBirthDate.isRunning,
            rightButtonIcon: const Icon(Icons.arrow_forward_ios),
          );
        },
      ),
    );
  }

  void _navBack() => context.pop();

  void _navToEmail() => context.pushNamed(RegisterRoutes.email.routeName);

  Future<void> _pickDate() async {
    final picked = await showDatePicker(
      context: context,
      initialDate: _selectedDate.value ?? DateTime(2000),
      firstDate: _minDate,
      lastDate: _maxDate,
    );
    if (picked != null) {
      _selectedDate.value = picked;
      _isDisabled.value = false;
    }
  }

  Future<void> _submitBirthDate() async {
    final selectedDate = _selectedDate.value;
    if (selectedDate == null) return;

    if (selectedDate.age < 18) {
      AppSnackbar.show(
        context,
        message: 'Você deve ser maior de 18 anos para criar uma conta.',
        type: SnackbarType.error,
      );
      return;
    }

    await _viewmodel.submitBirthDate.execute(selectedDate);

    final result = _viewmodel.submitBirthDate.result!;
    if (result.isFailure) {
      if (!mounted) return;
      AppSnackbar.show(
        context,
        message: result.error!.message,
        type: SnackbarType.error,
      );
      return;
    }

    _navToEmail();
  }

  void _initialize() {
    final initialDate = _viewmodel.state?.birthDate;
    if (initialDate != null) {
      _selectedDate.value = initialDate;
      _isDisabled.value = _selectedDate.value!.age < 18;
    }
  }
}
