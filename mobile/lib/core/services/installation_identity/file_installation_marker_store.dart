import 'dart:io';

import 'package:path_provider/path_provider.dart';

import '/core/extensions/storage_app_error.dart';
import '/core/resources/storage_keys.dart';
import '/core/result/result.dart';
import '/core/services/logging/console_log.dart';
import 'installation_marker_store.dart';

typedef InstallationMarkerDirectoryResolver = Future<Directory> Function();

class FileInstallationMarkerStore implements InstallationMarkerStore {
  final InstallationMarkerDirectoryResolver _directoryResolver;

  FileInstallationMarkerStore({
    InstallationMarkerDirectoryResolver? directoryResolver,
  }) : _directoryResolver = directoryResolver ?? getApplicationSupportDirectory;

  final _log = ConsoleLog('FileInstallationMarkerStore');

  @override
  AsyncResult<bool> hasMarker() async {
    try {
      final file = await _markerFile();
      return Success(await file.exists());
    } catch (err, stack) {
      _log.error(
        'Failed to read installation local marker.',
        error: err,
        stack: stack,
      );
      return Failure(
        StorageAppError.storage(
          message: 'Failed to read installation local marker',
          details: err,
        ),
      );
    }
  }

  @override
  AsyncResult<Unit> markResolved() async {
    try {
      final file = await _markerFile();
      await file.parent.create(recursive: true);
      await file.writeAsString('resolved', flush: true);
      return Success(unit);
    } catch (err, stack) {
      _log.error(
        'Failed to write installation local marker.',
        error: err,
        stack: stack,
      );
      return Failure(
        StorageAppError.storage(
          message: 'Failed to write installation local marker',
          details: err,
        ),
      );
    }
  }

  Future<File> _markerFile() async {
    final directory = await _directoryResolver();
    final path = '${directory.path}/${StorageKeys.installationLocalMarker}';
    return File(path);
  }
}
