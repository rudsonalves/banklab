import 'package:flutter/material.dart';

class BalanceTile extends StatelessWidget {
  final String balanceLabel;
  final String supportingLabel;
  final bool isLoading;

  const BalanceTile({
    super.key,
    required this.balanceLabel,
    required this.supportingLabel,
    required this.isLoading,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Container(
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [
            colorScheme.primary,
            colorScheme.primaryContainer,
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(24),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Saldo em conta',
            style: Theme.of(context).textTheme.titleMedium?.copyWith(
              color: colorScheme.onPrimary,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            balanceLabel,
            style: Theme.of(context).textTheme.headlineMedium?.copyWith(
              color: colorScheme.onPrimary,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 12),
          if (!isLoading) ...[
            Text(
              supportingLabel,
              style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                color: colorScheme.onPrimary.withValues(alpha: 0.86),
              ),
            ),
          ] else ...[
            Padding(
              padding: const EdgeInsets.symmetric(vertical: 7.0),
              child: LinearProgressIndicator(
                borderRadius: BorderRadius.circular(999),
                minHeight: 6,
              ),
            ),
          ],
        ],
      ),
    );
  }
}
