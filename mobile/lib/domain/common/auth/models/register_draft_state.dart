import '/core/extensions/string.dart';
import '../enums/register_draft_field.dart';
import 'register_draft_snapshot.dart';

class RegisterDraftState {
  String _cpf = '';
  String? _name;
  DateTime? _birthDate;
  String? _email;
  String? _phone;
  String? _emailVerificationId;
  String? _phoneVerificationId;
  bool _isEmailVerified = false;
  bool _isPhoneVerified = false;
  DateTime _createdAt;
  DateTime _updatedAt;

  final Set<RegisterDraftField> _dirtyFields = {};

  RegisterDraftState({
    DateTime? now,
  }) : this._withNow((now ?? DateTime.now()).toUtc());

  RegisterDraftState._withNow(DateTime now)
    : _createdAt = now,
      _updatedAt = now;

  String get cpf => _cpf;
  String? get name => _name;
  DateTime? get birthDate => _birthDate;
  String? get email => _email;
  String? get phone => _phone;
  String? get emailVerificationId => _emailVerificationId;
  String? get phoneVerificationId => _phoneVerificationId;
  bool get isEmailVerified => _isEmailVerified;
  bool get isPhoneVerified => _isPhoneVerified;
  DateTime get createdAt => _createdAt;
  DateTime get updatedAt => _updatedAt;
  Set<RegisterDraftField> get dirtyFields => Set.unmodifiable(_dirtyFields);
  bool get isDirty => _dirtyFields.isNotEmpty;

  void updateCPF(String value) {
    _set(
      RegisterDraftField.cpf,
      _cpf,
      value.onlyNumbers,
      (next) => _cpf = next,
    );
  }

  void updateName(String? value) {
    _set(
      RegisterDraftField.name,
      _name,
      value?.trimToNull(),
      (next) => _name = next,
    );
  }

  void updateBirthDate(DateTime? value) {
    _set(
      RegisterDraftField.birthDate,
      _birthDate,
      value,
      (next) => _birthDate = next,
    );
  }

  void updateEmail(String? value) {
    _set(
      RegisterDraftField.email,
      _email,
      value?.trimToNull(),
      (next) => _email = next,
    );
  }

  void updatePhone(String? value) {
    _set(
      RegisterDraftField.phone,
      _phone,
      value?.trimToNull(),
      (next) => _phone = next,
    );
  }

  void updateEmailVerificationId(String? value) {
    _set(
      RegisterDraftField.emailVerificationId,
      _emailVerificationId,
      value?.trimToNull(),
      (next) => _emailVerificationId = next,
    );
  }

  void updatePhoneVerificationId(String? value) {
    _set(
      RegisterDraftField.phoneVerificationId,
      _phoneVerificationId,
      value?.trimToNull(),
      (next) => _phoneVerificationId = next,
    );
  }

  void updateEmailVerified(bool value) {
    _set(
      RegisterDraftField.isEmailVerified,
      _isEmailVerified,
      value,
      (next) => _isEmailVerified = next,
    );
  }

  void updatePhoneVerified(bool value) {
    _set(
      RegisterDraftField.isPhoneVerified,
      _isPhoneVerified,
      value,
      (next) => _isPhoneVerified = next,
    );
  }

  RegisterDraftSnapshot toSnapshot() {
    return RegisterDraftSnapshot(
      cpf: _cpf,
      name: _name,
      birthDate: _birthDate,
      email: _email,
      phone: _phone,
      emailVerificationId: _emailVerificationId,
      phoneVerificationId: _phoneVerificationId,
      isEmailVerified: _isEmailVerified,
      isPhoneVerified: _isPhoneVerified,
      createdAt: _createdAt,
      updatedAt: _updatedAt,
    );
  }

  static RegisterDraftState fromSnapshot(RegisterDraftSnapshot snapshot) {
    final state = RegisterDraftState();
    state.hydrate(snapshot);
    return state;
  }

  void hydrate(RegisterDraftSnapshot snapshot) {
    _cpf = snapshot.cpf.onlyNumbers;
    _name = snapshot.name;
    _birthDate = snapshot.birthDate;
    _email = snapshot.email;
    _phone = snapshot.phone;
    _emailVerificationId = snapshot.emailVerificationId;
    _phoneVerificationId = snapshot.phoneVerificationId;
    _isEmailVerified = snapshot.isEmailVerified;
    _isPhoneVerified = snapshot.isPhoneVerified;
    _createdAt = snapshot.createdAt;
    _updatedAt = snapshot.updatedAt;
    markClean();
  }

  void markClean() {
    _dirtyFields.clear();
  }

  void markPersisted(DateTime persistedAt) {
    _updatedAt = persistedAt.toUtc();
    markClean();
  }

  void _markDirty(RegisterDraftField field) {
    _dirtyFields.add(field);
    _updatedAt = DateTime.now().toUtc();
  }

  void _set<T>(
    RegisterDraftField field,
    T current,
    T next,
    void Function(T value) apply,
  ) {
    if (current == next) return;
    apply(next);
    _markDirty(field);
  }
}
