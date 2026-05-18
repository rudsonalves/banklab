import 'package:bankflow/data/services/apis/core/api_envelope.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('ApiError.fromMap', () {
    test('parses code, message, and details when details exist', () {
      final error = ApiError.fromMap({
        'code': 'CONTACT_NOT_VERIFIED',
        'message': 'Contact verification is required',
        'details': {
          'email_verified': false,
          'phone_verified': true,
        },
      });

      expect(error.code, 'CONTACT_NOT_VERIFIED');
      expect(error.message, 'Contact verification is required');
      expect(error.details, isNotNull);
      expect(error.details?['email_verified'], isFalse);
      expect(error.details?['phone_verified'], isTrue);
    });

    test('keeps compatibility when details is absent', () {
      final error = ApiError.fromMap({
        'code': 'INVALID_CREDENTIALS',
        'message': 'invalid email or password',
      });

      expect(error.code, 'INVALID_CREDENTIALS');
      expect(error.message, 'invalid email or password');
      expect(error.details, isNull);
    });
  });

  group('ApiEnvelope.fromMap', () {
    test('parses error envelopes without details', () {
      final envelope = ApiEnvelope<_DummyData>.fromMap(
        {
          'data': null,
          'error': {
            'code': 'INVALID_CREDENTIALS',
            'message': 'invalid email or password',
          },
        },
        _DummyData.fromMap,
      );

      expect(envelope.data, isNull);
      expect(envelope.error, isNotNull);
      expect(envelope.error?.code, 'INVALID_CREDENTIALS');
      expect(envelope.error?.details, isNull);
    });

    test('parses error envelopes with details', () {
      final envelope = ApiEnvelope<_DummyData>.fromMap(
        {
          'data': null,
          'error': {
            'code': 'CONTACT_NOT_VERIFIED',
            'message': 'Contact verification is required',
            'details': {
              'email_verified': false,
              'phone_verified': false,
            },
          },
        },
        _DummyData.fromMap,
      );

      expect(envelope.data, isNull);
      expect(envelope.error, isNotNull);
      expect(envelope.error?.code, 'CONTACT_NOT_VERIFIED');
      expect(envelope.error?.details?['email_verified'], isFalse);
      expect(envelope.error?.details?['phone_verified'], isFalse);
    });
  });
}

class _DummyData {
  final String value;

  _DummyData({required this.value});

  factory _DummyData.fromMap(Map<String, dynamic> map) {
    return _DummyData(value: map['value'] as String);
  }
}
