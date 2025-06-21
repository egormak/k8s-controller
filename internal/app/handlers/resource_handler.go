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

	switch event.Type {
	case domain.ResourceEventCreated:
		return h.handleCreated(ctx, event.Resource)
	case domain.ResourceEventUpdated:
		return h.handleUpdated(ctx, event.Resource)
	case domain.ResourceEventDeleted:
		return h.handleDeleted(ctx, event.Resource)
	default:
		slog.Warn("Unknown event type", "type", event.Type)
		return nil
	}
}

// handleCreated processes a created resource
func (h *ResourceHandler) handleCreated(ctx context.Context, resource domain.Resource) error {
	slog.Info("Resource created",
		"kind", resource.Kind,
		"name", resource.Name,
		"namespace", resource.Namespace)

	// Implement business logic for creation events
	return nil
}

// handleUpdated processes an updated resource
func (h *ResourceHandler) handleUpdated(ctx context.Context, resource domain.Resource) error {
	slog.Info("Resource updated",
		"kind", resource.Kind,
		"name", resource.Name,
		"namespace", resource.Namespace)

	// Implement business logic for update events
	return nil
}

// handleDeleted processes a deleted resource
func (h *ResourceHandler) handleDeleted(ctx context.Context, resource domain.Resource) error {
	slog.Info("Resource deleted",
		"kind", resource.Kind,
		"name", resource.Name,
		"namespace", resource.Namespace)

	// Implement business logic for deletion events
	return nil
}
