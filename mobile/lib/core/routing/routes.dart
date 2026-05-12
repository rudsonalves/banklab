enum AuthRoutes {
  login('/login'),
  shortLogin('/short-login'),
  register('/register')
  ;

  const AuthRoutes(this.path);

  final String path;
}

enum BaseRoutes {
  home('/home'),
  splash('/splash'),
  statement('/statement')
  ;

  const BaseRoutes(this.path);

  final String path;
}

enum GeneralRoutes {
  splash('/splash'),
  receipt('/receipt')
  ;

  const GeneralRoutes(this.path);

  final String path;
}

enum TransferRoutes {
  recipient('/recipient'),
  payment('/payment'),
  confirmation('/confirmation'),
  statusSuccess('/status/success'),
  statusFailure('/status/failure')
  ;

  const TransferRoutes(this.path);

  final String path;
}

enum SharedRoutes {
  details('/details')
  ;

  const SharedRoutes(this.path);

  final String path;
}
