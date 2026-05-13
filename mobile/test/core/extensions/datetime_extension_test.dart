import 'package:bankflow/core/extensions/datetime_extension.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:intl/intl.dart';

void main() {
  late String? previousLocale;

  setUp(() {
    previousLocale = Intl.defaultLocale;
    Intl.defaultLocale = 'en_US';
  });

  tearDown(() {
    Intl.defaultLocale = previousLocale;
  });

  group('DateParser.parseOrNull', () {
    test('returns null when input is null or blank', () {
      expect(DateParser.parseOrNull(null), isNull);
      expect(DateParser.parseOrNull(''), isNull);
      expect(DateParser.parseOrNull('   '), isNull);
    });

    test('returns parsed DateTime for valid ISO string', () {
      expect(
        DateParser.parseOrNull('2026-05-13T14:30:00Z'),
        DateTime.parse('2026-05-13T14:30:00Z'),
      );
    });

    test('returns null for invalid date string', () {
      expect(DateParser.parseOrNull('invalid-date'), isNull);
    });
  });

  group('DateParser.parseOrNow', () {
    test('returns parsed DateTime for valid ISO string', () {
      expect(
        DateParser.parseOrNow('2026-05-13T14:30:00Z'),
        DateTime.parse('2026-05-13T14:30:00Z'),
      );
    });

    test('returns current time when input is null, blank, or invalid', () {
      final before = DateTime.now();

      final fromNull = DateParser.parseOrNow(null);
      final fromBlank = DateParser.parseOrNow('   ');
      final fromInvalid = DateParser.parseOrNow('invalid-date');

      final after = DateTime.now();

      expect(fromNull.isBefore(before), isFalse);
      expect(fromNull.isAfter(after), isFalse);
      expect(fromBlank.isBefore(before), isFalse);
      expect(fromBlank.isAfter(after), isFalse);
      expect(fromInvalid.isBefore(before), isFalse);
      expect(fromInvalid.isAfter(after), isFalse);
    });
  });

  group('DateTimeExtensions', () {
    testWidgets('format uses locale from BuildContext', (tester) async {
      BuildContext? context;

      await tester.pumpWidget(
        WidgetsApp(
          color: const Color(0xFFFFFFFF),
          locale: const Locale('en', 'US'),
          supportedLocales: const [Locale('en', 'US')],
          builder: (ctx, _) {
            context = ctx;
            return const SizedBox.shrink();
          },
        ),
      );

      final date = DateTime.utc(2026, 5, 13, 14, 30);

      expect(context, isNotNull);
      expect(date.format(context!), 'May 13, 2026');
    });

    test('formatMonthLabel formats as full month and year', () {
      final date = DateTime(2026, 5, 13);

      expect(date.formatMonthLabel, 'May 2026');
    });

    test('formatDayLabel formats as dd/MM/yyyy', () {
      final date = DateTime(2026, 5, 3);

      expect(date.formatDayLabel, '03/05/2026');
    });

    test('formatHour formats as HH:mm', () {
      final date = DateTime(2026, 5, 13, 9, 5);

      expect(date.formatHour, '09:05');
    });
  });
}
