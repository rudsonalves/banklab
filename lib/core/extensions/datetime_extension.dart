import 'package:flutter/widgets.dart';
import 'package:intl/intl.dart';

extension DateTimeExtensions on DateTime {
  String format(BuildContext context, [String pattern = 'yMMMd']) {
    final locale = Localizations.localeOf(context).toString();
    return DateFormat(pattern, locale).format(this);
  }

  static DateTime? parseOrNull(String? dateString) {
    if (dateString == null || dateString.isEmpty) return null;

    final str = dateString.toString();
    if (str.isEmpty) return null;

    return DateTime.tryParse(str);
  }
}
