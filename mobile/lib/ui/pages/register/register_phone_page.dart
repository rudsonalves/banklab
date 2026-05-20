import 'package:flutter/material.dart';

import '/ui/components/base/safe_scaffold.dart';
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

  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
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
          padding: const EdgeInsets.all(12),
          child: Container(),
        ),
      ),
    );
  }
}
