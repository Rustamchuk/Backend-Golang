package workerpool

import (
	"context"
	"hw7/internal/pipeline"
	"sync"

	"hw7/internal/model"
)

type OrderWorkerPool interface {
	StartWorkerPool(ctx context.Context, orders <-chan model.OrderInitialized, additionalActions model.OrderActions, workersCount int) <-chan model.OrderProcessFinished
}

type OrderWorkerPoolImplementation struct{}

func NewOrderWorkerPoolImplementation() *OrderWorkerPoolImplementation {
	return &OrderWorkerPoolImplementation{}
}

func (o *OrderWorkerPoolImplementation) StartWorkerPool(
	ctx context.Context,
	orders <-chan model.OrderInitialized,
	additionalActions model.OrderActions,
	workersCount int,
) <-chan model.OrderProcessFinished {
	outCh := make(chan model.OrderProcessFinished)
	pl := pipeline.NewOrderPipelineImplementation()
	wg := sync.WaitGroup{}

	wg.Add(workersCount)
	for i := 0; i < workersCount; i++ {
		go o.worker(ctx, &wg, pl, additionalActions, orders, outCh)
	}

	go func() {
		wg.Wait()
		close(outCh)
	}()
	return outCh
}

func (o *OrderWorkerPoolImplementation) worker(
	ctx context.Context,
	wg *sync.WaitGroup,
	pl pipeline.OrderPipeline,
	actions model.OrderActions,
	orders <-chan model.OrderInitialized,
	outCh chan<- model.OrderProcessFinished,
) {
	defer wg.Done()
	result := make(chan model.OrderProcessFinished)

	pl.Start(ctx, actions, orders, result)

	for r := range result {
		select {
		case <-ctx.Done():
			return
		case outCh <- r:
		}
	}
}
