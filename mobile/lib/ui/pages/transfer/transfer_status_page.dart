import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/double_bottom_buttons.dart';

class TransferStatusPage extends StatelessWidget {
  final bool isSuccess;
  final String? transactionReference;

  const TransferStatusPage({
    super.key,
    required this.isSuccess,
    this.transactionReference,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return SafeScaffold(
      appBar: AppBar(
        title: const Text('Status da transferência'),
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
              isSuccess ? 'Transferência realizada' : 'Transferência falhou',
              textAlign: TextAlign.center,
              style: textTheme.headlineSmall,
            ),
            const SizedBox(height: 12),
            Text(
              isSuccess
                  ? 'Sua transferência foi realizada com sucesso.'
                  : 'Não foi possível concluir a transferência. Tente novamente.',
              textAlign: TextAlign.center,
              style: textTheme.bodyLarge,
            ),
          ],
        ),
      ),
      bottomNavigationBar: isSuccess
          ? DoubleBottomButton(
              leftButtonLabel: 'Fechar',
              rightButtonLabel: 'Comprovante',
              leftOnPressed: () => _navHome(context),
              rightOnPressed: () => _showReceipt(context),
              isRightEnabled: transactionReference != null,
              rightButtonIcon: const Icon(Icons.receipt_long_rounded),
            )
          : DoubleBottomButton(
              leftButtonLabel: 'Cancelar',
              rightButtonLabel: 'Novamente',
              leftOnPressed: () => _navHome(context),
              rightOnPressed: () => _navBack(context),
              isRightEnabled: true,
              rightButtonIcon: const Icon(Icons.refresh_rounded),
            ),
    );
  }

  void _navBack(BuildContext context) => context.pop();

  void _navHome(BuildContext context) {
    context.goNamed(BaseRoutes.home.routeName);
  }

  void _showReceipt(BuildContext context) {
    final reference = transactionReference;
    if (reference == null) return;
    context.pushNamed(
      SharedRoutes.details.name,
      extra: reference,
    );
  }
}
