import 'package:money2/money2.dart';

import '/core/resources/app_currencies.dart';

abstract final class ApiParse {
  const ApiParse._();

  static BigInt toBigInt(
    dynamic value, {
    String fieldName = 'value',
  }) {
    if (value is BigInt) return value;

    if (value is int) {
      return BigInt.from(value);
    }

    if (value is double) {
      if (!value.isFinite || value.truncateToDouble() != value) {
        throw FormatException('Cannot parse BigInt from $fieldName: $value');
      }

      return BigInt.from(value);
    }

    if (value is String) {
      final normalized = value.trim();

      if (normalized.isEmpty) {
        throw FormatException('Cannot parse BigInt from empty $fieldName');
      }

      return BigInt.parse(normalized);
    }

    throw FormatException('Cannot parse BigInt from $fieldName: $value');
  }

  static Money toMoney(
    dynamic value, {
    Currency? currency,
    String fieldName = 'money value',
  }) {
    currency ??= appCurrency;

    final bigIntValue = toBigInt(value, fieldName: fieldName);
    return Money.fromBigIntWithCurrency(bigIntValue, currency);
  }

  static int toInt(Money value) => value.minorUnits.toInt();
}
