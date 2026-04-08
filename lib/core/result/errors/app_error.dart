import 'app_error_code.dart';

export 'app_error_code.dart';

class AppError implements Exception {
  final int? statusCode;
  final AppErrorCode code;
  final String message;
  final Object? details;

  const AppError({
    this.statusCode,
    required this.code,
    required this.message,
    this.details,
  });

  @override
  String toString() => 'AppError($statusCode, ${code.name}, $message)';
}
