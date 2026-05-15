import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/auth/api/dtos/login_request_dto.dart';
import '/domain/common/auth/models/auth_user.dart';

class ShortLoginViewModel {
  ShortLoginViewModel({
    required AuthRepository authRepository,
  }) {
    login = Command1(authRepository.login);
  }

  late final Command1<LoggedUser, LoginRequestDto> login;
}
