# BankFlow (mobile)

BankFlow is the Flutter client of the banklab monorepo. Its primary role is to validate end-to-end behavior of the banking flows exposed by the Go API.

This is an engineering-oriented app focused on integration quality, predictable state flow, and API contract validation.

## Stack

- Flutter
- Dart SDK ^3.11.4
- dio
- go_router
- flutter_secure_storage
- intl

## Main flows

- authentication with JWT
- account creation and account lifecycle interactions
- deposit and withdraw operations
- transfer between accounts
- transaction history visualization

## Local setup

From repository root:

```bash
cd mobile
flutter pub get
```

Run in debug mode:

```bash
flutter run
```

### Environment files

The project includes:

- dev.env
- staging.env
- prod.env

Adjust or load the environment expected by your run configuration so the app points to the correct API base URL.

## Running tests

From repository root:

```bash
make mobile-test
make mobile-test-unit
```

Or directly from the mobile directory:

```bash
cd mobile
flutter test
flutter test test/core
```

## Build helpers

From repository root:

```bash
make fclean
make fbuild
```

- fclean: flutter clean + flutter pub get
- fbuild: release APK build

## Project structure (summary)

```text
mobile/
|-- lib/
|   |-- core/
|   |-- data/
|   |-- domain/
|   `-- uis/
|-- test/
|-- android/
|-- ios/
|-- web/
`-- pubspec.yaml
```

## Related docs

- Monorepo overview: [../README.md](../README.md)
- API service guide: [../api/README.md](../api/README.md)
- Mobile architecture: [../docs/mobile/ARCHITECTURE.md](../docs/mobile/ARCHITECTURE.md)

## License

MIT. See [LICENSE](LICENSE).