class RestClientRequest {
  final String path;
  final Map<String, dynamic>? headers;
  final Map<String, dynamic>? queryParameters;
  final dynamic body;

  const RestClientRequest({
    required this.path,
    this.headers,
    this.queryParameters,
    this.body,
  });

  RestClientRequest copyWith({
    String? path,
    Map<String, dynamic>? headers,
    Map<String, dynamic>? queryParameters,
    dynamic body,
  }) {
    return RestClientRequest(
      path: path ?? this.path,
      headers: headers ?? this.headers,
      queryParameters: queryParameters ?? this.queryParameters,
      body: body ?? this.body,
    );
  }
}
