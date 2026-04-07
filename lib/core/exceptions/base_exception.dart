abstract class BaseException implements Exception {
  final dynamic data;
  final String message;
  final int? statusCode;
  final dynamic stackTracing;

  const BaseException({
    this.data,
    required this.message,
    this.statusCode,
    this.stackTracing,
  });
}
