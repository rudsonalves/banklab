import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/big_button.dart';

class RegisterStatusPage extends StatelessWidget {
  final bool isSuccess;

  const RegisterStatusPage({
    super.key,
    required this.isSuccess,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Status do registro'),
      ),
      body: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            CircleAvatar(
              radius: 42,
              backgroundColor: isSuccess
                  ? colorScheme.primaryContainer
                  : colorScheme.errorContainer,
              child: Icon(
                isSuccess ? Icons.check_rounded : Icons.close_rounded,
                size: 44,
                color: isSuccess ? colorScheme.primary : colorScheme.error,
              ),
            ),
            const SizedBox(height: 20),
            Text(
              isSuccess ? 'Registro concluido' : 'Falha no registro',
              textAlign: TextAlign.center,
              style: textTheme.headlineSmall,
            ),
            const SizedBox(height: 12),
            Text(
              isSuccess
                  ? 'Sua conta foi criada com sucesso.'
                  : 'Nao foi possivel concluir o registro. Tente novamente.',
              textAlign: TextAlign.center,
              style: textTheme.bodyLarge,
            ),
          ],
        ),
      ),
      bottomNavigationBar: isSuccess
          ? BigButton(
              label: 'Entrar',
              onPressed: () => _goToLogin(context),
              rightIcon: const Icon(Icons.login_rounded),
            )
          : BigButton(
              label: 'Tentar novamente',
              onPressed: () => _retryPassword(context),
              rightIcon: const Icon(Icons.refresh_rounded),
            ),
    );
  }

  void _goToLogin(BuildContext context) {
    context.goNamed(AuthRoutes.login.name);
  }

  void _retryPassword(BuildContext context) {
    context.pop();
  }
}
