package controller

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"k8s-controller/internal/infrastructure/config"
)

// ControllerRuntime encapsulates the controller-runtime manager and client
type ControllerRuntime struct {
	manager        manager.Manager
	client         client.Client
	scheme         *runtime.Scheme
	stopCh         chan struct{}
	stopped        bool
	mu             sync.Mutex
	reconcilers    map[string]reconcile.Reconciler
	metricsAddress string
	healthAddress  string
}

// NewControllerRuntime creates a new controller runtime instance
func NewControllerRuntime(cfg *config.Config) (*ControllerRuntime, error) {
	// Set up a Scheme
	scheme := runtime.NewScheme()
	if err := appsv1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("error adding apps/v1 to scheme: %w", err)
	}
	if err := corev1.AddToScheme(scheme); err != nil {
		return nil, fmt.Errorf("error adding core/v1 to scheme: %w", err)
	}

	metricsAddr := ":8081"
	healthAddr := ":8082"

	// Create manager options
	options := ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress: healthAddr,
		LeaderElection:         false,
	}

	// Create manager
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), options)
	if err != nil {
		return nil, fmt.Errorf("unable to create manager: %w", err)
	}

	// Add health and readiness check endpoints
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to set up health check: %w", err)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		return nil, fmt.Errorf("unable to set up ready check: %w", err)
	}

	return &ControllerRuntime{
		manager:        mgr,
		client:         mgr.GetClient(),
		scheme:         scheme,
		stopCh:         make(chan struct{}),
		reconcilers:    make(map[string]reconcile.Reconciler),
		metricsAddress: metricsAddr,
		healthAddress:  healthAddr,
	}, nil
}

// Start starts the controller manager
func (cr *ControllerRuntime) Start(ctx context.Context) error {
	slog.Info("Starting controller manager")
	return cr.manager.Start(ctx)
}

// Stop stops the controller manager
func (cr *ControllerRuntime) Stop() {
	cr.mu.Lock()
	defer cr.mu.Unlock()

	if !cr.stopped {
		close(cr.stopCh)
		cr.stopped = true
		slog.Info("Controller manager stopped")
	}
}

// GetClient returns the controller-runtime client
func (cr *ControllerRuntime) GetClient() client.Client {
	return cr.client
}

// GetManager returns the controller-runtime manager
func (cr *ControllerRuntime) GetManager() manager.Manager {
	return cr.manager
}

// RegisterDeploymentController registers a deployment controller
func (cr *ControllerRuntime) RegisterDeploymentController(reconciler reconcile.Reconciler) error {
	err := ctrl.NewControllerManagedBy(cr.manager).
		For(&appsv1.Deployment{}).
		Complete(reconciler)

	if err != nil {
		return fmt.Errorf("unable to create controller: %w", err)
	}

	cr.reconcilers["deployment"] = reconciler
	return nil
}

// RegisterPodController registers a pod controller
func (cr *ControllerRuntime) RegisterPodController(reconciler reconcile.Reconciler) error {
	err := ctrl.NewControllerManagedBy(cr.manager).
		For(&corev1.Pod{}).
		Complete(reconciler)

	if err != nil {
		return fmt.Errorf("unable to create controller: %w", err)
	}

	cr.reconcilers["pod"] = reconciler
	return nil
}

// GetMetricsEndpoint returns the metrics endpoint address
func (cr *ControllerRuntime) GetMetricsEndpoint() string {
	return cr.metricsAddress
}

// GetHealthEndpoint returns the health endpoint address
func (cr *ControllerRuntime) GetHealthEndpoint() string {
	return cr.healthAddress
}
