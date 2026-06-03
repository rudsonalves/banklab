import 'package:bankflow/ui/components/input_text/otp_input.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  testWidgets('fills all fields when pasting a full OTP', (tester) async {
    String? changedValue;
    String? completedValue;

    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: OtpInput(
            onChanged: (value) => changedValue = value,
            onCompleted: (value) => completedValue = value,
          ),
        ),
      ),
    );

    await tester.enterText(
      find.byType(TextField, skipOffstage: false),
      '123456',
    );
    await tester.pump();

    expect(changedValue, '123456');
    expect(completedValue, '123456');

    for (final digit in '123456'.characters) {
      expect(find.text(digit), findsOneWidget);
    }
  });

  testWidgets('fills all fields when long pressing with OTP in clipboard', (
    tester,
  ) async {
    String? changedValue;
    String? completedValue;

    TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
        .setMockMethodCallHandler(SystemChannels.platform, (call) async {
          if (call.method == 'Clipboard.getData') {
            return {'text': '987654'};
          }

          return null;
        });
    addTearDown(() {
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
          .setMockMethodCallHandler(SystemChannels.platform, null);
    });

    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: OtpInput(
            onChanged: (value) => changedValue = value,
            onCompleted: (value) => completedValue = value,
          ),
        ),
      ),
    );

    await tester.longPress(find.byType(OtpInput));
    await tester.pump();

    expect(changedValue, '987654');
    expect(completedValue, '987654');

    for (final digit in '987654'.characters) {
      expect(find.text(digit), findsOneWidget);
    }
  });
}
