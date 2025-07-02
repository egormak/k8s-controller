package domain

// Deployment represents a Kubernetes deployment
type Deployment struct {
	Name              string
	Namespace         string
	ReadyReplicas     int32
	UpdatedReplicas   int32
	AvailableReplicas int32
	Replicas          int32
	Labels            map[string]string
	CreationTimestamp string
}
