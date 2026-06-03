import 'package:flutter/material.dart';
import 'package:flutter/services.dart';

class BasicInputText extends StatelessWidget {
  final TextEditingController controller;
  final FocusNode? focusNode;
  final String? labelText;
  final String? hintText;
  final TextInputType? keyboardType;
  final bool enabled;
  final TextStyle style;
  final bool? obscureText;
  final Iterable<String>? autofillHints;
  final TextCapitalization textCapitalization;
  final List<TextInputFormatter>? inputFormatters;
  final TextInputAction? textInputAction;
  final String? Function(String?)? validator;
  final void Function(String)? onFieldSubmitted;
  final void Function(String)? onChanged;
  final Widget? prefixIcon;
  final Widget? suffixIcon;

  const BasicInputText({
    super.key,
    required this.controller,
    this.focusNode,
    this.labelText,
    this.hintText,
    this.keyboardType,
    this.enabled = true,
    this.style = const TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
    this.obscureText,
    this.autofillHints,
    this.textCapitalization = TextCapitalization.none,
    this.inputFormatters,
    this.textInputAction,
    this.validator,
    this.onFieldSubmitted,
    this.onChanged,
    this.prefixIcon,
    this.suffixIcon,
  });

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: controller,
      focusNode: focusNode,
      keyboardType: keyboardType,
      enabled: enabled,
      style: style,
      obscureText: obscureText ?? false,
      autofillHints: autofillHints,
      textCapitalization: textCapitalization,
      inputFormatters: inputFormatters,
      textInputAction: textInputAction,
      validator: validator,
      onFieldSubmitted: onFieldSubmitted,
      onChanged: onChanged,
      decoration: InputDecoration(
        labelText: labelText,
        hintText: hintText,
        hintStyle: const TextStyle(
          fontSize: 16,
          fontWeight: FontWeight.w400,
          color: Colors.grey,
        ),
        prefixIcon: prefixIcon,
        suffixIcon: suffixIcon,

        filled: true,
        border: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: BorderSide.none,
        ),
        enabledBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: BorderSide.none,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: BorderRadius.circular(8),
          borderSide: BorderSide.none,
        ),
      ),
    );
  }
}

// import 'package:flutter/material.dart';

// class BasicInputText extends TextFormField {
//   BasicInputText({
//     super.key,
//     required TextEditingController super.controller,
//     super.focusNode,
//     String? labelText,
//     String? hintText,
//     TextInputType super.keyboardType = TextInputType.text,
//     bool super.enabled = true,
//     FloatingLabelBehavior floatingLabelBehavior = FloatingLabelBehavior.never,
//     super.style = const TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
//     super.obscureText,
//     super.autofillHints,
//     super.textCapitalization,
//     super.inputFormatters,
//     super.textInputAction,
//     super.validator,
//     super.onFieldSubmitted,
//     super.onChanged,
//     Widget? prefixIcon,
//     Widget? suffixIcon,
//   }) : super(
//          decoration: InputDecoration(
//            labelText: labelText,
//            hintText: hintText,
//            hintStyle: const TextStyle(
//              fontSize: 16,
//              fontWeight: FontWeight.w400,
//              color: Colors.grey,
//            ),
//            prefixIcon: prefixIcon,
//            suffixIcon: suffixIcon,
//            floatingLabelBehavior: floatingLabelBehavior,

//            filled: true,
//            border: OutlineInputBorder(
//              borderRadius: BorderRadius.circular(8),
//              borderSide: BorderSide.none,
//            ),
//            enabledBorder: OutlineInputBorder(
//              borderRadius: BorderRadius.circular(8),
//              borderSide: BorderSide.none,
//            ),
//            focusedBorder: OutlineInputBorder(
//              borderRadius: BorderRadius.circular(8),
//              borderSide: BorderSide.none,
//            ),
//          ),
//        );
// }
