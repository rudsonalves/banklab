import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/ui/pages/auth/register/register_cpf_page.dart';
import '/ui/pages/auth/register/viewmodel/register_viewmodel.dart';

List<RouteBase> registerRoutes() => [
  GoRoute(
    path: RegisterRoutes.cpf.path,
    name: RegisterRoutes.cpf.name,
    builder: (context, state) => RegisterCpfPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),
];
