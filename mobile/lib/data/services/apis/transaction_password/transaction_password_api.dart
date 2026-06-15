import '/core/result/result.dart';
import '/core/services/client_http/client/rest_client.dart';
import '/core/services/client_http/client/rest_client_request.dart';
import '/core/services/client_http/client/rest_client_response.dart';
import '/core/services/logging/console_log.dart';
import '../core/api_envelope.dart';
import 'dtos/create_transaction_password_request_dto.dart';
import 'dtos/step_up_authorize_request_dto.dart';
import 'dtos/step_up_authorize_response_dto.dart';
import 'dtos/transaction_password_status_response_dto.dart';

const _transactionPasswordAlreadySetBackendCode =
    'TRANSACTION_PASSWORD_ALREADY_SET';

class TransactionPasswordApi {
  final RestClient _client;

  TransactionPasswordApi(this._client);

  final _log = ConsoleLog('TransactionPasswordApi');

  AsyncResult<TransactionPasswordStatusResponseDto> create(
    CreateTransactionPasswordRequestDto dto,
  ) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/security/transaction-password',
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) {
      final error = response.error!;
      _log.error(
        'HTTP error: ${error.message}',
        label: 'create',
      );
      return Result.failure(error);
    }

    try {
      final resp = response.value as RestClientResponse;
      if (resp.statusCode == null ||
          resp.statusCode! < 200 ||
          resp.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          label: 'create',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope =
          ApiEnvelope<TransactionPasswordStatusResponseDto>.fromMap(
            resp.data as Map<String, dynamic>,
            TransactionPasswordStatusResponseDto.fromApi,
          );

      if (envelope.error != null) {
        final error = envelope.error!;
        _log.error(
          'Backend error: ${error.code} ${error.message}',
          label: 'create',
        );
        return Failure(
          AppError(
            code: _mapBackendErrorCode(error.code),
            message: error.message,
            details: error.toAppErrorDetails(),
          ),
        );
      }

      if (envelope.data == null) {
        _log.error(
          'No data received from the server.',
          label: 'create',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      final transPasswdResp = envelope.data!;

      return Success(transPasswdResp);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'create',
      );
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  AsyncResult<StepUpAuthorizeResponseDto> stepUpAuthorize(
    StepUpAuthorizeRequestDto dto,
  ) async {
    final result = await _client.post(
      RestClientRequest(
        path: '/security/step-up/authorize',
        body: dto.toMap(),
      ),
    );

    if (result.isFailure) {
      final error = result.error!;
      _log.error(
        'HTTP error: ${error.message}',
        label: 'stepUpAuthorize',
      );
      return Failure(result.error!);
    }

    try {
      final response = result.value as RestClientResponse;
      if (response.statusCode == null ||
          response.statusCode! < 200 ||
          response.statusCode! >= 300) {
        _log.error(
          'HTTP error: ${response.statusCode} ${response.statusMessage}',
          label: 'stepUpAuthorize',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message:
                'HTTP error: ${response.statusCode} ${response.statusMessage}',
          ),
        );
      }

      final apiEnv = ApiEnvelope<StepUpAuthorizeResponseDto>.fromMap(
        response.data as Map<String, dynamic>,
        StepUpAuthorizeResponseDto.fromApi,
      );

      if (apiEnv.error case final error?) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: error.message,
            details: error.toAppErrorDetails(),
          ),
        );
      }

      final setUpAuthResp = apiEnv.data;
      if (setUpAuthResp == null) {
        _log.error(
          'No data received from the server.',
          label: 'stepUpAuthorize',
        );
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      return Success(setUpAuthResp);
    } catch (err, stack) {
      _log.error(
        'Error parsing response: $err',
        error: err,
        stack: stack,
        label: 'stepUpAuthorize',
      );
      return Failure(
        AppError(
          code: AppErrorCode.parsingError,
          message: 'Failed to parse the response from the server.',
        ),
      );
    }
  }

  // Additional methods for updating and deleting transaction password
  // can be added here
  AppErrorCode _mapBackendErrorCode(String code) {
    if (code == _transactionPasswordAlreadySetBackendCode) {
      return AppErrorCode.transactionPasswordAlreadySet;
    }

    return AppErrorCode.httpError;
  }
}
