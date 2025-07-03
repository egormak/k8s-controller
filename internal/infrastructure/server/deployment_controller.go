// package server provides HTTP server functionality using Fiber
package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/tools/cache"

	"k8s-controller/internal/domain"
	"k8s-controller/internal/infrastructure/kubernetes"
)

// DeploymentController handles deployment-related HTTP endpoints
type DeploymentController struct {
	client kubernetes.Client
}

// NewDeploymentController creates a new deployment controller
func NewDeploymentController(client kubernetes.Client) *DeploymentController {
	return &DeploymentController{
		client: client,
	}
}

// ListDeployments handles requests to list deployments
func (c *DeploymentController) ListDeployments(ctx *fiber.Ctx) error {
	// Get namespace from query param, default to "default"
	namespace := ctx.Query("namespace", "default")

	// Create a context with timeout
	reqCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var deployments []domain.Deployment
	var source string

	// Try to get the informer to use its store directly
	informer, err := c.client.GetDeploymentInformer(namespace)

	if err != nil {
		slog.Warn("Could not get deployment informer, falling back to client",
			"namespace", namespace, "error", err)

		// Fall back to the client's ListDeployments which will use API if no informer
		deployments, err = c.client.ListDeployments(reqCtx, namespace)
		if err != nil {
			slog.Error("Failed to list deployments", "error", err, "namespace", namespace)
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to list deployments",
				"error":   err.Error(),
			})
		}
		source = "api"
	} else {
		// Get items directly from the store (cached data)
		storeDeployments, storeErr := c.getDeploymentsFromStore(informer.GetStore(), namespace)
		if storeErr != nil {
			slog.Error("Failed to get deployments from informer store", "error", storeErr, "namespace", namespace)

			// Fall back to the client method
			apiDeployments, apiErr := c.client.ListDeployments(reqCtx, namespace)
			if apiErr != nil {
				return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to list deployments",
					"error":   apiErr.Error(),
				})
			}

			deployments = apiDeployments
			source = "api-fallback"
		} else {
			deployments = storeDeployments
			source = "informer-cache"
		}
	}

	// Return successful response with deployments
	return ctx.JSON(fiber.Map{
		"status":      "success",
		"namespace":   namespace,
		"deployments": deployments,
		"count":       len(deployments),
		"source":      source,
	})
}

// getDeploymentsFromStore converts informer store items to domain deployments
func (c *DeploymentController) getDeploymentsFromStore(store cache.Store, namespace string) ([]domain.Deployment, error) {
	// Get all items from the store
	objs := store.List()

	var deployments []domain.Deployment
	for _, obj := range objs {
		dep, ok := obj.(*appsv1.Deployment)
		if !ok {
			slog.Error("Failed to convert object to Deployment")
			continue
		}

		// Only include deployments from the requested namespace
		if namespace != "" && dep.Namespace != namespace {
			continue
		}

		// Convert to domain model
		deployment := domain.Deployment{
			Name:              dep.Name,
			Namespace:         dep.Namespace,
			ReadyReplicas:     dep.Status.ReadyReplicas,
			UpdatedReplicas:   dep.Status.UpdatedReplicas,
			AvailableReplicas: dep.Status.AvailableReplicas,
			Replicas:          *dep.Spec.Replicas,
			Labels:            dep.Labels,
			CreationTimestamp: dep.CreationTimestamp.Format("2006-01-02 15:04:05"),
		}
		deployments = append(deployments, deployment)
	}

	return deployments, nil
}
