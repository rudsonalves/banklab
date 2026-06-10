import '/core/result/command.dart';
import '/core/services/app_section/app_section.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/domain/common/auth/models/auth_user.dart';
import '../../models/post_login_destination.dart';

class ShortLoginViewModel {
  final AuthRepository _authRepository;
  final AppSection _appSection;

  ShortLoginViewModel({
    required AuthRepository authRepository,
    required AppSection appSection,
  }) : _authRepository = authRepository,
       _appSection = appSection {
    login = Command1(_authRepository.login);
  }

  late final Command1<LoggedUser, LoginRequestDto> login;

  PostLoginDestination resolvePostLoginDestination() {
    return PostLoginDestinationResolver.resolve(_appSection.currentSession);
  }
}
