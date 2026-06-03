import 'package:bankflow/core/result/result.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('backendErrorCode', () {
    test('returns null for null error', () {
      expect(backendErrorCode(null), isNull);
    });

    test('extracts code from details map', () {
      const error = AppError(
        code: AppErrorCode.httpError,
        message: 'Invalid data',
        details: {'code': 'INVALID_DATA'},
      );

      expect(backendErrorCode(error), 'INVALID_DATA');
    });

    test('extracts code from nested error map', () {
      const error = AppError(
        code: AppErrorCode.httpError,
        message: 'Forbidden',
        details: {
          'error': {
            'code': 'FORBIDDEN',
            'message': 'Access denied',
          },
        },
      );

      expect(backendErrorCode(error), 'FORBIDDEN');
    });

    test('extracts code from nested details map', () {
      const error = AppError(
        code: AppErrorCode.httpError,
        message: 'Unauthorized',
        details: {
          'details': {'code': 'INVALID_TOKEN'},
        },
      );

      expect(backendErrorCode(error), 'INVALID_TOKEN');
    });

    test('returns null when details has no backend code', () {
      const error = AppError(
        code: AppErrorCode.httpError,
        message: 'Request error',
        details: {'message': 'Request error'},
      );

      expect(backendErrorCode(error), isNull);
    });

    test('returns null when details is not a map', () {
      const error = AppError(
        code: AppErrorCode.unexpected,
        message: 'Unexpected',
        details: 'raw error',
      );

      expect(backendErrorCode(error), isNull);
    });

    test('ignores blank code values', () {
      const error = AppError(
        code: AppErrorCode.httpError,
        message: 'Request error',
        details: {'code': '   '},
      );

      expect(backendErrorCode(error), isNull);
    });
  });
}
