package controller

import (
	"context"
	"log/slog"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"k8s-controller/internal/domain"
)

// DeploymentReconciler reconciles Deployment objects
type DeploymentReconciler struct {
	client client.Client
	scheme *runtime.Scheme
	// Add a reference to the domain service if needed
	resourceService domain.ResourceService
}

// NewDeploymentReconciler creates a new deployment reconciler
func NewDeploymentReconciler(client client.Client, scheme *runtime.Scheme, resourceService domain.ResourceService) *DeploymentReconciler {
	return &DeploymentReconciler{
		client:          client,
		scheme:          scheme,
		resourceService: resourceService,
	}
}

// Reconcile implements the reconcile.Reconciler interface
func (r *DeploymentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// Get the Deployment object
	var deployment appsv1.Deployment
	if err := r.client.Get(ctx, req.NamespacedName, &deployment); err != nil {
		if errors.IsNotFound(err) {
			// The object was deleted
			slog.Info("Deployment was deleted", "name", req.Name, "namespace", req.Namespace)
			return ctrl.Result{}, nil
		}
		slog.Error("Failed to get Deployment", "name", req.Name, "namespace", req.Namespace, "error", err)
		return ctrl.Result{}, err
	}

	// Convert k8s deployment to domain deployment
	domainDeployment := domain.Deployment{
		Name:      deployment.Name,
		Namespace: deployment.Namespace,
		Replicas:  *deployment.Spec.Replicas,
		Status: domain.DeploymentStatus{
			AvailableReplicas:   deployment.Status.AvailableReplicas,
			UnavailableReplicas: deployment.Status.UnavailableReplicas,
			ReadyReplicas:       deployment.Status.ReadyReplicas,
		},
		Labels:    deployment.Labels,
		CreatedAt: deployment.CreationTimestamp.Time,
	}

	// Process the domain deployment using the resource service
	if r.resourceService != nil {
		if err := r.resourceService.ProcessDeployment(ctx, domainDeployment); err != nil {
			slog.Error("Failed to process deployment", "name", deployment.Name, "error", err)
			// Requeue after 30 seconds
			return ctrl.Result{RequeueAfter: 30 * time.Second}, nil
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *DeploymentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.Deployment{}).
		Complete(r)
}
