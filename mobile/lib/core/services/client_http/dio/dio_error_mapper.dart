import 'package:dio/dio.dart';

import '/core/result/result.dart';

AppError mapHttpError(Object err, [StackTrace? stack]) {
  if (err is DioException) {
    final response = err.response;
    final data = response?.data;

    // --- network / timeout handling ---
    switch (err.type) {
      case DioExceptionType.connectionTimeout:
      case DioExceptionType.receiveTimeout:
      case DioExceptionType.sendTimeout:
        return AppError(
          statusCode: response?.statusCode,
          code: AppErrorCode.timeout,
          message: 'Connection timeout',
        );

      case DioExceptionType.connectionError:
        return AppError(
          code: AppErrorCode.networkError,
          message: 'No internet connection',
        );

      case DioExceptionType.cancel:
        return AppError(
          code: AppErrorCode.unexpected,
          message: 'Request cancelled',
        );

      default:
        break;
    }

    // --- API error parsing ---
    if (data is Map<String, dynamic>) {
      final error = data['error'];

      if (error is Map<String, dynamic>) {
        return AppError(
          statusCode: response?.statusCode,
          code: AppErrorCode.httpError,
          message: error['message'] ?? err.message ?? 'Request error',
          details: error['details'],
        );
      }

      if (data['message'] is String) {
        return AppError(
          statusCode: response?.statusCode,
          code: AppErrorCode.httpError,
          message: data['message'],
          details: data,
        );
      }
    }

    // --- fallback HTTP ---
    return AppError(
      statusCode: response?.statusCode,
      code: AppErrorCode.httpError,
      message: err.message ?? 'Request error',
      details: data,
    );
  }

  // --- unexpected error ---
  return AppError(
    code: AppErrorCode.unexpected,
    message: err.toString(),
    details: err,
  );
}
