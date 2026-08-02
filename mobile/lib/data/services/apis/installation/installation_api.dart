import '/core/resources/app_http_headers.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/installation_identity/installation_identity.dart';
import '/core/services/logging/console_log.dart';
import '../core/api_envelope.dart';
import 'dtos/installation_registration_response_dto.dart';

class InstallationApi {
  InstallationApi({
    required RestClient client,
    required InstallationIdentityService installationIdentityService,
  }) : _client = client,
       _installationIdentityService = installationIdentityService;

  final RestClient _client;
  final InstallationIdentityService _installationIdentityService;
  final _log = ConsoleLog('InstallationApi');

  /// Registers the current installation with the backend using the provided tokens.
  /// Returns installation details on success.
  AsyncResult<InstallationRegistrationResponseDto> register({
    required String restrictedAccessToken,
    required String stepUpToken,
  }) async {
    final installationIdResult = await _installationIdentityService.resolve();
    if (installationIdResult.isFailure) {
      return Failure(installationIdResult.error!);
    }

    final result = await _client.post(
      RestClientRequest(
        path: '/security/installations',
        headers: {
          AppHttpHeaders.authorization: AppHttpHeaders.bearer(
            restrictedAccessToken,
          ),
          AppHttpHeaders.installationId: installationIdResult.value!,
          AppHttpHeaders.stepUpToken: stepUpToken,
        },
      ),
    );

    if (result.isFailure) {
      final error = result.error!;
      _log.warn(
        'Request failed (${error.statusCode ?? '-'}/${error.code.name}): '
        '${error.message}',
        label: 'register',
      );
      return Failure(error);
    }

    try {
      final response = result.value!;
      if (response.statusCode == null ||
          response.statusCode! < 200 ||
          response.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${response.statusCode} ${response.statusMessage}',
          label: 'register',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message:
                'HTTP error: ${response.statusCode} ${response.statusMessage}',
          ),
        );
      }

      final envelope = ApiEnvelope<InstallationRegistrationResponseDto>.fromMap(
        response.data as Map<String, dynamic>,
        InstallationRegistrationResponseDto.fromMap,
      );

      if (envelope.error case final error?) {
        _log.error('API error: ${error.message}', label: 'register');
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: error.message,
            details: error.toAppErrorDetails(),
          ),
        );
      }

      if (envelope.data == null) {
        _log.error('No data received from the server.', label: 'register');
        return const Failure(
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
        label: 'register',
      );
      return const Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }
}
