import 'dart:convert';
import 'dart:typed_data';

import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/interceptors/auth/auth_interceptor.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('AuthInterceptor', () {
    test(
      'serializes concurrent refresh attempts and retries all requests',
      () async {
        final storage = _MemorySecureStorage({
          StorageKeys.accessToken: 'old-access',
          StorageKeys.refreshToken: 'old-refresh',
        });
        var refreshCalls = 0;
        var protectedCalls = 0;

        final authDio = Dio(BaseOptions(baseUrl: 'https://api.test'));
        final refreshDio = Dio(BaseOptions(baseUrl: 'https://api.test'));

        authDio.httpClientAdapter = _FakeHttpClientAdapter((
          options,
          _,
          _,
        ) async {
          protectedCalls++;

          final authorization = options.headers['Authorization'];
          if (authorization == 'Bearer new-access') {
            return _jsonResponse(200, {
              'data': {'ok': true},
            });
          }

          return _jsonResponse(401, {
            'error': {'code': 'INVALID_TOKEN'},
          });
        });

        refreshDio.httpClientAdapter = _FakeHttpClientAdapter((
          options,
          _,
          _,
        ) async {
          refreshCalls++;
          await Future<void>.delayed(const Duration(milliseconds: 20));

          return _jsonResponse(200, {
            'data': {
              'access_token': 'new-access',
              'refresh_token': 'new-refresh',
            },
          });
        });

        authDio.interceptors.add(
          AuthInterceptor(
            authDio: authDio,
            refreshDio: refreshDio,
            secureStorage: storage,
            baseUrl: 'https://api.test',
          ),
        );

        final responses = await Future.wait([
          authDio.get('/protected'),
          authDio.get('/protected'),
        ]);

        expect(responses.map((response) => response.statusCode), [200, 200]);
        expect(refreshCalls, 1);
        expect(protectedCalls, 4);
        expect(storage.values[StorageKeys.accessToken], 'new-access');
        expect(storage.values[StorageKeys.refreshToken], 'new-refresh');
        expect(storage.deleteCalls, isEmpty);
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

class _MemorySecureStorage implements LocalSecureStorage {
  _MemorySecureStorage(this.values);

  final Map<String, String> values;
  final List<String> deleteCalls = [];

  @override
  Future<List<String>> keysWithPrefix(String pattern) async {
    return values.keys.where((key) => key.startsWith(pattern)).toList();
  }

  @override
  AsyncResult<Unit> delete(String key) async {
    deleteCalls.add(key);
    values.remove(key);
    return const Success(unit);
  }

  @override
  AsyncResult<Unit> deleteAll() async {
    deleteCalls.addAll(values.keys);
    values.clear();
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
  AsyncResult<Unit> write(String key, String value) async {
    values[key] = value;
    return const Success(unit);
  }
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
