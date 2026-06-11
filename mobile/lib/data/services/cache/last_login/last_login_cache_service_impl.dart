import '/core/resources/storage_keys.dart';
import '/core/result/result.dart';
import '/core/extensions/string.dart';
import '/core/services/logging/console_log.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import 'last_login_cache_service.dart';
import 'models/last_login_identity.dart';

class LastLoginCacheServiceImpl implements LastLoginCacheService {
  final LocalSecureStorage _secureStorage;

  LastLoginCacheServiceImpl(this._secureStorage);

  final _log = ConsoleLog('LastLoginCacheServiceImpl');

  @override
  AsyncResult<LastLoginIdentity> get() async {
    final nameResult = await _secureStorage.read(StorageKeys.lastLoginName);
    final identifierResult = await _secureStorage.read(
      StorageKeys.lastLoginIdentifier,
    );

    if (nameResult.isFailure || identifierResult.isFailure) {
      _log.error('Failed to retrieve last login identity');
      return Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Failed to retrieve last login identity',
        ),
      );
    }

    final name = nameResult.value?.trim();
    final identifier = identifierResult.value?.trim();

    if (name == null ||
        name.isEmpty ||
        identifier == null ||
        identifier.isEmpty ||
        !identifier.isValidEmail) {
      _log.info('No valid last login identity found');
      return Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Last login identity data is incomplete',
        ),
      );
    }

    return Success(LastLoginIdentity(name: name, identifier: identifier));
  }

  @override
  AsyncResult<Unit> save(LastLoginIdentity identity) async {
    if (identity.name.trim().isEmpty || identity.identifier.trim().isEmpty) {
      _log.error('Name and identifier cannot be empty');
      return Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Name and identifier cannot be empty',
        ),
      );
    }

    final nameResult = await _secureStorage.write(
      StorageKeys.lastLoginName,
      identity.name.trim(),
    );
    final identifierResult = await _secureStorage.write(
      StorageKeys.lastLoginIdentifier,
      identity.identifier.trim(),
    );

    if (nameResult.isFailure || identifierResult.isFailure) {
      _log.error('Failed to save last login identity');
      return Failure(
        AppError(
          code: AppErrorCode.storageError,
          message: 'Failed to save last login identity',
        ),
      );
    }

    return const Success(unit);
  }

  @override
  AsyncResult<Unit> clear() async {
    final nameResult = await _secureStorage.delete(StorageKeys.lastLoginName);
    final identifierResult = await _secureStorage.delete(
      StorageKeys.lastLoginIdentifier,
    );

    if (nameResult.isFailure || identifierResult.isFailure) {
      _log.error('Failed to clear last login identity');
      return Failure(
        AppError(
          code: AppErrorCode.storageError,
          message: 'Failed to clear last login identity',
        ),
      );
    }

    return const Success(unit);
  }
}
