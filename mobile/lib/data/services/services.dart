import 'package:auto_injector/auto_injector.dart';

import 'apis/account/balance_api.dart';
import 'apis/account/list_accounts_api.dart';
import 'apis/account/statement_api.dart';
import 'apis/auth/auth_api.dart';
import 'apis/contact_verification/contact_verification_api.dart';
import 'apis/installation/installation_api.dart';
import 'apis/receipt/api_receipt.dart';
import 'apis/registration/registration_api.dart';
import 'apis/transaction_password/transaction_password_api.dart';
import 'apis/transfer/api_transfer.dart';
import 'cache/last_login/last_login_cache_service.dart';
import 'cache/last_login/last_login_cache_service_impl.dart';
import 'cache/register_draft/register_draft_store.dart';

class Services {
  static void add(AutoInjector injector) {
    injector
      ..addSingleton<AuthApi>(AuthApi.new)
      ..addLazySingleton<RegistrationApi>(RegistrationApi.new)
      ..addLazySingleton<ContactVerificationApi>(ContactVerificationApi.new)
      ..addSingleton<BalanceApi>(BalanceApi.new)
      ..addSingleton<ListAccountsApi>(ListAccountsApi.new)
      ..addSingleton<StatementApi>(StatementApi.new)
      ..add<InstallationApi>(
        () => InstallationApi(
          client: injector.get(),
          installationIdentityService: injector.get(),
        ),
      )
      ..addSingleton<ApiTransfer>(ApiTransfer.new)
      ..addSingleton<ApiReceipt>(ApiReceipt.new)
      ..add<LastLoginCacheService>(LastLoginCacheServiceImpl.new)
      ..addLazySingleton<RegisterDraftStore>(RegisterDraftStore.new)
      ..add<TransactionPasswordApi>(TransactionPasswordApi.new);
  }
}
