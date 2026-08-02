import 'package:dio/dio.dart';

import '/core/result/result.dart';

const _accountApprovalRequiredBackendCode = 'ACCOUNT_APPROVAL_REQUIRED';
const _accountApprovalRequiredMessage =
    'Sua conta ainda está aguardando aprovação. Assim que ela for liberada,'
    ' você poderá acessar o app.';
const _contactNotVerifiedBackendCode = 'CONTACT_NOT_VERIFIED';
const _transactionPasswordLockedBackendCode = 'TRANSACTION_PASSWORD_LOCKED';
const _transactionPasswordNotSetBackendCode = 'TRANSACTION_PASSWORD_NOT_SET';
const _installationLimitReachedBackendCode = 'INSTALLATION_LIMIT_REACHED';
const _contactNotVerifiedGenericMessage =
    'Confirme seu e-mail e telefone antes de entrar.';
const _contactNotVerifiedEmailOnlyMessage =
    'Confirme seu e-mail antes de entrar.';
const _contactNotVerifiedPhoneOnlyMessage =
    'Confirme seu telefone antes de entrar.';

AppError mapHttpError(Object err, [StackTrace? stack]) {
  if (err is DioException) {
    final interceptorError = err.error;
    if (interceptorError is AppError) {
      return interceptorError;
    }

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
        final backendCode = error['code'];
        final details = error['details'];
        final appErrorDetails = _appErrorDetails(error);

        if (backendCode == _accountApprovalRequiredBackendCode) {
          return AppError(
            statusCode: response?.statusCode,
            code: AppErrorCode.accountApprovalRequired,
            message: _accountApprovalRequiredMessage,
            details: appErrorDetails,
          );
        }

        if (backendCode == _contactNotVerifiedBackendCode) {
          return AppError(
            statusCode: response?.statusCode,
            code: AppErrorCode.contactNotVerified,
            message: _contactNotVerifiedMessage(details),
            details: appErrorDetails,
          );
        }

        if (backendCode == _transactionPasswordLockedBackendCode) {
          return AppError(
            statusCode: response?.statusCode,
            code: AppErrorCode.transactionPasswordLocked,
            message: error['message'] ?? 'Senha transacional bloqueada.',
            details: appErrorDetails,
          );
        }

        if (backendCode == _transactionPasswordNotSetBackendCode) {
          return AppError(
            statusCode: response?.statusCode,
            code: AppErrorCode.transactionPasswordNotSet,
            message: error['message'] ?? 'Senha transacional não configurada.',
            details: appErrorDetails,
          );
        }

        if (backendCode == _installationLimitReachedBackendCode) {
          return AppError(
            statusCode: response?.statusCode,
            code: AppErrorCode.installationLimitReached,
            message: error['message'] ?? 'Limite de instalações atingido.',
            details: appErrorDetails,
          );
        }

        return AppError(
          statusCode: response?.statusCode,
          code: AppErrorCode.httpError,
          message: error['message'] ?? err.message ?? 'Request error',
          details: appErrorDetails,
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

Map<String, dynamic> _appErrorDetails(Map<String, dynamic> error) {
  final rawDetails = error['details'];
  if (rawDetails is Map) {
    return {
      ...rawDetails.map(
        (key, value) => MapEntry(key.toString(), value),
      ),
      if (error['code'] case final String code) 'code': code,
    };
  }

  return Map<String, dynamic>.from(error);
}

String _contactNotVerifiedMessage(Object? details) {
  if (details is! Map<String, dynamic>) {
    return _contactNotVerifiedGenericMessage;
  }

  final emailVerified = details['email_verified'];
  final phoneVerified = details['phone_verified'];

  if (emailVerified == false && phoneVerified == true) {
    return _contactNotVerifiedEmailOnlyMessage;
  }

  if (emailVerified == true && phoneVerified == false) {
    return _contactNotVerifiedPhoneOnlyMessage;
  }

  return _contactNotVerifiedGenericMessage;
}
