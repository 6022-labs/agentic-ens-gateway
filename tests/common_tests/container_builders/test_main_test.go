package container_builders

import (
	"os"
	"testing"
)

// TestMain ensures containers are stopped once this test
// package process finishes running tests. Each test package/process
// should manage stopping its own container lifecycle.
func TestMain(m *testing.M) {
	code := m.Run()
	stopRedis()
	os.Exit(code)
}
