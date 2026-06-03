import 'enums/transaction_password_status.dart';

class ReadinessSession {
  final bool onboardingCompleted;
  final bool approved;
  final bool hasOperationalAccount;
  final TransactionPasswordStatus transactionPasswordStatus;
  final bool canAccessHome;

  ReadinessSession({
    required this.onboardingCompleted,
    required this.approved,
    required this.hasOperationalAccount,
    required this.transactionPasswordStatus,
    required this.canAccessHome,
  });

  factory ReadinessSession.fromApi(Map<String, dynamic> map) {
    return ReadinessSession(
      onboardingCompleted: map['onboarding_completed'] as bool,
      approved: map['approved'] as bool,
      hasOperationalAccount: map['has_operational_account'] as bool,
      transactionPasswordStatus: TransactionPasswordStatus.byName(
        map['transaction_password_status'] as String,
      ),
      canAccessHome: map['can_access_home'] as bool,
    );
  }
}
