abstract class GenericException implements Exception {
  final String message;
  final Object? error;
  final StackTrace? stackTrace;

  const GenericException(
    this.message, {
    this.error,
    this.stackTrace,
  });
}
