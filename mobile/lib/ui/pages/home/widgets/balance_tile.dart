import 'package:flutter/material.dart';

import '../../../../data/services/apis/account/dtos/account_summary_response_dto.dart';
import '../../../../data/services/apis/account/dtos/balance_response_dto.dart';

class BalanceTile extends StatelessWidget {
  final BalanceResponseDto? balance;
  final AccountSummaryResponseDto? account;
  final bool isLoading;

  const BalanceTile({
    super.key,
    this.balance,
    required this.account,
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
            colorScheme.primary.withValues(alpha: .6),
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
            _balanceLabel(),
            style: Theme.of(context).textTheme.headlineMedium?.copyWith(
              color: colorScheme.onPrimary,
              fontWeight: FontWeight.w700,
            ),
          ),
          const SizedBox(height: 12),
          if (!isLoading) ...[
            Text(
              _supportingLabel(),
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

  String _balanceLabel() {
    if (balance == null) return 'R\$ ---';

    return balance!.balance.format();
  }

  String _supportingLabel() {
    if (isLoading || account == null) return 'Iniciando...';

    return 'Agência: ${account!.branch} - Conta: ${account!.number}';
  }
}
