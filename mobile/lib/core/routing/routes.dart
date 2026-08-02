abstract interface class AppRoute {
  String get routePath;
  String get routeName;
}

enum AuthRoutes implements AppRoute {
  login('/login'),
  shortLogin('/short-login'),
  installationCertification('/installation-certification');

  const AuthRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}

enum TransactionPasswordRoutes implements AppRoute {
  introduction('/transaction-password'),
  create('/transaction-password/create'),
  confirm('/transaction-password/confirm'),

  transactionPassword('/transaction-password/verify');

  const TransactionPasswordRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}

enum RegisterRoutes implements AppRoute {
  cpf('/register/cpf'),
  fullName('/register/name'),
  birthDate('/register/birth-date'),
  email('/register/email'),
  emailToken('/register/email-token'),
  phone('/register/phone'),
  phoneToken('/register/phone-token'),
  password('/register/password'),
  success('/register/success'),
  failure('/register/failure');

  const RegisterRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}

enum BaseRoutes implements AppRoute {
  home('/home'),
  splash('/splash'),
  statement('/statement');

  const BaseRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}

enum GeneralRoutes implements AppRoute {
  splash('/splash'),
  receipt('/receipt');

  const GeneralRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}

enum TransferRoutes implements AppRoute {
  recipient('/recipient'),
  payment('/payment'),
  confirmation('/confirmation'),
  statusSuccess('/status/success'),
  statusFailure('/status/failure');

  const TransferRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}

enum SharedRoutes implements AppRoute {
  details('/details');

  const SharedRoutes(this.routePath);

  @override
  final String routePath;

  @override
  String get routeName => name;
}
