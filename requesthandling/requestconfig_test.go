package requesthandling

import "testing"

func TestRequestConfigInit(t *testing.T) {
	config := NewRequestConfig("3", "com.test.go", "public")
	if config.Version != "3" {
		t.Errorf("Request Config Version has not been initialised correctly")
	}

	if config.ContainerID != "com.test.go" {
		t.Errorf("Request Container ID has not been initialised correctly")
	}

	if config.Database != "public" {
		t.Errorf("Request Database has not been initialised correctly")
	}
}
