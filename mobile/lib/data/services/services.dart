import 'package:auto_injector/auto_injector.dart';

import '/core/services/client_http/client/rest_client.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import 'apis/account/balance_api.dart';
import 'apis/account/list_accounts_api.dart';
import 'apis/account/statement_api.dart';
import 'apis/receipt/api_receipt.dart';
import 'apis/transfer/api_transfer.dart';
import 'auth/api/auth_api.dart';
import 'auth/cache/last_login_cache_service.dart';
import 'auth/cache/last_login_cache_service_impl.dart';

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
      )
      ..addSingleton<ApiTransfer>(
        () => ApiTransfer(injector.get<RestClient>()),
      )
      ..addSingleton<ApiReceipt>(
        () => ApiReceipt(injector.get<RestClient>()),
      )
      ..add<LastLoginCacheService>(
        () => LastLoginCacheServiceImpl(injector.get<LocalSecureStorage>()),
      );
  }
}
