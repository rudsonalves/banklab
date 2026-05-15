# Operational Account Endpoints Backlog

## 1. Objective

Move account provisioning and direct ledger mutation capabilities out of the
customer-facing account API surface.

The current API exposes three authenticated endpoints that should not be treated
as mobile or customer-facing web capabilities in the target product:

- `POST /accounts`
- `POST /terminal/accounts/{id}/deposit`
- `POST /terminal/accounts/{id}/withdraw`

These endpoints are useful during development and early validation, but they
represent operational actions. Account creation is a provisioning/onboarding
action, while deposit and withdraw directly mutate the ledger without modeling a
real cash-in or cash-out flow.

## 2. Problem

The mobile app and future customer-facing web app should not be able to create
new accounts or inject/remove balance directly through generic account routes.

Leaving these capabilities in the same surface as customer account listing,
balance, statement, and internal transfers creates unclear product boundaries:

- customers could appear to have account-provisioning capabilities;
- direct deposit/withdraw endpoints could be confused with real money movement;
- operational actions are protected only by general JWT authentication and
  account ownership checks;
- onboarding approval already provisions accounts, making `POST /accounts`
  redundant for the customer-facing flow.

## 3. Fixed Premises

This backlog should not reopen decisions that already belong to the current
architecture.

Fixed premises:

- the mobile app and future customer-facing web app must not call these three
  endpoints as product capabilities;
- account creation is tied to onboarding/provisioning, not to customer self-service;
- direct deposit and withdraw are operational ledger actions, not real
  cash-in/cash-out product flows;
- the existing ledger model should remain unchanged;
- internal transfer, balance, statement, and receipt behavior are outside this
  discussion;
- implementation tasks should be created only after the exact endpoint
  positioning is clear.

## 4. Scope

Included:

- examine the three endpoints and define their correct API positioning;
- clarify how each endpoint should be represented outside the customer-facing
  API surface;
- identify what must change in routing, authorization, and documentation once the
  positioning is accepted;
- preserve development or sandbox usability where needed;
- prepare a future implementation backlog after the positioning is defined.

Not included:

- implementation tasks;
- route changes;
- handler or use-case changes;
- test implementation;
- implementing real cash-in/cash-out integrations;
- changing the ledger model;
- changing internal transfer behavior;
- changing balance, statement, or receipt endpoints;
- building admin/backoffice UI screens;
- redesigning onboarding data collection.

## 5. Current API Surface

Current routes:

- `POST /accounts`
  - creates an account for the authenticated user's `customer_id`;
  - accepts an empty body;
  - requires JWT, but is not admin-only.

- `POST /terminal/accounts/{id}/deposit`
  - directly increases account balance;
  - requires JWT and account access;
  - does not represent a real external cash-in flow.

- `POST /terminal/accounts/{id}/withdraw`
  - directly decreases account balance;
  - requires JWT and account access;
  - does not represent a real external cash-out flow.

## 6. Target Direction

Account creation should be positioned as an administrative provisioning endpoint.
It should not remain in the same customer-facing account surface as account
listing, balance, statement, and transfers.

Deposit and withdraw should be positioned as protected operational endpoints.
They may remain useful for development, tests, demos, or controlled operational
flows, but they should not be exposed as customer product actions.

The current customer-facing account API should remain focused on:

- listing the authenticated customer's accounts;
- reading balances;
- reading statements;
- looking up transfer recipients;
- executing internal transfers;
- reading transfer receipts.

## 7. Endpoint Review: POST /accounts

### Current Behavior

- Creates an account for the authenticated user's `customer_id`.
- Accepts an empty body.
- Requires JWT, but is not admin-only.
- Overlaps with the existing onboarding approval flow that already provisions an
  account.

### Target Position

The customer-facing `POST /accounts` endpoint should be removed. Creating an
account from the customer account API is considered an architectural mistake
because account creation is a provisioning action, not a customer self-service
operation.

This capability is not disabled for future use; it is repositioned. After the
change, no route, documentation, or client contract should point to the old
customer-facing behavior.

The first account continues to be created automatically during onboarding
approval. In this application, an approved customer without an account is not a
valid business state.

Additional accounts for the same `customer_id` are allowed, but only through an
admin-controlled flow. The target route is:

```http
POST /admin/customers/{customer_id}/accounts
```

The request body should remain empty for now. No account type, currency, product,
or future optional field should be reserved before there is a concrete product
need. There is no account limit for now.

The response can use the existing `AccountData` shape.

### Expected Outcome

- Mobile/customer-facing clients cannot create accounts directly.
- Account creation is represented as provisioning/admin behavior.
- The onboarding approval flow creates the first account automatically.
- Admin users can create additional accounts for an existing customer through
  `POST /admin/customers/{customer_id}/accounts`.
- The old `POST /accounts` route is removed instead of deprecated.

## 8. Endpoint Review: POST /terminal/accounts/{id}/deposit

### Current Behavior

- Directly increases account balance.
- Requires JWT and account access.
- Does not represent a real external cash-in flow.

### Target Position

Direct deposit is conceptually a banking terminal operation, such as an ATM,
teller terminal, or equivalent controlled cash desk channel. Implementing a real
terminal channel is outside the current project scope.

The terminal HTTP endpoint should be disabled. It should not be moved to
`/admin`, `/dev`, `/test-support`, or any other temporary HTTP surface, because
that would keep a direct balance-mutation capability reachable through the API
router.

For this backlog, disabled means the route registration should be commented out
in the API wiring. The handler, use case, domain behavior, and tests may remain
in place.

The existing application/domain code may remain available for tests, fixtures,
local development support, or future terminal work, but it should not be exposed
as a callable HTTP endpoint.

### Expected Outcome

- Mobile/customer-facing clients cannot directly inject account balance.
- No HTTP route exposes direct deposit in the current project scope.
- Existing ledger invariants remain unchanged.
- Future real cash-in remains a separate product/API concern.

## 9. Endpoint Review: POST /terminal/accounts/{id}/withdraw

### Current Behavior

- Directly decreases account balance.
- Requires JWT and account access.
- Does not represent a real external cash-out flow.

### Target Position

Direct withdraw is conceptually a banking terminal operation, such as an ATM,
teller terminal, or equivalent controlled cash desk channel. Implementing a real
terminal channel is outside the current project scope.

The terminal HTTP endpoint should be disabled. It should not be moved to
`/admin`, `/dev`, `/test-support`, or any other temporary HTTP surface, because
that would keep a direct balance-mutation capability reachable through the API
router.

For this backlog, disabled means the route registration should be commented out
in the API wiring. The handler, use case, domain behavior, and tests may remain
in place.

The existing application/domain code may remain available for tests, fixtures,
local development support, or future terminal work, but it should not be exposed
as a callable HTTP endpoint.

### Expected Outcome

- Mobile/customer-facing clients cannot directly remove account balance.
- No HTTP route exposes direct withdraw in the current project scope.
- Existing ledger invariants remain unchanged.
- Future real cash-out remains a separate product/API concern.

## 10. Implementation Decisions

The endpoint positioning is closed enough to generate implementation tasks:

- `POST /accounts` is repositioned to
  `POST /admin/customers/{customer_id}/accounts`.
- The old customer-facing `POST /accounts` route should not remain registered.
- `POST /terminal/accounts/{id}/deposit` should be disabled by commenting out the route
  registration in API wiring.
- `POST /terminal/accounts/{id}/withdraw` should be disabled by commenting out the route
  registration in API wiring.
- No `/admin`, `/dev`, `/test-support`, or temporary HTTP replacement should be
  created for direct deposit or withdraw.
- Route boundary tests should prove customer-facing clients cannot access the old
  customer account creation behavior and cannot reach direct deposit/withdraw
  through HTTP.

## 11. Documentation Questions

The documentation should make the endpoint positioning explicit:

- `POST /accounts` should not be described as a mobile/customer account opening
  flow.
- `deposit` should not be described as customer cash-in.
- `withdraw` should not be described as customer cash-out.
- customer-facing docs should focus on account reading and internal transfer
  flows.
- operational docs should explicitly say when a route is for development,
  sandbox, admin, or controlled operational use.

## 12. Implementation Tasks

Implementation tasks are tracked in:

- `003 - operational_account_endpoints_tasks.md`

## 13. Definition Of Done For This Backlog

- The three endpoints have been reviewed individually.
- Each endpoint has a clear target position outside the customer-facing API
  surface.
- Compatibility behavior for the current routes is known.
- Implementation tasks can be executed without guessing route placement, access
  rules, or documentation intent.
