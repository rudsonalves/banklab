import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import 'viewmodel/splash_viewmodel.dart';

class SplashPage extends StatefulWidget {
  final SplashViewmodel viewModel;

  const SplashPage({
    super.key,
    required this.viewModel,
  });

  @override
  State<SplashPage> createState() => _SplashPageState();
}

class _SplashPageState extends State<SplashPage>
    with SingleTickerProviderStateMixin {
  static const Duration _minimumSplashDuration = Duration(milliseconds: 2200);

  SplashViewmodel get _viewModel => widget.viewModel;

  late final AnimationController _animationController;
  late final Animation<double> _logoScale;
  late final Animation<double> _logoOpacity;
  DateTime? _startedAt;

  @override
  void initState() {
    super.initState();

    _startedAt = DateTime.now();

    _animationController = AnimationController(
      vsync: this,
      duration: const Duration(milliseconds: 2200),
    );

    _logoScale = Tween<double>(begin: 0.8, end: 1).animate(
      CurvedAnimation(
        parent: _animationController,
        curve: Curves.easeOutCubic,
      ),
    );

    _logoOpacity = Tween<double>(begin: 0, end: 1).animate(
      CurvedAnimation(
        parent: _animationController,
        curve: const Interval(0.1, 0.75, curve: Curves.easeOut),
      ),
    );

    _viewModel.initialize.addListener(_onInitializeChanged);
    _animationController.forward();
    _viewModel.initialize.execute();
  }

  @override
  void dispose() {
    _viewModel.initialize.removeListener(_onInitializeChanged);
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;

    return Scaffold(
      body: Container(
        width: double.infinity,
        height: double.infinity,
        decoration: BoxDecoration(
          gradient: LinearGradient(
            colors: [
              colorScheme.surface,
              colorScheme.surfaceContainerHighest.withValues(alpha: 0.7),
            ],
            begin: Alignment.topCenter,
            end: Alignment.bottomCenter,
          ),
        ),
        child: Center(
          child: AnimatedBuilder(
            animation: _animationController,
            builder: (context, _) => Opacity(
              opacity: _logoOpacity.value,
              child: Transform.scale(
                scale: _logoScale.value,
                child: Card(
                  color: colorScheme.onPrimary,
                  child: Padding(
                    padding: const EdgeInsets.all(12),
                    child: Image.asset(
                      'assets/images/brand.png',
                      width: 200,
                    ),
                  ),
                ),
              ),
            ),
          ),
        ),
      ),
    );
  }

  Future<void> _onInitializeChanged() async {
    final command = _viewModel.initialize;
    if (!mounted || command.isRunning) return;

    final startedAt = _startedAt;
    if (startedAt != null) {
      final elapsed = DateTime.now().difference(startedAt);
      final remaining = _minimumSplashDuration - elapsed;
      if (remaining > Duration.zero) {
        await Future<void>.delayed(remaining);
      }
    }

    if (!mounted) return;

    if (command.isSuccess) {
      context.goNamed(
        AuthRoutes.shortLogin.name,
        extra: command.value,
      );
      return;
    }

    context.goNamed(AuthRoutes.login.name);
  }
}
