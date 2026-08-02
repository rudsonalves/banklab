import '/core/resources/app_http_headers.dart';
import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '../core/api_envelope.dart';
import 'dtos/recipient_info_dto.dart';
import 'dtos/recipient_request_dto.dart';
import 'dtos/recipient_response_dto.dart';
import 'dtos/transfer_request_dto.dart';
import 'dtos/transfer_response_dto.dart';

class ApiTransfer {
  final RestClient _client;

  ApiTransfer(this._client);

  final _log = ConsoleLog('ApiTransfer');

  AsyncResult<TransferResponseDto> transfer({
    required String token,
    required TransferRequestDto dto,
  }) async {
    final response = await _client.post(
      RestClientRequest(
        path: '/accounts/internal-transfers',
        headers: {AppHttpHeaders.stepUpToken: token},
        body: dto.toMap(),
      ),
    );

    if (response.isFailure) {
      final err = response.error!;
      _log.error('HTTP request failed: ${err.message}', error: err);
      return Result.failure(err);
    }

    try {
      final resp = response.value as RestClientResponse;
      if (resp.statusCode == null ||
          resp.statusCode! < 200 ||
          resp.statusCode! >= 300) {
        _log.error('HTTP error: ${resp.statusCode} ${resp.statusMessage}');
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'HTTP error: ${resp.statusCode} ${resp.statusMessage}',
          ),
        );
      }

      final envelope = ApiEnvelope<TransferResponseDto>.fromMap(
        resp.data as Map<String, dynamic>,
        TransferResponseDto.fromMap,
      );

      if (envelope.error case final error?) {
        _log.error('API error: ${error.message}', error: error);
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: error.message,
            details: error.toAppErrorDetails(),
          ),
        );
      }

      if (envelope.data == null) {
        _log.error('No data received from the server.');
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

  AsyncResult<List<RecipientInfoDto>> getInternalRecipient(
    RecipientRequestDto dto,
  ) async {
    final response = await _client.get(
      RestClientRequest(
        path: '/accounts/internal-transfers/recipients',
        queryParameters: dto.toMap(),
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

      final envelope = ApiEnvelope<RecipientResponseDto>.fromMap(
        resp.data as Map<String, dynamic>,
        RecipientResponseDto.fromMap,
      );

      if (envelope.error case final error?) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: error.message,
            details: error.toAppErrorDetails(),
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

      return Success(envelope.data!.accounts);
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
