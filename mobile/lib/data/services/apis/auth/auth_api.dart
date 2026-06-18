import '/core/resources/app_env.dart';
import '/core/resources/app_http_headers.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '/data/services/apis/core/api_envelope.dart';
import '/domain/common/auth/models/auth_session/auth_session.dart';
import '/domain/common/auth/models/auth_user.dart';
import 'dtos/login_request_dto.dart';

class AuthApi {
  final RestClient _client;

  AuthApi(this._client);

  final _log = ConsoleLog('AuthApi');

  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/auth/login',
        headers: {
          AppHttpHeaders.appToken: AppEnv.appToken,
        },
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) {
      final error = response.error!;
      _log.warn(
        'Request failed (${error.statusCode ?? '-'}/${error.code.name}): '
        '${error.message}',
        label: 'login',
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
          label: 'login',
        );
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
        _log.error('API error: ${envelope.error!.message}', label: 'login');
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
        label: 'login',
      );
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<AuthSession> getAuthSession() async {
    final result = await _client.get(
      RestClientRequest(
        path: '/auth/session',
      ),
    );

    if (result.isFailure) {
      final error = result.error!;
      _log.warn(
        'Request failed (${error.statusCode ?? '-'}/${error.code.name}): '
        '${error.message}',
        label: 'getAuthSession',
      );
      return Result.failure(result.error!);
    }

    try {
      final response = result.value as RestClientResponse;
      if (response.statusCode == null ||
          response.statusCode! < 200 ||
          response.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${response.statusCode} ${response.statusMessage}',
          label: 'getAuthSession',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message:
                'HTTP error: ${response.statusCode} ${response.statusMessage}',
          ),
        );
      }

      final apiEnv = ApiEnvelope<AuthSession>.fromMap(
        response.data as Map<String, dynamic>,
        AuthSession.fromApi,
      );

      final authSession = apiEnv.data;
      if (authSession == null) {
        _log.error(
          'API error: ${apiEnv.error?.message}',
          label: 'getAuthSession',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: apiEnv.error?.message ?? 'Unknown API error',
          ),
        );
      }

      return Success(authSession);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'getAuthSession',
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
