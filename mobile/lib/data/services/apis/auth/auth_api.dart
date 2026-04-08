import 'dart:developer';

import 'package:bankflow/domain/auth/models/user_profile.dart';

import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/data/services/apis/auth/dtos/login_request_dto.dart';
import '/domain/auth/models/auth_user.dart';
import '../core/api_envelope.dart';
import 'dtos/register_request_dto.dart';
import 'dtos/register_response_dto.dart';

class AuthApi {
  final RestClient _client;

  AuthApi(this._client);

  AsyncResult<Unit> register(RegisterRequestDto dto) async {
    _client.setBaseUrl('http://localhost:3000/api/v1');

    final response = await _client.post(
      RestClientRequest(
        path: '/auth/register',
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) return Result.failure(response.error!);

    try {
      final envelope = ApiEnvelope<RegisterResponseDto>.fromMap(
        response.value as Map<String, dynamic>,
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
    } catch (err) {
      log('Error parsing response: $err');
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<LoggedUser> login(LoginRequestDto dto) async {
    _client.setBaseUrl('http://localhost:3000/api/v1');

    final response = await _client.post(
      RestClientRequest(
        path: '/auth/login',
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) return Result.failure(response.error!);

    try {
      final envelope = ApiEnvelope<LoggedUser>.fromMap(
        response.value as Map<String, dynamic>,
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

      if (envelope.data == null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      return Success(envelope.data!);
    } catch (err) {
      log('Error parsing response: $err');
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<UserProfile> getProfile() async {
    _client.setBaseUrl('http://localhost:3000/api/v1');

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
    } catch (err) {
      log('Error parsing response: $err');
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }
}
