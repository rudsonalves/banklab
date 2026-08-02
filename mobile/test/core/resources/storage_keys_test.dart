import 'package:bankflow/core/resources/storage_keys.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('StorageKeys installation identity contract', () {
    test('exposes the durable installation id key', () {
      expect(StorageKeys.installationId, 'banklab.installation.id');
    });

    test('uses a separate non-secret local marker key', () {
      expect(
        StorageKeys.installationLocalMarker,
        'banklab.installation.marker',
      );
      expect(
        StorageKeys.installationLocalMarker,
        isNot(StorageKeys.installationId),
      );
    });
  });
}
