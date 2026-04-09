import '/core/result/result.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/data/services/apis/auth/dtos/register_request_dto.dart';
import '/domain/auth/models/auth_user.dart';
import '/domain/auth/models/user_profile.dart';

abstract class AuthRepository {
  AuthUser get currentUser;

  UserProfile? get userProfile;

  bool get isLoggedIn;

  AsyncResult<LoggedUser> login(LoginRequestDto dto);

  AsyncResult<Unit> logout();

  AsyncResult<Unit> register(RegisterRequestDto dto);

  AsyncResult<UserProfile> profile();
}
