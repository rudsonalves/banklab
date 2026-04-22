import 'dart:async';

import 'package:flutter/material.dart';

import '/core/result/result.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/services/apis/account/dtos/balance_response_dto.dart';

class HomeViewmodel extends ChangeNotifier {
  final AccountRepository _accountRepository;

  HomeViewmodel({
    required AccountRepository accountRepository,
  }) : _accountRepository = accountRepository;

  StreamSubscription<BalanceResponseDto>? _balanceSubscription;
  Timer? _refreshTimer;

  BalanceResponseDto? _balance;
  AppError? _balanceError;
  bool _isLoading = false;
  DateTime? _lastUpdatedAt;
  String? _accountId;

  BalanceResponseDto? get balance => _balance;
  AppError? get balanceError => _balanceError;
  bool get isLoading => _isLoading;
  DateTime? get lastUpdatedAt => _lastUpdatedAt;
  String? get accountId => _accountId;

  bool get hasAccountId => _accountId != null && _accountId!.isNotEmpty;

  Future<void> initialize({String? accountId}) async {
    final normalizedAccountId = _normalizeAccountId(accountId);
    if (_accountId == normalizedAccountId && _refreshTimer != null) return;

    await _balanceSubscription?.cancel();
    _refreshTimer?.cancel();

    _accountId = normalizedAccountId;
    _balance = null;
    _balanceError = null;
    _lastUpdatedAt = null;

    if (!hasAccountId) {
      notifyListeners();
      return;
    }

    _balance = _accountRepository.getCachedBalance();
    _balanceSubscription = _accountRepository.watchBalance().listen(
      (balance) {
        _balance = balance;
        _lastUpdatedAt = DateTime.now();
        _balanceError = null;
        notifyListeners();
      },
    );

    notifyListeners();
    await refreshBalance();

    _refreshTimer = Timer.periodic(
      const Duration(seconds: 10),
      (_) => refreshBalance(),
    );
  }

  Future<void> refreshBalance() async {
    if (!hasAccountId || _isLoading) return;

    _isLoading = true;
    notifyListeners();

    final result = await _accountRepository.getBalance(_accountId!);
    if (result.isFailure) {
      _balanceError = result.error!;
    } else {
      _balance = result.value!;
      _lastUpdatedAt = DateTime.now();
      _balanceError = null;
    }

    _isLoading = false;
    notifyListeners();
  }

  String? _normalizeAccountId(String? accountId) {
    final normalized = accountId?.trim();
    if (normalized == null || normalized.isEmpty) return null;
    return normalized;
  }

  @override
  void dispose() {
    _balanceSubscription?.cancel();
    _refreshTimer?.cancel();
    super.dispose();
  }
}
