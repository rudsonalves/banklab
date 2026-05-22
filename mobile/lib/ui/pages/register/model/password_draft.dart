import '/domain/usecases/register/models/password_model.dart';

class PasswordDraft {
  String? value;
  String? confirm;

  PasswordDraft({
    this.value,
    this.confirm,
  });

  static const int minLength = PasswordModel.minLength;

  PasswordModel get model => PasswordModel(
    value: value,
    confirm: confirm,
  );

  // Password rules
  bool get hasNumber => model.hasNumber;
  bool get hasUppercase => model.hasUppercase;
  bool get hasLowercase => model.hasLowercase;
  bool get hasMinLength => model.hasMinLength;

  bool get isValidPassword => model.isValidPassword;

  bool get hasEquals => model.hasEquals;

  bool get isValid => model.isValid;
}
