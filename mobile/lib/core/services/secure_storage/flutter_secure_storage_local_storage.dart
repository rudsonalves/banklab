import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '/core/result/result.dart';
import '/core/services/logging/console_log.dart';
import '../../extensions/storage_app_error.dart';
import 'local_secure_storage.dart';

class FlutterSecureStorageLocalStorage implements LocalSecureStorage {
  final FlutterSecureStorage storage;

  FlutterSecureStorageLocalStorage({required this.storage});

  final _log = ConsoleLog('FlutterSecureStorageLocalStorage');

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    try {
      final allKeys = await storage.readAll();
      final filteredKeys = allKeys.keys
          .where((key) => key.startsWith(pattern))
          .toList();
      return filteredKeys;
    } catch (err, stack) {
      _log.error('[readKeys]: $err', error: err, stack: stack);
      return [];
    }
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    try {
      await storage.delete(key: key);
      return Success(unit);
    } catch (err, stack) {
      _log.error('[delete]: $err', error: err, stack: stack);
      return Failure(
        StorageAppError.storage(
          message: 'Failed to delete key: $key',
          details: err,
        ),
      );
    }
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    try {
      await storage.deleteAll();
      return Success(unit);
    } catch (err, stack) {
      _log.error('[deleteAll]: $err', error: err, stack: stack);
      return Failure(
        StorageAppError.storage(
          message: 'Failed to delete all keys',
          details: err,
        ),
      );
    }
  }

  @override
  AsyncResult<String> read(String key) async {
    try {
      final value = await storage.read(key: key);
      if (value == null) {
        return Failure(StorageAppError.notFound(key));
      }
      return Success(value);
    } catch (err, stack) {
      _log.error('[read]: $err', error: err, stack: stack);
      return Failure(
        StorageAppError.storage(
          message: 'Failed to read key: $key',
          details: err,
        ),
      );
    }
  }

  @override
  AsyncResult<Unit> write(String key, String value) async {
    try {
      await storage.write(key: key, value: value);
      return Success(unit);
    } catch (err, stack) {
      _log.error('[write]: $err', error: err, stack: stack);
      return Failure(
        StorageAppError.storage(
          message: 'Failed to write key: $key',
          details: err,
        ),
      );
    }
  }
}
