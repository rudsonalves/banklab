import 'package:bankflow/core/routing/extensions/context_extencions.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('popUntil removes routes above the requested named route', (
    tester,
  ) async {
    final router = GoRouter(
      initialLocation: '/home',
      routes: [
        GoRoute(
          path: '/home',
          name: 'home',
          builder: (context, state) => Scaffold(
            body: TextButton(
              onPressed: () => context.pushNamed('first'),
              child: const Text('Open first'),
            ),
          ),
        ),
        GoRoute(
          path: '/first',
          name: 'first',
          builder: (context, state) => Scaffold(
            body: TextButton(
              onPressed: () => context.pushNamed('second'),
              child: const Text('Open second'),
            ),
          ),
        ),
        GoRoute(
          path: '/second',
          name: 'second',
          builder: (context, state) => Scaffold(
            body: TextButton(
              onPressed: () => context.popUntil('home'),
              child: const Text('Return home'),
            ),
          ),
        ),
      ],
    );
    addTearDown(router.dispose);

    await tester.pumpWidget(MaterialApp.router(routerConfig: router));

    await tester.tap(find.text('Open first'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Open second'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Return home'));
    await tester.pumpAndSettle();

    expect(find.text('Open first'), findsOneWidget);
    expect(router.canPop(), isFalse);
  });

  testWidgets('popUntil preserves the first route when the name is not found', (
    tester,
  ) async {
    final router = GoRouter(
      initialLocation: '/home',
      routes: [
        GoRoute(
          path: '/home',
          name: 'home',
          builder: (context, state) => Scaffold(
            body: TextButton(
              onPressed: () => context.pushNamed('child'),
              child: const Text('Open child'),
            ),
          ),
        ),
        GoRoute(
          path: '/child',
          name: 'child',
          builder: (context, state) => Scaffold(
            body: TextButton(
              onPressed: () => context.popUntil('unknown'),
              child: const Text('Return'),
            ),
          ),
        ),
      ],
    );
    addTearDown(router.dispose);

    await tester.pumpWidget(MaterialApp.router(routerConfig: router));

    await tester.tap(find.text('Open child'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Return'));
    await tester.pumpAndSettle();

    expect(find.text('Open child'), findsOneWidget);
    expect(router.canPop(), isFalse);
  });
}
