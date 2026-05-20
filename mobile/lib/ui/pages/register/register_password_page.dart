import 'package:flutter/material.dart';

import '/ui/components/base/safe_scaffold.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterPasswordPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterPasswordPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterPasswordPage> createState() => _RegisterPasswordPageState();
}

class _RegisterPasswordPageState extends State<RegisterPasswordPage> {
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
        title: const Text('Criar conta'),
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
