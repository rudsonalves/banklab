import '/core/resources/app_env.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '../../apis/core/api_envelope.dart';
import 'dtos/contact_verification_confirm_request_dto.dart';
import 'dtos/contact_verification_confirm_response_dto.dart';
import 'dtos/contact_verification_request_dto.dart';
import 'dtos/contact_verification_request_response_dto.dart';

class ContactVerificationApi {
  final RestClient _client;

  ContactVerificationApi(this._client);

  final _log = ConsoleLog('ContactVerificationApi');

  AsyncResult<ContactVerificationRequestResponseDto> requestContactVerification(
    ContactVerificationRequestDto dto,
  ) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/auth/contact-verifications',
        headers: {
          'X-App-Token': AppEnv.appToken,
        },
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) {
      _log.error(
        'Request failed: ${response.error}',
        label: 'requestContactVerification',
      );
      return Result.failure(response.error!);
    }

    try {
      final resp = response.value as RestClientResponse;
      if (resp.statusCode == null ||
          resp.statusCode! < 200 ||
          resp.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          label: 'requestContactVerification',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope =
          ApiEnvelope<ContactVerificationRequestResponseDto>.fromMap(
            resp.data as Map<String, dynamic>,
            ContactVerificationRequestResponseDto.fromMap,
          );

      if (envelope.error != null) {
        _log.error(
          'API error: ${envelope.error!.message}',
          label: 'requestContactVerification',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      if (envelope.data == null) {
        _log.error(
          'No data received from the server.',
          label: 'requestContactVerification',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      const appMode = String.fromEnvironment('APP_MODE', defaultValue: 'dev');
      if (appMode.toLowerCase() == 'dev') {
        final token =
            (resp.data as Map<String, dynamic>)['data']['debug_token']
                as String?;
        _log.info(
          'Contact verification token (${envelope.data!.channel}): $token',
          label: 'requestContactVerification',
        );
      }

      return Success(envelope.data!);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'requestContactVerification',
      );
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<ContactVerificationConfirmResponseDto> confirmContactVerification(
    ContactVerificationConfirmRequestDto dto,
  ) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/auth/contact-verifications/confirm',
        headers: {
          'X-App-Token': AppEnv.appToken,
        },
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) {
      _log.error(
        'Request failed: ${response.error}',
        label: 'confirmContactVerification',
      );
      return Result.failure(response.error!);
    }

    try {
      final resp = response.value as RestClientResponse;
      if (resp.statusCode == null ||
          resp.statusCode! < 200 ||
          resp.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          label: 'confirmContactVerification',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope =
          ApiEnvelope<ContactVerificationConfirmResponseDto>.fromMap(
            resp.data as Map<String, dynamic>,
            ContactVerificationConfirmResponseDto.fromMap,
          );

      if (envelope.error != null) {
        _log.error(
          'API error: ${envelope.error!.message}',
          label: 'confirmContactVerification',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      if (envelope.data == null) {
        _log.error(
          'No data received from the server.',
          label: 'confirmContactVerification',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      return Success(envelope.data!);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'confirmContactVerification',
      );
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }
}
