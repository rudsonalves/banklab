import '/domain/enums/user_role.dart';

class RegisterResponseDto {
  final String id;
  final String email;
  final UserRole role;
  final String customerId;

  RegisterResponseDto({
    required this.id,
    required this.email,
    required this.role,
    required this.customerId,
  });

  factory RegisterResponseDto.fromMap(Map<String, dynamic> map) {
    return RegisterResponseDto(
      id: map['id'] as String,
      email: map['email'] as String,
      role: UserRole.byName(map['role'] as String),
      customerId: map['customer_id'] as String,
    );
  }
}
