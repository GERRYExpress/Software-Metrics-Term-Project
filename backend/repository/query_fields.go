package repository

const (
	cartUserIDColumn    = "ca.user_id"
	cartProductIDColumn = "ca.product_id"
	cartDateColumn      = "ca.cart_date"

	wishlistUserIDColumn    = "w.user_id"
	wishlistProductIDColumn = "w.product_id"
	wishlistDateColumn      = "w.wishlist_date"

	orderIDColumn       = "o.order_id"
	orderUserIDColumn   = "o.user_id"
	orderDateColumn     = "o.order_date"
	paymentMethodColumn = "o.payment_method"
	finalPriceColumn    = "o.final_price"
)
