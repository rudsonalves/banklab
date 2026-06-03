import '../../../user/enums/user_role.dart';

class UserSession {
  final String userId;
  final String email;
  final String? phone;
  final UserRole role;

  UserSession({
    required this.userId,
    required this.email,
    this.phone,
    required this.role,
  });

  factory UserSession.fromApi(Map<String, dynamic> map) {
    return UserSession(
      userId: map['id'] as String,
      email: map['email'] as String,
      phone: map['phone'] as String?,
      role: UserRole.byName(map['role'] as String),
    );
  }
}
