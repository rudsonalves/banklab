import 'package:flutter/services.dart';

import '/core/extensions/string.dart';

class PhoneInputFormatter extends TextInputFormatter {
  @override
  TextEditingValue formatEditUpdate(
    TextEditingValue oldValue,
    TextEditingValue newValue,
  ) {
    final formatted = _format(newValue.text);

    return TextEditingValue(
      text: formatted,
      selection: TextSelection.collapsed(offset: formatted.length),
    );
  }

  String _format(String value) {
    final digits = value.onlyNumbers;

    if (digits.isEmpty) {
      return '';
    }

    final buffer = StringBuffer();

    // DDD
    if (digits.isNotEmpty) {
      buffer.write('(');
    }

    if (digits.length >= 2) {
      buffer.write(digits.substring(0, 2));
      buffer.write(') ');
    } else {
      buffer.write(digits);
      return buffer.toString();
    }

    final number = digits.substring(2);

    // Mobile: 9 digits
    if (number.length > 8) {
      for (int i = 0; i < number.length && i < 9; i++) {
        if (i == 5) {
          buffer.write('-');
        }

        buffer.write(number[i]);
      }
    }
    // Landline: 8 digits
    else {
      for (int i = 0; i < number.length && i < 8; i++) {
        if (i == 4) {
          buffer.write('-');
        }

        buffer.write(number[i]);
      }
    }

    return buffer.toString();
  }
}
