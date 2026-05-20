import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/ui/pages/register/register_birthdate_page.dart';
import '/ui/pages/register/register_cpf_page.dart';
import '/ui/pages/register/register_email_page.dart';
import '/ui/pages/register/register_name_page.dart';
import '/ui/pages/register/register_password_page.dart';
import '/ui/pages/register/register_phone_page.dart';
import '/ui/pages/register/register_token_page.dart';
import '/ui/pages/register/viewmodel/register_viewmodel.dart';

List<RouteBase> registerRoutes() => [
  GoRoute(
    path: RegisterRoutes.cpf.path,
    name: RegisterRoutes.cpf.name,
    builder: (context, state) => RegisterCpfPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.name.path,
    name: RegisterRoutes.name.name,
    builder: (context, state) => RegisterNamePage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.birthDate.path,
    name: RegisterRoutes.birthDate.name,
    builder: (context, state) => RegisterBirthdatePage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.email.path,
    name: RegisterRoutes.email.name,
    builder: (context, state) => RegisterEmailPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.emailToken.path,
    name: RegisterRoutes.emailToken.name,
    builder: (context, state) => RegisterTokenPage(
      viewmodel: injector.get<RegisterViewmodel>(),
      tokenType: TokenType.email,
    ),
  ),

  GoRoute(
    path: RegisterRoutes.phone.path,
    name: RegisterRoutes.phone.name,
    builder: (context, state) => RegisterPhonePage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.phoneToken.path,
    name: RegisterRoutes.phoneToken.name,
    builder: (context, state) => RegisterTokenPage(
      viewmodel: injector.get<RegisterViewmodel>(),
      tokenType: TokenType.phone,
    ),
  ),

  GoRoute(
    path: RegisterRoutes.password.path,
    name: RegisterRoutes.password.name,
    builder: (context, state) => RegisterPasswordPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),
];
