# Task: Add remembered login flow

## Problem Statement

Today, the mobile app always starts the authentication flow with the full login
screen, requiring the user to type their e-mail and password every time. This
creates unnecessary friction for returning users, especially in a banking-style
experience where the app should remember the last known user identity and ask
only for the password on future logins.

The goal is to introduce a remembered login flow. After a successful full login,
the app will load the authenticated user profile, cache only the minimum login
identity data, and use that cached identity on future app starts to decide
whether to show the full login screen or a shorter login screen.

For now, the remembered identity will store only:

```dart
StorageKeys.lastLoginName
StorageKeys.lastLoginIdentifier
```

The identifier will initially be the user e-mail from `UserProfile.email`. In
the future, this may become either e-mail or document number, depending on
backend support. No password, token, or full profile should be stored as part of
this cache.

If profile loading fails after a successful login, the login flow should fail
for now. This decision must be documented in the code because it may be
revisited later.

## Goal

Add a splash-based authentication entry flow that routes returning users to a
short login page when a remembered identity exists, while preserving the current
full login experience for first-time or uncached users.

## Non-Goals

- Do not implement document-based login yet.
- Do not rename `LoginRequestDto.email` yet.
- Do not change the backend login contract yet.
- Do not add a "forget this user" action yet.
- Do not implement logout behavior as part of this backlog.
- Do not persist passwords, access tokens, refresh tokens, or full profiles in
  the remembered login cache.

## Epic 1: Remembered Login Identity Cache

Create the foundation for storing and retrieving the last login identity.

### Scope

- Create `LastLoginIdentity`.
- Create `LastLoginCacheService`.
- Use `LocalSecureStorage`.
- Store `lastLoginName`.
- Store `lastLoginIdentifier`.
- Return no identity if either value is missing or empty.
- Register the service in dependency injection.

### Proposed Model

```dart
class LastLoginIdentity {
  final String name;
  final String identifier;
}
```

### Proposed Service API

```dart
class LastLoginCacheService {
  AsyncResult<LastLoginIdentity> get();
  AsyncResult<Unit> save(LastLoginIdentity identity);
  AsyncResult<Unit> clear();
}
```

## Epic 2: Profile Loading After Login

Update the login flow so every successful login also loads the user profile and
updates the remembered identity cache.

### Scope

- Update `AuthRepositoryImpl.login()`.
- After `_api.login()` succeeds, save auth tokens as today.
- Call `_api.getProfile()`.
- If profile loading fails, return login failure.
- Store the loaded profile in `_userProfile`.
- Save remembered identity using:

```dart
LastLoginIdentity(
  name: profile.name,
  identifier: profile.email,
)
```

- Add a TODO documenting that profile failure currently blocks login and should
  be reviewed later.

### Expected Flow

```text
_api.login()
save currentUser
save accessToken/refreshToken
load profile through _api.getProfile()
if profile fails: fail login
save _userProfile
save LastLoginIdentity
return LoggedUser
```

## Epic 3: Splash Decision Flow

Introduce a splash screen as the initial route and decision point for the
authentication entry flow.

### Scope

- Create `SplashPage`.
- Create `SplashViewModel`.
- Change the app initial route to `/splash`.
- On splash load, read `LastLoginCacheService.get()`.
- If identity exists, navigate to `ShortLoginPage`.
- If identity does not exist, navigate to the full `LoginPage`.

### Expected Flow

```text
Splash
  -> remembered identity found: ShortLoginPage
  -> remembered identity not found: LoginPage
```

## Epic 4: Short Login Flow

Create a login screen for returning users with a remembered identity.

### Scope

- Add `ShortLoginPage`.
- Add `ShortLoginViewModel`.
- Add a route such as `/login/short`.
- Display a welcome message using the cached name.
- Show only the password input.
- Submit login using:

```dart
LoginRequestDto(
  email: identity.identifier,
  password: password,
)
```

- On success, navigate to `HomeRoutes.home`.
- Add a link to use another account, navigating to the full login screen.
- Do not implement "Forget this user" yet.

## Epic 5: Full Login Compatibility

Keep the existing full login behavior while allowing it to feed the remembered
login flow.

### Scope

- Keep `LoginPage` navigation to home after success.
- Do not make `LoginPage` directly responsible for saving the remembered login
  cache.
- Keep cache saving centralized in `AuthRepositoryImpl.login()`.
- Keep using `profile.email` as the remembered identifier.
- Keep using `LoginRequestDto.email` for now.

## Epic 6: Flow Validation

Validate the expected user journeys.

### Scenarios

First app open without cache:

```text
Splash -> Login
```

Successful full login:

```text
Login -> login API -> profile API -> save cache -> Home
```

Next app open with cache:

```text
Splash -> ShortLogin
```

Successful short login:

```text
ShortLogin -> login API -> profile API -> update cache -> Home
```

Profile load failure after login:

```text
Login fails and shows an error
```

Use another account:

```text
ShortLogin -> Login
```

## Acceptance Criteria

- The app starts at the splash route.
- Splash routes to full login when no remembered identity exists.
- Splash routes to short login when both remembered identity fields exist.
- Full login loads profile after authentication.
- Full login saves remembered identity after profile loading succeeds.
- Short login authenticates using the cached identifier and typed password.
- Short login loads profile and refreshes the remembered identity after success.
- Profile loading failure blocks login for now and returns an error.
- The remembered login cache stores only name and identifier.
- Passwords are never persisted in the remembered login cache.
