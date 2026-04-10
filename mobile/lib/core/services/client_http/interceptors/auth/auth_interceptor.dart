import 'dart:developer';

import 'package:dio/dio.dart';

import '/core/config/storage_keys.dart';
import '/core/services/secure_storage/local_secure_storage.dart';

/// NOTE: Potential race condition during token refresh
///
/// If multiple requests receive a 401 response simultaneously,
/// each one may trigger its own refresh token request.
///
/// This can lead to:
/// - Duplicate refresh calls
/// - Token overwrites (last write wins)
/// - Unnecessary load on the backend
///
/// Recommended future improvement:
/// Introduce a refresh lock (e.g., using a Completer or mutex)
/// to ensure only one refresh request is executed at a time,
/// while other requests await its result.
///
/// This is intentionally not implemented now to keep the
/// current design simple, but should be addressed if
/// concurrent requests become frequent.

class AuthInterceptor extends Interceptor {
  static const _refreshPath = '/auth/refresh';

  final Dio _authDio;
  final Dio _refreshDio;
  final LocalSecureStorage _secureStorage;

  AuthInterceptor({
    required Dio authDio,
    required LocalSecureStorage secureStorage,
    required String baseUrl,
    Duration timeout = const Duration(seconds: 10),
  }) : _authDio = authDio,
       _secureStorage = secureStorage,
       _refreshDio = Dio(
         BaseOptions(
           baseUrl: baseUrl,
           connectTimeout: timeout,
           receiveTimeout: timeout,
           headers: const {
             'Accept': 'application/json',
           },
         ),
       );

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    // If the request already has an Authorization header, skip adding the token
    if (options.headers.containsKey('Authorization')) {
      return handler.next(options);
    }

    final tokenResult = await _secureStorage.read(StorageKeys.accessToken);

    tokenResult.fold(
      onSuccess: (token) {
        if (token.isNotEmpty) {
          options.headers['Authorization'] = 'Bearer $token';
        }
      },
      onFailure: (_) {},
    );

    handler.next(options);
  }

  @override
  Future<void> onError(
    DioException err,
    ErrorInterceptorHandler handler,
  ) async {
    final statusCode = err.response?.statusCode;
    final path = err.requestOptions.path;

    log(
      '[AuthInterceptor] ERROR $statusCode → ${err.requestOptions.method} $path',
    );

    // only attempt refresh if 401 from non-refresh endpoint
    if (statusCode != 401 || path.endsWith(_refreshPath)) {
      return handler.next(err);
    }

    final refreshResult = await _secureStorage.read(StorageKeys.refreshToken);

    if (refreshResult.isFailure) {
      await _clearSession();
      return handler.next(err);
    }

    final refreshToken = refreshResult.value!;

    if (refreshToken.isEmpty) {
      await _clearSession();
      return handler.next(err);
    }

    try {
      final response = await _refreshDio.post(
        _refreshPath,
        data: {'refresh_token': refreshToken},
      );

      final data = response.data;
      final payload = data is Map<String, dynamic> ? data : {};
      final inner = payload['data'] ?? payload;

      final newAccessToken = inner['access_token'] as String?;
      final newRefreshToken = inner['refresh_token'] as String?;

      if (newAccessToken == null || newAccessToken.isEmpty) {
        throw Exception('Invalid refresh response');
      }

      await _secureStorage.write(StorageKeys.accessToken, newAccessToken);

      if (newRefreshToken != null && newRefreshToken.isNotEmpty) {
        await _secureStorage.write(StorageKeys.refreshToken, newRefreshToken);
      }

      // retry original request with new token
      final request = err.requestOptions;
      final newRequest =
          Options(
            method: request.method,
            headers: {
              ...request.headers,
              'Authorization': 'Bearer $newAccessToken',
            },
            responseType: request.responseType,
            contentType: request.contentType,
            extra: request.extra,
            followRedirects: request.followRedirects,
            validateStatus: request.validateStatus,
          ).compose(
            _authDio.options,
            request.path,
            data: request.data,
            queryParameters: request.queryParameters,
          );

      final retryResponse = await _authDio.fetch(newRequest);

      return handler.resolve(retryResponse);
    } catch (e, s) {
      log('[AuthInterceptor] refresh failed: $e');
      log('$s');

      await _clearSession();
      return handler.next(err);
    }
  }

  Future<void> _clearSession() async {
    await _secureStorage.delete(StorageKeys.accessToken);
    await _secureStorage.delete(StorageKeys.refreshToken);
  }
}
