import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/core/services/client_http/client/rest_client.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_response.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/repositories/auth/auth_repository_impl.dart';
import 'package:bankflow/data/services/apis/auth/auth_api.dart';
import 'package:bankflow/data/services/apis/auth/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/cache/last_login/last_login_cache_service.dart';
import 'package:bankflow/data/services/cache/last_login/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/auth/models/auth_user.dart';
import 'package:bankflow/domain/common/user/enums/user_role.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('AuthRepositoryImpl.login', () {
    test(
      'does not persist tokens, load profile, or save remembered cache on accountApprovalRequired failure',
      () async {
        final api = _FakeAuthApi(
          loginResult: const Failure(
            AppError(
              code: AppErrorCode.accountApprovalRequired,
              message:
                  'Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você poderá acessar sua conta.',
            ),
          ),
          profileResult: Success(_profile()),
        );
        final storage = _FakeLocalSecureStorage();
        final cache = _FakeLastLoginCacheService();
        final appSection = AppSection();

        final repository = AuthRepositoryImpl(
          api: api,
          storage: storage,
          lastLoginCacheService: cache,
          appSection: appSection,
        );

        final result = await repository.login(
          LoginRequestDto(email: 'customer@example.com', password: '123456'),
        );

        expect(result, isA<Failure<LoggedUser>>());
        expect(result.error?.code, AppErrorCode.accountApprovalRequired);
        expect(api.loginCalls, 1);
        expect(api.getProfileCalls, 0);
        expect(storage.writeCalls, 0);
        expect(storage.writesByKey, isEmpty);
        expect(cache.saveCalls, 0);
        expect(repository.isLoggedIn, isFalse);
        expect(appSection.currentSession, isNull);
      },
    );

    test(
      'loads profile and updates remembered identity on successful login',
      () async {
        final api = _FakeAuthApi(
          loginResult: Success(_loggedUser()),
          profileResult: Success(_profile()),
        );
        final storage = _FakeLocalSecureStorage();
        final cache = _FakeLastLoginCacheService();
        final appSection = AppSection();

        final repository = AuthRepositoryImpl(
          api: api,
          storage: storage,
          lastLoginCacheService: cache,
          appSection: appSection,
        );

        final result = await repository.login(
          LoginRequestDto(email: 'customer@example.com', password: '123456'),
        );

        expect(result, isA<Success<LoggedUser>>());
        expect(api.loginCalls, 1);
        expect(api.getProfileCalls, 1);
        expect(storage.writeCalls, 2);
        expect(storage.writesByKey[StorageKeys.accessToken], 'access-token');
        expect(storage.writesByKey[StorageKeys.refreshToken], 'refresh-token');
        expect(cache.saveCalls, 1);
        expect(cache.lastSavedIdentity?.name, 'Maria Silva');
        expect(
          cache.lastSavedIdentity?.identifier,
          'customer@example.com',
        );
        expect(repository.isLoggedIn, isTrue);
        expect(appSection.currentSession?.customer?.name, 'Maria Silva');
      },
    );
  });

  group('AuthRepositoryImpl session lifecycle', () {
    test(
      'getAuthSession returns the loaded snapshot without calling API again',
      () async {
        final api = _FakeAuthApi(
          loginResult: Success(_loggedUser()),
          profileResult: Success(_profile()),
        );
        final appSection = AppSection();
        final repository = AuthRepositoryImpl(
          api: api,
          storage: _FakeLocalSecureStorage(),
          lastLoginCacheService: _FakeLastLoginCacheService(),
          appSection: appSection,
        );

        await repository.login(
          LoginRequestDto(email: 'customer@example.com', password: '123456'),
        );
        final result = await repository.getAuthSession();

        expect(result, isA<Success<AuthSession>>());
        expect(result.value, same(appSection.currentSession));
        expect(api.getProfileCalls, 1);
      },
    );

    test('logout clears the loaded AppSection snapshot', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_loggedUser()),
        profileResult: Success(_profile()),
      );
      final appSection = AppSection();
      final repository = AuthRepositoryImpl(
        api: api,
        storage: _FakeLocalSecureStorage(),
        lastLoginCacheService: _FakeLastLoginCacheService(),
        appSection: appSection,
      );

      await repository.login(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );
      expect(appSection.currentSession, isNotNull);

      final result = await repository.logout();

      expect(result, isA<Success<Unit>>());
      expect(appSection.currentSession, isNull);
      expect(repository.isLoggedIn, isFalse);
    });

    test(
      'logout clears AppSection even when repository is not logged in',
      () async {
        final appSection = AppSection()..setAuthSession(_profile());
        final repository = AuthRepositoryImpl(
          api: _FakeAuthApi(
            loginResult: Success(_loggedUser()),
            profileResult: Success(_profile()),
          ),
          storage: _FakeLocalSecureStorage(),
          lastLoginCacheService: _FakeLastLoginCacheService(),
          appSection: appSection,
        );

        final result = await repository.logout();

        expect(result, isA<Success<Unit>>());
        expect(appSection.currentSession, isNull);
      },
    );
  });
}

LoggedUser _loggedUser() {
  return LoggedUser(
    accessToken: 'access-token',
    refreshToken: 'refresh-token',
    userId: 'user-1',
    email: 'customer@example.com',
    role: UserRole.customer,
    customerId: 'customer-1',
  );
}

AuthSession _profile() {
  return AuthSession(
    user: UserSession(
      userId: 'user-1',
      email: 'customer@example.com',
      role: UserRole.customer,
    ),
    customer: CustomerSession(
      id: 'customer-1',
      name: 'Maria Silva',
      cpf: '12345678901',
      birthDate: DateTime(1990, 1, 1),
      createdAt: DateTime(2026, 5, 13),
    ),
    readiness: ReadinessSession(
      onboardingCompleted: true,
      approved: true,
      hasOperationalAccount: true,
      transactionPasswordStatus: TransactionPasswordStatus.active,
    ),
  );
}

class _FakeAuthApi extends AuthApi {
  Result<LoggedUser> loginResult;
  Result<AuthSession> profileResult;

  int loginCalls = 0;
  int getProfileCalls = 0;

  _FakeAuthApi({
    required this.loginResult,
    required this.profileResult,
  }) : super(_NoopRestClient());

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    loginCalls++;
    return loginResult;
  }

  @override
  AsyncResult<AuthSession> getAuthSession() async {
    getProfileCalls++;
    return profileResult;
  }
}

class _FakeLocalSecureStorage implements LocalSecureStorage {
  int writeCalls = 0;
  final Map<String, String> writesByKey = {};

  @override
  AsyncResult<Unit> write(String key, String value) async {
    writeCalls++;
    writesByKey[key] = value;
    return Success(unit);
  }

  @override
  AsyncResult<String> read(String key) async {
    final value = writesByKey[key];
    if (value == null) {
      return const Failure(
        AppError(code: AppErrorCode.storageNotFound, message: 'Not found'),
      );
    }

    return Success(value);
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    writesByKey.remove(key);
    return Success(unit);
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    writesByKey.clear();
    return Success(unit);
  }

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    return writesByKey.keys.where((key) => key.startsWith(pattern)).toList();
  }
}

class _FakeLastLoginCacheService implements LastLoginCacheService {
  int saveCalls = 0;
  LastLoginIdentity? lastSavedIdentity;

  @override
  AsyncResult<LastLoginIdentity> get() async {
    if (lastSavedIdentity == null) {
      return const Failure(
        AppError(code: AppErrorCode.storageNotFound, message: 'Not found'),
      );
    }

    return Success(lastSavedIdentity!);
  }

  @override
  AsyncResult<Unit> save(LastLoginIdentity identity) async {
    saveCalls++;
    lastSavedIdentity = identity;
    return Success(unit);
  }

  @override
  AsyncResult<Unit> clear() async {
    lastSavedIdentity = null;
    return Success(unit);
  }
}

class _NoopRestClient implements RestClient {
  const _NoopRestClient();

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );
}
