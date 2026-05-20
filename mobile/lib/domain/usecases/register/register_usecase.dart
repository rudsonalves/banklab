import '/core/result/result.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/repositories/register_draft/register_draft_repository.dart';
import '/data/services/auth/api/dtos/contact_verification_confirm_request_dto.dart';
import '/data/services/auth/api/dtos/contact_verification_request_dto.dart';
import '/data/services/auth/api/dtos/register_request_dto.dart';
import '/data/services/auth/cache/register_draft/register_draft_load_result.dart';
import '/domain/common/auth/models/register_draft_snapshot.dart';
import '/domain/common/auth/models/register_draft_state.dart';

class RegisterUsecase {
  final AuthRepository _authRepository;
  final RegisterDraftRepository _registerDraftRepository;

  RegisterUsecase({
    required AuthRepository authRepository,
    required RegisterDraftRepository registerDraftRepository,
  }) : _authRepository = authRepository,
       _registerDraftRepository = registerDraftRepository;

  // State Management
  RegisterDraftSnapshot? get draft => _registerDraftRepository.snapshot;

  RegisterDraftState? _state;
  RegisterDraftState? get state => _state;
  bool _started = false;
  String? _emailVerificationToken;
  String? _phoneVerificationToken;
  String? _password;

  AsyncResult<Unit> initialize(String? cpf) async {
    if (_started) return const Success(unit);

    if (cpf != null) {
      final result = await _registerDraftRepository.getByCPF(cpf);

      if (result.isSuccess) {
        final loadResult = result.value!;
        if (loadResult.isFound) {
          _state = RegisterDraftState.fromSnapshot(
            (loadResult as RegisterDraftFound).snapshot,
          );
          _started = true;
          return const Success(unit);
        }
      } else {
        return Failure(result.error!);
      }
    }

    _state = RegisterDraftState();
    _started = true;
    return const Success(unit);
  }

  AsyncResult<Unit> submitCPF(String cpf) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    final result = await _authRepository.cpfCheck(cpf);
    if (result.isFailure) return Result.failure(result.error!);
    final cpfCheck = result.value!;
    if (!cpfCheck.avaliable) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'CPF is already registered',
        ),
      );
    }

    _state!.updateCPF(cpf);
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> submitName(String name) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    _state!.updateName(name);
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> submitBirthDate(DateTime birthDate) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    _state!.updateBirthDate(birthDate);
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> submitAndRequestEmailToken(String email) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    final result = await _authRepository.requestContactVerification(
      ContactVerificationRequestDto(
        channel: 'email',
        target: email,
      ),
    );
    if (result.isFailure) return Result.failure(result.error!);

    _state!.updateEmail(email);
    _state!.updateEmailVerified(false);
    _state!.updateEmailVerificationId(result.value!.verificationId);
    _emailVerificationToken = null;
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> submitAndRequestPhoneToken(String phone) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    final result = await _authRepository.requestContactVerification(
      ContactVerificationRequestDto(
        channel: 'phone',
        target: phone,
      ),
    );
    if (result.isFailure) return Result.failure(result.error!);

    _state!.updatePhone(phone);
    _state!.updatePhoneVerified(false);
    _state!.updatePhoneVerificationId(result.value!.verificationId);
    _phoneVerificationToken = null;
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> confirmEmailToken(String token) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    if (_state!.emailVerificationId == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Email verification not requested',
        ),
      );
    }

    final result = await _authRepository.confirmContactVerification(
      ContactVerificationConfirmRequestDto(
        verificationId: _state!.emailVerificationId!,
        token: token,
      ),
    );
    if (result.isFailure) return Result.failure(result.error!);

    _emailVerificationToken = result.value!.verificationToken;
    _state!.updateEmailVerified(true);
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> confirmPhoneToken(String token) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    if (_state!.phoneVerificationId == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Phone verification not requested',
        ),
      );
    }

    final result = await _authRepository.confirmContactVerification(
      ContactVerificationConfirmRequestDto(
        verificationId: _state!.phoneVerificationId!,
        token: token,
      ),
    );
    if (result.isFailure) return Result.failure(result.error!);

    _phoneVerificationToken = result.value!.verificationToken;
    _state!.updatePhoneVerified(true);
    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> submitPassword((String, String) passWd) async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    final (password, confirmPassword) = passWd;
    if (password.trim().isEmpty || password.trim() != confirmPassword) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Password is empty or does not match confirmation',
        ),
      );
    }
    _password = password.trim();

    final saveResult = await _registerDraftRepository.save(
      _state!.toSnapshot(),
    );
    if (saveResult.isFailure) return Result.failure(saveResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> register() async {
    if (_state == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.unexpected,
          message: 'Usecase not initialized',
        ),
      );
    }

    if (_state!.cpf.isEmpty) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'CPF is required',
        ),
      );
    }

    if (_password == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Password is required',
        ),
      );
    }

    if (_emailVerificationToken == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Email verification is required',
        ),
      );
    }

    if (_phoneVerificationToken == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'Phone verification is required',
        ),
      );
    }

    if (_state!.name == null ||
        _state!.birthDate == null ||
        _state!.email == null ||
        _state!.phone == null) {
      return const Failure(
        AppError(
          code: AppErrorCode.invalidData,
          message: 'All fields are required',
        ),
      );
    }

    final request = RegisterRequestDto(
      cpf: _state!.cpf,
      name: _state!.name!,
      birthDate: _state!.birthDate!,
      email: _state!.email!,
      phone: _state!.phone!,
      password: _password!,
      emailVerificationToken: _emailVerificationToken!,
      phoneVerificationToken: _phoneVerificationToken!,
    );

    final result = await _authRepository.register(request);
    if (result.isFailure) return Result.failure(result.error!);

    await _deleteDraft();
    _clearMemory();

    return const Success(unit);
  }

  AsyncResult<Unit> reset() async {
    final deleteResult = await _deleteDraft();
    _clearMemory();
    if (deleteResult.isFailure) return Result.failure(deleteResult.error!);

    return const Success(unit);
  }

  AsyncResult<Unit> _deleteDraft() async {
    if (_state == null) {
      return const Success(unit);
    }

    if (_state!.cpf.isEmpty) {
      return const Success(unit);
    }

    final deleteResult = await _registerDraftRepository.deleteByCPF(
      _state!.cpf,
    );
    if (deleteResult.isFailure) return Result.failure(deleteResult.error!);

    return const Success(unit);
  }

  void _clearMemory() {
    _state = null;
    _started = false;
    _emailVerificationToken = null;
    _phoneVerificationToken = null;
    _password = null;
  }
}
