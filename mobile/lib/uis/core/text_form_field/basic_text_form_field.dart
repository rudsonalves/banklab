import 'package:flutter/material.dart';

class BasicTextFormField extends TextFormField {
  BasicTextFormField({
    super.key,
    required TextEditingController super.controller,
    required String labelText,
    String? hintText,
    TextInputType super.keyboardType = TextInputType.text,
    super.textCapitalization = TextCapitalization.none,
    bool super.enabled = true,
    FloatingLabelBehavior floatingLabelBehavior = FloatingLabelBehavior.never,
    super.obscureText,
    super.autofillHints,
    super.inputFormatters,
    super.textInputAction,
    super.validator,
    super.onFieldSubmitted,
    Widget? prefixIcon,
    Widget? suffixIcon,
  }) : super(
         decoration: InputDecoration(
           labelText: labelText,
           hintText: hintText,
           prefixIcon: prefixIcon,
           suffixIcon: suffixIcon,
           floatingLabelBehavior: floatingLabelBehavior,

           filled: true,
           border: OutlineInputBorder(
             borderRadius: BorderRadius.circular(24),
             borderSide: BorderSide.none,
           ),
           enabledBorder: OutlineInputBorder(
             borderRadius: BorderRadius.circular(24),
             borderSide: BorderSide.none,
           ),
           focusedBorder: OutlineInputBorder(
             borderRadius: BorderRadius.circular(24),
             borderSide: BorderSide.none,
           ),
         ),
       );
}
