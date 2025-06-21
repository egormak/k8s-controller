package domain

// ResourceEventType represents the type of resource event
type ResourceEventType string

const (
	// ResourceEventCreated indicates a resource was created
	ResourceEventCreated ResourceEventType = "CREATED"

	// ResourceEventUpdated indicates a resource was updated
	ResourceEventUpdated ResourceEventType = "UPDATED"

	// ResourceEventDeleted indicates a resource was deleted
	ResourceEventDeleted ResourceEventType = "DELETED"
)

// Resource represents a Kubernetes resource
type Resource struct {
	Kind       string
	Name       string
	Namespace  string
	APIVersion string
	Labels     map[string]string
	Data       map[string]interface{}
}

// ResourceEvent represents an event that occurred on a Kubernetes resource
type ResourceEvent struct {
	Type     ResourceEventType
	Resource Resource
}
