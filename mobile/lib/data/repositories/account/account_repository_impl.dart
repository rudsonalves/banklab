import 'dart:async';

import '/core/result/result.dart';
import '/data/repositories/account/account_repository.dart';
import '../../services/apis/account/balance_api.dart';
import '../../services/apis/account/dtos/account_summary_response_dto.dart';
import '../../services/apis/account/dtos/balance_response_dto.dart';
import '../../services/apis/account/dtos/statement_query_params_dto.dart';
import '../../services/apis/account/dtos/statement_response_dto.dart';
import '../../services/apis/account/list_accounts_api.dart';
import '../../services/apis/account/statement_api.dart';

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
  StatementResponseDto? _statementCache;

  @override
  BalanceResponseDto? get lastBalance => _balanceCache;

  @override
  AccountSummaryResponseDto? get selectedAccount => _selectedAccount;

  @override
  StatementResponseDto? get lastStatement => _statementCache;

  @override
  List<AccountSummaryResponseDto>? get accounts => _accountsCache;

  @override
  Stream<BalanceResponseDto> balance() async* {
    if (_balanceCache != null) {
      yield _balanceCache!;
    }

    yield* _balanceController.stream;
  }

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
  AsyncResult<Unit> selectAccount(String accountId) async {
    if (_accountsCache == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No accounts available.',
        ),
      );
    }

    for (final account in _accountsCache!) {
      if (account.id == accountId) {
        _selectedAccount = account;
        _statementCache = null;

        await loadBalance();
        return const Success(unit);
      }
    }

    return const Failure(
      AppError(
        code: AppErrorCode.unexpected,
        message: 'Account not found.',
      ),
    );
  }

  @override
  AsyncResult<StatementResponseDto> getStatement(
    StatementQueryParamsDto queryParams,
  ) async {
    _statementCache = null;
    if (_selectedAccount == null) {
      return Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'No account selected.',
        ),
      );
    }

    final result = await _statementApi.getStatement(
      _selectedAccount!.id,
      queryParams: queryParams,
    );

    if (result.isFailure) return Result.failure(result.error!);

    _statementCache = result.value!;
    return Success(_statementCache!);
  }
}
