import '/core/result/result.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';

abstract class AccountRepository {
  BalanceResponseDto? get lastBalance;

  AccountSummaryResponseDto? get selectedAccount;

  Stream<BalanceResponseDto> balance();

  AsyncResult<Unit> loadBalance();

  AsyncResult<List<AccountSummaryResponseDto>> listAccounts();

  void selectAccount(AccountSummaryResponseDto account);

  void clearSelectedAccount();

  AsyncResult<StatementResponseDto> getStatement(
    StatementQueryParamsDto queryParams,
  );
}
