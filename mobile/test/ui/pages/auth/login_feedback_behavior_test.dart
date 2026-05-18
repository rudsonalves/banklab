import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/data/repositories/auth/auth_repository.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_confirm_request_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_confirm_response_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_request_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/contact_verification_request_response_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/login_request_dto.dart';
import 'package:bankflow/data/services/auth/api/dtos/register_request_dto.dart';
import 'package:bankflow/data/services/auth/cache/models/last_login_identity.dart';
import 'package:bankflow/domain/common/auth/models/auth_user.dart';
import 'package:bankflow/domain/common/auth/models/user_profile.dart';
import 'package:bankflow/ui/components/text_form_field/basic_text_form_field.dart';
import 'package:bankflow/ui/pages/auth/login/login_page.dart';
import 'package:bankflow/ui/pages/auth/login/viewmodel/login_viewmodel.dart';
import 'package:bankflow/ui/pages/auth/short_login/short_login_page.dart';
import 'package:bankflow/ui/pages/auth/short_login/viewmodel/short_login_viewmodel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';

const _approvalRequiredMessage =
    'Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você poderá acessar sua conta.';

void main() {
  group('Login UI feedback behavior', () {
    testWidgets(
      'full login shows approval-required message and stays on login screen',
      (tester) async {
        final repository = _FakeAuthRepository(
          loginResult: const Failure(
            AppError(
              code: AppErrorCode.accountApprovalRequired,
              message: 'backend message should be ignored for this case',
            ),
          ),
        );

        await tester.pumpWidget(
          MaterialApp(
            home: LoginPage(
              viewModel: LoginViewModel(authRepository: repository),
            ),
          ),
        );

        await _submitFullLogin(tester);

        expect(find.text(_approvalRequiredMessage), findsOneWidget);
        expect(
          find.text('Acesse sua conta para continuar no BankFlow.'),
          findsOneWidget,
        );
      },
    );

    testWidgets(
      'short login shows approval-required message and keeps remembered identity visible',
      (tester) async {
        final repository = _FakeAuthRepository(
          loginResult: const Failure(
            AppError(
              code: AppErrorCode.accountApprovalRequired,
              message: 'backend message should be ignored for this case',
            ),
          ),
        );

        await tester.pumpWidget(
          MaterialApp(
            home: ShortLoginPage(
              viewModel: ShortLoginViewModel(authRepository: repository),
              identity: LastLoginIdentity(
                name: 'Maria Silva',
                identifier: 'maria@example.com',
              ),
            ),
          ),
        );

        await _submitShortLogin(tester);

        expect(find.text(_approvalRequiredMessage), findsOneWidget);
        expect(find.text('Maria Silva'), findsOneWidget);
        expect(find.text('Entrar com outra conta'), findsOneWidget);
      },
    );

    testWidgets(
      'full login keeps invalid credentials feedback unchanged',
      (tester) async {
        final repository = _FakeAuthRepository(
          loginResult: const Failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Credenciais invalidas',
            ),
          ),
        );

        await tester.pumpWidget(
          MaterialApp(
            home: LoginPage(
              viewModel: LoginViewModel(authRepository: repository),
            ),
          ),
        );

        await _submitFullLogin(tester);

        expect(find.text('Credenciais invalidas'), findsOneWidget);
        expect(find.text(_approvalRequiredMessage), findsNothing);
      },
    );

    testWidgets(
      'short login keeps generic failure feedback unchanged',
      (tester) async {
        final repository = _FakeAuthRepository(
          loginResult: const Failure(
            AppError(
              code: AppErrorCode.httpError,
              message: 'Falha temporaria no servidor',
            ),
          ),
        );

        await tester.pumpWidget(
          MaterialApp(
            home: ShortLoginPage(
              viewModel: ShortLoginViewModel(authRepository: repository),
              identity: LastLoginIdentity(
                name: 'Maria Silva',
                identifier: 'maria@example.com',
              ),
            ),
          ),
        );

        await _submitShortLogin(tester);

        expect(find.text('Falha temporaria no servidor'), findsOneWidget);
        expect(find.text(_approvalRequiredMessage), findsNothing);
      },
    );
  });
}

Future<void> _submitFullLogin(WidgetTester tester) async {
  final fields = find.byType(BasicTextFormField);
  expect(fields, findsNWidgets(2));
  await tester.enterText(fields.first, 'customer@example.com');
  await tester.enterText(fields.last, '123456');
  await tester.tap(find.widgetWithText(ElevatedButton, 'Entrar'));
  await tester.pump();
  await tester.pump(const Duration(milliseconds: 50));
}

Future<void> _submitShortLogin(WidgetTester tester) async {
  final passwordField = find.byType(BasicTextFormField);
  expect(passwordField, findsOneWidget);
  await tester.enterText(passwordField, '123456');
  await tester.tap(find.widgetWithText(ElevatedButton, 'Entrar'));
  await tester.pump();
  await tester.pump(const Duration(milliseconds: 50));
}

class _FakeAuthRepository implements AuthRepository {
  final Result<LoggedUser> loginResult;

  _FakeAuthRepository({required this.loginResult});

  @override
  AuthUser get currentUser => NotLoggedUser();

  @override
  bool get isLoggedIn => false;

  @override
  UserProfile? get userProfile => null;

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async {
    return const Failure(
      AppError(code: AppErrorCode.storageNotFound, message: 'Not found'),
    );
  }

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    return loginResult;
  }

  @override
  AsyncResult<Unit> logout() async {
    return Success(unit);
  }

  @override
  AsyncResult<UserProfile> profile() async {
    return Failure(
      AppError(
        code: AppErrorCode.unauthenticated,
        message: 'User is not logged in.',
      ),
    );
  }

  @override
  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    return Success(unit);
  }

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    return const Failure(
      AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
    );
  }

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    return const Failure(
      AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
    );
  }
}
