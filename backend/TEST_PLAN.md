# Backend Test Plan for SonarQube Coverage

## Target: ≥80% Line Coverage, ≥70% Branch Coverage

---

## 1. Test Infrastructure Setup

### Dependencies to Add
```bash
cd backend
go get -tool github.com/stretchr/testify@latest
go get -tool github.com/DATA-DOG/go-sqlmock@latest
go get -tool github.com/gin-gonic/gin@latest
go get -tool go.uber.org/mock@latest
go get -tool github.com/stretchr/testify/suite@latest
go get -tool gopkg.in/yaml.v3@latest
```

### Test Structure
```
backend/
├── config/
│   └── config_test.go          # Config loading tests
├── auth/
│   └── jwt_test.go             # JWT token tests
├── model/
│   └── *_test.go               # Model validation/constructors
├── repository/
│   ├── *_test.go               # Repository tests (with sqlmock)
│   └── testutil/
│       ├── db.go               # Test DB helpers
│       └── fixtures.go         # Test data builders
├── service/
│   └── *_test.go               # Service tests (mocked repos)
├── handler/
│   └── *_test.go               # HTTP handler tests (httptest)
└── integration/
    └── api_test.go             # Full integration tests
```

### Go Test Config (`.github/workflows/test.yml` snippet)
```yaml
- name: Run backend tests
  run: |
    cd backend
    go test -v -race -coverprofile=coverage.out ./...
- name: Upload coverage to SonarQube
  uses: SonarSource/sonarqube-scan-action@v5
  env:
    SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

---

## 2. Test Categories & Coverage Targets

| Layer | Files | Coverage Target | Priority |
|-------|-------|----------------|----------|
| **Models** | 9 files | 90%+ | P0 - Constructors, validation |
| **Auth/JWT** | 1 file | 90%+ | P0 - Token gen/parse/refresh |
| **Config** | 2 files | 80%+ | P1 - Env loading, DB connect |
| **Repositories** | 6 files | 85%+ | P0 - SQL, edge cases |
| **Services** | 6 files | 85%+ | P1 - Business logic, errors |
| **Handlers** | 6 files | 80%+ | P1 - HTTP, binding, status codes |
| **Integration** | 1 suite | 70%+ | P2 - Full flows |

---

## 3. Model Tests (P0)

### `model/account.go`
- `NewAccount(id, username)` → sets fields correctly
- `Validate()` → valid email, password length, username format
- `HashPassword()` → bcrypt cost, verify works
- `CheckPassword()` → correct/incorrect password
- JSON tags → serialize correctly

### `model/game.go`
- `NewGame`, `NewGameDetail` → all fields set
- `Validate()` → price ≥ 0, name/description non-empty, poster URL valid
- Categories/requirements/media → append correctly
- JSON serialization → omitempty works

### `model/category.go`
- `NewCategory(id, name, isHidden)` → sets fields
- `Validate()` → name unique, non-empty, length ≤ 50

### `model/hardware.go`
- `NewHardware(data)` → all fields mapped
- Type enum → only "min"/"recommended" accepted

### `model/media.go`, `model/cart.go`, `model/wishlist.go`, `model/order.go`
- Constructors, validation, JSON tags

---

## 4. Auth/JWT Tests (P0)

### `auth/jwt.go`
| Function | Test Cases |
|----------|------------|
| `GenerateToken(userID, username, isSeller)` | Valid claims, expiry, signing method |
| `ParseToken(tokenString)` | Valid token, expired, malformed, wrong signature, empty |
| `RefreshToken(oldToken)` | Valid refresh, expired refresh, invalid token |
| `ExtractClaims(token)` | All claim fields, type assertions |

### Edge Cases
- Clock skew tolerance (1-2 min)
- Concurrent token parsing
- Key rotation (if implemented)

---

## 5. Config Tests (P1)

### `config/env.go`
- `LoadEnv(path)` → loads .env, missing file, malformed lines
- `Get(key, fallback)` → returns value/fallback
- Required vars → panic on missing (JWT_SECRET, DB_DSN)

### `config/db.go`
- `ConnectDB()` → success, invalid DSN, connection pool config
- Test with `sqlmock` — no real DB needed

---

## 6. Repository Tests (P0) — Use `go-sqlmock`

### Pattern: Each repo method = 1 test file section

```go
// repository/category_test.go
func TestCategoryRepository_GetCategories(t *testing.T) {
    db, mock, _ := sqlmock.New()
    defer db.Close()
    repo := &CategoryRepository{DB: db}

    // Success
    mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"category_id", "name", "is_hidden"}).AddRow(1, "Action", false))
    cats, err := repo.GetCategories()
    assert.NoError(t, err)
    assert.Len(t, cats, 1)

    // DB error
    mock.ExpectQuery("SELECT").WillReturnError(errors.New("connection failed"))
    _, err = repo.GetCategories()
    assert.Error(t, err)
}
```

### Per Repository Test Matrix

| Repository | Methods | Key Scenarios |
|------------|---------|---------------|
| **Category** | GetCategories, Create, Update, Delete | Duplicate name (unique constraint), not found on update/delete, hidden flag |
| **Game** | GetGames, GetGamesByName, GetGame, Create, Update, Delete | Pagination, filters (keyword, category, sellerID), joined categories/requirements/media, owner check on delete, not found/owned errors |
| **Account** | GetAccounts, Create, Delete, Login | Email unique, password hash, login success/fail, not found |
| **Cart** | GetCarts, AddCart, DeleteCart | Duplicate prevention, user isolation, cascade delete |
| **Wishlist** | GetWishlists, AddWishlist, DeleteWishlist | Duplicate, user isolation |
| **Order** | Create, GetOrders | Transaction atomicity, order items, total calc |

### Repository Test Utilities (`repository/testutil/`)
```go
// db.go - sqlmock wrapper
func NewMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
    db, mock, _ := sqlmock.New()
    return db, mock, func() { db.Close() }
}

// fixtures.go - test data builders
func CategoryFixture(id int, name string) model.Category { ... }
func GameFixture(id int, overrides ...func(*model.GameDetail)) *model.GameDetail { ... }
func AccountFixture(id int, username string) model.Account { ... }
```

---

## 7. Service Tests (P1) — Mock Repositories

### Pattern: Test service logic, not SQL

```go
// service/category_test.go
func TestCategoryService_Create(t *testing.T) {
    mockRepo := new(MockCategoryRepository)
    svc := &CategoryService{Repo: mockRepo}

    // Success
    mockRepo.On("Create", "Action").Return(model.Category{ID: 1, Name: "Action"}, nil)
    cat, err := svc.Create("Action")
    assert.NoError(t, err)
    assert.Equal(t, "Action", cat.Name)

    // Repo error
    mockRepo.On("Create", "RPG").Return(model.Category{}, errors.New("duplicate"))
    _, err = svc.Create("RPG")
    assert.Error(t, err)
}
```

### Per Service Test Matrix

| Service | Methods | Key Scenarios |
|---------|---------|---------------|
| **Category** | GetCategories, Create, Update, Delete | Passthrough + validation |
| **Game** | GetGamesCount, GetGames, GetGamesByName, GetGame, Create, Update, Delete | Pagination calc, filter delegation, owner check, error wrapping |
| **Account** | GetAccounts, Create, Delete, Login | Password hashing, JWT generation, unique email |
| **Cart** | GetCarts, AddCart, DeleteCart | User isolation, duplicate check |
| **Wishlist** | GetWishlists, AddWishlist, DeleteWishlist | User isolation, duplicate check |
| **Order** | Create | Transaction: order + items + stock decrement (if any), rollback on fail |

---

## 8. Handler Tests (P1) — Use `httptest`

### Pattern: Test HTTP layer — binding, status codes, JSON shape

```go
// handler/category_test.go
func TestCategoryHandler_Delete(t *testing.T) {
    mockSvc := new(MockCategoryService)
    h := &CategoryHandler{Service: mockSvc}

    router := gin.New()
    router.DELETE("/api/v1/categories/:id", h.Delete)

    // Success
    mockSvc.On("Delete", 1).Return(nil)
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("DELETE", "/api/v1/categories/1", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusNoContent, w.Code)

    // Not found
    mockSvc.On("Delete", 999).Return(sql.ErrNoRows)
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/api/v1/categories/999", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusNotFound, w.Code)
    assert.JSONEq(t, `{"error":"category not found"}`, w.Body.String())

    // Bad ID
    w = httptest.NewRecorder()
    req, _ = http.NewRequest("DELETE", "/api/v1/categories/abc", nil)
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

### Per Handler Test Matrix

| Handler | Endpoints | Key Scenarios |
|---------|-----------|---------------|
| **Category** | GET/POST/PUT/DELETE /categories | Bind JSON, ID validation, 404/400/201/204/500 |
| **Game** | GET/POST/PUT/DELETE /games, GET /games/:id | Query params (pagination, filters), multipart? (images), owner check |
| **Account** | GET/POST/DELETE /accounts, POST /login | Login credentials, JWT in response, password not in JSON |
| **Cart** | GET/POST/DELETE /carts/:userID/:gameID | User ID from param, game ID, 404/409 |
| **Wishlist** | GET/POST/DELETE /wishlists/:userID/:gameID | Same as cart |
| **Order** | POST /orders | Order creation, validation, 201 with order ID |

### Test Auth Middleware
```go
func TestAuthMiddleware(t *testing.T) {
    // Valid token → sets user in context
    // Expired → 401
    // Missing → 401
    // Malformed → 401
    // Wrong signature → 401
}
```

---

## 9. Integration Tests (P2) — Real DB (Testcontainers or local)

### Run with `-tags=integration` to skip in unit test runs

```go
// integration/api_test.go
func TestFullGameLifecycle(t *testing.T) {
    if testing.Short() { t.Skip("skipping integration test") }
    // Setup test DB (docker-compose or testcontainers)
    // 1. Register seller
    // 2. Login → get JWT
    // 3. Create category
    // 4. Create game with category
    // 5. Get game → verify
    // 6. Update game
    // 7. Buyer: add to cart/wishlist
    // 8. Create order
    // 9. Delete game (seller)
    // 10. Verify cascade/cleanup
}
```

---

## 10. Test Utilities & Patterns

### `repository/testutil/fixtures.go`
```go
package testutil

func Category(name string, opts ...func(*model.Category)) model.Category {
    c := model.NewCategory(0, name, false)
    for _, o := range opts { o(&c) }
    return c
}

func Game(title string, opts ...func(*model.GameDetail)) *model.GameDetail {
    g := model.NewGameDetail(0, title, "desc", 9.99, "poster.jpg", time.Now())
    for _, o := range opts { o(g) }
    return g
}

func WithID(id int) func(*model.Category) {
    return func(c *model.Category) { c.ID = id }
}

func WithHidden(hidden bool) func(*model.Category) {
    return func(c *model.Category) { c.IsHidden = hidden }
}
```

### `handler/testutil/httptest.go`
```go
func NewTestRouter(handlers ...func(*gin.Engine)) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    for _, h := range handlers { h(r) }
    return r
}

func AuthHeader(token string) map[string]string {
    return map[string]string{"Authorization": "Bearer " + token}
}
```

---

## 11. CI/CD Coverage Enforcement

```yaml
# .github/workflows/test.yml
- name: Backend tests
  run: |
    cd backend
    go test -race -coverprofile=coverage.out -covermode=atomic ./...
    go tool cover -func=coverage.out | grep total

- name: Check coverage thresholds
  run: |
    cd backend
    go test -coverprofile=coverage.out ./...
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$COVERAGE < 80" | bc -l) )); then
      echo "Coverage $COVERAGE% below 80%"
      exit 1
    fi
```

---

## 12. SonarQube Properties Update

```properties
# sonar-project.properties (update)
sonar.go.coverage.reportPaths=backend/coverage.out
sonar.go.test.reportPaths=backend/test-report.xml
sonar.coverage.exclusions=**/*_test.go,**/testutil/**,**/integration/**
```

Generate JUnit XML:
```bash
go test -v ./... 2>&1 | go-junit-report > test-report.xml
# or use: gotestsum --junitfile test-report.xml
```

---

## 13. Risk Areas (Extra Attention)

| Area | Risk | Test Focus |
|------|------|------------|
| **GameRepository.GetGame** | Complex joins, duplicate dedup | Multi-row scan, category/req/media dedup |
| **GameRepository.Create/Update** | Transaction rollback | Partial failure → full rollback |
| **Auth JWT** | Token expiry, refresh race | Clock skew, concurrent refresh |
| **Order Create** | Atomicity | Rollback on item failure |
| **Handler binding** | Malformed JSON, missing fields | 400 with field-level errors |
| **Pagination** | Off-by-one, overflow | Page 0, huge limit, empty results |

---

## 14. Estimate

| Phase | Effort |
|-------|--------|
| Infrastructure (deps, mock setup, fixtures) | 2 hr |
| Model tests (9 files) | 2 hr |
| Auth/JWT tests | 1 hr |
| Config tests | 1 hr |
| Repository tests (6 files × ~8 methods) | 6 hr |
| Service tests (6 files × ~6 methods) | 4 hr |
| Handler tests (6 files × ~5 endpoints) | 5 hr |
| Integration test suite | 3 hr |
| CI/CD + SonarQube config | 1 hr |
| **Total** | **~25 hr** |

---

## 15. Start Now - Minimal Viable Setup

```bash
cd backend

# 1. Add test deps
go get github.com/stretchr/testify
go get github.com/DATA-DOG/go-sqlmock
go get go.uber.org/mock/mockgen

# 2. Generate mocks for all interfaces
go generate ./...  # add //go:generate in service/*.go

# 3. Create first test file (e.g., model/account_test.go)
# 4. Run: go test -v ./model/...
# 5. Add to CI, then expand layer by layer
```

### Add `//go:generate` to service interfaces:
```go
// service/game.go
//go:generate mockgen -destination=../mocks/mock_game_repo.go -package=mocks software-metrics-project/service GameRepository
```

Then `go generate ./service/...` creates all mocks.