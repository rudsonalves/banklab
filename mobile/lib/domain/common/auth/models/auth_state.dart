import '/domain/common/user/enums/user_role.dart';

sealed class AuthState {
  final String? userId;
  final String email;
  final UserRole role;

  AuthState({
    this.userId,
    required this.email,
    this.role = UserRole.none,
  });

  factory AuthState.fromLoginMap(Map<String, dynamic> map) {
    if (_hasOperationalTokens(map)) {
      return OperationalAuthState.fromMap(map);
    }

    if (_hasRestrictedToken(map)) {
      return RestrictedInstallationAuthState.fromMap(map);
    }

    throw FormatException('Unknown login response shape');
  }

  static bool _hasOperationalTokens(Map<String, dynamic> map) {
    return map['access_token'] is String && map['refresh_token'] is String;
  }

  static bool _hasRestrictedToken(Map<String, dynamic> map) {
    return map['restricted_access_token'] is String;
  }
}

class OperationalAuthState extends AuthState {
  final String accessToken;
  final String refreshToken;
  final String customerId;

  OperationalAuthState({
    required this.accessToken,
    required this.refreshToken,
    required String super.userId,
    required super.email,
    required super.role,
    required this.customerId,
  });

  factory OperationalAuthState.fromMap(Map<String, dynamic> map) {
    return OperationalAuthState(
      accessToken: map['access_token'] as String,
      refreshToken: map['refresh_token'] as String,
      userId: map['user_id'] as String,
      email: map['email'] as String,
      role: UserRole.byName(map['role'] as String),
      customerId: map['customer_id'] as String,
    );
  }
}

class AnonymousAuthState extends AuthState {
  AnonymousAuthState() : super(email: '', role: UserRole.none);
}

class RestrictedInstallationAuthState extends AuthState {
  final String restrictedAccessToken;
  final String restrictedTokenType;
  final String restrictedScope;
  final DateTime restrictedExpiresAt;
  final String customerId;

  RestrictedInstallationAuthState({
    required this.restrictedAccessToken,
    required this.restrictedTokenType,
    required this.restrictedScope,
    required this.restrictedExpiresAt,
    required String super.userId,
    required super.email,
    required super.role,
    required this.customerId,
  });

  factory RestrictedInstallationAuthState.fromMap(Map<String, dynamic> map) {
    return RestrictedInstallationAuthState(
      restrictedAccessToken: map['restricted_access_token'] as String,
      restrictedTokenType: map['restricted_token_type'] as String,
      restrictedScope: map['restricted_scope'] as String,
      restrictedExpiresAt: DateTime.parse(
        map['restricted_expires_at'] as String,
      ),
      userId: map['user_id'] as String,
      email: map['email'] as String,
      role: UserRole.byName(map['role'] as String),
      customerId: map['customer_id'] as String,
    );
  }
}
