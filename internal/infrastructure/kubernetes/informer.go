package kubernetes

import (
	"context"
	"log/slog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"

	"k8s-controller/internal/domain"
)

// ResourceEventHandler is an interface for handling resource events
type ResourceEventHandler interface {
	HandleEvent(ctx context.Context, event domain.ResourceEvent) error
}

// startInformers initializes and starts informers for the given resources
func (c *kubeClient) startInformers(ctx context.Context, namespaces []string, resources []string, handler ResourceEventHandler) error {
	slog.Info("Starting informers", "namespaces", namespaces, "resources", resources)

	// Create a factory for each namespace
	for _, namespace := range namespaces {
		factory := informers.NewSharedInformerFactoryWithOptions(
			c.clientset,
			30*time.Second, // resync period
			informers.WithNamespace(namespace),
		)

		// Set up informers for each resource type
		for _, resource := range resources {
			if err := c.setupInformer(ctx, factory, resource, namespace, handler); err != nil {
				return err
			}
		}

		// Start the informer factory
		factory.Start(ctx.Done())
	}

	return nil
}

// setupInformer creates an informer for a specific resource type
func (c *kubeClient) setupInformer(ctx context.Context, factory informers.SharedInformerFactory, resource string, namespace string, handler ResourceEventHandler) error {
	var informer cache.SharedIndexInformer

	// Configure the appropriate informer based on resource type
	switch resource {
	case "pods", "pod":
		informer = factory.Core().V1().Pods().Informer()
	case "services", "service":
		informer = factory.Core().V1().Services().Informer()
	case "deployments", "deployment":
		informer = factory.Apps().V1().Deployments().Informer()
	case "configmaps", "configmap":
		informer = factory.Core().V1().ConfigMaps().Informer()
	default:
		slog.Warn("Unsupported resource type", "resource", resource)
		return nil
	}

	// Add event handlers
	_, err := informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			c.handleAddEvent(ctx, obj, handler)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			c.handleUpdateEvent(ctx, oldObj, newObj, handler)
		},
		DeleteFunc: func(obj interface{}) {
			c.handleDeleteEvent(ctx, obj, handler)
		},
	})

	if err != nil {
		slog.Error("Failed to add event handler", "resource", resource, "namespace", namespace, "error", err)
		return err
	}

	slog.Info("Informer configured", "resource", resource, "namespace", namespace)
	return nil
}

// handleAddEvent processes resource creation events
func (c *kubeClient) handleAddEvent(ctx context.Context, obj interface{}, handler ResourceEventHandler) {
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		slog.Error("Failed to convert object to metav1.Object")
		return
	}

	// Convert to domain model
	resource := c.convertToDomainResource(obj)
	event := domain.ResourceEvent{
		Type:     domain.ResourceEventCreated,
		Resource: resource,
	}

	// Process the event
	if err := handler.HandleEvent(ctx, event); err != nil {
		slog.Error("Failed to handle add event",
			"name", metaObj.GetName(),
			"namespace", metaObj.GetNamespace(),
			"error", err)
	}
}

// handleUpdateEvent processes resource update events
func (c *kubeClient) handleUpdateEvent(ctx context.Context, oldObj, newObj interface{}, handler ResourceEventHandler) {
	metaObj, ok := newObj.(metav1.Object)
	if !ok {
		slog.Error("Failed to convert object to metav1.Object")
		return
	}

	// Convert to domain model
	resource := c.convertToDomainResource(newObj)
	event := domain.ResourceEvent{
		Type:     domain.ResourceEventUpdated,
		Resource: resource,
	}

	// Process the event
	if err := handler.HandleEvent(ctx, event); err != nil {
		slog.Error("Failed to handle update event",
			"name", metaObj.GetName(),
			"namespace", metaObj.GetNamespace(),
			"error", err)
	}
}

// handleDeleteEvent processes resource deletion events
func (c *kubeClient) handleDeleteEvent(ctx context.Context, obj interface{}, handler ResourceEventHandler) {
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		// Handle deleted objects that might be tombstones
		tombstone, ok := obj.(cache.DeletedFinalStateUnknown)
		if !ok {
			slog.Error("Failed to convert object to metav1.Object or tombstone")
			return
		}
		metaObj, ok = tombstone.Obj.(metav1.Object)
		if !ok {
			slog.Error("Failed to convert tombstone to metav1.Object")
			return
		}
	}

	// Convert to domain model
	resource := c.convertToDomainResource(obj)
	event := domain.ResourceEvent{
		Type:     domain.ResourceEventDeleted,
		Resource: resource,
	}

	// Process the event
	if err := handler.HandleEvent(ctx, event); err != nil {
		slog.Error("Failed to handle delete event",
			"name", metaObj.GetName(),
			"namespace", metaObj.GetNamespace(),
			"error", err)
	}
}

// convertToDomainResource converts a Kubernetes object to a domain resource
func (c *kubeClient) convertToDomainResource(obj interface{}) domain.Resource {
	// Handle tombstones from deletion events
	if tombstone, isTombstone := obj.(cache.DeletedFinalStateUnknown); isTombstone {
		obj = tombstone.Obj
	}

	// Get metadata from the object
	metaObj, ok := obj.(metav1.Object)
	if !ok {
		slog.Error("Failed to convert object to metav1.Object")
		return domain.Resource{}
	}

	// Extract kind using runtime.Object
	kind := "Unknown"
	if runtimeObj, ok := obj.(runtime.Object); ok {
		gvk := runtimeObj.GetObjectKind().GroupVersionKind()
		if gvk.Kind != "" {
			kind = gvk.Kind
		} else {
			// Use type reflection as a fallback for kind detection
			kind = c.getKindFromResourceType(obj)
		}
	}

	// Create the domain resource
	return domain.Resource{
		Kind:      kind,
		Name:      metaObj.GetName(),
		Namespace: metaObj.GetNamespace(),
		Labels:    metaObj.GetLabels(),
	}
}

// getKindFromResourceType attempts to determine the kind based on the type of the object
func (c *kubeClient) getKindFromResourceType(obj interface{}) string {
	// Simple resource type detection based on specific object types
	switch obj.(type) {
	case *corev1.Pod:
		return "Pod"
	case *corev1.Service:
		return "Service"
	case *appsv1.Deployment:
		return "Deployment"
	case *corev1.ConfigMap:
		return "ConfigMap"
	default:
		return "Unknown"
	}
}
