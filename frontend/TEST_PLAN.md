# Frontend Test Plan for SonarQube Coverage

## Target: ≥90% Line Coverage, ≥70% Branch Coverage

---

## 1. Test Infrastructure Setup

### Dependencies to Add
```json
{
  "devDependencies": {
    "vitest": "^2.0.0",
    "@vue/test-utils": "^2.4.0",
    "jsdom": "^24.0.0",
    "@vitest/coverage-v8": "^2.0.0",
    "@testing-library/vue": "^8.1.0",
    "vitest-mock-extended": "^2.0.0"
  }
}
```

### Scripts to Add
```json
{
  "scripts": {
    "test": "vitest run --coverage",
    "test:watch": "vitest",
    "test:ui": "vitest --ui"
  }
}
```

### Vitest Config (`vitest.config.ts`)
```typescript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    coverage: {
      provider: 'v8',
      reporter: ['text', 'json', 'html', 'lcov'],
      reportsDirectory: './coverage',
      exclude: ['node_modules/**', 'dist/**', '**/*.d.ts', '**/*.config.*', '**/main.ts'],
      thresholds: {
        lines: 80,
        functions: 80,
        branches: 70,
        statements: 80
      }
    },
    include: ['src/**/*.{test,spec}.{ts,tsx}'],
    setupFiles: ['./tests/setup.ts']
  }
})
```

---

## 2. Test Categories & Coverage Targets

| Layer | Files | Coverage Target | Priority |
|-------|-------|----------------|----------|
| **Services** | 6 files | 90%+ | P0 - Critical |
| **Stores (Pinia)** | 5 files | 85%+ | P0 - Critical |
| **Components** | 21 files | 75%+ | P1 - High |
| **Views/Pages** | 14 files | 70%+ | P2 - Medium |
| **Router/Utils** | 3 files | 70%+ | P2 - Medium |

---

## 3. Service Layer Tests (P0)

### `AccountService.ts`
- `login(email, password)` → success, invalid creds, network error
- `register(data)` → success, validation error, duplicate email
- `logout()` → clears token, redirects
- `getProfile()` → success, 401 handling
- `updateProfile(data)` → success, validation error
- `changePassword(old, new)` → success, wrong password

### `GameService.ts`
- `getGames(params)` → pagination, filters, empty response
- `getGame(id)` → success, 404, network error
- `createGame(data)` → success, validation error, auth error
- `updateGame(id, data)` → success, 404, auth error
- `deleteGame(id)` → success, 404, auth error

### `CategoryService.ts`
- `getCategories()` → success, empty, network error
- `createCategory(name, isHidden)` → success, duplicate, validation
- `updateCategory(id, name, isHidden)` → success, 404, validation
- `deleteCategory(id)` → success, 404, network error

### `CartService.ts`
- `getCarts(userId)` → success, empty, unauthorized
- `addCart(userId, gameId)` → success, duplicate, unauthorized
- `deleteCart(userId, gameId)` → success, not found

### `WishlistService.ts`
- `getWishlists(userId)` → success, empty
- `addWishlist(userId, gameId)` → success, duplicate
- `deleteWishlist(userId, gameId)` → success, not found

### `AxiosClient.ts`
- Interceptors: token injection, 401 refresh, error normalization
- Base URL config, timeout handling

---

## 4. Store Tests (P0)

### `auth.ts`
- State: `user`, `token`, `isAuthenticated`
- Actions: `login`, `register`, `logout`, `fetchProfile`, `refreshToken`
- Getters: `isLoggedIn`, `currentUser`
- Edge: token expiry, concurrent login, failed refresh

### `game.ts`
- State: `games`, `currentGame`, `loading`, `error`, `pagination`
- Actions: `fetchGames`, `fetchGame`, `createGame`, `updateGame`, `deleteGame`
- Getters: `gameList`, `gameById`, `isLoading`
- Edge: optimistic updates, error rollback, cache invalidation

### `seller.ts`
- State: `products`, `sales`, `stats`, `loading`
- Actions: `fetchSellerGames`, `fetchSales`, `fetchStats`
- Edge: empty seller, permission denied

### `userAction.ts`
- State: `cartCount`, `wishlistCount`
- Actions: `incrementCart`, `decrementCart`, `toggleWishlist`
- Getters: `totalCartItems`, `isInWishlist`

### `counter.ts` (simple)
- State: `count`
- Actions: `increment`, `decrement`, `reset`

---

## 5. Component Tests (P1)

### Base Components (High Reuse → High Priority)
| Component | Key Behaviors |
|-----------|--------------|
| `BaseButton.vue` | click emit, disabled, loading, variants, sizes |
| `BaseInput.vue` | v-model, validation, error display, types |
| `BaseSelect.vue` | options, v-model, disabled, placeholder |
| `BaseToggle.vue` | checked/unchecked, disabled, label |
| `AppCard.vue` | slots, hover, click, image |
| `AppNavbar.vue` | links, active state, mobile toggle |
| `AppFooter.vue` | links, copyright |
| `AppNavItem.vue` | active, icon, badge |

### Business Components
| Component | Key Behaviors |
|-----------|--------------|
| `GameCard.vue` | image, title, price, category, add to cart/wishlist |
| `ProductListCard.vue` | compact variant, seller badge |
| `ProductMainContainer.vue` | gallery, info, price, actions |
| `CategoryTable.vue` | list, create, edit, delete, pagination |
| `MediaCarousel.vue` | slides, navigation, autoplay, thumbnails |
| `SystemRequirement.vue` | min/rec tabs, formatted display |
| `GalleryUploadCard.vue` | drag-drop, preview, remove, validation |
| `PriceSummaryCard.vue` | subtotal, discount, total, promo code |
| `PaymentModal.vue` | method select, form validation, submit |
| `PaymentSuccessModal.vue` | order details, close, redirect |
| `NotificationToast.vue` | types (success/error/warning), auto-dismiss |
| `WelcomeModal.vue` | show once, dismiss, CTA |
| `SessionExpiredModal.vue` | forced logout, re-login |
| `SellerNavbar.vue` | seller links, profile dropdown |

---

## 6. View Tests (P2)

| View | Critical Paths |
|------|---------------|
| `LandingView.vue` | hero, featured games, categories, CTA |
| `LoginFormView.vue` | form validation, submit, error display, redirect |
| `RegisterFormView.vue` | validation, password match, submit, redirect |
| `ProductDetailView.vue` | load game, gallery, requirements, add to cart |
| `SearchResultView.vue` | filters, pagination, empty state, sort |
| `CartView.vue` | list items, quantity, remove, checkout |
| `ProfileView.vue` | tabs, edit profile, change password, orders |
| `SellerDashboardView.vue` | stats, charts, quick actions |
| `SellerProductFormView.vue` | create/edit, validation, images, requirements |
| `SellerProfileView.vue` | seller info, verification status |
| `WishlistView.vue` | list, remove, move to cart |
| `ContactView.vue` | form validation, submit |
| `IdentificationView.vue` | document upload, verification flow |

---

## 7. Test Patterns & Conventions

### Service Test Template
```typescript
// tests/services/CategoryService.test.ts
import { describe, it, expect, vi, beforeEach } from 'vitest'
import CategoryService from '@/service/CategoryService'
import apiClient from '@/service/AxiosClient'

vi.mock('@/service/AxiosClient')

describe('CategoryService', () => {
  beforeEach(() => vi.clearAllMocks())

  describe('getCategories', () => {
    it('returns categories on success', async () => {
      const mockData = [{ id: 1, name: 'Action', isHidden: false }]
      vi.mocked(apiClient.get).mockResolvedValue({ data: mockData })
      const result = await CategoryService.getCategories()
      expect(result).toEqual(mockData)
      expect(apiClient.get).toHaveBeenCalledWith('/api/v1/categories')
    })

    it('throws on network error', async () => {
      vi.mocked(apiClient.get).mockRejectedValue(new Error('Network'))
      await expect(CategoryService.getCategories()).rejects.toThrow('Network')
    })
  })

  describe('deleteCategory', () => {
    it('calls delete endpoint with id', async () => {
      vi.mocked(apiClient.delete).mockResolvedValue({})
      await CategoryService.deleteCategory(5)
      expect(apiClient.delete).toHaveBeenCalledWith('/api/v1/categories/5')
    })

    it('propagates 404 error', async () => {
      const err = { response: { status: 404 } }
      vi.mocked(apiClient.delete).mockRejectedValue(err)
      await expect(CategoryService.deleteCategory(999)).rejects.toMatchObject(err)
    })
  })
})
```

### Store Test Template
```typescript
// tests/stores/auth.test.ts
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '@/stores/auth'
import AccountService from '@/service/AccountService'

vi.mock('@/service/AccountService')

describe('auth store', () => {
  beforeEach(() => setActivePinia(createPinia()))

  it('login sets user and token', async () => {
    const store = useAuthStore()
    const mockUser = { id: 1, username: 'test' }
    const mockToken = 'jwt-token'
    vi.mocked(AccountService.login).mockResolvedValue({ user: mockUser, token: mockToken })

    await store.login('test@test.com', 'password')

    expect(store.user).toEqual(mockUser)
    expect(store.token).toBe(mockToken)
    expect(store.isAuthenticated).toBe(true)
  })

  it('logout clears state', async () => {
    const store = useAuthStore()
    store.user = { id: 1 } as any
    store.token = 'token'
    store.logout()
    expect(store.user).toBeNull()
    expect(store.token).toBeNull()
    expect(store.isAuthenticated).toBe(false)
  })
})
```

### Component Test Template
```typescript
// tests/components/BaseButton.test.ts
import { mount } from '@vue/test-utils'
import BaseButton from '@/components/BaseButton.vue'

describe('BaseButton', () => {
  it('emits click when clicked', () => {
    const wrapper = mount(BaseButton)
    wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('does not emit when disabled', () => {
    const wrapper = mount(BaseButton, { props: { disabled: true } })
    wrapper.trigger('click')
    expect(wrapper.emitted('click')).toBeFalsy()
  })

  it('renders loading state', () => {
    const wrapper = mount(BaseButton, { props: { loading: true } })
    expect(wrapper.find('.loading').exists()).toBe(true)
  })

  it('applies variant classes', () => {
    const wrapper = mount(BaseButton, { props: { variant: 'danger' } })
    expect(wrapper.classes()).toContain('btn-danger')
  })
})
```

---

## 8. Execution Order

1. **Setup** - Add deps, config, setup file
2. **Services** - 6 files, ~50 tests (fast, no DOM)
3. **Stores** - 5 files, ~40 tests (Pinia + mocked services)
4. **Base Components** - 8 files, ~60 tests (shallow mount)
5. **Business Components** - 13 files, ~80 tests (full mount, mocked stores)
6. **Views** - 14 files, ~40 tests (router, integration)
7. **Run Coverage** - `npm run test` → check `coverage/lcov.info`

---

## 9. CI Integration

```yaml
# .github/workflows/test.yml
- name: Run frontend tests
  run: cd frontend && npm ci && npm run test
- name: Upload coverage to SonarQube
  uses: SonarSource/sonarqube-scan-action@v5
  env:
    SONAR_TOKEN: ${{ secrets.SONAR_TOKEN }}
```

---

## 10. Risk Areas (Manual Verification Needed)

- Axios interceptor token refresh race conditions
- Pinia store persistence (localStorage)
- Router navigation guards with auth state
- File upload components (GalleryUploadCard)
- Payment flow integration (PaymentModal → backend)
- WebSocket/real-time features (if any)

---

## 11. Estimate

| Phase | Effort |
|-------|--------|
| Infrastructure setup | 1 hr |
| Service tests (6 files) | 3 hr |
| Store tests (5 files) | 2 hr |
| Component tests (21 files) | 6 hr |
| View tests (14 files) | 4 hr |
| CI config & tuning | 1 hr |
| **Total** | **~17 hr** |

---

## 12. Start Now - Minimal Viable Setup

```bash
cd frontend
npm i -D vitest @vue/test-utils jsdom @vitest/coverage-v8 @testing-library/vue
# Add vitest.config.ts, tests/setup.ts, scripts to package.json
npm run test  # verify pipeline works
# Then add tests incrementally by priority
```