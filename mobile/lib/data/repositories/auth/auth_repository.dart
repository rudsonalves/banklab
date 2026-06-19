import '/core/result/result.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/data/services/cache/last_login/models/last_login_identity.dart';
import '/domain/common/auth/models/auth_session/auth_session.dart';
import '/domain/common/auth/models/auth_state.dart';

abstract class AuthRepository {
  /// Returns the current authentication state for the app session.
  AuthState get currentUser;

  /// Indicates whether the current user is authenticated.
  bool get isLoggedIn;

  /// Authenticates the user with the provided credentials.
  ///
  /// If a user is already logged in, the current logged user is returned
  /// without making a new API request.
  AsyncResult<AuthState> login(LoginRequestDto dto);

  /// Completes the restricted installation registration flow.
  ///
  /// Requires a previous restricted login result in [currentUser].
  AsyncResult<OperationalAuthState> certifyInstallation(
    String transactionPassword,
  );

  /// Ends the current session and clears persisted auth tokens.
  ///
  /// If no user is logged in, it completes successfully without side effects.
  AsyncResult<Unit> logout();

  /// Returns the authenticated user's profile.
  ///
  /// If the user is not logged in, it returns an unauthenticated failure.
  /// When available, the cached profile is returned instead of fetching it
  /// again.
  AsyncResult<AuthSession> getAuthSession();

  /// Retrieves the last login identity (name and identifier) from cache.
  /// Returns a failure if there is an error accessing the cache or if the
  /// data is incomplete.
  AsyncResult<LastLoginIdentity> getLastLoginIdentity();
}
