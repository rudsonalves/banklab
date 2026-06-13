import 'package:bankflow/core/routing/extra_codec.dart';
import 'package:bankflow/core/routing/models/transaction_password_setup_origin.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('ExtraCodec', () {
    test('encodes and decodes string extras', () {
      const codec = ExtraCodec();

      final encoded = codec.encode('123456');
      final decoded = codec.decode(encoded);

      expect(decoded, '123456');
    });

    test('does not serialize arbitrary maps that may contain credentials', () {
      const codec = ExtraCodec();

      expect(
        () => codec.encode({'pin': '123456'}),
        throwsUnsupportedError,
      );
    });

    test('encodes and decodes the typed setup origin', () {
      const codec = ExtraCodec();

      final encoded = codec.encode(TransactionPasswordSetupOrigin.transfer);
      final decoded = codec.decode(encoded);

      expect(decoded, TransactionPasswordSetupOrigin.transfer);
    });

    test('decodes an unknown setup origin to the safe post-login origin', () {
      const codec = ExtraCodec();

      final decoded = codec.decode(
        '{"type":"transaction_password_setup_origin","data":"unknown"}',
      );

      expect(decoded, TransactionPasswordSetupOrigin.postLogin);
    });
  });
}
