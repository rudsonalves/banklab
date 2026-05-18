import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client/rest_client.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_response.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/repositories/auth/auth_repository_impl.dart';
import 'package:bankflow/data/services/auth/api/auth_api.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_confirm_request_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_confirm_response_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_request_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_request_response_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/auth/cache/last_login_cache_service.dart';
import 'package:bankflow/data/services/auth/cache/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_user.dart';
import 'package:bankflow/domain/common/auth/models/user_profile.dart';
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
          requestContactVerificationResult: Success(
            ContactVerificationRequestResponseDto(
              verificationId: 'verification-id-1',
              channel: 'email',
              target: 'customer@example.com',
              token: '123456',
              expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
            ),
          ),
          confirmContactVerificationResult: Success(
            ContactVerificationConfirmResponseDto(
              verificationToken: 'verified-token',
              channel: 'email',
              target: 'customer@example.com',
              verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
            ),
          ),
        );
        final storage = _FakeLocalSecureStorage();
        final cache = _FakeLastLoginCacheService();

        final repository = AuthRepositoryImpl(
          api: api,
          storage: storage,
          lastLoginCacheService: cache,
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
        expect(repository.userProfile, isNull);
      },
    );

    test(
      'loads profile and updates remembered identity on successful login',
      () async {
        final api = _FakeAuthApi(
          loginResult: Success(_loggedUser()),
          profileResult: Success(_profile()),
          requestContactVerificationResult: Success(
            ContactVerificationRequestResponseDto(
              verificationId: 'verification-id-1',
              channel: 'email',
              target: 'customer@example.com',
              token: '123456',
              expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
            ),
          ),
          confirmContactVerificationResult: Success(
            ContactVerificationConfirmResponseDto(
              verificationToken: 'verified-token',
              channel: 'email',
              target: 'customer@example.com',
              verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
            ),
          ),
        );
        final storage = _FakeLocalSecureStorage();
        final cache = _FakeLastLoginCacheService();

        final repository = AuthRepositoryImpl(
          api: api,
          storage: storage,
          lastLoginCacheService: cache,
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
        expect(cache.lastSavedIdentity?.identifier, 'customer@example.com');
        expect(repository.isLoggedIn, isTrue);
        expect(repository.userProfile?.name, 'Maria Silva');
      },
    );
  });

  group('AuthRepositoryImpl.contactVerification', () {
    test('request delegates to API and keeps auth session untouched', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_loggedUser()),
        profileResult: Success(_profile()),
        requestContactVerificationResult: Success(
          ContactVerificationRequestResponseDto(
            verificationId: 'verification-id-1',
            channel: 'email',
            target: 'customer@example.com',
            token: '123456',
            expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
          ),
        ),
        confirmContactVerificationResult: Success(
          ContactVerificationConfirmResponseDto(
            verificationToken: 'verified-token',
            channel: 'email',
            target: 'customer@example.com',
            verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
          ),
        ),
      );
      final storage = _FakeLocalSecureStorage();
      final cache = _FakeLastLoginCacheService();

      final repository = AuthRepositoryImpl(
        api: api,
        storage: storage,
        lastLoginCacheService: cache,
      );

      final result = await repository.requestContactVerification(
        ContactVerificationRequestDto(
          channel: 'email',
          target: 'customer@example.com',
        ),
      );

      expect(result, isA<Success<ContactVerificationRequestResponseDto>>());
      expect(api.requestContactVerificationCalls, 1);
      expect(storage.writeCalls, 0);
      expect(repository.isLoggedIn, isFalse);
      expect(repository.userProfile, isNull);
    });

    test('request propagates AppError failure', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_loggedUser()),
        profileResult: Success(_profile()),
        requestContactVerificationResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'target is invalid',
          ),
        ),
        confirmContactVerificationResult: Success(
          ContactVerificationConfirmResponseDto(
            verificationToken: 'verified-token',
            channel: 'email',
            target: 'customer@example.com',
            verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
          ),
        ),
      );

      final repository = AuthRepositoryImpl(
        api: api,
        storage: _FakeLocalSecureStorage(),
        lastLoginCacheService: _FakeLastLoginCacheService(),
      );

      final result = await repository.requestContactVerification(
        ContactVerificationRequestDto(
          channel: 'email',
          target: 'invalid',
        ),
      );

      expect(result, isA<Failure<ContactVerificationRequestResponseDto>>());
      expect(result.error?.message, 'target is invalid');
    });

    test('confirm delegates to API and keeps auth session untouched', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_loggedUser()),
        profileResult: Success(_profile()),
        requestContactVerificationResult: Success(
          ContactVerificationRequestResponseDto(
            verificationId: 'verification-id-1',
            channel: 'phone',
            target: '+5511999999999',
            token: '123456',
            expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
          ),
        ),
        confirmContactVerificationResult: Success(
          ContactVerificationConfirmResponseDto(
            verificationToken: 'verified-token',
            channel: 'phone',
            target: '+5511999999999',
            verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
          ),
        ),
      );
      final storage = _FakeLocalSecureStorage();
      final cache = _FakeLastLoginCacheService();

      final repository = AuthRepositoryImpl(
        api: api,
        storage: storage,
        lastLoginCacheService: cache,
      );

      final result = await repository.confirmContactVerification(
        ContactVerificationConfirmRequestDto(
          verificationId: 'verification-id-1',
          token: '123456',
        ),
      );

      expect(result, isA<Success<ContactVerificationConfirmResponseDto>>());
      expect(api.confirmContactVerificationCalls, 1);
      expect(storage.writeCalls, 0);
      expect(repository.isLoggedIn, isFalse);
      expect(repository.userProfile, isNull);
    });

    test('confirm propagates AppError failure', () async {
      final api = _FakeAuthApi(
        loginResult: Success(_loggedUser()),
        profileResult: Success(_profile()),
        requestContactVerificationResult: Success(
          ContactVerificationRequestResponseDto(
            verificationId: 'verification-id-1',
            channel: 'phone',
            target: '+5511999999999',
            token: '123456',
            expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
          ),
        ),
        confirmContactVerificationResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'invalid verification token',
          ),
        ),
      );

      final repository = AuthRepositoryImpl(
        api: api,
        storage: _FakeLocalSecureStorage(),
        lastLoginCacheService: _FakeLastLoginCacheService(),
      );

      final result = await repository.confirmContactVerification(
        ContactVerificationConfirmRequestDto(
          verificationId: 'verification-id-1',
          token: '000000',
        ),
      );

      expect(result, isA<Failure<ContactVerificationConfirmResponseDto>>());
      expect(result.error?.message, 'invalid verification token');
    });
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

UserProfile _profile() {
  return UserProfile(
    userId: 'user-1',
    customerId: 'customer-1',
    name: 'Maria Silva',
    email: 'customer@example.com',
    role: UserRole.customer,
    createdAt: DateTime(2026, 5, 13),
    updatedAt: DateTime(2026, 5, 13),
  );
}

class _FakeAuthApi extends AuthApi {
  Result<LoggedUser> loginResult;
  Result<UserProfile> profileResult;
  Result<ContactVerificationRequestResponseDto>
  requestContactVerificationResult;
  Result<ContactVerificationConfirmResponseDto>
  confirmContactVerificationResult;

  int loginCalls = 0;
  int getProfileCalls = 0;
  int requestContactVerificationCalls = 0;
  int confirmContactVerificationCalls = 0;

  _FakeAuthApi({
    required this.loginResult,
    required this.profileResult,
    required this.requestContactVerificationResult,
    required this.confirmContactVerificationResult,
  }) : super(_NoopRestClient());

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    loginCalls++;
    return loginResult;
  }

  @override
  AsyncResult<UserProfile> getProfile() async {
    getProfileCalls++;
    return profileResult;
  }

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    requestContactVerificationCalls++;
    return requestContactVerificationResult;
  }

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    confirmContactVerificationCalls++;
    return confirmContactVerificationResult;
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
