import '/core/result/result.dart';

abstract class LocalSecureStorage {
  AsyncResult<Unit> write(String key, String value);
  AsyncResult<String> read(String key);
  AsyncResult<Unit> delete(String key);
  AsyncResult<Unit> deleteAll();
  Future<List<String>> keysWithPrefix(String pattern);
}
