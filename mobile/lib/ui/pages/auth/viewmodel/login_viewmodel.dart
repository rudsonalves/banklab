import '/core/result/command.dart';
import '/core/services/app_section/app_section.dart';
import '/core/services/installation_identity/installation_identity.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/domain/common/auth/models/auth_user.dart';
import '../models/post_login_destination.dart';

class LoginViewModel {
  final AuthRepository _authRepository;
  final InstallationIdentityService _installationIdentityService;
  final AppSection _appSection;

  LoginViewModel({
    required AuthRepository authRepository,
    required InstallationIdentityService installationIdentityService,
    required AppSection appSection,
  }) : _authRepository = authRepository,
       _installationIdentityService = installationIdentityService,
       _appSection = appSection {
    login = Command1(_login);
  }

  late final Command1<LoggedUser, LoginRequestDto> login;

  AsyncResult<LoggedUser> _login(LoginRequestDto dto) async {
    final identityResult = await _installationIdentityService.resolve();
    if (identityResult.isFailure) {
      return Failure(identityResult.error!);
    }

    return _authRepository.login(dto);
  }

  PostLoginDestination resolvePostLoginDestination() {
    return PostLoginDestinationResolver.resolve(_appSection.currentSession);
  }
}
