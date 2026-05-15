import 'package:flutter/material.dart';
import 'package:money2/money2.dart';

class DayHeader extends StatelessWidget {
  final String label;
  final Color greyColor;
  final Money? balance;

  const DayHeader({
    super.key,
    required this.label,
    required this.greyColor,
    this.balance,
  });

  @override
  Widget build(BuildContext context) {
    final colorPrimary = Theme.of(context).colorScheme.primary;

    return Padding(
      padding: const EdgeInsets.only(bottom: 8, right: 24),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(
              fontSize: 14,
              fontWeight: FontWeight.w600,
              color: greyColor,
            ),
          ),
          if (balance != null)
            Text(
              balance!.format(),
              style: TextStyle(
                fontSize: 14,
                fontWeight: FontWeight.w700,
                color: colorPrimary,
              ),
            ),
        ],
      ),
    );
  }
}
