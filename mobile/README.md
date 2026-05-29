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

- authentication with JWT, including approval-pending login states before access is granted
- multi-step account creation onboarding and account lifecycle interactions
- transfer between accounts; the current API contract requires
  `X-Step-Up-Token` from `POST /security/step-up/authorize` for
  `POST /accounts/internal-transfers` with `endpoint_key=internal_transfer.create`
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
make mobile-tests
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
|   `-- ui/
|-- test/
|-- android/
|-- ios/
|-- web/
`-- pubspec.yaml
```

Dependency injection entrypoints currently live in:

- `lib/data/repositories.dart`
- `lib/domain/usecases/usecases.dart`
- `lib/ui/viewmodels.dart`

## Related docs

- Monorepo overview: [../README.md](../README.md)
- API service guide: [../api/README.md](../api/README.md)
- API getting started: [../api/docs/00-getting_started.md](../api/docs/00-getting_started.md)
- Mobile getting started: [docs/00-getting_started.md](docs/00-getting_started.md)
- Mobile architecture: [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)

## License

MIT. See [LICENSE](LICENSE).
