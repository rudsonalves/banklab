import 'package:flutter/material.dart';

class BigButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final bool enabled;
  final bool isRunning;
  final Widget? rightIcon;
  final Widget? leftIcon;

  const BigButton({
    super.key,
    required this.label,
    this.onPressed,
    this.enabled = true,
    this.isRunning = false,
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

    final leftWidget = isRunning
        ? SizedBox(
            width: 16,
            height: 16,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              valueColor: AlwaysStoppedAnimation(colorScheme.onPrimary),
            ),
          )
        : leftIcon;

    final rightWidget = isRunning
        ? SizedBox(
            width: 16,
            height: 16,
            child: CircularProgressIndicator(
              strokeWidth: 2,
              valueColor: AlwaysStoppedAnimation(colorScheme.onPrimary),
            ),
          )
        : rightIcon;

    return ElevatedButton(
      onPressed: enabled && !isRunning ? onPressed : null,
      style: style,
      child: Padding(
        padding: EdgeInsets.symmetric(vertical: 12),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            if (leftWidget != null) ...[
              leftWidget,
              SizedBox(width: 8),
            ],
            Text(label, style: textStyle),
            if (rightWidget != null) ...[
              SizedBox(width: 8),
              rightWidget,
            ],
          ],
        ),
      ),
    );
  }
}
