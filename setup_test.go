package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_setup(t *testing.T) {
	assert.NoError(t, os.Setenv(shipyardModulesEnvKey, "/default.hcl"))
	assert.True(t, gotoTargetDir())
}
func Test_filter(t *testing.T) {
	assert.Equal(t, []string{"3", "4"}, filterSlice([]string{"1", "2", "3", "4", "1"}, "1", "2"))
}
