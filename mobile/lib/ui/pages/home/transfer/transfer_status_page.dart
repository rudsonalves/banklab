import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/ui/components/base/safe_scaffold.dart';
import '/ui/components/buttons/big_button.dart';
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
                  ? Theme.of(context).colorScheme.primaryContainer
                  : Theme.of(context).colorScheme.errorContainer,
              child: Icon(
                isSuccess ? Icons.check_rounded : Icons.close_rounded,
                size: 44,
                color: isSuccess
                    ? Theme.of(context).colorScheme.primary
                    : Theme.of(context).colorScheme.error,
              ),
            ),
            const SizedBox(height: 20),
            Text(
              isSuccess ? 'Transferência realizada' : 'Transferência falhou',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.headlineSmall,
            ),
            const SizedBox(height: 12),
            Text(
              isSuccess
                  ? 'Sua transferência foi realizada com sucesso.'
                  : 'Não foi possível concluir a transferência. Tente novamente.',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
          ],
        ),
      ),
      bottomNavigationBar: isSuccess
          ? DoubleBottomButton(
              leftButtonLabel: 'Voltar',
              rightButtonLabel: 'Comprovante',
              leftOnPressed: () => _navBack(context),
              rightOnPressed: () => _showReceipt(context),
              isRightEnabled: transactionReference != null,
              rightButtonIcon: const Icon(Icons.receipt_long_rounded),
            )
          : BigButton(
              label: 'Voltar',
              onPressed: () => _navBack(context),
              enabled: true,
            ),
    );
  }

  void _navBack(BuildContext context) {
    context.goNamed(BaseRoutes.home.name);
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
