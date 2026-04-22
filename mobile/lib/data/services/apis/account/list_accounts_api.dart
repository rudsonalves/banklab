import '/core/result/result.dart';
import '/core/services/client_http/client_http.dart';
import '/core/services/logging/console_log.dart';
import 'dtos/account_summary_response_dto.dart';

class ListAccountsApi {
  final RestClient _client;

  ListAccountsApi(this._client);

  final _log = ConsoleLog('ListAccountsApi');

  AsyncResult<List<AccountSummaryResponseDto>> listAccounts() async {
    final response = await _client.get(
      RestClientRequest(
        path: '/accounts',
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

      final envelope = resp.data as Map<String, dynamic>;

      if (envelope['error'] != null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message:
                (envelope['error'] as Map<String, dynamic>)['message']
                    as String,
          ),
        );
      }

      final data = envelope['data'];
      if (data == null) {
        return Failure(
          AppError(
            code: AppErrorCode.httpError,
            message: 'No data received from the server.',
          ),
        );
      }

      if (data is! List) {
        return Failure(
          AppError(
            code: AppErrorCode.parsingError,
            message: 'Expected a list of accounts.',
          ),
        );
      }

      final accounts = data
          .map<AccountSummaryResponseDto>(
            (item) =>
                AccountSummaryResponseDto.fromMap(item as Map<String, dynamic>),
          )
          .toList();

      return Success(accounts);
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
