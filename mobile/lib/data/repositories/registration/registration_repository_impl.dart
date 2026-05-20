import '/core/result/result.dart';
import '/data/services/apis/registration/dtos/cpf_check_response_dto.dart';
import '/data/services/apis/registration/dtos/register_request_dto.dart';
import '../../services/apis/registration/registration_api.dart';
import 'registration_repository.dart';

class RegistrationRepositoryImpl implements RegistrationRepository {
  final RegistrationApi _api;

  RegistrationRepositoryImpl({
    required RegistrationApi api,
  }) : _api = api;

  // final _log = ConsoleLog('RegistrationRepositoryImpl');

  @override
  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    final result = await _api.register(dto);
    if (result.isFailure) return Result.failure(result.error!);

    return Success(unit);
  }

  @override
  AsyncResult<CpfCheckResponseDto> cpfCheck(String cpf) async {
    final result = await _api.cpfCheck(cpf);
    if (result.isFailure) return Result.failure(result.error!);
    return Success(result.value!);
  }
}
