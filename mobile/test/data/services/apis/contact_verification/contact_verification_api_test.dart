import 'package:bankflow/core/resources/app_env.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/contact_verification/contact_verification_api.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/enums/contact_verification_channel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('ContactVerificationApi.requestContactVerification', () {
    test(
      'calls POST /auth/contact-verifications with X-App-Token and request body',
      () async {
        final client = _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _requestSuccessEnvelope(),
            ),
          ),
        );
        final api = ContactVerificationApi(client);

        final result = await api.requestContactVerification(
          ContactVerificationRequestDto(
            channel: ContactVerificationChannel.email,
            target: 'user@example.com',
          ),
        );

        expect(result, isA<Success<ContactVerificationRequestResponseDto>>());
        expect(client.postCalls, 1);
        expect(client.lastPostRequest?.path, '/auth/contact-verifications');
        expect(
          client.lastPostRequest?.headers?['X-App-Token'],
          AppEnv.appToken,
        );
        expect(client.lastPostRequest?.body, {
          'channel': 'email',
          'target': 'user@example.com',
        });
      },
    );

    test('parses request verification success envelope', () async {
      final api = ContactVerificationApi(
        _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _requestSuccessEnvelope(),
            ),
          ),
        ),
      );

      final result = await api.requestContactVerification(
        ContactVerificationRequestDto(
          channel: ContactVerificationChannel.phone,
          target: '+5511999999999',
        ),
      );

      expect(result, isA<Success<ContactVerificationRequestResponseDto>>());
      final dto = result.value!;
      expect(dto.verificationId, 'a5d4f5f1-a1b0-4f58-9f74-123456789abc');
      expect(dto.channel, ContactVerificationChannel.email);
      expect(dto.target, 'user@example.com');
      // expect(dto.token, '123456');
      expect(dto.expiresAt, DateTime.parse('2026-05-18T12:10:00Z'));
    });

    test('maps backend error envelope to AppError failure', () async {
      final api = ContactVerificationApi(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'INVALID_CONTACT_TARGET',
                  'message': 'target is invalid',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.requestContactVerification(
        ContactVerificationRequestDto(
          channel: ContactVerificationChannel.email,
          target: 'invalid-email',
        ),
      );

      expect(result, isA<Failure<ContactVerificationRequestResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'target is invalid');
    });
  });

  group('ContactVerificationApi.confirmContactVerification', () {
    test(
      'calls POST /auth/contact-verifications/confirm with X-App-Token and request body',
      () async {
        final client = _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _confirmSuccessEnvelope(),
            ),
          ),
        );
        final api = ContactVerificationApi(client);

        final result = await api.confirmContactVerification(
          ContactVerificationConfirmRequestDto(
            verificationId: 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
            token: '123456',
          ),
        );

        expect(result, isA<Success<ContactVerificationConfirmResponseDto>>());
        expect(client.postCalls, 1);
        expect(
          client.lastPostRequest?.path,
          '/auth/contact-verifications/confirm',
        );
        expect(
          client.lastPostRequest?.headers?['X-App-Token'],
          AppEnv.appToken,
        );
        expect(client.lastPostRequest?.body, {
          'verification_id': 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
          'token': '123456',
        });
      },
    );

    test('parses confirm verification success envelope', () async {
      final api = ContactVerificationApi(
        _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(
              statusCode: 200,
              data: _confirmSuccessEnvelope(),
            ),
          ),
        ),
      );

      final result = await api.confirmContactVerification(
        ContactVerificationConfirmRequestDto(
          verificationId: 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
          token: '123456',
        ),
      );

      expect(result, isA<Success<ContactVerificationConfirmResponseDto>>());
      final dto = result.value!;
      expect(dto.verificationToken, 'token-confirmado');
      expect(dto.channel, ContactVerificationChannel.email);
      expect(dto.target, 'user@example.com');
      expect(dto.verifiedAt, DateTime.parse('2026-05-18T12:03:00Z'));
    });

    test('maps backend error envelope to AppError failure', () async {
      final api = ContactVerificationApi(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'CONTACT_VERIFICATION_INVALID',
                  'message': 'invalid verification token',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.confirmContactVerification(
        ContactVerificationConfirmRequestDto(
          verificationId: 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
          token: '000000',
        ),
      );

      expect(result, isA<Failure<ContactVerificationConfirmResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'invalid verification token');
    });
  });
}

Map<String, dynamic> _requestSuccessEnvelope() => {
  'data': {
    'verification_id': 'a5d4f5f1-a1b0-4f58-9f74-123456789abc',
    'channel': 'email',
    'target': 'user@example.com',
    'expires_at': '2026-05-18T12:10:00Z',
  },
  'error': null,
};

Map<String, dynamic> _confirmSuccessEnvelope() => {
  'data': {
    'verification_token': 'token-confirmado',
    'channel': 'email',
    'target': 'user@example.com',
    'verified_at': '2026-05-18T12:03:00Z',
  },
  'error': null,
};

class _FakeRestClient implements RestClient {
  _FakeRestClient({required this.postResult});

  final Result<RestClientResponse> postResult;
  RestClientRequest? lastPostRequest;
  int postCalls = 0;

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async {
    postCalls++;
    lastPostRequest = request;
    return postResult;
  }

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      throw UnimplementedError();
}
