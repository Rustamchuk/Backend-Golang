package main

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"hw7/internal/generator"
	"hw7/internal/model"
	"hw7/internal/workerpool"
	"sort"
	"testing"
)

func TestOrder(t *testing.T) {
	t.Parallel()
	fullStates := []model.OrderState{model.Initialized, model.ProcessStarted, model.FinishedExternalInteraction, model.ProcessFinished}

	tests := []struct {
		name       string
		orders     []model.OrderInitialized
		wc         int
		waitOrders []model.Order
		waitCount  CountChecker
	}{
		{
			name: "success_normal_test",
			orders: []model.OrderInitialized{
				*model.NewOrderInitialized(1, 1, nil),
				*model.NewOrderInitialized(2, 2, errors.New("")),
				*model.NewOrderInitialized(3, 3, nil),
				*model.NewOrderInitialized(4, 4, nil),
			},
			wc: 5,
			waitOrders: []model.Order{
				*model.NewOrder(1, 1, 2, 2, true, fullStates),
				*model.NewOrder(2, 2, 0, 0, false, []model.OrderState{model.Initialized}),
				*model.NewOrder(3, 3, 2, 1, true, fullStates),
				*model.NewOrder(4, 4, 1, 2, true, fullStates),
			},
			waitCount: CountChecker{3, 3, 3},
		},
		{
			name: "success_normal_test_workers_1",
			orders: []model.OrderInitialized{
				*model.NewOrderInitialized(1, 1, nil),
				*model.NewOrderInitialized(2, 2, errors.New("")),
				*model.NewOrderInitialized(3, 3, nil),
				*model.NewOrderInitialized(4, 4, nil),
			},
			wc: 1,
			waitOrders: []model.Order{
				*model.NewOrder(1, 1, 2, 2, true, fullStates),
				*model.NewOrder(2, 2, 0, 0, false, []model.OrderState{model.Initialized}),
				*model.NewOrder(3, 3, 2, 1, true, fullStates),
				*model.NewOrder(4, 4, 1, 2, true, fullStates),
			},
			waitCount: CountChecker{3, 3, 3},
		},
		{
			name: "success_normal_test_workers_0",
			orders: []model.OrderInitialized{
				*model.NewOrderInitialized(1, 1, nil),
				*model.NewOrderInitialized(2, 2, errors.New("")),
				*model.NewOrderInitialized(3, 3, nil),
				*model.NewOrderInitialized(4, 4, nil),
			},
			wc:         0,
			waitOrders: []model.Order{},
			waitCount:  CountChecker{0, 0, 0},
		},
		{
			name:       "success_normal_test_orders_empty",
			orders:     []model.OrderInitialized{},
			wc:         5,
			waitOrders: []model.Order{},
			waitCount:  CountChecker{0, 0, 0},
		},
		{
			name: "errors_orders_test",
			orders: []model.OrderInitialized{
				*model.NewOrderInitialized(1, 1, errors.New("")),
				*model.NewOrderInitialized(2, 2, errors.New("")),
				*model.NewOrderInitialized(3, 3, errors.New("")),
				*model.NewOrderInitialized(4, 4, errors.New("")),
			},
			wc: 5,
			waitOrders: []model.Order{
				*model.NewOrder(1, 1, 0, 0, false, []model.OrderState{model.Initialized}),
				*model.NewOrder(2, 2, 0, 0, false, []model.OrderState{model.Initialized}),
				*model.NewOrder(3, 3, 0, 0, false, []model.OrderState{model.Initialized}),
				*model.NewOrder(4, 4, 0, 0, false, []model.OrderState{model.Initialized}),
			},
			waitCount: CountChecker{0, 0, 0},
		},
	}

	ctx := context.Background()
	orderWorkerPool := workerpool.NewOrderWorkerPoolImplementation()
	orderGenerator := generator.NewOrderGeneratorImplementation()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actions, counter := getDefaultAdditionalActions()
			results := orderWorkerPool.StartWorkerPool(ctx, orderGenerator.GenerateOrdersStream(ctx, test.orders), actions, test.wc)

			resArr := make([]model.Order, 0)
			for res := range results {
				resArr = append(resArr, finishedStateToOrder(res))
			}

			sort.SliceStable(resArr, func(i, j int) bool {
				return resArr[i].OrderID < resArr[j].OrderID
			})

			require.Equal(t, len(test.waitOrders), len(resArr))
			require.Equal(t, test.waitCount, *counter)

			for i, res := range resArr {
				require.Equal(t, test.waitOrders[i].OrderID, res.OrderID)
				require.Equal(t, test.waitOrders[i].ProductID, res.ProductID)
				require.Equal(t, test.waitOrders[i].StorageID, res.StorageID)
				require.Equal(t, test.waitOrders[i].PickupPointID, res.PickupPointID)
				require.Equal(t, test.waitOrders[i].IsProcessed, res.IsProcessed)
				require.Equal(t, test.waitOrders[i].OrderStates, res.OrderStates)
			}
		})
	}
}
