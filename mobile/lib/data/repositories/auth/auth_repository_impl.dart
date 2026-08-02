import '/core/resources/storage_keys.dart';
import '/core/result/result.dart';
import '/core/services/app_section/app_section.dart';
import '/core/services/logging/console_log.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/auth_api.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/data/services/apis/installation/installation_api.dart';
import '/data/services/cache/last_login/last_login_cache_service.dart';
import '/data/services/cache/last_login/models/last_login_identity.dart';
import '/data/repositories/transaction_password/transaction_password_repository.dart';
import '/domain/common/auth/models/auth_session/auth_session.dart';
import '/domain/common/auth/models/auth_state.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthApi _api;
  final InstallationApi _installationApi;
  final TransactionPasswordRepository _transactionPasswordRepository;
  final LocalSecureStorage _storage;
  final LastLoginCacheService _lastLoginCacheService;
  final AppSection _appSection;

  AuthRepositoryImpl({
    required AuthApi api,
    required InstallationApi installationApi,
    required TransactionPasswordRepository transactionPasswordRepository,
    required LocalSecureStorage storage,
    required LastLoginCacheService lastLoginCacheService,
    required AppSection appSection,
  }) : _api = api,
       _installationApi = installationApi,
       _transactionPasswordRepository = transactionPasswordRepository,
       _storage = storage,
       _lastLoginCacheService = lastLoginCacheService,
       _appSection = appSection;

  AuthState _currentUser = AnonymousAuthState();

  final _log = ConsoleLog('AuthRepositoryImpl');

  @override
  AuthState get currentUser => _currentUser;

  @override
  bool get isLoggedIn => _currentUser is OperationalAuthState;

  @override
  AsyncResult<AuthState> login(LoginRequestDto dto) async {
    if (isLoggedIn) return Success(_currentUser as OperationalAuthState);

    final result = await _api.login(dto);
    if (result.isFailure) return Result.failure(result.error!);

    final loginResult = result.value!;
    if (loginResult is RestrictedInstallationAuthState) {
      _currentUser = loginResult;
      return Success(loginResult);
    }

    if (loginResult is! OperationalAuthState) {
      return const Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Unknown login result.',
        ),
      );
    }

    return _startOperationalSession(loginResult);
  }

  @override
  AsyncResult<OperationalAuthState> certifyInstallation(
    String transactionPassword,
  ) async {
    final restrictedUser = _currentUser;
    if (restrictedUser is! RestrictedInstallationAuthState) {
      return const Failure(
        AppError(
          code: AppErrorCode.unauthenticated,
          message: 'Restricted installation registration is not active.',
        ),
      );
    }

    final stepUpResult = await _transactionPasswordRepository
        .authorizeInstallationRegistration(transactionPassword);
    if (stepUpResult.isFailure) {
      _discardTemporaryAuthState();
      return Result.failure(stepUpResult.error!);
    }

    final registrationResult = await _installationApi.register(
      restrictedAccessToken: restrictedUser.restrictedAccessToken,
      stepUpToken: stepUpResult.value!.stepUpToken,
    );
    if (registrationResult.isFailure) {
      _discardTemporaryAuthState();
      return Result.failure(registrationResult.error!);
    }

    final registration = registrationResult.value!;
    final operationalUser = OperationalAuthState(
      accessToken: registration.accessToken,
      refreshToken: registration.refreshToken,
      userId: restrictedUser.userId!,
      email: restrictedUser.email,
      role: restrictedUser.role,
      customerId: restrictedUser.customerId,
    );

    return _startOperationalSession(operationalUser);
  }

  AsyncResult<OperationalAuthState> _startOperationalSession(
    OperationalAuthState user,
  ) async {
    _currentUser = user;

    await _storage.write(StorageKeys.accessToken, user.accessToken);
    await _storage.write(StorageKeys.refreshToken, user.refreshToken);

    final authSessionResult = await getAuthSession();
    if (authSessionResult.isFailure) {
      // If fetching the profile fails, we should log out to clear any partial state.
      await logout();
      return Result.failure(authSessionResult.error!);
    }

    _appSection.setAuthSession(authSessionResult.value!);
    final lastLoginIdentity = LastLoginIdentity(
      name: _appSection.customer?.name ?? '***',
      identifier: _appSection.user?.email ?? user.email,
    );
    final saveResult = await _lastLoginCacheService.save(lastLoginIdentity);

    if (saveResult.isFailure) {
      // TODO: Log the error but don't fail the login process since it's not
      //       critical. In a real app, consider using a logging service here.
      _log.warn('Failed to save last login identity: ${saveResult.error}');
    }

    return Success(user);
  }

  void _discardTemporaryAuthState() {
    if (_currentUser is RestrictedInstallationAuthState) {
      _currentUser = AnonymousAuthState();
    }
    _appSection.clear();
  }

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async {
    final result = await _lastLoginCacheService.get();
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }

  @override
  AsyncResult<Unit> logout() async {
    _appSection.clear();
    final wasLoggedIn = isLoggedIn;

    _currentUser = AnonymousAuthState();

    if (!wasLoggedIn) return Success(unit);

    await _storage.delete(StorageKeys.accessToken);
    await _storage.delete(StorageKeys.refreshToken);

    return Success(unit);
  }

  @override
  AsyncResult<AuthSession> getAuthSession() async {
    if (!isLoggedIn) {
      return Failure(
        AppError(
          code: AppErrorCode.unauthenticated,
          message: 'User is not logged in.',
        ),
      );
    }

    if (_appSection.isNotNull) return Success(_appSection.currentSession!);

    final result = await _api.getAuthSession();
    if (result.isFailure) return Result.failure(result.error!);

    _appSection.setAuthSession(result.value!);

    return Success(_appSection.currentSession!);
  }
}
