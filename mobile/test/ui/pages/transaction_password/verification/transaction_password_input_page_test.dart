import 'package:bankflow/ui/pages/transaction_password/verification/transaction_password_input_page.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('keeps submit disabled until all six digits are entered', (
    tester,
  ) async {
    await _pumpPromptHarness(tester);

    await tester.tap(find.text('Open PIN'));
    await tester.pumpAndSettle();

    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '12345',
    );
    await tester.pump();

    final button = tester.widget<ElevatedButton>(
      find.widgetWithText(ElevatedButton, 'Concluir'),
    );
    expect(button.onPressed, isNull);
    expect(find.text('1'), findsNothing);
    expect(find.text('*'), findsNWidgets(5));
  });

  testWidgets('returns a masked six-digit PIN to the caller', (tester) async {
    await _pumpPromptHarness(tester);

    await tester.tap(find.text('Open PIN'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '123456',
    );
    await tester.pump();

    expect(find.text('1'), findsNothing);
    expect(find.text('*'), findsNWidgets(6));

    await tester.tap(find.text('Concluir'));
    await tester.pumpAndSettle();

    expect(find.text('Result: 123456'), findsOneWidget);

    await tester.tap(find.text('Open PIN'));
    await tester.pumpAndSettle();

    final field = tester.widget<TextField>(
      find.byType(TextField, skipOffstage: false),
    );
    expect(field.controller?.text, isEmpty);
  });

  testWidgets('cancel returns null and a reopened prompt starts empty', (
    tester,
  ) async {
    await _pumpPromptHarness(tester);

    await tester.tap(find.text('Open PIN'));
    await tester.pumpAndSettle();
    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '123456',
    );
    await tester.pump();
    await tester.tap(find.text('Cancelar'));
    await tester.pumpAndSettle();

    expect(find.text('Result: canceled'), findsOneWidget);

    await tester.tap(find.text('Open PIN'));
    await tester.pumpAndSettle();

    final field = tester.widget<TextField>(
      find.byType(TextField, skipOffstage: false),
    );
    expect(field.controller?.text, isEmpty);
  });
}

Future<void> _pumpPromptHarness(WidgetTester tester) async {
  final router = GoRouter(
    initialLocation: '/',
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => const _PromptHarnessPage(),
        routes: [
          GoRoute(
            path: 'pin',
            name: 'pin',
            builder: (context, state) => const TransactionPasswordInputPage(),
          ),
        ],
      ),
    ],
  );
  addTearDown(router.dispose);

  await tester.pumpWidget(MaterialApp.router(routerConfig: router));
  await tester.pumpAndSettle();
}

class _PromptHarnessPage extends StatefulWidget {
  const _PromptHarnessPage();

  @override
  State<_PromptHarnessPage> createState() => _PromptHarnessPageState();
}

class _PromptHarnessPageState extends State<_PromptHarnessPage> {
  String? _result;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Column(
        children: [
          ElevatedButton(
            onPressed: () async {
              final result = await context.pushNamed<String?>('pin');
              if (!mounted) return;
              setState(() => _result = result ?? 'canceled');
            },
            child: const Text('Open PIN'),
          ),
          if (_result != null) Text('Result: $_result'),
        ],
      ),
    );
  }
}
