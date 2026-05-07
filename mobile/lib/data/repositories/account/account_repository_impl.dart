import 'dart:async';

import '/core/result/result.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/services/apis/account/balance_api.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';
import '/data/services/apis/account/list_accounts_api.dart';
import '/data/services/apis/account/statement_api.dart';

class AccountRepositoryImpl implements AccountRepository {
  final BalanceApi _balanceApi;
  final ListAccountsApi _listAccountsApi;
  final StatementApi _statementApi;

  AccountRepositoryImpl({
    required BalanceApi balanceApi,
    required ListAccountsApi listAccountsApi,
    required StatementApi statementApi,
  }) : _balanceApi = balanceApi,
       _listAccountsApi = listAccountsApi,
       _statementApi = statementApi;

  BalanceResponseDto? _balanceCache;
  final _balanceController = StreamController<BalanceResponseDto>.broadcast();

  AccountSummaryResponseDto? _selectedAccount;

  List<AccountSummaryResponseDto>? _accountsCache;

  @override
  BalanceResponseDto? get lastBalance => _balanceCache;

  @override
  AccountSummaryResponseDto? get selectedAccount => _selectedAccount;

  @override
  List<AccountSummaryResponseDto>? get accounts => _accountsCache;

  @override
  Stream<BalanceResponseDto> balance() => _balanceController.stream;

  @override
  AsyncResult<Unit> loadBalance() async {
    if (_selectedAccount == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No account selected.',
        ),
      );
    }

    final result = await _balanceApi.getBalance(_selectedAccount!.id);
    if (result.isSuccess) {
      _balanceCache = result.value;
      _balanceController.add(_balanceCache!);
      return const Success(unit);
    } else {
      return Failure(result.error!);
    }
  }

  @override
  AsyncResult<List<AccountSummaryResponseDto>> listAccounts() async {
    final result = await _listAccountsApi.listAccounts();
    if (result.isFailure) return Result.failure(result.error!);

    final accounts = result.value!;

    if (accounts.isNotEmpty) {
      _accountsCache = accounts;
      _selectedAccount = accounts[0];
      await loadBalance();
    }

    return Success(accounts);
  }

  @override
  void selectAccount(AccountSummaryResponseDto account) {
    _selectedAccount = account;
    loadBalance();
  }

  @override
  void clearSelectedAccount() {
    _selectedAccount = null;
  }

  @override
  AsyncResult<StatementResponseDto> getStatement(
    StatementQueryParamsDto queryParams,
  ) async {
    if (_selectedAccount == null) {
      return Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No account selected.',
        ),
      );
    }

    return _statementApi.getStatement(
      _selectedAccount!.id,
      queryParams: queryParams,
    );
  }
}
