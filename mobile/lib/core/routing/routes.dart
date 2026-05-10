enum AuthRoutes {
  login('/login'),
  register('/register')
  ;

  const AuthRoutes(this.path);

  final String path;
}

enum HomeRoutes {
  home('/home')
  ;

  const HomeRoutes(this.path);

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
  status('/status')
  ;

  const TransferRoutes(this.path);

  final String path;
}
