enum AppMode { dev, staging, prod }

class AppEnv {
  static const _baseUrl = String.fromEnvironment('BASE_URL');
  static const _appToken = String.fromEnvironment('APP_ACCESS_TOKEN');
  static const _connectTimeout = int.fromEnvironment(
    'CONNECT_TIMEOUT',
    defaultValue: 10000,
  );
  static const _receiveTimeout = int.fromEnvironment(
    'RECEIVE_TIMEOUT',
    defaultValue: 10000,
  );
  static const _appMode = String.fromEnvironment(
    'APP_MODE',
    defaultValue: 'dev',
  );

  static AppMode get _mode {
    final modeStr = _appMode.toLowerCase();

    switch (modeStr) {
      case 'prod':
        return AppMode.prod;
      case 'staging':
        return AppMode.staging;
      case 'dev':
        return AppMode.dev;
      default:
        throw StateError('Invalid APP_MODE: $modeStr');
    }
  }

  static String get baseUrl {
    if (_baseUrl.isEmpty) {
      throw StateError('BASE_URL not defined.');
    }

    final uri = Uri.tryParse(_baseUrl);
    if (uri == null || !uri.hasScheme) {
      throw StateError('BASE_URL is invalid: $_baseUrl');
    }

    return _baseUrl;
  }

  static String get appToken => _appToken;

  static int get connectTimeout => _connectTimeout;

  static int get receiveTimeout => _receiveTimeout;

  static bool get isProd => _mode == AppMode.prod;
}
