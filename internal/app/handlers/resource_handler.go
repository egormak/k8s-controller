package handlers

import (
	"context"
	"log/slog"

	"k8s-controller/internal/domain"
)

// ResourceHandler processes resource events based on business rules
type ResourceHandler struct {
	resourceService domain.ResourceService
}

// NewResourceHandler creates a new resource handler
func NewResourceHandler(resourceService domain.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		resourceService: resourceService,
	}
}

// HandleEvent processes a resource event
func (h *ResourceHandler) HandleEvent(ctx context.Context, event domain.ResourceEvent) error {
	slog.Debug("Processing resource event",
		"type", event.Type,
		"kind", event.Resource.Kind,
		"name", event.Resource.Name,
		"namespace", event.Resource.Namespace)

	// Forward event to the domain service
	return h.resourceService.HandleResourceEvent(ctx, event)
}
