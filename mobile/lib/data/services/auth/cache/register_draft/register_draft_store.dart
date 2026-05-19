import 'dart:convert';

import 'package:crypto/crypto.dart';

import '/core/extensions/string.dart';
import '/core/result/result.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';
import 'register_draft_load_result.dart';

export 'register_draft_load_result.dart';

class RegisterDraftStore {
  static const keyPrefix = 'onboarding_draft:';

  final LocalSecureStorage _secureStorage;

  const RegisterDraftStore(this._secureStorage);

  String keyForCPF(String cpf) {
    final normalizedCPF = cpf.onlyNumbers;
    final digest = sha256.convert(utf8.encode(normalizedCPF)).toString();
    return '$keyPrefix$digest';
  }

  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot) async {
    final key = keyForCPF(snapshot.cpf);
    final payload = jsonEncode(snapshot.toMap());
    return _secureStorage.write(key, payload);
  }

  AsyncResult<RegisterDraftLoadResult> getByCPF(String cpf) async {
    final key = keyForCPF(cpf);
    final result = await _secureStorage.read(key);

    if (result.isFailure) {
      if (result.error?.code == AppErrorCode.storageNotFound) {
        return const Success(RegisterDraftNotFound());
      }

      return Failure(result.error!);
    }

    final snapshot = _decodeSnapshot(result.value);
    if (snapshot == null) {
      final deleteResult = await _secureStorage.delete(key);
      if (deleteResult.isFailure) return Failure(deleteResult.error!);

      return const Success(RegisterDraftNotFound());
    }

    return Success(RegisterDraftFound(snapshot));
  }

  AsyncResult<Unit> deleteByCPF(String cpf) {
    return _secureStorage.delete(keyForCPF(cpf));
  }

  RegisterDraftSnapshot? _decodeSnapshot(String? payload) {
    if (payload == null || payload.trim().isEmpty) return null;

    try {
      final decoded = jsonDecode(payload);
      if (decoded is! Map<String, dynamic>) return null;
      return RegisterDraftSnapshot.fromMapOrNull(decoded);
    } catch (_) {
      return null;
    }
  }
}
