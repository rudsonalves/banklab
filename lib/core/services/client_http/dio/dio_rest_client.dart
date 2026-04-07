import 'package:dio/dio.dart';

import '/core/result/result.dart';
import '../client/rest_client.dart';
import '../client/rest_client_exception.dart';
import '../client/rest_client_request.dart';
import '../client/rest_client_response.dart';

class DioRestClient implements RestClient {
  final Dio _dio;

  DioRestClient({
    required Dio dio,
  }) : _dio = dio;

  @override
  void setBaseUrl(String url) {
    _dio.options.baseUrl = url;
  }

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

      return Success(
        RestClientResponse(
          data: response.data,
          statusCode: response.statusCode,
          statusMessage: response.statusMessage,
        ),
      );
    } on DioException catch (err) {
      final response = err.response;

      return Failure(
        RestClientException(
          message: err.message ?? 'HTTP error',
          statusCode: response?.statusCode,
          data: response?.data,
        ),
      );
    } catch (err) {
      return Failure(
        RestClientException(
          message: err.toString(),
        ),
      );
    }
  }
}
