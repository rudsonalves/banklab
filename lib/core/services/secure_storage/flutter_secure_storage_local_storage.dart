import 'dart:developer';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '/core/result/result.dart';
import '../../extensions/storage_app_error.dart';
import 'local_secure_storage.dart';

class FlutterSecureStorageLocalStorage implements LocalSecureStorage {
  final FlutterSecureStorage storage;

  const FlutterSecureStorageLocalStorage({required this.storage});

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    try {
      final allKeys = await storage.readAll();
      final filteredKeys = allKeys.keys
          .where((key) => key.startsWith(pattern))
          .toList();
      return filteredKeys;
    } catch (err) {
      log('[readKeys]: $err');
      return [];
    }
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    try {
      await storage.delete(key: key);
      return Success(unit);
    } catch (err) {
      log('[delete]: $err');
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
    } catch (err) {
      log('[deleteAll]: $err');
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
    } catch (err) {
      log('[read]: $err');
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
    } catch (err) {
      log('[write]: $err');
      return Failure(
        StorageAppError.storage(
          message: 'Failed to write key: $key',
          details: err,
        ),
      );
    }
  }
}
