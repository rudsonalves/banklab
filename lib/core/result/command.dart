import 'package:flutter/material.dart';

import 'result.dart';

typedef CommandAction0<Output extends Object> =
    Future<Result<Output>> Function();
typedef CommandAction1<Output extends Object, Input> =
    Future<Result<Output>> Function(Input);

abstract interface class Command<Output extends Object> extends ChangeNotifier {
  Command();

  bool _running = false;

  Result<Output>? _result;

  bool get isRunning => _running;
  bool get isSuccess => _result?.isSuccess ?? false;
  bool get isFailure => _result?.isFailure ?? false;

  Result<Output>? get value => _result;

  Future<void> _execute(CommandAction0<Output> action) async {
    if (_running) return;

    _running = true;
    _result = null;
    notifyListeners();

    try {
      _result = await action();
    } finally {
      _running = false;
      notifyListeners();
    }
  }
}

class Command0<Output extends Object> extends Command<Output> {
  final CommandAction0<Output> _action;

  Command0(this._action);

  Future<void> execute() async {
    await _execute(_action);
  }
}

class Command1<Output extends Object, Input> extends Command<Output> {
  final CommandAction1<Output, Input> _action;

  Command1(this._action);

  Future<void> execute(Input param) async {
    await _execute(() => _action(param));
  }
}
