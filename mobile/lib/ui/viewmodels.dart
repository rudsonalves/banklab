import 'package:auto_injector/auto_injector.dart';

import 'pages/auth/login/viewmodel/login_viewmodel.dart';
import 'pages/auth/register/viewmodel/register_viewmodel.dart';
import 'pages/auth/short_login/viewmodel/short_login_viewmodel.dart';
import 'pages/home/transfer/viewmodel/transfer_viewmodel.dart';
import 'pages/home/viewmodel/home_viewmodel.dart';
import 'pages/shared/details/viewmodel/details_viewmodel.dart';
import 'pages/splash/viewmodel/splash_viewmodel.dart';
import 'pages/statement/viewmodel/statement_viewmodel.dart';

class Viewmodels {
  static void add(AutoInjector injector) {
    injector
      ..add<HomeViewmodel>(HomeViewmodel.new)
      ..add<LoginViewModel>(LoginViewModel.new)
      ..addLazySingleton<RegisterViewmodel>(RegisterViewmodel.new)
      ..add<TransferViewmodel>(TransferViewmodel.new)
      ..add<DetailsViewmodel>(DetailsViewmodel.new)
      ..add<SplashViewmodel>(SplashViewmodel.new)
      ..add<ShortLoginViewModel>(ShortLoginViewModel.new)
      ..add<StatementViewmodel>(StatementViewmodel.new);
  }
}
