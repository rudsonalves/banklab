import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/dio/dio_error_mapper.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('mapHttpError', () {
    test(
      'maps ACCOUNT_APPROVAL_REQUIRED to accountApprovalRequired with approved message',
      () {
        final error = mapHttpError(
          _dioBadResponse(
            statusCode: 403,
            data: {
              'data': null,
              'error': {
                'code': 'ACCOUNT_APPROVAL_REQUIRED',
                'message': 'Account approval required',
              },
            },
          ),
        );

        expect(error.code, AppErrorCode.accountApprovalRequired);
        expect(
          error.message,
          'Sua conta ainda está aguardando aprovação. Assim que ela for liberada, você poderá acessar o app.',
        );
        expect(error.statusCode, 403);
      },
    );

    test('keeps INVALID_CREDENTIALS behavior unchanged', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'data': null,
            'error': {
              'code': 'INVALID_CREDENTIALS',
              'message': 'Invalid credentials',
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.httpError);
      expect(error.message, 'Invalid credentials');
      expect(error.statusCode, 403);
    });

    test('keeps generic forbidden as httpError when code is not approval', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'error': {
              'message': 'Forbidden',
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.httpError);
      expect(error.message, 'Forbidden');
      expect(error.statusCode, 403);
    });
  });
}

DioException _dioBadResponse({
  required int statusCode,
  required Map<String, dynamic> data,
}) {
  final requestOptions = RequestOptions(path: '/auth/login');

  return DioException(
    requestOptions: requestOptions,
    type: DioExceptionType.badResponse,
    response: Response(
      requestOptions: requestOptions,
      statusCode: statusCode,
      data: data,
    ),
  );
}
