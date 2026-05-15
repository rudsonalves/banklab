import 'package:flutter/material.dart';

class TextHeader extends StatelessWidget {
  final String text;

  const TextHeader(this.text, {super.key});

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 4),
      child: Text(
        text,
        style: Theme.of(context).textTheme.bodyLarge!.copyWith(
          fontWeight: FontWeight.w700,
        ),
      ),
    );
  }
}
