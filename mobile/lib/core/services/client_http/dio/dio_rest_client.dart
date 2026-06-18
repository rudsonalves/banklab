import 'package:dio/dio.dart';

import '/core/resources/app_http_headers.dart';
import '/core/result/result.dart';
import '../../logging/console_log.dart';
import '../client/rest_client.dart';
import '../client/rest_client_request.dart';
import '../client/rest_client_response.dart';
import 'dio_error_mapper.dart';

class DioRestClient implements RestClient {
  static const _redactedHeaderValue = '<redacted>';

  final Dio _dio;

  DioRestClient({
    required Dio dio,
  }) : _dio = dio;

  final _log = ConsoleLog('DioRestClient');

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) {
    _logRequest('GET', request);
    return _request(
      () => _dio.get(
        request.path,
        queryParameters: request.queryParameters,
        options: Options(headers: request.headers),
      ),
    );
  }

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) {
    _logRequest('POST', request);
    return _request(
      () => _dio.post(
        request.path,
        data: request.body,
        queryParameters: request.queryParameters,
        options: Options(headers: request.headers),
      ),
    );
  }

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) {
    _logRequest('PUT', request);
    return _request(
      () => _dio.put(
        request.path,
        data: request.body,
        queryParameters: request.queryParameters,
        options: Options(headers: request.headers),
      ),
    );
  }

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) {
    _logRequest('PATCH', request);
    return _request(
      () => _dio.patch(
        request.path,
        data: request.body,
        queryParameters: request.queryParameters,
        options: Options(headers: request.headers),
      ),
    );
  }

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) {
    _logRequest('DELETE', request);
    return _request(
      () => _dio.delete(
        request.path,
        data: request.body,
        queryParameters: request.queryParameters,
        options: Options(headers: request.headers),
      ),
    );
  }

  /// Logs the HTTP request method, path, and headers, redacting sensitive
  /// header values.
  void _logRequest(String method, RestClientRequest request) {
    final headers = request.headers;
    if (headers == null || headers.isEmpty) {
      _log.info('$method ${request.path}');
      return;
    }

    final sanitizedHeaders = headers.map(
      (name, value) => MapEntry(
        name,
        AppHttpHeaders.sensitiveLowercase.contains(name.toLowerCase())
            ? _redactedHeaderValue
            : value,
      ),
    );
    _log.info('$method ${request.path} - Headers: $sanitizedHeaders');
  }

  /// Executes the given HTTP call and maps the response to a [RestClientResponse].
  /// If an error occurs, it maps the error to an [AppError] and logs
  /// it appropriately.
  AsyncResult<RestClientResponse> _request(
    Future<Response> Function() call,
  ) async {
    try {
      final response = await call();

      return Result.success(
        RestClientResponse(
          data: response.data,
          statusCode: response.statusCode,
          statusMessage: response.statusMessage,
        ),
      );
    } catch (err, stack) {
      final error = mapHttpError(err, stack);

      if (_shouldLogAsTechnicalError(err, error)) {
        _log.error('Request error: ${error.message}', error: err, stack: stack);
      } else {
        _log.warn(
          'Request failed (${error.statusCode ?? '-'}/${error.code.name}): '
          '${error.message}',
        );
      }

      return Result.failure(error);
    }
  }

  /// Determines whether the given error should be logged as a technical error
  /// based on its type and status code.
  bool _shouldLogAsTechnicalError(Object err, AppError error) {
    if (err is DioException) {
      final statusCode = err.response?.statusCode;

      if (statusCode != null) {
        return statusCode >= 500;
      }

      return switch (err.type) {
        DioExceptionType.connectionTimeout ||
        DioExceptionType.sendTimeout ||
        DioExceptionType.receiveTimeout ||
        DioExceptionType.connectionError ||
        DioExceptionType.cancel => false,
        _ => true,
      };
    }

    return error.code == AppErrorCode.unexpected;
  }
}
