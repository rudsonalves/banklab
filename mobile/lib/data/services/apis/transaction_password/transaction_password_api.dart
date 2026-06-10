import '/core/result/result.dart';
import '/core/services/client_http/client/rest_client.dart';
import '/core/services/client_http/client/rest_client_request.dart';
import '/core/services/client_http/client/rest_client_response.dart';
import '/core/services/logging/console_log.dart';
import '../core/api_envelope.dart';
import 'dtos/create_transaction_password_request_dto.dart';
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

      final envelope =
          ApiEnvelope<TransactionPasswordStatusResponseDto>.fromMap(
            resp.data as Map<String, dynamic>,
            TransactionPasswordStatusResponseDto.fromApi,
          );

      if (envelope.error != null) {
        final error = envelope.error!;
        return Failure(
          AppError(
            code: _mapBackendErrorCode(error.code),
            message: error.message,
            details: {
              'code': error.code,
              if (error.details != null) 'details': error.details,
            },
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

      final transPasswdResp = envelope.data!;

      return Success(transPasswdResp);
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

  // Additional methods for updating and deleting transaction password
  // can be added here
  AppErrorCode _mapBackendErrorCode(String code) {
    if (code == _transactionPasswordAlreadySetBackendCode) {
      return AppErrorCode.transactionPasswordAlreadySet;
    }

    return AppErrorCode.httpError;
  }
}
