package pipeline

import (
	"context"
	"hw7/internal/model"
)

const fanLimit = 10

type OrderPipeline interface {
	Start(ctx context.Context, actions model.OrderActions, orders <-chan model.OrderInitialized, processed chan<- model.OrderProcessFinished)
}

type OrderPipelineImplementation struct{}

func NewOrderPipelineImplementation() *OrderPipelineImplementation {
	return &OrderPipelineImplementation{}
}

func (o *OrderPipelineImplementation) Start(
	ctx context.Context,
	actions model.OrderActions,
	orders <-chan model.OrderInitialized,
	processed chan<- model.OrderProcessFinished,
) {
	stateStart := model.OrderProcessStarted{}
	stateInteract := model.OrderFinishedExternalInteraction{}
	stateFinish := model.OrderProcessFinished{}

	started := stateStart.Pipeline(ctx, actions, orders)
	interact := stateInteract.PipelineFan(ctx, actions, fanLimit, started)
	stateFinish.Pipeline(ctx, actions, interact, processed)
}
