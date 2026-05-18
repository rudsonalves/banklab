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
import 'package:bankflow/ui/pages/auth/register/register_page.dart';
import 'package:bankflow/ui/pages/auth/register/viewmodel/register_viewmodel.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';

void main() {
  testWidgets('register page completes main multi-step flow', (tester) async {
    final repository = _FakeAuthRepository();
    final viewmodel = RegisterViewmodel(authRepository: repository);

    final router = GoRouter(
      initialLocation: '/register',
      routes: [
        GoRoute(
          path: '/register',
          name: 'register',
          builder: (context, state) => RegisterPage(viewmodel: viewmodel),
        ),
        GoRoute(
          path: '/login',
          name: 'login',
          builder: (context, state) =>
              const Scaffold(body: Text('Login Placeholder')),
        ),
      ],
    );

    await tester.pumpWidget(MaterialApp.router(routerConfig: router));
    await tester.pumpAndSettle();

    // Step 1: personal data + date picker.
    final personalFields = find.byType(BasicTextFormField);
    expect(personalFields, findsNWidgets(4));

    await tester.enterText(personalFields.at(0), 'Maria Silva');
    await tester.enterText(personalFields.at(1), '12345678901');
    await tester.enterText(personalFields.at(3), 'P@ssword123');

    await tester.ensureVisible(find.text('Selecionar data de nascimento'));
    await tester.tap(find.text('Selecionar data de nascimento'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('15').last);
    await tester.tap(find.text('OK'));
    await tester.pumpAndSettle();

    await tester.tap(find.widgetWithText(ElevatedButton, 'Continuar'));
    await tester.pumpAndSettle();

    // Step 2: contact data.
    final contactFields = find.byType(BasicTextFormField);
    expect(contactFields, findsNWidgets(2));

    await tester.enterText(contactFields.at(0), 'maria@example.com');
    await tester.enterText(contactFields.at(1), '27999999999');

    expect(find.text('(27) 99999-9999'), findsWidgets);

    await tester.tap(find.widgetWithText(ElevatedButton, 'Continuar'));
    await tester.pumpAndSettle();

    // Step 3: e-mail verification.
    await tester.ensureVisible(
      find.widgetWithText(ElevatedButton, 'Enviar codigo e-mail'),
    );
    await tester.tap(
      find.widgetWithText(ElevatedButton, 'Enviar codigo e-mail'),
    );
    await tester.pumpAndSettle();
    await tester.pump(const Duration(seconds: 4));

    final emailCodeField = find.byType(BasicTextFormField);
    await tester.enterText(emailCodeField, '123456');
    await tester.ensureVisible(
      find.widgetWithText(ElevatedButton, 'Confirmar e-mail'),
    );
    await tester.tap(find.widgetWithText(ElevatedButton, 'Confirmar e-mail'));
    await tester.pumpAndSettle();

    // Step 4: phone verification.
    final reviewButtonBeforeConfirm = tester.widget<ElevatedButton>(
      find.widgetWithText(ElevatedButton, 'Revisar cadastro'),
    );
    expect(reviewButtonBeforeConfirm.onPressed, isNull);

    await tester.ensureVisible(
      find.widgetWithText(ElevatedButton, 'Enviar codigo telefone'),
    );
    await tester.tap(
      find.widgetWithText(ElevatedButton, 'Enviar codigo telefone'),
    );
    await tester.pumpAndSettle();
    await tester.pump(const Duration(seconds: 4));

    final phoneCodeField = find.byType(BasicTextFormField);
    await tester.enterText(phoneCodeField, '654321');
    await tester.ensureVisible(
      find.widgetWithText(ElevatedButton, 'Confirmar telefone'),
    );
    await tester.tap(find.widgetWithText(ElevatedButton, 'Confirmar telefone'));
    await tester.pumpAndSettle();

    // Step 5: review + final register.
    expect(
      find.textContaining('Revise seus dados antes de concluir.'),
      findsOneWidget,
    );
    await tester.tap(find.widgetWithText(ElevatedButton, 'Concluir cadastro'));
    await tester.pumpAndSettle();

    expect(repository.registerCalls, 1);
    expect(repository.lastRegisterDto?.email, 'maria@example.com');
    expect(repository.lastRegisterDto?.phone, '+5527999999999');
    expect(
      repository.lastRegisterDto?.emailVerificationToken,
      'email-verified-token',
    );
    expect(
      repository.lastRegisterDto?.phoneVerificationToken,
      'phone-verified-token',
    );
    expect(find.text('Login Placeholder'), findsOneWidget);
  });
}

class _FakeAuthRepository implements AuthRepository {
  int registerCalls = 0;
  RegisterRequestDto? lastRegisterDto;

  @override
  AuthUser get currentUser => NotLoggedUser();

  @override
  bool get isLoggedIn => false;

  @override
  UserProfile? get userProfile => null;

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    final isEmail = dto.verificationId == 'email-verification-id';

    return Success(
      ContactVerificationConfirmResponseDto(
        verificationToken: isEmail
            ? 'email-verified-token'
            : 'phone-verified-token',
        channel: isEmail ? 'email' : 'phone',
        target: isEmail ? 'maria@example.com' : '+5527999999999',
        verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
      ),
    );
  }

  @override
  AsyncResult<LastLoginIdentity> getLastLoginIdentity() async {
    return const Failure(
      AppError(code: AppErrorCode.storageNotFound, message: 'Not found'),
    );
  }

  @override
  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    return const Failure(
      AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
    );
  }

  @override
  AsyncResult<Unit> logout() async {
    return Success(unit);
  }

  @override
  AsyncResult<UserProfile> profile() async {
    return const Failure(
      AppError(code: AppErrorCode.unauthenticated, message: 'Not logged in'),
    );
  }

  @override
  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    registerCalls++;
    lastRegisterDto = dto;
    return Success(unit);
  }

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    final isEmail = dto.channel == 'email';
    return Success(
      ContactVerificationRequestResponseDto(
        verificationId: isEmail
            ? 'email-verification-id'
            : 'phone-verification-id',
        channel: dto.channel,
        target: dto.target,
        token: isEmail ? '123456' : '654321',
        expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
      ),
    );
  }
}
