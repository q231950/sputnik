package requesthandling

// RequestConfig is used to initialise RequestManagers. It specifies the Cloudkit API version and container ID to use for requests
type RequestConfig struct {
	Version     string
	ContainerID string
	Database    string
}

// NewRequestConfig creates a fresh config with the given version and container ID
func NewRequestConfig(version string, containerID string, database string) RequestConfig {
	return RequestConfig{Version: version, ContainerID: containerID, Database: database}
}
