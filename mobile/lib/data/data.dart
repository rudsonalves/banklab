import 'package:auto_injector/auto_injector.dart';

import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/account/account_repository_impl.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/repositories/auth/auth_repository_impl.dart';
import '/data/services/apis/account/balance_api.dart';
import '/data/services/apis/account/list_accounts_api.dart';
import '/data/services/apis/account/statement_api.dart';
import '/data/services/apis/auth/auth_api.dart';

class Data {
  static void add(AutoInjector injector) {
    injector
      ..addSingleton<AuthRepository>(
        () => AuthRepositoryImpl(
          api: injector.get<AuthApi>(),
          storage: injector.get<LocalSecureStorage>(),
        ),
      )
      ..addSingleton<AccountRepository>(
        () => AccountRepositoryImpl(
          balanceApi: injector.get<BalanceApi>(),
          listAccountsApi: injector.get<ListAccountsApi>(),
          statementApi: injector.get<StatementApi>(),
        ),
      );
  }
}
