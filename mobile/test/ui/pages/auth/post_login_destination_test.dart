import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/app_section/app_section.dart';
import 'package:bankflow/data/repositories/auth/auth_repository.dart';
import 'package:bankflow/data/services/apis/auth/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/cache/last_login/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_session/auth_session.dart';
import 'package:bankflow/domain/common/auth/models/auth_user.dart';
import 'package:bankflow/domain/common/user/enums/user_role.dart';
import 'package:bankflow/ui/pages/auth/models/post_login_destination.dart';
import 'package:bankflow/ui/pages/auth/viewmodel/login_viewmodel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('LoginViewModel.resolvePostLoginDestination', () {
    test('returns sessionError when userProfile is missing', () {
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: AppSection(),
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.sessionError,
      );
    });

    test('returns home for active transaction password and home readiness', () {
      final appSection = AppSection()..setAuthSession(_session());
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: appSection,
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.home,
      );
    });

    test('blocks active transaction password when home is not allowed', () {
      final appSection = AppSection()
        ..setAuthSession(_session(hasOperationalAccount: false));
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: appSection,
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.blocked,
      );
    });

    test('returns transactionPassword for notSet status', () {
      final appSection = AppSection()
        ..setAuthSession(_session(status: TransactionPasswordStatus.notSet));
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: appSection,
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.transactionPassword,
      );
    });

    test('blocks locked status', () {
      final appSection = AppSection()
        ..setAuthSession(_session(status: TransactionPasswordStatus.locked));
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: appSection,
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.blocked,
      );
    });

    test('blocks unknown status', () {
      final appSection = AppSection()
        ..setAuthSession(_session(status: TransactionPasswordStatus.unknown));
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: appSection,
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.blocked,
      );
    });
  });

  group('ShortLoginViewModel.resolvePostLoginDestination', () {
    test('uses the same post-login destination rules', () {
      final appSection = AppSection()
        ..setAuthSession(_session(status: TransactionPasswordStatus.notSet));
      final viewModel = LoginViewModel(
        authRepository: _FakeAuthRepository(),
        appSection: appSection,
      );

      expect(
        viewModel.resolvePostLoginDestination(),
        PostLoginDestination.transactionPassword,
      );
    });
  });
}

AuthSession _session({
  TransactionPasswordStatus status = TransactionPasswordStatus.active,
  bool hasOperationalAccount = true,
}) {
  return AuthSession(
    user: UserSession(
      userId: 'user-1',
      email: 'customer@example.com',
      role: UserRole.customer,
    ),
    customer: CustomerSession(
      id: 'customer-1',
      name: 'Maria Silva',
      cpf: '12345678901',
      birthDate: DateTime(1990, 1, 1),
      createdAt: DateTime(2026, 5, 13),
    ),
    readiness: ReadinessSession(
      onboardingCompleted: true,
      approved: true,
      hasOperationalAccount: hasOperationalAccount,
      transactionPasswordStatus: status,
    ),
  );
}

class _FakeAuthRepository implements AuthRepository {
  _FakeAuthRepository();

  @override
  AuthUser get currentUser => NotLoggedUser();

  @override
  bool get isLoggedIn => currentUser is LoggedUser;

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async =>
      throw UnimplementedError();

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async =>
      throw UnimplementedError();

  @override
  AsyncResult<Unit> logout() async => throw UnimplementedError();

  @override
  AsyncResult<AuthSession> getAuthSession() async => throw UnimplementedError();
}
