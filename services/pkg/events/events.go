// Package events defines the shared data structures for events that are
// published to and consumed from the NATS message queue.
package events

// Defines the subjects for NATS messaging.
const (
	SubjectDeploymentRequested = "deployment.requested"
	SubjectBuildSucceeded      = "build.succeeded"
)

// DeploymentRequest is the event payload for a new deployment, published by the
// API service and consumed by the build-worker.
type DeploymentRequest struct {
	AppID         string `json:"app_id"`
	GitRepository string `json:"git_repository"`
	GitBranch     string `json:"git_branch"`
}

// BuildSucceeded is the event payload published by the build-worker when it
// successfully builds a container image. It is consumed by the oal-worker.
type BuildSucceeded struct {
	AppID        string `json:"app_id"`
	ImageURI     string `json:"image_uri"`
	GitCommitSHA string `json:"git_commit_sha"`
}