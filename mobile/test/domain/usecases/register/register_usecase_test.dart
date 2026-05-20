import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/data/repositories/contact_verification/contact_verification_repository.dart';
import 'package:bankflow/data/repositories/register_draft/register_draft_repository.dart';
import 'package:bankflow/data/repositories/registration/registration_repository.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_response_dto.dart';
import 'package:bankflow/data/services/apis/registration/dtos/cpf_check_response_dto.dart';
import 'package:bankflow/data/services/apis/registration/dtos/register_request_dto.dart';
import 'package:bankflow/domain/common/auth/models/register_draft_snapshot.dart';
import 'package:bankflow/domain/usecases/register/register_usecase.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RegisterUsecase', () {
    test('startEmptyRegisterState initializes empty state', () {
      final draftRepository = _FakeRegisterDraftRepository(
        loadResult: Success<RegisterDraftSnapshot>(_snapshot()),
      );
      final usecase = _usecase(draftRepository: draftRepository);

      usecase.startEmptyRegisterState();

      expect(usecase.state, isNotNull);
      expect(usecase.state?.cpf, '');
      expect(usecase.state?.name, isNull);
      expect(draftRepository.loadedCpfs, isEmpty);
    });

    test('startEmptyRegisterState does not load draft or fail', () {
      final draftRepository = _FakeRegisterDraftRepository(
        loadResult: const Failure<RegisterDraftSnapshot>(
          AppError(code: AppErrorCode.unexpected, message: 'load failed'),
        ),
      );
      final usecase = _usecase(draftRepository: draftRepository);

      usecase.startEmptyRegisterState();

      expect(usecase.state, isNotNull);
      expect(usecase.state?.cpf, '');
      expect(usecase.state?.name, isNull);
      expect(draftRepository.loadedCpfs, isEmpty);
    });

    test('submitCPF rejects unavailable cpf without saving draft', () async {
      final registrationRepository = _FakeRegistrationRepository(
        cpfCheckResult: Success(
          CpfCheckResponseDto(
            cpf: '12345678909',
            exists: true,
            available: false,
          ),
        ),
      );
      final draftRepository = _FakeRegisterDraftRepository();
      final usecase = _usecase(
        registrationRepository: registrationRepository,
        draftRepository: draftRepository,
      );
      usecase.startEmptyRegisterState();

      final result = await usecase.submitCPF('123.456.789-09');

      expect(result, isA<Failure<Unit>>());
      expect(result.error?.message, 'CPF is already registered');
      expect(draftRepository.savedSnapshots, isEmpty);
    });

    test('request email token saves email and verification id', () async {
      final contactVerificationRepository =
          _FakeContactVerificationRepository();
      final draftRepository = _FakeRegisterDraftRepository();
      final usecase = _usecase(
        contactVerificationRepository: contactVerificationRepository,
        draftRepository: draftRepository,
      );
      usecase.startEmptyRegisterState();

      final result = await usecase.submitAndRequestEmailToken(
        'maria@example.com',
      );

      expect(result, isA<Success<Unit>>());
      expect(usecase.state?.email, 'maria@example.com');
      expect(usecase.state?.emailVerificationId, 'email-verification-id');
      expect(usecase.state?.isEmailVerified, isFalse);
      expect(draftRepository.savedSnapshots.last.email, 'maria@example.com');
      expect(
        draftRepository.savedSnapshots.last.emailVerificationId,
        'email-verification-id',
      );
    });

    test(
      'requesting a new email token invalidates previous confirmed token',
      () async {
        final registrationRepository = _FakeRegistrationRepository();
        final contactVerificationRepository =
            _FakeContactVerificationRepository();
        final draftRepository = _FakeRegisterDraftRepository();
        final usecase = _usecase(
          registrationRepository: registrationRepository,
          contactVerificationRepository: contactVerificationRepository,
          draftRepository: draftRepository,
        );
        await _prepareValidRegistration(usecase);

        final requestResult = await usecase.submitAndRequestEmailToken(
          'other@example.com',
        );
        final registerResult = await usecase.register();

        expect(requestResult, isA<Success<Unit>>());
        expect(usecase.state?.email, 'other@example.com');
        expect(usecase.state?.isEmailVerified, isFalse);
        expect(registerResult, isA<Failure<Unit>>());
        expect(registerResult.error?.message, 'Email verification is required');
        expect(registrationRepository.registerRequests, isEmpty);
      },
    );

    test('confirm tokens stores verification tokens only in memory', () async {
      final contactVerificationRepository =
          _FakeContactVerificationRepository();
      final draftRepository = _FakeRegisterDraftRepository();
      final usecase = _usecase(
        contactVerificationRepository: contactVerificationRepository,
        draftRepository: draftRepository,
      );
      usecase.startEmptyRegisterState();
      await usecase.submitAndRequestEmailToken('maria@example.com');
      await usecase.submitAndRequestPhoneToken('(27) 99999-9999');

      final emailResult = await usecase.confirmEmailToken('111111');
      final phoneResult = await usecase.confirmPhoneToken('222222');

      expect(emailResult, isA<Success<Unit>>());
      expect(phoneResult, isA<Success<Unit>>());
      expect(usecase.state?.isEmailVerified, isTrue);
      expect(usecase.state?.isPhoneVerified, isTrue);
      expect(
        draftRepository.savedSnapshots.last.toMap(),
        isNot(contains('email_verification_token')),
      );
      expect(
        draftRepository.savedSnapshots.last.toMap(),
        isNot(contains('phone_verification_token')),
      );
    });

    test(
      'register sends complete request, deletes draft, and clears memory',
      () async {
        final registrationRepository = _FakeRegistrationRepository();
        final contactVerificationRepository =
            _FakeContactVerificationRepository();
        final draftRepository = _FakeRegisterDraftRepository();
        final usecase = _usecase(
          registrationRepository: registrationRepository,
          contactVerificationRepository: contactVerificationRepository,
          draftRepository: draftRepository,
        );
        await _prepareValidRegistration(usecase);

        final result = await usecase.register();

        expect(result, isA<Success<Unit>>());
        expect(registrationRepository.registerRequests, hasLength(1));
        final request = registrationRepository.registerRequests.single;
        expect(request.cpf, '12345678909');
        expect(request.name, 'Maria Silva');
        expect(request.email, 'maria@example.com');
        expect(request.phone, '(27) 99999-9999');
        expect(request.password, 'secret123');
        expect(request.emailVerificationToken, 'email-verification-token');
        expect(request.phoneVerificationToken, 'phone-verification-token');
        expect(draftRepository.deletedCpfs, ['12345678909']);
        expect(usecase.state, isNull);
      },
    );

    test(
      'register succeeds and clears memory even when draft delete fails',
      () async {
        final registrationRepository = _FakeRegistrationRepository();
        final contactVerificationRepository =
            _FakeContactVerificationRepository();
        final draftRepository = _FakeRegisterDraftRepository(
          deleteResult: const Failure<Unit>(
            AppError(code: AppErrorCode.unexpected, message: 'delete failed'),
          ),
        );
        final usecase = _usecase(
          registrationRepository: registrationRepository,
          contactVerificationRepository: contactVerificationRepository,
          draftRepository: draftRepository,
        );
        await _prepareValidRegistration(usecase);

        final result = await usecase.register();

        expect(result, isA<Success<Unit>>());
        expect(registrationRepository.registerRequests, hasLength(1));
        expect(usecase.state, isNull);
      },
    );
  });
}

RegisterUsecase _usecase({
  _FakeRegistrationRepository? registrationRepository,
  _FakeContactVerificationRepository? contactVerificationRepository,
  _FakeRegisterDraftRepository? draftRepository,
}) {
  return RegisterUsecase(
    registrationRepository:
        registrationRepository ?? _FakeRegistrationRepository(),
    contactVerificationRepository:
        contactVerificationRepository ?? _FakeContactVerificationRepository(),
    registerDraftRepository: draftRepository ?? _FakeRegisterDraftRepository(),
  );
}

Future<void> _prepareValidRegistration(RegisterUsecase usecase) async {
  usecase.startEmptyRegisterState();
  await usecase.submitCPF('123.456.789-09');
  await usecase.submitName('Maria Silva');
  await usecase.submitBirthDate(DateTime(1990, 1, 15));
  await usecase.submitAndRequestEmailToken('maria@example.com');
  await usecase.confirmEmailToken('111111');
  await usecase.submitAndRequestPhoneToken('(27) 99999-9999');
  await usecase.confirmPhoneToken('222222');
  await usecase.submitPassword(('secret123', 'secret123'));
}

RegisterDraftSnapshot _snapshot() {
  return RegisterDraftSnapshot(
    cpf: '123.456.789-09',
    name: 'Maria Silva',
    birthDate: DateTime(1990, 1, 15),
    email: 'maria@example.com',
    phone: '(27) 99999-9999',
    emailVerificationId: 'email-verification-id',
    phoneVerificationId: 'phone-verification-id',
    isEmailVerified: true,
    isPhoneVerified: false,
    createdAt: DateTime.utc(2026, 5, 19, 10),
    updatedAt: DateTime.utc(2026, 5, 19, 11),
  );
}

class _FakeRegisterDraftRepository implements RegisterDraftRepository {
  final Result<RegisterDraftSnapshot> loadResult;
  final Result<Unit> saveResult;
  final Result<Unit> deleteResult;

  @override
  RegisterDraftSnapshot? snapshot;

  final List<String> loadedCpfs = [];
  final List<RegisterDraftSnapshot> savedSnapshots = [];
  final List<String> deletedCpfs = [];

  _FakeRegisterDraftRepository({
    Result<RegisterDraftSnapshot>? loadResult,
    Result<Unit>? saveResult,
    Result<Unit>? deleteResult,
  }) : loadResult = loadResult ?? Success<RegisterDraftSnapshot>(_snapshot()),
       saveResult = saveResult ?? const Success<Unit>(unit),
       deleteResult = deleteResult ?? const Success<Unit>(unit);

  @override
  AsyncResult<RegisterDraftSnapshot> getByCPF(String cpf) async {
    loadedCpfs.add(cpf);
    if (loadResult.isSuccess) snapshot = loadResult.value;
    return loadResult;
  }

  @override
  AsyncResult<Unit> save(RegisterDraftSnapshot snapshot) async {
    savedSnapshots.add(snapshot);
    if (saveResult.isSuccess) this.snapshot = snapshot;
    return saveResult;
  }

  @override
  AsyncResult<Unit> deleteByCPF(String cpf) async {
    deletedCpfs.add(cpf);
    if (deleteResult.isSuccess) snapshot = null;
    return deleteResult;
  }
}

class _FakeRegistrationRepository implements RegistrationRepository {
  Result<CpfCheckResponseDto> cpfCheckResult;
  Result<Unit> registerResult;

  final List<String> checkedCpfs = [];
  final List<RegisterRequestDto> registerRequests = [];

  _FakeRegistrationRepository({
    Result<CpfCheckResponseDto>? cpfCheckResult,
    Result<Unit>? registerResult,
  }) : cpfCheckResult =
           cpfCheckResult ??
           Success(
             CpfCheckResponseDto(
               cpf: '12345678909',
               exists: false,
               available: true,
             ),
           ),
       registerResult = registerResult ?? const Success<Unit>(unit);

  @override
  AsyncResult<CpfCheckResponseDto> cpfCheck(String cpf) async {
    checkedCpfs.add(cpf);
    return cpfCheckResult;
  }

  @override
  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    registerRequests.add(dto);
    return registerResult;
  }
}

class _FakeContactVerificationRepository
    implements ContactVerificationRepository {
  Result<ContactVerificationRequestResponseDto> emailRequestResult;
  Result<ContactVerificationRequestResponseDto> phoneRequestResult;
  Result<ContactVerificationConfirmResponseDto> emailConfirmResult;
  Result<ContactVerificationConfirmResponseDto> phoneConfirmResult;

  final List<ContactVerificationRequestDto> verificationRequests = [];
  final List<ContactVerificationConfirmRequestDto> confirmationRequests = [];

  _FakeContactVerificationRepository({
    Result<ContactVerificationRequestResponseDto>? emailRequestResult,
    Result<ContactVerificationRequestResponseDto>? phoneRequestResult,
    Result<ContactVerificationConfirmResponseDto>? emailConfirmResult,
    Result<ContactVerificationConfirmResponseDto>? phoneConfirmResult,
  }) : emailRequestResult =
           emailRequestResult ??
           Success(
             ContactVerificationRequestResponseDto(
               verificationId: 'email-verification-id',
               channel: 'email',
               target: 'maria@example.com',
               token: '111111',
               expiresAt: DateTime.utc(2026, 5, 19, 12, 5),
             ),
           ),
       phoneRequestResult =
           phoneRequestResult ??
           Success(
             ContactVerificationRequestResponseDto(
               verificationId: 'phone-verification-id',
               channel: 'phone',
               target: '(27) 99999-9999',
               token: '222222',
               expiresAt: DateTime.utc(2026, 5, 19, 12, 5),
             ),
           ),
       emailConfirmResult =
           emailConfirmResult ??
           Success(
             ContactVerificationConfirmResponseDto(
               verificationToken: 'email-verification-token',
               channel: 'email',
               target: 'maria@example.com',
               verifiedAt: DateTime.utc(2026, 5, 19, 12),
             ),
           ),
       phoneConfirmResult =
           phoneConfirmResult ??
           Success(
             ContactVerificationConfirmResponseDto(
               verificationToken: 'phone-verification-token',
               channel: 'phone',
               target: '(27) 99999-9999',
               verifiedAt: DateTime.utc(2026, 5, 19, 12),
             ),
           );

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    verificationRequests.add(dto);
    if (dto.channel == 'email') return emailRequestResult;
    return phoneRequestResult;
  }

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    confirmationRequests.add(dto);
    if (dto.verificationId == 'email-verification-id') {
      return emailConfirmResult;
    }
    return phoneConfirmResult;
  }
}
