import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import 'extra_codec.dart';
import 'route_observer.dart';
import 'routes/auth_routes.dart';
import 'routes/base_routes.dart';
import 'routes/register_routes.dart';
import 'routes/shared_routes.dart';
import 'routes/transaction_password_routes.dart';
import 'routes/transfer_routes.dart';

GoRouter router() => GoRouter(
  initialLocation: BaseRoutes.splash.path,
  debugLogDiagnostics: kDebugMode,
  observers: [routeObserver],
  extraCodec: const ExtraCodec(),
  routes: [
    ...baseRoutes(),
    ...registerRoutes(),
    ...authRoutes(),
    ...transactionPasswordRoutes(),
    ...transferRoutes(),
    ...sharedRoutes(),
  ],
);
