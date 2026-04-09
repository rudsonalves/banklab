import 'package:auto_injector/auto_injector.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

import '/core/services/secure_storage/flutter_secure_storage_local_storage.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/repositories/auth/auth_repository_impl.dart';
import '/data/services/apis/auth/auth_api.dart';

class Data {
  static void add(AutoInjector injector) {
    injector
      ..addSingleton<LocalSecureStorage>(
        () => FlutterSecureStorageLocalStorage(
          storage: injector.get<FlutterSecureStorage>(),
        ),
      )
      ..addSingleton<AuthRepository>(
        () => AuthRepositoryImpl(
          api: injector.get<AuthApi>(),
          storage: injector.get<LocalSecureStorage>(),
        ),
      );
  }
}
