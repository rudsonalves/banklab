import '/core/result/result.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';
import '../../services/cache/last_login/register_draft/register_draft_store.dart';
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
  AsyncResult<RegisterDraftLoadResult> getByCPF(String cpf) async {
    final result = await _store.getByCPF(cpf);

    if (result.isFailure) {
      _snapshot = null;
      return Failure(result.error!);
    }

    final loadResult = result.value!;
    if (loadResult.isNotFound) {
      _snapshot = null;
      return result;
    }

    _snapshot = (loadResult as RegisterDraftFound).snapshot;

    final expiresAt = _snapshot!.updatedAt.toUtc().add(_ttl);
    final now = _now().toUtc();
    if (!now.isBefore(expiresAt)) {
      final deleteResult = await deleteByCPF(cpf);
      if (deleteResult.isFailure) {
        _snapshot = null;
        return Failure(deleteResult.error!);
      }

      _snapshot = null;
      return const Success(RegisterDraftNotFound());
    }

    return result;
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
