import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/data/services/apis/contact_verification/enums/contact_verification_channel.dart';
import '/ui/components/base/safe_scaffold.dart';
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
}
