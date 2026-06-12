import 'dart:convert';

import '/data/services/apis/transfer/dtos/recipient_info_dto.dart';
import '/data/services/cache/last_login/models/last_login_identity.dart';
import '/ui/pages/home/transfer/models/transfer_confirmation_data.dart';
import '/core/routing/models/transaction_password_setup_origin.dart';

class ExtraCodec extends Codec<Object?, String> {
  const ExtraCodec();

  @override
  Converter<Object?, String> get encoder => const _ExtraEncoder();

  @override
  Converter<String, Object?> get decoder => const _ExtraDecoder();
}

class _ExtraEncoder extends Converter<Object?, String> {
  const _ExtraEncoder();

  @override
  String convert(Object? extra) {
    if (extra is Map<String, Object?>) {
      return jsonEncode({
        'type': 'map',
        'data': extra,
      });
    }

    if (extra is List<Object?>) {
      return jsonEncode({
        'type': 'list',
        'data': extra,
      });
    }

    if (extra is String || extra is num || extra is bool) {
      return jsonEncode({
        'type': 'primitive',
        'data': extra,
      });
    }

    if (extra == null) {
      return jsonEncode({
        'type': 'primitive',
        'data': null,
      });
    }

    if (extra is RecipientInfoDto) {
      return jsonEncode({
        'type': 'recipient_info',
        'data': extra.toMap(),
      });
    }

    if (extra is TransferConfirmationData) {
      return jsonEncode({
        'type': 'transfer_confirmation_data',
        'data': extra.toMap(),
      });
    }

    if (extra is LastLoginIdentity) {
      return jsonEncode({
        'type': 'last_login_identity',
        'data': extra.toMap(),
      });
    }

    if (extra is TransactionPasswordSetupOrigin) {
      return jsonEncode({
        'type': 'transaction_password_setup_origin',
        'data': extra.name,
      });
    }

    throw UnsupportedError('Unsupported type: ${extra.runtimeType}');
  }
}

class _ExtraDecoder extends Converter<String, Object?> {
  const _ExtraDecoder();

  @override
  Object? convert(String extra) {
    final decoded = jsonDecode(extra) as Map<String, Object?>;
    final type = decoded['type'] as String;

    switch (type) {
      case 'map':
        return decoded['data'] as Map<String, Object?>;

      case 'list':
        return decoded['data'] as List<Object?>;

      case 'primitive':
        return decoded['data'];

      case 'recipient_info':
        return RecipientInfoDto.fromMap(
          decoded['data'] as Map<String, dynamic>,
        );

      case 'transfer_confirmation_data':
        return TransferConfirmationData.fromMap(
          decoded['data'] as Map<String, dynamic>,
        );

      case 'last_login_identity':
        return LastLoginIdentity.fromMap(
          decoded['data'] as Map<String, dynamic>,
        );

      case 'transaction_password_setup_origin':
        final originName = decoded['data'] as String;
        return TransactionPasswordSetupOrigin.fromName(originName);

      default:
        throw UnsupportedError('Unsupported type: $type');
    }
  }
}
