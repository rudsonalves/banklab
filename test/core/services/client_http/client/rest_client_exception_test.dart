import 'package:bankflow/core/result/result.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('AppError', () {
    test('should expose its fields', () {
      const error = AppError(
        message: 'Unauthorized',
        statusCode: 401,
        code: AppErrorCode.httpError,
        details: {'error': 'invalid_token'},
      );

      expect(error.statusCode, 401);
      expect(error.code, AppErrorCode.httpError);
      expect(error.message, 'Unauthorized');
      expect(error.details, {'error': 'invalid_token'});
    });

    test('toString should include statusCode, code and message', () {
      const error = AppError(
        message: 'Unauthorized',
        statusCode: 401,
        code: AppErrorCode.httpError,
      );

      expect(
        error.toString(),
        'AppError(401, httpError, Unauthorized)',
      );
    });
  });
}
