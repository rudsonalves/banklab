import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/domain/common/auth/models/auth_user.dart';
import '../../../../../data/services/auth/api/dtos/login_request_dto.dart';

class LoginViewModel {
  final AuthRepository _authRepository;

  LoginViewModel({
    required AuthRepository authRepository,
  }) : _authRepository = authRepository {
    login = Command1(_authRepository.login);
  }

  late final Command1<LoggedUser, LoginRequestDto> login;
}
