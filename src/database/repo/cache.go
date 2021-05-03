package repo

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var (
	TokensCache           = cache.New(5*time.Minute, 10*time.Minute)
	UploadsCache          = cache.New(5*time.Minute, 10*time.Minute)
	EmailHashUserCache    = cache.New(5*time.Minute, 10*time.Minute)
	pathFileCache         = cache.New(5*time.Minute, 10*time.Minute)
	FileIDChunkCountCache = cache.New(5*time.Minute, 10*time.Minute)
)
