class PasswordModel {
  final String _value;
  final String _confirm;

  PasswordModel(String value, String confirm)
    : _value = value.trim(),
      _confirm = confirm.trim();

  String get value => _value;
  String get confirm => _confirm;

  bool get isValid {
    return _value.isNotEmpty && _value == _confirm;
  }
}
