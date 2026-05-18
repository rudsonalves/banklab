import 'package:flutter/foundation.dart';

import '/core/extensions/string.dart';
import '/core/result/command.dart';
import '/data/repositories/auth/auth_repository.dart';
import '/data/services/auth/api/dtos/contact_verification_confirm_request_dto.dart';
import '/data/services/auth/api/dtos/contact_verification_request_dto.dart';
import '/data/services/auth/api/dtos/contact_verification_request_response_dto.dart';
import '../../../../../data/services/auth/api/dtos/register_request_dto.dart';

enum RegisterStep {
  personalData,
  contactData,
  emailVerification,
  phoneVerification,
  review,
}

class RegisterViewmodel extends ChangeNotifier {
  final AuthRepository _authRepository;

  RegisterStep _currentStep = RegisterStep.personalData;

  String _name = '';
  String _cpf = '';
  DateTime? _birthDate;
  String _email = '';
  String _phone = '';
  String _password = '';

  String? _emailVerificationId;
  String? _emailVerificationToken;
  String? _phoneVerificationId;
  String? _phoneVerificationToken;

  AppError? _stepError;
  RegisterStep? _stepErrorStep;

  RegisterViewmodel({
    required AuthRepository authRepository,
  }) : _authRepository = authRepository {
    requestEmailCode = Command0(_requestEmailCode);
    confirmEmailCode = Command1(_confirmEmailCode);
    requestPhoneCode = Command0(_requestPhoneCode);
    confirmPhoneCode = Command1(_confirmPhoneCode);
    register = Command0(_register);
  }

  RegisterStep get currentStep => _currentStep;

  String get name => _name;
  String get cpf => _cpf;
  DateTime? get birthDate => _birthDate;
  String get email => _email;
  String get phone => _phone;
  String get password => _password;

  String? get emailVerificationId => _emailVerificationId;
  String? get emailVerificationToken => _emailVerificationToken;
  String? get phoneVerificationId => _phoneVerificationId;
  String? get phoneVerificationToken => _phoneVerificationToken;

  bool get isEmailVerified =>
      _emailVerificationToken != null && _emailVerificationToken!.isNotEmpty;
  bool get isPhoneVerified =>
      _phoneVerificationToken != null && _phoneVerificationToken!.isNotEmpty;

  AppError? get stepError => _stepError;
  RegisterStep? get stepErrorStep => _stepErrorStep;

  bool get canAdvanceCurrentStep {
    return switch (_currentStep) {
      RegisterStep.personalData => _hasValidPersonalData,
      RegisterStep.contactData => _hasValidContactData,
      RegisterStep.emailVerification => isEmailVerified,
      RegisterStep.phoneVerification => isPhoneVerified,
      RegisterStep.review => false,
    };
  }

  bool get canRegister =>
      _hasValidPersonalData &&
      _hasValidContactData &&
      isEmailVerified &&
      isPhoneVerified;

  late final Command0<ContactVerificationRequestResponseDto> requestEmailCode;
  late final Command1<Unit, String> confirmEmailCode;
  late final Command0<ContactVerificationRequestResponseDto> requestPhoneCode;
  late final Command1<Unit, String> confirmPhoneCode;
  late final Command0<Unit> register;

  void updatePersonalData({
    String? name,
    String? cpf,
    DateTime? birthDate,
    String? password,
  }) {
    if (name != null) _name = name.trim();
    if (cpf != null) _cpf = cpf.onlyNumbers;
    if (birthDate != null) _birthDate = birthDate;
    if (password != null) _password = password;
    _clearStepError();
    notifyListeners();
  }

  void updateContactData({
    String? email,
    String? phone,
  }) {
    if (email != null) _email = email.trim();
    if (phone != null) _phone = phone.trim();
    _clearStepError();
    notifyListeners();
  }

  bool nextStep() {
    if (!canAdvanceCurrentStep) {
      _setValidationErrorForCurrentStep();
      return false;
    }

    _currentStep = switch (_currentStep) {
      RegisterStep.personalData => RegisterStep.contactData,
      RegisterStep.contactData => RegisterStep.emailVerification,
      RegisterStep.emailVerification => RegisterStep.phoneVerification,
      RegisterStep.phoneVerification => RegisterStep.review,
      RegisterStep.review => RegisterStep.review,
    };

    _clearStepError();
    notifyListeners();
    return true;
  }

  bool previousStep() {
    final previous = switch (_currentStep) {
      RegisterStep.personalData => null,
      RegisterStep.contactData => RegisterStep.personalData,
      RegisterStep.emailVerification => RegisterStep.contactData,
      RegisterStep.phoneVerification => RegisterStep.emailVerification,
      RegisterStep.review => RegisterStep.phoneVerification,
    };

    if (previous == null) return false;

    _currentStep = previous;
    _clearStepError();
    notifyListeners();
    return true;
  }

  void goToStep(RegisterStep step) {
    _currentStep = step;
    _clearStepError();
    notifyListeners();
  }

  AsyncResult<ContactVerificationRequestResponseDto> _requestEmailCode() async {
    _currentStep = RegisterStep.emailVerification;

    if (_email.isEmpty || !_isValidEmail(_email)) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message: 'Informe um e-mail valido para verificacao.',
      );
      _setStepError(error, RegisterStep.emailVerification);
      return Failure(error);
    }

    final result = await _authRepository.requestContactVerification(
      ContactVerificationRequestDto(
        channel: 'email',
        target: _email,
      ),
    );

    if (result.isFailure) {
      _setStepError(result.error!, RegisterStep.emailVerification);
      return Result.failure(result.error!);
    }

    final data = result.value!;
    _emailVerificationId = data.verificationId;
    _emailVerificationToken = null;
    _clearStepError();
    notifyListeners();
    return Success(data);
  }

  AsyncResult<Unit> _confirmEmailCode(String code) async {
    _currentStep = RegisterStep.emailVerification;

    if (_emailVerificationId == null || _emailVerificationId!.isEmpty) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message: 'Solicite o codigo de e-mail antes de confirmar.',
      );
      _setStepError(error, RegisterStep.emailVerification);
      return Failure(error);
    }

    final token = code.trim();
    if (token.isEmpty) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message: 'Informe o codigo de verificacao do e-mail.',
      );
      _setStepError(error, RegisterStep.emailVerification);
      return Failure(error);
    }

    final result = await _authRepository.confirmContactVerification(
      ContactVerificationConfirmRequestDto(
        verificationId: _emailVerificationId!,
        token: token,
      ),
    );

    if (result.isFailure) {
      _setStepError(result.error!, RegisterStep.emailVerification);
      return Result.failure(result.error!);
    }

    _emailVerificationToken = result.value!.verificationToken;
    _currentStep = RegisterStep.phoneVerification;
    _clearStepError();
    notifyListeners();
    return Success(unit);
  }

  AsyncResult<ContactVerificationRequestResponseDto> _requestPhoneCode() async {
    _currentStep = RegisterStep.phoneVerification;

    if (_phone.isEmpty) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message: 'Informe um telefone valido para verificacao.',
      );
      _setStepError(error, RegisterStep.phoneVerification);
      return Failure(error);
    }

    final result = await _authRepository.requestContactVerification(
      ContactVerificationRequestDto(
        channel: 'phone',
        target: _phone,
      ),
    );

    if (result.isFailure) {
      _setStepError(result.error!, RegisterStep.phoneVerification);
      return Result.failure(result.error!);
    }

    final data = result.value!;
    _phoneVerificationId = data.verificationId;
    _phoneVerificationToken = null;
    _clearStepError();
    notifyListeners();
    return Success(data);
  }

  AsyncResult<Unit> _confirmPhoneCode(String code) async {
    _currentStep = RegisterStep.phoneVerification;

    if (_phoneVerificationId == null || _phoneVerificationId!.isEmpty) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message: 'Solicite o codigo de telefone antes de confirmar.',
      );
      _setStepError(error, RegisterStep.phoneVerification);
      return Failure(error);
    }

    final token = code.trim();
    if (token.isEmpty) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message: 'Informe o codigo de verificacao do telefone.',
      );
      _setStepError(error, RegisterStep.phoneVerification);
      return Failure(error);
    }

    final result = await _authRepository.confirmContactVerification(
      ContactVerificationConfirmRequestDto(
        verificationId: _phoneVerificationId!,
        token: token,
      ),
    );

    if (result.isFailure) {
      _setStepError(result.error!, RegisterStep.phoneVerification);
      return Result.failure(result.error!);
    }

    _phoneVerificationToken = result.value!.verificationToken;
    _currentStep = RegisterStep.review;
    _clearStepError();
    notifyListeners();
    return Success(unit);
  }

  AsyncResult<Unit> _register() async {
    if (!canRegister) {
      final error = AppError(
        code: AppErrorCode.invalidData,
        message:
            'Preencha os dados obrigatorios e confirme e-mail e telefone antes de concluir o cadastro.',
      );
      _setStepError(error, _currentStep);
      return Failure(error);
    }

    final result = await _authRepository.register(
      RegisterRequestDto(
        name: _name,
        email: _email,
        phone: _phone,
        password: _password,
        birthDate: _birthDate!,
        cpf: _cpf,
        emailVerificationToken: _emailVerificationToken!,
        phoneVerificationToken: _phoneVerificationToken!,
      ),
    );

    if (result.isFailure) {
      _setStepError(result.error!, RegisterStep.review);
      return Result.failure(result.error!);
    }

    _currentStep = RegisterStep.review;
    _clearStepError();
    notifyListeners();
    return Success(unit);
  }

  bool get _hasValidPersonalData {
    return _name.trim().length >= 3 &&
        _cpf.onlyNumbers.length == 11 &&
        _birthDate != null &&
        _password.length >= 6;
  }

  bool get _hasValidContactData {
    return _isValidEmail(_email) && _phone.trim().isNotEmpty;
  }

  bool _isValidEmail(String value) {
    return RegExp(r'^[^@\s]+@[^@\s]+\.[^@\s]+$').hasMatch(value.trim());
  }

  void _setValidationErrorForCurrentStep() {
    final message = switch (_currentStep) {
      RegisterStep.personalData =>
        'Preencha nome, CPF, data de nascimento e senha validos para avancar.',
      RegisterStep.contactData =>
        'Preencha e-mail e telefone validos para avancar.',
      RegisterStep.emailVerification =>
        'Confirme o e-mail para avancar para a proxima etapa.',
      RegisterStep.phoneVerification =>
        'Confirme o telefone para avancar para a proxima etapa.',
      RegisterStep.review => 'Etapa de revisao nao possui avancar.',
    };

    _setStepError(
      AppError(
        code: AppErrorCode.invalidData,
        message: message,
      ),
      _currentStep,
    );
  }

  void _setStepError(AppError error, RegisterStep step) {
    _stepError = error;
    _stepErrorStep = step;
    notifyListeners();
  }

  void _clearStepError() {
    _stepError = null;
    _stepErrorStep = null;
  }
}
