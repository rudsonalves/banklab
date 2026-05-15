import 'package:flutter/material.dart';

class BigButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final bool enabled;
  final Widget? rightIcon;
  final Widget? leftIcon;

  const BigButton({
    super.key,
    required this.label,
    this.onPressed,
    this.enabled = true,
    this.rightIcon,
    this.leftIcon,
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
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (leftIcon != null) ...[
              leftIcon!,
              SizedBox(width: 8),
            ],
            Text(label, style: textStyle),
            if (rightIcon != null) ...[
              SizedBox(width: 8),
              rightIcon!,
            ],
          ],
        ),
      ),
    );
  }
}
