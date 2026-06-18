import 'dart:convert';
import 'dart:typed_data';

import 'package:bankflow/core/resources/app_http_headers.dart';
import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/dio/dio_rest_client.dart';
import 'package:bankflow/core/services/client_http/interceptors/auth/auth_interceptor.dart';
import 'package:bankflow/core/services/client_http/interceptors/installation/installation_interceptor.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:bankflow/core/services/installation_identity/installation_identity.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('InstallationInterceptor', () {
    test('adds installation id header to login requests', () async {
      const installationId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';
      late RequestOptions capturedOptions;
      final dio = Dio(BaseOptions(baseUrl: 'https://api.test'))
        ..interceptors.add(
          InstallationInterceptor(
            installationIdentityService: _FakeInstallationIdentityService(
              result: const Success(installationId),
            ),
          ),
        );

      dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
        capturedOptions = options;
        return _jsonResponse(200, {
          'data': {'ok': true},
        });
      });

      await dio.post('/auth/login');

      expect(
        capturedOptions.headers[InstallationInterceptor.headerName],
        installationId,
      );
    });

    test('adds installation id header without changing auth header', () async {
      const installationId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';
      late RequestOptions capturedOptions;
      final storage = _MemorySecureStorage({
        StorageKeys.accessToken: 'access-token',
      });
      final dio = Dio(BaseOptions(baseUrl: 'https://api.test'));

      dio.interceptors.add(
        InstallationInterceptor(
          installationIdentityService: _FakeInstallationIdentityService(
            result: const Success(installationId),
          ),
        ),
      );
      dio.interceptors.add(
        AuthInterceptor(
          authDio: dio,
          secureStorage: storage,
          baseUrl: 'https://api.test',
        ),
      );

      dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
        capturedOptions = options;
        return _jsonResponse(200, {
          'data': {'ok': true},
        });
      });

      await dio.get('/accounts');

      expect(
        capturedOptions.headers[InstallationInterceptor.headerName],
        installationId,
      );
      expect(
        capturedOptions.headers[AppHttpHeaders.authorization],
        'Bearer access-token',
      );
    });

    test(
      'blocks request when installation identity resolution fails',
      () async {
        var adapterCalls = 0;
        final dio = Dio(BaseOptions(baseUrl: 'https://api.test'))
          ..interceptors.add(
            InstallationInterceptor(
              installationIdentityService: _FakeInstallationIdentityService(
                result: const Failure(
                  AppError(
                    code: AppErrorCode.storageError,
                    message: 'identity failed',
                  ),
                ),
              ),
            ),
          );

        dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
          adapterCalls++;
          return _jsonResponse(200, {
            'data': {'ok': true},
          });
        });

        await expectLater(
          dio.get('/accounts'),
          throwsA(
            isA<DioException>()
                .having((error) => error.error, 'error', isA<AppError>())
                .having(
                  (error) => (error.error as AppError).message,
                  'message',
                  'identity failed',
                ),
          ),
        );
        expect(adapterCalls, 0);
      },
    );

    test(
      'RestClient returns the identity failure without sending request',
      () async {
        var adapterCalls = 0;
        final dio = Dio(BaseOptions(baseUrl: 'https://api.test'))
          ..interceptors.add(
            InstallationInterceptor(
              installationIdentityService: _FakeInstallationIdentityService(
                result: const Failure(
                  AppError(
                    code: AppErrorCode.storageError,
                    message: 'identity failed',
                  ),
                ),
              ),
            ),
          );

        dio.httpClientAdapter = _FakeHttpClientAdapter((options, _, _) async {
          adapterCalls++;
          return _jsonResponse(200, {
            'data': {'ok': true},
          });
        });

        final client = DioRestClient(dio: dio);

        final result = await client.get(
          const RestClientRequest(path: '/accounts'),
        );

        expect(result, isA<Failure>());
        expect(result.error?.code, AppErrorCode.storageError);
        expect(result.error?.message, 'identity failed');
        expect(adapterCalls, 0);
      },
    );
  });
}

ResponseBody _jsonResponse(int statusCode, Map<String, Object?> body) {
  return ResponseBody.fromString(
    jsonEncode(body),
    statusCode,
    headers: {
      Headers.contentTypeHeader: [Headers.jsonContentType],
    },
  );
}

class _FakeInstallationIdentityService extends InstallationIdentityService {
  _FakeInstallationIdentityService({required Result<String> result})
    : _result = result,
      super(
        secureStorage: _NoopSecureStorage(),
        markerStore: _NoopMarkerStore(),
      );

  final Result<String> _result;

  @override
  AsyncResult<String> resolve() async => _result;
}

class _NoopMarkerStore implements InstallationMarkerStore {
  @override
  AsyncResult<bool> hasMarker() async => const Success(true);

  @override
  AsyncResult<Unit> markResolved() async => const Success(unit);
}

class _MemorySecureStorage implements LocalSecureStorage {
  _MemorySecureStorage(this.values);

  final Map<String, String> values;

  @override
  AsyncResult<Unit> write(String key, String value) async {
    values[key] = value;
    return const Success(unit);
  }

  @override
  AsyncResult<String> read(String key) async {
    final value = values[key];
    if (value == null) {
      return Failure(
        AppError(
          code: AppErrorCode.storageNotFound,
          message: 'Key not found: $key',
        ),
      );
    }

    return Success(value);
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    values.remove(key);
    return const Success(unit);
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    values.clear();
    return const Success(unit);
  }

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    return values.keys.where((key) => key.startsWith(pattern)).toList();
  }
}

class _NoopSecureStorage implements LocalSecureStorage {
  @override
  AsyncResult<Unit> write(String key, String value) async =>
      const Success(unit);

  @override
  AsyncResult<String> read(String key) async =>
      const Success('018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8');

  @override
  AsyncResult<Unit> delete(String key) async => const Success(unit);

  @override
  AsyncResult<Unit> deleteAll() async => const Success(unit);

  @override
  Future<List<String>> keysWithPrefix(String pattern) async => const [];
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
