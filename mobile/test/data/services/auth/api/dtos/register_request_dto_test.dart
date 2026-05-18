import 'package:bankflow/data/services/auth/api/dtos/register_request_dto.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RegisterRequestDto.toMap', () {
    test('serializes full payload expected by register endpoint', () {
      final dto = RegisterRequestDto(
        name: 'Maria Silva',
        email: 'user@example.com',
        phone: '+5511999999999',
        password: 'P@ssword123',
        birthDate: DateTime(1990, 1, 15),
        cpf: '123.456.789-01',
        emailVerificationToken: 'token-confirmado-email',
        phoneVerificationToken: 'token-confirmado-phone',
      );

      final map = dto.toMap();

      expect(map['name'], 'Maria Silva');
      expect(map['email'], 'user@example.com');
      expect(map['phone'], '+5511999999999');
      expect(map['password'], 'P@ssword123');
      expect(map['birth_date'], '1990-01-15');
      expect(map['cpf'], '12345678901');
      expect(map['email_verification_token'], 'token-confirmado-email');
      expect(map['phone_verification_token'], 'token-confirmado-phone');
    });

    test('keeps cpf normalized with 11 digits in payload', () {
      final dto = RegisterRequestDto(
        name: 'Maria Silva',
        email: 'user@example.com',
        phone: '+5511999999999',
        password: 'P@ssword123',
        birthDate: DateTime(1990, 1, 15),
        cpf: '123.456.789-01',
        emailVerificationToken: 'token-confirmado-email',
        phoneVerificationToken: 'token-confirmado-phone',
      );

      final map = dto.toMap();

      expect(map['cpf'], '12345678901');
    });

    test('always includes required verification tokens', () {
      final dto = RegisterRequestDto(
        name: 'Maria Silva',
        email: 'user@example.com',
        phone: '+5511999999999',
        password: 'P@ssword123',
        birthDate: DateTime(1990, 1, 15),
        cpf: '12345678901',
        emailVerificationToken: 'token-confirmado-email',
        phoneVerificationToken: 'token-confirmado-phone',
      );

      final map = dto.toMap();

      expect(map['email_verification_token'], 'token-confirmado-email');
      expect(map['phone_verification_token'], 'token-confirmado-phone');
    });
  });

  group('RegisterRequestDto.fromMap', () {
    test('parses new register contract fields', () {
      final dto = RegisterRequestDto.fromMap({
        'name': 'Maria Silva',
        'email': 'user@example.com',
        'phone': '+5511999999999',
        'password': 'P@ssword123',
        'birth_date': '1990-01-15',
        'cpf': '12345678901',
        'email_verification_token': 'token-confirmado-email',
        'phone_verification_token': 'token-confirmado-phone',
      });

      expect(dto.name, 'Maria Silva');
      expect(dto.email, 'user@example.com');
      expect(dto.phone, '+5511999999999');
      expect(dto.password, 'P@ssword123');
      expect(dto.birthDate, DateTime.parse('1990-01-15'));
      expect(dto.cpf, '12345678901');
      expect(dto.emailVerificationToken, 'token-confirmado-email');
      expect(dto.phoneVerificationToken, 'token-confirmado-phone');
    });
  });
}
