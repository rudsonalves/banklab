import 'package:auto_injector/auto_injector.dart';

import '../../data/repositories.dart';
import '/data/services/services.dart';
import '../../ui/viewmodels.dart';
import '../../domain/usecases/usecases.dart';
import '../services/core_services.dart';

final injector = AutoInjector();
bool _initialized = false;

void setupDependencies() {
  if (_initialized) return;

  CoreServices.add(injector);
  Services.add(injector);
  Repositories.add(injector);
  Usecases.add(injector);
  Viewmodels.add(injector);

  injector.commit();
  _initialized = true;
}
