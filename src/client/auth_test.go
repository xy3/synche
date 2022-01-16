package client

import (
	"encoding/json"
	"github.com/xy3/synche/src/files"
	"github.com/xy3/synche/src/models"
	"os"
	"path/filepath"
	"testing"
)

var testAccessToken = T{AccessToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEsIkVtYWlsIjoiYWRtaW5AYWRtaW4uY29tIiwiTmFtZSI6ImFkbWluIiwiUGljdHVyZSI6IiIsIlJvbGUiOiJ1c2VyIiwiVG9rZW5UeXBlIjoiYWNjZXNzIiwiZXhwIjoxNjIwMDYyMzE5LCJpc3MiOiJzeW5jaGUuYXV0aC5zZXJ2aWNlIn0.GTR4xg5v7tIEb8e-leJLh9dTh6R5dWcfn9eOqSCDA7Y",
	AccessTokenExpiry: 1620062319,
	RefreshToken:      "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6ImFkbWluQGFkbWluLmNvbSIsIlRva2VuVHlwZSI6InJlZnJlc2giLCJDdXN0b21LZXkiOiI2M2FhY2EwNWI2NTcyYjkwYzI2ZDczYzU2ZjI1MzQzNCIsImlzcyI6InN5bmNoZS5hdXRoLnNlcnZpY2UifQ.FehB6zJ06QCfWA8cXXl5taGtMQHFK9JA7lo_o37Y5dU"}

type T struct {
	AccessToken       string `json:"accessToken"`
	AccessTokenExpiry int    `json:"accessTokenExpiry"`
	RefreshToken      string `json:"refreshToken"`
}

func testValidTokens(tok *models.AccessAndRefreshToken) bool {
	return tok.AccessToken == testAccessToken.AccessToken &&
		int(tok.AccessTokenExpiry) == testAccessToken.AccessTokenExpiry &&
		tok.RefreshToken == testAccessToken.RefreshToken
}

func createTestDir(dir string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	path := filepath.Join(cwd, "../data/testdata")
	return files.Afs.TempDir(path, dir)
}

func createTestJsonFile(testDir string) (string, error) {
	testFile, err := files.Afs.TempFile(testDir, "token.json")
	if err != nil {
		return "", err
	}
	fileData, _ := json.MarshalIndent(testAccessToken, "", " ")
	if err := files.Afs.WriteFile(testFile.Name(), fileData, 0644); err != nil {
		return "", err
	}

	return testFile.Name(), nil
}

func TestGetSavedToken(t *testing.T) {
	testDir, err := createTestDir("auth")
	if err != nil {
		t.Errorf("getSavedToken: Failed to create test directory: %v", err)
	}
	defer files.AppFS.Remove(testDir)

	testFile, err := createTestJsonFile(testDir)
	if err != nil {
		t.Errorf("getSavedToken: Failed to create test file: %v", err)
	}
	defer files.AppFS.Remove(testFile)

	getSavedTokenResult, err := getSavedToken(testFile)
	if err != nil {
		t.Errorf("getSavedToken failed: %v", err)
	}
	if !testValidTokens(getSavedTokenResult) {
		t.Errorf("getSavedToken: Invalid token")
	}

	_ = files.AppFS.Remove(testDir)
	_ = files.AppFS.Remove(testFile)

	t.Log("getSavedToken: All tests passed")
}
