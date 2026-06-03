import '/core/resources/storage_keys.dart';
import '/core/result/result.dart';
import '/core/services/logging/console_log.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/auth_api.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/data/services/cache/last_login/last_login_cache_service.dart';
import '/data/services/cache/last_login/models/last_login_identity.dart';
import '/domain/common/auth/models/auth_session/auth_session.dart';
import '/domain/common/auth/models/auth_user.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthApi _api;
  final LocalSecureStorage _storage;
  final LastLoginCacheService _lastLoginCacheService;

  AuthRepositoryImpl({
    required AuthApi api,
    required LocalSecureStorage storage,
    required LastLoginCacheService lastLoginCacheService,
  }) : _api = api,
       _storage = storage,
       _lastLoginCacheService = lastLoginCacheService;

  AuthUser _currentUser = NotLoggedUser();
  AuthSession? _authSession;

  final _log = ConsoleLog('AuthRepositoryImpl');

  @override
  AuthUser get currentUser => _currentUser;

  @override
  AuthSession? get userProfile => _authSession;

  @override
  bool get isLoggedIn => _currentUser is LoggedUser;

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    if (isLoggedIn) return Success(_currentUser as LoggedUser);

    final result = await _api.login(dto);
    if (result.isFailure) return Result.failure(result.error!);

    final user = result.value!;
    _currentUser = user;

    await _storage.write(StorageKeys.accessToken, user.accessToken);
    await _storage.write(StorageKeys.refreshToken, user.refreshToken);

    final profileResult = await profile();
    if (profileResult.isFailure) {
      // If fetching the profile fails, we should log out to clear any partial state.
      await logout();
      return Result.failure(profileResult.error!);
    }

    _authSession = profileResult.value!;
    final saveResult = await _lastLoginCacheService.save(
      LastLoginIdentity(
        name: _authSession!.customer?.name ?? '***',
        identifier: _authSession!.customer?.cpf ?? '***',
      ),
    );

    if (saveResult.isFailure) {
      // TODO: Log the error but don't fail the login process since it's not
      //       critical. In a real app, consider using a logging service here.
      _log.warn('Failed to save last login identity: ${saveResult.error}');
    }

    return Success(user);
  }

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async {
    final result = await _lastLoginCacheService.get();
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }

  @override
  AsyncResult<Unit> logout() async {
    if (!isLoggedIn) return Success(unit);

    _currentUser = NotLoggedUser();
    _authSession = null;

    await _storage.delete(StorageKeys.accessToken);
    await _storage.delete(StorageKeys.refreshToken);

    return Success(unit);
  }

  @override
  AsyncResult<AuthSession> profile() async {
    if (!isLoggedIn) {
      return Failure(
        AppError(
          code: AppErrorCode.unauthenticated,
          message: 'User is not logged in.',
        ),
      );
    }

    if (_authSession != null) return Success(_authSession!);

    final result = await _api.getAuthSession();
    if (result.isFailure) return Result.failure(result.error!);

    _authSession = result.value!;

    return Success(_authSession!);
  }
}
