import '/core/result/command.dart';
import '/core/services/app_section/app_section.dart';
import '/data/repositories/transaction_password/transaction_password_repository.dart';
import '/data/services/apis/transaction_password/dtos/step_up_authorize_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/step_up_authorize_response_dto.dart';

class VerifyTansactionPasswordViewmodel {
  final TransactionPasswordRepository _repository;

  VerifyTansactionPasswordViewmodel({
    required TransactionPasswordRepository repository,
    required AppSection appSection,
  }) : _repository = repository {
    stepUpAuthorize = Command1(_repository.stepUpAuthorize);
  }

  late final Command1<StepUpAuthorizeResponseDto, StepUpAuthorizeRequestDto>
  stepUpAuthorize;
}
