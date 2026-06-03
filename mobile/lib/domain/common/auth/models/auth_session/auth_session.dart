import 'customer_session.dart';
import 'readiness_session.dart';
import 'user_session.dart';

export 'customer_session.dart';
export 'enums/transaction_password_status.dart';
export 'readiness_session.dart';
export 'user_session.dart';

class AuthSession {
  final UserSession user;
  final CustommerSession? customer;
  final ReadinessSession readiness;

  AuthSession({
    required this.user,
    required this.customer,
    required this.readiness,
  });

  factory AuthSession.fromApi(Map<String, dynamic> map) {
    final userMap = map['user'] as Map<String, dynamic>?;
    final customerMap = map['customer'] as Map<String, dynamic>?;
    final readinessMap = map['readiness'] as Map<String, dynamic>?;

    if (userMap == null || readinessMap == null) {
      throw ArgumentError('Invalid API response: missing required fields');
    }

    return AuthSession(
      user: UserSession.fromApi(userMap),
      customer: customerMap != null
          ? CustommerSession.fromApi(customerMap)
          : null,
      readiness: ReadinessSession.fromApi(readinessMap),
    );
  }
}
