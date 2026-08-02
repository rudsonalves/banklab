abstract final class AppHttpHeaders {
  const AppHttpHeaders._();

  static const accept = 'Accept';
  static const authorization = 'Authorization';
  static const contentType = 'Content-Type';
  static const appToken = 'X-App-Token';
  static const installationId = 'X-Installation-Id';
  static const stepUpToken = 'X-Step-Up-Token';
  static const traceId = 'X-Trace-Id';

  static const sensitiveLowercase = {
    'authorization',
    'x-app-token',
    'x-installation-id',
    'x-step-up-token',
  };

  static String bearer(String token) => 'Bearer $token';
}
