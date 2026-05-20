import 'package:auto_injector/auto_injector.dart';

import '/data/repositories/account/account_repository.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/repositories/register_draft/register_draft_repository.dart';
import '/data/repositories/transaction/transaction_repository.dart';
import 'details/details_usecase.dart';
import 'register/register_usecase.dart';
import 'transfer/transfer_usecase.dart';

class Usecases {
  static void add(AutoInjector injector) {
    injector
      ..add<TransferUsecase>(
        () => TransferUsecase(
          accountRepo: injector.get<AccountRepository>(),
          transactionRepo: injector.get<TransactionRepository>(),
        ),
      )
      ..add<DetailsUsecase>(
        () => DetailsUsecase(
          accountRepo: injector.get<AccountRepository>(),
          transactionRepo: injector.get<TransactionRepository>(),
        ),
      )
      ..addLazySingleton<RegisterUsecase>(
        () => RegisterUsecase(
          authRepository: injector.get<AuthRepository>(),
          registerDraftRepository: injector.get<RegisterDraftRepository>(),
        ),
      );
  }
}
