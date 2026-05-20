import '/core/resources/app_env.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '../../apis/core/api_envelope.dart';
import 'dtos/cpf_check_response_dto.dart';
import 'dtos/register_request_dto.dart';
import 'dtos/register_response_dto.dart';

class RegistrationApi {
  final RestClient _client;

  RegistrationApi(this._client);

  final _log = ConsoleLog('RegistrationApi');

  AsyncResult<CpfCheckResponseDto> cpfCheck(String cpf) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/auth/cpf-check',
        headers: {
          'X-App-Token': AppEnv.appToken,
        },
        body: {
          'cpf': cpf,
        },
      ),
    );

    if (response.isFailure) {
      _log.error('Request failed: ${response.error}', label: 'cpfCheck');
      return Result.failure(response.error!);
    }

    try {
      final resp = response.value as RestClientResponse;
      if (resp.statusCode == null ||
          resp.statusCode! < 200 ||
          resp.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          label: 'cpfCheck',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope = ApiEnvelope<CpfCheckResponseDto>.fromMap(
        resp.data as Map<String, dynamic>,
        CpfCheckResponseDto.fromMap,
      );

      if (envelope.error != null) {
        _log.error('API error: ${envelope.error!.message}', label: 'cpfCheck');
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      return Success(envelope.data!);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'cpfCheck',
      );
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/auth/register',
        headers: {
          'X-App-Token': AppEnv.appToken,
        },
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) return Result.failure(response.error!);

    try {
      final resp = response.value as RestClientResponse;
      if (resp.statusCode == null ||
          resp.statusCode! < 200 ||
          resp.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          label: 'register',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope = ApiEnvelope<RegisterResponseDto>.fromMap(
        resp.data as Map<String, dynamic>,
        RegisterResponseDto.fromMap,
      );

      if (envelope.error != null) {
        _log.error('API error: ${envelope.error!.message}', label: 'register');
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      if (envelope.data == null) {
        _log.error('No data received from the server.', label: 'register');
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      return Success(unit);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'register',
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
