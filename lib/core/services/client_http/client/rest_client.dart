import 'package:bankflow/core/result/result.dart';

import 'rest_client_request.dart';
import 'rest_client_response.dart';

abstract class RestClient {
  void setBaseUrl(String url);

  AsyncResult<RestClientResponse> get(RestClientRequest request);
  AsyncResult<RestClientResponse> post(RestClientRequest request);
  AsyncResult<RestClientResponse> put(RestClientRequest request);
  AsyncResult<RestClientResponse> patch(RestClientRequest request);
  AsyncResult<RestClientResponse> delete(RestClientRequest request);
}
