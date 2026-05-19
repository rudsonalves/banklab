import '/domain/common/auth/models/register_draft_snapshot.dart';

sealed class RegisterDraftLoadResult {
  const RegisterDraftLoadResult();

  bool get isFound => this is RegisterDraftFound;
  bool get isNotFound => this is RegisterDraftNotFound;
}

final class RegisterDraftFound extends RegisterDraftLoadResult {
  final RegisterDraftSnapshot snapshot;

  const RegisterDraftFound(this.snapshot);
}

final class RegisterDraftNotFound extends RegisterDraftLoadResult {
  const RegisterDraftNotFound();
}
