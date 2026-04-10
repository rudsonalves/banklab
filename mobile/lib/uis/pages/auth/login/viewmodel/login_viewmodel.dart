import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/domain/auth/models/auth_user.dart';

class LoginViewModel {
  final AuthRepository _authRepository;

  LoginViewModel({
    required AuthRepository authRepository,
  }) : _authRepository = authRepository {
    login = Command1(_authRepository.login);
  }

  late final Command1<LoggedUser, LoginRequestDto> login;
}
