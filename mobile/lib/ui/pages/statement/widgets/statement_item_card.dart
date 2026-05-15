import 'package:flutter/material.dart';

import '/data/services/apis/account/dtos/statement_response_dto.dart';
import '/ui/components/transaction/transaction_movement.dart';

class StatementItemCard extends StatelessWidget {
  final StatementItemDto item;
  final String hourLabel;
  final VoidCallback? onTap;

  const StatementItemCard({
    super.key,
    required this.item,
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
            spacing: 8,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(
                    child: Text(
                      movement.label,
                      style: theme.textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                  Text(
                    '${movement.sign}${item.amount.format()}',
                    style: theme.textTheme.titleSmall?.copyWith(
                      fontWeight: FontWeight.w700,
                      color: amountColor,
                    ),
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

              Row(
                spacing: 4,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Expanded(
                    child: Text(
                      description,
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.bodyMedium?.copyWith(
                        color: colorScheme.onSurfaceVariant,
                      ),
                    ),
                  ),

                  Padding(
                    padding: const EdgeInsets.only(right: 8),
                    child: Text(
                      hourLabel,
                      style: theme.textTheme.bodySmall?.copyWith(
                        color: colorScheme.onSurfaceVariant,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
