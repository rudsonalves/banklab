import 'dart:convert';

import 'package:crypto/crypto.dart';

import '/core/extensions/string.dart';
import '/core/result/result.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';

class RegisterDraftStore {
  static const keyPrefix = 'onboarding_draft:';

  final LocalSecureStorage _secureStorage;

  const RegisterDraftStore(this._secureStorage);

  /// Generates a unique storage key for a given CPF by normalizing it
  /// and applying a SHA-256 hash.
  String keyForCPF(String cpf) {
    final normalizedCPF = cpf.onlyNumbers;
    final digest = sha256.convert(utf8.encode(normalizedCPF)).toString();
    return '$keyPrefix$digest';
  }

  /// Saves a [RegisterDraftSnapshot] to secure storage using a key derived
  /// from the CPF.
  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot) async {
    final key = keyForCPF(snapshot.cpf);
    final payload = jsonEncode(snapshot.toMap());
    return _secureStorage.write(key, payload);
  }

  /// Retrieves a [RegisterDraftSnapshot] from secure storage using a key derived
  /// from the CPF. If the snapshot is not found or cannot be decoded, it returns
  /// a storage not found failure.
  AsyncResult<RegisterDraftSnapshot> getByCPF(String cpf) async {
    final key = keyForCPF(cpf);
    final result = await _secureStorage.read(key);

    if (result.isFailure) return Failure(result.error!);

    final snapshot = _decodeSnapshot(result.value);
    if (snapshot == null) {
      final deleteResult = await _secureStorage.delete(key);
      if (deleteResult.isFailure) return Failure(deleteResult.error!);

      return const Failure(
        AppError(
          code: AppErrorCode.storageNotFound,
          message: 'Register draft not found',
        ),
      );
    }

    return Success(snapshot);
  }

  /// Deletes a [RegisterDraftSnapshot] from secure storage using a key derived
  /// from the CPF. Returns a [Unit] result indicating success or failure of the
  /// operation.
  AsyncResult<Unit> deleteByCPF(String cpf) {
    return _secureStorage.delete(keyForCPF(cpf));
  }

  /// Decodes a JSON string into a [RegisterDraftSnapshot]. If the input is
  /// null, empty, or cannot be decoded, it returns null.
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
