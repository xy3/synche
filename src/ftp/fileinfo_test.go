package ftp

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func TestFileInfo_Group(t *testing.T) {
	assert.Equal(t, "synche", (&FileInfo{}).Group())
}

func TestFileInfo_IsDir(t *testing.T) {
	assert.True(t, (&FileInfo{isDir: true}).IsDir())
}

func TestFileInfo_Mode(t *testing.T) {
	assert.Equal(t, os.ModePerm|os.ModeDir, (&FileInfo{isDir: true}).Mode())
	assert.Equal(t, os.ModePerm, (&FileInfo{isDir: false}).Mode())
}

func TestFileInfo_ModTime(t *testing.T) {
	modTime := time.Now()
	assert.Equal(t, modTime, (&FileInfo{modeTime: modTime}).ModTime())
}

func TestFileInfo_Name(t *testing.T) {
	assert.Equal(t, "test", (&FileInfo{name: "test"}).Name())
}

func TestFileInfo_Owner(t *testing.T) {
	assert.Equal(t, "synche", (&FileInfo{}).Owner())
}

func TestFileInfo_Size(t *testing.T) {
	assert.Equal(t, int64(222), (&FileInfo{size: 222}).Size())
}

func TestFileInfo_Sys(t *testing.T) {
	assert.Nil(t, (&FileInfo{}).Sys())
}
