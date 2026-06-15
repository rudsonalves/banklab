import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/dio/dio_error_mapper.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('mapHttpError', () {
    test(
      'maps TRANSACTION_PASSWORD_LOCKED to transactionPasswordLocked',
      () {
        final error = mapHttpError(
          DioException(
            requestOptions: RequestOptions(path: '/security/step-up/authorize'),
            type: DioExceptionType.badResponse,
            response: Response(
              requestOptions: RequestOptions(
                path: '/security/step-up/authorize',
              ),
              statusCode: 403,
              data: {
                'error': {
                  'code': 'TRANSACTION_PASSWORD_LOCKED',
                  'message': 'Transaction password is locked',
                },
              },
            ),
          ),
        );

        expect(error.code, AppErrorCode.transactionPasswordLocked);
        expect(error.statusCode, 403);
        expect(error.message, 'Transaction password is locked');
        expect(
          backendErrorCode(error),
          'TRANSACTION_PASSWORD_LOCKED',
        );
      },
    );

    test('preserves backend code when the error contains details', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 401,
          data: {
            'error': {
              'code': 'STEP_UP_TOKEN_EXPIRED',
              'message': 'Step-up token expired',
              'details': {'expired_at': '2026-06-12T12:00:00Z'},
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.httpError);
      expect(backendErrorCode(error), 'STEP_UP_TOKEN_EXPIRED');
      expect(
        error.details,
        {
          'code': 'STEP_UP_TOKEN_EXPIRED',
          'expired_at': '2026-06-12T12:00:00Z',
        },
      );
    });

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

    test('maps CONTACT_NOT_VERIFIED with both channels pending', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'data': null,
            'error': {
              'code': 'CONTACT_NOT_VERIFIED',
              'message': 'Contact not verified',
              'details': {
                'email_verified': false,
                'phone_verified': false,
              },
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.contactNotVerified);
      expect(error.message, 'Confirme seu e-mail e telefone antes de entrar.');
      expect(error.statusCode, 403);
      expect(error.details, isA<Map<String, dynamic>>());
      final details = error.details! as Map<String, dynamic>;
      expect(details['code'], 'CONTACT_NOT_VERIFIED');
      expect(details['email_verified'], isFalse);
      expect(details['phone_verified'], isFalse);
    });

    test('maps CONTACT_NOT_VERIFIED with only e-mail pending', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'data': null,
            'error': {
              'code': 'CONTACT_NOT_VERIFIED',
              'message': 'Contact not verified',
              'details': {
                'email_verified': false,
                'phone_verified': true,
              },
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.contactNotVerified);
      expect(error.message, 'Confirme seu e-mail antes de entrar.');
      expect(error.statusCode, 403);
    });

    test('maps CONTACT_NOT_VERIFIED with only phone pending', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'data': null,
            'error': {
              'code': 'CONTACT_NOT_VERIFIED',
              'message': 'Contact not verified',
              'details': {
                'email_verified': true,
                'phone_verified': false,
              },
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.contactNotVerified);
      expect(error.message, 'Confirme seu telefone antes de entrar.');
      expect(error.statusCode, 403);
    });

    test('maps CONTACT_NOT_VERIFIED without details to generic message', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'data': null,
            'error': {
              'code': 'CONTACT_NOT_VERIFIED',
              'message': 'Contact not verified',
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.contactNotVerified);
      expect(error.message, 'Confirme seu e-mail e telefone antes de entrar.');
      expect(error.statusCode, 403);
      expect(backendErrorCode(error), 'CONTACT_NOT_VERIFIED');
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

    test('does not map generic 403 to contactNotVerified', () {
      final error = mapHttpError(
        _dioBadResponse(
          statusCode: 403,
          data: {
            'data': null,
            'error': {
              'code': 'FORBIDDEN',
              'message': 'Forbidden',
            },
          },
        ),
      );

      expect(error.code, AppErrorCode.httpError);
      expect(error.code, isNot(AppErrorCode.contactNotVerified));
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
