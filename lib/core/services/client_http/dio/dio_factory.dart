import 'package:dio/dio.dart';

import '../client/rest_client.dart';
import 'dio_rest_client.dart';

class DioFactory {
  static RestClient create({
    required String baseUrl,
    Duration connectTimeout = const Duration(seconds: 10),
    Duration receiveTimeout = const Duration(seconds: 10),
    List<Interceptor> interceptors = const [],
  }) {
    final dio = Dio(
      BaseOptions(
        baseUrl: baseUrl,
        connectTimeout: connectTimeout,
        receiveTimeout: receiveTimeout,
        headers: const {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
      ),
    );

    _registerInterceptors(dio, interceptors);

    return DioRestClient(dio: dio);
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
    final alreadyRegistered = dio.interceptors.any(
      (existing) => existing.runtimeType == interceptor.runtimeType,
    );

    if (!alreadyRegistered) {
      dio.interceptors.add(interceptor);
    }
  }
}
