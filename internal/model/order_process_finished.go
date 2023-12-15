package model

import (
	"context"
)

func NewOrderProcessFinished(orderFinishedExternalInteraction OrderFinishedExternalInteraction, orderStates []OrderState, error error) *OrderProcessFinished {
	return &OrderProcessFinished{
		OrderFinishedExternalInteraction: orderFinishedExternalInteraction,
		OrderStates:                      orderStates,
		Error:                            error,
	}
}

func (o *OrderProcessFinished) Pipeline(ctx context.Context, actions OrderActions, orders <-chan OrderFinishedExternalInteraction, outCh chan<- OrderProcessFinished) {
	go func() {
		defer close(outCh)
		for order := range orders {
			states := order.OrderStates
			err := order.Error
			if err == nil {
				err = actions.stateToState(actions.FinishedExternalInteractionToProcessFinished)
				if err == nil {
					states = append(order.OrderStates, ProcessFinished)
				}
			}
			select {
			case <-ctx.Done():
				return
			case outCh <- *NewOrderProcessFinished(order, states, err):
			}
		}
	}()
}
