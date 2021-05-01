package cmd

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/models"
	"reflect"
	"testing"
)

func TestListDirectoryByID(t *testing.T) {
	type args struct {
		dirId uint64
	}
	tests := []struct {
		name string
		args args
		want *models.DirectoryContents
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ListDirectoryByID(tt.args.dirId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListDirectoryByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

