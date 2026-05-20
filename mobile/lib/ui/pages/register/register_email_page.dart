import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/ui/components/base/safe_scaffold.dart';
import '../../../core/routing/routes.dart';
import 'viewmodel/register_viewmodel.dart';

class RegisterEmailPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterEmailPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterEmailPage> createState() => _RegisterEmailPageState();
}

class _RegisterEmailPageState extends State<RegisterEmailPage> {
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

  void _navToEmailToken() => context.pushNamed(RegisterRoutes.emailToken.name);
}
