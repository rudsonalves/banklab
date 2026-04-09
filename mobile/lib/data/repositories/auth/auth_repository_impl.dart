import '/core/config/storage_keys.dart';
import '/core/result/result.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/auth_api.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/data/services/apis/auth/dtos/register_request_dto.dart';
import '/domain/auth/models/auth_user.dart';
import '/domain/auth/models/user_profile.dart';

class AuthRepositoryImpl implements AuthRepository {
  final AuthApi _api;
  final LocalSecureStorage _storage;

  AuthRepositoryImpl({
    required AuthApi api,
    required LocalSecureStorage storage,
  }) : _api = api,
       _storage = storage;

  AuthUser _currentUser = NotLoggedUser();
  UserProfile? _userProfile;

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
    // await _storage.write(StorageKeys.refreshToken, loggedUser.refreshToken);

    return Success(user);
  }

  @override
  AsyncResult<Unit> logout() async {
    if (!isLoggedIn) return Success(unit);

    _currentUser = NotLoggedUser();
    _userProfile = null;

    await _storage.delete(StorageKeys.accessToken);
    // await _storage.delete(StorageKeys.refreshToken);

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
}
