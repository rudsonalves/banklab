import 'package:bankflow/core/resources/app_http_headers.dart';
import 'package:bankflow/core/services/client_http/client/rest_client_request.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('RestClientRequest', () {
    test('copyWith should keep original values when not provided', () {
      const request = RestClientRequest(
        path: '/users',
        headers: {AppHttpHeaders.authorization: 'Bearer token'},
        queryParameters: {'page': 1},
        body: {'name': 'Ada'},
      );

      final copied = request.copyWith();

      expect(copied.path, '/users');
      expect(copied.headers, {AppHttpHeaders.authorization: 'Bearer token'});
      expect(copied.queryParameters, {'page': 1});
      expect(copied.body, {'name': 'Ada'});
    });

    test('copyWith should replace provided fields', () {
      const request = RestClientRequest(
        path: '/users',
        headers: {AppHttpHeaders.authorization: 'Bearer token'},
        queryParameters: {'page': 1},
        body: {'name': 'Ada'},
      );

      final copied = request.copyWith(
        path: '/accounts',
        headers: {AppHttpHeaders.traceId: 'abc'},
        queryParameters: {'limit': 10},
        body: {'active': true},
      );

      expect(copied.path, '/accounts');
      expect(copied.headers, {AppHttpHeaders.traceId: 'abc'});
      expect(copied.queryParameters, {'limit': 10});
      expect(copied.body, {'active': true});
    });
  });
}
