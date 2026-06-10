import '/core/result/command.dart';
import '/core/services/app_section/app_section.dart';
import '/data/repositories/transaction_password/transaction_password_repository.dart';
import '/data/services/apis/transaction_password/dtos/create_transaction_password_request_dto.dart';
import '/data/services/apis/transaction_password/dtos/transaction_password_status_response_dto.dart';
import '/domain/common/auth/models/auth_session/auth_session.dart';

class TransactionPasswordViewModel {
  final TransactionPasswordRepository _repository;
  final AppSection _appSection;

  TransactionPasswordViewModel({
    required TransactionPasswordRepository repository,
    required AppSection appSection,
  }) : _repository = repository,
       _appSection = appSection {
    create = Command1(_repository.create);
  }

  late final Command1<
    TransactionPasswordStatusResponseDto,
    CreateTransactionPasswordRequestDto
  >
  create;

  AuthSession? get currentSession => _appSection.currentSession;

  bool get canAccessHome => _appSection.canAccessHome;
}
