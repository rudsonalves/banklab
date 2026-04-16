import '/domain/enums/user_role.dart';

sealed class AuthUser {
  final String? userId;
  final String email;
  final UserRole role;

  AuthUser({
    this.userId,
    required this.email,
    this.role = UserRole.none,
  });
}

class LoggedUser extends AuthUser {
  final String accessToken;
  final String refreshToken;
  final String customerId;

  LoggedUser({
    required this.accessToken,
    required this.refreshToken,
    required String super.userId,
    required super.email,
    required super.role,
    required this.customerId,
  });

  factory LoggedUser.fromMap(Map<String, dynamic> map) {
    return LoggedUser(
      accessToken: map['access_token'] as String,
      refreshToken: map['refresh_token'] as String,
      userId: map['user_id'] as String,
      email: map['email'] as String,
      role: UserRole.byName(map['role'] as String),
      customerId: map['customer_id'] as String,
    );
  }
}

class NotLoggedUser extends AuthUser {
  NotLoggedUser() : super(email: '', role: UserRole.none);
}
