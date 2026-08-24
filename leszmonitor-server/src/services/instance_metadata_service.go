package services

import (
	"github.com/m-milek/leszmonitor/meta"
)

type InstanceMetadata struct {
	CIBuildNumber string `json:"ciBuildNumber"` // Build number from the CI/CD pipeline
	GitCommit     string `json:"gitCommit"`     // Git commit hash
	ImageTag      string `json:"imageTag"`      // Docker image tag
	Version       string `json:"version"`       // Version of the application
}

type IInstanceMetadataService interface {
	GetInstanceMetadata() InstanceMetadata
}

// InstanceMetadataService provides methods to retrieve instance metadata.
type InstanceMetadataService struct{}

type InstanceMetadataServiceDeps struct{}

// NewInstanceMetadataService creates a new instance of InstanceMetadataService
func NewInstanceMetadataService(deps InstanceMetadataServiceDeps) *InstanceMetadataService {
	return &InstanceMetadataService{}
}

func (s *InstanceMetadataService) GetInstanceMetadata() InstanceMetadata {
	return InstanceMetadata{
		CIBuildNumber: meta.CIBuildNumber,
		GitCommit:     meta.GitCommit,
		ImageTag:      meta.ImageTag,
		Version:       meta.Version,
	}
}
