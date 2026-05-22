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

  /// Creates a new draft state initialized with the provided [now] timestamp.
  RegisterDraftState({
    DateTime? now,
  }) : this._withNow((now ?? DateTime.now()).toUtc());

  /// Internal constructor used to initialize creation and update timestamps.
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

  /// Read-only view of fields changed since the last clean mark.
  Set<RegisterDraftField> get dirtyFields => Set.unmodifiable(_dirtyFields);

  /// Whether at least one field has changed since the last clean mark.
  bool get isDirty => _dirtyFields.isNotEmpty;

  set cpf(String value) {
    _set(
      RegisterDraftField.cpf,
      _cpf,
      value.onlyNumbers,
      (next) => _cpf = next,
    );
  }

  set name(String? value) {
    _set(
      RegisterDraftField.name,
      _name,
      value?.trimToNull(),
      (next) => _name = next,
    );
  }

  set birthDate(DateTime? value) {
    _set(
      RegisterDraftField.birthDate,
      _birthDate,
      value,
      (next) => _birthDate = next,
    );
  }

  set email(String? value) {
    _set(
      RegisterDraftField.email,
      _email,
      value?.trimToNull(),
      (next) => _email = next,
    );
  }

  set phone(String? value) {
    _set(
      RegisterDraftField.phone,
      _phone,
      value?.trimToNull(),
      (next) => _phone = next,
    );
  }

  set emailVerificationId(String? value) {
    _set(
      RegisterDraftField.emailVerificationId,
      _emailVerificationId,
      value?.trimToNull(),
      (next) => _emailVerificationId = next,
    );
  }

  set phoneVerificationId(String? value) {
    _set(
      RegisterDraftField.phoneVerificationId,
      _phoneVerificationId,
      value?.trimToNull(),
      (next) => _phoneVerificationId = next,
    );
  }

  set isEmailVerified(bool value) {
    _set(
      RegisterDraftField.isEmailVerified,
      _isEmailVerified,
      value,
      (next) => _isEmailVerified = next,
    );
  }

  set isPhoneVerified(bool value) {
    _set(
      RegisterDraftField.isPhoneVerified,
      _isPhoneVerified,
      value,
      (next) => _isPhoneVerified = next,
    );
  }

  /// Converts the current mutable state into an immutable snapshot.
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

  /// Recreates a mutable state from an existing [snapshot].
  static RegisterDraftState fromSnapshot(RegisterDraftSnapshot snapshot) {
    final state = RegisterDraftState();
    state.hydrate(snapshot);
    return state;
  }

  /// Replaces the current in-memory values with data from [snapshot].
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

  /// Clears dirty tracking for all fields.
  void markClean() {
    _dirtyFields.clear();
  }

  /// Marks this draft as persisted at [persistedAt] and clears dirty fields.
  void markPersisted(DateTime persistedAt) {
    _updatedAt = persistedAt.toUtc();
    markClean();
  }

  /// Marks [field] as dirty and refreshes the update timestamp.
  void _markDirty(RegisterDraftField field) {
    _dirtyFields.add(field);
    _updatedAt = DateTime.now().toUtc();
  }

  /// Applies [next] to [field] only when it differs from [current].
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
