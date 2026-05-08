enum AuthRoutes {
  login('/login'),
  register('/register')
  ;

  const AuthRoutes(this.path);

  final String path;
}

enum HomeRoutes {
  home('/home'),
  transfer('/transfer')
  ;

  const HomeRoutes(this.path);

  final String path;
}
