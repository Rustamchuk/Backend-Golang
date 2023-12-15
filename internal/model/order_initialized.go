package model

func NewOrderInitialized(orderID int, productID int, error error) *OrderInitialized {
	return &OrderInitialized{
		OrderID:     orderID,
		ProductID:   productID,
		OrderStates: []OrderState{Initialized},
		Error:       error,
	}
}
