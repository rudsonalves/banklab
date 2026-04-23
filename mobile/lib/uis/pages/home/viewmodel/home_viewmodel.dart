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

    startTimer();
    return const Success(unit);
  }

  Timer? _timer;
  final Duration _delay = const Duration(seconds: 20);

  void startTimer() {
    stopTimer();

    _timer = Timer.periodic(_delay, (_) => loadBalance());
  }

  void stopTimer() {
    _timer?.cancel();
    _timer = null;
  }

  void dispose() {
    stopTimer();
  }
}
