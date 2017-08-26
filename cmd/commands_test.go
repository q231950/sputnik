package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRootCommand(t *testing.T) {
	flag := RootCmd.Flag("config")
	fmt.Println(flag)
}

func TestRequestsCommand(t *testing.T) {
	run := requestsCmd.Run
	assert.NotNil(t, run)
}
