import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/ui/components/base/safe_scaffold.dart';
import '../../../core/routing/routes.dart';
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

  void _navToBirthdate() => context.pushNamed(RegisterRoutes.birthDate.name);
}
