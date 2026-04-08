import 'api_error.dart';

export 'api_error.dart';

class ApiEnvelope<T> {
  final T? data;
  final ApiError? error;

  ApiEnvelope({
    required this.data,
    this.error,
  });

  factory ApiEnvelope.fromMap(
    Map<String, dynamic> map,
    T Function(Map<String, dynamic>) fromMap,
  ) {
    return ApiEnvelope(
      data: map['data'] != null
          ? fromMap(map['data'] as Map<String, dynamic>)
          : null,
      error: map['error'] != null
          ? ApiError.fromMap(map['error'] as Map<String, dynamic>)
          : null,
    );
  }
}
