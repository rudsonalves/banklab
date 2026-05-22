class PasswordModel {
  final String _value;
  final String _confirm;

  static const int minLength = 8;

  PasswordModel({String? value, String? confirm})
    : _value = value?.trim() ?? '',
      _confirm = confirm?.trim() ?? '';

  String get value => _value;
  String get confirm => _confirm;

  bool get hasNumber => RegExp(r'\d').hasMatch(_value);
  bool get hasUppercase => RegExp(r'[A-Z]').hasMatch(_value);
  bool get hasLowercase => RegExp(r'[a-z]').hasMatch(_value);
  bool get hasMinLength => _value.length >= minLength;

  bool get isValidPassword =>
      hasNumber && hasUppercase && hasLowercase && hasMinLength;

  bool get hasEquals => _value.isNotEmpty && _value == _confirm;

  bool get isValid => isValidPassword && hasEquals;
}
