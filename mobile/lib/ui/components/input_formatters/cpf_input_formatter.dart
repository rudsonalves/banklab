import 'package:flutter/services.dart';

class CpfInputFormatter extends TextInputFormatter {
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
    // Remove tudo que não for dígito
    final digitsOnly = value.replaceAll(RegExp(r'\D'), '');

    // Aplica a formatação: 000.000.000-00
    final buffer = StringBuffer();
    for (int i = 0; i < digitsOnly.length && i < 11; i++) {
      if (i == 3 || i == 6) {
        buffer.write('.'); // Adiciona ponto após o terceiro e sexto dígito
      } else if (i == 9) {
        buffer.write('-'); // Adiciona hífen após o nono dígito
      }
      buffer.write(digitsOnly[i]);
    }

    return buffer.toString();
  }
}
