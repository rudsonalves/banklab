import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/repositories/register_draft/register_draft_repository_impl.dart';
import 'package:bankflow/data/services/cache/last_login/register_draft/register_draft_store.dart';
import 'package:bankflow/domain/common/auth/models/register_draft.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RegisterDraftRepositoryImpl', () {
    test('returns found draft inside TTL', () async {
      final store = _FakeRegisterDraftStore(
        lookupResult: Future.value(
          Success<RegisterDraftLoadResult>(
            RegisterDraftFound(_snapshot(updatedAt: _now)),
          ),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(
        store,
        now: () => _now.add(const Duration(hours: 1)),
      );

      final result = await repository.getByCPF('123.456.789-09');

      expect(result, isA<Success<RegisterDraftLoadResult>>());
      expect(result.value, isA<RegisterDraftFound>());
      expect(repository.snapshot, isNotNull);
      expect(store.deleteCalls, 0);
    });

    test('returns not found and removes expired draft', () async {
      final store = _FakeRegisterDraftStore(
        lookupResult: Future.value(
          Success<RegisterDraftLoadResult>(
            RegisterDraftFound(
              _snapshot(updatedAt: _now.subtract(const Duration(hours: 25))),
            ),
          ),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(
        store,
        now: () => _now,
      );

      final result = await repository.getByCPF('123.456.789-09');

      expect(result, isA<Success<RegisterDraftLoadResult>>());
      expect(result.value, isA<RegisterDraftNotFound>());
      expect(repository.snapshot, isNull);
      expect(store.deleteCalls, 1);
      expect(store.deletedCpfs, contains('123.456.789-09'));
    });

    test('expires draft exactly at TTL boundary', () async {
      final store = _FakeRegisterDraftStore(
        lookupResult: Future.value(
          Success<RegisterDraftLoadResult>(
            RegisterDraftFound(_snapshot(updatedAt: _now)),
          ),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(
        store,
        now: () => _now.add(const Duration(hours: 24)),
      );

      final result = await repository.getByCPF('123.456.789-09');

      expect(result, isA<Success<RegisterDraftLoadResult>>());
      expect(result.value, isA<RegisterDraftNotFound>());
      expect(repository.snapshot, isNull);
      expect(store.deleteCalls, 1);
    });

    test('returns not found when draft is absent', () async {
      final store = _FakeRegisterDraftStore(
        lookupResult: Future.value(
          const Success<RegisterDraftLoadResult>(RegisterDraftNotFound()),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(store);

      final result = await repository.getByCPF('123.456.789-09');

      expect(result, isA<Success<RegisterDraftLoadResult>>());
      expect(result.value, isA<RegisterDraftNotFound>());
      expect(repository.snapshot, isNull);
      expect(store.deleteCalls, 0);
    });

    test('propagates store failure on load', () async {
      final store = _FakeRegisterDraftStore(
        lookupResult: Future.value(
          const Failure<RegisterDraftLoadResult>(
            AppError(code: AppErrorCode.unexpected, message: 'load failed'),
          ),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(store);

      final result = await repository.getByCPF('123.456.789-09');

      expect(result, isA<Failure<RegisterDraftLoadResult>>());
      expect(result.error?.message, 'load failed');
      expect(repository.snapshot, isNull);
    });

    test('propagates store failure on save', () async {
      final store = _FakeRegisterDraftStore(
        saveResult: Future.value(
          const Failure<Unit>(
            AppError(code: AppErrorCode.unexpected, message: 'save failed'),
          ),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(store);

      final result = await repository.save(_snapshot(updatedAt: _now));

      expect(result, isA<Failure<Unit>>());
      expect(result.error?.message, 'save failed');
      expect(repository.snapshot, isNull);
    });

    test('propagates store failure on delete', () async {
      final store = _FakeRegisterDraftStore(
        deleteResult: Future.value(
          const Failure<Unit>(
            AppError(code: AppErrorCode.unexpected, message: 'delete failed'),
          ),
        ),
      );
      final repository = RegisterDraftRepositoryImpl(store);

      final result = await repository.deleteByCPF('123.456.789-09');

      expect(result, isA<Failure<Unit>>());
      expect(result.error?.message, 'delete failed');
      expect(repository.snapshot, isNull);
    });
  });
}

final _now = DateTime.utc(2026, 5, 19, 12);

RegisterDraftSnapshot _snapshot({required DateTime updatedAt}) {
  return RegisterDraftSnapshot(
    cpf: '123.456.789-09',
    name: 'Maria Silva',
    birthDate: DateTime(1990, 1, 15),
    email: 'maria@example.com',
    phone: '+5527999999999',
    emailVerificationId: 'email-id',
    phoneVerificationId: null,
    isEmailVerified: false,
    isPhoneVerified: false,
    createdAt: updatedAt.subtract(const Duration(hours: 1)),
    updatedAt: updatedAt,
  );
}

class _FakeRegisterDraftStore extends RegisterDraftStore {
  final AsyncResult<RegisterDraftLoadResult> lookupResult;
  final AsyncResult<Unit> saveResult;
  final AsyncResult<Unit> deleteResult;

  int deleteCalls = 0;
  final List<String> deletedCpfs = [];

  _FakeRegisterDraftStore({
    AsyncResult<RegisterDraftLoadResult>? lookupResult,
    AsyncResult<Unit>? saveResult,
    AsyncResult<Unit>? deleteResult,
  }) : lookupResult =
           lookupResult ??
           Future.value(
             Success<RegisterDraftLoadResult>(RegisterDraftNotFound()),
           ),
       saveResult = saveResult ?? Future.value(const Success<Unit>(unit)),
       deleteResult = deleteResult ?? Future.value(const Success<Unit>(unit)),
       super(_FakeSecureStorage());

  @override
  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot) async {
    return saveResult;
  }

  @override
  AsyncResult<RegisterDraftLoadResult> getByCPF(String cpf) async {
    return lookupResult;
  }

  @override
  AsyncResult<Unit> deleteByCPF(String cpf) async {
    deleteCalls += 1;
    deletedCpfs.add(cpf);
    return deleteResult;
  }
}

class _FakeSecureStorage implements LocalSecureStorage {
  @override
  AsyncResult<Unit> write(String key, String value) async {
    return const Success(unit);
  }

  @override
  AsyncResult<String> read(String key) async {
    return const Success('');
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    return const Success(unit);
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    return const Success(unit);
  }

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    return const [];
  }
}
