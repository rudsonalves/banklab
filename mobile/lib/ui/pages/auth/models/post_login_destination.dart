import '/domain/common/auth/models/auth_session/auth_session.dart';

enum PostLoginDestination {
  home,
  transactionPassword,
  blocked,
  sessionError,
}

class PostLoginDestinationResolver {
  const PostLoginDestinationResolver._();

  static PostLoginDestination resolve(AuthSession? session) {
    if (session == null) return PostLoginDestination.sessionError;

    final readiness = session.readiness;

    switch (readiness.transactionPasswordStatus) {
      case TransactionPasswordStatus.active:
        return readiness.canAccessHome
            ? PostLoginDestination.home
            : PostLoginDestination.blocked;
      case TransactionPasswordStatus.notSet:
        return PostLoginDestination.transactionPassword;
      case TransactionPasswordStatus.locked:
      case TransactionPasswordStatus.unknown:
        return PostLoginDestination.blocked;
    }
  }
}
