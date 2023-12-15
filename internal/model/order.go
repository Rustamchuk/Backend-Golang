package model

func NewOrder(orderID int, productID int, storageID int, pickupPointID int, isProcessed bool, orderStates []OrderState) *Order {
	return &Order{OrderID: orderID, ProductID: productID, StorageID: storageID, PickupPointID: pickupPointID, IsProcessed: isProcessed, OrderStates: orderStates}
}
