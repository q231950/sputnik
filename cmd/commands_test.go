package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	flag := RootCmd.Flag("config")
	assert.NotNil(t, flag)
}

func TestRequestsCommand(t *testing.T) {
	run := requestsCmd.Run
	assert.NotNil(t, run)
}

func TestPostRequestCommand(t *testing.T) {
	run := postCmd.Run
	assert.NotNil(t, run)
}

func TestPostRequestCommandJSONFlag(t *testing.T) {
	flag := postCmd.Flag("json-file-path")
	assert.NotNil(t, flag)
}

func TestPostRequestCommandPayloadFlag(t *testing.T) {
	flag := postCmd.Flag("payload")
	assert.NotNil(t, flag)
}

func TestPostRequestCommandOperationFlag(t *testing.T) {
	flag := postCmd.Flag("operation")
	assert.NotNil(t, flag)
}
