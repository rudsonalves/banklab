import '/core/result/command.dart';
import '/core/services/installation_identity/installation_identity.dart';
import '/data/repositories/auth/auth_repository.dart';
import '../models/splash_bootstrap_state.dart';

class SplashViewmodel {
  final AuthRepository _authRepository;
  final InstallationIdentityService _installationIdentityService;

  SplashViewmodel({
    required AuthRepository authRepository,
    required InstallationIdentityService installationIdentityService,
  }) : _authRepository = authRepository,
       _installationIdentityService = installationIdentityService {
    initialize = Command0(_initialize);
  }

  late final Command0<SplashBootstrapState> initialize;

  AsyncResult<SplashBootstrapState> _initialize() async {
    final identityResult = await _installationIdentityService.resolve();
    if (identityResult.isFailure) {
      return Failure(identityResult.error!);
    }

    final lastLoginResult = await _authRepository.getLastLoginIdentity();
    if (lastLoginResult.isFailure) {
      return const Success(SplashBootstrapState());
    }

    return Success(
      SplashBootstrapState(lastLoginIdentity: lastLoginResult.value),
    );
  }
}
