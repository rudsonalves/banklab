import 'dart:async';

import '/core/result/command.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/services/apis/account/dtos/account_summary_response_dto.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';

class HomeViewmodel {
  final AccountRepository _accountRepository;

  HomeViewmodel({
    required AccountRepository accountRepository,
  }) : _accountRepository = accountRepository {
    initialize = Command0(_initialize);
  }

  late final Command0<Unit> initialize;
  late final Command0<Unit> refreshBalance;

  BalanceResponseDto? get lastBalance => _accountRepository.lastBalance;

  AsyncResult<Unit> loadBalance() => _accountRepository.loadBalance();

  AccountSummaryResponseDto? get selectedAccount =>
      _accountRepository.selectedAccount;

  Future<Result<Unit>> _initialize() async {
    final result = await _accountRepository.listAccounts();
    if (result.isFailure) {
      return Failure(result.error!);
    }

    await loadBalance();
    return const Success(unit);
  }
}
