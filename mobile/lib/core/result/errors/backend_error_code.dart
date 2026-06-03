import 'app_error.dart';

String? backendErrorCode(AppError? error) {
  if (error == null) return null;

  return _readCode(error.details);
}

String? _readCode(Object? value) {
  if (value is! Map) return null;

  final code = value['code'];
  if (code is String && code.trim().isNotEmpty) {
    return code;
  }

  final error = value['error'];
  final errorCode = _readCode(error);
  if (errorCode != null) return errorCode;

  final details = value['details'];
  return _readCode(details);
}
