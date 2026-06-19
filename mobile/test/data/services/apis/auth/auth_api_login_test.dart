import 'package:bankflow/core/resources/app_env.dart';
import 'package:bankflow/core/resources/app_http_headers.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/data/services/apis/auth/auth_api.dart';
import 'package:bankflow/data/services/apis/auth/dtos/login_request_dto.dart';
import 'package:bankflow/domain/common/auth/models/auth_state.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('AuthApi.login', () {
    test('parses operational login response', () async {
      final client = _FakeRestClient(
        postResult: Result.success(
          RestClientResponse(statusCode: 200, data: _operationalEnvelope()),
        ),
      );
      final api = AuthApi(client);

      final result = await api.login(_loginRequest());

      expect(result, isA<Success<AuthState>>());
      expect(result.value, isA<OperationalAuthState>());
      final operational = result.value! as OperationalAuthState;
      expect(operational.accessToken, 'access-token');
      expect(operational.refreshToken, 'refresh-token');
      expect(client.lastPostRequest?.path, '/auth/login');
      expect(
        client.lastPostRequest?.headers?[AppHttpHeaders.appToken],
        AppEnv.appToken,
      );
    });

    test('parses restricted login response', () async {
      final api = AuthApi(
        _FakeRestClient(
          postResult: Result.success(
            RestClientResponse(statusCode: 200, data: _restrictedEnvelope()),
          ),
        ),
      );

      final result = await api.login(_loginRequest());

      expect(result, isA<Success<AuthState>>());
      expect(result.value, isA<RestrictedInstallationAuthState>());
      final restricted = result.value! as RestrictedInstallationAuthState;
      expect(restricted.restrictedAccessToken, 'restricted-token');
      expect(restricted.restrictedTokenType, 'restricted_access');
      expect(restricted.restrictedScope, 'installation.register');
      expect(
        restricted.restrictedExpiresAt,
        DateTime.parse('2026-06-17T10:05:00Z'),
      );
      expect(restricted.email, 'user@example.com');
    });

    test('maps installation limit reached envelope to typed error', () async {
      final api = AuthApi(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'INSTALLATION_LIMIT_REACHED',
                  'message': 'installation limit reached',
                },
              },
            ),
          ),
        ),
      );

      final result = await api.login(_loginRequest());

      expect(result, isA<Failure<AuthState>>());
      expect(result.error?.code, AppErrorCode.installationLimitReached);
      expect(result.error?.message, 'installation limit reached');
    });

    test('returns parsing error for unknown success shape', () async {
      final api = AuthApi(
        _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': {'email': 'user@example.com'},
                'error': null,
              },
            ),
          ),
        ),
      );

      final result = await api.login(_loginRequest());

      expect(result, isA<Failure<AuthState>>());
      expect(result.error?.code, AppErrorCode.parsingError);
    });
  });
}

LoginRequestDto _loginRequest() {
  return LoginRequestDto(email: 'user@example.com', password: '123456');
}

Map<String, dynamic> _operationalEnvelope() => {
  'data': {
    'access_token': 'access-token',
    'refresh_token': 'refresh-token',
    'user_id': 'd3de5f8b-4892-42e8-9680-979cf3f37844',
    'email': 'user@example.com',
    'role': 'customer',
    'customer_id': '6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3',
  },
  'error': null,
};

Map<String, dynamic> _restrictedEnvelope() => {
  'data': {
    'restricted_access_token': 'restricted-token',
    'restricted_token_type': 'restricted_access',
    'restricted_scope': 'installation.register',
    'restricted_expires_at': '2026-06-17T10:05:00Z',
    'user_id': 'd3de5f8b-4892-42e8-9680-979cf3f37844',
    'email': 'user@example.com',
    'role': 'customer',
    'customer_id': '6f3ebf86-bf82-4b75-a2ce-cd261ca47ec3',
  },
  'error': null,
};

class _FakeRestClient implements RestClient {
  _FakeRestClient({required this.postResult});

  final Result<RestClientResponse> postResult;
  RestClientRequest? lastPostRequest;

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async {
    lastPostRequest = request;
    return postResult;
  }

  @override
  AsyncResult<RestClientResponse> get(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> put(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> patch(RestClientRequest request) async =>
      throw UnimplementedError();

  @override
  AsyncResult<RestClientResponse> delete(RestClientRequest request) async =>
      throw UnimplementedError();
}
