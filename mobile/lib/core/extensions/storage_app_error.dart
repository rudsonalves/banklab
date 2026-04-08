import '../result/errors/app_error.dart';

extension StorageAppError on AppError {
  static AppError storage({
    required String message,
    Object? details,
  }) {
    return AppError(
      code: AppErrorCode.storageError,
      message: message,
      details: details,
    );
  }

  static AppError notFound(String message) {
    return AppError(
      code: AppErrorCode.storageNotFound,
      message: 'Key not found: $message',
      details: 'Key "$message" was not found in storage.',
    );
  }
}
