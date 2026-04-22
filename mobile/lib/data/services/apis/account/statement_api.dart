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
