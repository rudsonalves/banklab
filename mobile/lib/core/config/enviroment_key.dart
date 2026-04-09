class EnviromentKey {
  static const baseUrl = String.fromEnvironment('BASE_URL');

  static const connectTimeout = int.fromEnvironment('CONNECT_TIMEOUT');
  static const receiveTimeout = int.fromEnvironment('RECEIVE_TIMEOUT');

  static const appMode = String.fromEnvironment('APP_MODE');

  static const appAccessToken = String.fromEnvironment('APP_ACCESS_TOKEN');
}
