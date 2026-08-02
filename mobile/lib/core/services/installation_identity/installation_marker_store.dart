import '/core/result/result.dart';

abstract class InstallationMarkerStore {
  /// Checks if the installation marker exists by verifying the presence of
  /// a local marker file.
  /// Returns a [Success] with `true` if the marker exists, `false` if it does
  /// not, or a [Failure] if an error occurs during the check.
  AsyncResult<bool> hasMarker();

  /// Marks the installation as resolved by creating a local marker file.
  /// Returns a [Success] if the marker was successfully created, or a [Failure]
  /// if an error occurred during the process.
  AsyncResult<Unit> markResolved();
}
