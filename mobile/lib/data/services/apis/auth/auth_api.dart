import '/core/resources/app_env.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/domain/auth/models/auth_user.dart';
import '/domain/auth/models/user_profile.dart';
import '../core/api_envelope.dart';
import 'dtos/register_request_dto.dart';
import 'dtos/register_response_dto.dart';

class AuthApi {
  final RestClient _client;

  AuthApi(this._client);

  final _log = ConsoleLog('AuthApi');

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
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      if (envelope.data == null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      return Success(unit);
    } catch (err, stack) {
      _log.error('Error parsing response: $err', error: err, stack: stack);
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/auth/login',
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
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope = ApiEnvelope<LoggedUser>.fromMap(
        resp.data as Map<String, dynamic>,
        LoggedUser.fromMap,
      );

      if (envelope.error != null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      return Success(envelope.data!);
    } catch (err, stack) {
      _log.error('Error parsing response: $err', error: err, stack: stack);
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  // Remove this to a profile API service.
  AsyncResult<UserProfile> getProfile() async {
    final response = await _client.get(
      RestClientRequest(
        path: 'profile/me',
      ),
    );

    if (response.isFailure) return Result.failure(response.error!);

    try {
      final envelope = ApiEnvelope<UserProfile>.fromMap(
        response.value as Map<String, dynamic>,
        UserProfile.fromMap,
      );

      if (envelope.error != null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: envelope.error!.message,
          ),
        );
      }

      if (envelope.data == null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      return Success(envelope.data!);
    } catch (err, stack) {
      _log.error('Error parsing response: $err', error: err, stack: stack);
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }
}
