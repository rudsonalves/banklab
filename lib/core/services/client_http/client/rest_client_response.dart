class RestClientResponse {
  final dynamic data;
  final int? statusCode;
  final String? statusMessage;

  const RestClientResponse({
    this.data,
    this.statusCode,
    this.statusMessage,
  });

  bool get isSuccess =>
      statusCode != null && statusCode! >= 200 && statusCode! < 300;

  RestClientResponse copyWith({
    dynamic data,
    int? statusCode,
    String? statusMessage,
  }) {
    return RestClientResponse(
      data: data ?? this.data,
      statusCode: statusCode ?? this.statusCode,
      statusMessage: statusMessage ?? this.statusMessage,
    );
  }
}
