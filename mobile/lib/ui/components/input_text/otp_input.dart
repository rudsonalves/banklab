import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

import '/core/extensions/string.dart';

class OtpInput extends StatefulWidget {
  final int lenth;
  final String? initialValue;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onCompleted;
  final double spacing;
  final double fieldWidth;
  final double fieldHeight;
  final bool autoFocus;

  const OtpInput({
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
  State<OtpInput> createState() => _OtpInputState();
}

class _OtpInputState extends State<OtpInput> {
  late final TextEditingController _controller;
  late final FocusNode _focusNode;

  String get value => _controller.text;

  @override
  void initState() {
    super.initState();

    _controller = TextEditingController(
      text: (widget.initialValue ?? '').onlyNumbers.substring(
        0,
        ((widget.initialValue ?? '').onlyNumbers.length).clamp(
          0,
          widget.lenth,
        ),
      ),
    );
    _focusNode = FocusNode()..addListener(_onFocusChanged);
    _controller.addListener(_onTextChanged);

    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (widget.autoFocus) _focusNode.requestFocus();
    });
  }

  @override
  void dispose() {
    _controller.removeListener(_onTextChanged);
    _focusNode.removeListener(_onFocusChanged);
    _controller.dispose();
    _focusNode.dispose();

    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: () => _focusNode.requestFocus(),
      onLongPress: _pasteFromClipboard,
      child: SizedBox(
        height: widget.fieldHeight,
        child: Stack(
          alignment: Alignment.center,
          children: [
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: List.generate(
                widget.lenth,
                (index) => _OtpCell(
                  char: index < value.length ? value[index] : null,
                  isFocused:
                      _focusNode.hasFocus &&
                      index == value.length.clamp(0, widget.lenth - 1),
                  width: widget.fieldWidth,
                  height: widget.fieldHeight,
                ),
              ),
            ),
            Positioned.fill(
              child: IgnorePointer(
                child: Opacity(
                  opacity: 0,
                  child: TextField(
                    controller: _controller,
                    focusNode: _focusNode,
                    keyboardType: TextInputType.number,
                    autofillHints: const [AutofillHints.oneTimeCode],
                    enableInteractiveSelection: false,
                    showCursor: false,
                    textAlign: TextAlign.center,
                    inputFormatters: [
                      FilteringTextInputFormatter.digitsOnly,
                      LengthLimitingTextInputFormatter(widget.lenth),
                    ],
                    maxLength: widget.lenth,
                    decoration: const InputDecoration(
                      border: InputBorder.none,
                      counterText: '',
                      contentPadding: EdgeInsets.zero,
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _onFocusChanged() => setState(() {});

  void _onTextChanged() {
    setState(() {});

    widget.onChanged?.call(value);

    if (value.length == widget.lenth) {
      widget.onCompleted?.call(value);
      _focusNode.unfocus();
    }
  }

  Future<void> _pasteFromClipboard() async {
    final data = await Clipboard.getData(Clipboard.kTextPlain);
    final digits = data?.text?.onlyNumbers;

    if (digits == null || digits.isEmpty) return;

    _controller.text = digits.substring(
      0,
      digits.length.clamp(0, widget.lenth),
    );
    _controller.selection = TextSelection.collapsed(
      offset: _controller.text.length,
    );

    if (_controller.text.length < widget.lenth) {
      _focusNode.requestFocus();
    }
  }
}

class _OtpCell extends StatelessWidget {
  final String? char;
  final bool isFocused;
  final double width;
  final double height;

  const _OtpCell({
    required this.char,
    required this.isFocused,
    required this.width,
    required this.height,
  });

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return SizedBox(
      width: width,
      height: height,
      child: DecoratedBox(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(8),
          border: Border.all(
            color: isFocused
                ? theme.colorScheme.primary
                : theme.colorScheme.outline,
          ),
        ),
        child: Center(
          child: Text(
            char ?? '',
            style: const TextStyle(fontWeight: FontWeight.w700),
          ),
        ),
      ),
    );
  }
}
