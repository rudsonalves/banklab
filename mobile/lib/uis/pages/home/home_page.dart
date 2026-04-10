import 'package:flutter/material.dart';

import 'viewmodel/home_viewmodel.dart';

class HomePage extends StatefulWidget {
  final HomeViewmodel viewModel;

  const HomePage({super.key, required this.viewModel});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Type Home Page'),
      ),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Text('You have pushed the button this many times:'),
          ],
        ),
      ),
    );
  }
}
