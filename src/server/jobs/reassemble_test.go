package jobs_test

import (
	"github.com/stretchr/testify/require"
	"github.com/xy3/synche/src/server/jobs"
	"testing"
)

// TODO: Write tests using tempDirs
func TestCreateUniqueFilePath(t *testing.T) {
	for _, tt := range []struct {
		filePath string
		fileName string
		result   string
	}{
		{filePath: "data/", fileName: "test.mp4", result: "data/test.mp4"},
	} {
		if _, uniqueName := jobs.CreateUniqueFilePath(tt.filePath, tt.fileName); uniqueName != tt.result {
			require.Equal(t, tt.result, uniqueName)
		}
	}
}

// TODO: Test Reassemble File function
func TestReassembleFile(t *testing.T) {

}
