import 'package:auto_injector/auto_injector.dart';

import '/core/services/secure_storage/local_secure_storage.dart';
import '/data/repositories/account/account_repository.dart';
import '/data/repositories/account/account_repository_impl.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/repositories/auth/auth_repository_impl.dart';
import 'repositories/transaction/transaction_repository.dart';
import 'repositories/transaction/transaction_repository_impl.dart';
import 'services/apis/account/balance_api.dart';
import 'services/apis/account/list_accounts_api.dart';
import 'services/apis/account/statement_api.dart';
import 'services/apis/receipt/api_receipt.dart';
import 'services/apis/transfer/api_transfer.dart';
import 'services/auth/api/auth_api.dart';
import 'services/auth/cache/last_login_cache_service.dart';

class Repositories {
  static void add(AutoInjector injector) {
    injector
      ..addSingleton<AuthRepository>(
        () => AuthRepositoryImpl(
          api: injector.get<AuthApi>(),
          storage: injector.get<LocalSecureStorage>(),
          lastLoginCacheService: injector.get<LastLoginCacheService>(),
        ),
      )
      ..addSingleton<AccountRepository>(
        () => AccountRepositoryImpl(
          balanceApi: injector.get<BalanceApi>(),
          listAccountsApi: injector.get<ListAccountsApi>(),
          statementApi: injector.get<StatementApi>(),
        ),
      )
      ..addSingleton<TransactionRepository>(
        () => TransactionRepositoryImpl(
          apiTransfer: injector.get<ApiTransfer>(),
          apiReceipt: injector.get<ApiReceipt>(),
        ),
      );
  }
}
