// Package events defines the shared data structures for events that are
// published to and consumed from the NATS message queue.
package events

// Defines the subjects for NATS messaging.
const (
	SubjectDeploymentRequested = "v1.deployment.requested"
	SubjectBuildSucceeded      = "v1.build.succeeded"
)

// DeploymentRequest is the event payload for a new deployment, published by the
// API service and consumed by the build-worker.
type DeploymentRequest struct {
	AppID         string `json:"app_id" validate:"required"`
	GitRepository string `json:"git_repository" validate:"required,url"`
	GitBranch     string `json:"git_branch" validate:"required"`
}

// BuildSucceeded is the event payload published by the build-worker when it
// successfully builds a container image. It is consumed by the oal-worker.
type BuildSucceeded struct {
	AppID        string `json:"app_id" validate:"required"`
	ImageURI     string `json:"image_uri" validate:"required"`
	GitCommitSHA string `json:"git_commit_sha" validate:"required"`
}