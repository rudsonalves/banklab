import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import '/uis/core/base/safe_scaffold.dart';
import '/uis/core/buttons/big_button.dart';

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
                  ? 'Seu dinheiro foi enviado com sucesso.'
                  : 'Não foi possível concluir a transferência. Tente novamente.',
              textAlign: TextAlign.center,
              style: Theme.of(context).textTheme.bodyLarge,
            ),
            if (isSuccess && transactionReference != null) ...[
              const SizedBox(height: 24),
              Text(
                'Referência da transação',
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.titleSmall,
              ),
              const SizedBox(height: 6),
              Text(
                transactionReference!,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.bodyMedium,
              ),
            ],
          ],
        ),
      ),
      bottomNavigationBar: isSuccess
          ? Row(
              children: [
                Expanded(
                  child: TextButton(
                    onPressed: () => _goBack(context),
                    style: TextButton.styleFrom(
                      padding: const EdgeInsets.symmetric(vertical: 14),
                      textStyle: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.w700,
                      ),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(12),
                      ),
                    ),
                    child: const Text('Voltar'),
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: BigButton(
                    label: 'Ver comprovante',
                    onPressed: () => _showReceipt(context),
                    leftIcon: Icons.receipt_long_rounded,
                    enabled: transactionReference != null,
                  ),
                ),
              ],
            )
          : BigButton(
              label: 'Voltar',
              onPressed: () => _goBack(context),
              enabled: true,
            ),
    );
  }

  void _goBack(BuildContext context) {
    context.goNamed(HomeRoutes.home.name);
  }

  void _showReceipt(BuildContext context) {
    final reference = transactionReference;
    if (reference == null) return;

    showModalBottomSheet<void>(
      context: context,
      showDragHandle: true,
      builder: (context) => Padding(
        padding: const EdgeInsets.fromLTRB(16, 8, 16, 24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.stretch,
          children: [
            Text(
              'Comprovante',
              style: Theme.of(context).textTheme.titleLarge,
            ),
            const SizedBox(height: 12),
            Text(
              'Referência da transação',
              style: Theme.of(context).textTheme.titleSmall,
            ),
            const SizedBox(height: 4),
            Text(
              reference,
              style: Theme.of(context).textTheme.bodyMedium,
            ),
            const SizedBox(height: 16),
            BigButton(
              label: 'Fechar',
              onPressed: () => Navigator.of(context).pop(),
              enabled: true,
            ),
          ],
        ),
      ),
    );
  }
}
