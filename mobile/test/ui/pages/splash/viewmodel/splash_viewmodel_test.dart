import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/installation_identity/installation_identity.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/repositories/auth/auth_repository.dart';
import 'package:bankflow/data/services/apis/auth/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/cache/last_login/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/auth/models/auth_state.dart';
import 'package:bankflow/ui/pages/splash/viewmodel/splash_viewmodel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('SplashViewmodel.initialize', () {
    test('blocks bootstrap when installation identity fails', () async {
      final authRepository = _FakeAuthRepository(
        lastLoginResult: Success(_identity()),
      );
      final viewModel = SplashViewmodel(
        authRepository: authRepository,
        installationIdentityService: _FakeInstallationIdentityService(
          result: const Failure(
            AppError(
              code: AppErrorCode.storageError,
              message: 'identity failed',
            ),
          ),
        ),
      );

      await viewModel.initialize.execute();

      expect(viewModel.initialize.isFailure, isTrue);
      expect(viewModel.initialize.error?.message, 'identity failed');
      expect(authRepository.getLastLoginCalls, 0);
    });

    test(
      'returns remembered identity after installation identity resolves',
      () async {
        final identity = _identity();
        final viewModel = SplashViewmodel(
          authRepository: _FakeAuthRepository(
            lastLoginResult: Success(identity),
          ),
          installationIdentityService: _FakeInstallationIdentityService(),
        );

        await viewModel.initialize.execute();

        expect(viewModel.initialize.isSuccess, isTrue);
        expect(viewModel.initialize.value?.lastLoginIdentity, identity);
      },
    );

    test('continues bootstrap without remembered identity', () async {
      final viewModel = SplashViewmodel(
        authRepository: _FakeAuthRepository(
          lastLoginResult: const Failure(
            AppError(code: AppErrorCode.storageNotFound, message: 'not found'),
          ),
        ),
        installationIdentityService: _FakeInstallationIdentityService(),
      );

      await viewModel.initialize.execute();

      expect(viewModel.initialize.isSuccess, isTrue);
      expect(viewModel.initialize.value?.lastLoginIdentity, isNull);
    });
  });
}

LastLoginIdentity _identity() {
  return LastLoginIdentity(
    name: 'Maria Silva',
    identifier: 'customer@example.com',
  );
}

class _FakeAuthRepository implements AuthRepository {
  _FakeAuthRepository({required this.lastLoginResult});

  final Result<LastLoginIdentity> lastLoginResult;
  int getLastLoginCalls = 0;

  @override
  AuthState get currentUser => AnonymousAuthState();

  @override
  bool get isLoggedIn => false;

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async {
    getLastLoginCalls++;
    return lastLoginResult;
  }

  @override
  AsyncResult<OperationalAuthState> login(LoginRequestDto dto) async =>
      throw UnimplementedError();

  @override
  AsyncResult<Unit> logout() async => throw UnimplementedError();

  @override
  AsyncResult<AuthSession> getAuthSession() async => throw UnimplementedError();
}

class _FakeInstallationIdentityService extends InstallationIdentityService {
  _FakeInstallationIdentityService({
    Result<String> result = const Success(
      '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8',
    ),
  }) : _result = result,
       super(
         secureStorage: _NoopSecureStorage(),
         markerStore: _NoopMarkerStore(),
       );

  final Result<String> _result;

  @override
  AsyncResult<String> resolve() async => _result;
}

class _NoopMarkerStore implements InstallationMarkerStore {
  @override
  AsyncResult<bool> hasMarker() async => const Success(true);

  @override
  AsyncResult<Unit> markResolved() async => const Success(unit);
}

class _NoopSecureStorage implements LocalSecureStorage {
  @override
  AsyncResult<Unit> write(String key, String value) async =>
      const Success(unit);

  @override
  AsyncResult<String> read(String key) async =>
      const Success('018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8');

  @override
  AsyncResult<Unit> delete(String key) async => const Success(unit);

  @override
  AsyncResult<Unit> deleteAll() async => const Success(unit);

  @override
  Future<List<String>> keysWithPrefix(String pattern) async => const [];
}
