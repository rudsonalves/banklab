import 'dart:developer' as developer;

import 'package:flutter/foundation.dart';

enum LogLevel {
  error,
  warn,
  info,
}

class ConsoleLog {
  final String context;

  const ConsoleLog(this.context);

  static const bool _enabled = kDebugMode;

  String _context([String? label]) =>
      label == null || label.isEmpty ? context : '$context.$label';

  void _separator() {
    if (!_enabled) return;
    debugPrint('--------------------------------------------------');
  }

  void _printBlock(String header, String message) {
    if (!_enabled) return;
    _separator();
    debugPrint(header);
    debugPrint(message);
    _separator();
  }

  void error(
    String message, {
    String? label,
    Object? error,
    StackTrace? stack,
  }) {
    if (!_enabled) return;

    final ctx = _context(label);

    _printBlock(
      '[ERROR][$ctx]',
      message,
    );

    if (error != null) {
      debugPrint('[ERROR][$ctx][exception]');
      debugPrint(error.toString());
    }

    if (stack != null) {
      debugPrint('[ERROR][$ctx][stacktrace]');
      debugPrint(stack.toString());
    }
  }

  void warn(String message, {String? label}) {
    if (!_enabled) return;

    final ctx = _context(label);

    _printBlock(
      '[WARN][$ctx]',
      message,
    );
  }

  void info(String message, {String? label}) {
    if (!_enabled) return;

    final ctx = _context(label);

    _printBlock(
      '[INFO][$ctx]',
      message,
    );
  }

  void log(String message, {String? label}) {
    if (!_enabled) return;

    final ctx = _context(label);

    developer.log('[$ctx]');
    developer.log('[$ctx] $message');
    developer.log('[$ctx]');
  }
}
