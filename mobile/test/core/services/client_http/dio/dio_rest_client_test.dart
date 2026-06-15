import 'dart:convert';

import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_response.dart';
import 'package:bankflow/core/services/client_http/dio/dio_rest_client.dart';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('DioRestClient', () {
    test(
      'get should map request and return Success with response data',
      () async {
        late RequestOptions capturedOptions;

        final dio = Dio();
        dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
          capturedOptions = options;

          return ResponseBody.fromString(
            jsonEncode({'ok': true}),
            200,
            headers: {
              Headers.contentTypeHeader: [Headers.jsonContentType],
            },
            statusMessage: 'OK',
          );
        });

        final client = DioRestClient(dio: dio);

        final result = await client.get(
          const RestClientRequest(
            path: '/users',
            headers: {'Authorization': 'Bearer token'},
            queryParameters: {'page': 1},
          ),
        );

        expect(result, isA<Success<RestClientResponse>>());
        expect(capturedOptions.method, 'GET');
        expect(capturedOptions.path, '/users');
        expect(capturedOptions.queryParameters, {'page': 1});
        expect(capturedOptions.headers['Authorization'], 'Bearer token');
        expect(result.value?.statusCode, 200);
        expect(result.value?.statusMessage, 'OK');
        expect(result.value?.data, {'ok': true});
      },
    );

    test('post should forward body, query and headers', () async {
      late RequestOptions capturedOptions;

      final dio = Dio();
      dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
        capturedOptions = options;

        return ResponseBody.fromString(
          jsonEncode({'created': true}),
          201,
          headers: {
            Headers.contentTypeHeader: [Headers.jsonContentType],
          },
          statusMessage: 'Created',
        );
      });

      final client = DioRestClient(dio: dio);

      final result = await client.post(
        const RestClientRequest(
          path: '/accounts',
          headers: {'X-Tenant': 'bankflow'},
          queryParameters: {'sync': true},
          body: {'name': 'Main account'},
        ),
      );

      expect(result, isA<Success<RestClientResponse>>());
      expect(capturedOptions.method, 'POST');
      expect(capturedOptions.path, '/accounts');
      expect(capturedOptions.queryParameters, {'sync': true});
      expect(capturedOptions.headers['X-Tenant'], 'bankflow');
      expect(capturedOptions.data, {'name': 'Main account'});
      expect(result.value?.statusCode, 201);
      expect(result.value?.data, {'created': true});
    });

    test(
      'redacts credentials and does not log step-up authorization body',
      () async {
        final logs = <String>[];
        final originalDebugPrint = debugPrint;
        debugPrint = (message, {wrapWidth}) {
          if (message != null) logs.add(message);
        };
        addTearDown(() => debugPrint = originalDebugPrint);

        final dio = Dio();
        dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
          return ResponseBody.fromString(
            jsonEncode({'ok': true}),
            200,
            headers: {
              Headers.contentTypeHeader: [Headers.jsonContentType],
            },
          );
        });
        final client = DioRestClient(dio: dio);

        await client.post(
          const RestClientRequest(
            path: '/security/step-up/authorize',
            headers: {
              'Authorization': 'Bearer access-secret',
              'X-Step-Up-Token': 'step-up-secret',
              'X-App-Token': 'app-secret',
              'X-Trace-Id': 'trace-123',
            },
            body: {'transaction_password': '123456'},
          ),
        );

        final output = logs.join('\n');
        expect(output, contains('POST /security/step-up/authorize'));
        expect(output, contains('X-Trace-Id: trace-123'));
        expect(output, contains('Authorization: <redacted>'));
        expect(output, contains('X-Step-Up-Token: <redacted>'));
        expect(output, isNot(contains('access-secret')));
        expect(output, isNot(contains('step-up-secret')));
        expect(output, isNot(contains('app-secret')));
        expect(output, isNot(contains('transaction_password')));
        expect(output, isNot(contains('123456')));
      },
    );

    test('does not expose the step-up token in transfer logs', () async {
      final logs = <String>[];
      final originalDebugPrint = debugPrint;
      debugPrint = (message, {wrapWidth}) {
        if (message != null) logs.add(message);
      };
      addTearDown(() => debugPrint = originalDebugPrint);

      final dio = Dio();
      dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
        return ResponseBody.fromString(
          jsonEncode({'ok': true}),
          200,
          headers: {
            Headers.contentTypeHeader: [Headers.jsonContentType],
          },
        );
      });
      final client = DioRestClient(dio: dio);

      await client.post(
        const RestClientRequest(
          path: '/accounts/internal-transfers',
          headers: {'X-Step-Up-Token': 'single-use-step-up-secret'},
          body: {'amount': 2500},
        ),
      );

      final output = logs.join('\n');
      expect(output, contains('POST /accounts/internal-transfers'));
      expect(output, contains('X-Step-Up-Token: <redacted>'));
      expect(output, isNot(contains('single-use-step-up-secret')));
    });

    test('should use dio baseUrl configuration', () async {
      late RequestOptions capturedOptions;

      final dio = Dio(
        BaseOptions(baseUrl: 'https://api.bankflow.dev'),
      );
      dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
        capturedOptions = options;

        return ResponseBody.fromString(
          jsonEncode({'ok': true}),
          200,
          headers: {
            Headers.contentTypeHeader: [Headers.jsonContentType],
          },
        );
      });

      final client = DioRestClient(dio: dio);

      await client.get(const RestClientRequest(path: '/health'));

      expect(capturedOptions.baseUrl, 'https://api.bankflow.dev');
    });

    test(
      'should map DioException to Failure with AppError',
      () async {
        final dio = Dio();
        dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
          throw DioException(
            requestOptions: options,
            type: DioExceptionType.badResponse,
            message: 'Unauthorized',
            response: Response(
              requestOptions: options,
              statusCode: 401,
              data: {'error': 'invalid_token'},
            ),
          );
        });

        final client = DioRestClient(dio: dio);

        final result = await client.get(const RestClientRequest(path: '/me'));

        expect(result, isA<Failure<RestClientResponse>>());
        expect(result.error, isA<AppError>());

        final error = result.error!;
        expect(error.code, AppErrorCode.httpError);
        expect(error.message, 'Unauthorized');
        expect(error.statusCode, 401);
        expect(error.details, {'error': 'invalid_token'});
      },
    );

    test(
      'should map connectionError to network error',
      () async {
        final dio = Dio();
        dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
          throw DioException(
            requestOptions: options,
            type: DioExceptionType.connectionError,
          );
        });

        final client = DioRestClient(dio: dio);

        final result = await client.get(const RestClientRequest(path: '/me'));

        expect(result, isA<Failure<RestClientResponse>>());
        expect(result.error, isA<AppError>());

        final error = result.error!;
        expect(error.code, AppErrorCode.networkError);
        expect(error.statusCode, isNull);
        expect(error.message, 'No internet connection');
      },
    );
  });
}

class _FakeHttpClientAdapter implements HttpClientAdapter {
  _FakeHttpClientAdapter(this._handler);

  final Future<ResponseBody> Function(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  )
  _handler;

  @override
  void close({bool force = false}) {}

  @override
  Future<ResponseBody> fetch(
    RequestOptions options,
    Stream<Uint8List>? requestStream,
    Future<void>? cancelFuture,
  ) {
    return _handler(options, requestStream, cancelFuture);
  }
}
