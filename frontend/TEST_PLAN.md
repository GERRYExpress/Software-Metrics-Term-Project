# Frontend Test Plan

## Purpose

Create a focused Vue 3 test suite that verifies user-visible behavior,
service/API contracts, and Pinia state transitions. The plan uses the tooling
already installed in `frontend`; it does not add a second test framework or
snapshot-heavy tests.

## Current test stack

- Vue 3 with `<script setup>`
- Vitest
- `@vue/test-utils`
- `@testing-library/vue` where semantic DOM queries are useful
- jsdom
- V8 coverage via `@vitest/coverage-v8`
- Tailwind CSS 4 through `@tailwindcss/vite`

The existing test configuration includes both `src/**/*.{test,spec}.ts` and
`tests/**/*.{test,spec}.ts`, and already excludes tests from production
coverage.

## Commands

Run from `frontend`:

```powershell
npm test
npm test -- --run tests/stores/auth.test.ts
npm run type-check
npm run build
```

`npm test` runs Vitest with coverage and enforces the configured thresholds:
80% lines, 80% statements, 80% functions, and 70% branches.

Coverage output is written to `frontend/coverage/`; it is a generated artifact
and must not be committed.

## Test rules

- Test observable behavior, not Vue internals.
- Prefer real components and Pinia stores; mock only network/API boundaries.
- Use `vi.clearAllMocks()` and reset Pinia before each test.
- Assert accessible roles, labels, text, emitted events, model updates, and
  navigation effects.
- Avoid snapshots unless the rendered structure itself is an API contract.
- Do not test that Tailwind generated CSS has a particular pixel value in jsdom.

## Coverage priorities

| Area | Initial goal | Priority |
| --- | ---: | --- |
| API services | 90% | P0 |
| Pinia stores | 85% | P0 |
| Shared/base components | 85% | P0 |
| Authentication and checkout views | 75% | P1 |
| Remaining views/components | 70% | P2 |
| Overall configured coverage | 80% | P0 |

## Test layers

### 1. API services

Mock `@/service/AxiosClient` only. For every public service method test:

- successful response and returned data;
- exact HTTP method, URL, and request body/query;
- rejected network/API error is propagated;
- empty response where the UI supports it;
- `401`, `404`, and validation errors where callers branch on them.

Priority services:

- `AccountService`;
- `CategoryService`;
- `GameService`;
- `CartService`;
- `WishlistService`.

`AxiosClient` itself should be tested at the interceptor boundary: token
injection, configured base URL, and the observable behavior of rejected
responses. Do not call a real backend.

### 2. Pinia stores

Start with the existing auth store and add equivalent coverage to the other
stores as their public actions are confirmed.

Auth cases:

- initial state is unauthenticated;
- successful login sets user, token, and authentication state;
- failed login rejects and leaves state unchanged;
- successful registration sets user and token;
- `fetchProfile` updates the user;
- a `401` profile response logs out;
- logout clears all session state.

For game, seller, and user-action stores cover loading/error transitions,
successful updates, empty responses, and API failures. Mock services, not the
store itself.

### 3. Shared components

Prioritize components reused by multiple views:

| Component | Required behavior |
| --- | --- |
| `BaseButton.vue` | slot content, click forwarding, each supported `type`, attribute forwarding |
| `BaseInput.vue` | label/placeholder/type, `v-model` update, arbitrary input attributes |
| `BaseSelect.vue` | options, selected value, change/update behavior |
| `BaseToggle.vue` | checked state, toggle event, disabled behavior |
| `AppCard.vue` | slot rendering and link/click behavior |
| `AppNavItem.vue` | label/link and active state |

Use the actual prop names in each component. Tests must not assume props or
CSS classes that the component does not expose.

### 4. Tailwind styling contract

Tailwind is presentation, so tests should verify only the class contracts that
are intentionally selected by component logic:

- a `BaseButton` type adds the expected semantic class (`bg-action`,
  `bg-red`, `bg-secondary`, or danger-link classes);
- full-width/shared layout classes are present when required;
- disabled state is represented by the native `disabled` attribute or the
  component's documented accessible state;
- conditional classes change when the relevant prop changes.

Do not test Tailwind's compiled CSS, generated utility order, responsive
breakpoints, or visual spacing in unit tests. Use browser/manual checks for
those concerns if a UI change specifically affects them.

### 5. Views and user flows

Add focused tests for high-risk flows, mocking services and router navigation:

- login: validation, successful redirect, failed-login message;
- registration: required fields, password mismatch, successful submission;
- product detail: loading, loaded product, add-to-cart action;
- cart: empty state, item rendering, removal, checkout action;
- wishlist: empty state, item removal, move-to-cart action;
- seller category/product forms: validation and save failure.

Do not attempt to fully test every visual branch of every page before P0
services, stores, and shared components are reliable.

## Test organization

```text
frontend/
├── tests/
│   ├── components/
│   ├── services/
│   ├── stores/
│   └── views/
└── src/
    └── **/*.{test,spec}.ts
```

Keep tests near the existing `tests/` structure. Name tests by behavior, for
example `renders the danger button styling`, not `test danger class`.

## Execution order

1. Run the focused baseline and full coverage command from `frontend`.
2. Correct tests that reference non-existent props/classes or assert only
   “not 500”.
3. Complete P0 service and auth-store cases.
4. Complete shared component and Tailwind class-contract cases.
5. Add P1 authentication, cart, and seller flows.
6. Run type-check, production build, and full coverage.
7. Review uncovered branches and add only tests tied to user-visible behavior.

## Exit criteria

- `npm test` passes with configured coverage thresholds.
- Focused tests can be run independently.
- `npm run type-check` passes.
- `npm run build` passes.
- No test calls a real API or depends on a developer's local storage state.
- Tailwind assertions check intentional class selection, not generated CSS.
- Coverage and build artifacts remain untracked.
