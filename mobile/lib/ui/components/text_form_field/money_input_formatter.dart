import 'package:flutter/services.dart';
import 'package:money2/money2.dart';

import '/core/extensions/string.dart';
import '/core/resources/app_currencies.dart';

class MoneyInputFormatter extends TextInputFormatter {
  final Currency currency;
  final int maxDigits;

  MoneyInputFormatter({
    Currency? currency,
    this.maxDigits = 15,
  }) : currency = currency ?? appCurrency;

  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    // Remove everything except digits
    final digits = newValue.text.onlyNumbers;

    if (digits.length > maxDigits) {
      return oldValue;
    }

    // POS behavior: empty input = zero
    final int minorUnits = digits.isEmpty ? 0 : int.parse(digits);

    final money = Money.fromIntWithCurrency(minorUnits, currency);

    final formatted = money.format(currency.pattern);

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(
        offset: formatted.length,
      ),
    );
  }
}
