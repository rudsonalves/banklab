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
      (index) => FocusNode(
        onKeyEvent: (_, event) => _handleKeyEvent(index, event),
      ),
    );

    for (int i = 0; i < widget.lenth; i++) {
      _focusNodes[i].addListener(() {
        if (_focusNodes[i].hasFocus) {
          _selectCellContent(i);
        }
      });
    }

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
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: List.generate(
        widget.lenth,
        (index) => SizedBox(
          width: widget.fieldWidth,
          height: widget.fieldHeight,
          child: TextField(
            controller: _controllers[index],
            focusNode: _focusNodes[index],
            keyboardType: TextInputType.number,
            textAlign: TextAlign.center,
            textInputAction: TextInputAction.next,
            autofillHints: const [AutofillHints.oneTimeCode],
            style: const TextStyle(fontWeight: FontWeight.w700),
            decoration: const InputDecoration(
              counterText: '',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.all(Radius.circular(8)),
              ),
            ),
            inputFormatters: [
              FilteringTextInputFormatter.digitsOnly,
              LengthLimitingTextInputFormatter(1),
            ],
            maxLength: 1,
            maxLines: 1,
            onChanged: (text) => _handleChanged(index, text),
            onTap: _onTap,
          ),
        ),
      ),
    );
  }

  Future<void> _onTap() async {
    final data = await Clipboard.getData(Clipboard.kTextPlain);
    final text = data?.text;
    if (text == null || text.length <= 1) return;

    await _handlePaste(text);
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

    if (firstEmptyIndex == -1) return;

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
      if (index > 0) {
        _focusNodes[index - 1].requestFocus();
      }

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
      if (index > 0) {
        _focusNodes[index - 1].requestFocus();
      }
      _notifyChanges();
      return KeyEventResult.handled;
    }

    if (index > 0) {
      _focusAndSelectCellContent(index - 1);
    }

    return KeyEventResult.handled;
  }

  void _focusAndSelectCellContent(int index) {
    _focusNodes[index].requestFocus();
    _selectCellContent(index);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (!mounted) return;

      _selectCellContent(index);
    });
  }

  void _selectCellContent(int index) {
    _controllers[index].selection = TextSelection(
      baseOffset: 0,
      extentOffset: _controllers[index].text.length,
    );
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
