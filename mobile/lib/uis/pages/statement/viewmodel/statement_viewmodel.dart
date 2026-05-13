import '/core/result/command.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';

class StatementViewmodel {
  final AccountRepository _accountRepository;

  StatementViewmodel(this._accountRepository) {
    getStatement = Command1(_accountRepository.getStatement);
    selectAccount = Command1(_accountRepository.selectAccount);
  }

  late final Command1<StatementResponseDto, StatementQueryParamsDto>
  getStatement;
  late final Command1<void, String> selectAccount;

  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepository.selectedAccount;

  StatementResponseDto? get lastStatement => _accountRepository.lastStatement;

  List<AccountSummaryResponseDto>? get accounts => _accountRepository.accounts;
}
