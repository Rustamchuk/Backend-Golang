package model

import (
	"context"
	"sync"
)

func NewOrderFinishedExternalInteraction(orderProcessStarted OrderProcessStarted, storageID int, pickupPointID int, orderStates []OrderState, error error) *OrderFinishedExternalInteraction {
	return &OrderFinishedExternalInteraction{
		OrderProcessStarted: orderProcessStarted,
		StorageID:           storageID,
		PickupPointID:       pickupPointID,
		OrderStates:         orderStates,
		Error:               error,
	}
}

func (o *OrderFinishedExternalInteraction) PipelineFan(ctx context.Context, actions OrderActions, k int, orders <-chan OrderProcessStarted) <-chan OrderFinishedExternalInteraction {
	fanOutInteract := make([]<-chan OrderFinishedExternalInteraction, k)
	for i := 0; i < k; i++ {
		fanOutInteract[i] = o.process(ctx, actions, orders)
	}
	outCh := o.fanIn(fanOutInteract)
	return outCh
}

func (o *OrderFinishedExternalInteraction) process(ctx context.Context, actions OrderActions, orders <-chan OrderProcessStarted) <-chan OrderFinishedExternalInteraction {
	outCh := make(chan OrderFinishedExternalInteraction)
	go func() {
		defer close(outCh)
		for order := range orders {
			storageID, pickupPointID := 0, 0
			states := order.OrderStates
			err := order.Error
			if err == nil {
				err = actions.stateToState(actions.StartedToFinishedExternalInteraction)
				if err == nil {
					states = append(order.OrderStates, FinishedExternalInteraction)
					storageID = o.countStorageID(order.OrderInitialized.ProductID)
					pickupPointID = o.countPickUpID(order.OrderInitialized.ProductID)
				}
			}
			select {
			case <-ctx.Done():
				return
			case outCh <- *NewOrderFinishedExternalInteraction(order, storageID, pickupPointID, states, err):
			}
		}
	}()
	return outCh
}

func (o *OrderFinishedExternalInteraction) fanIn(fanOutInteract []<-chan OrderFinishedExternalInteraction) <-chan OrderFinishedExternalInteraction {
	outCh := make(chan OrderFinishedExternalInteraction)
	wgFan := sync.WaitGroup{}
	for _, fanIn := range fanOutInteract {
		wgFan.Add(1)
		go func(fanIn <-chan OrderFinishedExternalInteraction) {
			defer wgFan.Done()
			for t := range fanIn {
				outCh <- t
			}
		}(fanIn)
	}
	go func() {
		wgFan.Wait()
		close(outCh)
	}()
	return outCh
}

func (o *OrderFinishedExternalInteraction) countStorageID(productID int) int {
	return productID%2 + 1
}

func (o *OrderFinishedExternalInteraction) countPickUpID(productID int) int {
	return productID%3 + 1
}
