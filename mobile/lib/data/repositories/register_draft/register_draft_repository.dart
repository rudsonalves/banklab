import '/core/result/result.dart';
import '/data/services/auth/cache/register_draft/register_draft_load_result.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';

abstract class RegisterDraftRepository {
  RegisterDraftSnapshot? get snapshot;

  AsyncResult<RegisterDraftLoadResult> getByCPF(String cpf);

  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot);

  AsyncResult<Unit> deleteByCPF(String cpf);
}
