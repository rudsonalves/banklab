import '/core/resources/app_env.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '/domain/common/auth/models/auth_user.dart';
import '/domain/common/auth/models/user_profile.dart';
import '../../apis/core/api_envelope.dart';
import 'dtos/auth_me_response_dto.dart';
import 'dtos/customer_me_response_dto.dart';
import 'dtos/login_request_dto.dart';
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

    if (response.isFailure) {
      _log.error('Request failed: ${response.error}', label: 'login');
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

  // Remove this to a profile API service.
  AsyncResult<UserProfile> getProfile() async {
    final respCustomer = await _client.get(
      RestClientRequest(
        path: '/customers/me',
      ),
    );

    if (respCustomer.isFailure) {
      _log.error('Request failed: ${respCustomer.error}', label: 'getProfile');
      return Result.failure(respCustomer.error!);
    }

    final respMe = await _client.get(
      RestClientRequest(
        path: '/auth/me',
      ),
    );

    if (respMe.isFailure) {
      _log.error('Request failed: ${respMe.error}', label: 'getProfile');
      return Result.failure(respMe.error!);
    }

    try {
      final clientCustomer = respCustomer.value as RestClientResponse;
      if (clientCustomer.statusCode == null ||
          clientCustomer.statusCode! < 200 ||
          clientCustomer.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${clientCustomer.statusCode} ${clientCustomer.statusMessage}',
          label: 'getProfile',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message:
                'HTTP error: ${clientCustomer.statusCode} ${clientCustomer.statusMessage}',
          ),
        );
      }

      final envCustomer = ApiEnvelope<CustomerMeResponseDto>.fromMap(
        clientCustomer.data as Map<String, dynamic>,
        CustomerMeResponseDto.fromMap,
      );

      final clientMe = respMe.value as RestClientResponse;
      if (clientMe.statusCode == null ||
          clientMe.statusCode! < 200 ||
          clientMe.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${clientMe.statusCode} ${clientMe.statusMessage}',
          label: 'getProfile',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message:
                'HTTP error: ${clientMe.statusCode} ${clientMe.statusMessage}',
          ),
        );
      }

      final envUserMe = ApiEnvelope<AuthMeResponseDto>.fromMap(
        clientMe.data as Map<String, dynamic>,
        AuthMeResponseDto.fromMap,
      );

      final authMe = envUserMe.data!;
      final customer = envCustomer.data!;

      final userProfile = UserProfile.fromMe(
        userMe: authMe,
        customer: customer,
      );

      return Success(userProfile);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'getProfile',
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
