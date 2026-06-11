import 'package:go_router/go_router.dart';

import '/core/config/dependencies.dart';
import '/core/routing/routes.dart';
import '/data/services/apis/contact_verification/enums/contact_verification_channel.dart';
import '/ui/pages/register/register_birthdate_page.dart';
import '/ui/pages/register/register_cpf_page.dart';
import '/ui/pages/register/register_email_page.dart';
import '/ui/pages/register/register_name_page.dart';
import '/ui/pages/register/register_password_page.dart';
import '/ui/pages/register/register_phone_page.dart';
import '/ui/pages/register/register_status_page.dart';
import '/ui/pages/register/register_token_page.dart';
import '/ui/pages/register/viewmodel/register_viewmodel.dart';

List<RouteBase> registerRoutes() => [
  GoRoute(
    path: RegisterRoutes.cpf.routePath,
    name: RegisterRoutes.cpf.routeName,
    builder: (context, state) => RegisterCpfPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.fullName.routePath,
    name: RegisterRoutes.fullName.routeName,
    builder: (context, state) => RegisterNamePage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.birthDate.routePath,
    name: RegisterRoutes.birthDate.routeName,
    builder: (context, state) => RegisterBirthdatePage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.email.routePath,
    name: RegisterRoutes.email.routeName,
    builder: (context, state) => RegisterEmailPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.emailToken.routePath,
    name: RegisterRoutes.emailToken.routeName,
    builder: (context, state) => RegisterTokenPage(
      viewmodel: injector.get<RegisterViewmodel>(),
      channel: ContactVerificationChannel.email,
    ),
  ),

  GoRoute(
    path: RegisterRoutes.phone.routePath,
    name: RegisterRoutes.phone.routeName,
    builder: (context, state) => RegisterPhonePage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.phoneToken.routePath,
    name: RegisterRoutes.phoneToken.routeName,
    builder: (context, state) => RegisterTokenPage(
      viewmodel: injector.get<RegisterViewmodel>(),
      channel: ContactVerificationChannel.phone,
    ),
  ),

  GoRoute(
    path: RegisterRoutes.password.routePath,
    name: RegisterRoutes.password.routeName,
    builder: (context, state) => RegisterPasswordPage(
      viewmodel: injector.get<RegisterViewmodel>(),
    ),
  ),

  GoRoute(
    path: RegisterRoutes.success.routePath,
    name: RegisterRoutes.success.routeName,
    builder: (context, state) => const RegisterStatusPage(
      isSuccess: true,
    ),
  ),

  GoRoute(
    path: RegisterRoutes.failure.routePath,
    name: RegisterRoutes.failure.routeName,
    builder: (context, state) => const RegisterStatusPage(
      isSuccess: false,
    ),
  ),
];
