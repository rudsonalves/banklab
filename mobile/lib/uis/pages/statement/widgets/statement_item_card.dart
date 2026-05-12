import 'package:flutter/material.dart';

import '/data/services/apis/account/dtos/statement_response_dto.dart';
import '/uis/core/transaction/transaction_movement.dart';

class StatementItemCard extends StatelessWidget {
  final StatementItemDto item;
  // final bool showConsolidatedBalance;
  final String hourLabel;
  final VoidCallback? onTap;

  const StatementItemCard({
    super.key,
    required this.item,
    // required this.showConsolidatedBalance,
    required this.hourLabel,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = Theme.of(context).colorScheme;
    final movement = TransactionMovement.fromType(item.type);
    final amountColor = movement.isDebit
        ? colorScheme.error
        : const Color(0xFF2E7D5B);
    final description = item.description.trim();
    final reference = item.referenceId?.trim();
    final hasDetails =
        onTap != null && reference != null && reference.isNotEmpty;

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          movement.label,
                          style: theme.textTheme.titleSmall?.copyWith(
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                        if (description.isNotEmpty) ...[
                          const SizedBox(height: 4),
                          Text(
                            description,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: theme.textTheme.bodyMedium?.copyWith(
                              color: colorScheme.onSurfaceVariant,
                            ),
                          ),
                        ],
                        const SizedBox(height: 4),
                        Text(
                          hourLabel,
                          style: theme.textTheme.bodySmall?.copyWith(
                            color: colorScheme.onSurfaceVariant,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 12),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.end,
                    children: [
                      Text(
                        '${movement.sign}${item.amount.format()}',
                        style: theme.textTheme.titleSmall?.copyWith(
                          fontWeight: FontWeight.w700,
                          color: amountColor,
                        ),
                      ),
                      // if (showConsolidatedBalance) ...[
                      //   const SizedBox(height: 6),
                      //   Text(
                      //     'Saldo após',
                      //     style: theme.textTheme.labelSmall?.copyWith(
                      //       color: colorScheme.primary,
                      //       fontWeight: FontWeight.w700,
                      //     ),
                      //   ),
                      //   Text(
                      //     item.balanceAfter.format(),
                      //     style: theme.textTheme.bodySmall?.copyWith(
                      //       color: colorScheme.primary,
                      //       fontWeight: FontWeight.w700,
                      //     ),
                      //   ),
                      // ],
                    ],
                  ),
                  if (hasDetails) ...[
                    const SizedBox(width: 4),
                    Icon(
                      Icons.chevron_right_rounded,
                      color: colorScheme.onSurfaceVariant,
                      size: 20,
                    ),
                  ],
                ],
              ),
              // if (reference != null && reference.isNotEmpty) ...[
              //   const SizedBox(height: 12),
              //   Text(
              //     'Ref: $reference',
              //     maxLines: 1,
              //     overflow: TextOverflow.ellipsis,
              //     style: theme.textTheme.labelSmall?.copyWith(
              //       color: colorScheme.onSurfaceVariant,
              //     ),
              //   ),
              // ],
            ],
          ),
        ),
      ),
    );
  }
}
