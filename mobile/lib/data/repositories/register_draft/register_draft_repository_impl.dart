import '/core/result/result.dart';
import '/data/services/cache/register_draft/register_draft_store.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';
import 'register_draft_repository.dart';

class RegisterDraftRepositoryImpl implements RegisterDraftRepository {
  static const Duration _defaultTTL = Duration(hours: 24);

  final RegisterDraftStore _store;
  final Duration _ttl;
  final DateTime Function() _now;

  RegisterDraftRepositoryImpl(
    this._store, {
    Duration ttl = _defaultTTL,
    DateTime Function()? now,
  }) : _ttl = ttl,
       _now = now ?? (() => DateTime.now().toUtc());

  RegisterDraftSnapshot? _snapshot;

  @override
  RegisterDraftSnapshot? get snapshot => _snapshot;

  @override
  AsyncResult<RegisterDraftSnapshot> getByCPF(String cpf) async {
    final result = await _store.getByCPF(cpf);

    if (result.isFailure) {
      if (result.error?.code == AppErrorCode.storageNotFound) {
        return await _createASnapshotForCPF(cpf);
      }

      _snapshot = null;
      return Failure(result.error!);
    }

    final snapshot = result.value!;
    if (_isOld(snapshot)) {
      return await _createASnapshotForCPF(cpf);
    }

    _snapshot = snapshot;
    return Success(snapshot);
  }

  AsyncResult<RegisterDraftSnapshot> _createASnapshotForCPF(String cpf) async {
    final snapshot = RegisterDraftSnapshot.empty(cpf);
    final result = await save(snapshot);
    if (result.isSuccess) _snapshot = snapshot;

    if (result.isSuccess) return Success(snapshot);

    return Failure(result.error!);
  }

  bool _isOld(RegisterDraftSnapshot snapshot) {
    final expiresAt = snapshot.updatedAt.toUtc().add(_ttl);
    final now = _now().toUtc();
    return !now.isBefore(expiresAt);
  }

  @override
  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot) async {
    final result = await _store.save(snapshot);

    if (result.isFailure) {
      _snapshot = null;
      return Failure(result.error!);
    }

    _snapshot = snapshot;
    return result;
  }

  @override
  AsyncResult<Unit> deleteByCPF(String cpf) async {
    final result = await _store.deleteByCPF(cpf);

    if (result.isFailure) {
      return Failure(result.error!);
    }

    _snapshot = null;
    return result;
  }
}
