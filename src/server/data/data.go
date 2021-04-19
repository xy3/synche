package data

import (
	log "github.com/sirupsen/logrus"
	c "gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/config"
	"gitlab.computing.dcu.ie/collint9/2021-ca400-collint9-coynemt2/src/server/data/schema"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func MigrateAll() error {
	return DB.AutoMigrate(
		&schema.Chunk{},
		&schema.File{},
		&schema.FileChunk{},
		&schema.Directory{},
		&schema.Upload{},
		&schema.User{},
	)
}

func InitDatabase() error {
	// set up the Synche database with Redis and SQL
	db, err := gorm.Open(mysql.Open(NewDSN(c.Config.Database)), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.WithError(err).Error("Failed to open gorm DB connection")
		return err
	}
	DB = db
	return nil
}

func InitSyncheData() error {
	err := InitDatabase()
	if err != nil {
		return err
	}

	if err = MigrateAll(); err != nil {
		log.WithError(err).Error("Failed to migrate database")
		return err
	}

	if err = Configure(); err != nil {
		log.WithError(err).Error("Failed to configure the SQL database")
		return err
	}

	return nil
}

// func GetUpload(uploadId uint) (upload *schema.Upload, err error) {
// 	if res, found := Cache.Uploads.Get(strconv.Itoa(int(uploadId))); found {
// 		return res.(*schema.Upload), nil
// 	}
// 	log.WithFields(log.Fields{"UploadID": uploadId}).Debug("Upload was not retrieved from cache", uploadId)
// 	res := DataBase.First(&upload, uploadId)
// 	if res.Error != nil {
// 		return upload, res.Error
// 	}
// 	Cache.Uploads.Set(strconv.Itoa(int(uploadId)), &upload, cache.DefaultExpiration)
// 	return upload, nil
