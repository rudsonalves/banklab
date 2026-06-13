import 'customer_session.dart';
import 'readiness_session.dart';
import 'user_session.dart';

export '../../../../../data/services/apis/transaction_password/enums/transaction_password_status.dart';
export 'customer_session.dart';
export 'readiness_session.dart';
export 'user_session.dart';

class AuthSession {
  final UserSession user;
  final CustomerSession? customer;
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
          ? CustomerSession.fromApi(customerMap)
          : null,
      readiness: ReadinessSession.fromApi(readinessMap),
    );
  }

  AuthSession copyWith({
    UserSession? user,
    CustomerSession? customer,
    ReadinessSession? readiness,
  }) {
    return AuthSession(
      user: user ?? this.user,
      customer: customer ?? this.customer,
      readiness: readiness ?? this.readiness,
    );
  }
}
