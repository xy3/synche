package repo

import (
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"strconv"
)

func GetDirPath(fileId int64) (dirPath string, err error) {
	var directory schema.Directory
	tx := data.DB.Begin()
	res := tx.Table("files").Select(
		"directories.path").Joins(
		"left join directories on directories.id = files.directory_id").Where(
		"files.id = ?", strconv.FormatInt(fileId, 10)).Find(
		&directory)
	if res.Error != nil {
		return "", res.Error
	}

	return  directory.Path, nil
}

