import 'package:flutter/material.dart';

class BigButton extends StatelessWidget {
  final String label;
  final VoidCallback onPressed;
  final bool enabled;

  const BigButton({
    super.key,
    required this.label,
    required this.onPressed,
    required this.enabled,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final style = enabled
        ? ButtonStyle(
            backgroundColor: WidgetStateProperty.all(
              colorScheme.onPrimaryFixedVariant,
            ),
          )
        : Theme.of(context).elevatedButtonTheme.style;
    final textStyle = TextStyle(fontSize: 16, fontWeight: FontWeight.w700);

    return ElevatedButton(
      onPressed: enabled ? onPressed : null,
      style: style,
      child: Padding(
        padding: EdgeInsets.symmetric(vertical: 12),
        child: Text(label, style: textStyle),
      ),
    );
  }
}
