package domain

import (
	"context"
	"testing"
)

// MockResourceClient is a mock implementation of the ResourceClient interface
type MockResourceClient struct {
	ConnectFunc        func(ctx context.Context) error
	WatchResourcesFunc func(ctx context.Context) error
	GetResourceFunc    func(ctx context.Context, kind, name, namespace string) (Resource, error)
	ApplyResourceFunc  func(ctx context.Context, resource Resource) error
}

func (m *MockResourceClient) Connect(ctx context.Context) error {
	if m.ConnectFunc != nil {
		return m.ConnectFunc(ctx)
	}
	return nil
}

func (m *MockResourceClient) WatchResources(ctx context.Context) error {
	if m.WatchResourcesFunc != nil {
		return m.WatchResourcesFunc(ctx)
	}
	return nil
}

func (m *MockResourceClient) GetResource(ctx context.Context, kind, name, namespace string) (Resource, error) {
	if m.GetResourceFunc != nil {
		return m.GetResourceFunc(ctx, kind, name, namespace)
	}
	return Resource{}, nil
}

func (m *MockResourceClient) ApplyResource(ctx context.Context, resource Resource) error {
	if m.ApplyResourceFunc != nil {
		return m.ApplyResourceFunc(ctx, resource)
	}
	return nil
}

func TestHandleResourceEvent(t *testing.T) {
	mockClient := &MockResourceClient{}
	service := NewResourceService(mockClient)

	event := ResourceEvent{
		Type: "ADDED",
		Resource: Resource{
			Kind:      "Pod",
			Name:      "test-pod",
			Namespace: "default",
		},
	}

	err := service.HandleResourceEvent(context.Background(), event)
	if err != nil {
		t.Errorf("HandleResourceEvent failed: %v", err)
	}

	// Here you could add assertions to check if the mock client's methods were called, for example.
	// For this simple test, we just check that no error is returned.
}
