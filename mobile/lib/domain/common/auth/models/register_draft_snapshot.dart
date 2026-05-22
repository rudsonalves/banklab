import '/core/extensions/datetime_extension.dart';
import '/core/extensions/string.dart';

class RegisterDraftSnapshot {
  final String cpf;
  final String? name;
  final DateTime? birthDate;
  final String? email;
  final String? phone;
  final String? emailVerificationId;
  final String? phoneVerificationId;
  final bool isEmailVerified;
  final bool isPhoneVerified;
  final DateTime createdAt;
  final DateTime updatedAt;

  const RegisterDraftSnapshot({
    required this.cpf,
    required this.isEmailVerified,
    required this.isPhoneVerified,
    required this.createdAt,
    required this.updatedAt,
    this.name,
    this.birthDate,
    this.email,
    this.phone,
    this.emailVerificationId,
    this.phoneVerificationId,
  });

  factory RegisterDraftSnapshot.empty(String cpf) {
    final now = DateTime.now().toUtc();
    return RegisterDraftSnapshot(
      cpf: cpf,
      isEmailVerified: false,
      isPhoneVerified: false,
      createdAt: now,
      updatedAt: now,
    );
  }

  Map<String, dynamic> toMap() {
    return {
      'cpf': cpf.onlyNumbers,
      'name': name,
      'birth_date': birthDate?.dateOnly,
      'email': email,
      'phone': phone,
      'email_verification_id': emailVerificationId,
      'phone_verification_id': phoneVerificationId,
      'is_email_verified': isEmailVerified,
      'is_phone_verified': isPhoneVerified,
      'created_at': createdAt.toUtc().toIso8601String(),
      'updated_at': updatedAt.toUtc().toIso8601String(),
    };
  }

  static RegisterDraftSnapshot? fromMapOrNull(Map<String, dynamic>? map) {
    if (map == null) return null;

    try {
      final cpf = (map['cpf'] as String?)?.onlyNumbers ?? '';
      final createdAt = DateTime.tryParse(map['created_at'] as String? ?? '');
      final updatedAt = DateTime.tryParse(map['updated_at'] as String? ?? '');

      if (cpf.isEmpty || createdAt == null || updatedAt == null) {
        return null;
      }

      return RegisterDraftSnapshot(
        cpf: cpf,
        name: (map['name'] as String?)?.trimToNull(),
        birthDate: DateParser.parseDateOnly(map['birth_date'] as String?),
        email: (map['email'] as String?)?.trimToNull(),
        phone: (map['phone'] as String?)?.trimToNull(),
        emailVerificationId: (map['email_verification_id'] as String?)
            ?.trimToNull(),
        phoneVerificationId: (map['phone_verification_id'] as String?)
            ?.trimToNull(),
        isEmailVerified: map['is_email_verified'] == true,
        isPhoneVerified: map['is_phone_verified'] == true,
        createdAt: createdAt,
        updatedAt: updatedAt,
      );
    } catch (_) {
      return null;
    }
  }
}
