import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/auth/cache/models/last_login_identity.dart';

class SplashViewmodel {
  final AuthRepository _authRepository;

  SplashViewmodel(this._authRepository) {
    initialize = Command0(_initialize);
  }

  late final Command0<LastLoginIdentity> initialize;

  AsyncResult<LastLoginIdentity> _initialize() async {
    return await _authRepository.getLastLoginIdentity();
  }
}
