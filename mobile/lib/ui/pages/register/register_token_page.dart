import 'package:flutter/material.dart';

import '/ui/components/base/safe_scaffold.dart';
import 'viewmodel/register_viewmodel.dart';

enum TokenType {
  email,
  phone,
}

class RegisterTokenPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;
  final TokenType tokenType;

  const RegisterTokenPage({
    super.key,
    required this.viewmodel,
    required this.tokenType,
  });

  @override
  State<RegisterTokenPage> createState() => _RegisterTokenPageState();
}

class _RegisterTokenPageState extends State<RegisterTokenPage> {
  RegisterViewmodel get _viewmodel => widget.viewmodel;
  TokenType get _tokenType => widget.tokenType;

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
          padding: const EdgeInsets.all(12),
          child: Container(),
        ),
      ),
    );
  }
}
