import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/router.dart';
import '/uis/core/themes/material_theme.dart';
import '/uis/core/themes/text_theme.dart';

class AppWidget extends StatefulWidget {
  const AppWidget({super.key});

  @override
  State<AppWidget> createState() => _AppWidgetState();
}

class _AppWidgetState extends State<AppWidget> {
  late final MaterialTheme _materialTheme;
  final GoRouter _router = router();

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    // Initialize once with context-dependent resources
    final textTheme = createTextTheme(context, "Quicksand", "EB Garamond");

    _materialTheme = MaterialTheme(textTheme);
  }

  @override
  Widget build(BuildContext context) {
    final brightness = View.of(context).platformDispatcher.platformBrightness;

    final baseTheme = _resolveBaseTheme(brightness);
    final theme = _buildAppTheme(baseTheme);

    return MaterialApp.router(
      theme: theme,
      debugShowCheckedModeBanner: false,
      routerConfig: _router,
    );
  }

  ThemeData _resolveBaseTheme(Brightness brightness) {
    return brightness == Brightness.light
        ? _materialTheme.light()
        : _materialTheme.dark();
  }

  ThemeData _buildAppTheme(ThemeData base) {
    return base.copyWith(
      appBarTheme: base.appBarTheme.copyWith(
        backgroundColor: base.colorScheme.primaryContainer,
        foregroundColor: base.colorScheme.onPrimaryContainer,
        titleTextStyle: base.textTheme.titleLarge?.copyWith(
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}
