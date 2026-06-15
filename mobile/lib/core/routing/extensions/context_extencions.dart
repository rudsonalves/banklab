import 'package:flutter/widgets.dart';

extension NavigatorPopUntilExtension on BuildContext {
  void popUntil(String routeName) {
    Navigator.of(this).popUntil(
      (route) => route.settings.name == routeName || route.isFirst,
    );
  }
}
