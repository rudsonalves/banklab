import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client/rest_client.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_response.dart';
import 'package:bankflow/data/repositories/contact_verification/contact_verification_repository_impl.dart';
import 'package:bankflow/data/services/apis/contact_verification/contact_verification_api.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_confirm_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/dtos/contact_verification_request_response_dto.dart';
import 'package:bankflow/data/services/apis/contact_verification/enums/contact_verification_channel.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('ContactVerificationRepositoryImpl.requestContactVerification', () {
    test('delegates to API and returns success response', () async {
      final api = _FakeContactVerificationApi(
        requestResult: Success(
          ContactVerificationRequestResponseDto(
            verificationId: 'verification-id-1',
            channel: ContactVerificationChannel.email,
            target: 'customer@example.com',
            // token: '123456',
            expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
          ),
        ),
        confirmResult: Success(
          ContactVerificationConfirmResponseDto(
            verificationToken: 'verified-token',
            channel: ContactVerificationChannel.email,
            target: 'customer@example.com',
            verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
          ),
        ),
      );

      final repository = ContactVerificationRepositoryImpl(api: api);

      final result = await repository.requestContactVerification(
        ContactVerificationRequestDto(
          channel: ContactVerificationChannel.email,
          target: 'customer@example.com',
        ),
      );

      expect(result, isA<Success<ContactVerificationRequestResponseDto>>());
      expect(result.value?.verificationId, 'verification-id-1');
      expect(result.value?.channel, ContactVerificationChannel.email);
      expect(result.value?.target, 'customer@example.com');
      expect(api.requestCalls, 1);
    });

    test('propagates AppError failure from API', () async {
      final api = _FakeContactVerificationApi(
        requestResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'target is invalid',
          ),
        ),
        confirmResult: Success(
          ContactVerificationConfirmResponseDto(
            verificationToken: 'verified-token',
            channel: ContactVerificationChannel.email,
            target: 'customer@example.com',
            verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
          ),
        ),
      );

      final repository = ContactVerificationRepositoryImpl(api: api);

      final result = await repository.requestContactVerification(
        ContactVerificationRequestDto(
          channel: ContactVerificationChannel.email,
          target: 'invalid',
        ),
      );

      expect(result, isA<Failure<ContactVerificationRequestResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'target is invalid');
      expect(api.requestCalls, 1);
    });
  });

  group('ContactVerificationRepositoryImpl.confirmContactVerification', () {
    test('delegates to API and returns success response', () async {
      final api = _FakeContactVerificationApi(
        requestResult: Success(
          ContactVerificationRequestResponseDto(
            verificationId: 'verification-id-1',
            channel: ContactVerificationChannel.phone,
            target: '+5511999999999',
            // token: '123456',
            expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
          ),
        ),
        confirmResult: Success(
          ContactVerificationConfirmResponseDto(
            verificationToken: 'verified-token',
            channel: ContactVerificationChannel.phone,
            target: '+5511999999999',
            verifiedAt: DateTime.parse('2026-05-18T12:03:00Z'),
          ),
        ),
      );

      final repository = ContactVerificationRepositoryImpl(api: api);

      final result = await repository.confirmContactVerification(
        ContactVerificationConfirmRequestDto(
          verificationId: 'verification-id-1',
          token: '123456',
        ),
      );

      expect(result, isA<Success<ContactVerificationConfirmResponseDto>>());
      expect(result.value?.verificationToken, 'verified-token');
      expect(result.value?.channel, ContactVerificationChannel.phone);
      expect(result.value?.target, '+5511999999999');
      expect(api.confirmCalls, 1);
    });

    test('propagates AppError failure from API', () async {
      final api = _FakeContactVerificationApi(
        requestResult: Success(
          ContactVerificationRequestResponseDto(
            verificationId: 'verification-id-1',
            channel: ContactVerificationChannel.phone,
            target: '+5511999999999',
            // token: '123456',
            expiresAt: DateTime.parse('2026-05-18T12:10:00Z'),
          ),
        ),
        confirmResult: const Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'invalid verification token',
          ),
        ),
      );

      final repository = ContactVerificationRepositoryImpl(api: api);

      final result = await repository.confirmContactVerification(
        ContactVerificationConfirmRequestDto(
          verificationId: 'verification-id-1',
          token: '000000',
        ),
      );

      expect(result, isA<Failure<ContactVerificationConfirmResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'invalid verification token');
      expect(api.confirmCalls, 1);
    });
  });
}

class _FakeContactVerificationApi extends ContactVerificationApi {
  Result<ContactVerificationRequestResponseDto> requestResult;
  Result<ContactVerificationConfirmResponseDto> confirmResult;

  int requestCalls = 0;
  int confirmCalls = 0;

  _FakeContactVerificationApi({
    required this.requestResult,
    required this.confirmResult,
  }) : super(_NoopRestClient());

  @override
  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    requestCalls++;
    return requestResult;
  }

  @override
  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    confirmCalls++;
    return confirmResult;
  }
}

class _NoopRestClient implements RestClient {
  const _NoopRestClient();

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      const Failure(
        AppError(code: AppErrorCode.unexpected, message: 'Not implemented'),
      );
}
