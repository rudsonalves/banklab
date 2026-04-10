import 'package:auto_injector/auto_injector.dart';

import 'pages/auth/login/viewmodel/login_viewmodel.dart';
import 'pages/auth/register/viewmodel/register_viewmodel.dart';
import 'pages/home/viewmodel/home_viewmodel.dart';

class Uis {
  static void add(AutoInjector injector) {
    injector
      ..add<HomeViewmodel>(HomeViewmodel.new)
      ..add<LoginViewModel>(LoginViewModel.new)
      ..add<RegisterViewmodel>(RegisterViewmodel.new);
  }
}
