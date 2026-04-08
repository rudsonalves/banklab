class ApiError {
  final String code;
  final String message;

  ApiError({
    required this.code,
    required this.message,
  });

  factory ApiError.fromMap(Map<String, dynamic> map) {
    return ApiError(
      code: map['code'] as String,
      message: map['message'] as String,
    );
  }
}
