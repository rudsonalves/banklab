import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/installation_identity/installation_identity.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('InstallationIdentityService', () {
    test('first run creates, persists, and marks a UUID v4', () async {
      final storage = _MemorySecureStorage();
      final markerStore = _FakeInstallationMarkerStore(markerPresent: false);
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Success<String>>());
      expect(result.value, _isCanonicalUuidV4);
      expect(storage.values[StorageKeys.installationId], result.value);
      expect(markerStore.markCalls, 1);
      expect(markerStore.markerPresent, isTrue);
    });

    test('reuses stored UUID v4 when local marker is present', () async {
      const storedId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';
      final storage = _MemorySecureStorage(
        initialValues: {StorageKeys.installationId: storedId},
      );
      final markerStore = _FakeInstallationMarkerStore(markerPresent: true);
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Success<String>>());
      expect(result.value, storedId);
      expect(storage.writeCalls, 0);
      expect(markerStore.markCalls, 1);
    });

    test('replaces stored value when local marker is missing', () async {
      const oldId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';
      final storage = _MemorySecureStorage(
        initialValues: {StorageKeys.installationId: oldId},
      );
      final markerStore = _FakeInstallationMarkerStore(markerPresent: false);
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Success<String>>());
      expect(result.value, _isCanonicalUuidV4);
      expect(result.value, isNot(oldId));
      expect(storage.values[StorageKeys.installationId], result.value);
      expect(markerStore.markCalls, 1);
    });

    test('replaces invalid stored value', () async {
      final storage = _MemorySecureStorage(
        initialValues: {StorageKeys.installationId: 'not-a-uuid'},
      );
      final markerStore = _FakeInstallationMarkerStore(markerPresent: true);
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Success<String>>());
      expect(result.value, _isCanonicalUuidV4);
      expect(storage.values[StorageKeys.installationId], result.value);
      expect(storage.writeCalls, 1);
    });

    test('returns failure when marker lookup fails', () async {
      final markerStore = _FakeInstallationMarkerStore(
        markerPresent: false,
        hasMarkerResult: const Failure(
          AppError(
            code: AppErrorCode.storageError,
            message: 'marker read failed',
          ),
        ),
      );
      final service = InstallationIdentityService(
        secureStorage: _MemorySecureStorage(),
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Failure<String>>());
      expect(result.error?.message, 'marker read failed');
    });

    test('returns failure when secure storage read fails', () async {
      final storage = _MemorySecureStorage(
        readResult: const Failure(
          AppError(
            code: AppErrorCode.storageError,
            message: 'secure read failed',
          ),
        ),
      );
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: _FakeInstallationMarkerStore(markerPresent: true),
      );

      final result = await service.resolve();

      expect(result, isA<Failure<String>>());
      expect(result.error?.message, 'secure read failed');
      expect(storage.writeCalls, 0);
    });

    test('returns failure when secure storage write fails', () async {
      final storage = _MemorySecureStorage(
        writeResult: const Failure(
          AppError(
            code: AppErrorCode.storageError,
            message: 'secure write failed',
          ),
        ),
      );
      final markerStore = _FakeInstallationMarkerStore(markerPresent: false);
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Failure<String>>());
      expect(result.error?.message, 'secure write failed');
      expect(markerStore.markCalls, 0);
    });

    test('returns failure when marker refresh fails', () async {
      const storedId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';
      final markerStore = _FakeInstallationMarkerStore(
        markerPresent: true,
        markResult: const Failure(
          AppError(
            code: AppErrorCode.storageError,
            message: 'marker write failed',
          ),
        ),
      );
      final service = InstallationIdentityService(
        secureStorage: _MemorySecureStorage(
          initialValues: {StorageKeys.installationId: storedId},
        ),
        markerStore: markerStore,
      );

      final result = await service.resolve();

      expect(result, isA<Failure<String>>());
      expect(result.error?.message, 'marker write failed');
    });

    test('logs lifecycle context without full installation ids', () async {
      const oldId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';
      final messages = <String>[];
      final originalDebugPrint = debugPrint;
      debugPrint = (message, {wrapWidth}) {
        if (message != null) messages.add(message);
      };
      addTearDown(() => debugPrint = originalDebugPrint);

      final storage = _MemorySecureStorage(
        initialValues: {StorageKeys.installationId: oldId},
      );
      final service = InstallationIdentityService(
        secureStorage: storage,
        markerStore: _FakeInstallationMarkerStore(markerPresent: false),
      );

      final result = await service.resolve();

      expect(result, isA<Success<String>>());
      final logOutput = messages.join('\n');
      expect(logOutput, contains('InstallationIdentityService.resolve'));
      expect(logOutput, contains('Installation marker is missing'));
      expect(logOutput, isNot(contains(oldId)));
      expect(logOutput, isNot(contains(result.value)));
    });
  });
}

final _isCanonicalUuidV4 = matches(
  RegExp(
    r'^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$',
  ),
);

class _FakeInstallationMarkerStore implements InstallationMarkerStore {
  _FakeInstallationMarkerStore({
    required this.markerPresent,
    Result<bool>? hasMarkerResult,
    Result<Unit>? markResult,
  }) : _hasMarkerResult = hasMarkerResult,
       _markResult = markResult;

  bool markerPresent;
  int markCalls = 0;

  final Result<bool>? _hasMarkerResult;
  final Result<Unit>? _markResult;

  @override
  AsyncResult<bool> hasMarker() async {
    return _hasMarkerResult ?? Success(markerPresent);
  }

  @override
  AsyncResult<Unit> markResolved() async {
    markCalls++;
    if (_markResult != null) return _markResult;

    markerPresent = true;
    return Success(unit);
  }
}

class _MemorySecureStorage implements LocalSecureStorage {
  _MemorySecureStorage({
    Map<String, String>? initialValues,
    Result<String>? readResult,
    Result<Unit>? writeResult,
  }) : values = {...?initialValues},
       _readResult = readResult,
       _writeResult = writeResult;

  final Map<String, String> values;
  final Result<String>? _readResult;
  final Result<Unit>? _writeResult;

  int writeCalls = 0;

  @override
  AsyncResult<Unit> write(String key, String value) async {
    writeCalls++;
    if (_writeResult != null) return _writeResult;

    values[key] = value;
    return Success(unit);
  }

  @override
  AsyncResult<String> read(String key) async {
    if (_readResult != null) return _readResult;

    final value = values[key];
    if (value == null) {
      return const Failure(
        AppError(code: AppErrorCode.storageNotFound, message: 'Not found'),
      );
    }

    return Success(value);
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    values.remove(key);
    return Success(unit);
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    values.clear();
    return Success(unit);
  }

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    return values.keys.where((key) => key.startsWith(pattern)).toList();
  }
}
