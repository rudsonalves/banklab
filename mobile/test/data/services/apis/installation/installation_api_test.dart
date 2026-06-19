import 'package:bankflow/core/resources/app_http_headers.dart';
import 'package:bankflow/core/result/result.dart';
import 'package:bankflow/core/services/client_http/client_http.dart';
import 'package:bankflow/core/services/installation_identity/installation_identity.dart';
import 'package:bankflow/core/services/secure_storage/local_secure_storage.dart';
import 'package:bankflow/data/services/apis/installation/dtos/installation_registration_response_dto.dart';
import 'package:bankflow/data/services/apis/installation/installation_api.dart';
import 'package:flutter_test/flutter_test.dart';

void main() {
  group('InstallationApi.register', () {
    test('calls POST /security/installations with required headers', () async {
      final client = _FakeRestClient(
        postResult: Result.success(
          RestClientResponse(statusCode: 201, data: _successEnvelope()),
        ),
      );
      final api = InstallationApi(
        client: client,
        installationIdentityService: _FakeInstallationIdentityService(),
      );

      final result = await api.register(
        restrictedAccessToken: 'restricted-token',
        stepUpToken: 'step-up-token',
      );

      expect(result, isA<Success<InstallationRegistrationResponseDto>>());
      expect(client.postCalls, 1);
      expect(client.lastPostRequest?.path, '/security/installations');
      expect(client.lastPostRequest?.headers, {
        AppHttpHeaders.authorization: 'Bearer restricted-token',
        AppHttpHeaders.installationId: _installationId,
        AppHttpHeaders.stepUpToken: 'step-up-token',
      });
      expect(client.lastPostRequest?.body, isNull);
      expect(
        client.lastPostRequest?.body,
        isNot(contains('transaction_password')),
      );

      final dto = result.value!;
      expect(dto.accessToken, 'access-token');
      expect(dto.refreshToken, 'refresh-token');
      expect(
        dto.installationResourceId,
        '2e4a8e20-272a-4e7b-b782-bc7f6b1d0442',
      );
      expect(dto.installationStatus, 'known');
    });

    test('does not call RestClient when installation identity fails', () async {
      final client = _FakeRestClient(
        postResult: Result.success(
          RestClientResponse(statusCode: 201, data: _successEnvelope()),
        ),
      );
      final api = InstallationApi(
        client: client,
        installationIdentityService: _FakeInstallationIdentityService(
          result: const Failure(
            AppError(
              code: AppErrorCode.storageError,
              message: 'identity failed',
            ),
          ),
        ),
      );

      final result = await api.register(
        restrictedAccessToken: 'restricted-token',
        stepUpToken: 'step-up-token',
      );

      expect(result, isA<Failure<InstallationRegistrationResponseDto>>());
      expect(result.error?.code, AppErrorCode.storageError);
      expect(client.postCalls, 0);
    });

    for (final errorCode in [
      'STEP_UP_TOKEN_REQUIRED',
      'STEP_UP_TOKEN_INVALID',
      'INSTALLATION_MISMATCH',
      'INVALID_INSTALLATION_ID',
    ]) {
      test('preserves $errorCode from RestClient failure', () async {
        final api = InstallationApi(
          client: _FakeRestClient(
            postResult: Result.failure(
              AppError(
                statusCode: errorCode == 'INSTALLATION_MISMATCH' ? 403 : 401,
                code: AppErrorCode.httpError,
                message: 'registration failed',
                details: {'code': errorCode},
              ),
            ),
          ),
          installationIdentityService: _FakeInstallationIdentityService(),
        );

        final result = await api.register(
          restrictedAccessToken: 'restricted-token',
          stepUpToken: 'step-up-token',
        );

        expect(result, isA<Failure<InstallationRegistrationResponseDto>>());
        expect(result.error?.message, 'registration failed');
        expect(backendErrorCode(result.error), errorCode);
      });
    }

    test('maps backend envelope error to AppError failure', () async {
      final api = InstallationApi(
        client: _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 200,
              data: {
                'data': null,
                'error': {
                  'code': 'STEP_UP_TOKEN_INVALID',
                  'message': 'invalid step-up token',
                },
              },
            ),
          ),
        ),
        installationIdentityService: _FakeInstallationIdentityService(),
      );

      final result = await api.register(
        restrictedAccessToken: 'restricted-token',
        stepUpToken: 'step-up-token',
      );

      expect(result, isA<Failure<InstallationRegistrationResponseDto>>());
      expect(result.error?.code, AppErrorCode.httpError);
      expect(result.error?.message, 'invalid step-up token');
      expect(backendErrorCode(result.error), 'STEP_UP_TOKEN_INVALID');
    });

    test('returns parsing error when response data is malformed', () async {
      final api = InstallationApi(
        client: _FakeRestClient(
          postResult: const Result.success(
            RestClientResponse(
              statusCode: 201,
              data: {
                'data': {'access_token': 'access-token'},
                'error': null,
              },
            ),
          ),
        ),
        installationIdentityService: _FakeInstallationIdentityService(),
      );

      final result = await api.register(
        restrictedAccessToken: 'restricted-token',
        stepUpToken: 'step-up-token',
      );

      expect(result, isA<Failure<InstallationRegistrationResponseDto>>());
      expect(result.error?.code, AppErrorCode.parsingError);
    });
  });
}

const _installationId = '018f7b82-4a3d-4f71-9ad7-dedc8e5b10c8';

Map<String, dynamic> _successEnvelope() => {
  'data': {
    'access_token': 'access-token',
    'refresh_token': 'refresh-token',
    'installation_resource_id': '2e4a8e20-272a-4e7b-b782-bc7f6b1d0442',
    'installation_status': 'known',
  },
  'error': null,
};

class _FakeInstallationIdentityService extends InstallationIdentityService {
  _FakeInstallationIdentityService({
    Result<String> result = const Success(_installationId),
  }) : _result = result,
       super(
         secureStorage: _NoopSecureStorage(),
         markerStore: _NoopMarkerStore(),
       );

  final Result<String> _result;

  @override
  AsyncResult<String> resolve() async => _result;
}

class _NoopMarkerStore implements InstallationMarkerStore {
  @override
  AsyncResult<bool> hasMarker() async => const Success(true);

  @override
  AsyncResult<Unit> markResolved() async => const Success(unit);
}

class _NoopSecureStorage implements LocalSecureStorage {
  @override
  AsyncResult<Unit> write(String key, String value) async =>
      const Success(unit);

  @override
  AsyncResult<String> read(String key) async => const Success(_installationId);

  @override
  AsyncResult<Unit> delete(String key) async => const Success(unit);

  @override
  AsyncResult<Unit> deleteAll() async => const Success(unit);

  @override
  Future<List<String>> keysWithPrefix(String pattern) async => const [];
}

class _FakeRestClient implements RestClient {
  _FakeRestClient({required this.postResult});

  final Result<RestClientResponse> postResult;
  RestClientRequest? lastPostRequest;
  int postCalls = 0;

  @override
  AsyncResult<RestClientResponse> post(RestClientRequest request) async {
    postCalls++;
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
