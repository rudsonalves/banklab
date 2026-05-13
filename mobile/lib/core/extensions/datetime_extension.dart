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
}
