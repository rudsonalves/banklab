import '/core/result/result.dart';
import 'models/last_login_identity.dart';

abstract class LastLoginCacheService {
  AsyncResult<LastLoginIdentity> get();
  AsyncResult<Unit> save(LastLoginIdentity identity);
  AsyncResult<Unit> clear();
}
