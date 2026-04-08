import 'package:flutter/material.dart';

import '/uis/pages/auth/register/viewmodel/register_viewmodel.dart';

class RegisterPage extends StatefulWidget {
  final RegisterViewmodel viewmodel;

  const RegisterPage({
    super.key,
    required this.viewmodel,
  });

  @override
  State<RegisterPage> createState() => _RegisterPageState();
}

class _RegisterPageState extends State<RegisterPage> {
  @override
  Widget build(BuildContext context) {
    return const Placeholder();
  }
}
