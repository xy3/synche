package data

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCRC32Hash(t *testing.T) {
	testCases := []struct {
		Name     string
		Bytes    []byte
		Expected string
	}{
		{"CRC32 string hashing", []byte("Hello world"), "8bd69e52"},
		{"CRC32 empty byte array hashing", []byte{}, "00000000"},
		{"CRC32 basic byte array hashing", []byte{0x10, 0x50, 0x00}, "59dc2736"},
		{"CRC32 character byte array hashing", []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}, "11e71165"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.Name), func(t *testing.T) {
			assert.Equal(t, tc.Expected, CRC32Hash(tc.Bytes))
		})
	}
}


func TestImoHash(t *testing.T) {
	AppFS = afero.NewMemMapFs()
	Afs = &afero.Afero{Fs: AppFS}

	testCases := []struct {
		Name     string
		Bytes    []byte
		Expected string
	}{
		{"ImoHash string hashing", []byte("Hello world"), "0b656d29d21ac0094ab5b5cd54eb2cb3"},
		{"ImoHash empty byte array hashing", []byte{}, "00000000000000000000000000000000"},
		{"ImoHash basic byte array hashing", []byte{0x10, 0x50, 0x00}, "034f88f16120c659259b18f665979be6"},
		{"ImoHash character byte array hashing", []byte{'h', 'e', 'l', 'l', '\xc3', '\xb8'}, "0671737b46eed644e6522b3cead59131"},
	}
	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%v", tc.Name), func(t *testing.T) {
			fileName := "testdata/hash/test"
			err := Afs.WriteFile(fileName, tc.Bytes, 0644)
			if err != nil {
				t.Error(err)
			}

			_, err = AppFS.Stat(fileName)
			if os.IsNotExist(err) {
				t.Errorf("file \"%s\" does not exist.\n", fileName)
			}

			actual, err := ImoHash(fileName)
			if err != nil {
				t.Errorf("Failed to hash the file: %v", err)
			}

			assert.Equal(t, tc.Expected, actual)
		})
	}
}