package inapppay

// Product represents an in-app product.
type Product struct {
	ID          string
	Title       string
	Description string
	Price       string
	// OriginalPrice might be useful for promotions
	OriginalPrice string
	CurrencyCode  string
}

// ProductDetailsEvent lists the products available.
type ProductDetailsEvent struct {
	Products []Product
}

func (ProductDetailsEvent) ImplementsEvent() {}

// PaymentResultEvent contains the result of a payment.
type PaymentResultEvent struct {
	ProductID        string
	PurchaseID       string // Or OrderID
	Status           PaymentStatus
	DeveloperPayload string
	OriginalJSON     string // Verify signature
	Signature        string
}

func (PaymentResultEvent) ImplementsEvent() {}

// PaymentStatus indicates the status of the purchase.
type PaymentStatus int

const (
	PaymentStatusPending PaymentStatus = iota
	PaymentStatusPurchased
	PaymentStatusCancelled
	PaymentStatusError
)

// ErrorEvent reports an error.
type ErrorEvent struct {
	Error error
}

func (ErrorEvent) ImplementsEvent() {}

// Event is the general interface for all events.
type Event interface {
	ImplementsEvent()
}
