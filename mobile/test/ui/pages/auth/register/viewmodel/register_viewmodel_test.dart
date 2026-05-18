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
import 'package:bankflow/ui/pages/auth/register/viewmodel/register_viewmodel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RegisterViewmodel', () {
    test('starts at personalData with empty verification state', () {
      final viewmodel = RegisterViewmodel(
        authRepository: _FakeAuthRepository(),
      );

      expect(viewmodel.currentStep, RegisterStep.personalData);
      expect(viewmodel.isEmailVerified, isFalse);
      expect(viewmodel.isPhoneVerified, isFalse);
      expect(viewmodel.canRegister, isFalse);
    });

    test('blocks nextStep when current step data is invalid', () {
      final viewmodel = RegisterViewmodel(
        authRepository: _FakeAuthRepository(),
      );

      final advanced = viewmodel.nextStep();

      expect(advanced, isFalse);
      expect(viewmodel.currentStep, RegisterStep.personalData);
      expect(viewmodel.stepError?.code, AppErrorCode.invalidData);
      expect(viewmodel.stepErrorStep, RegisterStep.personalData);
    });

    test('keeps entered data when navigating back and forth', () {
      final viewmodel = RegisterViewmodel(
        authRepository: _FakeAuthRepository(),
      );

      viewmodel.updatePersonalData(
        name: 'Maria Silva',
        cpf: '123.456.789-01',
        birthDate: DateTime(1990, 1, 15),
        password: 'P@ssword123',
      );
      viewmodel.updateContactData(
        email: 'maria@example.com',
        phone: '+5511999999999',
      );

      expect(viewmodel.nextStep(), isTrue);
      expect(viewmodel.nextStep(), isTrue);
      expect(viewmodel.currentStep, RegisterStep.emailVerification);

      expect(viewmodel.previousStep(), isTrue);
      expect(viewmodel.currentStep, RegisterStep.contactData);

      expect(viewmodel.name, 'Maria Silva');
      expect(viewmodel.cpf, '12345678901');
      expect(viewmodel.birthDate, DateTime(1990, 1, 15));
      expect(viewmodel.email, 'maria@example.com');
      expect(viewmodel.phone, '+5511999999999');
      expect(viewmodel.password, 'P@ssword123');
    });

    test(
      'confirms email and phone and transitions through verification steps',
      () async {
        final repository = _FakeAuthRepository();
        final viewmodel = RegisterViewmodel(authRepository: repository);

        viewmodel.updatePersonalData(
          name: 'Maria Silva',
          cpf: '12345678901',
          birthDate: DateTime(1990, 1, 15),
          password: 'P@ssword123',
        );
        viewmodel.updateContactData(
          email: 'maria@example.com',
          phone: '+5511999999999',
        );

        viewmodel.goToStep(RegisterStep.emailVerification);
        await viewmodel.requestEmailCode.execute();
        await viewmodel.confirmEmailCode.execute('123456');

        expect(viewmodel.emailVerificationId, 'email-verification-id');
        expect(viewmodel.isEmailVerified, isTrue);
        expect(viewmodel.currentStep, RegisterStep.phoneVerification);

        await viewmodel.requestPhoneCode.execute();
        await viewmodel.confirmPhoneCode.execute('654321');

        expect(viewmodel.phoneVerificationId, 'phone-verification-id');
        expect(viewmodel.isPhoneVerified, isTrue);
        expect(viewmodel.currentStep, RegisterStep.review);
      },
    );

    test(
      'blocks final register until both verification tokens are confirmed',
      () async {
        final repository = _FakeAuthRepository();
        final viewmodel = RegisterViewmodel(authRepository: repository);

        viewmodel.updatePersonalData(
          name: 'Maria Silva',
          cpf: '12345678901',
          birthDate: DateTime(1990, 1, 15),
          password: 'P@ssword123',
        );
        viewmodel.updateContactData(
          email: 'maria@example.com',
          phone: '+5511999999999',
        );

        await viewmodel.register.execute();

        expect(viewmodel.register.isFailure, isTrue);
        expect(viewmodel.register.error?.code, AppErrorCode.invalidData);
        expect(repository.registerCalls, 0);
      },
    );

    test('registers with full payload after both confirmations', () async {
      final repository = _FakeAuthRepository();
      final viewmodel = RegisterViewmodel(authRepository: repository);

      viewmodel.updatePersonalData(
        name: 'Maria Silva',
        cpf: '123.456.789-01',
        birthDate: DateTime(1990, 1, 15),
        password: 'P@ssword123',
      );
      viewmodel.updateContactData(
        email: 'maria@example.com',
        phone: '+5511999999999',
      );

      await viewmodel.requestEmailCode.execute();
      await viewmodel.confirmEmailCode.execute('123456');
      await viewmodel.requestPhoneCode.execute();
      await viewmodel.confirmPhoneCode.execute('654321');
      await viewmodel.register.execute();

      expect(viewmodel.register.isSuccess, isTrue);
      expect(repository.registerCalls, 1);
      final dto = repository.lastRegisterDto;
      expect(dto, isNotNull);
      expect(dto!.name, 'Maria Silva');
      expect(dto.email, 'maria@example.com');
      expect(dto.phone, '+5511999999999');
      expect(dto.password, 'P@ssword123');
      expect(dto.birthDate, DateTime(1990, 1, 15));
      expect(dto.cpf, '12345678901');
      expect(dto.emailVerificationToken, 'email-verified-token');
      expect(dto.phoneVerificationToken, 'phone-verified-token');
    });
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
        target: isEmail ? 'maria@example.com' : '+5511999999999',
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
