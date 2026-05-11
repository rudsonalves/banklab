import '/domain/common/user/enums/user_role.dart';

class AuthMeResponseDto {
  final String id;
  final UserRole role;
  final String email;
  final String customerId;

  AuthMeResponseDto({
    required this.id,
    required this.role,
    required this.email,
    required this.customerId,
  });

  factory AuthMeResponseDto.fromMap(Map<String, dynamic> map) {
    return AuthMeResponseDto(
      id: map['id'] as String,
      role: UserRole.byName(map['role'] as String),
      email: map['email'] as String,
      customerId: map['customerId'] as String,
    );
  }
}
