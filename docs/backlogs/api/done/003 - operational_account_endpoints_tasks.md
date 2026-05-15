# Operational Account Endpoints Tasks

Parent backlog:

- `003 - operational_account_endpoints_backlog.md`

Suggested project fields for all tasks:

- Status: Backlog
- Priority: High
- Area: API
- Type: Architecture/Security

## Task 1/7: Reposition account creation under the admin API

Status: Done

### Goal

Move account creation out of the customer-facing account API and expose it as an
admin-controlled provisioning action.

### Scope

- Add the admin route:
  `POST /admin/customers/{customer_id}/accounts`.
- Reuse the existing account creation behavior where appropriate.
- Make `customer_id` come from the admin route path.
- Keep the request body empty.
- Return the existing `AccountData` response shape.
- Preserve onboarding approval as the flow that creates the first account.
- Allow additional accounts for the same `customer_id`.
- Do not add an account limit.

### Acceptance Criteria

- Admin users can create an additional account for an existing customer through
  `POST /admin/customers/{customer_id}/accounts`.
- The response uses `AccountData`.
- The request does not reserve future fields.
- Onboarding approval continues to create the first account automatically.
- A customer-facing client cannot create accounts directly.

### Depends On

- None.

## Task 2/7: Remove the customer-facing POST /accounts route

Status: Done

### Goal

Remove the old customer-facing account creation behavior from the public account
API surface.

### Scope

- Remove the `POST /accounts` route registration from the customer-facing account
  router.
- Ensure no route, documentation, or client contract points to customer
  self-service account creation.
- Keep reusable account creation code only insofar as it supports onboarding and
  the new admin route.
- Update or remove tests that expected customer-facing account creation.

### Acceptance Criteria

- `POST /accounts` is no longer registered as a customer-facing route.
- Customer account routes only expose account reading and transfer-related
  capabilities.
- Existing tests no longer describe `POST /accounts` as valid customer behavior.
- Account creation remains available through onboarding approval and the new admin
  route.

### Depends On

- Task 1.

## Task 3/7: Disable the direct deposit HTTP endpoint

Status: Done

### Goal

Disable direct deposit as a callable HTTP endpoint without discarding the internal
application/domain behavior.

### Scope

- Comment out the `POST /terminal/accounts/{id}/deposit` route registration in API wiring.
- Do not move the endpoint to `/admin`, `/dev`, `/test-support`, or any temporary
  HTTP surface.
- Keep the handler, use case, domain behavior, and tests if still useful.
- Ensure mobile/customer-facing clients cannot call direct deposit through HTTP.

### Acceptance Criteria

- `POST /terminal/accounts/{id}/deposit` is not reachable through the API router.
- No replacement HTTP route is introduced for direct deposit.
- Existing internal deposit logic remains available for tests, fixtures, local
  support, or future terminal work.
- Direct deposit is not presented as customer cash-in.

### Depends On

- None.

## Task 4/7: Disable the direct withdraw HTTP endpoint

Status: Done

### Goal

Disable direct withdraw as a callable HTTP endpoint without discarding the
internal application/domain behavior.

### Scope

- Comment out the `POST /terminal/accounts/{id}/withdraw` route registration in API wiring.
- Do not move the endpoint to `/admin`, `/dev`, `/test-support`, or any temporary
  HTTP surface.
- Keep the handler, use case, domain behavior, and tests if still useful.
- Ensure mobile/customer-facing clients cannot call direct withdraw through HTTP.

### Acceptance Criteria

- `POST /terminal/accounts/{id}/withdraw` is not reachable through the API router.
- No replacement HTTP route is introduced for direct withdraw.
- Existing internal withdraw logic remains available for tests, fixtures, local
  support, or future terminal work.
- Direct withdraw is not presented as customer cash-out.

### Depends On

- None.

## Task 5/7: Update API documentation for the new endpoint surface

Status: Done

Note: implemented API report files were intentionally not edited in this pass.

### Goal

Make the documented REST API match the new route positioning.

### Scope

- Remove customer-facing documentation for `POST /accounts`.
- Add documentation for `POST /admin/customers/{customer_id}/accounts`.
- Remove or clearly mark `POST /terminal/accounts/{id}/deposit` as disabled.
- Remove or clearly mark `POST /terminal/accounts/{id}/withdraw` as disabled.
- Update route lists in `api/README.md`.
- Keep implemented API reports unchanged in this pass.
- Make clear that direct deposit and withdraw are conceptual terminal operations,
  but no terminal HTTP API exists in the current project scope.

### Acceptance Criteria

- Documentation does not describe account creation as customer self-service.
- Documentation describes account creation as admin provisioning.
- Documentation does not expose direct deposit/withdraw as callable customer,
  admin, dev, or test-support HTTP endpoints.
- Route lists match the registered API surface.

### Depends On

- Task 1.
- Task 2.
- Task 3.
- Task 4.

## Task 6/7: Add and update route boundary tests

Status: Done

### Goal

Prove the API route boundary after account creation is repositioned and direct
ledger mutation endpoints are disabled.

### Scope

- Add tests for successful admin account creation.
- Add tests that customer-role users cannot create accounts through the admin
  endpoint.
- Update or remove tests for customer-facing `POST /accounts`.
- Add route-level coverage proving direct deposit is no longer reachable through
  HTTP.
- Add route-level coverage proving direct withdraw is no longer reachable through
  HTTP.
- Keep application/domain tests for account creation, deposit, and withdraw
  behavior where they still provide useful internal coverage.

### Acceptance Criteria

- Admin account provisioning is tested.
- Customer account self-service creation is not treated as valid behavior.
- Direct deposit and withdraw are not reachable through registered HTTP routes.
- Existing account, onboarding, and transaction behavior remains covered.

### Depends On

- Task 1.
- Task 2.
- Task 3.
- Task 4.

## Task 7/7: Run affected API tests

Status: Done

### Goal

Validate the changed API surface without regressing account, onboarding, and
transaction behavior.

### Scope

- Run account/bankaccount tests.
- Run admin/onboarding approval tests.
- Run transaction tests related to deposit and withdraw internal behavior.
- Run any handler or route tests affected by the route changes.
- Fix regressions introduced by the implementation.

### Acceptance Criteria

- Affected Go tests pass.
- Onboarding approval still provisions the first account.
- Admin account provisioning works for additional accounts.
- Direct deposit and withdraw are disabled at the HTTP route level.
- Internal ledger behavior remains covered by tests.

### Depends On

- Task 5.
- Task 6.

## Suggested GitHub Project Order

1. Reposition account creation under the admin API.
2. Remove the customer-facing `POST /accounts` route.
3. Disable the direct deposit HTTP endpoint.
4. Disable the direct withdraw HTTP endpoint.
5. Update API documentation for the new endpoint surface.
6. Add and update route boundary tests.
7. Run affected API tests.
