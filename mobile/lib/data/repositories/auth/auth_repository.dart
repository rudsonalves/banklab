import '/core/result/result.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/data/services/apis/auth/dtos/register_request_dto.dart';
import '/domain/common/auth/models/auth_user.dart';
import '/domain/common/auth/models/user_profile.dart';

abstract class AuthRepository {
  /// Returns the current authentication state for the app session.
  AuthUser get currentUser;

  /// Returns the cached user profile for the logged-in user, if one has already
  /// been loaded.
  UserProfile? get userProfile;

  /// Indicates whether the current user is authenticated.
  bool get isLoggedIn;

  /// Authenticates the user with the provided credentials.
  ///
  /// If a user is already logged in, the current logged user is returned
  /// without making a new API request.
  AsyncResult<LoggedUser> login(LoginRequestDto dto);

  /// Ends the current session and clears persisted auth tokens.
  ///
  /// If no user is logged in, it completes successfully without side effects.
  AsyncResult<Unit> logout();

  /// Registers a new user account.
  ///
  /// Returns a failure when called while there is an active logged-in session.
  AsyncResult<Unit> register(RegisterRequestDto dto);

  /// Returns the authenticated user's profile.
  ///
  /// If the user is not logged in, it returns an unauthenticated failure.
  /// When available, the cached profile is returned instead of fetching it
  /// again.
  AsyncResult<UserProfile> profile();
}
