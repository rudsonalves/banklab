import '/core/result/result.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';

abstract class RegisterDraftRepository {
  RegisterDraftSnapshot? get snapshot;

  AsyncResult<RegisterDraftSnapshot> getByCPF(String cpf);

  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot);

  AsyncResult<Unit> deleteByCPF(String cpf);
}
