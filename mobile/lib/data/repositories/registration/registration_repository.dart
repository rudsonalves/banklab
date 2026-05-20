import '/core/result/result.dart';
import '/data/services/apis/registration/dtos/cpf_check_response_dto.dart';
import '/data/services/apis/registration/dtos/register_request_dto.dart';

abstract class RegistrationRepository {
  /// Registers a new user with the provided registration data.
  AsyncResult<Unit> register(RegisterRequestDto dto);

  /// Checks if a CPF is already associated with an existing account.
  AsyncResult<CpfCheckResponseDto> cpfCheck(String cpf);
}
