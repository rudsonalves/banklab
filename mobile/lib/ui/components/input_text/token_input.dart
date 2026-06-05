import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class TokenInput extends StatefulWidget {
  final int length;
  final String? initialValue;
  final ValueChanged<String>? onChanged;
  final ValueChanged<String>? onCompleted;
  final bool visible;
  final bool autoFocus;
  final double spacing;
  final double cellSize;

  const TokenInput({
    super.key,
    this.length = 6,
    this.initialValue,
    this.onChanged,
    this.onCompleted,
    this.visible = false,
    this.autoFocus = true,
    this.spacing = 12,
    this.cellSize = 48,
  });

  @override
  State<TokenInput> createState() => _TokenInputState();
}

class _TokenInputState extends State<TokenInput> {
  late final TextEditingController _controller;
  late final FocusNode _focusNode;

  String get _value => _controller.text;

  @override
  void initState() {
    super.initState();
    _controller = TextEditingController(
      text: widget.initialValue?.substring(
        0,
        (widget.initialValue!.length).clamp(0, widget.length),
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

  void _onTextChanged() {
    setState(() {});
    widget.onChanged?.call(_value);
    if (_value.length == widget.length) widget.onCompleted?.call(_value);
  }

  void _onFocusChanged() => setState(() {});

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      behavior: HitTestBehavior.opaque,
      onTap: () => _focusNode.requestFocus(),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Offstage(
            child: TextField(
              controller: _controller,
              focusNode: _focusNode,
              keyboardType: TextInputType.number,
              autofillHints: const [AutofillHints.oneTimeCode],
              inputFormatters: [
                FilteringTextInputFormatter.digitsOnly,
                LengthLimitingTextInputFormatter(widget.length),
              ],
              maxLength: widget.length,
              decoration: const InputDecoration(counterText: ''),
            ),
          ),
          Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              for (int i = 0; i < widget.length; i++) ...[
                if (i > 0) SizedBox(width: widget.spacing),
                _TokenCell(
                  char: i < _value.length ? _value[i] : null,
                  isFocused:
                      _focusNode.hasFocus &&
                      i == _value.length.clamp(0, widget.length - 1),
                  visible: widget.visible,
                  size: widget.cellSize,
                ),
              ],
            ],
          ),
        ],
      ),
    );
  }
}

class _TokenCell extends StatelessWidget {
  final String? char;
  final bool isFocused;
  final bool visible;
  final double size;

  const _TokenCell({
    required this.char,
    required this.isFocused,
    required this.visible,
    required this.size,
  });

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final isFilled = char != null;

    return AnimatedContainer(
      duration: const Duration(milliseconds: 150),
      width: size,
      height: size,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(8),
        color: colorScheme.surface,
        border: Border.all(
          color: isFocused ? colorScheme.primary : colorScheme.outlineVariant,
          width: isFocused ? 2 : 1,
        ),
      ),
      child: Center(
        child: isFilled
            ? Text(
                visible ? char! : '*',
                style: TextStyle(
                  fontSize: size * 0.42,
                  fontWeight: FontWeight.w700,
                  color: colorScheme.onSurface,
                ),
              )
            : Container(
                width: size * 0.22,
                height: size * 0.22,
                decoration: BoxDecoration(
                  shape: BoxShape.circle,
                  color: colorScheme.outlineVariant,
                ),
              ),
      ),
    );
  }
}
