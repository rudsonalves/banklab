import 'package:dio/dio.dart';

import '/core/result/result.dart';
import '../../logging/console_log.dart';
import '../client/rest_client.dart';
import '../client/rest_client_request.dart';
import '../client/rest_client_response.dart';
import 'dio_error_mapper.dart';

class DioRestClient implements RestClient {
  final Dio _dio;

  DioRestClient({
    required Dio dio,
  }) : _dio = dio;

  final _log = ConsoleLog('DioRestClient');

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) {
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
    _log.info('POST ${request.path} - Headers: ${request.headers}');
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
    _log.info('PUT ${request.path} - Headers: ${request.headers}');
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
    return _request(
      () => _dio.delete(
        request.path,
        data: request.body,
        queryParameters: request.queryParameters,
        options: Options(headers: request.headers),
      ),
    );
  }

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
      _log.error('Request error: $err', error: err, stack: stack);
      return Result.failure(mapHttpError(err, stack));
    }
  }
}
