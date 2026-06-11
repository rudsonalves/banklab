import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/services/cache/last_login/last_login_cache_service_impl.dart';
import 'package:bankflow/data/services/cache/last_login/models/last_login_identity.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('LastLoginCacheServiceImpl', () {
    test(
      'returns the remembered identity when identifier is an email',
      () async {
        final storage = _MemorySecureStorage({
          StorageKeys.lastLoginName: 'Maria Silva',
          StorageKeys.lastLoginIdentifier: 'customer@example.com',
        });
        final service = LastLoginCacheServiceImpl(storage);

        final result = await service.get();

        expect(result, isA<Success<LastLoginIdentity>>());
        expect(result.value?.name, 'Maria Silva');
        expect(result.value?.identifier, 'customer@example.com');
      },
    );

    test('rejects a legacy CPF identifier unsupported by login API', () async {
      final storage = _MemorySecureStorage({
        StorageKeys.lastLoginName: 'Maria Silva',
        StorageKeys.lastLoginIdentifier: '12345678901',
      });
      final service = LastLoginCacheServiceImpl(storage);

      final result = await service.get();

      expect(result, isA<Failure<LastLoginIdentity>>());
      expect(result.error?.code, AppErrorCode.invalidData);
    });
  });
}

class _MemorySecureStorage implements LocalSecureStorage {
  _MemorySecureStorage(this.values);

  final Map<String, String> values;

  @override
  AsyncResult<Unit> delete(String key) async {
    values.remove(key);
    return const Success(unit);
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    values.clear();
    return const Success(unit);
  }

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    return values.keys.where((key) => key.startsWith(pattern)).toList();
  }

  @override
  AsyncResult<String> read(String key) async {
    final value = values[key];
    if (value == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.storageNotFound,
          message: 'Not found',
        ),
      );
    }

    return Success(value);
  }

  @override
  AsyncResult<Unit> write(String key, String value) async {
    values[key] = value;
    return const Success(unit);
  }
}
