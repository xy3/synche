package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidDirName(t *testing.T) {
	validDirNames := []string{"DirName", "dir123name", "abcabc😄"}
	invalidDirNames := []string{"/", "dirname*", "dirname;", "dirname:", "dirname?"}

	for _, name := range validDirNames {
		assert.True(t, isValidDirName(name))
	}

	for _, name := range invalidDirNames {
		assert.False(t, isValidDirName(name))
	}

	t.Log("isValidDirName: All tests passed")
}
