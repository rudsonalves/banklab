import 'package:auto_injector/auto_injector.dart';
import 'package:dio/dio.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '/core/services/client_http/client/rest_client.dart';
import '/core/services/client_http/dio/dio_factory.dart';
import '/core/services/client_http/interceptors/interceptors.dart';
import '../resources/app_env.dart';
import 'app_section/app_section.dart';
import 'client_http/dio/dio_rest_client.dart';
import 'installation_identity/installation_identity.dart';
import 'secure_storage/flutter_secure_storage_local_storage.dart';
import 'secure_storage/local_secure_storage.dart';

class CoreServices {
  static void add(AutoInjector injector) {
    injector
      // 1. Secure Storage
      ..add<FlutterSecureStorage>(FlutterSecureStorage.new)
      // 2. Local Secure Storage (wrapper around FlutterSecureStorage)
      ..add<LocalSecureStorage>(
        () => FlutterSecureStorageLocalStorage(
          storage: injector.get<FlutterSecureStorage>(),
        ),
      )
      // 3. Installation identity
      ..add<InstallationMarkerStore>(FileInstallationMarkerStore.new)
      ..add<InstallationIdentityService>(
        () => InstallationIdentityService(
          secureStorage: injector.get<LocalSecureStorage>(),
          markerStore: injector.get<InstallationMarkerStore>(),
        ),
      )
      // 4. Main Dio instance without request interceptors
      ..addSingleton<Dio>(() => DioFactory.create())
      // 5. InstallationInterceptor
      ..addSingleton<InstallationInterceptor>(
        () => InstallationInterceptor(
          installationIdentityService: injector
              .get<InstallationIdentityService>(),
        ),
      )
      // 6. AuthInterceptor with its own Dio instance (no interceptors to avoid recursion)
      ..addSingleton<AuthInterceptor>(
        () => AuthInterceptor(
          authDio: injector.get<Dio>(),
          secureStorage: injector.get<LocalSecureStorage>(),
          baseUrl: AppEnv.baseUrl,
        ),
      )
      // 7. RestClient with installation and auth interceptors
      ..addSingleton<RestClient>(() {
        final dio = injector.get<Dio>();
        dio.interceptors.add(injector.get<InstallationInterceptor>());
        dio.interceptors.add(injector.get<AuthInterceptor>());
        return DioRestClient(dio: dio);
      })
      ..addSingleton<AppSection>(AppSection.new);
  }
}
