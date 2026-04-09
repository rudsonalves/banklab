import 'package:flutter/material.dart';

import '/uis/core/themes/material_theme.dart';
import '/uis/core/themes/text_theme.dart';
import 'uis/pages/home/home_page.dart';

class MainApp extends StatefulWidget {
  const MainApp({super.key});

  @override
  State<MainApp> createState() => _MainAppState();
}

class _MainAppState extends State<MainApp> {
  late final MaterialTheme _materialTheme;

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

    return MaterialApp(
      theme: theme,
      debugShowCheckedModeBanner: false,
      home: const HomePage(),
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
