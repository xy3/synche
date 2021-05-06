package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidDirName(t *testing.T) {
	validDirNames := []string{"DirName", "nnnnn", "😄"}
	invalidDirNames := []string{"/", "name*", "name;", "name:", "name?"}

	for _, name := range validDirNames {
		assert.True(t, isValidDirName(name))
	}

	for _, name := range invalidDirNames {
		assert.False(t, isValidDirName(name))
	}

	t.Log("isValidDirName: All tests passed")
}
