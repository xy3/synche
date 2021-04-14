package data

import (
	"github.com/patrickmn/go-cache"
	"time"
)

var Cache = struct {
	Tokens  *cache.Cache
	Uploads *cache.Cache
	Users   *cache.Cache
}{
	cache.New(5*time.Minute, 10*time.Minute),
	cache.New(5*time.Minute, 10*time.Minute),
	cache.New(5*time.Minute, 10*time.Minute),
}
