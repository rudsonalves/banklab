import 'package:dio/dio.dart';

import '/core/config/app_env.dart';

class DioFactory {
  static Dio create({
    Map<String, String>? defaultHeaders,
    List<Interceptor> interceptors = const [],
  }) {
    final dio = Dio(
      BaseOptions(
        baseUrl: AppEnv.baseUrl,
        connectTimeout: Duration(milliseconds: AppEnv.connectTimeout),
        receiveTimeout: Duration(milliseconds: AppEnv.receiveTimeout),
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          ...?defaultHeaders,
        },
      ),
    );

    _registerInterceptors(dio, interceptors);

    return dio;
  }

  static void _registerInterceptors(
    Dio dio,
    List<Interceptor> interceptors,
  ) {
    for (final interceptor in interceptors) {
      _addIfAbsent(dio, interceptor);
    }
  }

  static void _addIfAbsent(
    Dio dio,
    Interceptor interceptor,
  ) {
    final exists = dio.interceptors.any(
      (i) => i.runtimeType == interceptor.runtimeType,
    );

    if (!exists) {
      dio.interceptors.add(interceptor);
    }
  }
}
