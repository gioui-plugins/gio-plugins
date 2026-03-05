package gioinapppay

import "reflect"

var wantCommands = []reflect.Type{
	reflect.TypeOf(ListProductsCmd{}),
	reflect.TypeOf(PurchaseCmd{}),
}

// ListProductsCmd commands to list products.
type ListProductsCmd struct {
	// ProductIDs is the list of product IDs to list.
	ProductIDs []string
}

func (c ListProductsCmd) ImplementsCommand() {}

// PurchaseCmd commands to purchase a product.
type PurchaseCmd struct {
	// ProductID is the ID of the product to purchase.
	ProductID string
	// CustomPayload is the custom payload to send with the purchase.
	CustomPayload string
	// IsPersonalizedPrice is whether the price is personalized for the user.
	IsPersonalizedPrice bool
}

func (c PurchaseCmd) ImplementsCommand() {}
