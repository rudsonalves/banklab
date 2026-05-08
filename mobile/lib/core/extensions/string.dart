extension StringExtension on String {
  String get onlyNumbers => replaceAll(RegExp(r'[^\d]'), '');

  bool get isValidCpf {
    final cleaned = onlyNumbers;
    final isValidLength = cleaned.length == 11;

    if (!isValidLength) return false;
    if (RegExp(r'^(\d)\1{10}$').hasMatch(cleaned)) return false;

    final firstNineDigits = cleaned.substring(0, 9);
    final firstCheckDigit = _calculateCPFCheckDigit(firstNineDigits);
    final secondCheckDigit = _calculateCPFCheckDigit(
      firstNineDigits + firstCheckDigit,
    );

    return cleaned.endsWith(firstCheckDigit + secondCheckDigit);
  }

  bool get isValidCnpj {
    final cleaned = onlyNumbers;
    final isValidLength = cleaned.length == 14;

    if (!isValidLength) return false;
    if (RegExp(r'^(\d)\1{13}$').hasMatch(cleaned)) return false;

    final firstTwelveDigits = cleaned.substring(0, 12);
    final firstCheckDigit = _calculateCNPJCheckDigit(firstTwelveDigits);
    final secondCheckDigit = _calculateCNPJCheckDigit(
      firstTwelveDigits + firstCheckDigit,
    );

    return cleaned.endsWith(firstCheckDigit + secondCheckDigit);
  }

  String _calculateCNPJCheckDigit(String digits) {
    final length = digits.length;
    final weights = switch (length) {
      12 => [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2],
      13 => [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2],
      _ => throw ArgumentError('CNPJ digits length must be 12 or 13'),
    };

    final sum = List.generate(
      length,
      (index) => int.parse(digits[index]) * weights[index],
    ).reduce((a, b) => a + b);
    final remainder = sum % 11;
    return (remainder < 2) ? '0' : (11 - remainder).toString();
  }

  String _calculateCPFCheckDigit(String digits) {
    final length = digits.length;
    final sum = List.generate(
      length,
      (index) => int.parse(digits[index]) * (length + 1 - index),
    ).reduce((a, b) => a + b);
    final remainder = sum % 11;
    return (remainder < 2) ? '0' : (11 - remainder).toString();
  }
}
