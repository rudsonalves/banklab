import '/domain/common/auth/models/auth_session/auth_session.dart';

class AppSection {
  AuthSession? _current;

  AuthSession? get currentSession => _current;

  UserSession? get user => _current?.user;

  CustomerSession? get customer => _current?.customer;

  ReadinessSession? get readiness => _current?.readiness;

  bool get canAccessHome =>
      _current != null && _current!.readiness.canAccessHome;

  bool get isNotNull => _current != null;

  void setAuthSession(AuthSession session) => _current = session;

  void clear() => _current = null;

  void markTransactionPasswordAsActive() {
    final session = _current;
    if (session == null) return;

    final readiness = session.readiness.copyWith(
      transactionPasswordStatus: TransactionPasswordStatus.active,
    );
    _current = session.copyWith(
      readiness: readiness,
    );
  }
}
