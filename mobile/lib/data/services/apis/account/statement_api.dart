import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import '../core/api_envelope.dart';
import 'dtos/statement_query_params_dto.dart';
import 'dtos/statement_response_dto.dart';

class StatementApi {
  final RestClient _client;

  StatementApi(this._client);

  final _log = ConsoleLog('StatementApi');

  /// Fetches the transaction statement for a specific account.
  ///
  /// Makes a GET request to `/accounts/{accountId}/statement` to retrieve
  /// the account's transaction history.
  ///
  /// Parameters:
  ///   - [accountId]: The unique identifier of the account.
  ///   - [queryParams]: Optional filtering parameters (e.g., date range, limit).
  ///     Defaults to [StatementQueryParamsDto()] with no filters.
  ///
  /// Returns:
  ///   An [AsyncResult] containing either:
  ///   - [Success] with [StatementResponseDto] containing the statement data.
  ///   - [Failure] with [AppError] if the request fails, parsing fails,
  ///     or the server returns an error.
  ///
  /// Error handling:
  ///   - HTTP errors (status codes outside 200-299): Returns [AppErrorCode.httpError].
  ///   - Missing response data: Returns [AppErrorCode.httpError].
  ///   - Response parsing failures: Returns [AppErrorCode.parsingError].
  AsyncResult<StatementResponseDto> getStatement(
    String accountId, {
    StatementQueryParamsDto queryParams = const StatementQueryParamsDto(),
  }) async {
    final response = await _client.get(
      RestClientRequest(
        path: '/accounts/$accountId/statement',
        queryParameters: queryParams.toMap(),
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

      final envelope = ApiEnvelope<StatementResponseDto>.fromMap(
        resp.data as Map<String, dynamic>,
        StatementResponseDto.fromMap,
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
