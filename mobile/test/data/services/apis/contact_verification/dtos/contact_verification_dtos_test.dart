import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/enums/contact_verification_channel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('ContactVerificationRequestDto', () {
    test('serializes request body with channel and target', () {
      final dto = ContactVerificationRequestDto(
        channel: ContactVerificationChannel.email,
        target: 'user@example.com',
      );

      final map = dto.toMap();

      expect(map['channel'], 'email');
      expect(map['target'], 'user@example.com');
    });

    test('parses request map', () {
      final dto = ContactVerificationRequestDto.fromMap({
        'channel': 'phone',
        'target': '+5511999999999',
      });

      expect(dto.channel, ContactVerificationChannel.phone);
      expect(dto.target, '+5511999999999');
    });
  });

  group('ContactVerificationRequestResponseDto.fromMap', () {
    test('parses verification request response payload', () {
      final dto = ContactVerificationRequestResponseDto.fromMap({
        'verification_id': 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
        'channel': 'email',
        'target': 'user@example.com',
        'expires_at': '2026-05-18T12:10:00Z',
      });

      expect(dto.verificationId, 'a5d4f5f1-a1b0-4f58-9f74-123456789abc');
      expect(dto.channel, ContactVerificationChannel.email);
      expect(dto.target, 'user@example.com');
      // expect(dto.token, '123456');
      expect(dto.expiresAt, DateTime.parse('2026-05-18T12:10:00Z'));
    });
  });

  group('ContactVerificationConfirmRequestDto', () {
    test('serializes confirm request body with verification id and token', () {
      final dto = ContactVerificationConfirmRequestDto(
        verificationId: 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
        token: '123456',
      );

      final map = dto.toMap();

      expect(map['verification_id'], 'a5d4f5f1-a1b0-4f58-9f74-123456789abc');
      expect(map['token'], '123456');
    });

    test('parses confirm request map', () {
      final dto = ContactVerificationConfirmRequestDto.fromMap({
        'verification_id': 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
        'token': '123456',
      });

      expect(dto.verificationId, 'a5d4f5f1-a1b0-4f58-9f74-123456789abc');
      expect(dto.token, '123456');
    });
  });

  group('ContactVerificationConfirmResponseDto.fromMap', () {
    test('parses verification confirm response payload', () {
      final dto = ContactVerificationConfirmResponseDto.fromMap({
        'verification_token': 'token-confirmado',
        'channel': 'phone',
        'target': '+5511999999999',
        'verified_at': '2026-05-18T12:03:00Z',
      });

      expect(dto.verificationToken, 'token-confirmado');
      expect(dto.channel, ContactVerificationChannel.phone);
      expect(dto.target, '+5511999999999');
      expect(dto.verifiedAt, DateTime.parse('2026-05-18T12:03:00Z'));
    });
  });
}
