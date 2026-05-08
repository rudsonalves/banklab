import 'package:bankflow/core/extensions/string.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('StringExtension.onlyNumbers', () {
    test('should keep only digits from mixed content', () {
      expect('abc123.45-6'.onlyNumbers, '123456');
    });

    test('should return empty string when there are no digits', () {
      expect('abc-xyz'.onlyNumbers, isEmpty);
    });

    test('should keep original value when all chars are digits', () {
      expect('001234'.onlyNumbers, '001234');
    });
  });

  group('StringExtension.isCPF', () {
    test('should validate formatted CPF', () {
      expect('529.982.247-25'.isValidCpf, isTrue);
    });

    test('should validate unformatted CPF', () {
      expect('52998224725'.isValidCpf, isTrue);
    });

    test('should invalidate CPF with incorrect check digits', () {
      expect('529.982.247-26'.isValidCpf, isFalse);
    });

    test('should invalidate CPF with invalid length', () {
      expect('5299822472'.isValidCpf, isFalse);
      expect('529982247250'.isValidCpf, isFalse);
    });

    test('should invalidate CPF with all repeated digits', () {
      expect('111.111.111-11'.isValidCpf, isFalse);
    });
  });

  group('StringExtension.isValidCNPJ', () {
    test('should validate formatted CNPJ', () {
      expect('04.252.011/0001-10'.isValidCnpj, isTrue);
    });

    test('should validate unformatted CNPJ', () {
      expect('04252011000110'.isValidCnpj, isTrue);
    });

    test('should invalidate CNPJ with incorrect check digits', () {
      expect('04.252.011/0001-11'.isValidCnpj, isFalse);
    });

    test('should invalidate CNPJ with invalid length', () {
      expect('0425201100011'.isValidCnpj, isFalse);
      expect('042520110001100'.isValidCnpj, isFalse);
    });

    test('should invalidate CNPJ with all repeated digits', () {
      expect('11.111.111/1111-11'.isValidCnpj, isFalse);
    });
  });
}
