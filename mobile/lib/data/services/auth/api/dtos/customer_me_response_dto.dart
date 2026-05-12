import '/core/extensions/datetime_extension.dart';

class CustomerMeResponseDto {
  final String id;
  final String name;
  final String cpf;
  final String email;
  final DateTime createdAt;

  CustomerMeResponseDto({
    required this.id,
    required this.name,
    required this.cpf,
    required this.email,
    required this.createdAt,
  });

  factory CustomerMeResponseDto.fromMap(Map<String, dynamic> map) {
    return CustomerMeResponseDto(
      id: map['id'] as String,
      name: map['name'] as String,
      cpf: map['cpf'] as String,
      email: map['email'] as String,
      createdAt:
          DateTimeExtensions.parseOrNull(map['createdAt']) ?? DateTime.now(),
    );
  }
}
