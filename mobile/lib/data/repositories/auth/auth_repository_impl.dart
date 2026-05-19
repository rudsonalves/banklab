import '/core/resources/storage_keys.dart';
import '/core/result/result.dart';
import '/core/services/logging/console_log.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/auth/api/auth_api.dart';
import '/data/services/auth/api/dtos/contact_verification_confirm_request_dto.dart';
import '/data/services/auth/api/dtos/contact_verification_confirm_response_dto.dart';
import '/data/services/auth/api/dtos/contact_verification_request_dto.dart';
import '/data/services/auth/api/dtos/contact_verification_request_response_dto.dart';
import '/data/services/auth/api/dtos/login_request_dto.dart';
import '/data/services/auth/api/dtos/register_request_dto.dart';
import '/data/services/auth/cache/last_login_cache_service.dart';
import '/data/services/auth/cache/models/last_login_identity.dart';
import '/domain/common/auth/models/auth_user.dart';
import '/domain/common/auth/models/user_profile.dart';
import '../../services/auth/api/dtos/cpf_check_response_dto.dart';

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
  UserProfile? _userProfile;

  final _log = ConsoleLog('AuthRepositoryImpl');

  @override
  AuthUser get currentUser => _currentUser;

  @override
  UserProfile? get userProfile => _userProfile;

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

    _userProfile = profileResult.value!;
    final saveResult = await _lastLoginCacheService.save(
      LastLoginIdentity(
        name: _userProfile!.name,
        identifier: _userProfile!.email,
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
    _userProfile = null;

    await _storage.delete(StorageKeys.accessToken);
    await _storage.delete(StorageKeys.refreshToken);

    return Success(unit);
  }

  @override
  AsyncResult<UserProfile> profile() async {
    if (!isLoggedIn) {
      return Failure(
        AppError(
          code: AppErrorCode.unauthenticated,
          message: 'User is not logged in.',
        ),
      );
    }

    if (_userProfile != null) return Success(_userProfile!);

    final result = await _api.getProfile();
    if (result.isFailure) return Result.failure(result.error!);

    _userProfile = result.value!;

    return Success(_userProfile!);
  }

  @override
  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    if (isLoggedIn) {
      return Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'User is already logged in.',
        ),
      );
    }

    final result = await _api.register(dto);
    if (result.isFailure) return Result.failure(result.error!);

    return Success(unit);
  }

  AsyncResult<CpfCheckResponseDto> cpfCheck(String cpf) async {
    final result = await _api.cpfCheck(cpf);
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    final result = await _api.requestContactVerification(dto);
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    final result = await _api.confirmContactVerification(dto);
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }
}
