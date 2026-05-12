import '/core/result/result.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';

abstract class AccountRepository {
  /// Returns the last fetched balance from the cache. If no balance has been
  /// fetched yet, it returns null.
  BalanceResponseDto? get lastBalance;

  /// Returns the currently selected account. If no account is selected, it
  /// returns null.
  AccountSummaryResponseDto? get selectedAccount;

  /// Returns the last fetched statement from the cache. If no statement has
  /// been fetched yet, it returns null.
  StatementResponseDto? get lastStatement;

  /// Returns the list of accounts from the cache. If the accounts have not been
  /// fetched yet, it returns null.
  List<AccountSummaryResponseDto>? get accounts;

  /// Provides a stream of balance updates. Whenever the balance is updated, the
  /// new balance is emitted to the stream.
  Stream<BalanceResponseDto> balance();

  /// Loads the balance for the currently selected account. If no account
  /// is selected, it returns a failure result.
  AsyncResult<Unit> loadBalance();

  /// Fetches the list of accounts from the API. If successful, it caches
  /// the accounts and selects the first account by default. If the API
  /// call fails, it returns a failure result.
  AsyncResult<List<AccountSummaryResponseDto>> listAccounts();

  /// Selects an account by its ID.
  /// - If the account is found in the cache, it sets it as the selected account
  ///   and loads its balance.
  /// - If the account is not found or if the accounts have not been loaded yet,
  ///   it returns a failure result.
  AsyncResult<Unit> selectAccount(String accountId);

  /// Retrieves the account statement based on the provided query parameters.
  /// If no account is selected, it returns a failure result.
  AsyncResult<StatementResponseDto> getStatement(
    StatementQueryParamsDto queryParams,
  );
}
