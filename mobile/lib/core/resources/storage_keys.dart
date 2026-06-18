final class StorageKeys {
  static const accessToken = 'auth_token';
  static const refreshToken = 'refresh_token';

  static const lastLoginName = 'last_login_name';
  static const lastLoginIdentifier = 'last_login_identifier';

  /// Durable installation identity stored in secure storage.
  ///
  /// Logout, credential cleanup, and user switching must not delete this key.
  static const installationId = 'banklab.installation.id';

  /// Non-secret local marker used to validate that secure storage still belongs
  /// to the current app installation.
  ///
  /// This marker is not an identity and must be stored outside secure storage.
  /// If it is missing, the installation identity service must treat the app as a
  /// new installation and replace any existing [installationId] value.
  static const installationLocalMarker = 'banklab.installation.marker';
}
