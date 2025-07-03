package domain

import (
	"context"
	"log/slog"
)

// ResourceClient defines the interface for interacting with Kubernetes resources
type ResourceClient interface {
	Connect(ctx context.Context) error
	WatchResources(ctx context.Context) error
	GetResource(ctx context.Context, kind, name, namespace string) (Resource, error)
	ApplyResource(ctx context.Context, resource Resource) error
}

// ResourceService is the domain service for handling Kubernetes resource operations
type ResourceService interface {
	WatchResources(ctx context.Context) error
	HandleResourceEvent(ctx context.Context, event ResourceEvent) error
}

// resourceService implements the ResourceService interface
type resourceService struct {
	client ResourceClient
}

// NewResourceService creates a new resource service
func NewResourceService(client ResourceClient) ResourceService {
	return &resourceService{
		client: client,
	}
}

// WatchResources starts watching for resource events
func (s *resourceService) WatchResources(ctx context.Context) error {
	slog.Info("Starting to watch resources")
	return s.client.WatchResources(ctx)
}

// HandleResourceEvent processes a resource event
func (s *resourceService) HandleResourceEvent(ctx context.Context, event ResourceEvent) error {
	slog.Info("Handling resource event",
		"kind", event.Resource.Kind,
		"name", event.Resource.Name,
		"namespace", event.Resource.Namespace,
		"eventType", event.Type)

	// Apply business logic based on the resource event
	switch event.Resource.Kind {
	case "Deployment":
		slog.Info("Deployment event detected",
			"name", event.Resource.Name,
			"namespace", event.Resource.Namespace,
			"eventType", event.Type)
		// Add specific deployment processing logic here if needed
	case "Service":
		slog.Info("Service event detected",
			"name", event.Resource.Name,
			"namespace", event.Resource.Namespace,
			"eventType", event.Type)
	case "Pod":
		slog.Info("Pod event detected",
			"name", event.Resource.Name,
			"namespace", event.Resource.Namespace,
			"eventType", event.Type)
	}

	return nil
}
