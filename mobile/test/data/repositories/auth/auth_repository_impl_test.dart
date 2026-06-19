import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/core/services/client_http/client/rest_client.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_response.dart';
import 'package:bankflow/core/services/installation_identity/installation_identity.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/repositories/auth/auth_repository_impl.dart';
import 'package:bankflow/data/repositories/transaction_password/transaction_password_repository.dart';
import 'package:bankflow/data/services/apis/auth/auth_api.dart';
import 'package:bankflow/data/services/apis/auth/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/apis/installation/dtos/installation_registration_response_dto.dart';
import 'package:bankflow/data/services/apis/installation/installation_api.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';
import 'package:bankflow/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import 'package:bankflow/data/services/cache/last_login/last_login_cache_service.dart';
import 'package:bankflow/data/services/cache/last_login/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/auth/models/auth_state.dart';
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

        final repository = _repository(
          api: api,
          storage: storage,
          lastLoginCacheService: cache,
          appSection: appSection,
        );

        final result = await repository.login(
          LoginRequestDto(email: 'customer@example.com', password: '123456'),
        );

        expect(result, isA<Failure<AuthState>>());
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

        final repository = _repository(
          api: api,
          storage: storage,
          lastLoginCacheService: cache,
          appSection: appSection,
        );

        final result = await repository.login(
          LoginRequestDto(email: 'customer@example.com', password: '123456'),
        );

        expect(result, isA<Success<OperationalAuthState>>());
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

    test('restricted login does not persist tokens or load profile', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_restrictedLogin()),
        profileResult: Success(_profile()),
      );
      final storage = _FakeLocalSecureStorage();
      final cache = _FakeLastLoginCacheService();
      final appSection = AppSection();

      final repository = _repository(
        api: api,
        storage: storage,
        lastLoginCacheService: cache,
        appSection: appSection,
      );

      final result = await repository.login(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );

      expect(result, isA<Success<AuthState>>());
      expect(result.value, isA<RestrictedInstallationAuthState>());
      expect(api.loginCalls, 1);
      expect(api.getProfileCalls, 0);
      expect(storage.writeCalls, 0);
      expect(storage.writesByKey, isEmpty);
      expect(cache.saveCalls, 0);
      expect(repository.isLoggedIn, isFalse);
      expect(appSection.currentSession, isNull);
    });

    test('installation limit reached does not persist tokens', () async {
      final api = _FakeAuthApi(
        loginResult: const Failure(
          AppError(
            code: AppErrorCode.installationLimitReached,
            message: 'Installation limit reached.',
          ),
        ),
        profileResult: Success(_profile()),
      );
      final storage = _FakeLocalSecureStorage();
      final cache = _FakeLastLoginCacheService();
      final appSection = AppSection();

      final repository = _repository(
        api: api,
        storage: storage,
        lastLoginCacheService: cache,
        appSection: appSection,
      );

      final result = await repository.login(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );

      expect(result, isA<Failure<AuthState>>());
      expect(result.error?.code, AppErrorCode.installationLimitReached);
      expect(api.loginCalls, 1);
      expect(api.getProfileCalls, 0);
      expect(storage.writeCalls, 0);
      expect(storage.writesByKey, isEmpty);
      expect(cache.saveCalls, 0);
      expect(repository.isLoggedIn, isFalse);
      expect(appSection.currentSession, isNull);
    });
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
        final repository = _repository(
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
      final storage = _FakeLocalSecureStorage();
      final appSection = AppSection();
      final repository = _repository(
        api: api,
        storage: storage,
        lastLoginCacheService: _FakeLastLoginCacheService(),
        appSection: appSection,
      );

      await repository.login(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );
      await storage.write(
        StorageKeys.installationId,
        '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8',
      );
      expect(appSection.currentSession, isNotNull);

      final result = await repository.logout();

      expect(result, isA<Success<Unit>>());
      expect(appSection.currentSession, isNull);
      expect(repository.isLoggedIn, isFalse);
      expect(storage.writesByKey, contains(StorageKeys.installationId));
      expect(storage.deleteCallsByKey, [
        StorageKeys.accessToken,
        StorageKeys.refreshToken,
      ]);
    });

    test(
      'logout clears AppSection even when repository is not logged in',
      () async {
        final appSection = AppSection()..setAuthSession(_profile());
        final repository = _repository(
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

  group('AuthRepositoryImpl.certifyInstallation', () {
    test('registers installation and starts an operational session', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_restrictedLogin()),
        profileResult: Success(_profile()),
      );
      final installationApi = _FakeInstallationApi(
        result: const Success(
          InstallationRegistrationResponseDto(
            accessToken: 'registered-access-token',
            refreshToken: 'registered-refresh-token',
            installationResourceId: 'installation-resource-1',
            installationStatus: 'active',
          ),
        ),
      );
      final transactionPasswordRepository = _FakeTransactionPasswordRepository(
        result: Success(
          StepUpAuthorizeResponseDto(
            stepUpToken: 'step-up-token',
            expiresIn: 120,
          ),
        ),
      );
      final storage = _FakeLocalSecureStorage();
      final cache = _FakeLastLoginCacheService();
      final appSection = AppSection();
      final repository = _repository(
        api: api,
        installationApi: installationApi,
        transactionPasswordRepository: transactionPasswordRepository,
        storage: storage,
        lastLoginCacheService: cache,
        appSection: appSection,
      );

      await repository.login(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );
      final result = await repository.certifyInstallation('654321');

      expect(result, isA<Success<OperationalAuthState>>());
      expect(result.value?.accessToken, 'registered-access-token');
      expect(transactionPasswordRepository.authorizeInstallationCalls, 1);
      expect(transactionPasswordRepository.lastTransactionPassword, '654321');
      expect(installationApi.registerCalls, 1);
      expect(installationApi.lastRestrictedAccessToken, 'restricted-token');
      expect(installationApi.lastStepUpToken, 'step-up-token');
      expect(
        storage.writesByKey[StorageKeys.accessToken],
        'registered-access-token',
      );
      expect(
        storage.writesByKey[StorageKeys.refreshToken],
        'registered-refresh-token',
      );
      expect(api.getProfileCalls, 1);
      expect(cache.saveCalls, 1);
      expect(repository.isLoggedIn, isTrue);
      expect(appSection.currentSession?.customer?.name, 'Maria Silva');
    });

    test(
      'does not register installation when step-up authorization fails',
      () async {
        final installationApi = _FakeInstallationApi(
          result: const Failure(
            AppError(code: AppErrorCode.unexpected, message: 'Should not call'),
          ),
        );
        final repository = _repository(
          api: _FakeAuthApi(
            loginResult: Success(_restrictedLogin()),
            profileResult: Success(_profile()),
          ),
          installationApi: installationApi,
          transactionPasswordRepository: _FakeTransactionPasswordRepository(
            result: const Failure(
              AppError(
                code: AppErrorCode.httpError,
                message: 'Invalid transaction password',
              ),
            ),
          ),
          storage: _FakeLocalSecureStorage(),
          lastLoginCacheService: _FakeLastLoginCacheService(),
          appSection: AppSection(),
        );

        await repository.login(
          LoginRequestDto(email: 'customer@example.com', password: '123456'),
        );
        final result = await repository.certifyInstallation('000000');

        expect(result, isA<Failure<OperationalAuthState>>());
        expect(result.error?.message, 'Invalid transaction password');
        expect(installationApi.registerCalls, 0);
        expect(repository.currentUser, isA<AnonymousAuthState>());
      },
    );
  });
}

AuthRepositoryImpl _repository({
  required AuthApi api,
  InstallationApi? installationApi,
  TransactionPasswordRepository? transactionPasswordRepository,
  required LocalSecureStorage storage,
  required LastLoginCacheService lastLoginCacheService,
  required AppSection appSection,
}) {
  return AuthRepositoryImpl(
    api: api,
    installationApi: installationApi ?? _FakeInstallationApi(),
    transactionPasswordRepository:
        transactionPasswordRepository ?? _FakeTransactionPasswordRepository(),
    storage: storage,
    lastLoginCacheService: lastLoginCacheService,
    appSection: appSection,
  );
}

OperationalAuthState _loggedUser() {
  return OperationalAuthState(
    accessToken: 'access-token',
    refreshToken: 'refresh-token',
    userId: 'user-1',
    email: 'customer@example.com',
    role: UserRole.customer,
    customerId: 'customer-1',
  );
}

RestrictedInstallationAuthState _restrictedLogin() {
  return RestrictedInstallationAuthState(
    restrictedAccessToken: 'restricted-token',
    restrictedTokenType: 'restricted_access',
    restrictedScope: 'installation.register',
    restrictedExpiresAt: DateTime.parse('2026-06-17T10:05:00Z'),
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
  Result<AuthState> loginResult;
  Result<AuthSession> profileResult;

  int loginCalls = 0;
  int getProfileCalls = 0;

  _FakeAuthApi({
    required this.loginResult,
    required this.profileResult,
  }) : super(_NoopRestClient());

  @override
  AsyncResult<AuthState> login(LoginRequestDto dto) async {
    loginCalls++;
    return loginResult;
  }

  @override
  AsyncResult<AuthSession> getAuthSession() async {
    getProfileCalls++;
    return profileResult;
  }
}

class _FakeInstallationApi extends InstallationApi {
  _FakeInstallationApi({
    this.result = const Success(
      InstallationRegistrationResponseDto(
        accessToken: 'access-token',
        refreshToken: 'refresh-token',
        installationResourceId: 'installation-resource-1',
        installationStatus: 'active',
      ),
    ),
  }) : super(
         client: const _NoopRestClient(),
         installationIdentityService: _NoopInstallationIdentityService(),
       );

  final Result<InstallationRegistrationResponseDto> result;
  int registerCalls = 0;
  String? lastRestrictedAccessToken;
  String? lastStepUpToken;

  @override
  AsyncResult<InstallationRegistrationResponseDto> register({
    required String restrictedAccessToken,
    required String stepUpToken,
  }) async {
    registerCalls++;
    lastRestrictedAccessToken = restrictedAccessToken;
    lastStepUpToken = stepUpToken;
    return result;
  }
}

class _FakeTransactionPasswordRepository
    implements TransactionPasswordRepository {
  _FakeTransactionPasswordRepository({
    this.result = const Failure(
      AppError(code: AppErrorCode.unexpected, message: 'Not configured'),
    ),
  });

  final Result<StepUpAuthorizeResponseDto> result;
  int authorizeInstallationCalls = 0;
  String? lastTransactionPassword;

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInstallationRegistration(
    String transactionPassword,
  ) async {
    authorizeInstallationCalls++;
    lastTransactionPassword = transactionPassword;
    return result;
  }

  @override
  AsyncResult<StepUpAuthorizeResponseDto> authorizeInternalTransfer(
    String transactionPassword,
  ) {
    throw UnimplementedError();
  }

  @override
  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) {
    throw UnimplementedError();
  }
}

class _NoopInstallationIdentityService extends InstallationIdentityService {
  _NoopInstallationIdentityService()
    : super(
        secureStorage: _FakeLocalSecureStorage(),
        markerStore: _NoopInstallationMarkerStore(),
      );

  @override
  AsyncResult<String> resolve() async {
    return const Success('018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8');
  }
}

class _NoopInstallationMarkerStore implements InstallationMarkerStore {
  @override
  AsyncResult<bool> hasMarker() async => const Success(true);

  @override
  AsyncResult<Unit> markResolved() async => const Success(unit);
}

class _FakeLocalSecureStorage implements LocalSecureStorage {
  int writeCalls = 0;
  final Map<String, String> writesByKey = {};
  final List<String> deleteCallsByKey = [];

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
    deleteCallsByKey.add(key);
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
