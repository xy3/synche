package jobs

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/files"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type NumericalFilename []os.FileInfo

func (nFile NumericalFilename) Len() int      { return len(nFile) }
func (nFile NumericalFilename) Swap(i, j int) { nFile[i], nFile[j] = nFile[j], nFile[i] }
func (nFile NumericalFilename) Less(i, j int) bool {
	// Compare file names
	pathX := nFile[i].Name()
	pathY := nFile[j].Name()

	// Extract integer value from filename
	x, errX := strconv.ParseInt(pathX[0:strings.LastIndex(pathX, "_")], 10, 64)
	y, errY := strconv.ParseInt(pathY[0:strings.LastIndex(pathY, "_")], 10, 64)

	// Lexicographical sort in the case that no int was present
	if errX != nil || errY != nil {
		return pathX < pathY
	}

	return x < y
}

func CreateUniqueFilePath(filePath string, fileName string) (uniqueFilePath string) {
	extension := filepath.Ext(fileName)
	nameWithoutExtension := fileName[0 : len(fileName)-len(extension)]
	newFilePath := filepath.Join(filePath, fileName)
	_, err := files.Afs.Stat(newFilePath)

	for counter := 1; err == nil; counter++ {
		// Create unique filepath
		newFileName := fmt.Sprintf("%s(%d)%s ", nameWithoutExtension, counter, extension)
		newFilePath = filepath.Join(filePath, newFileName)
		_, err = files.Afs.Stat(newFilePath)
	}

	return newFilePath
}

func ReassembleFile(chunkDir, fileName string, uploadRequestId uint) error {
	chunkFileNames, err := files.Afs.ReadDir(chunkDir)
	if err != nil {
		return err
	}

	// Sort files so that they are reassembled in the correct order
	sort.Sort(NumericalFilename(chunkFileNames))

	filePath := c.Config.Server.StorageDir
	reassembledFileLocation := filepath.Join(filePath, fileName)

	// Rename file if there is a file name collision
	if _, err = os.Stat(reassembledFileLocation); err == nil {
		reassembledFileLocation = CreateUniqueFilePath(filePath, fileName)
	}

	reassembledFile, err := files.AppFS.OpenFile(reassembledFileLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer reassembledFile.Close()

	for _, file := range chunkFileNames {
		// Open chunk file and get data
		fileData, err := files.Afs.ReadFile(filepath.Join(chunkDir, file.Name()))
		if err != nil {
			return err
		}

		_, err = reassembledFile.Write(fileData)
		if err != nil {
			return err
		}
	}

	// Remove the upload from the cache
	data.Cache.Uploads.Delete(strconv.Itoa(int(uploadRequestId)))

	log.WithFields(log.Fields{"name": fileName, "location": reassembledFileLocation}).Info("File successfully uploaded")
	return nil
}
