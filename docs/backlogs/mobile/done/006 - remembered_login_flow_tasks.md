# Remembered Login Flow Tasks

These tasks split the remembered login flow into implementation-sized steps.

The goal is to let the app start from a splash decision point and route
returning users to a short login page when the last login identity is available.
The cache must store only the minimum identity data required for this UX:
name and login identifier.

Profile loading currently blocks login success. This is intentional for this
backlog and must be documented in code so it can be revisited later.

## Task 1/10: Add remembered login identity model

### Goal

Define the small model used by the app to represent the last known login
identity.

### Scope

- Add `LastLoginIdentity`.
- Include `name`.
- Include `identifier`.
- Keep the model independent from tokens, session state, and full profile data.
- Do not include customer ID, user ID, role, access token, refresh token, or
  password.

### Acceptance Criteria

- `LastLoginIdentity` has only `name` and `identifier`.
- `name` is the display name used by the short login screen.
- `identifier` is the value used to authenticate through the existing
  `LoginRequestDto.email` field for now.
- The model does not represent an authenticated user.

### Depends On

- None.

## Task 2/10: Add last login cache service

### Goal

Centralize remembered login storage access behind a dedicated service.

### Scope

- Add `LastLoginCacheService`.
- Use `LocalSecureStorage`.
- Read and write `StorageKeys.lastLoginName`.
- Read and write `StorageKeys.lastLoginIdentifier`.
- Add `get()`.
- Add `save(LastLoginIdentity identity)`.
- Add `clear()`.
- Return a failure/no identity result from `get()` when either stored value is
  missing or blank.
- Do not use `FlutterSecureStorage` directly.

### Acceptance Criteria

- The service reads storage only through `LocalSecureStorage`.
- `get()` returns a valid `LastLoginIdentity` only when both fields exist and
  are not blank.
- `save()` persists both name and identifier.
- `clear()` removes only the remembered login keys.
- The service does not read or write auth tokens.

### Depends On

- Task 1.

## Task 3/10: Register last login cache service in DI

### Goal

Make `LastLoginCacheService` available to repositories and view models through
the existing dependency injection setup.

### Scope

- Register `LastLoginCacheService`.
- Inject `LocalSecureStorage` into the service.
- Keep the registration consistent with existing core/data service patterns.
- Do not instantiate the service manually in pages or view models.

### Acceptance Criteria

- The app dependency setup can resolve `LastLoginCacheService`.
- The service receives the existing `LocalSecureStorage` instance.
- No UI class constructs `LastLoginCacheService` directly.

### Depends On

- Task 2.

## Task 4/10: Load profile and update remembered identity after login

### Goal

Ensure every successful login loads the profile and updates the remembered login
cache before the app navigates forward.

### Scope

- Update `AuthRepositoryImpl`.
- Inject `LastLoginCacheService`.
- Keep the existing `_api.login()` behavior.
- Keep storing access and refresh tokens after login succeeds.
- Call `_api.getProfile()` after tokens are stored.
- If profile loading fails, return login failure for now.
- Store the loaded profile in `_userProfile`.
- Save remembered identity using `LastLoginCacheService.save()`.
- Add a TODO documenting that profile failure currently blocks login and should
  be reviewed later.

### Acceptance Criteria

- `AuthRepositoryImpl.login()` loads profile after successful authentication.
- `_userProfile` is populated after login succeeds.
- Remembered login cache is updated after profile loading succeeds.
- Login returns failure when profile loading fails.
- The failure behavior is documented with a short TODO.
- Existing token persistence remains in place.

### Depends On

- Task 3.

## Task 5/10: Add splash route and page shell

### Goal

Introduce a splash screen as the initial route and authentication entry point.

### Scope

- Add `SplashPage`.
- Add a route for `GeneralRoutes.splash`.
- Change router `initialLocation` to `GeneralRoutes.splash.path`.
- Keep the splash page focused on startup routing, not authentication itself.
- Show a simple loading state while the decision is being resolved.

### Acceptance Criteria

- The app starts at `/splash`.
- The splash route is registered in the router.
- The splash page can render a loading state.
- The splash page does not perform login or profile loading.

### Depends On

- Task 3.

## Task 6/10: Add splash view model and decision logic

### Goal

Route users from splash to the correct login entry screen based on remembered
identity availability.

### Scope

- Add `SplashViewModel`.
- Inject `LastLoginCacheService`.
- Load remembered identity on startup.
- Expose a loading state.
- Expose a decision state for full login.
- Expose a decision state for short login.
- Navigate to full login when no remembered identity exists.
- Navigate to short login when remembered identity exists.
- Register `SplashViewModel` in DI.

### Acceptance Criteria

- Splash can resolve whether a remembered identity exists.
- Missing or incomplete remembered identity routes to `LoginPage`.
- Valid remembered identity routes to `ShortLoginPage`.
- Decision logic lives in the view model, not directly inside storage code.
- `SplashViewModel` is resolved through DI.

### Depends On

- Task 5.

## Task 7/10: Add short login route and view model

### Goal

Prepare the short login flow for returning users.

### Scope

- Add an auth route for `ShortLoginPage`, such as `/login/short`.
- Add `ShortLoginViewModel`.
- Inject `AuthRepository`.
- Inject `LastLoginCacheService`.
- Load the remembered identity.
- Expose loading, ready, missing-identity, running, success, and failure states
  as needed by the UI.
- Submit login using the cached identifier and typed password:

```dart
LoginRequestDto(
  email: identity.identifier,
  password: password,
)
```

- Register `ShortLoginViewModel` in DI.

### Acceptance Criteria

- The short login route is registered.
- `ShortLoginViewModel` can load the remembered identity.
- Missing identity can be handled by routing back to full login.
- Login submission uses `identity.identifier` through the existing
  `LoginRequestDto.email` field.
- Successful login uses the same repository flow as full login.
- The view model is resolved through DI.

### Depends On

- Task 4.
- Task 6.

## Task 8/10: Build short login page UI

### Goal

Provide the returning-user login screen that asks only for the password.

### Scope

- Add `ShortLoginPage`.
- Display a welcome message using the remembered name.
- Add a password field.
- Add password visibility toggle.
- Add submit button.
- Add loading state while login is running.
- Show login errors through the existing app snackbar/message pattern.
- Navigate to `HomeRoutes.home` after successful login.
- Add a "Use another account" link that navigates to the full login page.
- Do not add "Forget this user" yet.

### Acceptance Criteria

- The page shows the remembered user's name.
- The page does not show an editable e-mail field.
- The page validates that password is present.
- The page submits through `ShortLoginViewModel`.
- The page navigates to home on success.
- The page can route to the full login screen.
- The page does not clear the remembered identity when using another account.

### Depends On

- Task 7.

## Task 9/10: Keep full login compatible with remembered login

### Goal

Preserve the existing full login page behavior while allowing it to populate
the remembered login cache through the repository.

### Scope

- Keep `LoginPage` navigation to `HomeRoutes.home` after login success.
- Do not make `LoginPage` write remembered login storage directly.
- Keep using the existing `LoginRequestDto.email` field.
- Keep the current e-mail validator for now.
- Ensure full login success goes through the updated `AuthRepositoryImpl.login()`
  flow.

### Acceptance Criteria

- Full login still works as before from the user's perspective.
- Full login updates remembered login cache indirectly through the repository.
- Full login does not depend on `ShortLoginPage`.
- Full login does not manually read or write `StorageKeys.lastLoginName` or
  `StorageKeys.lastLoginIdentifier`.

### Depends On

- Task 4.
- Task 8.

## Task 10/10: Validate remembered login journeys

### Goal

Verify the end-to-end behavior of the remembered login flow.

### Scope

- Validate first app open without cache.
- Validate full login success.
- Validate remembered cache creation.
- Validate next app open with cache.
- Validate short login success.
- Validate remembered cache refresh after short login.
- Validate profile failure after login returns failure.
- Validate "Use another account" navigation.
- Run focused mobile checks available for the project.

### Acceptance Criteria

- Without remembered identity:

```text
Splash -> Login
```

- After full login succeeds:

```text
Login -> login API -> profile API -> save cache -> Home
```

- With remembered identity:

```text
Splash -> ShortLogin
```

- After short login succeeds:

```text
ShortLogin -> login API -> profile API -> update cache -> Home
```

- When profile loading fails after login:

```text
Login fails and shows an error
```

- "Use another account" routes from short login to full login.

### Depends On

- Task 9.
