import 'package:bankflow/core/routing/routes.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

import 'extra_codec.dart';

GoRouter router() => GoRouter(
  initialLocation: Routes.home.name,
  debugLogDiagnostics: kDebugMode,
  extraCodec: const ExtraCodec(),
  routes: [
    GoRoute(
      path: Routes.home.name,
      name: Routes.home.name,
      builder: (context, state) => const SizedBox.shrink(),
    ),
  ],
);
