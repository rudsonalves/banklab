class RestClientException implements Exception {
  final String message;
  final int? statusCode;
  final dynamic data;

  const RestClientException({
    required this.message,
    this.statusCode,
    this.data,
  });

  bool get isUnauthorized => statusCode == 401;
  bool get isForbidden => statusCode == 403;
  bool get isNotFound => statusCode == 404;
  bool get isServerError => statusCode != null && statusCode! >= 500;

  @override
  String toString() {
    return 'RestClientException(statusCode: $statusCode, message: $message)';
  }
}
