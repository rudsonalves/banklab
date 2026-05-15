import 'package:flutter/material.dart';

class LoadStatementError extends StatelessWidget {
  final bool isLoading;
  final VoidCallback onRetry;

  const LoadStatementError({
    super.key,
    required this.isLoading,
    required this.onRetry,
  });

  @override
  Widget build(BuildContext context) {
    final greyColor = Theme.of(context).colorScheme.onSurface;

    return Center(
      child: Column(
        spacing: 16,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.error_outline,
            size: 64,
            color: greyColor,
          ),

          const Text(
            'Não foi possível carregar o extrato',
            textAlign: TextAlign.center,
            style: TextStyle(
              fontSize: 16,
              fontWeight: FontWeight.w500,
            ),
          ),
          const SizedBox(height: 8),
          ElevatedButton(
            onPressed: isLoading ? null : onRetry,
            child: isLoading
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(
                      strokeWidth: 2,
                    ),
                  )
                : const Text('Tentar novamente'),
          ),
        ],
      ),
    );
  }
}
