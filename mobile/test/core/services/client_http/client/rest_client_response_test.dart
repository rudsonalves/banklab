import 'package:bankflow/core/services/client_http/client/rest_client_response.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RestClientResponse', () {
    test('isSuccess should be true for 2xx status codes', () {
      expect(const RestClientResponse(statusCode: 200).isSuccess, isTrue);
      expect(const RestClientResponse(statusCode: 201).isSuccess, isTrue);
      expect(const RestClientResponse(statusCode: 299).isSuccess, isTrue);
    });

    test('isSuccess should be false for non 2xx status codes', () {
      expect(const RestClientResponse(statusCode: 199).isSuccess, isFalse);
      expect(const RestClientResponse(statusCode: 300).isSuccess, isFalse);
      expect(const RestClientResponse(statusCode: null).isSuccess, isFalse);
    });

    test('copyWith should keep original values when not provided', () {
      const response = RestClientResponse(
        data: {'ok': true},
        statusCode: 200,
        statusMessage: 'OK',
      );

      final copied = response.copyWith();

      expect(copied.data, {'ok': true});
      expect(copied.statusCode, 200);
      expect(copied.statusMessage, 'OK');
    });

    test('copyWith should replace provided fields', () {
      const response = RestClientResponse(
        data: {'ok': true},
        statusCode: 200,
        statusMessage: 'OK',
      );

      final copied = response.copyWith(
        data: {'ok': false},
        statusCode: 202,
        statusMessage: 'Accepted',
      );

      expect(copied.data, {'ok': false});
      expect(copied.statusCode, 202);
      expect(copied.statusMessage, 'Accepted');
    });
  });
}
