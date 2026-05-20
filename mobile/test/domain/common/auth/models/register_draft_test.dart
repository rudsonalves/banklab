import 'package:bankflow/domain/common/auth/models/register_draft.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RegisterDraftSnapshot', () {
    test('serializes persistable fields only', () {
      final snapshot = RegisterDraftSnapshot(
        cpf: '123.456.789-09',
        name: 'Maria Silva',
        birthDate: DateTime(1990, 1, 15),
        email: 'maria@example.com',
        phone: '+5527999999999',
        emailVerificationId: 'email-verification-id',
        phoneVerificationId: 'phone-verification-id',
        isEmailVerified: true,
        isPhoneVerified: false,
        createdAt: DateTime.utc(2026, 5, 19, 10),
        updatedAt: DateTime.utc(2026, 5, 19, 11),
      );

      final map = snapshot.toMap();

      expect(map['cpf'], '12345678909');
      expect(map['name'], 'Maria Silva');
      expect(map['birth_date'], '1990-01-15');
      expect(map['email'], 'maria@example.com');
      expect(map['phone'], '+5527999999999');
      expect(map['current_step'], 'phoneToken');
      expect(map['email_verification_id'], 'email-verification-id');
      expect(map['phone_verification_id'], 'phone-verification-id');
      expect(map['is_email_verified'], isTrue);
      expect(map['is_phone_verified'], isFalse);
      expect(map['created_at'], '2026-05-19T10:00:00.000Z');
      expect(map['updated_at'], '2026-05-19T11:00:00.000Z');
      expect(map.containsKey('password'), isFalse);
      expect(map.containsKey('password_confirmation'), isFalse);
      expect(map.containsKey('email_verification_token'), isFalse);
      expect(map.containsKey('phone_verification_token'), isFalse);
      expect(map.containsKey('access_token'), isFalse);
      expect(map.containsKey('refresh_token'), isFalse);
    });

    test('parses valid payload', () {
      final snapshot = RegisterDraftSnapshot.fromMapOrNull({
        'cpf': '123.456.789-09',
        'name': ' Maria Silva ',
        'birth_date': '1990-01-15',
        'email': ' maria@example.com ',
        'phone': ' +5527999999999 ',
        'current_step': 'emailToken',
        'email_verification_id': ' email-id ',
        'phone_verification_id': ' phone-id ',
        'is_email_verified': true,
        'is_phone_verified': false,
        'created_at': '2026-05-19T10:00:00.000Z',
        'updated_at': '2026-05-19T11:00:00.000Z',
      });

      expect(snapshot, isNotNull);
      expect(snapshot!.cpf, '12345678909');
      expect(snapshot.name, 'Maria Silva');
      expect(snapshot.birthDate, DateTime(1990, 1, 15));
      expect(snapshot.email, 'maria@example.com');
      expect(snapshot.phone, '+5527999999999');
      expect(snapshot.emailVerificationId, 'email-id');
      expect(snapshot.phoneVerificationId, 'phone-id');
      expect(snapshot.isEmailVerified, isTrue);
      expect(snapshot.isPhoneVerified, isFalse);
      expect(snapshot.createdAt, DateTime.parse('2026-05-19T10:00:00.000Z'));
      expect(snapshot.updatedAt, DateTime.parse('2026-05-19T11:00:00.000Z'));
    });

    test('returns null for invalid payload', () {
      expect(RegisterDraftSnapshot.fromMapOrNull(null), isNull);
      expect(RegisterDraftSnapshot.fromMapOrNull({}), isNull);
      expect(
        RegisterDraftSnapshot.fromMapOrNull({
          'cpf': '12345678909',
          'current_step': 'unknown',
          'created_at': '2026-05-19T10:00:00.000Z',
          'updated_at': '2026-05-19T11:00:00.000Z',
        }),
        isNull,
      );
      expect(
        RegisterDraftSnapshot.fromMapOrNull({
          'cpf': '12345678909',
          'current_step': 'cpf',
          'created_at': 'not-a-date',
          'updated_at': '2026-05-19T11:00:00.000Z',
        }),
        isNull,
      );
    });
  });

  group('RegisterDraftState', () {
    test('tracks dirty fields when values change', () {
      final state = RegisterDraftState(now: DateTime.utc(2026, 5, 19, 10));

      expect(state.isDirty, isFalse);

      state.updateCPF('123.456.789-09');
      state.updateName(' Maria Silva ');

      expect(state.cpf, '12345678909');
      expect(state.name, 'Maria Silva');
      expect(state.isDirty, isTrue);
      expect(
        state.dirtyFields,
        containsAll([
          RegisterDraftField.cpf,
          RegisterDraftField.name,
          RegisterDraftField.currentStep,
        ]),
      );
    });

    test('does not mark dirty when value does not change', () {
      final state = RegisterDraftState(now: DateTime.utc(2026, 5, 19, 10));

      state.updateCPF('12345678909');
      state.markClean();
      state.updateCPF('123.456.789-09');

      expect(state.isDirty, isFalse);
      expect(state.dirtyFields, isEmpty);
    });

    test('clears dirty tracking after markClean', () {
      final state = RegisterDraftState(now: DateTime.utc(2026, 5, 19, 10));

      state.updateEmail('maria@example.com');
      expect(state.isDirty, isTrue);

      state.markClean();

      expect(state.isDirty, isFalse);
      expect(state.dirtyFields, isEmpty);
    });

    test('hydrates from snapshot without dirty fields', () {
      final state = RegisterDraftState(now: DateTime.utc(2026, 5, 19, 9));
      final snapshot = RegisterDraftSnapshot(
        cpf: '12345678909',
        name: 'Maria Silva',
        birthDate: DateTime(1990, 1, 15),
        email: 'maria@example.com',
        phone: '+5527999999999',
        emailVerificationId: 'email-id',
        phoneVerificationId: 'phone-id',
        isEmailVerified: true,
        isPhoneVerified: false,
        createdAt: DateTime.utc(2026, 5, 19, 10),
        updatedAt: DateTime.utc(2026, 5, 19, 11),
      );

      state.updateCPF('00000000000');
      state.hydrate(snapshot);

      expect(state.cpf, '12345678909');
      expect(state.name, 'Maria Silva');
      expect(state.birthDate, DateTime(1990, 1, 15));
      expect(state.email, 'maria@example.com');
      expect(state.phone, '+5527999999999');
      expect(state.emailVerificationId, 'email-id');
      expect(state.phoneVerificationId, 'phone-id');
      expect(state.isEmailVerified, isTrue);
      expect(state.isPhoneVerified, isFalse);
      expect(state.createdAt, DateTime.utc(2026, 5, 19, 10));
      expect(state.updatedAt, DateTime.utc(2026, 5, 19, 11));
      expect(state.isDirty, isFalse);
    });

    test('creates persistable snapshot from current state', () {
      final state = RegisterDraftState(now: DateTime.utc(2026, 5, 19, 10));

      state.updateCPF('123.456.789-09');
      state.updateName('Maria Silva');
      state.updateBirthDate(DateTime(1990, 1, 15));
      state.updateEmail('maria@example.com');
      state.updatePhone('+5527999999999');
      state.updateEmailVerificationId('email-id');
      state.updatePhoneVerificationId('phone-id');
      state.updateEmailVerified(true);

      final snapshot = state.toSnapshot();

      expect(snapshot.cpf, '12345678909');
      expect(snapshot.name, 'Maria Silva');
      expect(snapshot.birthDate, DateTime(1990, 1, 15));
      expect(snapshot.email, 'maria@example.com');
      expect(snapshot.phone, '+5527999999999');
      expect(snapshot.emailVerificationId, 'email-id');
      expect(snapshot.phoneVerificationId, 'phone-id');
      expect(snapshot.isEmailVerified, isTrue);
      expect(snapshot.isPhoneVerified, isFalse);
      expect(snapshot.toMap().containsKey('password'), isFalse);
      expect(
        snapshot.toMap().containsKey('email_verification_token'),
        isFalse,
      );
    });

    test('updates persisted timestamp and clears dirty tracking', () {
      final state = RegisterDraftState(now: DateTime.utc(2026, 5, 19, 10));
      final persistedAt = DateTime.utc(2026, 5, 19, 12);

      state.updatePhone('+5527999999999');
      state.markPersisted(persistedAt);

      expect(state.updatedAt, persistedAt);
      expect(state.isDirty, isFalse);
    });
  });
}
