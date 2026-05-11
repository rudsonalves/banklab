import 'package:flutter/material.dart';
import 'package:money2/money2.dart';

import '/core/resources/app_currencies.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';

class BalanceCard extends StatelessWidget {
  final Stream<BalanceResponseDto> balance;
  final bool isVisible;
  final VoidCallback? onToggleVisibility;
  final VoidCallback? onTapAccount;
  final AccountSummaryResponseDto? selectedAccount;

  const BalanceCard({
    super.key,
    required this.balance,
    this.isVisible = true,
    this.onToggleVisibility,
    this.onTapAccount,
    this.selectedAccount,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Card(
      margin: EdgeInsets.zero,
      color: colorScheme.onPrimary,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Saldo',
                        style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                          fontWeight: FontWeight.w700,
                        ),
                      ),
                      const SizedBox(height: 8),
                      StreamBuilder<BalanceResponseDto>(
                        stream: balance,
                        builder: (context, snapshot) {
                          final balanceValue =
                              snapshot.data?.balance ??
                              Money.fromIntWithCurrency(0, appCurrency);

                          return Text(
                            isVisible ? balanceValue.format() : '••••••',
                            style: Theme.of(context).textTheme.headlineMedium,
                            overflow: TextOverflow.ellipsis,
                          );
                        },
                      ),
                    ],
                  ),
                ),
                SizedBox(
                  width: 48,
                  height: 48,
                  child: IconButton(
                    onPressed: onToggleVisibility,
                    icon: Icon(
                      isVisible ? Icons.visibility : Icons.visibility_off,
                    ),
                  ),
                ),
              ],
            ),
            // const SizedBox(height: 12),
            InkWell(
              onTap: onTapAccount,
              borderRadius: BorderRadius.circular(12),
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 4),
                child: Row(
                  children: [
                    Expanded(
                      child: Text(
                        selectedAccount != null
                            ? 'Conta: ${selectedAccount!.branch} - ${selectedAccount!.number}'
                            : 'Conta: •••• - ••••••',
                        style: Theme.of(context).textTheme.bodyMedium,
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    const SizedBox(width: 4),
                    if (onTapAccount != null)
                      const Icon(
                        Icons.arrow_drop_down,
                        size: 16,
                      ),
                    const SizedBox(width: 10),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
