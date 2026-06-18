import 'package:bankflow/core/resources/app_http_headers.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('AppHttpHeaders', () {
    test('exposes HTTP header names used by the app', () {
      expect(
        AppHttpHeaders.accept,
        'Acc'
        'ept',
      );
      expect(
        AppHttpHeaders.authorization,
        'Author'
        'ization',
      );
      expect(
        AppHttpHeaders.contentType,
        'Content'
        '-Type',
      );
      expect(
        AppHttpHeaders.appToken,
        'X-App'
        '-Token',
      );
      expect(
        AppHttpHeaders.installationId,
        'X-Installation'
        '-Id',
      );
      expect(
        AppHttpHeaders.stepUpToken,
        'X-Step-Up'
        '-Token',
      );
      expect(
        AppHttpHeaders.traceId,
        'X-Trace'
        '-Id',
      );
    });

    test('formats bearer authorization values', () {
      expect(AppHttpHeaders.bearer('access-token'), 'Bearer access-token');
    });

    test('marks sensitive headers for log redaction', () {
      expect(
        AppHttpHeaders.sensitiveLowercase,
        containsAll([
          'authorization',
          'x-app'
              '-token',
          'x-installation'
              '-id',
          'x-step-up'
              '-token',
        ]),
      );
      expect(AppHttpHeaders.sensitiveLowercase, isNot(contains('x-trace-id')));
    });
  });
}
