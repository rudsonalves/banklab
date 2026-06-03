import 'package:bankflow/core/routing/extra_codec.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('ExtraCodec', () {
    test('encodes and decodes string extras', () {
      const codec = ExtraCodec();

      final encoded = codec.encode('123456');
      final decoded = codec.decode(encoded);

      expect(decoded, '123456');
    });
  });
}
