import 'package:flutter/material.dart';

class BigTextButton extends StatelessWidget {
  final VoidCallback? onPressed;
  final String label;
  final Widget? rightIcon;
  final Widget? leftIcon;

  const BigTextButton({
    super.key,
    required this.onPressed,
    required this.label,
    this.rightIcon,
    this.leftIcon,
  });

  @override
  Widget build(BuildContext context) {
    return TextButton(
      onPressed: onPressed,
      style: TextButton.styleFrom(
        padding: const EdgeInsets.symmetric(vertical: 14),
        textStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w700,
        ),
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
        ),
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          if (leftIcon != null) ...[
            leftIcon!,
            const SizedBox(width: 8),
          ],
          Text(label),
          if (rightIcon != null) ...[
            const SizedBox(width: 8),
            rightIcon!,
          ],
        ],
      ),
    );
  }
}
