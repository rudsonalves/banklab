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
            animation: Listenable.merge([
              _animationController,
              _viewModel.initialize,
            ]),
            builder: (context, _) {
              final command = _viewModel.initialize;
              final showRecoverableError =
                  command.isFailure && !command.isRunning;

              return Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  mainAxisSize: MainAxisSize.min,
                  children: [
                    Opacity(
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
                    if (showRecoverableError) ...[
                      const SizedBox(height: 28),
                      ConstrainedBox(
                        constraints: const BoxConstraints(maxWidth: 420),
                        child: Column(
                          spacing: 14,
                          children: [
                            Text(
                              'Não foi possível preparar este app para acesso.',
                              style: Theme.of(context).textTheme.titleMedium,
                              textAlign: TextAlign.center,
                            ),
                            Text(
                              'Verifique o armazenamento do dispositivo e tente novamente.',
                              style: Theme.of(context).textTheme.bodyMedium
                                  ?.copyWith(
                                    color: colorScheme.onSurfaceVariant,
                                  ),
                              textAlign: TextAlign.center,
                            ),
                            FilledButton.icon(
                              onPressed: _retryInitialize,
                              icon: const Icon(Icons.refresh_rounded),
                              label: const Text('Tentar novamente'),
                            ),
                          ],
                        ),
                      ),
                    ],
                  ],
                ),
              );
            },
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

    if (command.isFailure) return;

    if (command.isSuccess) {
      final lastLoginIdentity = command.value?.lastLoginIdentity;
      if (lastLoginIdentity == null) {
        context.goNamed(AuthRoutes.login.routeName);
        return;
      }

      context.goNamed(
        AuthRoutes.shortLogin.name,
        extra: lastLoginIdentity,
      );
      return;
    }

    context.goNamed(AuthRoutes.login.routeName);
  }

  void _retryInitialize() {
    _startedAt = DateTime.now();
    _viewModel.initialize.execute();
  }
}
