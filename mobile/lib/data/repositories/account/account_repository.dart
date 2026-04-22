import '/core/result/result.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';

abstract class AccountRepository {
  BalanceResponseDto? getCachedBalance(String accountId);

  Stream<BalanceResponseDto> watchBalance(String accountId);

  AsyncResult<BalanceResponseDto> getBalance(String accountId);

  AsyncResult<StatementResponseDto> getStatement(
    String accountId, {
    StatementQueryParamsDto queryParams,
  });
}
