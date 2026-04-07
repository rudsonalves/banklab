import 'dart:convert';

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
      default:
        throw UnsupportedError('Unsupported type: $type');
    }
  }
}
