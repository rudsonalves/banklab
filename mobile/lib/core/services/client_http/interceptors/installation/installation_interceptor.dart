import 'package:dio/dio.dart';

import '/core/resources/app_http_headers.dart';
import '/core/services/installation_identity/installation_identity.dart';
import '/core/services/logging/console_log.dart';

class InstallationInterceptor extends Interceptor {
  static const headerName = AppHttpHeaders.installationId;

  InstallationInterceptor({
    required InstallationIdentityService installationIdentityService,
  }) : _installationIdentityService = installationIdentityService;

  final InstallationIdentityService _installationIdentityService;
  final _log = ConsoleLog('InstallationInterceptor');

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final result = await _installationIdentityService.resolve();

    if (result.isFailure) {
      _log.warn('Blocking request because installation identity failed.');
      return handler.reject(
        DioException(
          requestOptions: options,
          type: DioExceptionType.unknown,
          message: 'Installation identity unavailable',
          error: result.error,
        ),
      );
    }

    options.headers[headerName] = result.value;
    return handler.next(options);
  }
}
