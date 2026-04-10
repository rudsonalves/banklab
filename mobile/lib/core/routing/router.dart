import 'package:flutter/foundation.dart';
import 'package:go_router/go_router.dart';

import '/core/routing/routes.dart';
import 'extra_codec.dart';
import 'routes/auth_routes.dart';
import 'routes/home_routes.dart';

GoRouter router() => GoRouter(
  initialLocation: AuthRoutes.login.path,
  debugLogDiagnostics: kDebugMode,
  extraCodec: const ExtraCodec(),
  routes: [
    ...homeRoutes(),
    ...authRoutes(),
  ],
);
