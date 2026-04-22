import 'package:auto_injector/auto_injector.dart';

import '/core/services/client_http/client/rest_client.dart';
import 'apis/account/balance_api.dart';
import 'apis/account/list_accounts_api.dart';
import 'apis/account/statement_api.dart';
import 'apis/auth/auth_api.dart';

class Services {
  static void add(AutoInjector injector) {
    injector
      ..addSingleton<AuthApi>(
        () => AuthApi(injector.get<RestClient>()),
      )
      ..addSingleton<BalanceApi>(
        () => BalanceApi(injector.get<RestClient>()),
      )
      ..addSingleton<ListAccountsApi>(
        () => ListAccountsApi(injector.get<RestClient>()),
      )
      ..addSingleton<StatementApi>(
        () => StatementApi(injector.get<RestClient>()),
      );
  }
}
