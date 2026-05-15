import 'package:flutter/material.dart';

class NoTransactionsCard extends StatelessWidget {
  const NoTransactionsCard({
    super.key,
    required this.greyColor,
  });

  final Color greyColor;

  @override
  Widget build(BuildContext context) {
    return Center(
      child: Column(
        spacing: 16,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.receipt_long_rounded,
            size: 64,
            color: greyColor,
          ),

          const Text(
            'Nenhuma transação encontrada',
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w500,
            ),
          ),
        ],
      ),
    );
  }
}
