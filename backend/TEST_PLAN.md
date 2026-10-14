# Backend Test Plan

## Purpose

Build a repeatable Go test suite for the HTTP API and its business layers. The
suite should prove behavior at the model, authentication, repository, service,
and handler boundaries without requiring a running MySQL server for unit tests.

## Current test stack

- Go module: `software-metrics-project`
- Test runner: `go test`
- Assertions: `testify`
- SQL isolation: `go-sqlmock`
- HTTP isolation: `gin` with `httptest`
- Coverage output: Go cover profile
- Existing test helpers: `repository/repotestutil` and `handler/handlertestutil`

Do not add a mocking framework or a second assertion library. Reuse the
helpers and patterns already present in the repository.

## Commands

Run from `backend`:

```powershell
go test ./...
go test -race ./...
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

The normal PR check is:

```powershell
go test -race -coverprofile=coverage.out -covermode=atomic ./...
```

`coverage.out` and `coverage.html` are generated artifacts and must not be
committed.

## Coverage goals

Coverage is a signal, not a reason to test implementation details.

| Area | Initial goal | Priority |
| --- | ---: | --- |
| Models and value objects | 90% | P0 |
| JWT/authentication | 90% | P0 |
| Repositories | 85% | P0 |
| Services | 85% | P1 |
| HTTP handlers and middleware | 80% | P0 |
| Overall package coverage | 80% | P0 |

Raise a threshold only after the current suite meets it consistently. Exclude
generated files and test helpers from SonarQube rather than writing artificial
tests for them.

## Test layers

### 1. Models

Cover constructors, conversion methods, validation, enum/value-object
behavior, and JSON serialization for the model packages that exist in the
repository.

Required cases:

- valid construction preserves every supplied field;
- zero values and invalid values are rejected or represented as documented;
- hidden flags, IDs, names, prices, and dates round-trip correctly;
- password values are never serialized into API response objects;
- invalid input does not panic.

### 2. Authentication

Cover `auth` with table-driven tests:

- token creation with a configured `JWT_SECRET`;
- missing or empty secret returns an error;
- valid token is accepted;
- malformed, expired, and incorrectly signed tokens are rejected;
- missing, invalid, and non-positive user IDs are rejected where applicable;
- middleware returns the expected `401` response and does not invoke the
  protected handler when authentication fails;
- valid middleware requests receive the authenticated user in Gin context.

Tests must set and restore `JWT_SECRET` per test. Do not use a real production
secret.

### 3. Repositories

Use `go-sqlmock` and assert both returned values and expected SQL arguments.
Every repository method needs at least:

1. a successful result;
2. a database error;
3. an empty or not-found result where the method can produce one.

For write methods also cover:

- affected-row count of zero;
- transaction rollback after a partial failure;
- duplicate/constraint errors;
- correct argument order and boolean values.

For joined game queries, include multiple rows for one game and verify related
categories, requirements, and media are not duplicated.

Always call `mock.ExpectationsWereMet()` through the existing helper or test
cleanup so a test cannot pass while silently skipping an expected query.

### 4. Services

Exercise service behavior through the real repository or a small local fake
when the service depends on an interface. Focus on business behavior:

- successful delegation and result mapping;
- validation before persistence;
- repository errors are returned, not swallowed;
- not-found and ownership rules;
- transaction-related errors remain visible to callers;
- authentication services never return a password.

Do not duplicate SQL assertions in service tests; those belong in repository
tests.

### 5. HTTP handlers

Use `gin.TestMode`, `httptest.NewRecorder`, and a fresh router per handler
scenario. Assert status, JSON body, and service/database interaction.

For each endpoint cover:

- valid success response;
- malformed JSON or invalid path/query input;
- missing required fields;
- not-found response;
- service/database failure response;
- authentication/authorization failure where protected.

Category endpoints are the reference implementation:

| Endpoint | Required cases |
| --- | --- |
| `GET /api/v1/categories` | list success, empty list, database error |
| `POST /api/v1/categories` | valid body, malformed body, validation failure, database error |
| `PUT /api/v1/categories/:id` | valid body, invalid ID/body, not found, database error |
| `DELETE /api/v1/categories/:id` | success `204`, invalid ID `400`, not found `404`, database error |

Use exact JSON assertions for error contracts. Do not assert only that a
response is “not 500”.

### 6. Integration tests

Keep database-backed tests separate from unit tests using an explicit build
tag such as `integration`. They should not run in `go test ./...`.

The first integration flow should cover:

1. register/login;
2. authenticated category creation;
3. category listing;
4. update;
5. delete;
6. verification that unauthorized requests are rejected.

Use an isolated database/schema and deterministic cleanup. Do not make unit
tests depend on Docker, MySQL, network access, or developer `.env` files.

## Test organization

Keep tests next to the package they exercise:

```text
backend/
├── auth/*_test.go
├── handler/*_test.go
├── model/*_test.go
├── repository/*_test.go
├── service/*_test.go
└── integration/*_test.go       # integration build tag only
```

Use table-driven tests for input matrices. Name tests by observable behavior,
for example `TestCategoryHandler_Delete_InvalidID`, not by implementation
details.

## Execution order

1. Run the baseline commands and record current failures.
2. Fix or replace misleading tests that assert only “not 500”.
3. Complete P0 authentication, model, repository, and handler cases.
4. Add P1 service error and validation cases.
5. Add tagged integration tests.
6. Run race and coverage commands.
7. Review the HTML report for untested branches, then add only behavior-based
   tests that cover meaningful gaps.

## Exit criteria

- `go test ./...` passes.
- `go test -race ./...` passes.
- No test requires a production secret, live database, or network.
- HTTP error responses have exact contract assertions.
- Coverage profile is generated successfully and overall coverage is at least
  80%.
- Integration tests are explicit and isolated from the unit-test command.
