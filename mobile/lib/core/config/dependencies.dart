import 'package:auto_injector/auto_injector.dart';

import '/data/data.dart';
import '/data/services/services.dart';
import '/uis/uis.dart';
import '../services/core_services.dart';

final injector = AutoInjector();
bool _initialized = false;

void setupDependencies() {
  if (_initialized) return;

  CoreServices.add(injector);
  Services.add(injector);
  Data.add(injector);
  // Usecase.add(injector);
  Uis.add(injector);

  injector.commit();
  _initialized = true;
}
