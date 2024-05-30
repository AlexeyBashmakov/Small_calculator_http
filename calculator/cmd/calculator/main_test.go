package main

import (
	"os"
	"testing"
)

func TestEnvironmentVariables(t *testing.T) {
	val, exist := os.LookupEnv("COMPUTING_POWER")
	if exist {
		prnt(val)
	} else {
		t.Errorf("Expected %d, got '%s'", 8, val)
	}
}
