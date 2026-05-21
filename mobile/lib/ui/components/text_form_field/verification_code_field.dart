import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '/core/extensions/string.dart';

class VerificationCodeField extends StatefulWidget {
  final int lenth;
  final String? initialValue;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onCompleted;
  final double spacing;
  final double fieldWidth;
  final double fieldHeight;
  final bool autoFocus;

  const VerificationCodeField({
    super.key,
    this.lenth = 6,
    this.initialValue,
    this.onChanged,
    this.onCompleted,
    this.spacing = 12,
    this.fieldWidth = 53,
    this.fieldHeight = 60,
    this.autoFocus = true,
  });

  @override
  State<VerificationCodeField> createState() => _VerificationCodeFieldState();
}

class _VerificationCodeFieldState extends State<VerificationCodeField> {
  late final List<TextEditingController> _controllers;
  late final List<FocusNode> _focusNodes;

  String get value => _controllers.map((c) => c.text).join();

  @override
  void initState() {
    super.initState();

    _controllers = List.generate(
      widget.lenth,
      (index) => TextEditingController(),
    );

    _focusNodes = List.generate(
      widget.lenth,
      (index) => FocusNode(),
    );

    _fillInitialValue();

    WidgetsBinding.instance.addPostFrameCallback((_) {
      _requestInitialFocus();
    });
  }

  @override
  void dispose() {
    for (final controller in _controllers) {
      controller.dispose();
    }

    for (final focusNode in _focusNodes) {
      focusNode.dispose();
    }

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: List.generate(
        widget.lenth,
        (index) => Padding(
          padding: EdgeInsets.only(
            right: index == widget.lenth - 1 ? 0 : widget.spacing,
          ),
          child: Focus(
            onKeyEvent: (_, event) => _handleKeyEvent(index, event),
            child: SizedBox(
              width: widget.fieldWidth,
              height: widget.fieldHeight,
              child: TextField(
                controller: _controllers[index],
                focusNode: _focusNodes[index],
                keyboardType: TextInputType.number,
                textAlign: TextAlign.center,
                textInputAction: TextInputAction.next,
                autofillHints: const [AutofillHints.oneTimeCode],
                style: const TextStyle(
                  fontSize: 24,
                  fontWeight: FontWeight.w600,
                ),
                decoration: const InputDecoration(
                  counterText: '',
                ),
                inputFormatters: [
                  FilteringTextInputFormatter.digitsOnly,
                  LengthLimitingTextInputFormatter(1),
                ],
                maxLines: 1,
                onChanged: (text) => _handleChanged(index, text),
                onTap: () async {
                  final data = await Clipboard.getData(Clipboard.kTextPlain);

                  final text = data?.text;

                  if (text == null || text.length <= 1) return;

                  await _handlePaste(text);
                },
              ),
            ),
          ),
        ),
      ),
    );
  }

  void _fillInitialValue() {
    final initial = widget.initialValue ?? '';

    for (int i = 0; i < widget.lenth; i++) {
      if (i < initial.length) {
        _controllers[i].text = initial[i];
      }
    }
  }

  void _requestInitialFocus() {
    if (!widget.autoFocus) return;

    final firstEmptyIndex = _controllers.indexWhere(
      (controller) => controller.text.isEmpty,
    );

    if (firstEmptyIndex != -1) return;

    _focusNodes[firstEmptyIndex].requestFocus();
  }

  void _notifyChanges() {
    final currentValue = value;

    widget.onChanged?.call(currentValue);

    if (currentValue.length == widget.lenth &&
        !_controllers.any((controller) => controller.text.isEmpty)) {
      widget.onCompleted?.call(currentValue);
    }
  }

  void _handleChanged(int index, String text) {
    if (text.isEmpty) {
      _notifyChanges();
      return;
    }

    _controllers[index].text = text;
    _controllers[index].selection = const TextSelection.collapsed(offset: 1);

    if (index < widget.lenth - 1) {
      _focusNodes[index + 1].requestFocus();
    } else {
      _focusNodes[index].unfocus();
    }

    _notifyChanges();
  }

  KeyEventResult _handleKeyEvent(int index, KeyEvent event) {
    if (event is! KeyDownEvent) return KeyEventResult.ignored;

    if (event.logicalKey != LogicalKeyboardKey.backspace) {
      return KeyEventResult.ignored;
    }

    final controller = _controllers[index];

    if (controller.text.isNotEmpty) {
      controller.clear();
      _notifyChanges();
      return KeyEventResult.handled;
    }

    if (index > 0) {
      _focusNodes[index - 1].requestFocus();
      _controllers[index - 1].clear();
      _notifyChanges();
    }

    return KeyEventResult.handled;
  }

  Future<void> _handlePaste(String text) async {
    final digits = text.onlyNumbers;

    if (digits.isEmpty) return;

    for (int i = 0; i < widget.lenth; i++) {
      _controllers[i].text = i < digits.length ? digits[i] : '';
    }

    final nextFocus = digits.length.clamp(0, widget.lenth - 1);

    if (digits.length >= widget.lenth) {
      _focusNodes.last.unfocus();
    } else {
      _focusNodes[nextFocus].requestFocus();
    }

    _notifyChanges();
  }
}
