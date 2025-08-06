package docs

import "github.com/kimenyu/executive/types"

// This file is only for Swagger documentation generation purposes.
// It helps ensure all relevant models are included in the Swagger schema,
// even if they are not directly referenced in annotations.

var (
	_ = types.User{}
	_ = types.RegisterUserPayload{}
	_ = types.LoginUserPayload{}

	_ = types.Category{}
	_ = types.CreateCategoryPayload{}

	_ = types.Product{}
	_ = types.CreateProductPayload{}

	_ = types.Cart{}
	_ = types.CartItem{}
	_ = types.AddToCartPayload{}

	_ = types.Order{}
	_ = types.OrderItem{}
	_ = types.CreateOrderPayload{}
	_ = types.CreateOrderItemDTO{}

	_ = types.Review{}
	_ = types.CreateReviewPayload{}

	_ = types.Address{}
)
