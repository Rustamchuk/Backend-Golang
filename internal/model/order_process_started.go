package model

import "context"

func NewOrderProcessStarted(orderInitialized OrderInitialized, orderStates []OrderState, error error) *OrderProcessStarted {
	return &OrderProcessStarted{
		OrderInitialized: orderInitialized,
		OrderStates:      orderStates,
		Error:            error,
	}
}

func (o *OrderProcessStarted) Pipeline(ctx context.Context, actions OrderActions, orders <-chan OrderInitialized) <-chan OrderProcessStarted {
	outCh := make(chan OrderProcessStarted)
	go func() {
		defer close(outCh)
		for order := range orders {
			states := order.OrderStates
			err := order.Error
			if err == nil {
				err = actions.stateToState(actions.InitToStarted)
				if err == nil {
					states = append(order.OrderStates, ProcessStarted)
				}
			}
			select {
			case <-ctx.Done():
				return
			case outCh <- *NewOrderProcessStarted(order, states, err):
			}
		}
	}()
	return outCh
}
