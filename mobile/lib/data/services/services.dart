import 'package:auto_injector/auto_injector.dart';

import '/core/services/client_http/client/rest_client.dart';
import 'apis/auth/auth_api.dart';

class Services {
  static void add(AutoInjector injector) {
    injector.addSingleton<AuthApi>(
      () => AuthApi(injector.get<RestClient>()),
    );
  }
}
