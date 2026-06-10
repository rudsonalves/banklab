import 'enums/transaction_password_status.dart';

class ReadinessSession {
  final bool onboardingCompleted;
  final bool approved;
  final bool hasOperationalAccount;
  final TransactionPasswordStatus transactionPasswordStatus;

  bool get hasActiveTransactionPassword =>
      transactionPasswordStatus == TransactionPasswordStatus.active;

  bool get canAccessHome =>
      onboardingCompleted &&
      approved &&
      hasOperationalAccount &&
      hasActiveTransactionPassword;

  ReadinessSession({
    required this.onboardingCompleted,
    required this.approved,
    required this.hasOperationalAccount,
    required this.transactionPasswordStatus,
  });

  factory ReadinessSession.fromApi(Map<String, dynamic> map) {
    return ReadinessSession(
      onboardingCompleted: map['onboarding_completed'] as bool,
      approved: map['approved'] as bool,
      hasOperationalAccount: map['has_operational_account'] as bool,
      transactionPasswordStatus: TransactionPasswordStatus.byName(
        map['transaction_password_status'] as String,
      ),
    );
  }

  ReadinessSession copyWith({
    bool? onboardingCompleted,
    bool? approved,
    bool? hasOperationalAccount,
    TransactionPasswordStatus? transactionPasswordStatus,
  }) {
    return ReadinessSession(
      onboardingCompleted: onboardingCompleted ?? this.onboardingCompleted,
      approved: approved ?? this.approved,
      hasOperationalAccount:
          hasOperationalAccount ?? this.hasOperationalAccount,
      transactionPasswordStatus:
          transactionPasswordStatus ?? this.transactionPasswordStatus,
    );
  }
}
