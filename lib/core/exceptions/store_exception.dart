import 'generic_exception.dart';

class StorageException extends GenericException {
  const StorageException(
    super.message, {
    super.error,
    super.stackTrace,
  });

  @override
  String toString() {
    return 'StorageException: $message,'
        ' error: $error,'
        ' stackTrace: $stackTrace';
  }
}

class StorageNotFoundException extends StorageException {
  const StorageNotFoundException(
    super.message, {
    super.error,
    super.stackTrace,
  });
}

class StorageConflictException extends StorageException {
  const StorageConflictException(
    super.message, {
    super.error,
    super.stackTrace,
  });
}

class StorageCorruptedException extends StorageException {
  const StorageCorruptedException(
    super.message, {
    super.error,
    super.stackTrace,
  });
}

class StorageExpiredException extends StorageException {
  const StorageExpiredException(
    super.message, {
    super.error,
    super.stackTrace,
  });
}

class StorageOperationException extends StorageException {
  const StorageOperationException(
    super.message, {
    super.error,
    super.stackTrace,
  });
}
