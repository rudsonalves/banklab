import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';
import '/core/routing/models/transaction_password_setup_origin.dart';
import 'widgets/information_item.dart';

class IntroductionTransactionPasswordPage extends StatelessWidget {
  final TransactionPasswordSetupOrigin origin;

  const IntroductionTransactionPasswordPage({
    super.key,
    required this.origin,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Senha transacional'),
      ),
      body: SingleChildScrollView(
        padding: const EdgeInsets.symmetric(vertical: 24, horizontal: 16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            CircleAvatar(
              radius: 40,
              backgroundColor: colorScheme.primaryContainer,
              child: Icon(
                Icons.shield_outlined,
                size: 42,
                color: colorScheme.primary,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              'Proteção extra para suas transações',
              textAlign: TextAlign.center,
              style: textTheme.headlineSmall?.copyWith(
                fontWeight: FontWeight.w700,
              ),
            ),
            const SizedBox(height: 12),
            Text(
              'A senha transacional é um código de 6 dígitos usado para '
              'confirmar operações financeiras dentro do app.',
              textAlign: TextAlign.center,
              style: textTheme.bodyLarge?.copyWith(
                color: colorScheme.onSurfaceVariant,
              ),
            ),
            const SizedBox(height: 32),
            const InformationItem(
              icon: Icons.lock_outline,
              title: 'É diferente da senha de acesso',
              description:
                  'Sua senha de acesso permite entrar no app. A senha '
                  'transacional autoriza movimentações na sua conta.',
            ),
            const SizedBox(height: 20),
            const InformationItem(
              icon: Icons.payments_outlined,
              title: 'Confirma operações financeiras',
              description:
                  'Ela poderá ser solicitada ao realizar transferências e '
                  'outras transações que exigem confirmação.',
            ),
            const SizedBox(height: 20),
            const InformationItem(
              icon: Icons.visibility_off_outlined,
              title: 'É pessoal e intransferível',
              description:
                  'Não compartilhe o código e evite sequências fáceis, datas '
                  'ou números que outras pessoas possam descobrir.',
            ),
          ],
        ),
      ),
      bottomNavigationBar: DoubleBottomButton(
        leftButtonLabel: 'Agora não',
        rightButtonLabel: 'Criar senha',
        leftOnPressed: () => _navBack(context),
        rightOnPressed: () => _navNext(context),
        isRightEnabled: true,
        rightButtonIcon: const Icon(Icons.arrow_forward_ios),
      ),
    );
  }

  void _navBack(BuildContext context) {
    context.pop();
  }

  Future<void> _navNext(BuildContext context) async {
    final result = await context.pushNamed(
      TransactionPasswordRoutes.create.routeName,
      extra: origin,
    );

    if (!context.mounted) return;
    if (result == null) context.pop(result);
  }
}
