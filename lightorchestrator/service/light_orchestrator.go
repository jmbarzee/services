package service

import (
	"context"
	"time"

	"github.com/jmbarzee/dominion/service"
	"github.com/jmbarzee/dominion/service/config"
	"github.com/jmbarzee/dominion/system"
	"github.com/jmbarzee/services/lightorchestrator/grpc"
)

const (
	displayFPS                = 30
	displayRate time.Duration = time.Second / displayFPS
)

type LightOrch struct {
	// UnimplementedLightOrchestratorServer is embedded to enable forwards compatability
	grpc.UnimplementedLightOrchestratorServer

	*service.Service

	Subscribers *SubscriberList

	NodeTree *NodeTree
}

func NewLightOrch(config config.ServiceConfig) (*LightOrch, error) {
	service, err := service.NewService(config)
	if err != nil {
		return nil, err
	}

	subscriberList, deviceNodeTree := NewStructs()

	lightOrch := &LightOrch{
		Service:     service,
		Subscribers: subscriberList,
		NodeTree:    deviceNodeTree,
	}

	grpc.RegisterLightOrchestratorServer(service.Server, lightOrch)
	return lightOrch, nil
}

func (l *LightOrch) Run(ctx context.Context) error {
	system.Logf("I seek to join the Dominion\n")
	system.Logf(l.ServiceIdentity.String())
	system.Logf("The Dominion ever expands!\n")

	go system.RoutineOperation(ctx, "allocateVibe", tickLength, l.allocateVibe)
	go system.RoutineOperation(ctx, "dispatchRender", displayRate, l.dispatchRender)

	return l.Service.HostService(ctx)
}
