import 'package:flutter/material.dart';

import 'big_button.dart';
import 'big_text_button.dart';

class DoubleBottomButton extends StatelessWidget {
  final String leftButtonLabel;
  final String rightButtonLabel;
  final Widget? rightButtonIcon;
  final VoidCallback? leftOnPressed;
  final VoidCallback? rightOnPressed;
  final bool isRightEnabled;

  const DoubleBottomButton({
    super.key,
    required this.leftButtonLabel,
    required this.rightButtonLabel,
    this.leftOnPressed,
    this.rightOnPressed,
    required this.isRightEnabled,
    this.rightButtonIcon,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      spacing: 12,
      children: [
        Expanded(
          flex: 1,
          child: BigTextButton(
            onPressed: leftOnPressed,
            label: leftButtonLabel,
          ),
        ),
        Expanded(
          flex: 2,
          child: BigButton(
            label: rightButtonLabel,
            onPressed: rightOnPressed,
            rightIcon: rightButtonIcon,
            enabled: isRightEnabled,
          ),
        ),
      ],
    );
  }
}
