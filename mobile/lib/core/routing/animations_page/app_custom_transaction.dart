import 'package:flutter/widgets.dart';
import 'package:go_router/go_router.dart';

class AppCustomTransactionPage<T> extends CustomTransitionPage<T> {
  AppCustomTransactionPage({
    required LocalKey super.key,
    required super.child,
    super.transitionDuration = const Duration(milliseconds: 400),
    super.reverseTransitionDuration = const Duration(milliseconds: 300),
  }) : super(
         transitionsBuilder: (context, animation, secondaryAnimation, child) {
           final curved = CurvedAnimation(
             parent: animation,
             curve: Curves.easeOutCubic,
           );

           return FadeTransition(
             opacity: CurvedAnimation(
               parent: animation,
               curve: Curves.easeOut,
             ),
             child: ScaleTransition(
               scale: Tween(begin: 0.95, end: 1.0).animate(curved),
               child: child,
             ),
           );
         },
       );
}
