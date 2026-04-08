import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/secure_storage/flutter_secure_storage_local_storage.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('FlutterSecureStorageLocalStorage', () {
    test('keysWithPrefix should filter matching keys', () async {
      final storage = _TestFlutterSecureStorage(
        onReadAll: () async => {
          'auth_token': 'abc',
          'auth_refresh': 'xyz',
          'profile_name': 'Ada',
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final keys = await localStorage.keysWithPrefix('auth_');

      expect(keys, ['auth_token', 'auth_refresh']);
    });

    test('keysWithPrefix should return empty list on error', () async {
      final storage = _TestFlutterSecureStorage(
        onReadAll: () async => throw Exception('readAll failure'),
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final keys = await localStorage.keysWithPrefix('auth_');

      expect(keys, isEmpty);
    });

    test('write should return Success on write completion', () async {
      String? capturedKey;
      String? capturedValue;
      final storage = _TestFlutterSecureStorage(
        onWrite: ({required key, required value}) async {
          capturedKey = key;
          capturedValue = value;
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.write('token', 'abc');

      expect(result, isA<Success<Unit>>());
      expect(capturedKey, 'token');
      expect(capturedValue, 'abc');
    });

    test('write should map errors to storage AppError', () async {
      final storage = _TestFlutterSecureStorage(
        onWrite: ({required key, required value}) async {
          throw StateError('disk full');
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.write('token', 'abc');

      expect(result, isA<Failure<Unit>>());
      final error = result.error;
      expect(error?.code, AppErrorCode.storageError);
      expect(error?.message, 'Failed to write key: token');
      expect(error?.details, isA<StateError>());
    });

    test('read should return Success when value exists', () async {
      final storage = _TestFlutterSecureStorage(
        onRead: ({required key}) async => key == 'token' ? 'abc' : null,
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.read('token');

      expect(result, isA<Success<String>>());
      expect(result.value, 'abc');
    });

    test('read should return notFound AppError when key is absent', () async {
      final storage = _TestFlutterSecureStorage(
        onRead: ({required key}) async => null,
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.read('missing_key');

      expect(result, isA<Failure<String>>());
      final error = result.error;
      expect(error?.code, AppErrorCode.storageNotFound);
      expect(error?.message, 'Key not found: missing_key');
      expect(
        error?.details,
        'Key "missing_key" was not found in storage.',
      );
    });

    test('delete should return Success on delete completion', () async {
      String? capturedKey;
      final storage = _TestFlutterSecureStorage(
        onDelete: ({required key}) async {
          capturedKey = key;
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.delete('token');

      expect(result, isA<Success<Unit>>());
      expect(capturedKey, 'token');
    });

    test('delete should map errors to storage AppError', () async {
      final storage = _TestFlutterSecureStorage(
        onDelete: ({required key}) async {
          throw Exception('delete failure');
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.delete('token');

      expect(result, isA<Failure<Unit>>());
      final error = result.error;
      expect(error?.code, AppErrorCode.storageError);
      expect(error?.message, 'Failed to delete key: token');
    });

    test('deleteAll should return Success on completion', () async {
      var called = false;
      final storage = _TestFlutterSecureStorage(
        onDeleteAll: () async {
          called = true;
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.deleteAll();

      expect(result, isA<Success<Unit>>());
      expect(called, isTrue);
    });

    test('deleteAll should map errors to storage AppError', () async {
      final storage = _TestFlutterSecureStorage(
        onDeleteAll: () async {
          throw Exception('deleteAll failure');
        },
      );
      final localStorage = FlutterSecureStorageLocalStorage(storage: storage);

      final result = await localStorage.deleteAll();

      expect(result, isA<Failure<Unit>>());
      final error = result.error;
      expect(error?.code, AppErrorCode.storageError);
      expect(error?.message, 'Failed to delete all keys');
    });
  });
}

typedef _ReadAll = Future<Map<String, String>> Function();
typedef _Read = Future<String?> Function({required String key});
typedef _Write =
    Future<void> Function({
      required String key,
      required String value,
    });
typedef _Delete = Future<void> Function({required String key});
typedef _DeleteAll = Future<void> Function();

class _TestFlutterSecureStorage extends FlutterSecureStorage {
  _TestFlutterSecureStorage({
    _ReadAll? onReadAll,
    _Read? onRead,
    _Write? onWrite,
    _Delete? onDelete,
    _DeleteAll? onDeleteAll,
  }) : _onReadAll = onReadAll,
       _onRead = onRead,
       _onWrite = onWrite,
       _onDelete = onDelete,
       _onDeleteAll = onDeleteAll;

  final _ReadAll? _onReadAll;
  final _Read? _onRead;
  final _Write? _onWrite;
  final _Delete? _onDelete;
  final _DeleteAll? _onDeleteAll;

  @override
  Future<Map<String, String>> readAll({
    AppleOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    AppleOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    if (_onReadAll != null) {
      return _onReadAll();
    }
    return {};
  }

  @override
  Future<String?> read({
    required String key,
    AppleOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    AppleOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    if (_onRead != null) {
      return _onRead(key: key);
    }
    return null;
  }

  @override
  Future<void> write({
    required String key,
    required String? value,
    AppleOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    AppleOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    if (_onWrite != null) {
      await _onWrite(key: key, value: value ?? '');
    }
  }

  @override
  Future<void> delete({
    required String key,
    AppleOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    AppleOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    if (_onDelete != null) {
      await _onDelete(key: key);
    }
  }

  @override
  Future<void> deleteAll({
    AppleOptions? iOptions,
    AndroidOptions? aOptions,
    LinuxOptions? lOptions,
    WebOptions? webOptions,
    AppleOptions? mOptions,
    WindowsOptions? wOptions,
  }) async {
    if (_onDeleteAll != null) {
      await _onDeleteAll();
    }
  }
}
