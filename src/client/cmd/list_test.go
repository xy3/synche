package cmd

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/client/apiclient/files"
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

func TestListDirectoryByName(t *testing.T) {
	type args struct {
		name string
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
			if got := ListDirectoryByName(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListDirectoryByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listDirectoryJob(t *testing.T) {
	type args struct {
		params *files.ListDirectoryParams
	}
	tests := []struct {
		name string
		args args
		want *files.ListDirectoryOK
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := listDirectoryJob(tt.args.params); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listDirectoryJob() = %v, want %v", got, tt.want)
			}
		})
	}
}
