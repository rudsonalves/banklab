import '/core/extensions/datetime_extension.dart';
import '/core/extensions/string.dart';

class RegisterRequestDto {
  final String name;
  final String email;
  final String phone;
  final String password;
  final DateTime birthDate;
  final String cpf;
  final String emailVerificationToken;
  final String phoneVerificationToken;

  RegisterRequestDto({
    required this.name,
    required this.email,
    required this.phone,
    required this.password,
    required this.birthDate,
    required this.cpf,
    required this.emailVerificationToken,
    required this.phoneVerificationToken,
  });

  Map<String, dynamic> toMap() {
    return {
      'name': name,
      'email': email,
      'phone': phone,
      'password': password,
      'birth_date': birthDate.toIso8601String().split('T').first,
      'cpf': cpf.onlyNumbers,
      'email_verification_token': emailVerificationToken,
      'phone_verification_token': phoneVerificationToken,
    };
  }

  factory RegisterRequestDto.fromMap(Map<String, dynamic> map) {
    final birthDate = DateParser.parseOrNull(map['birth_date'] as String?);
    if (birthDate == null) {
      throw const FormatException('birth_date is required');
    }

    return RegisterRequestDto(
      name: map['name'] as String,
      email: map['email'] as String,
      phone: map['phone'] as String,
      password: map['password'] as String,
      birthDate: birthDate,
      cpf: map['cpf'] as String,
      emailVerificationToken: map['email_verification_token'] as String,
      phoneVerificationToken: map['phone_verification_token'] as String,
    );
  }
}
