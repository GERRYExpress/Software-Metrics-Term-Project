# Seller Frontend Conversion Brief

## Role

You are a senior Vue 3 and TypeScript engineer working inside the repository
`Software-Metrics-Term-Project`.

Convert the seller-side frontend from the legacy project into this repository's
existing Vue frontend. The buyer-side frontend has already been converted and
must remain working. Do not rebuild or redesign buyer pages unless a shared
component change is strictly required by the seller implementation.

## Main objective

Implement a complete seller experience that is consistent with the existing
buyer experience and the legacy seller templates:

1. Seller layout and navigation.
2. Seller dashboard.
3. Seller category management.
4. Seller product listing management.
5. Seller product add/edit form.
6. Seller gallery and media upload controls.
7. Seller profile and account actions.
8. Seller routes and role-aware navigation.

Inspect the repository before editing. Reuse existing components, services,
types, stores, CSS classes, API client configuration, and visual patterns
wherever possible. Do not introduce a second design system or duplicate
existing buyer components.

## Scope boundaries

### In scope

- Seller-only views, components, services, types, routing, and state.
- Connecting seller screens to the existing backend where endpoints exist.
- Adding the minimum backend/API support required for seller mutations if the
  current backend does not expose the required operations.
- Validation, loading states, empty states, error states, and success
  notifications for seller actions.
- Responsive behavior for desktop, tablet, and mobile layouts.
- Type-safe Vue and TypeScript implementation.

### Out of scope

- Replacing the existing buyer pages.
- Changing the buyer visual language, cart, wishlist, or checkout behavior.
- Replacing the existing authentication architecture unless it is required for
  seller routing.
- Adding unrelated features.
- Removing existing shared components merely because the legacy project used a
  different structure.

## Current project context

### Frontend stack

- Vue 3 with `<script setup lang="ts">`.
- Vite.
- TypeScript.
- Vue Router.
- Pinia.
- Tailwind CSS v4 through the Vite plugin.
- Axios through the shared API client.
- Lucide Vue icons.
- Existing reusable UI components such as `BaseButton`, `BaseInput`,
  `BaseSelect`, `AppCard`, `AppNavbar`, `AppFooter`, `NotificationToast`, and
  modal components.

### Important frontend locations

- `frontend/src/views/` - route-level pages.
- `frontend/src/components/` - reusable UI components.
- `frontend/src/service/` - API service modules.
- `frontend/src/stores/` - Pinia stores.
- `frontend/src/router/index.ts` - current route definitions.
- `frontend/src/types.ts` - shared domain types.
- `frontend/src/main.css` - global styling and Tailwind setup.
- `frontend/public/Landing/` - landing-page carousel images.

### Existing buyer implementation to preserve

- `BuyerLayoutView.vue`
- `LandingView.vue`
- `SearchResultView.vue`
- `ProductDetailView.vue`
- `CartView.vue`
- `WishlistView.vue`
- `ProfileView.vue`
- `ContactView.vue`
- `AppNavbar.vue`
- `AppFooter.vue`
- `GameCard.vue`
- `ProductListCard.vue`
- `MediaCarousel.vue`
- `AppPagination.vue`
- `PaymentModal.vue`
- `PaymentSuccessModal.vue`
- `SessionExpiredModal.vue`
- `WelcomeModal.vue`

Some shared components already exist even if the old conversion plan marked
them as missing. Verify the current repository state instead of relying on
that older plan.

## Seller screens to implement

### 1. Seller layout

Create a seller layout, preferably:

- `frontend/src/views/SellerLayoutView.vue`
- `frontend/src/components/SellerNavbar.vue`

The seller layout should:

- Render seller navigation and a seller-specific active state.
- Reuse `AppFooter.vue` unless a genuine seller-specific footer is required.
- Render nested seller routes with `<RouterView />`.
- Preserve the existing page background, spacing, typography, and responsive
  behavior.
- Expose links to the seller dashboard, categories, products, and profile.
- Provide a clear way to return to the buyer storefront.
- Show the current seller identity when account data is available.
- Use the existing session-expired and notification patterns.

### 2. Seller dashboard

Create:

- `frontend/src/views/SellerDashboardView.vue`

The dashboard should contain the following sections.

#### Category management

- Display categories in a table or responsive card list.
- Show category name and visibility status.
- Add a category.
- Rename or edit a category.
- Delete a category with confirmation.
- Toggle category visibility.
- Show loading, empty, validation-error, API-error, and success states.
- Prevent accidental duplicate submissions while a request is pending.

#### Product management

- Display the seller's products.
- Show poster, name, price, categories, and relevant status information.
- Link to add a product.
- Link to edit an existing product.
- Delete a product with confirmation.
- Handle an empty product list clearly.
- Do not expose another seller's products in the seller UI.

If the backend currently lacks write endpoints, do not hide the problem with
silent local-only updates. First inspect the backend conventions, then add
the smallest type-safe API support necessary or clearly isolate an adapter
until the endpoint contract is available. Surface failures explicitly.

### 3. Seller product form

Create:

- `frontend/src/views/SellerProductFormView.vue`

The form must support both add and edit modes. It should include:

- Product name.
- Price.
- Category selection.
- Main poster image URL or upload control with preview.
- Up to four gallery image slots.
- Video/media URL inputs where supported by the existing model.
- Description.
- Release date if required by the backend model.
- Minimum hardware requirements.
- Recommended hardware requirements.
- Save/submit action.
- Cancel/back action.
- Clear validation messages.
- Loading state while loading or saving.
- Error state that does not discard the user's entered values.
- Success feedback followed by a safe redirect.

The form should load existing product data in edit mode and avoid treating
missing optional media or hardware values as fatal errors. Keep the model
compatible with the existing `Game`, `GameDetail`, `Hardware`, `Media`, and
`Category` types, updating those types only when the actual API contract
requires it.

### 4. Gallery/media uploader

Create a reusable component, preferably:

- `frontend/src/components/GalleryUploadCard.vue`

Requirements:

- Support existing media URLs in edit mode.
- Support replacing or removing an existing media item.
- Support adding new gallery items.
- Limit the gallery to four image slots.
- Provide a preview or clear visual representation of each slot.
- Validate image URLs/files and video URLs according to the chosen API
  contract.
- Emit typed updates to the parent form.
- Do not silently lose existing media when the user edits unrelated fields.

### 5. Seller profile

Create:

- `frontend/src/views/SellerProfileView.vue`

The seller profile should show:

- Username.
- Account creation date.
- Total products or total games, if available.
- Total sales, if available.
- Seller role information.
- Logout action.
- Delete-account action with an explicit confirmation modal.
- Loading and API-error states.

Reuse the existing account/profile service and authentication conventions
where possible. Do not display password data. Do not use hardcoded seller
statistics in the final implementation.

## Routing requirements

Update `frontend/src/router/index.ts` with a seller route tree similar to:

```text
/seller
/seller/dashboard
/seller/categories
/seller/products/add
/seller/products/:id/edit
/seller/profile
```

Use route names rather than hardcoded navigation where practical. Redirect
`/seller` to the seller dashboard if that matches the existing router style.

The final routing behavior must:

- Keep all existing buyer routes working.
- Keep login and registration routes working.
- Route seller users to the seller area after seller login.
- Prevent buyer-only users from opening seller pages when role information is
  available.
- Avoid redirect loops when the account/session is not yet loaded.
- Use the existing page-title and progress-indicator conventions.

If authentication is currently a placeholder, preserve that behavior while
creating a clear guard seam for the real account/session implementation. Do
not invent insecure client-only authorization as a substitute for backend
authorization.

## API and data-model requirements

Before implementing services, inspect:

- `frontend/src/service/AxiosClient.ts`
- Existing service modules.
- `frontend/src/types.ts`.
- `backend/main.go`.
- Backend handlers, services, repositories, models, and SQL schema.

The current database model includes accounts, products, categories,
hardware requirements, media, and product-category relationships. Relevant
fields include:

- `accounts.is_seller`
- `accounts.total_game`
- `accounts.total_sale`
- `products.user_id`
- `categories.is_hidden`
- `hardware_requirement.type`
- `media.product_id`

The current backend may expose read-only endpoints for accounts, games, and
categories while seller create/update/delete endpoints may still be absent.
Confirm the actual implementation before assuming an endpoint exists. If
backend changes are required for the seller UI, keep them minimal and follow
the existing Go Gin/service/repository patterns.

Use typed request/response shapes. Avoid `any`, broad casts, silent catches,
and fake success responses. Show a user-visible error when an operation
fails.

## Visual direction

Match the existing KD5F storefront visual language:

- Dark blue/purple backgrounds.
- Bright cyan and pink accent colors.
- Rounded cards and controls.
- Clear contrast and readable labels.
- Consistent button, input, modal, table, and notification styles.
- Responsive layouts rather than fixed-width desktop-only screens.

The following existing landing-page assets show the intended branding and
color direction. They are references only; do not change buyer landing-page
behavior while implementing the seller area.

![KD5F Festival landing banner](frontend/public/Landing/Landing-Page-Carousel-1.png)

![KD5F discount landing banner](frontend/public/Landing/Landing-Page-Carousel-2.png)

![Resident Evil promotion banner](frontend/public/Landing/Landing-Page-Carousel-3.png)

![Monster Hunter promotion banner](frontend/public/Landing/Landing-Page-Carousel-4.png)

![The Walking Jord promotion banner](frontend/public/Landing/Landing-Page-Carousel-5.png)

## Implementation rules

1. Read the current file before modifying it.
2. Follow the existing formatting and naming conventions.
3. Prefer existing shared components over duplicated markup.
4. Keep seller-specific logic out of buyer components unless it is genuinely
   shared.
5. Use typed emits, props, refs, computed values, and service responses.
6. Include disabled/loading states for every mutating action.
7. Confirm destructive actions before sending them.
8. Preserve entered form data when a request fails.
9. Make API errors visible to the user and useful during development.
10. Do not add placeholder buttons that appear functional but do nothing.
11. Do not use hardcoded seller data in place of API integration.
12. Keep the implementation accessible: labels, keyboard operation, visible
    focus states, useful alt text, and sensible dialog semantics.
13. Do not modify generated `frontend/dist/` files manually.
14. Do not modify `frontend/package-lock.json` unless dependencies actually
    change.

## Suggested implementation order

1. Inspect the current frontend, backend, SQL schema, and any legacy seller
   templates supplied with this task.
2. Document the actual API contract and identify missing seller endpoints.
3. Add or update shared seller types and services.
4. Build `SellerNavbar.vue` and `SellerLayoutView.vue`.
5. Add seller routes and role-aware navigation.
6. Implement the seller dashboard and category/product management.
7. Implement `GalleryUploadCard.vue`.
8. Implement add/edit product form behavior.
9. Implement seller profile and account actions.
10. Add backend support only for operations that are required and missing.
11. Run type-check, build, lint, and focused manual verification.
12. Review the diff to ensure buyer functionality was not unintentionally
    changed.

## Acceptance criteria

The work is complete only when all of the following are true:

- A seller can open the seller dashboard from the application.
- Seller navigation is separate from buyer navigation and is responsive.
- Categories can be viewed and, where supported by the backend, created,
  edited, hidden/unhidden, and deleted.
- A seller sees only their own products.
- A seller can open an add-product form.
- A seller can open an existing product in edit mode.
- Product validation prevents invalid submissions.
- Poster, gallery, video/media, description, category, price, and hardware
  fields are represented without losing existing values.
- Product save, edit, and delete operations show loading, success, and error
  feedback.
- Seller profile data comes from the account/session source rather than
  hardcoded values.
- Logout and account deletion require the appropriate confirmation/flow.
- Existing buyer pages and routes continue to work.
- `npm run type-check` passes from `frontend/`.
- `npm run build` passes from `frontend/`.
- `npm run lint` passes, or any pre-existing lint failures are identified
  separately from failures introduced by this work.

## Verification commands

Run from `frontend/`:

```sh
npm run type-check
npm run build
npm run lint
```

Also manually verify:

1. Seller login or seller-role navigation.
2. Seller dashboard loading and empty states.
3. Category add/edit/hide/delete flows.
4. Product add flow with validation errors.
5. Product edit flow with existing poster, media, and hardware values.
6. Product deletion confirmation.
7. Seller profile, logout, and account deletion confirmation.
8. Buyer landing page, search, product detail, cart, wishlist, and profile
   routes after seller changes.

## Expected final response from the implementing AI

Report:

- Files created and changed.
- Seller features implemented.
- Backend/API changes, if any.
- Tests, type-check, build, and lint commands run.
- Any remaining limitation caused by an unavailable legacy template or
  missing backend endpoint.

