import 'package:flutter/material.dart';

enum SnackbarType {
  success,
  error,
  info,
}

class AppSnackbar {
  static void show(
    BuildContext context, {
    String? title,
    required String message,
    SnackbarType type = SnackbarType.info,
    Duration duration = const Duration(seconds: 3),
  }) {
    final backgroundColor = switch (type) {
      SnackbarType.error => Theme.of(context).colorScheme.error,
      SnackbarType.success => Theme.of(context).colorScheme.primary,
      SnackbarType.info => Theme.of(context).colorScheme.secondary,
    };

    ScaffoldMessenger.of(context)
      ..hideCurrentSnackBar()
      ..showSnackBar(
        SnackBar(
          duration: duration,
          backgroundColor: backgroundColor,
          shape: RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(12),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            crossAxisAlignment: CrossAxisAlignment.center,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              if (title != null)
                Padding(
                  padding: const EdgeInsets.only(bottom: 8.0),
                  child: Text(
                    title,
                    style: Theme.of(context).textTheme.titleLarge!.copyWith(
                      color: Theme.of(context).colorScheme.onPrimary,
                    ),
                  ),
                ),
              Text(message),
            ],
          ),
          behavior: SnackBarBehavior.floating,
        ),
      );
  }
}
