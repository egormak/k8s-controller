package domain

import "time"

// DeploymentStatus contains status information for a deployment
type DeploymentStatus struct {
	ReadyReplicas       int32
	UpdatedReplicas     int32
	AvailableReplicas   int32
	UnavailableReplicas int32
}

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
	Status            DeploymentStatus
	CreatedAt         time.Time
}
