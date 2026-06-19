import 'package:bankflow/ui/pages/auth/widgets/installation_limit_dialog.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('shows approved installation limit copy and action', (
    tester,
  ) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Builder(
          builder: (context) => Scaffold(
            body: TextButton(
              onPressed: () => showInstallationLimitDialog(context),
              child: const Text('Abrir'),
            ),
          ),
        ),
      ),
    );

    await tester.tap(find.text('Abrir'));
    await tester.pumpAndSettle();

    expect(find.text(installationLimitDialogTitle), findsOneWidget);
    expect(
      find.textContaining(installationLimitDialogPrimaryMessage),
      findsOneWidget,
    );
    expect(
      find.textContaining(installationLimitDialogSecondaryMessage),
      findsOneWidget,
    );
    expect(find.text(installationLimitDialogButtonLabel), findsOneWidget);

    await tester.tap(find.text(installationLimitDialogButtonLabel));
    await tester.pumpAndSettle();

    expect(find.text(installationLimitDialogTitle), findsNothing);
  });
}
