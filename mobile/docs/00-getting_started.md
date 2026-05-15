# Getting Started — Bank Mobile

## 1. Overview

This document describes how to run the mobile app locally from a clean environment.

The setup assumes:

- Flutter SDK installed and configured
- Device or emulator available (Android/iOS)
- Project API running (backend required for app flows)

The mobile app depends on the API for authentication, accounts, and transactions, so the backend must be reachable before starting the app.

---

## 2. Prerequisites

### 2.1 Environment variables

The app uses `dart-define` at build/run time.

Inside the `mobile/` directory, create (or adjust) the desired environment file.

Development example (`mobile/dev.env`):

```env
BASE_URL=http://192.168.0.17:8080

CONNECT_TIMEOUT=30000
RECEIVE_TIMEOUT=30000

APP_MODE=dev

APP_ACCESS_TOKEN=your_app_token_here
```

Description:

- **BASE_URL**
	- Base API URL consumed by the app
	- Must be reachable by the device/emulator

- **CONNECT_TIMEOUT**
	- HTTP connection timeout in milliseconds

- **RECEIVE_TIMEOUT**
	- HTTP response timeout in milliseconds

- **APP_MODE**
	- Execution mode (`dev`, `staging`, `prod`)

- **APP_ACCESS_TOKEN**
	- Token sent in `X-App-Token` header on onboarding routes

---

### 2.2 API and infrastructure

Start the infrastructure/API in the monorepo root:

```bash
make run
```

This ensures:

1. Docker validation
2. PostgreSQL startup
3. Database readiness wait
4. Migration application
5. API startup

---

## 3. Bootstrap (first run)

Inside the `mobile/` directory, install dependencies:

```bash
flutter pub get
```

---

## 4. Run the mobile app

Inside the `mobile/` directory, run with the environment file:

```bash
flutter run --dart-define-from-file=dev.env
```

For other environments:

```bash
flutter run --dart-define-from-file=staging.env
flutter run --dart-define-from-file=prod.env
```

---

## 5. Reset environment

To reset the Flutter environment (local cache/build):

```bash
make fclean
```

This command runs:

1. `flutter clean`
2. `flutter pub get`

If you need a full reset with a clean backend state, also run:

```bash
make reset
```

---

## 6. Notes

- The app fails fast if `BASE_URL` is missing or invalid
- On a physical device, `localhost` points to the device itself, not your machine
- Keep `APP_MODE` consistent with the chosen environment file
- Before changing architecture or adding new code, read the closest `AGENT.md`
  guide for the folder being edited.

Relevant agent guides:

- [mobile/AGENT.md](../AGENT.md): top-level mobile guidance
- [mobile/lib/core/AGENT.md](../lib/core/AGENT.md): core infrastructure
- [mobile/lib/data/AGENT.md](../lib/data/AGENT.md): data layer
- [mobile/lib/data/repositories/AGENT.md](../lib/data/repositories/AGENT.md): repositories
- [mobile/lib/data/services/apis/AGENT.md](../lib/data/services/apis/AGENT.md): API services and DTOs
- [mobile/lib/domain/AGENT.md](../lib/domain/AGENT.md): domain models and enums
- [mobile/lib/domain/usecases/AGENT.md](../lib/domain/usecases/AGENT.md): use cases for complex workflows
- [mobile/lib/ui/AGENT.md](../lib/ui/AGENT.md): UI layer
- [mobile/lib/ui/pages/AGENT.md](../lib/ui/pages/AGENT.md): pages and view models
- [mobile/lib/ui/components/AGENT.md](../lib/ui/components/AGENT.md): shared UI components

---

## 7. Architecture Notes

Simple user flows currently follow:

```text
UI -> ViewModel -> Repository -> API/Service -> RestClient -> Dio
```

When a view model would need to coordinate multiple repositories or a reusable
multi-step workflow, add a use case under `mobile/lib/domain/usecases` and use:

```text
UI -> ViewModel -> UseCase -> Repository -> API/Service -> RestClient -> Dio
```

Keep pages focused on Flutter concerns, view models focused on commands and
presentation orchestration, repositories focused on app data operations, and API
services focused on endpoint calls and DTO parsing.

---

## 8. Troubleshooting

### Startup failure due to missing variable

Expected error:

```text
BASE_URL not defined.
```

Or:

```text
BASE_URL is invalid: ...
```

Check whether the app was started with `--dart-define-from-file` and if the file contains a valid `BASE_URL`.

---

### App cannot connect to API

1. Confirm the API is running (`make run`)
2. Confirm `BASE_URL` is reachable from the device/emulator
3. On Android Emulator, prefer `10.0.2.2` to reach the local host machine

---

### Inconsistent Flutter dependencies

```bash
make fclean
```

If needed, run again:

```bash
cd mobile
flutter pub get
flutter run --dart-define-from-file=dev.env
```
