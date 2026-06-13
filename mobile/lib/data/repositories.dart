import 'package:auto_injector/auto_injector.dart';

import '/data/repositories/account/account_repository.dart';
import '/data/repositories/account/account_repository_impl.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/repositories/auth/auth_repository_impl.dart';
import '/data/repositories/register_draft/register_draft_repository.dart';
import '/data/repositories/register_draft/register_draft_repository_impl.dart';
import 'repositories/contact_verification/contact_verification_repository.dart';
import 'repositories/contact_verification/contact_verification_repository_impl.dart';
import 'repositories/registration/registration_repository.dart';
import 'repositories/registration/registration_repository_impl.dart';
import 'repositories/transaction_password/transaction_password_repository.dart';
import 'repositories/transaction_password/transaction_password_repository_impl.dart';
import 'repositories/transfer/transfer_repository.dart';
import 'repositories/transfer/transfer_repository_impl.dart';

class Repositories {
  static void add(AutoInjector injector) {
    injector
      ..addSingleton<AuthRepository>(AuthRepositoryImpl.new)
      ..addSingleton<AccountRepository>(AccountRepositoryImpl.new)
      ..addSingleton<TransferRepository>(TransferRepositoryImpl.new)
      ..addLazySingleton<RegisterDraftRepository>(
        RegisterDraftRepositoryImpl.new,
      )
      ..addLazySingleton<ContactVerificationRepository>(
        ContactVerificationRepositoryImpl.new,
      )
      ..addLazySingleton<RegistrationRepository>(
        RegistrationRepositoryImpl.new,
      )
      ..add<TransactionPasswordRepository>(
        TransactionPasswordRepositoryImpl.new,
      );
  }
}
