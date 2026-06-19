import 'package:uuid/uuid.dart';

import '/core/resources/storage_keys.dart';
import '/core/result/result.dart';
import '/core/services/logging/console_log.dart';
import '/core/services/secure_storage/local_secure_storage.dart';
import 'installation_marker_store.dart';

class InstallationIdentityService {
  final LocalSecureStorage _secureStorage;
  final InstallationMarkerStore _markerStore;
  final Uuid _uuid;

  InstallationIdentityService({
    required LocalSecureStorage secureStorage,
    required InstallationMarkerStore markerStore,
    Uuid? uuid,
  }) : _secureStorage = secureStorage,
       _markerStore = markerStore,
       _uuid = uuid ?? const Uuid();

  final _log = ConsoleLog('InstallationIdentityService');

  static final RegExp _canonicalUuidV4 = RegExp(
    r'^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$',
  );

  /// Resolves the installation identity by checking for the presence of a local
  /// marker and a stored installation ID. If the marker is missing or the stored
  /// ID is invalid, a new installation ID is generated and stored, and the marker
  /// is updated accordingly. Returns a [Success] with the resolved installation ID,
  /// or a [Failure] if an error occurs during the process.
  AsyncResult<String> resolve() async {
    final markerResult = await _markerStore.hasMarker();
    if (markerResult.isFailure) {
      _log.warn(
        'Installation marker check failed.',
        label: 'resolve',
      );
      return Failure(markerResult.error!);
    }

    if (!markerResult.value!) {
      _log.info(
        'Installation marker is missing. A new identity will be generated.',
        label: 'resolve',
      );
      return _replaceIdentity();
    }

    final storedResult = await _secureStorage.read(StorageKeys.installationId);
    if (storedResult.isFailure) {
      final error = storedResult.error!;
      if (error.code == AppErrorCode.storageNotFound) {
        return _replaceIdentity();
      }

      _log.warn(
        'Installation identity read failed.',
        label: 'resolve',
      );
      return Failure(error);
    }

    final storedInstallationId = storedResult.value!;
    if (!_isCanonicalUuidV4(storedInstallationId)) {
      _log.warn(
        'Stored installation identity is invalid and will be replaced.',
        label: 'resolve',
      );
      return _replaceIdentity();
    }

    final markerWriteResult = await _markerStore.markResolved();
    if (markerWriteResult.isFailure) {
      _log.warn(
        'Installation marker refresh failed.',
        label: 'resolve',
      );
      return Failure(markerWriteResult.error!);
    }

    return Success(storedInstallationId);
  }

  /// Replaces the existing installation identity with a new one by generating a
  /// new UUID, storing it securely, and updating the installation marker. Returns a
  /// [Success] with the new installation ID if successful, or a [Failure] if an error
  /// occurs during the process.
  AsyncResult<String> _replaceIdentity() async {
    final installationId = _uuid.v4().toLowerCase();

    final writeResult = await _secureStorage.write(
      StorageKeys.installationId,
      installationId,
    );
    if (writeResult.isFailure) {
      _log.warn(
        'Installation identity write failed.',
        label: 'replace',
      );
      return Failure(writeResult.error!);
    }

    final markerWriteResult = await _markerStore.markResolved();
    if (markerWriteResult.isFailure) {
      _log.warn(
        'Installation marker write failed.',
        label: 'replace',
      );
      return Failure(markerWriteResult.error!);
    }

    return Success(installationId);
  }

  /// Checks if the provided value is a canonical UUID version 4 string.
  bool _isCanonicalUuidV4(String value) {
    return _canonicalUuidV4.hasMatch(value);
  }
}
