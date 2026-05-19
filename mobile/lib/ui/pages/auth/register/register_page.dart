import 'package:flutter/material.dart';

import '/ui/components/base/safe_scaffold.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterCpfPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterCpfPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterCpfPage> createState() => _RegisterCpfPageState();
}

class _RegisterCpfPageState extends State<RegisterCpfPage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;

  @override
  void initState() {
    super.initState();

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
          padding: const EdgeInsets.all(24),
          child: Container(),
        ),
      ),
    );
  }
}
