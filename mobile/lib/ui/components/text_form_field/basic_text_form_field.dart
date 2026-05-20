import 'package:flutter/material.dart';

class BasicTextFormField extends TextFormField {
  BasicTextFormField({
    super.key,
    required TextEditingController super.controller,
    String? labelText,
    String? hintText,
    TextInputType super.keyboardType = TextInputType.text,
    bool super.enabled = true,
    FloatingLabelBehavior floatingLabelBehavior = FloatingLabelBehavior.never,
    super.style = const TextStyle(fontSize: 16, fontWeight: FontWeight.w700),
    super.obscureText,
    super.autofillHints,
    super.textCapitalization,
    super.inputFormatters,
    super.textInputAction,
    super.validator,
    super.onFieldSubmitted,
    super.onChanged,
    Widget? prefixIcon,
    Widget? suffixIcon,
  }) : super(
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
           floatingLabelBehavior: floatingLabelBehavior,

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
