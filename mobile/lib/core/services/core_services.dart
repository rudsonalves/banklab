import 'package:auto_injector/auto_injector.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '/core/config/enviroment_key.dart';
import '/core/services/client_http/client/rest_client.dart';
import '/core/services/client_http/dio/dio_factory.dart';

class CoreServices {
  static void add(AutoInjector injector) {
    injector
      ..add(FlutterSecureStorage.new)
      ..addSingleton<RestClient>(() {
        final client = DioFactory.create(
          baseUrl: EnviromentKey.baseUrl,
          connectTimeout: const Duration(seconds: 10),
          receiveTimeout: const Duration(seconds: 10),
        );
        return client;
      });
  }
}
