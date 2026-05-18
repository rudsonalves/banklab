class ApiError {
  final String code;
  final String message;
  final Map<String, dynamic>? details;

  ApiError({
    required this.code,
    required this.message,
    this.details,
  });

  factory ApiError.fromMap(Map<String, dynamic> map) {
    final rawDetails = map['details'];

    return ApiError(
      code: map['code'] as String,
      message: map['message'] as String,
      details: rawDetails is Map
          ? rawDetails.map(
              (key, value) => MapEntry(key.toString(), value),
            )
          : null,
    );
  }
}
