import 'package:flutter/widgets.dart';
import 'package:intl/intl.dart';

extension DateTimeExtensions on DateTime {
  String format(BuildContext context, [String pattern = 'yMMMd']) {
    final locale = Localizations.localeOf(context).toString();
    return DateFormat(pattern, locale).format(this);
  }

  String get formatMonthLabel => DateFormat('MMMM yyyy').format(this);

  String get formatDayLabel => DateFormat('dd/MM/yyyy').format(this);

  String get formatHour => DateFormat('HH:mm').format(this);

  String get dateOnly => toIso8601String().split('T').first;

  int get age {
    final now = DateTime.now();
    int age = now.year - year;
    if (now.month < month || (now.month == month && now.day < day)) {
      age--;
    }
    return age;
  }
}

final class DateParser {
  static DateTime? parseOrNull(String? dateString) {
    final str = dateString?.trim();
    if (str == null || str.isEmpty) return null;

    return DateTime.tryParse(str);
  }

  static DateTime parseOrNow(String? dateString) {
    final str = dateString?.trim();
    if (str == null || str.isEmpty) return DateTime.now();

    return DateTime.tryParse(str) ?? DateTime.now();
  }

  static DateTime? parseDateOnly(String? value) {
    final str = value?.trim();
    if (str == null || str.isEmpty) return null;

    return DateTime.tryParse(str);
  }
}
