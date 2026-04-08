import 'package:flutter/material.dart';

import 'result.dart';

export 'result.dart';

typedef CommandAction0<Output extends Object> =
    Future<Result<Output>> Function();

typedef CommandAction1<Output extends Object, Input> =
    Future<Result<Output>> Function(Input);

enum CommandState {
  idle,
  running,
  success,
  failure,
}

abstract class Command<Output extends Object> extends ChangeNotifier {
  Command();

  CommandState _state = CommandState.idle;
  Result<Output>? _result;

  int _executionId = 0;

  CommandState get state => _state;

  bool get isIdle => _state == CommandState.idle;
  bool get isRunning => _state == CommandState.running;
  bool get isSuccess => _state == CommandState.success;
  bool get isFailure => _state == CommandState.failure;

  Result<Output>? get result => _result;

  Output? get data => _result?.value;
  AppError? get error => _result?.error;

  Future<void> _execute(
    AsyncResult<Output> Function() action,
  ) async {
    final currentExecution = ++_executionId;

    _state = CommandState.running;
    notifyListeners();

    try {
      final result = await action();

      if (currentExecution != _executionId) return;

      _result = result;
      _state = result.isSuccess ? CommandState.success : CommandState.failure;
    } catch (err) {
      // CORE fallback (infra-agnostic)
      _result = Result.failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: err.toString(),
          details: err,
        ),
      );
      _state = CommandState.failure;
    } finally {
      notifyListeners();
    }
  }
}

class Command0<Output extends Object> extends Command<Output> {
  final CommandAction0<Output> _action;

  Command0(this._action);

  Future<void> execute() async {
    if (isRunning) return;
    await _execute(_action);
  }
}

class Command1<Output extends Object, Input> extends Command<Output> {
  final CommandAction1<Output, Input> _action;

  Command1(this._action);

  Future<void> execute(Input param) async {
    if (isRunning) return;
    await _execute(() => _action(param));
  }
}
