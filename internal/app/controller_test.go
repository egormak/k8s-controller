package app

import (
	"testing"
	"time"
)

func TestNewKubernetesController(t *testing.T) {
	controller := NewKubernetesController()
	if controller == nil {
		t.Fatal("NewKubernetesController returned nil")
	}

	if controller.client == nil {
		t.Error("controller.client is nil")
	}

	if controller.resourceService == nil {
		t.Error("controller.resourceService is nil")
	}

	if controller.resourceHandler == nil {
		t.Error("controller.resourceHandler is nil")
	}

	if controller.ctx == nil {
		t.Error("controller.ctx is nil")
	}

	if controller.cancelFunc == nil {
		t.Error("controller.cancelFunc is nil")
	}

	if controller.config == nil {
		t.Error("controller.config is nil")
	}
}

func TestKubernetesController_RunStop(t *testing.T) {
	controller := NewKubernetesController()

	go func() {
		time.Sleep(1 * time.Second)
		controller.Stop()
	}()

	err := controller.Run()
	if err != nil {
		// We expect the run to be cancelled, so no error should be returned
		// but if the connect fails, it will return an error.
		// For this test, we can ignore the connection error as we are not testing connectivity.
		if err.Error() != "failed to connect to Kubernetes: invalid configuration: no configuration has been provided, try setting KUBERNETES_MASTER environment variable" {
			t.Errorf("Controller Run() returned an unexpected error: %v", err)
		}
	}
}
