package jobs

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/repo"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/database/schema"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files/hash"
	"os"
	"path/filepath"
)

// CreateUniqueFilePath Creates a unique file name if the file name already exists on the server
func CreateUniqueFilePath(storageDir string, fileName string) (uniqueFilename string, uniqueFilePath string) {
	extension := filepath.Ext(fileName)
	nameWithoutExtension := fileName[0 : len(fileName)-len(extension)]
	newFilePath := filepath.Join(storageDir, fileName)
	_, err := files.Afs.Stat(newFilePath)

	var newFilename string
	for counter := 1; err == nil; counter++ {
		// Create unique filepath
		newFilename = fmt.Sprintf("%s(%d)%s", nameWithoutExtension, counter, extension)
		newFilePath = filepath.Join(storageDir, newFilename)
		_, err = files.Afs.Stat(newFilePath)
	}

	return newFilename, newFilePath
}

// ReassembleFile Retrieves all the chunk data relating to a file and reassembles the file
func ReassembleFile(chunkDir string, file *schema.File) error {
	var (
		fileChunks       []schema.FileChunk
		chunkData        []byte
		filename         string
		existingFileHash string
	)

	storageDir, err := repo.GetDirectoryForFileID(file.ID, database.DB)
	if err != nil {
		return err
	}

	filename = file.Name
	reassembledFileLocation := filepath.Join(storageDir.Path, filename)

	// Rename file if there is a file name collision
	if _, err = files.Afs.Stat(reassembledFileLocation); err == nil {
		existingFileHash, err = hash.File(reassembledFileLocation)
		if file.Hash != existingFileHash {
			filename, reassembledFileLocation = CreateUniqueFilePath(storageDir.Path, filename)
			if err = repo.RenameFile(file.ID, filename, database.DB); err != nil {
				return err
			}
		}
	}

	reassembledFile, err := files.AppFS.OpenFile(reassembledFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer reassembledFile.Close()

	fileChunks, err = repo.GetFileChunksInOrder(file.ID)
	if err != nil {
		return err
	}

	for index, chunk := range fileChunks {
		expectedNumber := int64(index + 1)
		if chunk.Number != expectedNumber {
			log.Errorf("missing chunk number: %d for file ID: %d", expectedNumber, file.ID)
			return errors.New("missing chunk for file reassembly")
		}
		// Open chunk file and get data
		chunkData, err = files.Afs.ReadFile(filepath.Join(chunkDir, chunk.Chunk.Hash))
		if err != nil {
			return err
		}

		_, err = reassembledFile.Write(chunkData)
		if err != nil {
			return err
		}
	}

	// !Important! The file must be set to available once its reassembled in order for it to appear for download
	if err = file.SetAvailable(database.DB); err != nil {
		return err
	}

	log.WithFields(log.Fields{"name": filename, "location": reassembledFileLocation}).Info("File successfully uploaded")
	return nil
}
