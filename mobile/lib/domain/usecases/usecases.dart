import 'package:auto_injector/auto_injector.dart';

import 'details/details_usecase.dart';
import 'register/register_usecase.dart';
import 'transfer/transfer_usecase.dart';

class Usecases {
  static void add(AutoInjector injector) {
    injector
      ..add<TransferUsecase>(TransferUsecase.new)
      ..add<DetailsUsecase>(DetailsUsecase.new)
      ..add<RegisterUsecase>(RegisterUsecase.new);
  }
}
