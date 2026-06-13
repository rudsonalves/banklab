import 'package:auto_injector/auto_injector.dart';

import 'pages/auth/viewmodel/login_viewmodel.dart';
import 'pages/home/transfer/viewmodel/transfer_viewmodel.dart';
import 'pages/home/viewmodel/home_viewmodel.dart';
import 'pages/register/viewmodel/register_viewmodel.dart';
import 'pages/shared/details/viewmodel/details_viewmodel.dart';
import 'pages/splash/viewmodel/splash_viewmodel.dart';
import 'pages/statement/viewmodel/statement_viewmodel.dart';
import 'pages/transaction_password/setup/viewmodel/transaction_password_viewmodel.dart';
import 'pages/transaction_password/verification/viewmodel/verify_tansaction_password_viewmodel.dart';

class Viewmodels {
  static void add(AutoInjector injector) {
    injector
      ..add<HomeViewmodel>(HomeViewmodel.new)
      ..add<LoginViewModel>(LoginViewModel.new)
      ..addLazySingleton<RegisterViewmodel>(RegisterViewmodel.new)
      ..add<TransferViewmodel>(TransferViewmodel.new)
      ..add<DetailsViewmodel>(DetailsViewmodel.new)
      ..add<SplashViewmodel>(SplashViewmodel.new)
      ..add<TransactionPasswordViewModel>(TransactionPasswordViewModel.new)
      ..add<StatementViewmodel>(StatementViewmodel.new)
      ..add<VerifyTansactionPasswordViewmodel>(
        VerifyTansactionPasswordViewmodel.new,
      );
  }
}
