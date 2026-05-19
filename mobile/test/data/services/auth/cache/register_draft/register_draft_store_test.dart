import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/services/auth/cache/register_draft/register_draft_store.dart';
import 'package:bankflow/domain/common/auth/models/register_draft.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RegisterDraftStore', () {
    test('builds storage key from normalized cpf hash', () {
      final store = RegisterDraftStore(_FakeSecureStorage());

      final key = store.keyForCPF('123.456.789-09');

      expect(
        key,
        'onboarding_draft:'
        '7ec94663084bd506d4f0c3e21042df233681fd7426e93f397c921b1d3e397bba',
      );
      expect(key.contains('12345678909'), isFalse);
      expect(key.contains('123.456.789-09'), isFalse);
    });

    test('saves and retrieves draft snapshot', () async {
      final storage = _FakeSecureStorage();
      final store = RegisterDraftStore(storage);
      final snapshot = _snapshot();

      final saveResult = await store.save(snapshot);
      final lookupResult = await store.getByCPF('123.456.789-09');

      expect(saveResult.isSuccess, isTrue);
      expect(lookupResult.isSuccess, isTrue);
      expect(lookupResult.value, isA<RegisterDraftFound>());
      final found = lookupResult.value as RegisterDraftFound;
      expect(found.snapshot.cpf, '12345678909');
      expect(found.snapshot.name, 'Maria Silva');
      expect(found.snapshot.currentStep, RegisterDraftStep.email);
      expect(storage.values.keys.single, store.keyForCPF('12345678909'));
    });

    test('returns empty lookup when draft is absent', () async {
      final store = RegisterDraftStore(_FakeSecureStorage());

      final result = await store.getByCPF('12345678909');

      expect(result.isSuccess, isTrue);
      expect(result.value, isA<RegisterDraftNotFound>());
    });

    test('deletes invalid json and returns empty lookup', () async {
      final storage = _FakeSecureStorage();
      final store = RegisterDraftStore(storage);
      final key = store.keyForCPF('12345678909');
      storage.values[key] = '{invalid-json';

      final result = await store.getByCPF('12345678909');

      expect(result.isSuccess, isTrue);
      expect(result.value, isA<RegisterDraftNotFound>());
      expect(storage.values.containsKey(key), isFalse);
      expect(storage.deletedKeys, contains(key));
    });

    test('deletes invalid snapshot payload and returns empty lookup', () async {
      final storage = _FakeSecureStorage();
      final store = RegisterDraftStore(storage);
      final key = store.keyForCPF('12345678909');
      storage.values[key] = '{"cpf":"12345678909","current_step":"unknown"}';

      final result = await store.getByCPF('12345678909');

      expect(result.isSuccess, isTrue);
      expect(result.value, isA<RegisterDraftNotFound>());
      expect(storage.values.containsKey(key), isFalse);
    });

    test('deletes draft by cpf', () async {
      final storage = _FakeSecureStorage();
      final store = RegisterDraftStore(storage);
      await store.save(_snapshot());

      final result = await store.deleteByCPF('123.456.789-09');

      expect(result.isSuccess, isTrue);
      expect(storage.values, isEmpty);
    });
  });
}

RegisterDraftSnapshot _snapshot() {
  return RegisterDraftSnapshot(
    cpf: '123.456.789-09',
    name: 'Maria Silva',
    birthDate: DateTime(1990, 1, 15),
    email: 'maria@example.com',
    phone: '+5527999999999',
    currentStep: RegisterDraftStep.email,
    emailVerificationId: 'email-id',
    phoneVerificationId: null,
    isEmailVerified: false,
    isPhoneVerified: false,
    createdAt: DateTime.utc(2026, 5, 19, 10),
    updatedAt: DateTime.utc(2026, 5, 19, 11),
  );
}

class _FakeSecureStorage implements LocalSecureStorage {
  final Map<String, String> values = {};
  final List<String> deletedKeys = [];

  @override
  AsyncResult<Unit> write(String key, String value) async {
    values[key] = value;
    return const Success(unit);
  }

  @override
  AsyncResult<String> read(String key) async {
    final value = values[key];
    if (value == null) {
      return Failure(
        AppError(
          code: AppErrorCode.storageNotFound,
          message: 'Key not found: $key',
        ),
      );
    }

    return Success(value);
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    deletedKeys.add(key);
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
}
