import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/apis/auth/dtos/register_request_dto.dart';

class RegisterViewmodel {
  final AuthRepository _authRepository;

  RegisterViewmodel({
    required AuthRepository authRepository,
  }) : _authRepository = authRepository {
    register = Command1(_authRepository.register);
  }

  late final Command1<Unit, RegisterRequestDto> register;
}
