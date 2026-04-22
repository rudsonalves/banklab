import 'dart:async';

import '/core/result/result.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/services/apis/account/balance_api.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';
import '/data/services/apis/account/dtos/statement_query_params_dto.dart';
import '/data/services/apis/account/dtos/statement_response_dto.dart';
import '/data/services/apis/account/statement_api.dart';

class AccountRepositoryImpl implements AccountRepository {
  final BalanceApi _balanceApi;
  final StatementApi _statementApi;

  AccountRepositoryImpl({
    required BalanceApi balanceApi,
    required StatementApi statementApi,
  }) : _balanceApi = balanceApi,
       _statementApi = statementApi;

  BalanceResponseDto? _balanceCache;
  final _balanceController = StreamController<BalanceResponseDto>.broadcast();

  @override
  BalanceResponseDto? getCachedBalance() {
    return _balanceCache;
  }

  @override
  Stream<BalanceResponseDto> watchBalance() async* {
    final cached = getCachedBalance();
    if (cached != null) {
      yield cached;
    }

    yield* _balanceController.stream;
  }

  @override
  AsyncResult<BalanceResponseDto> getBalance(String accountId) async {
    final result = await _balanceApi.getBalance(accountId);
    if (result.isFailure) return Result.failure(result.error!);

    final balance = result.value!;
    _balanceCache = balance;
    _balanceController.add(balance);

    return Success(balance);
  }

  @override
  AsyncResult<StatementResponseDto> getStatement(
    String accountId, {
    StatementQueryParamsDto queryParams = const StatementQueryParamsDto(),
  }) {
    return _statementApi.getStatement(accountId, queryParams: queryParams);
  }
}
