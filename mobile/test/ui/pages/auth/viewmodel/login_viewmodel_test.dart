import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/core/services/installation_identity/installation_identity.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/repositories/auth/auth_repository.dart';
import 'package:bankflow/data/services/apis/auth/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/cache/last_login/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/auth/models/auth_user.dart';
import 'package:bankflow/domain/common/user/enums/user_role.dart';
import 'package:bankflow/ui/pages/auth/viewmodel/login_viewmodel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('LoginViewModel.login', () {
    test('does not call login API when installation identity fails', () async {
      final authRepository = _FakeAuthRepository();
      final viewModel = LoginViewModel(
        authRepository: authRepository,
        installationIdentityService: _FakeInstallationIdentityService(
          result: const Failure(
            AppError(
              code: AppErrorCode.storageError,
              message: 'identity failed',
            ),
          ),
        ),
        appSection: AppSection(),
      );

      await viewModel.login.execute(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );

      expect(viewModel.login.isFailure, isTrue);
      expect(viewModel.login.error?.message, 'identity failed');
      expect(authRepository.loginCalls, 0);
    });

    test('calls login API after installation identity resolves', () async {
      final authRepository = _FakeAuthRepository();
      final viewModel = LoginViewModel(
        authRepository: authRepository,
        installationIdentityService: _FakeInstallationIdentityService(),
        appSection: AppSection(),
      );

      await viewModel.login.execute(
        LoginRequestDto(email: 'customer@example.com', password: '123456'),
      );

      expect(viewModel.login.isSuccess, isTrue);
      expect(authRepository.loginCalls, 1);
    });
  });
}

class _FakeAuthRepository implements AuthRepository {
  int loginCalls = 0;

  @override
  AuthUser get currentUser => NotLoggedUser();

  @override
  bool get isLoggedIn => false;

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async =>
      throw UnimplementedError();

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    loginCalls++;
    return Success(
      LoggedUser(
        accessToken: 'access-token',
        refreshToken: 'refresh-token',
        userId: 'user-1',
        email: dto.email,
        role: UserRole.customer,
        customerId: 'customer-1',
      ),
    );
  }

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
